// +build !ci

package glfw

import (
	"image/color"
	"os"
	"runtime"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/layout"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/go-gl/glfw/v3.3/glfw"
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
	w.Resize(fyne.NewSize(20, 10))

	repaintWindow(w)
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

	repaintWindow(w)
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

	repaintWindow(w)
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

	repaintWindow(w)
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
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(5, 5), AbsolutePosition: fyne.NewPos(9, 9)},
			Button:     desktop.LeftMouseButton,
		},
		d1.popMouseEvent(),
	)
	assert.Equal(t,
		&desktop.MouseEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(15, 5), AbsolutePosition: fyne.NewPos(19, 9)},
			Button:     desktop.LeftMouseButton,
		},
		d1.popMouseEvent(),
	)
	assert.Nil(t, d1.popMouseEvent())

	// we should have no mouse events on d2
	assert.Nil(t, d2.popMouseEvent())
}

func TestWindow_NoDragEndWithoutDraggedEvent(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	do := &draggableMouseableObject{Rectangle: canvas.NewRectangle(color.White)}
	do.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(do)

	repaintWindow(w)
	require.Equal(t, fyne.NewPos(4, 4), do.Position())

	w.mouseMoved(w.viewport, 9, 9)
	// mouse down (potential drag)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	// mouse release without move (not really a drag)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.waitForEvents()

	assert.Nil(t, do.popDragEvent(), "no drag event without move")
	assert.Nil(t, do.popDragEndEvent(), "no drag end event without drag event")
}

func TestWindow_HoverableOnDragging(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	dh := &draggableHoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	dh.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(dh)

	repaintWindow(w)
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

func TestWindow_Tapped(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	o := &tappableObject{Rectangle: canvas.NewRectangle(color.White)}
	o.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(widget.NewVBox(rect, o))

	w.mousePos = fyne.NewPos(50, 160)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.waitForEvents()

	assert.Nil(t, o.popSecondaryTapEvent(), "no secondary tap")
	if e, _ := o.popTapEvent().(*fyne.PointEvent); assert.NotNil(t, e, "tapped") {
		assert.Equal(t, fyne.NewPos(50, 160), e.AbsolutePosition)
		assert.Equal(t, fyne.NewPos(46, 52), e.Position)
	}
}

func TestWindow_TappedIgnoresScrollerClip(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	tapped := false
	button := widget.NewButton("Tap", func() {
		tapped = true
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

	w.waitForEvents()
	assert.False(t, tapped, "Tapped button that was clipped")

	w.mousePos = fyne.NewPos(10, 120)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.waitForEvents()
	assert.True(t, tapped, "Tapped button that was clipped")
}

func TestWindow_TappedIgnoredWhenMovedOffOfTappable(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.Canvas().SetScale(1.0)
	tapped := 0
	b1 := widget.NewButton("Tap", func() { tapped = 1 })
	b2 := widget.NewButton("Tap", func() { tapped = 2 })
	w.SetContent(widget.NewVBox(b1, b2))

	w.mouseMoved(w.viewport, 10, 20)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.waitForEvents()
	assert.Equal(t, 1, tapped, "Button 1 should be tapped")
	tapped = 0

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 10, 40)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.waitForEvents()
	assert.Equal(t, 0, tapped, "button was tapped without mouse press & release on it %d", tapped)

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.waitForEvents()
	assert.Equal(t, 2, tapped, "Button 2 should be tapped")
}

func TestWindow_MouseEventContainsModifierKeys(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	m := &mouseableObject{Rectangle: canvas.NewRectangle(color.White)}
	m.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(m)

	w.mouseMoved(w.viewport, 5, 5)
	w.waitForEvents()

	// On OS X a Ctrl+Click is normally translated into a Right-Click.
	// The well-known Ctrl+Click for extending a selection is a Cmd+Click there.
	var superModifier, ctrlModifier desktop.Modifier
	if runtime.GOOS == "darwin" {
		superModifier = desktop.ControlModifier
		ctrlModifier = 0
	} else {
		superModifier = desktop.SuperModifier
		ctrlModifier = desktop.ControlModifier
	}

	tests := map[string]struct {
		modifier              glfw.ModifierKey
		expectedEventModifier desktop.Modifier
	}{
		"no modifier key": {
			modifier:              0,
			expectedEventModifier: 0,
		},
		"shift": {
			modifier:              glfw.ModShift,
			expectedEventModifier: desktop.ShiftModifier,
		},
		"ctrl": {
			modifier:              glfw.ModControl,
			expectedEventModifier: ctrlModifier,
		},
		"alt": {
			modifier:              glfw.ModAlt,
			expectedEventModifier: desktop.AltModifier,
		},
		"super": {
			modifier:              glfw.ModSuper,
			expectedEventModifier: superModifier,
		},
		"shift+ctrl": {
			modifier:              glfw.ModShift | glfw.ModControl,
			expectedEventModifier: desktop.ShiftModifier | ctrlModifier,
		},
		"shift+alt": {
			modifier:              glfw.ModShift | glfw.ModAlt,
			expectedEventModifier: desktop.ShiftModifier | desktop.AltModifier,
		},
		"shift+super": {
			modifier:              glfw.ModShift | glfw.ModSuper,
			expectedEventModifier: desktop.ShiftModifier | superModifier,
		},
		"ctrl+alt": {
			modifier:              glfw.ModControl | glfw.ModAlt,
			expectedEventModifier: ctrlModifier | desktop.AltModifier,
		},
		"ctrl+super": {
			modifier:              glfw.ModControl | glfw.ModSuper,
			expectedEventModifier: ctrlModifier | superModifier,
		},
		"alt+super": {
			modifier:              glfw.ModAlt | glfw.ModSuper,
			expectedEventModifier: desktop.AltModifier | superModifier,
		},
		"shift+ctrl+alt": {
			modifier:              glfw.ModShift | glfw.ModControl | glfw.ModAlt,
			expectedEventModifier: desktop.ShiftModifier | ctrlModifier | desktop.AltModifier,
		},
		"shift+ctrl+super": {
			modifier:              glfw.ModShift | glfw.ModControl | glfw.ModSuper,
			expectedEventModifier: desktop.ShiftModifier | ctrlModifier | superModifier,
		},
		"shift+alt+super": {
			modifier:              glfw.ModShift | glfw.ModAlt | glfw.ModSuper,
			expectedEventModifier: desktop.ShiftModifier | desktop.AltModifier | superModifier,
		},
		"ctrl+alt+super": {
			modifier:              glfw.ModControl | glfw.ModAlt | glfw.ModSuper,
			expectedEventModifier: ctrlModifier | desktop.AltModifier | superModifier,
		},
		"shift+ctrl+alt+super": {
			modifier:              glfw.ModShift | glfw.ModControl | glfw.ModAlt | glfw.ModSuper,
			expectedEventModifier: desktop.ShiftModifier | ctrlModifier | desktop.AltModifier | superModifier,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Nil(t, m.popMouseEvent(), "no initial mouse event")
			w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, tt.modifier)
			w.waitForEvents()
			me, _ := m.popMouseEvent().(*desktop.MouseEvent)
			if assert.NotNil(t, me, "mouse event triggered") {
				assert.Equal(t, tt.expectedEventModifier, me.Modifier, "expect modifier to be correct")
			}
		})
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

			repaintWindow(w)
			contentSize := content.Size()
			expectedCanvasSize := contentSize.
				Add(fyne.NewSize(2*tt.expectedPad, 2*tt.expectedPad)).
				Add(fyne.NewSize(0, tt.expectedMenuHeight))

			w.SetPadded(tt.padding)
			repaintWindow(w)
			assert.Equal(t, contentSize, content.Size())
			assert.Equal(t, fyne.NewPos(tt.expectedPad, tt.expectedPad+tt.expectedMenuHeight), content.Position())
			assert.Equal(t, expectedCanvasSize, w.Canvas().Size())
		})
	}
}

