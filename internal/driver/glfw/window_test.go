//go:build !ci && !mobile
// +build !ci,!mobile

package glfw

import (
	"image/color"
	"net/url"
	"os"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

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
	d.(*gLDriver).initGLFW()
	go func() {
		// Wait for GLFW loop to be running.
		// If we try to create windows before the context is created, this will fail with an exception.
		run.Lock()
		for !run.flag {
			run.cond.Wait()
		}
		run.Unlock()

		initMainMenu()
		os.Exit(m.Run())
	}()

	master := createWindow("Master")
	master.SetOnClosed(func() {
		// we do not close, keeping the driver running
	})
	d.Run()
}

func TestGLDriver_CreateWindow(t *testing.T) {
	w := createWindow("Test").(*window)
	w.create()

	assert.Equal(t, 1, w.viewport.GetAttrib(glfw.Decorated))
	assert.True(t, w.Padded())
	assert.False(t, w.centered)
}

func TestGLDriver_CreateWindow_EmptyTitle(t *testing.T) {
	w := createWindow("").(*window)
	assert.Equal(t, w.Title(), "Fyne Application")
}

func TestGLDriver_CreateSplashWindow(t *testing.T) {
	d := NewGLDriver().(desktop.Driver)
	w := d.CreateSplashWindow().(*window)
	w.create()

	assert.Equal(t, 0, w.viewport.GetAttrib(glfw.Decorated))
	assert.False(t, w.Padded())
	assert.True(t, w.centered)
}

func TestWindow_MinSize_Fixed(t *testing.T) {
	w := createWindow("Test").(*window)
	r := canvas.NewRectangle(color.White)
	r.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(r)
	w.SetFixedSize(true)

	assert.Equal(t, float32(100)+theme.Padding()*2, w.Canvas().Size().Width)

	w = createWindow("Test").(*window)
	r.SetMinSize(fyne.NewSize(100, 100))
	w.SetFixedSize(true)
	w.SetContent(r)

	assert.Equal(t, float32(100)+theme.Padding()*2, w.Canvas().Size().Width)
}

func TestWindow_ToggleMainMenuByKeyboard(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	m := fyne.NewMainMenu(fyne.NewMenu("File"), fyne.NewMenu("Edit"), fyne.NewMenu("Help"))
	menuBar := buildMenuOverlay(m, w).(*MenuBar)
	c.Lock()
	c.setMenuOverlay(menuBar)
	c.Unlock()
	w.SetContent(canvas.NewRectangle(color.Black))

	altPressingMod := glfw.ModAlt
	altReleasingMod := glfw.ModifierKey(0)
	// Simulate known issue with GLFW inconsistency https://github.com/glfw/glfw/issues/1630
	if runtime.GOOS == "linux" {
		altPressingMod = 0
		altReleasingMod = glfw.ModAlt
	}

	require.False(t, menuBar.IsActive())
	t.Run("toggle via left Alt", func(t *testing.T) {
		w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Press, altPressingMod)
		assert.False(t, menuBar.IsActive())
		w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Release, altReleasingMod)
		assert.True(t, menuBar.IsActive())

		w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Press, altPressingMod)
		assert.True(t, menuBar.IsActive())
		w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Release, altReleasingMod)
		assert.False(t, menuBar.IsActive())
	})

	require.False(t, menuBar.IsActive())
	t.Run("toggle via right Alt", func(t *testing.T) {
		w.keyPressed(w.viewport, glfw.KeyRightAlt, 0, glfw.Press, altPressingMod)
		assert.False(t, menuBar.IsActive())
		w.keyPressed(w.viewport, glfw.KeyRightAlt, 0, glfw.Release, altReleasingMod)
		assert.True(t, menuBar.IsActive())

		w.keyPressed(w.viewport, glfw.KeyRightAlt, 0, glfw.Press, altPressingMod)
		assert.True(t, menuBar.IsActive())
		w.keyPressed(w.viewport, glfw.KeyRightAlt, 0, glfw.Release, altReleasingMod)
		assert.False(t, menuBar.IsActive())
	})

	require.False(t, menuBar.IsActive())
	t.Run("press non-special key after pressing Alt and release it before releasing Alt", func(t *testing.T) {
		w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Press, altPressingMod)
		assert.False(t, menuBar.IsActive())
		w.keyPressed(w.viewport, glfw.KeyA, 0, glfw.Press, glfw.ModAlt)
		w.keyPressed(w.viewport, glfw.KeyA, 0, glfw.Release, glfw.ModAlt)
		w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Release, altReleasingMod)
		assert.False(t, menuBar.IsActive())
	})

	for name, tt := range map[string]struct {
		key glfw.Key
		mod glfw.ModifierKey
	}{
		"left shift":    {key: glfw.KeyLeftShift, mod: glfw.ModShift},
		"right shift":   {key: glfw.KeyRightShift, mod: glfw.ModShift},
		"left control":  {key: glfw.KeyLeftControl, mod: glfw.ModControl},
		"right control": {key: glfw.KeyRightControl, mod: glfw.ModControl},
		"left super":    {key: glfw.KeyLeftSuper, mod: glfw.ModSuper},
		"right super":   {key: glfw.KeyRightSuper, mod: glfw.ModSuper},
	} {
		require.False(t, menuBar.IsActive())
		t.Run("press and release "+name+" after pressing Alt and before releasing it", func(t *testing.T) {
			w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Press, altPressingMod)
			assert.False(t, menuBar.IsActive())
			w.keyPressed(w.viewport, tt.key, 0, glfw.Press, glfw.ModAlt|tt.mod)
			w.keyPressed(w.viewport, tt.key, 0, glfw.Release, glfw.ModAlt)
			w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Release, altReleasingMod)
			assert.False(t, menuBar.IsActive())
		})

		require.False(t, menuBar.IsActive())
		t.Run("press "+name+" before pressing Alt and release it before releasing Alt", func(t *testing.T) {
			w.keyPressed(w.viewport, tt.key, 0, glfw.Press, tt.mod)
			w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Press, tt.mod|altPressingMod)
			assert.False(t, menuBar.IsActive())
			w.keyPressed(w.viewport, tt.key, 0, glfw.Release, glfw.ModAlt)
			w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Release, altReleasingMod)
			assert.False(t, menuBar.IsActive())
		})

		require.False(t, menuBar.IsActive())
		t.Run("press "+name+" after pressing Alt and release it after releasing Alt", func(t *testing.T) {
			w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Press, altPressingMod)
			assert.False(t, menuBar.IsActive())
			w.keyPressed(w.viewport, tt.key, 0, glfw.Press, glfw.ModAlt|tt.mod)
			w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Release, altReleasingMod|tt.mod)
			w.keyPressed(w.viewport, tt.key, 0, glfw.Release, 0)
			assert.False(t, menuBar.IsActive())
		})

		require.False(t, menuBar.IsActive())
		t.Run("press "+name+" before pressing Alt and release it after releasing Alt", func(t *testing.T) {
			w.keyPressed(w.viewport, tt.key, 0, glfw.Press, tt.mod)
			w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Press, tt.mod|altPressingMod)
			assert.False(t, menuBar.IsActive())
			w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Release, tt.mod|altReleasingMod)
			w.keyPressed(w.viewport, tt.key, 0, glfw.Release, 0)
			assert.False(t, menuBar.IsActive())
		})
	}

	require.False(t, menuBar.IsActive())
	t.Run("toggle via Escape", func(t *testing.T) {
		w.keyPressed(w.viewport, glfw.KeyEscape, 0, glfw.Press, 0)
		assert.False(t, menuBar.IsActive())
		w.keyPressed(w.viewport, glfw.KeyEscape, 0, glfw.Release, 0)
		assert.False(t, menuBar.IsActive(), "Escape does not activate the menu")

		c.ToggleMenu()
		require.True(t, menuBar.IsActive())

		w.keyPressed(w.viewport, glfw.KeyEscape, 0, glfw.Press, 0)
		assert.True(t, menuBar.IsActive())
		w.keyPressed(w.viewport, glfw.KeyEscape, 0, glfw.Release, 0)
		assert.False(t, menuBar.IsActive())
	})

	t.Run("when canvas has no menu", func(t *testing.T) {
		w = createWindow("Test").(*window)
		w.SetContent(canvas.NewRectangle(color.Black))

		w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Press, altPressingMod)
		w.keyPressed(w.viewport, glfw.KeyLeftAlt, 0, glfw.Release, altReleasingMod)
		// does not crash :)
	})
}

