package icns

import (
	"fmt"
	"image"
)

// ErrImageTooSmall is returned when the image is too small to process.
type ErrImageTooSmall struct {
	need  int
	image image.Image
}

func (err ErrImageTooSmall) Error() string {
	b := err.image.Bounds().Max
	format := "image is too small: %dx%d, need at least %dx%d"
	return fmt.Sprintf(format, b.X, b.Y, err.need, err.need)
}

func panicf(format string, values ...interface{}) {
	panic(fmt.Sprintf(format, values...))
}
