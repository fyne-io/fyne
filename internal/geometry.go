package internal

// PixelSize describe an object width and height in Pixels
type PixelSize struct {
	WidthPx  int // The number of pixels along the X axis.
	HeightPx int // The number of pixels along the Y axis.
}

// NewPixelSize return a newly allocated PixelSize of the specified dimensions.
func NewPixelSize(w int, h int) PixelSize {
	return PixelSize{w, h}
}