func TestWindow_Cursor(t *testing.T) {
	w := createWindow("Test").(*window)
	e := widget.NewEntry()
	u, _ := url.Parse("https://testing.fyne")
	h := widget.NewHyperlink("Testing", u)
	b := widget.NewButton("Test", nil)

	w.SetContent(container.NewVBox(e, h, b))

	w.mouseMoved(w.viewport, 10, float64(e.Position().Y+10))
	textCursor := desktop.TextCursor
	assert.Equal(t, textCursor, w.cursor)

	w.mouseMoved(w.viewport, 10, float64(h.Position().Y+10))
	pointerCursor := desktop.PointerCursor
	assert.Equal(t, pointerCursor, w.cursor)

	w.mouseMoved(w.viewport, 10, float64(b.Position().Y+10))
	defaultCursor := desktop.DefaultCursor
	assert.Equal(t, defaultCursor, w.cursor)
}

func TestWindow_HandleHoverable(t *testing.T) {
	w := createWindow("Test").(*window)
	h1 := &hoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	h1.SetMinSize(fyne.NewSize(10, 10))
	h2 := &hoverableObject{Rectangle: canvas.NewRectangle(color.Black)}
	h2.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(container.NewHBox(h1, h2))
	w.Resize(fyne.NewSize(30, 20))

	repaintWindow(w)
	require.Equal(t, fyne.NewPos(0, 0), h1.Position())
	require.Equal(t, fyne.NewPos(16, 0), h2.Position())

	w.mouseMoved(w.viewport, 9, 9)
	w.WaitForEvents()
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(3, 3),
		AbsolutePosition: fyne.NewPos(9, 9)}}, h1.popMouseInEvent())
	assert.Nil(t, h1.popMouseMovedEvent())
	assert.Nil(t, h1.popMouseOutEvent())

	w.mouseMoved(w.viewport, 9, 8)
	w.WaitForEvents()
	assert.Nil(t, h1.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(3, 2),
		AbsolutePosition: fyne.NewPos(9, 8)}}, h1.popMouseMovedEvent())
	assert.Nil(t, h1.popMouseOutEvent())

	w.mouseMoved(w.viewport, 23, 11)
	w.WaitForEvents()
	assert.Nil(t, h1.popMouseInEvent())
	assert.Nil(t, h1.popMouseMovedEvent())
	assert.NotNil(t, h1.popMouseOutEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(1, 5),
		AbsolutePosition: fyne.NewPos(23, 11)}}, h2.popMouseInEvent())
	assert.Nil(t, h2.popMouseMovedEvent())
	assert.Nil(t, h2.popMouseOutEvent())

	w.mouseMoved(w.viewport, 23, 10)
	w.WaitForEvents()
	assert.Nil(t, h2.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(1, 4),
		AbsolutePosition: fyne.NewPos(23, 10)}}, h2.popMouseMovedEvent())
	assert.Nil(t, h2.popMouseOutEvent())
}

func TestWindow_HandleOutsideHoverableObject(t *testing.T) {
	w := createWindow("Test").(*window)
	test.ApplyTheme(t, theme.DarkTheme())
	l := widget.NewList(
		func() int { return 2 },
		func() fyne.CanvasObject { return widget.NewEntry() },
		func(lii widget.ListItemID, co fyne.CanvasObject) {},
	)
	l.Resize(fyne.NewSize(200, 300))
	w.SetContent(l)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(200, 300))
	repaintWindow(w)

	w.mouseMoved(w.viewport, 7, 42)
	w.WaitForEvents()
	repaintWindow(w)
	w.mouseLock.RLock()
	assert.NotNil(t, w.mouseOver)
	w.mouseLock.RUnlock()
	test.AssertRendersToMarkup(t, "windows_hover_object.xml", w.Canvas())

	w.mouseMoved(w.viewport, 42, 42)
	w.WaitForEvents()
	repaintWindow(w)
	w.mouseLock.RLock()
	assert.NotNil(t, w.mouseOver)
	w.mouseLock.RUnlock()
	test.AssertRendersToMarkup(t, "windows_hover_object.xml", w.Canvas())

	w.mouseMoved(w.viewport, 42, 100)
	w.WaitForEvents()
	repaintWindow(w)
	w.mouseLock.RLock()
	assert.Nil(t, w.mouseOver)
	w.mouseLock.RUnlock()
	test.AssertRendersToMarkup(t, "windows_no_hover_outside_object.xml", w.Canvas())
}

