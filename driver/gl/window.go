package gl

import (
	"bytes"
	"image"
	_ "image/png" // for the icon
	"log"
	"os"
	"runtime"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	scrollSpeed = 10
)

var (
	defaultCursor, entryCursor, hyperlinkCursor *glfw.Cursor
)

func initCursors() {
	defaultCursor = glfw.CreateStandardCursor(glfw.ArrowCursor)
	entryCursor = glfw.CreateStandardCursor(glfw.IBeamCursor)
	hyperlinkCursor = glfw.CreateStandardCursor(glfw.HandCursor)
}

type window struct {
	viewport *glfw.Window
	painted  int // part of the macOS GL fix, updated GLFW should fix this
	canvas   *glCanvas
	title    string
	icon     fyne.Resource

	clipboard fyne.Clipboard

	master     bool
	fullScreen bool
	fixedSize  bool
	padded     bool
	visible    bool

	mousePos fyne.Position
	onClosed func()

	xpos, ypos int
}

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title
	runOnMainAsync(func() {
		w.viewport.SetTitle(title)
	})
}

func (w *window) FullScreen() bool {
	return w.fullScreen
}

func (w *window) SetFullScreen(full bool) {
	w.fullScreen = full
	if !w.visible {
		return
	}
	runOnMainAsync(func() {
		monitor := w.getMonitorForWindow()
		mode := monitor.GetVideoMode()

		if full {
			w.viewport.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
		} else {
			min := w.canvas.content.MinSize()
			winWidth, winHeight := scaleInt(w.canvas, min.Width), scaleInt(w.canvas, min.Height)

			w.viewport.SetMonitor(nil, 0, 0, winWidth, winHeight, 0) // TODO remember position?
		}
	})
}

func (w *window) CenterOnScreen() {
	viewWidth, viewHeight := w.sizeOnScreen()

	runOnMainAsync(func() {
		// get window dimensions in pixels
		monitor := w.getMonitorForWindow()
		monMode := monitor.GetVideoMode()

		// these come into play when dealing with multiple monitors
		monX, monY := monitor.GetPos()

		// math them to the middle
		newX := (monMode.Width / 2) - (viewWidth / 2) + monX
		newY := (monMode.Height / 2) - (viewHeight / 2) + monY

		// set new window coordinates
		w.viewport.SetPos(newX, newY)
	})
}

// sizeOnScreen gets the size of a window content in screen pixels
func (w *window) sizeOnScreen() (int, int) {
	// get current size of content inside the window
	winContentSize := w.canvas.content.MinSize()
	// add padding, if required
	if w.Padded() {
		pad := theme.Padding() * 2
		winContentSize = fyne.NewSize(winContentSize.Width+pad, winContentSize.Height+pad)
	}

	// calculate how many pixels will be used at this scale
	viewWidth := scaleInt(w.canvas, winContentSize.Width)
	viewHeight := scaleInt(w.canvas, winContentSize.Height)

	return viewWidth, viewHeight
}

func (w *window) Resize(size fyne.Size) {
	runOnMainAsync(func() {
		scale := w.canvas.Scale()
		w.viewport.SetSize(int(float32(size.Width)*scale), int(float32(size.Height)*scale))
	})
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
	runOnMainAsync(w.fitContent)
}

func (w *window) Padded() bool {
	return w.padded
}

func (w *window) SetPadded(padded bool) {
	w.padded = padded
	if w.canvas.content == nil {
		return
	}

	if padded {
		w.canvas.content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	} else {
		w.canvas.content.Move(fyne.NewPos(0, 0))
	}

	runOnMainAsync(w.fitContent)
}

func (w *window) Icon() fyne.Resource {
	if w.icon == nil {
		return fyne.CurrentApp().Icon()
	}

	return w.icon
}

func (w *window) SetIcon(icon fyne.Resource) {
	w.icon = icon
}

func (w *window) fitContent() {
	if w.canvas.content == nil {
		return
	}

	runOnMainAsync(func() {
		winWidth, winHeight := w.sizeOnScreen()
		if w.fixedSize {
			w.viewport.SetSizeLimits(winWidth, winHeight, winWidth, winHeight)
		} else {
			w.viewport.SetSizeLimits(winWidth, winHeight, glfw.DontCare, glfw.DontCare)
		}

		width, height := w.viewport.GetSize()
		if width < winWidth || height < winHeight {
			w.viewport.SetSize(winWidth, winHeight)
		}
	})
}

