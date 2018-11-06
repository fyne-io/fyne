package apps

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/widget"
)

// Widget shows a window containing widget demos
func Widget(app fyne.App) {
	w := app.NewWindow("TabContainer")

	tabs := widget.NewTabContainer(
		&widget.TabItem{
			Text:    "Item 2",
			Content: canvas.NewImageFromResource(theme.FyneLogo()),
		},
		&widget.TabItem{
			Text:    "Item 1",
			Content: widget.NewLabel("Label"),
		})
	w.SetContent(tabs)

	w.Show()
}
