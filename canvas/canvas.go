package canvas

import "fyne.io/fyne/v2"

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