func TestWindow_Focus(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("Test").(*window)

	e1 := widget.NewEntry()
	e2 := widget.NewEntry()

	w.SetContent(widget.NewVBox(e1, e2))
	w.Canvas().Focus(e1)

	w.charModInput(w.viewport, 'a', 0)
	w.charModInput(w.viewport, 'b', 0)
	w.charModInput(w.viewport, 'c', 0)
	w.charModInput(w.viewport, 'd', 0)
	w.keyPressed(w.viewport, glfw.KeyTab, 0, glfw.Press, 0)
	w.keyPressed(w.viewport, glfw.KeyTab, 0, glfw.Release, 0)
	w.charModInput(w.viewport, 'e', 0)
	w.charModInput(w.viewport, 'f', 0)

	w.waitForEvents()
	assert.Equal(t, "abcd", e1.Text)
	assert.Equal(t, "ef", e2.Text)
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

	// Restore clipboardContent, if any
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

type mouseableObject struct {
	*canvas.Rectangle
	mouseable
}

var _ desktop.Mouseable = (*mouseable)(nil)

type mouseable struct {
	mouseEvents []interface{}
}

func (m *mouseable) MouseDown(e *desktop.MouseEvent) {
	m.mouseEvents = append(m.mouseEvents, e)
}

func (m *mouseable) MouseUp(e *desktop.MouseEvent) {
	m.mouseEvents = append(m.mouseEvents, e)
}

func (m *mouseable) popMouseEvent() (e interface{}) {
	e, m.mouseEvents = pop(m.mouseEvents)
	return
}

type draggableMouseableObject struct {
	*canvas.Rectangle
	draggable
	mouseable
}

type tappableObject struct {
	*canvas.Rectangle
	tappable
}

var _ fyne.Tappable = (*tappable)(nil)

type tappable struct {
	tapEvents          []interface{}
	secondaryTapEvents []interface{}
}

func (t *tappable) Tapped(e *fyne.PointEvent) {
	t.tapEvents = append(t.tapEvents, e)
}

func (t *tappable) TappedSecondary(e *fyne.PointEvent) {
	t.secondaryTapEvents = append(t.secondaryTapEvents, e)
}

func (t *tappable) popTapEvent() (e interface{}) {
	e, t.tapEvents = pop(t.tapEvents)
	return
}

func (t *tappable) popSecondaryTapEvent() (e interface{}) {
	e, t.secondaryTapEvents = pop(t.secondaryTapEvents)
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
