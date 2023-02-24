package scale

import (
	"errors"
	"math"

	"fyne.io/fyne/v2"
)

// ToScreenCoordinate converts a fyne coordinate in the given canvas to a screen coordinate
func ToScreenCoordinate(c fyne.Canvas, v float32) int {
	return int(math.Ceil(float64(v * c.Scale())))
}

// ToFyneCoordinate converts a screen coordinate for a given canvas to a fyne coordinate
func ToFyneCoordinate(c fyne.Canvas, v int) float32 {
	switch c.Scale() {
	case 0.0:
		panic("Incorrect scale most likely not set.")
	case 1.0:
		return float32(v)
	default:
		return float32(v) / c.Scale()
	}
}

func ToFyneSize(obj fyne.CanvasObject, width, height int) (fyne.Size, error) {
	app := fyne.CurrentApp()
	if app == nil {
		return fyne.NewSize(0, 0), errors.New("no current app")
	}
	driver := app.Driver()
	if driver == nil {
		return fyne.NewSize(0, 0), errors.New("no current driver")
	}
	c := driver.CanvasForObject(obj)
	if c == nil {
		return fyne.NewSize(0, 0), errors.New("object is not attached to a canvas yet")
	}
	dpSize := fyne.NewSize(ToFyneCoordinate(c, width), ToFyneCoordinate(c, height))

	return dpSize, nil
}
