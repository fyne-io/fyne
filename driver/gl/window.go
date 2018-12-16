// +build !ci,gl

package gl

import (
	"bytes"
	"image"
	_ "image/png" // for the icon
	"os"
	"strconv"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type window struct {
	viewport *glfw.Window
	painted  int // part of the macOS GL fix, updated GLFW should fix this
	canvas   *glCanvas
	title    string
	icon     fyne.Resource

	master     bool
	fullScreen bool
	fixedSize  bool

	mouseX, mouseY float64
	onClosed       func()
}

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title
	w.SetTitle(title)
}

func (w *window) FullScreen() bool {
	return w.fullScreen
}

func (w *window) SetFullScreen(full bool) {
	runOnMainAsync(func() {
		w.fullScreen = full
		monitor := glfw.GetPrimaryMonitor() // TODO detect if the window is on this one...
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
		winWidth := scaleInt(w.canvas, min.Width+theme.Padding())
		winHeight := scaleInt(w.canvas, min.Height+theme.Padding())
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
	if xdpi > 192 {
		return float32(1.5)
	} else if xdpi > 144 {
		return float32(1.35)
	} else if xdpi > 120 {
		return float32(1.2)
	}

	return float32(1.0)
}

func detectScale(win *glfw.Window) float32 {
	env := os.Getenv("FYNE_SCALE")
	if env != "" {
		scale, _ := strconv.ParseFloat(env, 32)
		return float32(scale)
	}

	monitor := glfw.GetPrimaryMonitor() // TODO detect if the window is on this one...
	widthMm, _ := monitor.GetPhysicalSize()
	widthPx := monitor.GetVideoMode().Width

	dpi := float32(widthPx) / (float32(widthMm) / 25.4)
	return scaleForDpi(int(dpi))
}

func (w *window) Show() {
	runOnMainAsync(func() {
		w.viewport.Show()
	})
}

func (w *window) Hide() {
	runOnMainAsync(func() {
		w.viewport.Hide()
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
	w.canvas.content.Resize(size)
	w.canvas.setDirty()
}

func (w *window) SetContent(content fyne.CanvasObject) {
	w.canvas.SetContent(content)
	min := content.MinSize()
	w.canvas.SetScale(detectScale(w.viewport))

	runOnMainAsync(w.fitContent)
	w.resize(min)
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

func (w *window) resized(viewport *glfw.Window, width, height int) {
	w.resize(fyne.NewSize(unscaleInt(w.canvas, width)-theme.Padding()*2, unscaleInt(w.canvas, height)-theme.Padding()*2))
}

func (w *window) frameSized(viewport *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func (w *window) refresh(viewport *glfw.Window) {
	viewport.MakeContextCurrent()

	size := w.canvas.Size()
	w.canvas.paint(size)

	viewport.SwapBuffers()
	glfw.DetachCurrentContext()
}

func (w *window) mouseMoved(viewport *glfw.Window, xpos float64, ypos float64) {
	w.mouseX = xpos
	w.mouseY = ypos
}

func findMouseObj(obj fyne.CanvasObject, x, y int) fyne.CanvasObject {
	found := obj
	walkObjects(obj, fyne.NewPos(0, 0), func(walked fyne.CanvasObject, pos fyne.Position) {
		if x < pos.X || y < pos.Y {
			return
		}

		x2 := pos.X + walked.Size().Width
		y2 := pos.Y + walked.Size().Height
		if x >= x2 || y >= y2 {
			return
		}

		switch walked.(type) {
		case fyne.ClickableObject:
			found = walked
		case fyne.FocusableObject:
			found = walked
		}
	})

	return found
}

func (w *window) mouseClicked(viewport *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	current := w.canvas

	pos := fyne.NewPos(unscaleInt(current, int(w.mouseX)), unscaleInt(current, int(w.mouseY)))
	pos = pos.Subtract(fyne.NewPos(theme.Padding(), theme.Padding())) // TODO within parent
	co := findMouseObj(w.canvas.content, pos.X, pos.Y)

	ev := new(fyne.MouseEvent)
	ev.Position = pos
	switch button {
	case glfw.MouseButtonRight:
		ev.Button = fyne.RightMouseButton
		//	case glfw.MouseButtonMiddle:
		//		ev.Button = fyne.middleMouseButton
	default:
		ev.Button = fyne.LeftMouseButton
	}

	switch w := co.(type) {
	case fyne.ClickableObject:
		if action == glfw.Press {
			w.OnMouseDown(ev)
		}
	case fyne.FocusableObject:
		current.Focus(w)
	}
}

func keyToName(key glfw.Key) fyne.KeyName {
	switch key {
	// printable
	case glfw.KeySpace:
		return fyne.KeySpace

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
		return fyne.KeyShift
	case glfw.KeyLeftControl:
		fallthrough
	case glfw.KeyRightControl:
		return fyne.KeyControl
	case glfw.KeyLeftAlt:
		fallthrough
	case glfw.KeyRightAlt:
		return fyne.KeyAlt
	case glfw.KeyLeftSuper:
		fallthrough
	case glfw.KeyRightSuper:
		return fyne.KeySuper
	case glfw.KeyMenu:
		return fyne.KeyMenu

	case glfw.KeyKPEnter:
		return fyne.KeyEnter
	}
	return ""
}

func charToName(char rune) fyne.KeyName {
	switch char {
	case ' ':
		return fyne.KeySpace

	}
	return ""
}

func (w *window) keyPressed(viewport *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if w.canvas.Focused() == nil && w.canvas.onKeyDown == nil {
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
	if (mods & glfw.ModShift) != 0 {
		ev.Modifiers |= fyne.ShiftModifier
	}
	if (mods & glfw.ModControl) != 0 {
		ev.Modifiers |= fyne.ControlModifier
	}
	if (mods & glfw.ModAlt) != 0 {
		ev.Modifiers |= fyne.AltModifier
	}

	if w.canvas.Focused() != nil {
		w.canvas.Focused().OnKeyDown(ev)
	}
	if w.canvas.onKeyDown != nil {
		w.canvas.onKeyDown(ev)
	}
}

func (w *window) charModInput(viewport *glfw.Window, char rune, mods glfw.ModifierKey) {
	if w.canvas.Focused() == nil && w.canvas.onKeyDown == nil {
		return
	}

	ev := new(fyne.KeyEvent)
	ev.Name = charToName(char)
	ev.String = string(char)
	if (mods & glfw.ModShift) != 0 {
		ev.Modifiers |= fyne.ShiftModifier
	}
	if (mods & glfw.ModControl) != 0 {
		ev.Modifiers |= fyne.ControlModifier
	}
	if (mods & glfw.ModAlt) != 0 {
		ev.Modifiers |= fyne.AltModifier
	}

	if w.canvas.Focused() != nil {
		w.canvas.Focused().OnKeyDown(ev)
	}
	if w.canvas.onKeyDown != nil {
		w.canvas.onKeyDown(ev)
	}
}

func (d *gLDriver) CreateWindow(title string) fyne.Window {
	var ret *window
	runOnMain(func() {
		master := len(d.windows) == 0
		if master {
			glfw.Init()
		}

		// make the window hidden, we will set it up and then show it later
		glfw.WindowHint(glfw.Visible, 0)

		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		win, _ := glfw.CreateWindow(100, 100, title, nil, nil)
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
		d.windows = append(d.windows, ret)

		win.SetCloseCallback(ret.closed)
		win.SetSizeCallback(ret.resized)
		win.SetFramebufferSizeCallback(ret.frameSized)
		win.SetRefreshCallback(ret.refresh)
		win.SetCursorPosCallback(ret.mouseMoved)
		win.SetMouseButtonCallback(ret.mouseClicked)
		win.SetKeyCallback(ret.keyPressed)
		win.SetCharModsCallback(ret.charModInput)
		glfw.DetachCurrentContext()
	})
	return ret
}

func (d *gLDriver) AllWindows() []fyne.Window {
	return d.windows
}
