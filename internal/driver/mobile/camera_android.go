//go:build android

package mobile

import (
	"bytes"
	"encoding/base64"
	"image"
	"sync"

	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

func (*device) CapturePhoto() (image.Image, bool, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	var result string
	app.NativeShowCameraOpen(func(data string) {
		result = data
		wg.Done()
	})
	wg.Wait()

	if result == "" {
		return nil, false, nil
	}

	data, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		return nil, true, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, true, err
	}

	return img, true, nil
}
