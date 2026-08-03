//go:build windows

package directx

import (
	"context"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
	fdriver "fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter/dx"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
)

const (
	dragMoveThreshold = 2
	wheelDelta        = 120

	// Scroll feel, matching the GLFW driver's non-darwin values.
	scrollSpeed            = 25
	scrollAccelerateRate   = 125
	scrollAccelerateCutoff = 10

	// Icon sizes Windows asks for, small first so it wins the title bar.
	iconSizeSmall = 16
	iconSizeBig   = 32
)

// Declare conformity with the Window interfaces.

// Declare conformity with the Window interfaces.
var (
	_ fyne.Window          = (*window)(nil)
	_ desktop.Window       = (*window)(nil)
	_ fdriver.NativeWindow = (*window)(nil)
)

type action int

const (
	press action = iota
	release
)

type window struct {
	hwnd   windows.Handle
	driver *dxDriver
	canvas *dxCanvas
	gpu    *dx.GPU
	paint  *dx.Painter

	title      string
	visible    bool
	fixedSize  bool
	fullScreen bool
	padded     bool
	decorate   bool
	master     bool
	closing    bool
	centered   bool
	// resizing is true between WM_ENTERSIZEMOVE and WM_EXITSIZEMOVE, i.e. while
	// the user is dragging a border. fitContent stands down for the duration.
	resizing bool
	// painting guards against a nested repaint: SetWindowPos delivers WM_SIZE
	// synchronously, so fitContent can re-enter the paint path from inside a paint.
	painting bool
	// highSurrogate holds the first half of a UTF-16 surrogate pair between the
	// two WM_CHAR messages that deliver a non-BMP character (emoji and friends).
	highSurrogate rune

	// customCursor is the HCURSOR built from a non-standard desktop.Cursor's
	// image, kept until the cursor changes; customCursorFor is the cursor value
	// it was built from.
	customCursor    windows.Handle
	customCursorFor desktop.Cursor

	// bgBrush paints newly exposed client area during a resize, cached against the
	// theme colour it was made from so it survives a theme change.
	bgBrush  windows.Handle
	bgColour uint32

	width, height int // requested client size in screen pixels
	restore       rect
	restoreStyle  uint32

	icon fyne.Resource
	// iconHandles are the live HICONs, indexed by ICON_SMALL/ICON_BIG, and
	// iconApplied is the resource they were made from so Show does not rebuild
	// them on every call.
	iconHandles [2]windows.Handle
	iconApplied fyne.Resource

	mainMenu     *fyne.MainMenu
	menu         *menuState
	onClosed     func()
	onCloseInter func()
	onDropped    func(fyne.Position, []fyne.URI)

	// mouse state, mirroring the GLFW driver's tap/drag/hover machine
	mousePos             fyne.Position
	mouseButton          desktop.MouseButton
	mouseDragged         fyne.Draggable
	mouseDraggedOffset   fyne.Position
	mouseDraggedObjStart fyne.Position
	mouseDragPos         fyne.Position
	mouseDragStarted     bool
	mouseOver            desktop.Hoverable
	mousePressed         fyne.CanvasObject
	mouseClickCount      int
	mouseLastClick       fyne.CanvasObject
	mouseCancelFunc      context.CancelFunc
	mouseTracking        bool
	cursor               desktop.Cursor

	lastWalkedTime time.Time
}

// frameSizeFor grows a client size into the full window size for the current style.
//
// The border is measured from the live window rather than predicted with
// AdjustWindowRect. The prediction has to be handed a DPI and a menu bar line
// count that both match what the OS actually drew, and every disagreement lands
// as a client area smaller than the one asked for - content clipped off the
// bottom, with the frame appearing to have gained a gap. A menu bar that wraps to
// two lines does it, and so does any host whose per-window DPI disagrees with the
// metrics it draws with, which is exactly what Wine does. Measuring cannot
// disagree, because it is reading back the OS's own answer.
func (w *window) frameSizeFor(width, height int) (int32, int32) {
	frame := getWindowRect(w.hwnd)
	client := getClientRect(w.hwnd)
	if client.Right > 0 && client.Bottom > 0 { // zero while minimised
		return int32(width) + (frame.Right - frame.Left) - client.Right,
			int32(height) + (frame.Bottom - frame.Top) - client.Bottom
	}

	r := rect{Right: int32(width), Bottom: int32(height)}
	adjustWindowRect(&r, w.style(), getDpiForWindow(w.hwnd), hasMenu(w.hwnd))
	return r.Right - r.Left, r.Bottom - r.Top
}