func TestWindow_HandleDragging(t *testing.T) {
	w := createWindow("Test").(*window)
	d1 := &draggableObject{Rectangle: canvas.NewRectangle(color.White)}
	d1.SetMinSize(fyne.NewSize(10, 10))
	d2 := &draggableObject{Rectangle: canvas.NewRectangle(color.Black)}
	d2.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(container.NewHBox(d1, d2))

	repaintWindow(w)
	require.Equal(t, fyne.NewPos(0, 0), d1.Position())
	require.Equal(t, fyne.NewPos(16, 0), d2.Position())

	// no drag event in simple move
	w.mouseMoved(w.viewport, 9, 9)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag event on secondary mouseDown
	w.mouseClicked(w.viewport, glfw.MouseButton2, glfw.Press, 0)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag start and no drag event with pressed secondary mouse button
	w.mouseMoved(w.viewport, 8, 8)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag end event on secondary mouseUp
	w.mouseClicked(w.viewport, glfw.MouseButton2, glfw.Release, 0)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEndEvent())
	assert.Nil(t, d2.popDragEndEvent())

	// no drag event in simple move
	w.mouseMoved(w.viewport, 9, 9)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag event on secondary mouseDown
	w.mouseClicked(w.viewport, glfw.MouseButton2, glfw.Press, 0)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag start and no drag event with pressed secondary mouse button
	w.mouseMoved(w.viewport, 8, 8)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag end event on secondary mouseUp
	w.mouseClicked(w.viewport, glfw.MouseButton2, glfw.Release, 0)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEndEvent())
	assert.Nil(t, d2.popDragEndEvent())

	// no drag event in simple move
	w.mouseMoved(w.viewport, 10, 10)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag event on mouseDown
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// drag start and drag event with pressed mouse button
	w.mouseMoved(w.viewport, 8, 8)
	w.WaitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(2, 2),
				AbsolutePosition: fyne.NewPos(8, 8)},
			Dragged: fyne.NewDelta(-2, -2),
		},
		d1.popDragEvent(),
	)
	assert.Nil(t, d1.popDragEndEvent())
	assert.Nil(t, d2.popDragEvent())

	// drag event going outside the widget's area
	w.mouseMoved(w.viewport, 16, 8)
	w.WaitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(10, 2),
				AbsolutePosition: fyne.NewPos(16, 8)},
			Dragged: fyne.NewDelta(8, 0),
		},
		d1.popDragEvent(),
	)
	assert.Nil(t, d1.popDragEndEvent())
	assert.Nil(t, d2.popDragEvent())

	// drag event entering a _different_ widget's area still for the widget dragged initially
	w.mouseMoved(w.viewport, 22, 6)
	w.WaitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(16, 0),
				AbsolutePosition: fyne.NewPos(22, 6)},
			Dragged: fyne.NewDelta(6, -2),
		},
		d1.popDragEvent(),
	)
	assert.Nil(t, d2.popDragEvent())

	// drag end event on mouseUp
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.NotNil(t, d1.popDragEndEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag event on further mouse move
	w.mouseMoved(w.viewport, 22, 6)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// no drag event on mouseDown
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Nil(t, d2.popDragEvent())

	// drag event for other widget
	w.mouseMoved(w.viewport, 26, 9)
	w.WaitForEvents()
	assert.Nil(t, d1.popDragEvent())
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 3),
				AbsolutePosition: fyne.NewPos(26, 9)},
			Dragged: fyne.NewDelta(4, 3),
		},
		d2.popDragEvent(),
	)
}

func TestWindow_DragObjectThatMoves(t *testing.T) {
	w := createWindow("Test").(*window)
	d1 := &draggableObject{Rectangle: canvas.NewRectangle(color.White)}
	d1.SetMinSize(fyne.NewSize(20, 20))
	w.SetContent(container.NewHBox(d1))

	repaintWindow(w)
	require.Equal(t, fyne.NewPos(0, 0), d1.Position())

	// drag -1,-1
	w.mouseMoved(w.viewport, 12, 12)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 10, 10)
	w.WaitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 4),
				AbsolutePosition: fyne.NewPos(10, 10)},
			Dragged: fyne.NewDelta(-2, -2),
		},
		d1.popDragEvent(),
	)
	assert.Nil(t, d1.popDragEndEvent())

	// element follows
	d1.Move(fyne.NewPos(-1, -1))

	// drag again -> position is relative to new element position
	w.mouseMoved(w.viewport, 12, 12)
	w.WaitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(7, 7),
				AbsolutePosition: fyne.NewPos(12, 12)},
			Dragged: fyne.NewDelta(2, 2),
		},
		d1.popDragEvent(),
	)
}

func TestWindow_DragIntoNewObjectKeepingFocus(t *testing.T) {
	w := createWindow("Test").(*window)
	d1 := &draggableMouseableObject{Rectangle: canvas.NewRectangle(color.White)}
	d1.SetMinSize(fyne.NewSize(20, 20))
	d2 := &draggableMouseableObject{Rectangle: canvas.NewRectangle(color.White)}
	d2.SetMinSize(fyne.NewSize(20, 20))
	w.SetContent(container.NewHBox(d1, d2))

	repaintWindow(w)
	require.Equal(t, fyne.NewPos(0, 0), d1.Position())

	// drag from d1 into d2
	w.mouseMoved(w.viewport, 11, 11)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 21, 11)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()

	// we should only have 2 mouse events on d1
	assert.Equal(t,
		&desktop.MouseEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(5, 5), AbsolutePosition: fyne.NewPos(11, 11)},
			Button:     desktop.MouseButtonPrimary,
		},
		d1.popMouseEvent(),
	)
	assert.Equal(t,
		&desktop.MouseEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(15, 5), AbsolutePosition: fyne.NewPos(21, 11)},
			Button:     desktop.MouseButtonPrimary,
		},
		d1.popMouseEvent(),
	)
	assert.Nil(t, d1.popMouseEvent())

	// we should have no mouse events on d2
	assert.Nil(t, d2.popMouseEvent())
}

func TestWindow_NoDragEndWithoutDraggedEvent(t *testing.T) {
	w := createWindow("Test").(*window)
	do := &draggableMouseableObject{Rectangle: canvas.NewRectangle(color.White)}
	do.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(do)

	repaintWindow(w)
	require.Equal(t, fyne.NewPos(6, 6), do.Position())

	w.mouseMoved(w.viewport, 9, 9)
	// mouse down (potential drag)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	// mouse release without move (not really a drag)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()

	assert.Nil(t, do.popDragEvent(), "no drag event without move")
	assert.Nil(t, do.popDragEndEvent(), "no drag end event without drag event")
}

