package main

import "math"
import "time"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/theme"
import "github.com/fyne-io/fyne-app"

var window ui.Window

type clockLayout struct {
	hour, minute, second     *canvas.LineObject
	hourdot, seconddot, face *canvas.CircleObject

	canvas ui.CanvasObject
}

func (c *clockLayout) rotate(hand ui.CanvasObject, middle ui.Position, facePosition float64, offset, length int) {
	rotation := math.Pi * 2 / 60 * float64(facePosition)
	x2 := int(float64(length) * math.Sin(rotation))
	y2 := int(float64(-length) * math.Cos(rotation))

	offX := 0
	offY := 0
	if offset > 0 {
		offX += int(float64(offset) * math.Sin(rotation))
		offY += int(float64(-offset) * math.Cos(rotation))
	}

	hand.Move(ui.NewPos(middle.X+offX, middle.Y+offY))
	hand.Resize(ui.NewSize(x2, y2))
}

func (c *clockLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	diameter := ui.Min(size.Width, size.Height)
	radius := diameter / 2
	dotRadius := radius / 12
	smallDotRadius := dotRadius / 8

	size = ui.NewSize(diameter, diameter)
	middle := ui.NewPos(size.Width/2, size.Height/2)
	topleft := ui.NewPos(middle.X-radius, middle.Y-radius)

	c.face.Resize(size)
	c.face.Move(topleft)

	c.rotate(c.hour, middle, float64((time.Now().Hour()%12)*5)+(float64(time.Now().Minute())/12), dotRadius, radius/2)
	c.rotate(c.minute, middle, float64(time.Now().Minute())+(float64(time.Now().Second())/60), dotRadius, int(float64(radius)*.9))
	c.rotate(c.second, middle, float64(time.Now().Second()), 0, radius-3)

	c.hourdot.Resize(ui.NewSize(dotRadius*2, dotRadius*2))
	c.hourdot.Move(ui.NewPos(middle.X-dotRadius, middle.Y-dotRadius))
	c.seconddot.Resize(ui.NewSize(smallDotRadius*2, smallDotRadius*2))
	c.seconddot.Move(ui.NewPos(middle.X-smallDotRadius, middle.Y-smallDotRadius))
}

func (c *clockLayout) MinSize(objects []ui.CanvasObject) ui.Size {
	return ui.NewSize(200, 200)
}

func (c *clockLayout) render() *ui.Container {
	c.hourdot = &canvas.CircleObject{Color: theme.TextColor(), Width: 5}
	c.seconddot = &canvas.CircleObject{Color: theme.PrimaryColor(), Width: 3}

	c.face = &canvas.CircleObject{Color: theme.TextColor(), Width: 1}
	c.hour = &canvas.LineObject{Color: theme.TextColor(), Width: 5}
	c.minute = &canvas.LineObject{Color: theme.TextColor(), Width: 3}
	c.second = &canvas.LineObject{Color: theme.PrimaryColor(), Width: 1}

	container := ui.NewContainer(c.hourdot, c.seconddot, c.face, c.hour, c.minute, c.second)
	container.Layout = c

	c.canvas = container
	return container
}

func (c *clockLayout) animate() {
	tick := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-tick.C:
				c.Layout(nil, window.Canvas().Size())
				window.Canvas().Refresh(c.canvas)
			}
		}
	}()
}

func main() {
	app := fyneapp.NewApp()
	window = app.NewWindow("Clock")
	clock := &clockLayout{}

	canvas := clock.render()
	go clock.animate()

	window.Canvas().SetContent(canvas)
	window.Show()
}