func (w *window) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func scaleForDpi(xdpi int) float32 {
	if xdpi > 1000 { // assume that this is a mistake and bail
		return float32(1.0)
	}

	if xdpi > 192 {
		return float32(1.5)
	} else if xdpi > 144 {
		return float32(1.35)
	} else if xdpi > 120 {
		return float32(1.2)
	}

	return float32(1.0)
}

func (w *window) getMonitorForWindow() *glfw.Monitor {
	for _, monitor := range glfw.GetMonitors() {
		x, y := monitor.GetPos()

		if x > w.xpos || y > w.ypos {
			continue
		}
		if x+monitor.GetVideoMode().Width <= w.xpos || y+monitor.GetVideoMode().Height <= w.ypos {
			continue
		}

		return monitor
	}

	// try built-in function to detect monitor if above logic didn't succeed
	// if it doesn't work then return primary monitor as default
	monitor := w.viewport.GetMonitor()
	if monitor == nil {
		monitor = glfw.GetPrimaryMonitor()
	}
	return monitor
}

func (w *window) detectScale() float32 {
	env := os.Getenv("FYNE_SCALE")
	if env != "" {
		scale, err := strconv.ParseFloat(env, 32)
		if err != nil {
			log.Println("Error reading scale:", err)
		} else if scale != 0 {
			return float32(scale)
		}
	}

	monitor := w.getMonitorForWindow()
	widthMm, _ := monitor.GetPhysicalSize()
	widthPx := monitor.GetVideoMode().Width

	dpi := float32(widthPx) / (float32(widthMm) / 25.4)
	return scaleForDpi(int(dpi))
}

func (w *window) Show() {
	runOnMainAsync(func() {
		w.visible = true
		w.viewport.Show()

		if w.fullScreen {
			w.SetFullScreen(true)
		}
	})
}

func (w *window) Hide() {
	runOnMainAsync(func() {
		w.viewport.Hide()
		w.visible = false
	})
}

func (w *window) Close() {
	w.closed(w.viewport)
}

func (w *window) ShowAndRun() {
	w.Show()
	fyne.CurrentApp().Driver().Run()
}

//Clipboard returns the system clipboard
func (w *window) Clipboard() fyne.Clipboard {
	if w.clipboard == nil {
		w.clipboard = &clipboard{window: w.viewport}
	}
	return w.clipboard
}

func (w *window) Content() fyne.CanvasObject {
	return w.canvas.content
}

func (w *window) resize(size fyne.Size) {
	if w.Padded() {
		pad := theme.Padding() * 2
		size = fyne.NewSize(size.Width-pad, size.Height-pad)
	}

	w.canvas.content.Resize(size)
	w.canvas.setDirty(true)
}

func (w *window) SetContent(content fyne.CanvasObject) {
	w.canvas.SetContent(content)
	min := content.MinSize()

	if w.Padded() {
		pad := theme.Padding() * 2
		min = fyne.NewSize(min.Width+pad, min.Height+pad)
	}
	runOnMain(func() {
		w.fitContent()
		w.resize(min)
	})
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(true)

	// trigger callbacks
	if w.onClosed != nil {
		w.onClosed()
	}
}

func (w *window) moved(viewport *glfw.Window, x, y int) {
	// save coordinates
	w.xpos, w.ypos = x, y
	scale := w.canvas.scale
	newScale := w.detectScale()

	if scale == newScale {
		return
	}

	ratio := scale / newScale
	newWidth, newHeight := viewport.GetSize()
	newWidth = int(float32(newWidth) / ratio)
	newHeight = int(float32(newHeight) / ratio)

	w.canvas.SetScale(newScale)
	viewport.SetSize(newWidth, newHeight)
}

func (w *window) resized(viewport *glfw.Window, width, height int) {
	w.resize(fyne.NewSize(unscaleInt(w.canvas, width), unscaleInt(w.canvas, height)))
}

