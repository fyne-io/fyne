//go:build wasm || js

package dialog

import (
	"bytes"
	"errors"
	"strings"
	"syscall/js"

	"fyne.io/fyne/v2"
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
		// TODO folders cannot be chosen in a browser yet
		return true
	}

	showBrowserFileOpen(f.callback.(func(fyne.URIReadCloser, error)), f.filter)
	return true
}

// browserFile is a file the user chose with the browser's file picker.
// The browser does not tell us where it lives, so its content is held in memory.
type browserFile struct {
	*bytes.Reader
	uri fyne.URI
}

func (f *browserFile) Close() error {
	return nil
}

func (f *browserFile) URI() fyne.URI {
	return f.uri
}

// browserURI names a file chosen in the browser, reporting the type the browser gave it.
type browserURI struct {
	fyne.URI
	mime string
}

func (u *browserURI) MimeType() string {
	return u.mime
}

// showBrowserFileOpen asks the browser to show its file picker and passes the chosen file to callback.
// This must be called soon after a user interaction or the browser will refuse to show the picker.
func showBrowserFileOpen(callback func(fyne.URIReadCloser, error), filter storage.FileFilter) {
	document := js.Global().Get("document")
	input := document.Call("createElement", "input")
	input.Set("type", "file")
	input.Get("style").Set("display", "none")
	if accept := browserAccept(filter); accept != "" {
		input.Set("accept", accept)
	}

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

			if files.Length() == 0 {
				fyne.Do(func() { callback(nil, nil) })
				return
			}

			file, err := readBrowserFile(files.Index(0))
			fyne.Do(func() {
				if err != nil {
					callback(nil, err)
					return
				}
				callback(file, nil)
			})
		}()
		return nil
	})
	onCancel = js.FuncOf(func(js.Value, []js.Value) any {
		go func() {
			cleanup()
			fyne.Do(func() { callback(nil, nil) })
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

// readBrowserFile loads the content of a JavaScript File. It waits for the browser and so must
// not be called from an event handler.
func readBrowserFile(file js.Value) (fyne.URIReadCloser, error) {
	loaded := make(chan js.Value, 1)
	failed := make(chan error, 1)
	onLoad := js.FuncOf(func(_ js.Value, args []js.Value) any {
		loaded <- args[0]
		return nil
	})
	defer onLoad.Release()
	onFail := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) > 0 && args[0].Truthy() {
			failed <- js.Error{Value: args[0]}
		} else {
			failed <- errors.New("failed to read the chosen file")
		}
		return nil
	})
	defer onFail.Release()
	file.Call("arrayBuffer").Call("then", onLoad, onFail)

	var buffer js.Value
	select {
	case buffer = <-loaded:
	case err := <-failed:
		return nil, err
	}

	array := js.Global().Get("Uint8Array").New(buffer)
	data := make([]byte, array.Length())
	js.CopyBytesToGo(data, array)

	var uri fyne.URI = storage.NewFileURI("/" + file.Get("name").String())
	if mime := file.Get("type").String(); mime != "" {
		uri = &browserURI{URI: uri, mime: mime}
	}
	return &browserFile{Reader: bytes.NewReader(data), uri: uri}, nil
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
