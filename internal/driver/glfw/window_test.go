// +build !ci

package glfw

import (
	"image/color"
	"os"
	"runtime"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"fyne.io/fyne/layout"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var d = NewGLDriver()

func init() {
	runtime.LockOSThread()
}

// TestMain makes sure that our driver is running on the main thread.
// This must be done for some of our tests to function correctly.
func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
	d.Run()
}

func TestWindow_HandleHoverable(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	h1 := &hoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	h1.SetMinSize(fyne.NewSize(10, 10))
	h2 := &hoverableObject{Rectangle: canvas.NewRectangle(color.Black)}
	h2.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(widget.NewHBox(h1, h2))

	// wait for canvas to get its size right
	for s := w.Canvas().Size(); s != fyne.NewSize(32, 18); s = w.Canvas().Size() {
		time.Sleep(time.Millisecond * 10)
	}

	require.Equal(t, fyne.NewPos(0, 0), h1.Position())
	require.Equal(t, fyne.NewPos(14, 0), h2.Position())

	w.mouseMoved(w.viewport, 9, 9)
	w.waitForEvents()
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(5, 5)}}, h1.popMouseInEvent())
	assert.Nil(t, h1.popMouseMovedEvent())
	assert.Nil(t, h1.popMouseOutEvent())

	w.mouseMoved(w.viewport, 9, 8)
	w.waitForEvents()
	assert.Nil(t, h1.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(5, 4)}}, h1.popMouseMovedEvent())
	assert.Nil(t, h1.popMouseOutEvent())

	w.mouseMoved(w.viewport, 19, 9)
	w.waitForEvents()
	assert.Nil(t, h1.popMouseInEvent())
	assert.Nil(t, h1.popMouseMovedEvent())
	assert.NotNil(t, h1.popMouseOutEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(1, 5)}}, h2.popMouseInEvent())
	assert.Nil(t, h2.popMouseMovedEvent())
	assert.Nil(t, h2.popMouseOutEvent())

	w.mouseMoved(w.viewport, 19, 8)
	w.waitForEvents()
	assert.Nil(t, h2.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(1, 4)}}, h2.popMouseMovedEvent())
	assert.Nil(t, h2.popMouseOutEvent())
}

func TestWindow_HandleDragging(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	d1 := &draggableObject{Rectangle: canvas.NewRectangle(color.White)}
	d1.SetMinSize(fyne.NewSize(10, 10))
	d2 := &draggableObject{Rectangle: canvas.NewRectangle(color.Black)}
	d2.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(widget.NewHBox(d1, d2))

	// wait for canvas to get its size right
	for s := w.Canvas().Size(); s != fyne.NewSize(32, 18); s = w.Canvas().Size() {
		time.Sleep(time.Millisecond * 10)
	}

	require.Equal(t, fyne.NewPos(0, 0), d1.Position())
	require.Equal(t, fyne.NewPos(14, 0), d2.Position())

	// no drag event in simple move
	w.mouseMoved(w.viewport, 9, 9)
	w.waitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag event on mouseDown
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.waitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// drag start and drag event with pressed mouse button
	w.mouseMoved(w.viewport, 8, 8)
	w.waitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 4)},
			DraggedX:   -1,
			DraggedY:   -1,
		},
		d1.popDragEvent(),
	)
	assert.Nil(t, d1.popDragEndEvent())
	assert.Nil(t, d2.popDragEvent())

	// drag event going outside the widget's area
	w.mouseMoved(w.viewport, 16, 8)
	w.waitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(12, 4)},
			DraggedX:   8,
			DraggedY:   0,
		},
		d1.popDragEvent(),
	)
	assert.Nil(t, d1.popDragEndEvent())
	assert.Nil(t, d2.popDragEvent())

	// drag event entering a _different_ widget's area still for the widget dragged initially
	w.mouseMoved(w.viewport, 22, 5)
	w.waitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(18, 1)},
			DraggedX:   6,
			DraggedY:   -3,
		},
		d1.popDragEvent(),
	)
	assert.Nil(t, d2.popDragEvent())

	// drag end event on mouseUp
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.waitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.NotNil(t, d1.popDragEndEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag event on further mouse move
	w.mouseMoved(w.viewport, 22, 6)
	w.waitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag event on mouseDown
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.waitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// drag event for other widget
	w.mouseMoved(w.viewport, 22, 7)
	w.waitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 3)},
			DraggedX:   0,
			DraggedY:   1,
		},
		d2.popDragEvent(),
	)
}

