package fyne

// Canvas defines a graphical canvas to which a CanvasObject or Container can be added.
// Each canvas has a scale which is automatically applied during the render process.
type Canvas interface {
	Content() CanvasObject
	SetContent(CanvasObject)
	Refresh(CanvasObject)
	Contains(CanvasObject) bool
	Focus(FocusableObject)
	Focused() FocusableObject

	Size() Size
	Scale() float32
	SetScale(float32)

	SetOnKeyDown(func(*KeyEvent))
}
