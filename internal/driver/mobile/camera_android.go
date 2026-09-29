//go:build android

package mobile

import (
	"bytes"
	"image"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

var previewFrames = make(chan image.Image)

func (*device) CapturePhoto() (image.Image, error) {
	var result image.Image
	var resultError error

	var wg sync.WaitGroup
	wg.Add(1)
	app.NativeCapturePhoto(func(data []byte) {
		defer wg.Done()

		if len(data) == 0 {
			fyne.LogError("empty data slice returned to capture photo callback", nil)
			return
		}

		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			fyne.LogError("error decoding data returned to capture photo callback", err)
			return
		}

		result = img
		resultError = err
	})
	wg.Wait()

	return result, resultError
}

func (*device) StartPreview() {
	app.NativeStartPreview(func(data []byte) {
		if len(data) == 0 {
			fyne.LogError("empty data slice returned as camera preview frame", nil)
			return
		}

		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			fyne.LogError("error decoding data returned as camera preview frame", err)
			return
		}

		select {
		case previewFrames <- img:
		default:
		}
	})
}

func (*device) StopPreview() {
	app.NativeStopPreview()
}

func (*device) Preview() chan image.Image {
	return previewFrames
}
