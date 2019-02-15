package test

import "fyne.io/fyne"

// Tap simulates a mouse click on the specified object
func Tap(obj fyne.Tappable) {
	ev := &fyne.PointEvent{Position: fyne.NewPos(1, 1)}
	obj.Tapped(ev)
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
