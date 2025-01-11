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

func openFolder(parentWindowHandle string, options *filechooser.OpenFileOptions) (fyne.ListableURI, error) {
	uris, err := filechooser.OpenFile(parentWindowHandle, "Open Folder", options)
	if err != nil {
		return nil, err
	}

	if len(uris) == 0 {
		return nil, nil
	}

	uri, err := storage.ParseURI(uris[0])
	if err != nil {
		return nil, err
	}

	return storage.ListerForURI(uri)
}

func openFile(parentWindowHandle string, options *filechooser.OpenFileOptions) (fyne.URIReadCloser, error) {
	uris, err := filechooser.OpenFile(parentWindowHandle, "Open File", options)
	if err != nil {
		return nil, err
	}

	if len(uris) == 0 {
		return nil, nil
	}

	uri, err := storage.ParseURI(uris[0])
	if err != nil {
		return nil, err
	}

	return storage.Reader(uri)
}

func saveFile(parentWindowHandle string, options *filechooser.SaveFileOptions) (fyne.URIWriteCloser, error) {
	uris, err := filechooser.SaveFile(parentWindowHandle, "Save File", options)
	if err != nil {
		return nil, err
	}

	if len(uris) == 0 {
		return nil, nil
	}

	uri, err := storage.ParseURI(uris[0])
	if err != nil {
		return nil, err
	}

	return storage.Writer(uri)
}

func fileOpenOSOverride(d *FileDialog) bool {
	folderCallback, folder := d.callback.(func(fyne.ListableURI, error))
	fileCallback, _ := d.callback.(func(fyne.URIReadCloser, error))
	options := &filechooser.OpenFileOptions{
		Directory:   folder,
		AcceptLabel: d.confirmText,
	}

	if d.startingLocation != nil {
		options.CurrentFolder = d.startingLocation.Path()
	}

	windowHandle := windowHandleForPortal(d.parent)

	if folder {
		go func() {
			folder, err := openFolder(windowHandle, options)
			fyne.CurrentApp().Driver().CallFromGoroutine(func() {
				folderCallback(folder, err)
			})
		}()
	} else {
		go func() {
			file, err := openFile(windowHandle, options)
			fyne.CurrentApp().Driver().CallFromGoroutine(func() {
				fileCallback(file, err)
			})
		}()
	}

	return true
}

func fileSaveOSOverride(d *FileDialog) bool {
	options := &filechooser.SaveFileOptions{
		AcceptLabel: d.confirmText,
		CurrentName: d.initialFileName,
	}
	if d.startingLocation != nil {
		options.CurrentFolder = d.startingLocation.Path()
	}

	callback := d.callback.(func(fyne.URIWriteCloser, error))
	windowHandle := windowHandleForPortal(d.parent)

	go func() {
		file, err := saveFile(windowHandle, options)
		fyne.CurrentApp().Driver().CallFromGoroutine(func() {
			callback(file, err)
		})
	}()

	return true
}

func x11WindowHandleToString(handle uintptr) string {
	return "x11:" + strconv.FormatUint(uint64(handle), 16)
}

func windowHandleForPortal(window fyne.Window) string {
	windowHandle := ""
	if !build.IsWayland {
		window.(driver.NativeWindow).RunNative(func(context any) {
			handle := context.(driver.X11WindowContext).WindowHandle
			windowHandle = x11WindowHandleToString(handle)
		})
	}

	// TODO: We need to get the Wayland handle from the xdg_foreign protocol and convert to string on the form "wayland:{id}".
	return windowHandle
}
