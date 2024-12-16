package fyne

// Layout defines how [CanvasObject]s may be laid out in a specified Size.
type Layout interface {
	// Layout will manipulate the listed [CanvasObject]s Size and Position
	// to fit within the specified size.
	Layout([]CanvasObject, Size)
	// MinSize calculates the smallest size that will fit the listed
	// [CanvasObject]s using this Layout algorithm.
	MinSize(objects []CanvasObject) Size
}
