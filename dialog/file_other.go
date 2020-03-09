// +build !windows

package dialog

import (
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return []fyne.CanvasObject{widget.NewButton("Computer", func() {
		f.setDirectory("/")
	})}
}

func isHidden(file os.FileInfo, dir string) bool {
	return len(file.Name()) == 0 || file.Name()[0] == '.'
}
