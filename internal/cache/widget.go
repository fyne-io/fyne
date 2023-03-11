package cache

import (
	"sync"

	"fyne.io/fyne/v2"
)

var renderersLock sync.RWMutex
var renderers = map[fyne.Widget]*rendererInfo{}

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

	renderersLock.RLock()
	rinfo, ok := renderers[wid]
	renderersLock.RUnlock()
	if !ok {
		rinfo = &rendererInfo{renderer: wid.CreateRenderer()}
		renderersLock.Lock()
		renderers[wid] = rinfo
		renderersLock.Unlock()
	}

	if rinfo == nil {
		return nil
	}

	rinfo.setAlive()

	return rinfo.renderer
}

// DestroyRenderer frees a render implementation for a widget.
// This is typically for internal use only.
func DestroyRenderer(wid fyne.Widget) {
	renderersLock.RLock()
	rinfo, ok := renderers[wid]
	renderersLock.RUnlock()
	if !ok {
		return
	}
	if rinfo != nil {
		rinfo.renderer.Destroy()
	}
	renderersLock.Lock()
	delete(renderers, wid)
	renderersLock.Unlock()
}

// IsRendered returns true of the widget currently has a renderer.
// One will be created the first time a widget is shown but may be removed after it is hidden.
func IsRendered(wid fyne.Widget) bool {
	renderersLock.RLock()
	_, found := renderers[wid]
	renderersLock.RUnlock()
	return found
}

type rendererInfo struct {
	expiringCache
	renderer fyne.WidgetRenderer
}
