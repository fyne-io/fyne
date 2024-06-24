//go:build flatpak && !windows && !android && !ios && !wasm && !js

package dialog

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/internal/build"
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

		parentWindowHandle := windowHandleForPortal(d.parent)

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

		parentWindowHandle := windowHandleForPortal(d.parent)

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

func x11WindowHandleToString(handle uintptr) string {
	return "x11:" + strconv.FormatUint(uint64(handle), 16)
}

func windowHandleForPortal(window fyne.Window) string {
	parentWindowHandle := ""
	if !build.IsWayland {
		window.(driver.NativeWindow).RunNative(func(context any) {
			handle := context.(driver.X11WindowContext).WindowHandle
			parentWindowHandle = x11WindowHandleToString(handle)
		})
	}

	// TODO: We need to get the Wayland handle from the xdg_foreign protocol and convert to string on the form "wayland:{id}".
	return parentWindowHandle
}
