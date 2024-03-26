//go:build flatpak && !windows && !android && !ios && !wasm && !js

package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/rymdport/portal/filechooser"
)

func openFolder(parentWindowHandle string, callback func(fyne.ListableURI, error), options *filechooser.OpenFileOptions) {
	uris, err := filechooser.OpenFile(parentWindowHandle, "Open Folder", options)
	if err != nil {
		callback(nil, err)
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

	callback(storage.ListerForURI(uri))
}

func openFile(parentWindowHandle string, callback func(fyne.URIReadCloser, error), options *filechooser.OpenFileOptions) {
	uris, err := filechooser.OpenFile(parentWindowHandle, "Open File", options)
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

	callback(storage.Reader(uri))
}

func fileOpenOSOverride(d *FileDialog) bool {
	go func() {
		folderCallback, folder := d.callback.(func(fyne.ListableURI, error))
		options := &filechooser.OpenFileOptions{
			Directory:   folder,
			AcceptLabel: d.confirmText,
		}
		if d.startingLocation != nil {
			options.CurrentFolder = d.startingLocation.Path()
		}

		parentWindowHandle := d.parent.(interface{ GetWindowHandle() string }).GetWindowHandle()

		if folder {
			openFolder(parentWindowHandle, folderCallback, options)
			return
		}

		fileCallback := d.callback.(func(fyne.URIReadCloser, error))
		openFile(parentWindowHandle, fileCallback, options)
	}()
	return true
}

func fileSaveOSOverride(d *FileDialog) bool {
	go func() {
		options := &filechooser.SaveFileOptions{
			AcceptLabel: d.confirmText,
			CurrentName: d.initialFileName,
		}
		if d.startingLocation != nil {
			options.CurrentFolder = d.startingLocation.Path()
		}

		parentWindowHandle := d.parent.(interface{ GetWindowHandle() string }).GetWindowHandle()

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
