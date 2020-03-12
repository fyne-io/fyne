// +build ios android

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

func fileOSOverride(save bool, callback func(string), parent fyne.Window) bool {
	if save {
		ShowInformation("File Save", "File save not available on mobile", parent)
	} else {
		ShowInformation("File Open", "File open not available on mobile", parent)
	}

	if callback != nil {
		callback("")
	}
	return true
}
