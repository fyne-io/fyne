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
	curves := makeAnimateCurves()
	curves.Move(fyne.NewPos(0, 140+theme.Padding()))
	return fyne.NewContainerWithoutLayout(makeAnimateCanvas(), curves)
}

func makeAnimateCanvas() fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Black)
	rect.Resize(fyne.NewSize(410, 140))

	a := canvas.NewColorAnimation(theme.PrimaryColorNamed(theme.ColorBlue), theme.PrimaryColorNamed(theme.ColorGreen),
		time.Second*3, func(c color.Color) {
			rect.FillColor = c
			canvas.Refresh(rect)
		})
	a.Repeat = true
	a.Start()

	var a2 *fyne.Animation
	i := widget.NewIcon(theme.CheckButtonCheckedIcon())
	a2 = canvas.NewPositionAnimation(fyne.NewPos(0, 0), fyne.NewPos(350, 80), time.Second*3, func(p fyne.Position) {
		i.Move(p)

		width := int(10 + (float64(p.X) / 7))
		i.Resize(fyne.NewSize(width, width))
	})
	a2.Repeat = true
	a2.Curve = fyne.AnimationLinear
	a2.Start()

	return fyne.NewContainerWithoutLayout(rect, i)
}

func makeAnimateCurves() fyne.CanvasObject {
	label1 := widget.NewLabel("EaseInOut")
	label1.Alignment = fyne.TextAlignCenter
	label1.Resize(fyne.NewSize(380, 30))
	box1 := canvas.NewRectangle(theme.TextColor())
	box1.Resize(fyne.NewSize(30, 30))
	box1.Move(fyne.NewPos(0, 0))
	a1 := canvas.NewPositionAnimation(
		fyne.NewPos(0, 0), fyne.NewPos(380, 0), time.Second, func(p fyne.Position) {
			box1.Move(p)
		})

	label2 := widget.NewLabel("Linear")
	label2.Alignment = fyne.TextAlignCenter
	label2.Move(fyne.NewPos(0, 30+theme.Padding()))
	label2.Resize(fyne.NewSize(380, 30))
	box2 := canvas.NewRectangle(theme.TextColor())
	box2.Resize(fyne.NewSize(30, 30))
	box2.Move(fyne.NewPos(0, 30+theme.Padding()))

	a2 := canvas.NewPositionAnimation(
		fyne.NewPos(0, 32+theme.Padding()), fyne.NewPos(380, 32+theme.Padding()), time.Second,
		func(p fyne.Position) {
			box2.Move(p)
		})
	a2.Curve = fyne.AnimationLinear

	start := widget.NewButton("Compare", func() {
		a1.Start()
		a2.Start()
	})
	start.Resize(start.MinSize())
	start.Move(fyne.NewPos(0, 60+theme.Padding()*2))
	return fyne.NewContainerWithoutLayout(label1, label2, box1, box2, start)
}
