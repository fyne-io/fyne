package cache

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var renderers async.Map[fyne.Widget, *rendererInfo]

type isBaseWidget interface {
	ExtendBaseWidget(fyne.Widget)
	super() fyne.Widget
}

// Renderer looks up the render implementation for a widget
// If one does not exist, it creates and caches a renderer for the widget.
func Renderer(wid fyne.Widget) fyne.WidgetRenderer {
	r, _ := rendererEntry(wid)
	return r
}

// RendererRef holds a widget's cached renderer between walks.
type RendererRef struct {
	info *rendererInfo
}

// RendererFor returns the renderer for wid, skipping the cache lookup while
// ref holds one that is not destroyed.
func RendererFor(ref *RendererRef, wid fyne.Widget) fyne.WidgetRenderer {
	if info := ref.info; info != nil && !info.destroyed.Load() {
		info.setAlive()
		return info.renderer
	}
	r, info := rendererEntry(wid)
	ref.info = info
	return r
}

func rendererEntry(wid fyne.Widget) (fyne.WidgetRenderer, *rendererInfo) {
	if rinfo := cachedRendererInfo(wid); rinfo != nil {
		return rinfo.renderer, rinfo
	}
	if wid == nil {
		return nil, nil
	}
	rinfo := &rendererInfo{renderer: wid.CreateRenderer()}
	rinfo.setAlive()
	renderers.Store(wid, rinfo)
	return rinfo.renderer, rinfo
}

// CachedRenderer looks up the cached render implementation for a widget
// If a renderer does not exist in the cache, it returns nil, false.
func CachedRenderer(wid fyne.Widget) (fyne.WidgetRenderer, bool) {
	rinfo := cachedRendererInfo(wid)
	if rinfo == nil {
		return nil, false
	}
	return rinfo.renderer, true
}

func cachedRendererInfo(wid fyne.Widget) *rendererInfo {
	if wid == nil {
		return nil
	}

	if wd, ok := wid.(isBaseWidget); ok {
		if wd.super() != nil {
			wid = wd.super()
		}
	}

	rinfo, ok := renderers.Load(wid)
	if !ok {
		return nil
	}

	rinfo.setAlive()
	return rinfo
}

// DestroyRenderer frees a render implementation for a widget.
// This is typically for internal use only.
func DestroyRenderer(wid fyne.Widget) {
	rinfo, ok := renderers.LoadAndDelete(wid)
	if !ok {
		return
	}
	if rinfo != nil {
		rinfo.destroyed.Store(true)
		rinfo.renderer.Destroy()
	}
	overrides.Delete(wid)
}

// IsRendered returns true of the widget currently has a renderer.
// One will be created the first time a widget is shown but may be removed after it is hidden.
func IsRendered(wid fyne.Widget) bool {
	_, found := renderers.Load(wid)
	return found
}

type rendererInfo struct {
	expiringCache
	renderer  fyne.WidgetRenderer
	destroyed atomic.Bool
}