func (w *window) style() uint32 {
	// No wsVisible in any creation style: windows appear on Show, so a splash
	// window does not flash at CW_USEDEFAULT before CenterOnScreen moves it.
	if w.fullScreen || !w.decorate {
		return wsPopup
	}
	if w.fixedSize {
		return wsCaption | wsSysMenu | wsMinimizeBox
	}
	return wsOverlappedWindow
}

// resized handles WM_SIZE: it resizes the swap chain to the new pixel size and
// pushes the equivalent Fyne size into the canvas.

// resized handles WM_SIZE: it resizes the swap chain to the new pixel size and
// pushes the equivalent Fyne size into the canvas.
// syncSurface makes the swap chain and the canvas agree with the window's actual
// client area, and reports whether anything changed.
//
// The client rect is the authority rather than the size carried by WM_SIZE. Not
// every change of client area arrives as one - adding or re-laying out a native
// menu bar moves the client edge on its own, and rapid drags coalesce messages -
// and any disagreement is directly visible, because Present only covers the
// window as far as the back buffer reaches. The rest is left undrawn.
func (w *window) syncSurface() bool {
	if w.hwnd == 0 || w.gpu == nil {
		return false
	}
	cr := getClientRect(w.hwnd)
	width, height := cr.Right, cr.Bottom
	if width <= 0 || height <= 0 { // minimised
		return false
	}
	if gw, gh := w.gpu.Size(); gw == uint32(width) && gh == uint32(height) {
		return false
	}

	if err := w.gpu.Resize(uint32(width), uint32(height)); err != nil {
		fyne.LogError("directx: resizing swap chain", err)
		return false
	}
	if w.paint != nil {
		// ResizeBuffers leaves the new buffers undefined. Clear now so a frame
		// presented before the repaint lands shows the background rather than
		// uninitialised memory.
		w.paint.Clear()
	}

	w.width, w.height = int(width), int(height)
	w.canvas.Resize(fyne.NewSize(
		scale.ToFyneCoordinate(w.canvas, int(width)),
		scale.ToFyneCoordinate(w.canvas, int(height)),
	))
	w.canvas.SetDirty()
	return true
}

func (w *window) resized(width, height int32) {
	if width <= 0 || height <= 0 { // minimised
		return
	}
	w.syncSurface()

	// Draw the new size now rather than waiting for the next loop tick. Dragging a
	// border puts Windows into a modal message loop inside DefWindowProc, which
	// blocks the run loop entirely - without this the swap chain is resized but
	// never redrawn, so the window shows a stale frame for the whole drag.
	if w.visible && !w.closing && w.gpu != nil && w.paint != nil {
		w.driver.repaintWindow(w)
	}
}

func (w *window) Title() string { return w.title }

func (w *window) SetTitle(title string) {
	async.EnsureMain(func() {
		w.title = title
		if w.hwnd != 0 {
			setWindowText(w.hwnd, title)
		}
	})
}

func (w *window) FullScreen() bool { return w.fullScreen }

// SetFullScreen switches between a borderless monitor-sized window and the
// previous framed geometry, the standard Win32 approach.
func (w *window) SetFullScreen(full bool) {
	async.EnsureMain(func() {
		if w.fullScreen == full || w.hwnd == 0 {
			w.fullScreen = full
			return
		}

		if full {
			w.enterFullScreen(monitorRectFor(w.hwnd))
			return
		}

		w.fullScreen = false
		setWindowLongPtr(w.hwnd, gwlStyle, uintptr(w.restoreStyle))
		setWindowPos(w.hwnd, w.restore.Left, w.restore.Top,
			w.restore.Right-w.restore.Left, w.restore.Bottom-w.restore.Top,
			swpNoZOrder|swpFrameChanged)
	})
}

