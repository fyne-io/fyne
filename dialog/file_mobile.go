// +build ios android

package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/gomobile"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return nil
}

func isHidden(file, _ string) bool {
	return false
}

func fileOpenOSOverride(callback func(fyne.FileReader, error), _ fyne.Window) bool {
	gomobile.ShowFileOpenPicker(callback)
	return true
}

func fileSaveOSOverride(callback func(fyne.FileWriter, error), parent fyne.Window) bool {
	ShowInformation("File Save", "File save not available on mobile", parent)

	if callback != nil {
		callback(nil, nil)
	}

	return true
}
