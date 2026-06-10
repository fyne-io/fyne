package tutorials

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
)

func rgbGradient(x, y, w, h int) color.Color {
	g := int(float32(x) / float32(w) * float32(255))
	b := int(float32(y) / float32(h) * float32(255))

	return color.NRGBA{uint8(255 - b), uint8(g), uint8(b), 0xff}
}

// canvasScreen loads a graphics example panel for the demo app
func canvasScreen(_ fyne.Window) fyne.CanvasObject {
	gradient := canvas.NewHorizontalGradient(color.NRGBA{0x80, 0, 0, 0xff}, color.NRGBA{0, 0x80, 0, 0xff})
	ticker := time.NewTicker(time.Second)

	OnChangeFuncs = append(OnChangeFuncs, ticker.Stop)

	go func() {
		for range ticker.C {
			fyne.Do(func() {
				gradient.Angle += 45
				if gradient.Angle >= 360 {
					gradient.Angle -= 360
				}
				canvas.Refresh(gradient)
			})
		}
	}()

	return container.NewGridWrap(
		fyne.NewSize(90, 90),
		canvas.NewImageFromResource(data.FyneLogo),
		&canvas.Rectangle{
			FillColor:   color.NRGBA{0x80, 0, 0, 0xff},
			StrokeColor: color.NRGBA{R: 255, G: 120, B: 0, A: 255},
			StrokeWidth: 1,
		},
		&canvas.Rectangle{
			FillColor:    color.NRGBA{R: 255, G: 200, B: 0, A: 180},
			StrokeColor:  color.NRGBA{R: 255, G: 120, B: 0, A: 255},
			StrokeWidth:  4.0,
			CornerRadius: 20,
		},
		&canvas.RegularPolygon{
			FillColor: color.NRGBA{B: 0x80, A: 0xff}, Sides: 3, Angle: -30,
			StrokeColor: color.NRGBA{R: 0xff, G: 0x80, A: 0xff}, StrokeWidth: 2,
		},
		&canvas.ArbitraryPolygon{
			FillColor:   color.NRGBA{G: 0x80, A: 0xff},
			StrokeColor: color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xff},
			StrokeWidth: 4,
			Points: []fyne.Position{
				fyne.NewPos(4, 4),
				fyne.NewPos(4, 84),
				fyne.NewPos(84, 84),
			},
		},
		&canvas.ArbitraryPolygon{
			FillColor: color.NRGBA{R: 0x80, A: 0xff},
			Points: []fyne.Position{
				fyne.NewPos(58, 4),
				fyne.NewPos(86, 32),
				fyne.NewPos(44, 74),
				fyne.NewPos(2, 32),
				fyne.NewPos(30, 4),
				fyne.NewPos(44, 32),
			},
			CornerRadii: []float32{20, 20, 0, 20, 20, 0},
		},
		&canvas.Line{StrokeColor: color.NRGBA{0, 0, 0x80, 0xff}, StrokeWidth: 5},
		&canvas.Circle{
			StrokeColor: color.NRGBA{0, 0, 0x80, 0xff},
			FillColor:   color.NRGBA{0x30, 0x30, 0x30, 0x60},
			StrokeWidth: 2,
		},
		&canvas.Arc{
			StrokeColor:  color.NRGBA{R: 0x80, A: 0xff},
			FillColor:    color.NRGBA{G: 0x80, A: 0xff},
			StrokeWidth:  4,
			StartAngle:   -240,
			EndAngle:     60,
			CornerRadius: 6,
			CutoutRatio:  0.2,
		},
		canvas.NewText("Text", color.NRGBA{0, 0x80, 0, 0xff}),
		canvas.NewRasterWithPixels(rgbGradient),
		gradient,
		canvas.NewRadialGradient(color.NRGBA{0x80, 0, 0, 0xff}, color.NRGBA{0, 0x80, 0x80, 0xff}),
	)
}
