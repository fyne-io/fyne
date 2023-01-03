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
	myWindow := myApp.NewWindow("Renlite - Vector/GPU based Rectangle and RoundRectangle")
	green := color.NRGBA{R: 0, G: 180, B: 0, A: 150}
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.NRGBA{R: 0, G: 0, B: 255, A: 100}
	orange := color.NRGBA{R: 255, G: 120, B: 0, A: 255}
	yellow := color.NRGBA{R: 255, G: 200, B: 0, A: 180}

	// RRect1
	rr1 := canvas.Rectangle{FillColor: green, StrokeColor: blue, StrokeWidth: 10.0, Radius: fyne.Radius{Pixel: 50.0}}
	rr1.Resize((fyne.NewSize(300, 150)))
	rr1.Move(fyne.NewPos(10, 0))
	// RRect2
	rr2 := canvas.Rectangle{FillColor: yellow, StrokeColor: orange, StrokeWidth: 8.0, Radius: fyne.Radius{Pixel: 25.0}}
	rr2.Resize((fyne.NewSize(300, 150)))
	rr2.Move(fyne.NewPos(360, 0))
	// RRect3
	rr3 := canvas.Rectangle{FillColor: orange, StrokeColor: blue, StrokeWidth: 15.0}
	rr3.Resize((fyne.NewSize(300, 150)))
	rr3.Move(fyne.NewPos(10, 400))
	// RRect4
	rr4 := canvas.Rectangle{FillColor: orange, StrokeColor: blue, StrokeWidth: 15.0, Radius: fyne.Radius{Pixel: 0.01}}
	rr4.Resize((fyne.NewSize(300, 150)))
	rr4.Move(fyne.NewPos(360, 400))
	// RRect5 as Circle emulation
	rr5 := canvas.Rectangle{FillColor: yellow, StrokeColor: orange, StrokeWidth: 8.0, Radius: fyne.Radius{Percent: 0.50}}
	rr5.Resize((fyne.NewSize(150, 150)))
	rr5.Move(fyne.NewPos(700, 0))
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

	cont := container.NewWithoutLayout(
		&rr1,
		&rr2,
		&rr5,
		line1,
		line2,
		&rr3,
		&rr4,
		btn1,
	)
	myWindow.SetContent(cont)
	myWindow.Resize(fyne.NewSize(1000, 600))

	myWindow.ShowAndRun()
}
