package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Lines")

	c := container.NewWithoutLayout()
	lps := canvas.NewText("- lps", color.White)
	sig := canvas.NewRectangle(theme.BackgroundColor())
	myWindow.SetContent(container.NewBorder(lps, nil, nil, nil, container.NewMax(sig, c)))
	myWindow.Resize(fyne.NewSize(400, 400))
	myWindow.SetPadded(false)

	drawn := 0
	r := uint8(0)
	g := uint8(0)
	go func() {
		pos := fyne.NewPos(0, 0)
		for {
			newPos := fyne.NewPos(float32(rand.Intn(400)), float32(rand.Intn(400)))

			line := canvas.NewLine(color.NRGBA{R: r, G: g, B: 250, A: 255})
			g++
			if g >= 255 {
				g = 0
			}
			if g%10 == 0 {
				r++
				if r >= 255 {
					r = 0
				}
			}
			line.Position1 = pos
			line.Position2 = newPos

			c.Objects = append(c.Objects, line)
			if len(c.Objects) > 1024 {
				c.Objects = c.Objects[1:1025]
			}
			pos = newPos

			canvas.Refresh(sig)
			drawn++
			time.Sleep(time.Millisecond / 8)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			lps.Text = fmt.Sprintf("%d lps", drawn)
			lps.Refresh()
			drawn = 0
		}
	}()

	myWindow.ShowAndRun()
}