// enterFullScreen makes the window borderless and monitor-sized, remembering
// the framed geometry so SetFullScreen(false) can restore it.
func (w *window) enterFullScreen(mon rect) {
	w.fullScreen = true
	w.restore = getWindowRect(w.hwnd)
	w.restoreStyle = uint32(getWindowLongPtr(w.hwnd, gwlStyle))
	setWindowLongPtr(w.hwnd, gwlStyle, uintptr(wsPopup|wsVisible))
	setWindowPos(w.hwnd, mon.Left, mon.Top, mon.Right-mon.Left, mon.Bottom-mon.Top,
		swpNoZOrder|swpFrameChanged)
}

func (w *window) Resize(size fyne.Size) {
	async.EnsureMain(func() {
		w.canvas.Resize(size)
		width := scale.ToScreenCoordinate(w.canvas, size.Width)
		height := scale.ToScreenCoordinate(w.canvas, size.Height)
		w.width, w.height = width, height
		if w.hwnd == 0 || w.fullScreen {
			return
		}
		fw, fh := w.frameSizeFor(width, height)
		setWindowPos(w.hwnd, 0, 0, fw, fh, swpNoMove|swpNoZOrder)
	})
}

func (w *window) RequestFocus() {
	async.EnsureMain(func() {
		if w.hwnd == 0 {
			return
		}
		procSetForegroundWindow.Call(uintptr(w.hwnd))
		procSetFocus.Call(uintptr(w.hwnd))
	})
}

func (w *window) FixedSize() bool { return w.fixedSize }

func (w *window) SetFixedSize(fixed bool) {
	async.EnsureMain(func() {
		w.fixedSize = fixed
		if w.hwnd == 0 {
			return
		}
		// Keep the live visibility bit: w.style() never carries wsVisible, and
		// writing it verbatim would hide a shown window.
		style := uintptr(w.style()) | getWindowLongPtr(w.hwnd, gwlStyle)&wsVisible
		setWindowLongPtr(w.hwnd, gwlStyle, style)
		setWindowPos(w.hwnd, 0, 0, 0, 0, swpNoMove|swpNoSize|swpNoZOrder|swpFrameChanged)
	})
}

func (w *window) CenterOnScreen() {
	async.EnsureMain(func() {
		w.centered = true
		if w.hwnd == 0 {
			return
		}
		work := monitorWorkRectFor(w.hwnd)
		frame := getWindowRect(w.hwnd)
		fw := frame.Right - frame.Left
		fh := frame.Bottom - frame.Top
		x := work.Left + (work.Right-work.Left-fw)/2
		y := work.Top + (work.Bottom-work.Top-fh)/2
		setWindowPos(w.hwnd, x, y, 0, 0, swpNoSize|swpNoZOrder)
	})
}

func (w *window) Padded() bool { return w.canvas.Padded() }

func (w *window) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)
	w.canvas.SetDirty()
}

func (w *window) Icon() fyne.Resource {
	if w.icon == nil {
		return fyne.CurrentApp().Icon()
	}
	return w.icon
}

func (w *window) SetIcon(icon fyne.Resource) {
	async.EnsureMain(func() {
		w.icon = icon
		w.applyIcon()
	})
}

// applyIcon pushes the window icon to Win32, which wants a separate handle per
// size. It is also called from Show, because a window created before the app icon
// is set would otherwise keep the default Windows one.
func (w *window) applyIcon() {
	icon := w.Icon()
	if w.hwnd == 0 || icon == nil || icon == w.iconApplied {
		return
	}
	w.iconApplied = icon

	for slot, size := range [2]int{iconSmall: iconSizeSmall, iconBig: iconSizeBig} {
		h := iconFromResource(icon, size)
		if h == 0 {
			continue
		}
		// Replace first, destroy after: the shell may still paint with the old
		// handle until WM_SETICON lands, and a destroyed handle value can be
		// recycled by any other icon in the process.
		old := w.iconHandles[slot]
		w.iconHandles[slot] = h
		procSendMessageW.Call(uintptr(w.hwnd), wmSetIcon, uintptr(slot), uintptr(h))
		if old != 0 {
			procDestroyIcon.Call(uintptr(old))
		}
	}
}

