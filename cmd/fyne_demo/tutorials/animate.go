package tutorials

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func makeAnimate(_ fyne.Window) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Black)
	rect.Resize(fyne.NewSize(200, 200))

	a := canvas.NewColorAnimation(theme.PrimaryColorNamed(theme.ColorBlue), theme.PrimaryColorNamed(theme.ColorYellow),
		time.Second*3, func(c color.Color) {
			rect.FillColor = c
			canvas.Refresh(rect)
		})
	a.Repeat = true
	a.Start()

	var a2 *fyne.Animation
	btn := widget.NewButton("Slide", func() {
		a2.Start()
	})
	btn.Resize(btn.MinSize())
	btn.Move(fyne.NewPos(20, 220))

	a2 = canvas.NewPositionAnimation(fyne.NewPos(20, 220), fyne.NewPos(320, 220), canvas.DurationStandard, func(p fyne.Position) {
		btn.Move(p)
	})

	return fyne.NewContainerWithoutLayout(rect, btn)
}
