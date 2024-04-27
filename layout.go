package fyne

// Layout defines how CanvasObjects may be laid out in a specified Size.
type Layout interface {
	// Layout will manipulate the listed CanvasObjects Size and Position
	// to fit within the specified size.
	Layout([]CanvasObject, Size)
	// MinSize calculates the smallest size that will fit the listed
	// CanvasObjects using this Layout algorithm.
	MinSize(objects []CanvasObject) Size
	// GetPaddings returns the top, bottom, left and right paddings used
	// by this layout.
	GetPaddings() (float32, float32, float32, float32)
	// SetPaddingFn sets the function that will be called to get the
	// padding values for this layout.
	SetPaddingFn(fn func() (float32, float32, float32, float32))
}
