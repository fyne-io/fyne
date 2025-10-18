package test

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	intdriver "fyne.io/fyne/v2/internal/driver"
)

// RenderObjectToMarkup renders the given [fyne.io/fyne/v2.CanvasObject] to a markup string.
//
// Since: 2.6
func RenderObjectToMarkup(o fyne.CanvasObject) string {
	c := NewCanvas()
	c.SetPadded(false)
	size := o.MinSize().Max(o.Size())
	c.SetContent(o)
	c.Resize(size) // ensure we are large enough for current size

	return snapshot(c)
}

// RenderToMarkup renders the given [fyne.io/fyne/v2.Canvas] to a markup string.
//
// Since: 2.6
func RenderToMarkup(c fyne.Canvas) string {
	return snapshot(c)
}

// Drag drags at an absolute position on the canvas.
// deltaX/Y is the dragging distance: <0 for dragging up/left, >0 for dragging down/right.
func Drag(c fyne.Canvas, pos fyne.Position, deltaX, deltaY float32) {
	matches := func(object fyne.CanvasObject) bool {
		_, ok := object.(fyne.Draggable)
		return ok
	}
	o, p, _ := intdriver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o == nil {
		return
	}
	e := &fyne.DragEvent{
		PointEvent: fyne.PointEvent{Position: p},
		Dragged:    fyne.Delta{DX: deltaX, DY: deltaY},
	}
	o.(fyne.Draggable).Dragged(e)
	o.(fyne.Draggable).DragEnd()
}

// FocusNext focuses the next focusable on the canvas.
func FocusNext(c fyne.Canvas) {
	if tc, ok := c.(*canvas); ok {
		tc.focusManager().FocusNext()
	} else {
		fyne.LogError("FocusNext can only be called with a test canvas", nil)
	}
}

// FocusPrevious focuses the previous focusable on the canvas.
func FocusPrevious(c fyne.Canvas) {
	if tc, ok := c.(*canvas); ok {
		tc.focusManager().FocusPrevious()
	} else {
		fyne.LogError("FocusPrevious can only be called with a test canvas", nil)
	}
}

// LaidOutObjects returns all fyne.CanvasObject starting at the given fyne.CanvasObject which is laid out previously.
func LaidOutObjects(o fyne.CanvasObject) (objects []fyne.CanvasObject) {
	if o != nil {
		objects = layoutAndCollect(objects, o, o.MinSize().Max(o.Size()))
	}
	return objects
}

// MoveMouse simulates a mouse movement to the given position.
func MoveMouse(c fyne.Canvas, pos fyne.Position) {
	if fyne.CurrentDevice().IsMobile() {
		return
	}

	tc, _ := c.(*canvas)
	var oldHovered, hovered desktop.Hoverable
	if tc != nil {
		oldHovered = tc.hovered
	}
	matches := func(object fyne.CanvasObject) bool {
		_, ok := object.(desktop.Hoverable)
		return ok
	}
	o, p, _ := intdriver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o != nil {
		hovered = o.(desktop.Hoverable)
		me := &desktop.MouseEvent{
			PointEvent: fyne.PointEvent{
				AbsolutePosition: pos,
				Position:         p,
			},
		}
		if hovered == oldHovered {
			hovered.MouseMoved(me)
		} else {
			if oldHovered != nil {
				oldHovered.MouseOut()
			}
			hovered.MouseIn(me)
		}
	} else if oldHovered != nil {
		oldHovered.MouseOut()
	}
	if tc != nil {
		tc.hovered = hovered
	}
}

// Scroll scrolls at an absolute position on the canvas.
// deltaX/Y is the scrolling distance: <0 for scrolling up/left, >0 for scrolling down/right.
func Scroll(c fyne.Canvas, pos fyne.Position, deltaX, deltaY float32) {
	matches := func(object fyne.CanvasObject) bool {
		_, ok := object.(fyne.Scrollable)
		return ok
	}
	o, _, _ := intdriver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o == nil {
		return
	}

	e := &fyne.ScrollEvent{Scrolled: fyne.Delta{DX: deltaX, DY: deltaY}}
	o.(fyne.Scrollable).Scrolled(e)
}

