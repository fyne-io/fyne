package icns

import (
	"github.com/nfnt/resize"
)

// InterpolationFunction is the algorithm used to resize the image.
type InterpolationFunction = resize.InterpolationFunction

// InterpolationFunction constants.
const (
	// Nearest-neighbor interpolation
	NearestNeighbor InterpolationFunction = iota
	// Bilinear interpolation
	Bilinear
	// Bicubic interpolation (with cubic hermite spline)
	Bicubic
	// Mitchell-Netravali interpolation
	MitchellNetravali
	// Lanczos interpolation (a=2)
	Lanczos2
	// Lanczos interpolation (a=3)
	Lanczos3
)
