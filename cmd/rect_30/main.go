package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Renlite - Flex RoundRectangle Time")
	green := color.NRGBA{R: 0, G: 180, B: 0, A: 150}
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.NRGBA{R: 0, G: 0, B: 255, A: 100}
	purple := color.NRGBA{R: 150, G: 0, B: 205, A: 255}

	// RRect1
	rr1 := canvas.Rectangle{FillColor: green, StrokeColor: purple, StrokeWidth: 40.0}
	rr1.Resize((fyne.NewSize(300, 150)))
	rr1.Move(fyne.NewPos(10, 0))
	// RRect2
	rr2 := canvas.Rectangle{FillColor: purple, StrokeColor: green, StrokeWidth: 5.0}
	rr2.Resize((fyne.NewSize(300, 150)))
	rr2.Move(fyne.NewPos(360, 0))
	// RRect3
	rr3 := canvas.Rectangle{FillColor: red, StrokeColor: green, StrokeWidth: 5.0}
	rr3.Resize((fyne.NewSize(300, 150)))
	rr3.Move(fyne.NewPos(10, 400))
	// Line1
	line1 := canvas.NewLine(blue)
	line1.Position1.X = 10.0
	line1.Position1.Y = 180.0
	line1.Position2.X = 310.0
	line1.Position2.Y = 330.0
	line1.StrokeWidth = 5.0
	// Line2
	line2 := canvas.NewLine(red)
	line2.Position1.X = 360.0
	line2.Position1.Y = 330.0
	line2.Position2.X = 660.0
	line2.Position2.Y = 180.0
	line2.StrokeWidth = 5.0

	// Widgets
	btn1 := widget.NewButton("click me", func() {
		log.Println("*** tapped ***")
	})
	btn1.Move(fyne.NewPos(700, 450))
	btn1.Resize((fyne.NewSize(140.0, 40.0)))
	txtSeg1 := &widget.TextSegment{Text: "composition"}
	txt1 := widget.NewRichText(txtSeg1)
	txt1.Move(fyne.NewPos(380, 380))
	txtSeg2 := &widget.TextSegment{Text: "one GL stream to GPU"}
	txt2 := widget.NewRichText(txtSeg2)
	txt2.Move(fyne.NewPos(70, 380))

	cont := container.NewWithoutLayout(
		&rr1,
		&rr2,
		&rr3,
		line1,
		line2,
		//btn1,
	)
	myWindow.SetContent(cont)
	myWindow.Resize(fyne.NewSize(900, 600))

	myWindow.ShowAndRun()
}
