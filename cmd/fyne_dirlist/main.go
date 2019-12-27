package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dataapi"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func welcomeScreen() fyne.CanvasObject {
	fl := NewFileList().Source(dataapi.NewDirectoryDataSource("."))
	scroller := widget.NewScrollContainer(fl)
	top := widget.NewLabelWithStyle(
		"Fyne File List DataSource Demo",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	w := fyne.NewContainerWithLayout(layout.NewBorderLayout(top, nil, nil, nil),
		top,
		scroller,
	)
	return w
}

func main() {
	a := app.NewWithID("io.fyne.file_source")
	w := a.NewWindow("Fyne File Source Demo")
	w.Resize(fyne.NewSize(800, 800))
	w.SetMaster()
	w.SetContent(welcomeScreen())
	w.ShowAndRun()
}
