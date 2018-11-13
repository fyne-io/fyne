// +build !ci,gl

package gl

import (
	"os"
	"strconv"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type window struct {
	viewport *glfw.Window
	canvas   *canvas
	title    string
	content  fyne.CanvasObject

	master     bool
	fullscreen bool
	fixedSize  bool

	mouseX, mouseY float64
	onClosed func()
}

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title
	w.SetTitle(title)
}

func (w *window) Fullscreen() bool {
	return w.fullscreen
}

func (w *window) SetFullscreen(full bool) {
	w.fullscreen = full
	monitor := glfw.GetPrimaryMonitor() // TODO detect if the window is on this one...
	mode := monitor.GetVideoMode()

	if full {
		w.viewport.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		min := w.content.MinSize()
		winWidth, winHeight := scaleInt(w.canvas, min.Width), scaleInt(w.canvas, min.Height)

		w.viewport.SetMonitor(nil, 0, 0, winWidth, winHeight, 0) // TODO remember position?
	}
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed

	min := w.content.MinSize()
	winWidth, winHeight := scaleInt(w.canvas, min.Width), scaleInt(w.canvas, min.Height)
	if fixed {
		w.viewport.SetSizeLimits(winWidth, winHeight, winWidth, winHeight)
	} else {
		w.viewport.SetSizeLimits(winWidth, winHeight, 10000, 10000)
	}
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

func detectScale(_ *glfw.Window) float32 {
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
	w.viewport.Show()
}

func (w *window) Hide() {
	w.viewport.Hide()
}

func (w *window) Close() {
	w.viewport.SetShouldClose(true)

	if w.onClosed != nil {
		w.onClosed()
	}
}

func (w *window) ShowAndRun() {
	w.Show()
	fyne.GetDriver().Run()
}

func (w *window) Content() fyne.CanvasObject {
	return w.content
}

func (w *window) resize(size fyne.Size) {
	w.content.Resize(size)
}

func (w *window) SetContent(content fyne.CanvasObject) {
	w.content = content
	min := content.MinSize()
	w.canvas.SetScale(detectScale(w.viewport))

	w.SetFixedSize(w.fixedSize)
	w.resize(min)
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(true)

	if w.onClosed != nil {
		w.onClosed()
	}
}

func (w *window) resized(viewport *glfw.Window, width, height int) {
	w.resize(fyne.NewSize(width, height))
}

func (w *window) mouseMoved(viewport *glfw.Window, xpos float64, ypos float64) {
	w.mouseX = xpos
	w.mouseY = ypos
}

func (w *window) mouseClicked(viewport *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	current := w.canvas
	co := w.content // TODO find correct object

	pos := fyne.NewPos(unscaleInt(current, int(w.mouseX)), unscaleInt(current, int(w.mouseY)))
	pos = pos.Subtract(fyne.NewPos(theme.Padding(), theme.Padding())) // TODO within parent

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
		w.OnMouseDown(ev)
	case fyne.FocusableObject:
		current.Focus(w)
	}
}

func keyToName(glfw.Key) string {
	return "" // TODO
}

func (w *window) keyPressed(viewport *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if w.canvas.Focused() == nil && w.canvas.onKeyDown == nil {
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
	ev.Name = string(char)
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
	master := len(d.windows) == 0
	if master {
		glfw.Init()
	}
	win, _ := glfw.CreateWindow(100, 100, title, nil, nil)
	win.MakeContextCurrent()

	if master {
		gl.Init()
	}
	ret := &window{viewport: win, title: title}
	ret.canvas = newCanvas(ret)
	ret.master = master
	d.windows = append(d.windows, ret)

	win.SetCloseCallback(ret.closed)
	win.SetSizeCallback(ret.resized)
	win.SetCursorPosCallback(ret.mouseMoved)
	win.SetMouseButtonCallback(ret.mouseClicked)
	win.SetKeyCallback(ret.keyPressed)
	win.SetCharModsCallback(ret.charModInput)
	return ret
}

func (d *gLDriver) AllWindows() []fyne.Window {
	return d.windows
}