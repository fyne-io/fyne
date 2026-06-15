package sensor

import (
	"image"
)

// A camera device is a device that has hardware support for taking photos.
//
// Since: 2.9
type CameraDevice interface {
	// CapturePhoto requests a single high quality photo from the camera.
	//
	// Since: 2.9
	CapturePhoto() (image.Image, error)

	// StartPreview tells the camera to begin capturing lower-resolution frames from the
	// camera and stream them into the channel returned by Preview().
	//
	// Since: 2.9
	StartPreview()

	// StopPreview tells the camera to stop sending frames to the preview channel
	// accessible via Preview().
	//
	// Since: 2.9
	StopPreview()

	// Preview provides a channel for image frames that can be used in a view finder or
	// barcode scanner.  Images start appearing when StartPreview() is called, and will
	// continue until StopPreview() is called.
	//
	// Since: 2.9
	Preview() chan image.Image
}
