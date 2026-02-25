package cache

import (
	"strconv"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

var (
	overrides     async.Map[fyne.CanvasObject, *overrideScope]
	overrideCount atomic.Uint32
)

type overrideScope struct {
	th      fyne.Theme
	cacheID string
}

// OverrideTheme allows an app to specify that a single object should use a different theme to the app.
// This should be used sparingly to avoid a jarring user experience.
// If the object is a container it will theme the children, if it is a canvas primitive it will do nothing.
//
// Since: 2.5
func OverrideTheme(o fyne.CanvasObject, th fyne.Theme) {
	id := overrideCount.Add(1)
	s := &overrideScope{th: th, cacheID: strconv.Itoa(int(id))}
	overrideTheme(o, s)
}

func OverrideThemeMatchingScope(o, parent fyne.CanvasObject) bool {
	scope, ok := overrides.Load(parent)
	if !ok { // not overridden in parent
		return false
	}

	overrideTheme(o, scope)
	return true
}

func WidgetScopeID(o fyne.CanvasObject) string {
	scope, ok := overrides.Load(o)
	if !ok {
		return ""
	}

	return scope.cacheID
}

func WidgetTheme(o fyne.CanvasObject) fyne.Theme {
	scope, ok := overrides.Load(o)
	if !ok {
		return nil
	}

	return scope.th
}

func overrideContainer(c *fyne.Container, s *overrideScope) {
	for _, o := range c.Objects {
		overrideTheme(o, s)
	}
}

func overrideTheme(o fyne.CanvasObject, s *overrideScope) {
	if _, ok := o.(interface{ SetDeviceIsMobile(bool) }); ok { // ThemeOverride without the import loop
		return // do not apply this theme over a new scope
	}

	switch c := o.(type) {
	case fyne.Widget:
		overrideWidget(c, s)
	case *fyne.Container:
		overrideContainer(c, s)
	default:
		overrides.Store(c, s)
	}
}

func overrideWidget(w fyne.Widget, s *overrideScope) {
	ResetThemeCaches()
	overrides.Store(w, s)

	r := Renderer(w)
	if r == nil {
		return
	}

	for _, o := range r.Objects() {
		overrideTheme(o, s)
	}
}
