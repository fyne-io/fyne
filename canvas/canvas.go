package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/svg"
)

// Refresh instructs the containing canvas to refresh the specified obj.
func Refresh(obj fyne.CanvasObject) {
	app := fyne.CurrentApp()
	if app == nil || app.Driver() == nil {
		return
	}

	c := app.Driver().CanvasForObject(obj)
	if c != nil {
		c.Refresh(obj)
	}
}

// RecolorSVG takes a []byte containing SVG content, and returns
// new SVG content, re-colorized to be monochrome with the given color.
// The content can be assigned to a new fyne.StaticResource with an appropriate name
// to be used in a widget.Button, canvas.Image, etc.
//
// If an error occurs, the returned content will be the original un-modified content,
// and a non-nil error is returned.
//
// Since: 2.6
func RecolorSVG(svgContent []byte, color color.Color) ([]byte, error) {
	return svg.Colorize(svgContent, color)
}

// repaint instructs the containing canvas to redraw, even if nothing changed.
func repaint(obj fyne.CanvasObject) {
	app := fyne.CurrentApp()
	if app == nil || app.Driver() == nil {
		return
	}

	c := app.Driver().CanvasForObject(obj)
	if c != nil {
		if paint, ok := c.(interface{ SetDirty() }); ok {
			paint.SetDirty()
		}
	}
}