func (w *window) frameSized(viewport *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func (w *window) refresh(viewport *glfw.Window) {
	updateWinSize(w)
	w.canvas.setDirty(true)
}

func findMouseObj(canvas *glCanvas, mouse fyne.Position) (fyne.CanvasObject, int, int) {
	found := canvas.content
	foundX, foundY := 0, 0
	canvas.walkObjects(canvas.content, fyne.NewPos(0, 0), func(walked fyne.CanvasObject, pos fyne.Position) {
		if mouse.X < pos.X || mouse.Y < pos.Y {
			return
		}

		x2 := pos.X + walked.Size().Width
		y2 := pos.Y + walked.Size().Height
		if mouse.X >= x2 || mouse.Y >= y2 {
			return
		}

		if !walked.Visible() {
			return
		}

		switch walked.(type) {
		case fyne.Tappable, desktop.Mouseable:
			found = walked
			foundX, foundY = pos.X, pos.Y
		case fyne.Focusable:
			found = walked
			foundX, foundY = pos.X, pos.Y
		case fyne.Scrollable:
			found = walked
			foundX, foundY = pos.X, pos.Y
		}
	})

	return found, foundX, foundY
}

func (w *window) mouseMoved(viewport *glfw.Window, xpos float64, ypos float64) {
	w.mousePos = fyne.NewPos(unscaleInt(w.canvas, int(xpos)), unscaleInt(w.canvas, int(ypos)))

	co, _, _ := findMouseObj(w.canvas, w.mousePos)
	cursor := defaultCursor
	switch wid := co.(type) {
	case *widget.Entry:
		if !wid.ReadOnly {
			cursor = entryCursor
		}
	case *widget.Hyperlink:
		cursor = hyperlinkCursor
	}
	runOnMainAsync(func() {
		viewport.SetCursor(cursor)
	})
}

func (w *window) mouseClicked(viewport *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	co, x, y := findMouseObj(w.canvas, w.mousePos)
	ev := new(fyne.PointEvent)
	pad := 0
	if w.padded {
		pad = theme.Padding()
	}
	ev.Position = fyne.NewPos(w.mousePos.X-pad-x, w.mousePos.Y-pad-y)

	if wid, ok := co.(desktop.Mouseable); ok {
		mev := new(desktop.MouseEvent)
		mev.Position = ev.Position
		mev.Button = convertMouseButton(button)
		if action == glfw.Press {
			go wid.MouseDown(mev)
		} else if action == glfw.Release {
			go wid.MouseUp(mev)
		}
	}

	switch wid := co.(type) {
	case fyne.Tappable:
		if action == glfw.Press {
			switch button {
			case glfw.MouseButtonRight:
				go wid.TappedSecondary(ev)
			default:
				go wid.Tapped(ev)
			}
		}
	case fyne.Focusable:
		w.canvas.Focus(wid)
	}
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff float64, yoff float64) {
	co, _, _ := findMouseObj(w.canvas, w.mousePos)

	switch wid := co.(type) {
	case fyne.Scrollable:
		ev := &fyne.ScrollEvent{}
		ev.DeltaX = int(xoff * scrollSpeed)
		ev.DeltaY = int(yoff * scrollSpeed)
		wid.Scrolled(ev)
	}
}

func convertMouseButton(button glfw.MouseButton) desktop.MouseButton {
	switch button {
	case glfw.MouseButton1:
		return desktop.LeftMouseButton
	case glfw.MouseButton2:
		return desktop.RightMouseButton
	default:
		return 0
	}
}

func keyToName(key glfw.Key) fyne.KeyName {
	switch key {
	// non-printable
	case glfw.KeyEscape:
		return fyne.KeyEscape
	case glfw.KeyEnter:
		return fyne.KeyReturn
	case glfw.KeyTab:
		return fyne.KeyTab
	case glfw.KeyBackspace:
		return fyne.KeyBackspace
	case glfw.KeyInsert:
		return fyne.KeyInsert
	case glfw.KeyDelete:
		return fyne.KeyDelete
	case glfw.KeyRight:
		return fyne.KeyRight
	case glfw.KeyLeft:
		return fyne.KeyLeft
	case glfw.KeyDown:
		return fyne.KeyDown
	case glfw.KeyUp:
		return fyne.KeyUp
	case glfw.KeyPageUp:
		return fyne.KeyPageUp
	case glfw.KeyPageDown:
		return fyne.KeyPageDown
	case glfw.KeyHome:
		return fyne.KeyHome
	case glfw.KeyEnd:
		return fyne.KeyEnd

	case glfw.KeyF1:
		return fyne.KeyF1
	case glfw.KeyF2:
		return fyne.KeyF2
	case glfw.KeyF3:
		return fyne.KeyF3
	case glfw.KeyF4:
		return fyne.KeyF4
	case glfw.KeyF5:
		return fyne.KeyF5
	case glfw.KeyF6:
		return fyne.KeyF6
	case glfw.KeyF7:
		return fyne.KeyF7
	case glfw.KeyF8:
		return fyne.KeyF8
	case glfw.KeyF9:
		return fyne.KeyF9
	case glfw.KeyF10:
		return fyne.KeyF10
	case glfw.KeyF11:
		return fyne.KeyF11
	case glfw.KeyF12:
		return fyne.KeyF12

	case glfw.KeyKPEnter:
		return fyne.KeyEnter

	// printable
	case glfw.KeyA:
		return fyne.KeyA
	case glfw.KeyB:
		return fyne.KeyB
	case glfw.KeyC:
		return fyne.KeyC
	case glfw.KeyD:
		return fyne.KeyD
	case glfw.KeyE:
		return fyne.KeyE
	case glfw.KeyF:
		return fyne.KeyF
	case glfw.KeyG:
		return fyne.KeyG
	case glfw.KeyH:
		return fyne.KeyH
	case glfw.KeyI:
		return fyne.KeyI
	case glfw.KeyJ:
		return fyne.KeyJ
	case glfw.KeyK:
		return fyne.KeyK
	case glfw.KeyL:
		return fyne.KeyL
	case glfw.KeyM:
		return fyne.KeyM
	case glfw.KeyN:
		return fyne.KeyN
	case glfw.KeyO:
		return fyne.KeyO
	case glfw.KeyP:
		return fyne.KeyP
	case glfw.KeyQ:
		return fyne.KeyQ
	case glfw.KeyR:
		return fyne.KeyR
	case glfw.KeyS:
		return fyne.KeyS
	case glfw.KeyT:
		return fyne.KeyT
	case glfw.KeyU:
		return fyne.KeyU
	case glfw.KeyV:
		return fyne.KeyV
	case glfw.KeyW:
		return fyne.KeyW
	case glfw.KeyX:
		return fyne.KeyX
	case glfw.KeyY:
		return fyne.KeyY
	case glfw.KeyZ:
		return fyne.KeyZ
	case glfw.Key0:
		return fyne.Key0
	case glfw.Key1:
		return fyne.Key1
	case glfw.Key2:
		return fyne.Key2
	case glfw.Key3:
		return fyne.Key3
	case glfw.Key4:
		return fyne.Key4
	case glfw.Key5:
		return fyne.Key5
	case glfw.Key6:
		return fyne.Key6
	case glfw.Key7:
		return fyne.Key7
	case glfw.Key8:
		return fyne.Key8
	case glfw.Key9:
		return fyne.Key9

	// desktop
	case glfw.KeyLeftShift:
		return desktop.KeyShiftLeft
	case glfw.KeyRightShift:
		return desktop.KeyShiftRight
	case glfw.KeyLeftControl:
		return desktop.KeyControlLeft
	case glfw.KeyRightControl:
		return desktop.KeyControlRight
	case glfw.KeyLeftAlt:
		return desktop.KeyAltLeft
	case glfw.KeyRightAlt:
		return desktop.KeyAltRight
	case glfw.KeyLeftSuper:
		return desktop.KeySuperLeft
	case glfw.KeyRightSuper:
		return desktop.KeySuperRight
	case glfw.KeyMenu:
		return desktop.KeyMenu
	}
	return ""
}

func (w *window) keyPressed(viewport *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key)
	keyEvent := &fyne.KeyEvent{Name: keyName}

	if action == glfw.Press {
		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				go focused.KeyDown(keyEvent)
			}
		} else if w.canvas.onKeyDown != nil {
			go w.canvas.onKeyDown(keyEvent)
		}
	} else { // ignore key up / repeat in core events
		if action == glfw.Release {
			if w.canvas.Focused() != nil {
				if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
					go focused.KeyUp(keyEvent)
				}
			} else if w.canvas.onKeyDown != nil {
				go w.canvas.onKeyUp(keyEvent)
			}
		}
		return
	}

	keyDesktopModifier := desktopModifier(mods)
	var shortcut fyne.Shortcut
	ctrlMod := desktop.ControlModifier
	if runtime.GOOS == "darwin" {
		ctrlMod = desktop.SuperModifier
	}
	if keyDesktopModifier == ctrlMod {
		switch keyName {
		case fyne.KeyV:
			// detect paste shortcut
			shortcut = &fyne.ShortcutPaste{
				Clipboard: w.Clipboard(),
			}
		case fyne.KeyC:
			// detect copy shortcut
			shortcut = &fyne.ShortcutCopy{
				Clipboard: w.Clipboard(),
			}
		case fyne.KeyX:
			// detect cut shortcut
			shortcut = &fyne.ShortcutCut{
				Clipboard: w.Clipboard(),
			}
		default:
			shortcut = &desktop.CustomShortcut{
				KeyName:  keyName,
				Modifier: keyDesktopModifier,
			}
		}
	}

	if shortcutable, ok := w.canvas.Focused().(fyne.Shortcutable); ok {
		if shortcutable.TypedShortcut(shortcut) {
			return
		}
	} else if w.canvas.shortcut.TypedShortcut(shortcut) {
		return
	}

	// No shortcut detected, pass down to TypedKey
	if w.canvas.Focused() != nil {
		go w.canvas.Focused().TypedKey(keyEvent)
	} else if w.canvas.onTypedKey != nil {
		go w.canvas.onTypedKey(keyEvent)
	}
}