// DoubleTap simulates a double left mouse click on the specified object.
func DoubleTap(obj fyne.DoubleTappable) {
	ev, c := prepareTap(obj, fyne.NewPos(1, 1))
	handleFocusOnTap(c, obj)
	obj.DoubleTapped(ev)
}

// Tap simulates a left mouse click on the specified object.
func Tap(obj fyne.Tappable) {
	TapAt(obj, fyne.NewPos(1, 1))
}

// TapAt simulates a left mouse click on the passed object at a specified place within it.
func TapAt(obj fyne.Tappable, pos fyne.Position) {
	ev, c := prepareTap(obj, pos)
	tap(c, obj, ev)
}

// TapCanvas taps at an absolute position on the canvas.
func TapCanvas(c fyne.Canvas, pos fyne.Position) {
	if o, p := findTappable(c, pos); o != nil {
		tap(c, o.(fyne.Tappable), &fyne.PointEvent{AbsolutePosition: pos, Position: p})
	}
}

// TapSecondary simulates a right mouse click on the specified object.
func TapSecondary(obj fyne.SecondaryTappable) {
	TapSecondaryAt(obj, fyne.NewPos(1, 1))
}

// TapSecondaryAt simulates a right mouse click on the passed object at a specified place within it.
func TapSecondaryAt(obj fyne.SecondaryTappable, pos fyne.Position) {
	ev, c := prepareTap(obj, pos)
	handleFocusOnTap(c, obj)
	obj.TappedSecondary(ev)
}

// Type performs a series of key events to simulate typing of a value into the specified object.
// The focusable object will be focused before typing begins.
// The chars parameter will be input one rune at a time to the focused object.
func Type(obj fyne.Focusable, chars string) {
	obj.FocusGained()

	typeChars([]rune(chars), obj.TypedRune)
}

// TypeOnCanvas is like the Type function but it passes the key events to the canvas object
// rather than a focusable widget.
func TypeOnCanvas(c fyne.Canvas, chars string) {
	typeChars([]rune(chars), c.OnTypedRune())
}

// WidgetRenderer allows test scripts to gain access to the current renderer for a widget.
// This can be used for verifying correctness of rendered components for a widget in unit tests.
func WidgetRenderer(wid fyne.Widget) fyne.WidgetRenderer {
	return cache.Renderer(wid)
}

func findTappable(c fyne.Canvas, pos fyne.Position) (o fyne.CanvasObject, p fyne.Position) {
	matches := func(object fyne.CanvasObject) bool {
		_, ok := object.(fyne.Tappable)
		return ok
	}
	o, p, _ = intdriver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	return o, p
}

func prepareTap(obj any, pos fyne.Position) (*fyne.PointEvent, fyne.Canvas) {
	d := fyne.CurrentApp().Driver()
	ev := &fyne.PointEvent{Position: pos}
	var c fyne.Canvas
	if co, ok := obj.(fyne.CanvasObject); ok {
		c = d.CanvasForObject(co)
		ev.AbsolutePosition = d.AbsolutePositionForObject(co).Add(pos)
	}
	return ev, c
}

func tap(c fyne.Canvas, obj fyne.Tappable, ev *fyne.PointEvent) {
	handleFocusOnTap(c, obj)
	obj.Tapped(ev)
}

func handleFocusOnTap(c fyne.Canvas, obj any) {
	if c == nil {
		return
	}

	if focus, ok := obj.(fyne.Focusable); ok {
		dis, ok := obj.(fyne.Disableable)
		if (!ok || !dis.Disabled()) && focus == c.Focused() {
			return
		}
	}

	c.Unfocus()
}

func typeChars(chars []rune, keyDown func(rune)) {
	for _, char := range chars {
		keyDown(char)
	}
}

func writeMarkup(path string, markup string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(markup), 0o644)
}
