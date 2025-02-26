//go:build flatpak && !windows && !android && !ios && !wasm && !js

package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/storage"

	"github.com/rymdport/portal"
	"github.com/rymdport/portal/filechooser"
)

func openFile(parentWindowHandle string, options *filechooser.OpenFileOptions) (fyne.URIReadCloser, error) {
	title := lang.L("Open") + " " + lang.L("File")
	uri, err := open(parentWindowHandle, title, options)
	if err != nil || uri == nil {
		return nil, err
	}

	return storage.Reader(uri)
}

func openFolder(parentWindowHandle string, options *filechooser.OpenFileOptions) (fyne.ListableURI, error) {
	title := lang.L("Open") + " " + lang.L("Folder")
	uri, err := open(parentWindowHandle, title, options)
	if err != nil || uri == nil {
		return nil, err
	}

	return storage.ListerForURI(uri)
}

func open(parentWindowHandle, title string, options *filechooser.OpenFileOptions) (fyne.URI, error) {
	uris, err := filechooser.OpenFile(parentWindowHandle, title, options)
	if err != nil {
		return nil, err
	}

	if len(uris) == 0 {
		return nil, nil
	}

	return storage.ParseURI(uris[0])
}

func saveFile(parentWindowHandle string, options *filechooser.SaveFileOptions) (fyne.URIWriteCloser, error) {
	title := lang.L("Save") + " " + lang.L("File")
	uris, err := filechooser.SaveFile(parentWindowHandle, title, options)
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

	go func() {
		if folder {
			folder, err := openFolder(windowHandle, options)
			fyne.Do(func() {
				folderCallback(folder, err)
			})
		} else {
			file, err := openFile(windowHandle, options)
			fyne.Do(func() {
				fileCallback(file, err)
			})
		}
	}()

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
		fyne.Do(func() {
			callback(file, err)
		})
	}()

	return true
}

func windowHandleForPortal(window fyne.Window) string {
	windowHandle := ""
	if !build.IsWayland {
		window.(driver.NativeWindow).RunNative(func(context any) {
			handle := context.(driver.X11WindowContext).WindowHandle
			windowHandle = portal.FormatX11WindowHandle(handle)
		})
	}

	// TODO: We need to get the Wayland handle from the xdg_foreign protocol and convert to string on the form "wayland:{id}".
	return windowHandle
}
