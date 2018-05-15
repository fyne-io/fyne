package ui

// Canvas defines a graphical canvas to which a CanvasObject or Container can be added.
// Each canvas has a scale which is automatically applied during the render process.
type Canvas interface {
	Content() CanvasObject
	SetContent(CanvasObject)
	Refresh(CanvasObject)
	Contains(CanvasObject) bool

	Size() Size
	Scale() float32
	SetScale(float32)

	SetOnKeyDown(func(*KeyEvent))
}

// GetCanvas returns the canvas containing the passed CanvasObject.
// It will return nil if the CanvasObject has not been added to any Canvas.
func GetCanvas(obj CanvasObject) Canvas {
	for _, window := range GetDriver().AllWindows() {
		if window.Canvas() != nil && window.Canvas().Contains(obj) {
			return window.Canvas()
		}
	}

	return nil
}
