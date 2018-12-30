package test

import "github.com/fyne-io/fyne"

// Click simulates a mouse click on the specified object
func Click(obj fyne.ClickableObject) {
	ev := &fyne.MouseEvent{Button: fyne.LeftMouseButton, Position: fyne.NewPos(1, 1)}
	obj.OnMouseDown(ev)
}

// Type performs a series of key events to simulate typing of a value into the specified object.
// The focusable object will be focused before typing begins.
// The chars parameter will be input one rune at a time to the focused object.
func Type(obj fyne.FocusableObject, chars string) {
	obj.OnFocusGained()

	for _, char := range chars {
		ev := &fyne.KeyEvent{String: string(char), Name: fyne.KeyName(char)}
		obj.OnKeyDown(ev)
	}
}
