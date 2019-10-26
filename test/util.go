package test

import "fyne.io/fyne"

// Tap simulates a left mouse click on the specified object.
func Tap(obj fyne.Tappable) {
	TapAt(obj, fyne.NewPos(1, 1))
}

// TapAt simulates a left mouse click on the passed object at a specified place within it.
func TapAt(obj fyne.Tappable, pos fyne.Position) {
	if focus, ok := obj.(fyne.Focusable); ok {
		if focus != Canvas().Focused() {
			Canvas().Focus(focus)
		}
	}

	ev := &fyne.PointEvent{Position: pos}
	obj.Tapped(ev)
}

// TapSecondary simulates a right mouse click on the specified object.
func TapSecondary(obj fyne.Tappable) {
	TapSecondaryAt(obj, fyne.NewPos(1, 1))
}

// TapSecondaryAt simulates a right mouse click on the passed object at a specified place within it.
func TapSecondaryAt(obj fyne.Tappable, pos fyne.Position) {
	if focus, ok := obj.(fyne.Focusable); ok {
		if focus != Canvas().Focused() {
			Canvas().Focus(focus)
		}
	}

	ev := &fyne.PointEvent{Position: pos}
	obj.TappedSecondary(ev)
}

func typeChars(chars []rune, keyDown func(rune)) {
	for _, char := range chars {
		keyDown(char)
	}
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
