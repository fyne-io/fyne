// +build !windows,!android,!ios

package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return []fyne.CanvasObject{widget.NewButton("Computer", func() {
		f.setDirectory("/")
	})}
}

func isHidden(file, _ string) bool {
	return len(file) == 0 || file[0] == '.'
}

func fileOSOverride(bool, func(string), fyne.Window) bool {
	return false
}
