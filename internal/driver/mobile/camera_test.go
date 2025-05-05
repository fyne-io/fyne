package mobile

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/sensor"
	"github.com/stretchr/testify/assert"
)

func TestDevice_CapturePhoto(t *testing.T) {
	var dev fyne.Device
	dev = &device{}

	camdev, ok := dev.(sensor.CameraDevice)
	assert.True(t, ok)

	// TODO change this assert ;)
	assert.Panics(t, func() {
		_ = camdev.CapturePhoto()
	})
}