func desktopModifier(mods glfw.ModifierKey) desktop.Modifier {
	var m desktop.Modifier
	if (mods & glfw.ModShift) != 0 {
		m |= desktop.ShiftModifier
	}
	if (mods & glfw.ModControl) != 0 {
		m |= desktop.ControlModifier
	}
	if (mods & glfw.ModAlt) != 0 {
		m |= desktop.AltModifier
	}
	if (mods & glfw.ModSuper) != 0 {
		m |= desktop.SuperModifier
	}
	return m
}

// charModInput defines the character with modifiers callback which is called when a
// Unicode character is input regardless of what modifier keys are used.
//
// The character with modifiers callback is intended for implementing custom
// Unicode character input. Characters do not map 1:1 to physical keys,
// as a key may produce zero, one or more characters.
func (w *window) charModInput(viewport *glfw.Window, char rune, mods glfw.ModifierKey) {
	if mods != 0 && mods != glfw.ModShift { // don't progress if it's part of a combination
		return
	}
	if w.canvas.Focused() == nil && w.canvas.onTypedRune == nil {
		return
	}

	if w.canvas.Focused() != nil {
		w.canvas.Focused().TypedRune(char)
	} else if w.canvas.onTypedRune != nil {
		w.canvas.onTypedRune(char)
	}
}

