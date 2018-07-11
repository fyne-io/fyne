package apps

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/widget"

// Layout loads a window that shows the layouts available for a container
func Layout(app fyne.App) {
	w := app.NewWindow("Layout")

	top := widget.NewEntry()
	bottom := widget.NewEntry()
	left := widget.NewEntry()
	right := widget.NewEntry()
	middle := widget.NewLabel("BorderLayout")
	middle.Alignment = fyne.TextAlignCenter

	borderLayout := layout.NewBorderLayout(top, bottom, left, right)
	container := fyne.NewContainerWithLayout(borderLayout,
		top, bottom, left, right, middle)
	w.SetContent(container)

	w.Show()
}
