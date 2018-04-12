package examples

import "image/color"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/layout"

func Canvas(app app.App) {
	w := app.NewWindow("Main")

	content := ui.NewContainer(
		canvas.NewText("Resize me"),
		&canvas.RectangleObject{FillColor: color.RGBA{0x80, 0, 0, 0xff},
			StrokeColor: color.RGBA{0xff, 0xff, 0xff, 0xff},
			StrokeWidth: 1},
		&canvas.LineObject{StrokeColor: color.RGBA{0, 0x80, 0, 0xff}, StrokeWidth: 5},
		&canvas.CircleObject{StrokeColor: color.RGBA{0, 0, 0x80, 0xff},
			FillColor:   color.RGBA{0x30, 0x30, 0x30, 0x60},
			StrokeWidth: 2})
	content.Layout = layout.NewFixedGridLayout(ui.NewSize(93, 93))

	w.Canvas().SetContent(content)
	w.Show()
}