func TestWindow_HoverableOnDragging(t *testing.T) {
	w := createWindow("Test").(*window)
	dh := &draggableHoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	c := container.NewWithoutLayout(dh)
	dh.Resize(fyne.NewSize(20, 20))
	w.SetContent(c)

	repaintWindow(w)
	w.mouseMoved(w.viewport, 10, 10)
	w.WaitForEvents()
	assert.Equal(t,
		&desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 4),
			AbsolutePosition: fyne.NewPos(10, 10)}},
		dh.popMouseInEvent(),
	)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 12, 12)
	w.WaitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(6, 6),
				AbsolutePosition: fyne.NewPos(12, 12)},
			Dragged: fyne.NewDelta(2, 2),
		},
		dh.popDragEvent(),
	)

	// drag event going outside the widget's area
	w.mouseMoved(w.viewport, 20, 12)
	w.WaitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(14, 6),
				AbsolutePosition: fyne.NewPos(20, 12)},
			Dragged: fyne.NewDelta(8, 0),
		},
		dh.popDragEvent(),
	)
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())

	// drag event going inside the widget's area again
	w.mouseMoved(w.viewport, 12, 12)
	w.WaitForEvents()
	assert.Equal(t,
		&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(6, 6),
				AbsolutePosition: fyne.NewPos(12, 12)},
			Dragged: fyne.NewDelta(-8, 0),
		},
		dh.popDragEvent(),
	)
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())

	// no hover events on end of drag event
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())

	// mouseOut on mouse release after dragging out of area
	w.mouseMoved(w.viewport, 12, 12)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 28, 12) // outside the 20x20 object
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()
	assert.NotNil(t, dh.popMouseOutEvent())
}

func TestWindow_HoverableUnderDraggable(t *testing.T) {
	w := createWindow("Test").(*window)
	h := &hoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	d := &draggableObject{Rectangle: canvas.NewRectangle(color.White)}
	dh := &draggableHoverableObject{Rectangle: canvas.NewRectangle(color.White)}

	c := container.NewWithoutLayout(h, d, dh)
	h.Resize(fyne.NewSize(50, 50))
	h.Move(fyne.NewPos(0, 0))
	d.Resize(fyne.NewSize(30, 30))
	d.Move(fyne.NewPos(10, 10))
	dh.Resize(fyne.NewSize(10, 10))
	dh.Move(fyne.NewPos(20, 20))

	w.SetContent(c)

	repaintWindow(w)

	// 1. move over to hoverableObject and verify
	//  - mouseIn received by hoverableObject
	//  - no events by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseMoved(w.viewport, 7, 7)
	w.WaitForEvents()
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(1, 1),
		AbsolutePosition: fyne.NewPos(7, 7)}}, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 1.a test for drag in non-draggable
	//  - move in hoverableObject
	//  - no events by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 8, 8)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(2, 2),
		AbsolutePosition: fyne.NewPos(8, 8)}, Button: 1, Modifier: 0}, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	// 2. move over to draggableObject and verify
	//  - mouseMoved by hoverableObject
	//  - no events by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseMoved(w.viewport, 16, 16)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(10, 10),
		AbsolutePosition: fyne.NewPos(16, 16)}}, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 2.a test for drag in draggable
	//  - move in hoverableObject
	//  - drag begin by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 18, 18)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(12, 12),
		AbsolutePosition: fyne.NewPos(18, 18)}, Button: 1, Modifier: 0}, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Equal(t, &fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(2, 2),
		AbsolutePosition: fyne.NewPos(18, 18)}, Dragged: fyne.NewDelta(2, 2)}, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 2.b drag end
	//  - no events by hoverableObject
	//  - drag end by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.NotNil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 3. move to draggableHoverableObject and verify
	// - mouseOut received by hoverableObject
	// - no events by draggableObject
	// - mouseIn by draggableHoverableObject
	w.mouseMoved(w.viewport, 27, 27)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.NotNil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(1, 1),
		AbsolutePosition: fyne.NewPos(27, 27)}}, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 4. move over to draggableObject and verify
	// - mouseIn received by hoverableObject
	// - no events by draggableObject
	// - mouseOut by draggableHoverableObject
	w.mouseMoved(w.viewport, 37, 37)
	w.WaitForEvents()
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(31, 31),
		AbsolutePosition: fyne.NewPos(37, 37)}}, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.NotNil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 5. move to  hoverableObject and verify
	//  - mouseMoved by hoverableObject
	//  - no events by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseMoved(w.viewport, 47, 47)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(41, 41),
		AbsolutePosition: fyne.NewPos(47, 47)}}, h.popMouseMovedEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 6. move ouside hoverableObject and verify
	//  - mouseOut by hoverableObject
	//  - no events by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseMoved(w.viewport, 57, 57)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.NotNil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())
}

func TestWindow_HoverableUnderDraggable_DragAcross(t *testing.T) {
	w := createWindow("Test").(*window)
	h := &hoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	d := &draggableObject{Rectangle: canvas.NewRectangle(color.White)}
	dh := &draggableHoverableObject{Rectangle: canvas.NewRectangle(color.White)}

	c := container.NewWithoutLayout(h, d, dh)
	h.Resize(fyne.NewSize(50, 50))
	h.Move(fyne.NewPos(0, 0))
	d.Resize(fyne.NewSize(30, 30))
	d.Move(fyne.NewPos(10, 10))
	dh.Resize(fyne.NewSize(10, 10))
	dh.Move(fyne.NewPos(20, 20))

	w.SetContent(c)

	repaintWindow(w)

	// 1. drag across hoverable
	//  - mouseIn by hoverableObject
	//  - no events by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseMoved(w.viewport, 16, 16)
	w.WaitForEvents()
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(10, 10),
		AbsolutePosition: fyne.NewPos(16, 16)}}, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 2 start drag in draggable
	//  - move in hoverableObject
	//  - drag begin by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 18, 18)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(12, 12),
		AbsolutePosition: fyne.NewPos(18, 18)}, Button: 1, Modifier: 0}, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Equal(t, &fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(2, 2),
		AbsolutePosition: fyne.NewPos(18, 18)}, Dragged: fyne.NewDelta(2, 2)}, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 3 drag to draggable+hoverable
	//  - moveOut hoverableObject
	//  - drag events by draggableObject
	//  - moveIn by draggableHoverableObject
	w.mouseMoved(w.viewport, 27, 27)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.NotNil(t, h.popMouseOutEvent())
	assert.Equal(t, &fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(11, 11),
		AbsolutePosition: fyne.NewPos(27, 27)}, Dragged: fyne.NewDelta(9, 9)}, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(1, 1),
		AbsolutePosition: fyne.NewPos(27, 27)}, Button: 1, Modifier: 0}, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 4 move over to draggableObject and verify
	// - mouseIn received by hoverableObject
	// - drag event by draggableObject
	// - mouseOut by draggableHoverableObject
	w.mouseMoved(w.viewport, 37, 37)
	w.WaitForEvents()
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(31, 31),
		AbsolutePosition: fyne.NewPos(37, 37)}, Button: 1, Modifier: 0}, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Equal(t, &fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(21, 21),
		AbsolutePosition: fyne.NewPos(37, 37)}, Dragged: fyne.NewDelta(10, 10)}, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.NotNil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 5 drag end
	//  - no events by hoverableObject
	//  - drag end by draggableObject
	//  - no events by draggableHoverableObject
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.NotNil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())
}

