//go:build freebsd || linux || openbsd

package app

import "fyne.io/fyne/v2/internal/driver/mobile/event/size"

// TODO: do this for all build targets, not just linux (x11 and Android)? If
// so, should package gl instead of this package call RegisterFilter??
//
// TODO: does Android need this?? It seems to work without it (Nexus 7,
// KitKat). If only x11 needs this, should we move this to x11.go??
func (a *app) registerGLViewportFilter() {
	a.RegisterFilter(func(e any) any {
		if e, ok := e.(size.Event); ok {
			a.glctx.Viewport(0, 0, e.WidthPx, e.HeightPx)
		}
		return e
	})
}
