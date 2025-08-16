//go:build android

package mobile

import (
	"bytes"
	"encoding/base64"
	"image"
	"io"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

func (*device) CapturePhoto() (image.Image, bool, error) {
	slog.Debug("camera.Open")

	var r io.Reader
	done := make(chan struct{})
	defer close(done)
	var openErr error
	ShowCameraOpen(func(reader io.Reader, err error) {
		slog.Debug("camera.Open: in callback", "reader", reader, "error", err)
		if err != nil {
			fyne.LogError("camera open: fyne returns error:", err)
			openErr = err
			done <- struct{}{}
			return
		}

		r = reader
		done <- struct{}{}
		slog.Debug("camera.Open: in callback: done")
	})

	slog.Debug("camera.Open: <-done")
	<-done

	if openErr != nil {
		return nil, false, openErr
	}

	if r == nil {
		return nil, false, nil
	}

	encoded, err := io.ReadAll(r)
	if err != nil {
		return nil, true, err
	}

	if string(encoded) == "" {
		return nil, false, nil
	}

	data, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return nil, true, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, true, err
	}

	return img, true, nil
}

func ShowCameraOpen(callback func(r io.Reader, err error)) {
	app.NativeShowCameraOpen(func(path string, closer func()) {
		slog.Debug("app show camera open", "path", path)
		if path == "" {
			callback(nil, nil)
			return
		}

		buf := bytes.NewBufferString(path)

		slog.Debug("app show camera open: call callback")

		callback(buf, nil)
	})
}