func TestWindow_HoverableUnderDraggable_Drag_draggableHoverable(t *testing.T) {
	w := createWindow("Test").(*window)
	h := &hoverableObject{Rectangle: canvas.NewRectangle(color.White)}
	d := &draggableObject{Rectangle: canvas.NewRectangle(color.White)}
	dh := &draggableHoverableObject{Rectangle: canvas.NewRectangle(color.White)}

	c := container.NewWithoutLayout(h, d, dh)
	h.Resize(fyne.NewSize(50, 50))
	h.Move(fyne.NewPos(0, 0))
	d.Resize(fyne.NewSize(30, 30))
	d.Move(fyne.NewPos(10, 10))
	dh.Resize(fyne.NewSize(10, 10))
	dh.Move(fyne.NewPos(20, 20))

	w.SetContent(c)

	repaintWindow(w)

	// 1. drag of draggableHoverable
	//  - no event by hoverableObject
	//  - no event by draggableObject
	//  - moveIn event by draggableHoverableObject
	w.mouseMoved(w.viewport, 28, 28)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(2, 2),
		AbsolutePosition: fyne.NewPos(28, 28)}}, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 2. start drag in draggableHoverable
	//  - no events by hoverableObject
	//  - no events by draggableObject
	//  - drag begin by draggableHoverableObject
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 30, 30)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Equal(t, &fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(4, 4),
		AbsolutePosition: fyne.NewPos(30, 30)}, Dragged: fyne.NewDelta(2, 2)}, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 3. drag to hoverable
	//  - moveIn by hoverableObject
	//  - no events by draggableObject
	//  - drag and moveOut by draggableHoverableObject
	w.mouseMoved(w.viewport, 47, 47)
	w.WaitForEvents()
	assert.Equal(t, &desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(41, 41),
		AbsolutePosition: fyne.NewPos(47, 47)}, Button: 1, Modifier: 0}, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.NotNil(t, dh.popMouseOutEvent())
	assert.Equal(t, &fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(21, 21),
		AbsolutePosition: fyne.NewPos(47, 47)}, Dragged: fyne.NewDelta(17, 17)}, dh.popDragEvent())
	assert.Nil(t, dh.popDragEndEvent())

	// 4. drag end
	//  - no events by hoverableObject
	//  - no events by draggableObject
	//  - drag end by draggableHoverableObject
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()
	assert.Nil(t, h.popMouseInEvent())
	assert.Nil(t, h.popMouseMovedEvent())
	assert.Nil(t, h.popMouseOutEvent())
	assert.Nil(t, d.popDragEvent())
	assert.Nil(t, d.popDragEndEvent())
	assert.Nil(t, dh.popMouseInEvent())
	assert.Nil(t, dh.popMouseMovedEvent())
	assert.Nil(t, dh.popMouseOutEvent())
	assert.Nil(t, dh.popDragEvent())
	assert.NotNil(t, dh.popDragEndEvent())
}

func TestWindow_DragEndWithoutTappedEvent(t *testing.T) {
	w := createWindow("Test").(*window)
	do := &draggableTappableObject{Rectangle: canvas.NewRectangle(color.White)}
	do.SetMinSize(fyne.NewSize(14, 14))
	w.SetContent(do)

	repaintWindow(w)
	require.Equal(t, fyne.NewPos(6, 6), do.Position())

	w.mouseMoved(w.viewport, 11, 11)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 10, 10) // Less than drag threshold
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.WaitForEvents()
	assert.NotNil(t, do.popTapEvent()) // it was slight drag, so call it a tap

	w.mouseMoved(w.viewport, 7, 7)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.WaitForEvents()

	assert.Nil(t, do.popTapEvent())
}

func TestWindow_Scrolled(t *testing.T) {
	w := createWindow("Test").(*window)
	o := &scrollable{Rectangle: canvas.NewRectangle(color.White)}
	o.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(o)

	w.mousePos = fyne.NewPos(50, 60)
	w.mouseScrolled(w.viewport, 10, 10)

	if e, _ := o.popScrollEvent().(*fyne.ScrollEvent); assert.NotNil(t, e, "scroll event") {
		assert.Equal(t, fyne.NewPos(50, 60), e.AbsolutePosition)
		assert.Equal(t, fyne.NewPos(44, 54), e.Position)
	}
}

func TestWindow_Tapped(t *testing.T) {
	w := createWindow("Test").(*window)
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	o := &tappableObject{Rectangle: canvas.NewRectangle(color.White)}
	o.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(container.NewVBox(rect, o))

	w.mousePos = fyne.NewPos(50, 160)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()

	assert.Nil(t, o.popSecondaryTapEvent(), "no secondary tap")
	if e, _ := o.popTapEvent().(*fyne.PointEvent); assert.NotNil(t, e, "tapped") {
		assert.Equal(t, fyne.NewPos(50, 160), e.AbsolutePosition)
		assert.Equal(t, fyne.NewPos(44, 48), e.Position)
	}
}

func TestWindow_TappedSecondary(t *testing.T) {
	w := createWindow("Test").(*window)
	o := &tappableObject{Rectangle: canvas.NewRectangle(color.White)}
	o.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(o)

	w.mousePos = fyne.NewPos(50, 60)
	w.mouseClicked(w.viewport, glfw.MouseButton2, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton2, glfw.Release, 0)
	w.WaitForEvents()

	assert.Nil(t, o.popTapEvent(), "no primary tap")
	if e, _ := o.popSecondaryTapEvent().(*fyne.PointEvent); assert.NotNil(t, e, "tapped secondary") {
		assert.Equal(t, fyne.NewPos(50, 60), e.AbsolutePosition)
		assert.Equal(t, fyne.NewPos(44, 54), e.Position)
	}
}

