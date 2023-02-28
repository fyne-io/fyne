package scale

import (
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

// ToFyneSize returns the scaled size of an object based on pixel coordinates, typically for images.
// This method will attempt to find the canvas for an object to get its scale.
// In the event that this fails it will assume a 1:1 mapping (scale=1 or low DPI display).
func ToFyneSize(obj fyne.CanvasObject, width, height int) fyne.Size {
	app := fyne.CurrentApp()
	if app == nil {
		return fyne.NewSize(float32(width), float32(height)) // can occur if called before app.New
	}
	driver := app.Driver()
	if driver == nil {
		return fyne.NewSize(float32(width), float32(height))
	}
	c := driver.CanvasForObject(obj)
	if c == nil {
		return fyne.NewSize(float32(width), float32(height)) // this will happen a lot during init
	}
	return fyne.NewSize(ToFyneCoordinate(c, width), ToFyneCoordinate(c, height))
}
