package fyne

import "image"

// Canvas defines a graphical canvas to which a CanvasObject or Container can be added.
// Each canvas has a scale which is automatically applied during the render process.
type Canvas interface {
	Content() CanvasObject
	SetContent(CanvasObject)
	Refresh(CanvasObject)
	Focus(Focusable)
	Unfocus()
	Focused() Focusable

	// Size returns the current size of this canvas
	Size() Size
	// Scale returns the current scale (multiplication factor) this canvas uses to render
	// The pixel size of a CanvasObject can be found by multiplying by this value.
	Scale() float32
	// SetScale sets ths scale for this canvas only, overriding system and user settings.
	//
	// Deprecated: SetScale will be replaced by system wide settings in the future
	SetScale(float32)

	Overlay() CanvasObject
	SetOverlay(CanvasObject)

	OnTypedRune() func(rune)
	SetOnTypedRune(func(rune))
	OnTypedKey() func(*KeyEvent)
	SetOnTypedKey(func(*KeyEvent))
	AddShortcut(shortcut Shortcut, handler func(shortcut Shortcut))

	Capture() image.Image
}