func TestWindow_TappedSecondary_OnPrimaryOnlyTarget(t *testing.T) {
	w := createWindow("Test").(*window)
	tapped := false
	o := widget.NewButton("Test", func() {
		tapped = true
	})
	w.SetContent(o)

	w.mousePos = fyne.NewPos(10, 25)
	w.mouseClicked(w.viewport, glfw.MouseButton2, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton2, glfw.Release, 0)
	w.WaitForEvents()

	assert.False(t, tapped)

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()

	assert.True(t, tapped)
}

func TestWindow_TappedIgnoresScrollerClip(t *testing.T) {
	w := createWindow("Test").(*window)
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	tapped := false
	button := widget.NewButton("Tap", func() {
		tapped = true
	})
	rect2 := canvas.NewRectangle(color.Black)
	rect2.SetMinSize(fyne.NewSize(100, 100))
	child := container.NewGridWithColumns(1, button, rect2)
	scroll := container.NewScroll(child)
	scroll.Offset = fyne.NewPos(0, 50)

	base := container.New(layout.NewGridLayout(1), rect, scroll)
	w.SetContent(base)
	refreshWindow(w) // ensure any async resize is done

	w.mousePos = fyne.NewPos(10, 80)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.WaitForEvents()

	assert.False(t, tapped, "Tapped button that was clipped")

	w.mousePos = fyne.NewPos(10, 120)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.WaitForEvents()

	assert.True(t, tapped, "Tapped button that was clipped")
}

func TestWindow_TappedIgnoredWhenMovedOffOfTappable(t *testing.T) {
	w := createWindow("Test").(*window)
	tapped := 0
	b1 := widget.NewButton("Tap", func() { tapped = 1 })
	b2 := widget.NewButton("Tap", func() { tapped = 2 })
	w.SetContent(container.NewVBox(b1, b2))

	w.mouseMoved(w.viewport, 17, 27)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.WaitForEvents()

	assert.Equal(t, 1, tapped, "Button 1 should be tapped")
	tapped = 0

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseMoved(w.viewport, 17, 59)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.WaitForEvents()

	assert.Equal(t, 0, tapped, "button was tapped without mouse press & release on it %d", tapped)

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	w.WaitForEvents()

	assert.Equal(t, 2, tapped, "Button 2 should be tapped")
}

func TestWindow_TappedAndDoubleTapped(t *testing.T) {
	w := createWindow("Test").(*window)
	waitSingleTapped := make(chan struct{})
	waitDoubleTapped := make(chan struct{})
	tapped := int32(0) // atomic
	but := newDoubleTappableButton()
	but.OnTapped = func() {
		atomic.StoreInt32(&tapped, 1)
		waitSingleTapped <- struct{}{}
	}
	but.onDoubleTap = func() {
		atomic.StoreInt32(&tapped, 2)
		waitDoubleTapped <- struct{}{}
	}
	w.SetContent(container.NewBorder(nil, nil, nil, nil, but))

	w.mouseMoved(w.viewport, 15, 25)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	<-waitSingleTapped
	w.WaitForEvents()
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, int32(1), atomic.LoadInt32(&tapped), "Single tap should have fired")
	atomic.StoreInt32(&tapped, 0)

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)

	<-waitDoubleTapped
	w.WaitForEvents()
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, int32(2), atomic.LoadInt32(&tapped), "Double tap should have fired")
}

func TestWindow_MouseEventContainsModifierKeys(t *testing.T) {
	w := createWindow("Test").(*window)
	m := &mouseableObject{Rectangle: canvas.NewRectangle(color.White)}
	m.SetMinSize(fyne.NewSize(20, 20))
	w.SetContent(m)

	w.mouseMoved(w.viewport, 7, 7)
	w.WaitForEvents()

	// On OS X a Ctrl+Click is normally translated into a Right-Click.
	// The well-known Ctrl+Click for extending a selection is a Cmd+Click there.
	var superModifier, ctrlModifier fyne.KeyModifier
	if runtime.GOOS == "darwin" {
		superModifier = fyne.KeyModifierControl
		ctrlModifier = 0
	} else {
		superModifier = fyne.KeyModifierSuper
		ctrlModifier = fyne.KeyModifierControl
	}

	tests := map[string]struct {
		modifier              glfw.ModifierKey
		expectedEventModifier fyne.KeyModifier
	}{
		"no modifier key": {
			modifier:              0,
			expectedEventModifier: 0,
		},
		"shift": {
			modifier:              glfw.ModShift,
			expectedEventModifier: fyne.KeyModifierShift,
		},
		"ctrl": {
			modifier:              glfw.ModControl,
			expectedEventModifier: ctrlModifier,
		},
		"alt": {
			modifier:              glfw.ModAlt,
			expectedEventModifier: fyne.KeyModifierAlt,
		},
		"super": {
			modifier:              glfw.ModSuper,
			expectedEventModifier: superModifier,
		},
		"shift+ctrl": {
			modifier:              glfw.ModShift | glfw.ModControl,
			expectedEventModifier: fyne.KeyModifierShift | ctrlModifier,
		},
		"shift+alt": {
			modifier:              glfw.ModShift | glfw.ModAlt,
			expectedEventModifier: fyne.KeyModifierShift | fyne.KeyModifierAlt,
		},
		"shift+super": {
			modifier:              glfw.ModShift | glfw.ModSuper,
			expectedEventModifier: fyne.KeyModifierShift | superModifier,
		},
		"ctrl+alt": {
			modifier:              glfw.ModControl | glfw.ModAlt,
			expectedEventModifier: ctrlModifier | fyne.KeyModifierAlt,
		},
		"ctrl+super": {
			modifier:              glfw.ModControl | glfw.ModSuper,
			expectedEventModifier: ctrlModifier | superModifier,
		},
		"alt+super": {
			modifier:              glfw.ModAlt | glfw.ModSuper,
			expectedEventModifier: fyne.KeyModifierAlt | superModifier,
		},
		"shift+ctrl+alt": {
			modifier:              glfw.ModShift | glfw.ModControl | glfw.ModAlt,
			expectedEventModifier: fyne.KeyModifierShift | ctrlModifier | fyne.KeyModifierAlt,
		},
		"shift+ctrl+super": {
			modifier:              glfw.ModShift | glfw.ModControl | glfw.ModSuper,
			expectedEventModifier: fyne.KeyModifierShift | ctrlModifier | superModifier,
		},
		"shift+alt+super": {
			modifier:              glfw.ModShift | glfw.ModAlt | glfw.ModSuper,
			expectedEventModifier: fyne.KeyModifierShift | fyne.KeyModifierAlt | superModifier,
		},
		"ctrl+alt+super": {
			modifier:              glfw.ModControl | glfw.ModAlt | glfw.ModSuper,
			expectedEventModifier: ctrlModifier | fyne.KeyModifierAlt | superModifier,
		},
		"shift+ctrl+alt+super": {
			modifier:              glfw.ModShift | glfw.ModControl | glfw.ModAlt | glfw.ModSuper,
			expectedEventModifier: fyne.KeyModifierShift | ctrlModifier | fyne.KeyModifierAlt | superModifier,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Nil(t, m.popMouseEvent(), "no initial mouse event")
			w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, tt.modifier)
			w.WaitForEvents()
			me, _ := m.popMouseEvent().(*desktop.MouseEvent)
			if assert.NotNil(t, me, "mouse event triggered") {
				assert.Equal(t, tt.expectedEventModifier, me.Modifier, "expect modifier to be correct")
			}
		})
	}
}

