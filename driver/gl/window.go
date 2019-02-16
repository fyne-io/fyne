package gl

import (
	"bytes"
	"image"
	_ "image/png" // for the icon
	"log"
	"os"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
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
	w.viewport.SetTitle(title)
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

	runOnMain(func() {
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
	winContentSize := w.Content().MinSize()
	// content size can be scaled, so factor that in to determining window size
	scale := w.canvas.Scale()

	// calculate how many pixels will be used at this scale
	viewWidth := int(float32(winContentSize.Width) * scale)
	viewHeight := int(float32(winContentSize.Height) * scale)

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
		min := w.canvas.content.MinSize()
		if w.Padded() {
			pad := theme.Padding() * 2
			min = fyne.NewSize(min.Width+pad, min.Height+pad)
		}
		winWidth := scaleInt(w.canvas, min.Width)
		winHeight := scaleInt(w.canvas, min.Height)
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

func (w *window) Content() fyne.CanvasObject {
	return w.canvas.content
}

func (w *window) resize(size fyne.Size) {
	if w.Padded() {
		pad := theme.Padding() * 2
		size = fyne.NewSize(size.Width-pad, size.Height-pad)
	}

	w.canvas.content.Resize(size)
	w.canvas.setDirty()
}

func (w *window) SetContent(content fyne.CanvasObject) {
	w.canvas.SetContent(content)
	min := content.MinSize()
	w.canvas.SetScale(w.detectScale())

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
	w.canvas.setDirty()
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
		case fyne.Tappable:
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
	viewport.SetCursor(cursor)
}

func (w *window) mouseClicked(viewport *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	co, x, y := findMouseObj(w.canvas, w.mousePos)
	ev := new(fyne.PointEvent)
	ev.Position = fyne.NewPos(w.mousePos.X-x, w.mousePos.Y-y)

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
		ev.DeltaX = int(xoff)
		ev.DeltaY = int(yoff)
		wid.Scrolled(ev)
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

	case glfw.KeyLeftShift:
		fallthrough
	case glfw.KeyRightShift:
		return desktop.KeyShift
	case glfw.KeyLeftControl:
		fallthrough
	case glfw.KeyRightControl:
		return desktop.KeyControl
	case glfw.KeyLeftAlt:
		fallthrough
	case glfw.KeyRightAlt:
		return desktop.KeyAlt
	case glfw.KeyLeftSuper:
		fallthrough
	case glfw.KeyRightSuper:
		return desktop.KeySuper
	case glfw.KeyMenu:
		return desktop.KeyMenu

	case glfw.KeyKPEnter:
		return fyne.KeyEnter
	}
	return ""
}

func (w *window) keyPressed(viewport *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if w.canvas.Focused() == nil && w.canvas.onTypedKey == nil {
		return
	}
	if action != glfw.Press { // ignore key up
		return
	}

	if key <= glfw.KeyWorld1 { // filter printable characters handled in charModInput
		return
	}

	ev := new(fyne.KeyEvent)
	ev.Name = keyToName(key)

	if ev.Name <= fyne.KeyF12 {
		if w.canvas.Focused() != nil {
			go w.canvas.Focused().TypedKey(ev)
		}
		if w.canvas.onTypedKey != nil {
			go w.canvas.onTypedKey(ev)
		}
	}
	// TODO handle desktop keys
}

func (w *window) charModInput(viewport *glfw.Window, char rune, mods glfw.ModifierKey) {
	if w.canvas.Focused() == nil && w.canvas.onTypedRune == nil {
		return
	}

	if mods == 0 || mods == glfw.ModShift {
		if w.canvas.Focused() != nil {
			w.canvas.Focused().TypedRune(char)
		}
		if w.canvas.onTypedRune != nil {
			w.canvas.onTypedRune(char)
		}

		return
	}

	// TODO handle shortcuts
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
