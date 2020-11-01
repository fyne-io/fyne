package tutorials

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

func makeAnimate(_ fyne.Window) fyne.CanvasObject {
	r1 := canvas.NewRectangle(color.Black)
	r1.Resize(fyne.NewSize(200, 200))

	a := canvas.NewColorAnimator(theme.PrimaryColorNamed(theme.ColorBlue), theme.PrimaryColorNamed(theme.ColorYellow),
		time.Second*3, func(c color.Color) {
			r1.FillColor = c
			canvas.Refresh(r1)
		})
	a.Repeat = true
	a.Start()


	return fyne.NewContainerWithoutLayout(r1)
}