func TestWindow_SetTitle(t *testing.T) {
	w := createWindow("Test")

	title := "My title"
	w.SetTitle(title)

	assert.Equal(t, title, w.Title())
}

func TestWindow_SetIcon(t *testing.T) {
	w := createWindow("Test")
	assert.Equal(t, fyne.CurrentApp().Icon(), w.Icon())

	newIcon := theme.FyneLogo()
	w.SetIcon(newIcon)
	assert.Equal(t, newIcon, w.Icon())
}

func TestWindow_PixelSize(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)

	rect := &canvas.Rectangle{}
	rect.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(container.NewWithoutLayout(rect))
	w.Canvas().Refresh(w.Content())

	winW, winH := w.(*window).minSizeOnScreen()
	assert.Equal(t, internal.ScaleInt(w.Canvas(), 100), winW)
	assert.Equal(t, internal.ScaleInt(w.Canvas(), 100), winH)
}

var scaleTests = []struct {
	user, system, detected, expected float32
	name                             string
}{
	{1.0, 1.0, 1.0, 1.0, "Windows with user setting 1.0"},
	{1.5, 1.0, 1.0, 1.5, "Windows with user setting 1.5"},

	{1.0, scaleAuto, 1.0, 1.0, "Linux lowDPI with user setting 1.0"},
	{1.5, scaleAuto, 1.0, 1.5, "Linux lowDPI with user setting 1.5"},

	{1.0, scaleAuto, 2.0, 2.0, "Linux highDPI with user setting 1.0"},
	{1.5, scaleAuto, 2.0, 3.0, "Linux highDPI with user setting 1.5"},
}

func TestWindow_calculateScale(t *testing.T) {
	for _, tt := range scaleTests {
		t.Run(tt.name, func(t *testing.T) {
			calculated := calculateScale(tt.user, tt.system, tt.detected)
			assert.Equal(t, tt.expected, calculated)
		})
	}
}

func TestWindow_Padded(t *testing.T) {
	w := createWindow("Test")
	content := canvas.NewRectangle(color.White)
	w.SetContent(content)

	width, _ := w.(*window).minSizeOnScreen()
	assert.Equal(t, int(theme.Padding()*2+content.MinSize().Width), width)
	assert.Equal(t, theme.Padding(), content.Position().X)
}

func TestWindow_SetPadded(t *testing.T) {
	var menuHeight float32
	if hasNativeMenu() {
		menuHeight = 0
	} else {
		menuHeight = canvas.NewText("", color.Black).MinSize().Height + theme.Padding()*2
	}
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	tests := []struct {
		name               string
		padding            bool
		menu               bool
		expectedPad        float32
		expectedMenuHeight float32
	}{
		{"window without padding", false, false, 0, 0},
		{"window with padding", true, false, theme.Padding(), 0},
		{"window with menu without padding", false, true, 0, menuHeight},
		{"window with menu and padding", true, true, theme.Padding(), menuHeight},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createWindow("Test").(*window)
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
	w := createWindow("Test").(*window)

	e1 := widget.NewEntry()
	e2 := widget.NewEntry()

	w.SetContent(container.NewVBox(e1, e2))
	w.Canvas().Focus(e1)

	w.charInput(w.viewport, 'a')
	w.charInput(w.viewport, 'b')
	w.charInput(w.viewport, 'c')
	w.charInput(w.viewport, 'd')
	w.keyPressed(w.viewport, glfw.KeyTab, 0, glfw.Press, 0)
	w.WaitForEvents()

	w.keyPressed(w.viewport, glfw.KeyTab, 0, glfw.Release, 0)
	w.charInput(w.viewport, 'e')
	w.charInput(w.viewport, 'f')
	w.WaitForEvents()

	assert.Equal(t, "abcd", e1.Text)
	assert.Equal(t, "ef", e2.Text)
}

func TestWindow_CaptureTypedShortcut(t *testing.T) {
	w := createWindow("Test").(*window)
	content := &typedShortcutable{}
	content.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(content)
	repaintWindow(w)

	w.Canvas().Focus(content)

	w.keyPressed(nil, glfw.KeyLeftControl, 0, glfw.Press, glfw.ModControl)
	w.keyPressed(nil, glfw.KeyLeftShift, 0, glfw.Press, glfw.ModControl)
	w.keyPressed(nil, glfw.KeyF, 0, glfw.Press, glfw.ModControl)
	w.keyPressed(nil, glfw.KeyLeftShift, 0, glfw.Press, glfw.ModControl)
	w.keyPressed(nil, glfw.KeyLeftControl, 0, glfw.Release, glfw.ModControl)
	w.keyPressed(nil, glfw.KeyF, 0, glfw.Release, glfw.ModControl)

	w.WaitForEvents()

	assert.Equal(t, 1, len(content.capturedShortcuts))
	assert.Equal(t, "CustomDesktop:Control+F", content.capturedShortcuts[0].ShortcutName())
}

func TestWindow_ManualFocus(t *testing.T) {
	w := createWindow("Test").(*window)
	content := &focusable{}
	content.SetMinSize(fyne.NewSize(10, 10))
	w.SetContent(content)
	repaintWindow(w)

	w.mouseMoved(w.viewport, 9, 9)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Release, 0)
	w.WaitForEvents()
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	w.canvas.Focus(content)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	w.canvas.Unfocus()
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)

	content.Disable()
	w.canvas.Focus(content)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)

	w.mouseClicked(w.viewport, glfw.MouseButton1, glfw.Press, 0)
	w.WaitForEvents()
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)
}