func (d *gLDriver) CreateWindow(title string) fyne.Window {
	var ret *window
	runOnMain(func() {
		master := len(d.windows) == 0
		if master {
			glfw.Init()
			initCursors()
		}

		// make the window hidden, we will set it up and then show it later
		glfw.WindowHint(glfw.Visible, 0)

		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		win, _ := glfw.CreateWindow(10, 10, title, nil, nil)
		win.MakeContextCurrent()

		iconRes := fyne.CurrentApp().Icon()
		if iconRes != nil {
			icon, _, _ := image.Decode(bytes.NewReader(iconRes.Content()))
			win.SetIcon([]image.Image{icon})
		}

		if master {
			gl.Init()
			gl.Disable(gl.DEPTH_TEST)
		}
		ret = &window{viewport: win, title: title}
		ret.canvas = newCanvas(ret)
		ret.master = master
		ret.padded = true
		ret.canvas.SetScale(ret.detectScale())
		d.windows = append(d.windows, ret)

		win.SetCloseCallback(ret.closed)
		win.SetPosCallback(ret.moved)
		win.SetSizeCallback(ret.resized)
		win.SetFramebufferSizeCallback(ret.frameSized)
		win.SetRefreshCallback(ret.refresh)
		win.SetCursorPosCallback(ret.mouseMoved)
		win.SetMouseButtonCallback(ret.mouseClicked)
		win.SetScrollCallback(ret.mouseScrolled)
		win.SetKeyCallback(ret.keyPressed)
		win.SetCharModsCallback(ret.charModInput)
		glfw.DetachCurrentContext()
	})
	return ret
}

func (d *gLDriver) AllWindows() []fyne.Window {
	return d.windows
}