// setDarkMode matches the title bar to the system theme, the way the GLFW driver
// does. The attribute is Windows 10 1809 and later; older builds keep a light
// frame, and pre-20H1 builds reject the ordinal with E_INVALIDARG.
func (w *window) setDarkMode() {
	if w.hwnd == 0 || procDwmSetWindowAttribute.Find() != nil {
		return
	}
	var dark int32
	if isDark() {
		dark = 1
	}
	procDwmSetWindowAttribute.Call(uintptr(w.hwnd), dwmwaUseImmersiveDarkMode,
		uintptr(unsafe.Pointer(&dark)), unsafe.Sizeof(dark))
}

func (w *window) SetMaster() { w.master = true }

func (w *window) MainMenu() *fyne.MainMenu { return w.mainMenu }

func (w *window) SetMainMenu(menu *fyne.MainMenu) {
	async.EnsureMain(func() {
		if menu != nil {
			addMissingQuitForMainMenu(menu, w)
		}
		w.mainMenu = menu
		w.setNativeMenu(menu)
	})
}

func (w *window) SetOnClosed(closed func()) { w.onClosed = closed }

func (w *window) SetCloseIntercept(action func()) { w.onCloseInter = action }

func (w *window) SetOnDropped(dropped func(fyne.Position, []fyne.URI)) {
	async.EnsureMain(func() {
		w.onDropped = dropped
		if w.hwnd == 0 {
			return
		}
		accept := uintptr(0)
		if dropped != nil {
			accept = 1
		}
		procDragAcceptFiles.Call(uintptr(w.hwnd), accept)
	})
}

// processDropped turns a WM_DROPFILES handle into the URIs and drop position the
// callback expects. The handle must be freed with DragFinish either way.
func (w *window) processDropped(drop uintptr) {
	defer procDragFinish.Call(drop)
	if w.onDropped == nil {
		return
	}

	count, _, _ := procDragQueryFileW.Call(drop, dragQueryCount, 0, 0)
	uris := make([]fyne.URI, 0, count)
	for i := uintptr(0); i < count; i++ {
		// A zero buffer asks for the length in characters, excluding the null.
		length, _, _ := procDragQueryFileW.Call(drop, i, 0, 0)
		if length == 0 {
			continue
		}
		buf := make([]uint16, length+1)
		if n, _, _ := procDragQueryFileW.Call(drop, i,
			uintptr(unsafe.Pointer(&buf[0])), length+1); n != 0 {
			uris = append(uris, storage.NewFileURI(windows.UTF16ToString(buf)))
		}
	}
	if len(uris) == 0 {
		return
	}

	var pt point
	procDragQueryPoint.Call(drop, uintptr(unsafe.Pointer(&pt)))
	w.onDropped(fyne.NewPos(
		scale.ToFyneCoordinate(w.canvas, int(pt.X)),
		scale.ToFyneCoordinate(w.canvas, int(pt.Y)),
	), uris)
}

func (w *window) Show() {
	async.EnsureMain(func() {
		if w.hwnd == 0 {
			return
		}
		w.visible = true
		w.applyIcon()
		showWindow(w.hwnd, swShow)
		if w.centered {
			w.CenterOnScreen()
		}
		w.canvas.reloadScale()
		w.canvas.SetDirty()
		w.RequestFocus()
	})
}

func (w *window) Hide() {
	async.EnsureMain(func() {
		if w.hwnd == 0 {
			return
		}
		w.visible = false
		showWindow(w.hwnd, swHide)
	})
}

func (w *window) Close() {
	async.EnsureMain(func() {
		if w.closing {
			return
		}
		w.closing = true
		w.visible = false
		if w.hwnd != 0 {
			showWindow(w.hwnd, swHide)
		}
	})
}

