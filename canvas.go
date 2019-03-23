package fyne

// Canvas defines a graphical canvas to which a CanvasObject or Container can be added.
// Each canvas has a scale which is automatically applied during the render process.
type Canvas interface {
	Content() CanvasObject
	SetContent(CanvasObject)
	Refresh(CanvasObject)
	Focus(Focusable)
	Unfocus()
	Focused() Focusable

	Size() Size
	Scale() float32
	SetScale(float32)

	OnTypedRune() func(rune)
	SetOnTypedRune(func(rune))
	OnTypedKey() func(*KeyEvent)
	SetOnTypedKey(func(*KeyEvent))
	AddShortcut(shortcut Shortcut, handler func(shortcut Shortcut))
}
