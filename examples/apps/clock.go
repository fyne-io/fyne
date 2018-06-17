package apps

import "math"
import "time"

import "github.com/fyne-io/fyne/api"
import "github.com/fyne-io/fyne/api/ui"
import "github.com/fyne-io/fyne/api/ui/canvas"
import "github.com/fyne-io/fyne/api/ui/theme"

type clockLayout struct {
	hour, minute, second     *canvas.Line
	hourdot, seconddot, face *canvas.Circle

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
	// TODO scale width to clock face size
	c.hourdot = &canvas.Circle{StrokeColor: theme.TextColor(), StrokeWidth: 5}
	c.seconddot = &canvas.Circle{StrokeColor: theme.PrimaryColor(), StrokeWidth: 3}
	c.face = &canvas.Circle{StrokeColor: theme.TextColor(), StrokeWidth: 1}

	c.hour = &canvas.Line{StrokeColor: theme.TextColor(), StrokeWidth: 5}
	c.minute = &canvas.Line{StrokeColor: theme.TextColor(), StrokeWidth: 3}
	c.second = &canvas.Line{StrokeColor: theme.PrimaryColor(), StrokeWidth: 1}

	container := ui.NewContainer(c.hourdot, c.seconddot, c.face, c.hour, c.minute, c.second)
	container.Layout = c

	c.canvas = container
	return container
}

func (c *clockLayout) animate(canvas ui.Canvas) {
	tick := time.NewTicker(time.Second)
	go func() {
		for {
			<-tick.C
			c.Layout(nil, canvas.Size())
			canvas.Refresh(c.canvas)
		}
	}()
}

func (c *clockLayout) applyTheme(setting fyne.Settings) {
	c.hour.StrokeColor = theme.TextColor()
	c.minute.StrokeColor = theme.TextColor()
	c.second.StrokeColor = theme.PrimaryColor()
}

// Clock loads a clock example window for the specified app context
func Clock(app fyne.App) {
	clockWindow := app.NewWindow("Clock")
	clock := &clockLayout{}

	canvas := clock.render()
	go clock.animate(clockWindow.Canvas())

	listener := make(chan fyne.Settings)
	fyne.GetSettings().AddChangeListener(listener)
	go func() {
		for {
			settings := <-listener
			clock.applyTheme(settings)
		}
	}()

	clockWindow.Canvas().SetContent(canvas)
	clockWindow.Show()
}