// ShowAndRun starts the application, not just the driver: App.Run is what spawns
// the lifecycle event queue goroutine and the settings watcher. Calling
// driver.Run directly leaves the event queue with no consumer, which deadlocks
// on shutdown when Run waits for it to drain.
// toggleVisible flips the window between shown and hidden, used by the system
// tray icon's tap handler.
func (w *window) toggleVisible() {
	if w.visible {
		w.Hide()
		return
	}
	w.Show()
}

func (w *window) ShowAndRun() {
	w.Show()
	fyne.CurrentApp().Run()
}

func (w *window) Content() fyne.CanvasObject { return w.canvas.Content() }

func (w *window) SetContent(obj fyne.CanvasObject) {
	async.EnsureMain(func() {
		w.canvas.SetContent(obj)
		w.fitContent()
	})
}

func (w *window) Canvas() fyne.Canvas { return w.canvas }

func (w *window) Clipboard() fyne.Clipboard { return NewClipboard() }

// applyClientSize resizes the frame so the client area comes back to the size
// that was asked for, picking up any change in border or menu bar height since.
// fitContent cannot stand in for this: it compares against the requested size,
// which has not moved, and returns early.
func (w *window) applyClientSize() {
	if w.hwnd == 0 || w.fullScreen {
		return
	}
	fw, fh := w.frameSizeFor(w.width, w.height)
	setWindowPos(w.hwnd, 0, 0, fw, fh, swpNoMove|swpNoZOrder)
}

// rescale resizes the window around a canvas whose scale has just changed. The
// canvas keeps its size in Fyne coordinates and the frame grows or shrinks to
// match, which is how the GLFW driver's RescaleContext behaves.
func (w *window) rescale() {
	if w.hwnd == 0 || w.closing || w.fullScreen {
		return
	}

	// Glyphs were rasterised for the old scale, so they would be resampled rather
	// than re-rendered and come out soft. Drop them and let the next frame rebuild.
	cache.DeleteTextTexturesFor(w.canvas)

	w.Resize(w.canvas.Size().Max(w.canvas.MinSize()))
	if content := w.canvas.Content(); content != nil {
		content.Refresh()
	}
}

// fitContent grows the window so the canvas minimum size fits.
func (w *window) fitContent() {
	// Never resize the window out from under a drag: Windows enforces the minimum
	// through WM_GETMINMAXINFO, so there is nothing to correct here mid-gesture.
	if w.hwnd == 0 || w.fullScreen || w.resizing {
		return
	}
	minSize := w.canvas.MinSize()
	minWidth := scale.ToScreenCoordinate(w.canvas, minSize.Width)
	minHeight := scale.ToScreenCoordinate(w.canvas, minSize.Height)

	width, height := w.width, w.height
	if width < minWidth {
		width = minWidth
	}
	if height < minHeight {
		height = minHeight
	}
	if width == w.width && height == w.height && !w.fixedSize {
		return
	}
	w.width, w.height = width, height
	fw, fh := w.frameSizeFor(width, height)
	setWindowPos(w.hwnd, 0, 0, fw, fh, swpNoMove|swpNoZOrder)
}

// RunNative exposes the native handle, satisfying driver.NativeWindow.

// RunNative exposes the native handle, satisfying driver.NativeWindow.
func (w *window) RunNative(f func(context any)) {
	f(fdriver.WindowsWindowContext{HWND: uintptr(w.hwnd)})
}

// eraseBackground fills the client area with the theme background.
//
// Without this the window class has no background brush, so DefWindowProc leaves
// newly exposed area untouched and it shows as black until the next frame lands.
// During a fast drag the OS composites the larger window many times before the
// repaint catches up, which is exactly when the black edges appear.
func (w *window) eraseBackground(hdc uintptr) {
	r, g, b, _ := theme.Color(theme.ColorNameBackground).RGBA()
	// COLORREF is 0x00BBGGRR.
	colour := uint32(r>>8) | uint32(g>>8)<<8 | uint32(b>>8)<<16

	if w.bgBrush == 0 || w.bgColour != colour {
		if w.bgBrush != 0 {
			deleteObject(w.bgBrush)
		}
		w.bgBrush = createSolidBrush(colour)
		w.bgColour = colour
	}
	if w.bgBrush == 0 {
		return
	}
	client := getClientRect(w.hwnd)
	fillRect(hdc, &client, w.bgBrush)
}

