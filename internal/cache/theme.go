package cache

import (
	"strconv"

	"fyne.io/fyne/v2"
)

var overrides = make(map[fyne.CanvasObject]*overrideScope)

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
	s := &overrideScope{th: th, cacheID: strconv.Itoa(len(overrides))}
	overrideTheme(o, s)
}

func OverrideThemeMatchingScope(o, parent fyne.CanvasObject) bool {
	scope, ok := overrides[parent]
	if !ok { // not overridden in parent
		return false
	}

	overrideTheme(o, scope)
	return true
}

func WidgetScopeID(o fyne.CanvasObject) string {
	scope, ok := overrides[o]
	if !ok {
		return ""
	}

	return scope.cacheID
}

func WidgetTheme(o fyne.CanvasObject) fyne.Theme {
	scope, ok := overrides[o]
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
	switch c := o.(type) {
	case fyne.Widget:
		overrideWidget(c, s)
	case *fyne.Container:
		overrideContainer(c, s)
	default:
		overrides[c] = s
	}
}

func overrideWidget(w fyne.Widget, s *overrideScope) {
	ResetThemeCaches()
	overrides[w] = s

	r := Renderer(w)
	if r == nil {
		return
	}

	for _, o := range r.Objects() {
		overrideTheme(o, s)
	}
}
