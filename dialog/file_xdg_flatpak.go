//go:build flatpak && !windows && !android && !ios && !wasm && !js
// +build flatpak,!windows,!android,!ios,!wasm,!js

package dialog

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/rymdport/portal/filechooser"
)

func fileOpenOSOverride(d *FileDialog) bool {
	go func() {
		folderCallback, folder := d.callback.(func(fyne.ListableURI, error))
		options := &filechooser.OpenOptions{
			Modal:       true,
			Directory:   folder,
			AcceptLabel: d.confirmText,
		}
		if d.startingLocation != nil {
			options.Location = d.startingLocation.Path()
		}

		parentWindowHandle := d.parent.(interface{ GetWindowHandle() string }).GetWindowHandle()

		if folder {
			uris, err := filechooser.OpenFile(parentWindowHandle, "Open Folder", options)
			if err != nil {
				folderCallback(nil, err)
			}

			if len(uris) == 0 {
				folderCallback(nil, nil)
				return
			}

			uri, err := storage.ParseURI(uris[0])
			if err != nil {
				folderCallback(nil, err)
				return
			}

			folderCallback(storage.ListerForURI(uri))
			return
		}

		uris, err := filechooser.OpenFile(parentWindowHandle, "Open File", options)
		fileCallback := d.callback.(func(fyne.URIReadCloser, error))
		if err != nil {
			fileCallback(nil, err)
			return
		}

		if len(uris) == 0 {
			fileCallback(nil, nil)
			return
		}

		uri, err := storage.ParseURI(uris[0])
		if err != nil {
			fileCallback(nil, err)
			return
		}

		fileCallback(storage.Reader(uri))
	}()
	return true
}

func fileSaveOSOverride(d *FileDialog) bool {
	go func() {
		options := &filechooser.SaveSingleOptions{
			Modal:       true,
			AcceptLabel: d.confirmText,
			FileName:    d.initialFileName,
		}
		if d.startingLocation != nil {
			options.Location = d.startingLocation.Path()
		}

		parentWindowHandle := d.parent.(interface{ GetWindowHandle() string }).GetWindowHandle()
		fmt.Println(parentWindowHandle)

		callback := d.callback.(func(fyne.URIWriteCloser, error))
		uris, err := filechooser.SaveFile(parentWindowHandle, "Open File", options)
		if err != nil {
			callback(nil, err)
			return
		}

		if len(uris) == 0 {
			callback(nil, nil)
			return
		}

		uri, err := storage.ParseURI(uris[0])
		if err != nil {
			callback(nil, err)
			return
		}

		callback(storage.Writer(uri))
	}()
	return true
}
