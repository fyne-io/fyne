package software

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/app"
)

// RenderCanvas takes a canvas and renders it to a regular Go image using the provided Theme.
// This is the same as setting the application theme and then calling Canvas.Capture().
func RenderCanvas(c fyne.Canvas, t fyne.Theme) image.Image {
	fyne.CurrentApp().Settings().SetTheme(t)
	app.ApplyThemeTo(c.Content(), c)

	return c.Capture()
}

// Render takes a canvas object and renders it to a regular Go image using the provided Theme.
// The returned image will be set to the object's minimum size.
// Use the theme.LightTheme() or theme.DarkTheme() to access the builtin themes.
func Render(obj fyne.CanvasObject, t fyne.Theme) image.Image {
	fyne.CurrentApp().Settings().SetTheme(t)

	c := NewCanvas()
	c.SetPadded(false)
	c.SetContent(obj)

	app.ApplyThemeTo(obj, c)
	return c.Capture()
}
