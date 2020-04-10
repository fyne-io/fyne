// +build wasm js

package dialog


import (
        "fyne.io/fyne"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
        return nil
}

func isHidden(file, _ string) bool {
        return false
}

func fileOpenOSOverride(func(fyne.FileReadCloser, error), fyne.Window) bool {
        return true
}

func fileSaveOSOverride(func(fyne.FileWriteCloser, error), fyne.Window) bool {
        return true
}