// recoverDevice rebuilds the Direct3D device, swap chain and painter after
// Present reports the device removed. Cached texture ids point into the dead
// painter's map, so the fresh painter simply re-uploads on the next paint.
func (w *window) recoverDevice() {
	cache.RangeTexturesFor(w.canvas, w.paint.Free)
	w.paint.Release()
	w.gpu.Release()
	w.gpu = nil
	w.paint = nil

	client := getClientRect(w.hwnd)
	g, err := dx.NewGPU(w.hwnd, uint32(client.Right), uint32(client.Bottom))
	if err != nil {
		// Leave gpu nil; the draw loop skips the window rather than crashing.
		fyne.LogError("directx: recreating Direct3D 11 device", err)
		return
	}
	w.gpu = g
	w.paint = dx.NewPainter(w.canvas, g)
	w.canvas.SetPainter(w.paint)
	w.paint.Init()
	w.canvas.SetDirty()
}

func (w *window) destroy() {
	if w.paint != nil {
		// Drop this canvas's entries from the shared texture cache before the
		// painter goes away, or they linger pointing at released GPU objects.
		cache.RangeTexturesFor(w.canvas, w.paint.Free)
		w.paint.Release()
		w.paint = nil
	}
	if w.gpu != nil {
		w.gpu.Release()
		w.gpu = nil
	}
	if w.bgBrush != 0 {
		deleteObject(w.bgBrush)
		w.bgBrush = 0
	}
	for i, h := range w.iconHandles {
		if h != 0 {
			procDestroyIcon.Call(uintptr(h))
			w.iconHandles[i] = 0
		}
	}
	if w.menu != nil {
		w.menu.release()
		w.menu = nil
	}
	w.freeCustomCursor()
	cache.CleanCanvas(w.canvas)
	if w.hwnd != 0 {
		delete(windowsByHandle, w.hwnd)
		if wakeTarget == w.hwnd {
			// Repoint the loop's wake handle at any surviving window, or clear it.
			wakeTarget = 0
			for h := range windowsByHandle {
				wakeTarget = h
				break
			}
		}
		destroyWindow(w.hwnd)
		w.hwnd = 0
	}
	if w.onClosed != nil {
		w.onClosed()
	}
}

func (w *window) processFocused(focus bool) {
	if focus {
		if curWindow == nil {
			if f := fyne.CurrentApp().Lifecycle().(*app.Lifecycle).OnEnteredForeground(); f != nil {
				f()
			}
		}
		curWindow = w
		w.canvas.FocusGained()
		return
	}

	w.canvas.FocusLost()
	w.mousePos = fyne.Position{}
	if curWindow != w {
		return
	}
	curWindow = nil
	if f := fyne.CurrentApp().Lifecycle().(*app.Lifecycle).OnExitedForeground(); f != nil {
		f()
	}
}

// RequestFullScreenSecondary fullscreens on the first non-primary monitor,
// falling back to ordinary fullscreen when only one monitor exists.
func (w *window) RequestFullScreenSecondary() {
	async.EnsureMain(func() {
		if w.hwnd == 0 || w.fullScreen {
			return
		}
		mon, ok := secondaryMonitorRect()
		if !ok {
			w.enterFullScreen(monitorRectFor(w.hwnd))
			return
		}
		w.enterFullScreen(mon)
	})
}

// RequestAlwaysOnTop keeps the window above others.
func (w *window) RequestAlwaysOnTop() {
	async.EnsureMain(func() {
		if w.hwnd == 0 {
			return
		}
		const hwndTopMost = ^uintptr(0) // (HWND)-1
		procSetWindowPos.Call(uintptr(w.hwnd), hwndTopMost, 0, 0, 0, 0, swpNoMove|swpNoSize)
	})
}

// RequestPosition moves the window to the given screen coordinates.
func (w *window) RequestPosition(x, y int) {
	async.EnsureMain(func() {
		if w.hwnd == 0 {
			return
		}
		setWindowPos(w.hwnd, int32(x), int32(y), 0, 0, swpNoSize|swpNoZOrder)
	})
}
