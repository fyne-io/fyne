package cache

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async/migration"
)

var (
	renderers     = make(map[fyne.Widget]*rendererInfo)
	renderersLock migration.Mutex
)

type isBaseWidget interface {
	ExtendBaseWidget(fyne.Widget)
	super() fyne.Widget
}

// Renderer looks up the render implementation for a widget
// If one does not exist, it creates and caches a renderer for the widget.
func Renderer(wid fyne.Widget) fyne.WidgetRenderer {
	renderer, ok := CachedRenderer(wid)
	if !ok && wid != nil {
		renderer = wid.CreateRenderer()
		rinfo := &rendererInfo{renderer: renderer}
		rinfo.setAlive()

		renderersLock.Lock()
		renderers[wid] = rinfo
		renderersLock.Unlock()
	}

	return renderer
}

// CachedRenderer looks up the cached render implementation for a widget
// If a renderer does not exist in the cache, it returns nil, false.
func CachedRenderer(wid fyne.Widget) (fyne.WidgetRenderer, bool) {
	if wid == nil {
		return nil, false
	}

	if wd, ok := wid.(isBaseWidget); ok {
		if wd.super() != nil {
			wid = wd.super()
		}
	}

	renderersLock.Lock()
	defer renderersLock.Unlock()

	rinfo, ok := renderers[wid]
	if !ok {
		return nil, false
	}

	rinfo.setAlive()
	return rinfo.renderer, true
}

// DestroyRenderer frees a render implementation for a widget.
// This is typically for internal use only.
func DestroyRenderer(wid fyne.Widget) {
	renderersLock.Lock()
	defer renderersLock.Unlock()

	rinfo, ok := renderers[wid]
	if !ok {
		return
	}

	delete(renderers, wid)
	if rinfo != nil {
		rinfo.renderer.Destroy()
	}

	overridesLock.Lock()
	delete(overrides, wid)
	overridesLock.Unlock()
}

// IsRendered returns true of the widget currently has a renderer.
// One will be created the first time a widget is shown but may be removed after it is hidden.
func IsRendered(wid fyne.Widget) bool {
	renderersLock.Lock()
	defer renderersLock.Unlock()

	_, found := renderers[wid]
	return found
}

type rendererInfo struct {
	expiringCache
	renderer fyne.WidgetRenderer
}