func TestWindow_Clipboard(t *testing.T) {
	w := createWindow("Test")

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

func TestWindow_ClipboardCopy_DisabledEntry(t *testing.T) {
	w := createWindow("Test").(*window)
	e := widget.NewEntry()
	e.SetText("Testing")
	e.Disable()
	w.SetContent(e)
	repaintWindow(w)

	w.canvas.Focus(e)
	e.DoubleTapped(nil)
	assert.Equal(t, "Testing", e.SelectedText())

	ctrlMod := glfw.ModControl
	if runtime.GOOS == "darwin" {
		ctrlMod = glfw.ModSuper
	}
	w.keyPressed(nil, glfw.KeyC, 0, glfw.Repeat, ctrlMod)
	w.WaitForEvents()

	assert.Equal(t, "Testing", w.Clipboard().Content())

	e.SetText("Testing2")
	e.DoubleTapped(nil)
	assert.Equal(t, "Testing2", e.SelectedText())

	// any other shortcut should be forbidden (Cut)
	w.keyPressed(nil, glfw.KeyX, 0, glfw.Repeat, ctrlMod)
	w.WaitForEvents()

	assert.Equal(t, "Testing2", e.Text)
	assert.Equal(t, "Testing", w.Clipboard().Content())

	// any other shortcut should be forbidden (Paste)
	w.keyPressed(nil, glfw.KeyV, 0, glfw.Repeat, ctrlMod)
	w.WaitForEvents()

	assert.Equal(t, "Testing2", e.Text)
	assert.Equal(t, "Testing", w.Clipboard().Content())
}

func TestWindow_CloseInterception(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("test").(*window)
	w.create()

	onIntercepted := false
	onClosed := false
	w.SetCloseIntercept(func() {
		onIntercepted = true
	})
	w.SetOnClosed(func() {
		onClosed = true
	})
	w.Close()
	w.WaitForEvents()
	assert.False(t, onIntercepted) // The interceptor is not called by the Close.
	assert.True(t, onClosed)
	assert.True(t, w.viewport.ShouldClose()) // For #2694

	w.closing = false // fake a fresh window
	onIntercepted = false
	onClosed = false
	w.closed(w.viewport)
	w.WaitForEvents()
	assert.True(t, onIntercepted) // The interceptor is called by the closed.
	assert.False(t, onClosed)     // If the interceptor is set Close is not called.

	w.closing = false // fake a fresh window
	onClosed = false
	w.SetCloseIntercept(nil)
	w.closed(w.viewport)
	w.WaitForEvents()
	assert.True(t, onClosed) // Close is called if the interceptor is not set.
}

func TestWindow_SetContent_Twice(t *testing.T) {
	w := createWindow("Test").(*window)

	e1 := widget.NewLabel("1")
	e2 := widget.NewLabel("2")

	w.SetContent(e1)
	assert.True(t, e1.Visible())
	w.SetContent(e2)
	assert.True(t, e2.Visible())
	w.SetContent(e1)
	assert.True(t, e1.Visible())
}

func TestWindow_SetFullScreen(t *testing.T) {
	w := createWindow("Full").(*window)
	w.SetFullScreen(true)
	w.create()
	w.doShow()
	waitForMain()

	// initial state - no window size set (except darwin?)
	if runtime.GOOS != "darwin" {
		assert.Zero(t, w.width)
		assert.Zero(t, w.height)
	}

	w.SetFullScreen(false)
	waitForMain()
	// ensure we realised size now!
	assert.NotZero(t, w.width)
	assert.NotZero(t, w.height)
}

// This test makes our developer screens flash, let's not run it regularly...
// func TestWindow_Shortcut(t *testing.T) {
//	w := createWindow("Test")
//
//	shortcutFullScreenWindow := &desktop.CustomShortcut{
//		KeyName: fyne.KeyF12,
//	}
//
//	w.Canvas().AddShortcut(shortcutFullScreenWindow, func(sc fyne.Shortcut) {
//		w.SetFullScreen(true)
//	})
//
//	assert.False(t, w.FullScreen())
//
//	w.Canvas().(*glCanvas).shortcut.TypedShortcut(shortcutFullScreenWindow)
//	assert.True(t, w.FullScreen())
// }

func createWindow(title string) fyne.Window {
	w := d.CreateWindow(title)
	w.(*window).create()
	return w
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

type draggableTappableObject struct {
	*canvas.Rectangle
	draggable
	tappable
}

var _ fyne.Focusable = (*focusable)(nil)
var _ fyne.Disableable = (*focusable)(nil)

type focusable struct {
	canvas.Rectangle
	id             string // helps identifying instances in comparisons
	focused        bool
	focusedTimes   int
	unfocusedTimes int
	disabled       bool
}

func (f *focusable) Tapped(*fyne.PointEvent) {
	d.CanvasForObject(f).Focus(f)
}

func (f *focusable) TypedRune(rune) {
}

func (f *focusable) TypedKey(*fyne.KeyEvent) {
}

func (f *focusable) FocusGained() {
	f.focusedTimes++
	if f.Disabled() {
		return
	}
	f.focused = true
}

func (f *focusable) FocusLost() {
	f.unfocusedTimes++
	f.focused = false
}

func (f *focusable) Enable() {
	f.disabled = false
}

func (f *focusable) Disable() {
	f.disabled = true
}

func (f *focusable) Disabled() bool {
	return f.disabled
}

type typedShortcutable struct {
	focusable
	capturedShortcuts []fyne.Shortcut
}

func (ts *typedShortcutable) TypedShortcut(s fyne.Shortcut) {
	ts.capturedShortcuts = append(ts.capturedShortcuts, s)
}

var _ fyne.Scrollable = (*scrollable)(nil)

type scrollable struct {
	*canvas.Rectangle
	events []interface{}
}

func (s *scrollable) Scrolled(e *fyne.ScrollEvent) {
	s.events = append(s.events, e)
}

func (s *scrollable) popScrollEvent() (e interface{}) {
	e, s.events = pop(s.events)
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

type doubleTappableButton struct {
	widget.Button

	onDoubleTap func()
}

func (t *doubleTappableButton) DoubleTapped(_ *fyne.PointEvent) {
	t.onDoubleTap()
}

func newDoubleTappableButton() *doubleTappableButton {
	but := &doubleTappableButton{}
	but.ExtendBaseWidget(but)

	return but
}

func waitForMain() {
	runOnMain(func() {}) // this blocks until processed
}
