package sensor

import (
	"image"
)

// A camera device is a device that has hardware support for taking photos.
//
// Since: 2.7
type CameraDevice interface {
	// CapturePhoto opens a native photo capture view that allows the user to optionally
	// take a photograph.  It returns the image, a boolean to indicate if an image was
	// taken, and any errors.
	//
	// Since: 2.7
	CapturePhoto() (image.Image, bool, error)
}
