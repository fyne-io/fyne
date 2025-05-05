package sensor

import (
	"image"
)

type CameraDevice interface {
	CapturePhoto() image.Image
}