func TestWindow_DragObjectThatMoves(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	d1 := &draggableObject{Rectangle: canvas.NewRectangle(color.White)}
	d1.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(widget.NewHBox(d1))

	// wait for canvas to get its size right
	for s := w.Canvas().Size(); s != fyne.NewSize(18, 18); s = w.Canvas().Size() {
		time.Sleep(time.Millisecond * 10)
	}

	require.Equal(t, fyne.NewPos(0, 0), d1.Position())

	// drag -1,-1
	w.mouseMoved(w.viewport, 9, 9)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 8, 8)
	w.waitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 4)},
			DraggedX:   -1,
			DraggedY:   -1,
		},
		d1.popDragEvent(),
	)
	assert.Nil(t, d1.popDragEndEvent())

	// element follows
	d1.Move(fyne.NewPos(-1, -1))

	// drag again -> position is relative to new element position
	w.mouseMoved(w.viewport, 10, 10)
	w.waitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(7, 7)},
			DraggedX:   2,
			DraggedY:   2,
		},
		d1.popDragEvent(),
	)
}

func TestWindow_DragIntoNewObjectKeepingFocus(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	d1 := &draggableMouseableObject{Rectangle: canvas.NewRectangle(color.White)}
	d1.SetMinSize(fyne.NewSize(10, 10))
	d2 := &draggableMouseableObject{Rectangle: canvas.NewRectangle(color.White)}
	d2.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(widget.NewHBox(d1, d2))

	// wait for canvas to get its size right
	for s := w.Canvas().Size(); s != fyne.NewSize(32, 18); s = w.Canvas().Size() {
		time.Sleep(time.Millisecond * 10)
	}

	require.Equal(t, fyne.NewPos(0, 0), d1.Position())

	// drag from d1 into d2
	w.mouseMoved(w.viewport, 9, 9)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 19, 9)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.waitForEvents()

	// we should only have 2 mouse events on d1
	assert.Equal(t,
		&desktop.MouseEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(5, 5)},
			Button:     desktop.LeftMouseButton,
		},
		d1.popMouseEvent(),
	)
	assert.Equal(t,
		&desktop.MouseEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(15, 5)},
			Button:     desktop.LeftMouseButton,
		},
		d1.popMouseEvent(),
	)
	assert.Nil(t, d1.popMouseEvent())

	// we should have no mouse events on d2
	assert.Nil(t, d2.popMouseEvent())
}

func TestWindow_HoverableOnDragging(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	dh := &draggableHoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	dh.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(dh)

	// wait for canvas to get its size right
	for s := w.Canvas().Size(); s != fyne.NewSize(18, 18); s = w.Canvas().Size() {
		time.Sleep(time.Millisecond * 10)
	}

	w.mouseMoved(w.viewport, 8, 8)
	w.waitForEvents()
	assert.Equal(t,
		&desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 4)}},
		dh.popMouseInEvent(),
	)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 8, 8)
	w.waitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 4)},
			DraggedX:   0,
			DraggedY:   0,
		},
		dh.popDragEvent(),
	)

	// drag event going outside the widget's area
	w.mouseMoved(w.viewport, 16, 8)
	w.waitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(12, 4)},
			DraggedX:   8,
			DraggedY:   0,
		},
		dh.popDragEvent(),
	)
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())

	// drag event going inside the widget's area again
	w.mouseMoved(w.viewport, 8, 8)
	w.waitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 4)},
			DraggedX:   -8,
			DraggedY:   0,
		},
		dh.popDragEvent(),
	)
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())

	// no hover events on end of drag event
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.waitForEvents()
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())

	// mouseOut on mouse release after dragging out of area
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 8, 8)
	w.mouseMoved(w.viewport, 16, 8)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.waitForEvents()
	assert.NotNil(t, dh.popMouseOutEvent())
}

