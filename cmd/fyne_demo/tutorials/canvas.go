package tutorials

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func rgbGradient(x, y, w, h int) color.Color {
	g := int(float32(x) / float32(w) * float32(255))
	b := int(float32(y) / float32(h) * float32(255))

	return color.NRGBA{uint8(255 - b), uint8(g), uint8(b), 0xff}
}

// canvasScreen loads a graphics example panel for the demo app
func canvasScreen(_ fyne.Window) fyne.CanvasObject {
	gradient := canvas.NewHorizontalGradient(color.NRGBA{0x80, 0, 0, 0xff}, color.NRGBA{0, 0x80, 0, 0xff})
	go func() {
		for {
			time.Sleep(time.Second)

			gradient.Angle += 45
			if gradient.Angle >= 360 {
				gradient.Angle -= 360
			}
			canvas.Refresh(gradient)
		}
	}()

	return container.NewGridWrap(fyne.NewSize(90, 90),
		canvas.NewImageFromResource(theme.FyneLogo()),
		&canvas.Rectangle{FillColor: color.NRGBA{0x80, 0, 0, 0xff},
			StrokeColor: color.NRGBA{0xff, 0xff, 0xff, 0xff},
			StrokeWidth: 1},
		&canvas.Line{StrokeColor: color.NRGBA{0, 0, 0x80, 0xff}, StrokeWidth: 5},
		&canvas.Circle{StrokeColor: color.NRGBA{0, 0, 0x80, 0xff},
			FillColor:   color.NRGBA{0x30, 0x30, 0x30, 0x60},
			StrokeWidth: 2},
		canvas.NewText("Text", color.NRGBA{0, 0x80, 0, 0xff}),
		canvas.NewRasterWithPixels(rgbGradient),
		gradient,
		canvas.NewRadialGradient(color.NRGBA{0x80, 0, 0, 0xff}, color.NRGBA{0, 0x80, 0x80, 0xff}),
	)
}
