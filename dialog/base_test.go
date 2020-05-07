package dialog_test

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"testing"
)

func TestDialog_Background_Group(t *testing.T) {
	test.NewApp().Settings().SetTheme(theme.LightTheme())
	defer test.NewApp()
	w := test.NewWindowWithPainter(canvas.NewRectangle(theme.BackgroundColor()), software.NewPainter())
	w.Resize(fyne.NewSize(400, 300))

	group := widget.NewGroup("Foo", &widget.Button{Text: "Foo"}, layout.NewSpacer(), &widget.Button{Text: "Bar"})
	dialog := dialog.NewCustom("Foo", "Cancel", group, w)
	dialog.Show()

	test.AssertImageMatches(t, "dialog_background_group.png", w.Canvas().Capture())

	w.Close()
}