func TestWindow_TappedIgnoresScrollerClip(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	tapped := make(chan bool, 1)
	button := widget.NewButton("Tap", func() {
		tapped <- true
	})
	rect2 := canvas.NewRectangle(color.Black)
	rect2.SetMinSize(fyne.NewSize(100, 100))
	child := fyne.NewContainerWithLayout(layout.NewGridLayout(1), button, rect2)
	scroll := widget.NewScrollContainer(child)
	scroll.Offset = fyne.NewPos(0, 50)

	base := fyne.NewContainerWithLayout(layout.NewGridLayout(1), rect, scroll)
	w.SetContent(base)

	w.mousePos = fyne.NewPos(10, 80)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	select {
	case _ = <-tapped:
		t.Error("Tapped button that was clipped")
	case <-time.After(100 * time.Millisecond):
		// didn't tap in a decent time
	}

	w.mousePos = fyne.NewPos(10, 120)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	select {
	case _ = <-tapped:
		// button tapped
	case <-time.After(3 * time.Second):
		t.Error("waiting for button to be tapped")
	}
}

func TestWindow_TappedIgnoredWhenMovedOffOfTappable(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	tapped := make(chan int, 1)
	b1 := widget.NewButton("Tap", func() { tapped <- 1 })
	b2 := widget.NewButton("Tap", func() { tapped <- 2 })
	w.SetContent(widget.NewVBox(b1, b2))

	w.mouseMoved(w.viewport, 10, 20)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	select {
	case b := <-tapped:
		assert.Equal(t, 1, b, "button 1 should be tapped")
	case <-time.After(1 * time.Second):
		t.Error("waiting for button to be tapped")
	}

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 10, 40)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	select {
	case b := <-tapped:
		t.Error("button was tapped without mouse press & release on it:", b)
	case <-time.After(100 * time.Millisecond):
		// didn't tap in a decent time
	}

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	select {
	case b := <-tapped:
		assert.Equal(t, 2, b, "button 2 should be tapped")
	case <-time.After(1 * time.Second):
		t.Error("waiting for button to be tapped")
	}
}

func TestWindow_SetTitle(t *testing.T) {
	w := d.CreateWindow("Test")

	title := "My title"
	w.SetTitle(title)

	assert.Equal(t, title, w.Title())
}

func TestWindow_SetIcon(t *testing.T) {
	w := d.CreateWindow("Test")
	assert.Equal(t, fyne.CurrentApp().Icon(), w.Icon())

	newIcon := theme.CancelIcon()
	w.SetIcon(newIcon)
	assert.Equal(t, newIcon, w.Icon())
}

func TestWindow_PixelSize(t *testing.T) {
	w := d.CreateWindow("Test")
	w.SetPadded(false)

	rect := &canvas.Rectangle{}
	rect.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(fyne.NewContainer(rect))
	w.Canvas().Refresh(w.Content())

	winW, winH := w.(*window).minSizeOnScreen()
	assert.Equal(t, internal.ScaleInt(w.Canvas(), 100), winW)
	assert.Equal(t, internal.ScaleInt(w.Canvas(), 100), winH)
}

func TestWindow_Padded(t *testing.T) {
	w := d.CreateWindow("Test")
	content := canvas.NewRectangle(color.White)
	w.Canvas().SetScale(1.0)
	w.SetContent(content)

	width, _ := w.(*window).minSizeOnScreen()
	assert.Equal(t, theme.Padding()*2+content.MinSize().Width, width)
	assert.Equal(t, theme.Padding(), content.Position().X)
}

func TestWindow_SetPadded(t *testing.T) {
	var menuHeight int
	if hasNativeMenu() {
		menuHeight = 0
	} else {
		menuHeight = widget.NewToolbar(widget.NewToolbarAction(theme.ContentCutIcon(), func() {})).MinSize().Height
	}
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	tests := []struct {
		name               string
		padding            bool
		menu               bool
		expectedPad        int
		expectedMenuHeight int
	}{
		{"window without padding", false, false, 0, 0},
		{"window with padding", true, false, 4, 0},
		{"window with menu without padding", false, true, 0, menuHeight},
		{"window with menu and padding", true, true, 4, menuHeight},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := d.CreateWindow("Test").(*window)
			w.Canvas().SetScale(1)
			w.SetPadded(tt.padding)
			if tt.menu {
				w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Test", func() {}))))
			}
			content := canvas.NewRectangle(color.White)
			w.SetContent(content)
			oldCanvasSize := fyne.NewSize(100, 100)
			w.Resize(oldCanvasSize)

			// wait for canvas to get its size right
			for w.Canvas().Size() != oldCanvasSize {
				time.Sleep(time.Millisecond * 10)
			}
			contentSize := content.Size()
			expectedCanvasSize := contentSize.
				Add(fyne.NewSize(2*tt.expectedPad, 2*tt.expectedPad)).
				Add(fyne.NewSize(0, tt.expectedMenuHeight))

			w.SetPadded(tt.padding)
			// wait for canvas resize
			for w.Canvas().Size() != expectedCanvasSize {
				time.Sleep(time.Millisecond * 10)
			}
			assert.Equal(t, contentSize, content.Size())
			assert.Equal(t, fyne.NewPos(tt.expectedPad, tt.expectedPad+tt.expectedMenuHeight), content.Position())
			assert.Equal(t, expectedCanvasSize, w.Canvas().Size())
		})
	}
}

