package test

import "fyne.io/fyne"

// Click simulates a mouse click on the specified object
func Click(obj fyne.TappableObject) {
	ev := &fyne.PointerEvent{Button: fyne.LeftMouseButton, Position: fyne.NewPos(1, 1)}
	obj.OnTap(ev)
}

func typeChars(chars string, keyDown func(event *fyne.KeyEvent)) {
	for _, char := range chars {
		ev := &fyne.KeyEvent{String: string(char), Name: fyne.KeyName(char)}
		keyDown(ev)
	}
}

// Type performs a series of key events to simulate typing of a value into the specified object.
// The focusable object will be focused before typing begins.
// The chars parameter will be input one rune at a time to the focused object.
func Type(obj fyne.FocusableObject, chars string) {
	obj.OnFocusGained()

	typeChars(chars, obj.OnKeyDown)
}

// TypeOnCanvas is like the Type function but it passes the key events to the canvas object
// rather than a focusable widget.
func TypeOnCanvas(c fyne.Canvas, chars string) {
	typeChars(chars, c.OnKeyDown())
}
