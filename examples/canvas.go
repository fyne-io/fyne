package examples

import "image/color"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/layout"

func rgbGradient(x, y, w, h int) color.RGBA {
	g := int(float32(x) / float32(h) * float32(255))
	b := int(float32(y) / float32(w) * float32(255))

	return color.RGBA{uint8(255 - b), uint8(g), uint8(b), 0xff}
}

func Canvas(app app.App) {
	w := app.NewWindow("Main")

	content := ui.NewContainer(
		canvas.NewText("Resize me"),
		&canvas.Rectangle{FillColor: color.RGBA{0x80, 0, 0, 0xff},
			StrokeColor: color.RGBA{0xff, 0xff, 0xff, 0xff},
			StrokeWidth: 1},
		canvas.NewRaster(rgbGradient),
		&canvas.Line{StrokeColor: color.RGBA{0, 0x80, 0, 0xff}, StrokeWidth: 5},
		&canvas.Circle{StrokeColor: color.RGBA{0, 0, 0x80, 0xff},
			FillColor:   color.RGBA{0x30, 0x30, 0x30, 0x60},
			StrokeWidth: 2})
	content.Layout = layout.NewFixedGridLayout(ui.NewSize(93, 93))

	w.Canvas().SetContent(content)
	w.Show()
}
