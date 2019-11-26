package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// ApplyThemeTo ensures that the specified canvasobject and all widgets and themeable objects will
// be updated for the current theme.
func ApplyThemeTo(content fyne.CanvasObject, canv fyne.Canvas) {
	if wid, ok := content.(fyne.Widget); ok {
		for _, o := range widget.Renderer(wid).Objects() {
			ApplyThemeTo(o, canv)
		}
	}
	if c, ok := content.(*fyne.Container); ok {
		for _, o := range c.Objects {
			ApplyThemeTo(o, canv)
		}
	}

	content.Refresh()
}

// ApplySettings ensures that all widgets and themeable objects in an application will be updated for the current theme.
// It also checks that scale changes are reflected if required
func ApplySettings(set fyne.Settings, app fyne.App) {
	for _, window := range app.Driver().AllWindows() {
		ApplyThemeTo(window.Content(), window.Canvas())
		window.Canvas().SetScale(set.Scale())

		if window.Canvas().Overlay() != nil {
			ApplyThemeTo(window.Canvas().Overlay(), window.Canvas())
		}
	}
}