func TestWindow_Clipboard(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("Test")

	text := "My content from test window"
	cb := w.Clipboard()

	cliboardContent := cb.Content()
	if cliboardContent != "" {
		// Current environment has some content stored in clipboard,
		// set temporary to an empty string to allow test and restore later.
		cb.SetContent("")
	}

	assert.Empty(t, cb.Content())

	cb.SetContent(text)
	assert.Equal(t, text, cb.Content())

	// Restore cliboardContent, if any
	cb.SetContent(cliboardContent)
}

func TestWindow_Shortcut(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("Test")

	shortcutFullScreenWindow := &desktop.CustomShortcut{
		KeyName: fyne.KeyF12,
	}

	w.Canvas().AddShortcut(shortcutFullScreenWindow, func(sc fyne.Shortcut) {
		w.SetFullScreen(true)
	})

	assert.False(t, w.FullScreen())

	w.Canvas().(*glCanvas).shortcut.TypedShortcut(shortcutFullScreenWindow)
	assert.True(t, w.FullScreen())
}

//
// Test structs
//

type hoverableObject struct {
	*canvas.Rectangle
	hoverable
}

var _ desktop.Hoverable = (*hoverable)(nil)

type hoverable struct {
	mouseInEvents    []interface{}
	mouseOutEvents   []interface{}
	mouseMovedEvents []interface{}
}

func (h *hoverable) MouseIn(e *desktop.MouseEvent) {
	h.mouseInEvents = append(h.mouseInEvents, e)
}

func (h *hoverable) MouseMoved(e *desktop.MouseEvent) {
	h.mouseMovedEvents = append(h.mouseMovedEvents, e)
}

func (h *hoverable) MouseOut() {
	h.mouseOutEvents = append(h.mouseOutEvents, true)
}

func (h *hoverable) popMouseInEvent() (e interface{}) {
	e, h.mouseInEvents = pop(h.mouseInEvents)
	return
}

func (h *hoverable) popMouseMovedEvent() (e interface{}) {
	e, h.mouseMovedEvents = pop(h.mouseMovedEvents)
	return
}

func (h *hoverable) popMouseOutEvent() (e interface{}) {
	e, h.mouseOutEvents = pop(h.mouseOutEvents)
	return
}

type draggableObject struct {
	*canvas.Rectangle
	draggable
}

var _ fyne.Draggable = (*draggable)(nil)

type draggable struct {
	events    []interface{}
	endEvents []interface{}
}

func (d *draggable) Dragged(e *fyne.DragEvent) {
	d.events = append(d.events, e)
}

func (d *draggable) DragEnd() {
	d.endEvents = append(d.endEvents, true)
}

func (d *draggable) popDragEvent() (e interface{}) {
	e, d.events = pop(d.events)
	return
}

func (d *draggable) popDragEndEvent() (e interface{}) {
	e, d.endEvents = pop(d.endEvents)
	return
}

type draggableHoverableObject struct {
	*canvas.Rectangle
	draggable
	hoverable
}

var _ desktop.Mouseable = (*draggableMouseableObject)(nil)

type draggableMouseableObject struct {
	*canvas.Rectangle
	draggable
	mouseEvents []interface{}
}

func (d *draggableMouseableObject) MouseDown(e *desktop.MouseEvent) {
	d.mouseEvents = append(d.mouseEvents, e)
}

func (d *draggableMouseableObject) MouseUp(e *desktop.MouseEvent) {
	d.mouseEvents = append(d.mouseEvents, e)
}

func (d *draggableMouseableObject) popMouseEvent() (e interface{}) {
	e, d.mouseEvents = pop(d.mouseEvents)
	return
}

//
// Test helper
//

func pop(s []interface{}) (interface{}, []interface{}) {
	if len(s) == 0 {
		return nil, s
	}
	return s[0], s[1:]
}
