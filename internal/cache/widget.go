package cache

import (
	"sync"

	"fyne.io/fyne"
)

var renderers sync.Map

type isBaseWidget interface {
	ExtendBaseWidget(fyne.Widget)
	super() fyne.Widget
}

// Renderer looks up the render implementation for a widget
func Renderer(wid fyne.Widget) fyne.WidgetRenderer {
	if wid == nil {
		return nil
	}

	if wd, ok := wid.(isBaseWidget); ok {
		if wd.super() != nil {
			wid = wd.super()
		}
	}
	renderer, ok := renderers.Load(wid)
	if !ok {
		renderer = wid.CreateRenderer()
		renderers.Store(wid, renderer)
	}

	if renderer == nil {
		return nil
	}
	return renderer.(fyne.WidgetRenderer)
}

// DestroyRenderer frees a render implementation for a widget.
// This is typically for internal use only.
func DestroyRenderer(wid fyne.Widget) {
	Renderer(wid).Destroy()

	renderers.Delete(wid)
}

// IsRendered returns true of the widget currently has a renderer.
// One will be created the first time a widget is shown but may be removed after it is hidden.
func IsRendered(wid fyne.Widget) bool {
	_, found := renderers.Load(wid)
	return found
}
