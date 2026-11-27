//go:build wasm || js

package dialog

import (
	"strings"
	"syscall/js"

	"fyne.io/fyne/v2"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return nil
}

func isHidden(file fyne.URI) bool {
	return false
}

func fileOpenOSOverride(f *FileDialog) bool {
	if f.isDirectory() {
		showBrowserFolderOpen(f.callback.(func(fyne.ListableURI, error)))
		return true
	}

	showBrowserFileOpen(f.callback.(func(fyne.URIReadCloser, error)), f.filter)
	return true
}

// showBrowserFileOpen asks the browser to show its file picker and passes the chosen file to callback.
// This must be called soon after a user interaction or the browser will refuse to show the picker.
func showBrowserFileOpen(callback func(fyne.URIReadCloser, error), filter storage.FileFilter) {
	showBrowserPicker(func(input js.Value) {
		if accept := browserAccept(filter); accept != "" {
			input.Set("accept", accept)
		}
	}, func(files js.Value) {
		if files.IsUndefined() || files.Length() == 0 {
			fyne.Do(func() { callback(nil, nil) })
			return
		}

		// the content is read now, as the callback cannot wait for the browser. The file
		// can also be opened again later from its URI, for as long as the page is open.
		chosen := files.Index(0)
		file, err := intRepo.ReadBrowserFile(chosen, intRepo.BrowserFiles().AddFile(chosen))
		fyne.Do(func() {
			if err != nil {
				callback(nil, err)
				return
			}
			callback(file, nil)
		})
	})
}

// showBrowserFolderOpen asks the browser to show its folder picker and passes the chosen folder to callback.
// This must be called soon after a user interaction or the browser will refuse to show the picker.
func showBrowserFolderOpen(callback func(fyne.ListableURI, error)) {
	showBrowserPicker(func(input js.Value) {
		input.Set("webkitdirectory", true)
	}, func(files js.Value) {
		if files.IsUndefined() || files.Length() == 0 {
			// cancelled, or an empty folder - which the browser tells us nothing about
			fyne.Do(func() { callback(nil, nil) })
			return
		}

		folder, err := storage.ListerForURI(intRepo.BrowserFiles().AddFolder(files))
		fyne.Do(func() { callback(folder, err) })
	})
}

// showBrowserPicker shows the browser's file picker, set up by configure, and passes what was chosen
// to done - which is called in a goroutine, with an undefined value if the user cancelled.
func showBrowserPicker(configure func(input js.Value), done func(files js.Value)) {
	document := js.Global().Get("document")
	input := document.Call("createElement", "input")
	input.Set("type", "file")
	input.Get("style").Set("display", "none")
	configure(input)

	var onChange, onCancel js.Func
	cleanup := func() {
		onChange.Release()
		onCancel.Release()
		input.Call("remove")
	}

	// The handlers must not block the browser, so the work is done in goroutines.
	onChange = js.FuncOf(func(js.Value, []js.Value) any {
		files := input.Get("files")
		go func() {
			defer cleanup()
			done(files)
		}()
		return nil
	})
	onCancel = js.FuncOf(func(js.Value, []js.Value) any {
		go func() {
			defer cleanup()
			done(js.Undefined())
		}()
		return nil
	})

	input.Call("addEventListener", "change", onChange)
	input.Call("addEventListener", "cancel", onCancel)
	document.Get("body").Call("appendChild", input)
	input.Call("click")
}

// browserAccept converts a filter to the value of a file input's "accept" attribute.
func browserAccept(filter storage.FileFilter) string {
	switch f := filter.(type) {
	case nil:
		return ""
	case *storage.ExtensionFileFilter:
		return strings.Join(f.Extensions, ",")
	case *storage.MimeTypeFileFilter:
		return strings.Join(f.MimeTypes, ",")
	}

	fyne.LogError("Custom filter types not supported in the browser", nil)
	return ""
}

func fileSaveOSOverride(f *FileDialog) bool {
	// TODO #2738
	return true
}

func (f *fileDialog) getPlaces() []favoriteItem {
	return []favoriteItem{}
}

func getFavoriteLocation(homeURI fyne.URI, name string) (fyne.URI, error) {
	return nil, nil
}
