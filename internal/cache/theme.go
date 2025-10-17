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

	// Provide an extra fast means of accessing the override
	// for a given object multiple times in succession.
	// (e.g. repeated calls to theme.ColorForWidget(colorName, wid))
	//
	// This should be set to the last queried entry in the overrides map
	//
	// Whenever an entry is deleted from overrides during cache clean,
	// if that entry is the current one here, both of these fields
	// should be set to nil as well
	currentOverrideObj   fyne.CanvasObject
	currentOverrideScope *overrideScope
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
	scope, ok := lookupOverride(o)
	if !ok { // not overridden in parent
		return false
	}

	overrideTheme(o, scope)
	return true
}

func WidgetScopeID(o fyne.CanvasObject) string {
	scope, ok := lookupOverride(o)
	if !ok {
		return ""
	}

	return scope.cacheID
}

func WidgetTheme(o fyne.CanvasObject) fyne.Theme {
	scope, ok := lookupOverride(o)
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
		storeOverride(c, s)
	}
}

func overrideWidget(w fyne.Widget, s *overrideScope) {
	ResetThemeCaches()
	storeOverride(w, s)

	r := Renderer(w)
	if r == nil {
		return
	}

	for _, o := range r.Objects() {
		overrideTheme(o, s)
	}
}

func lookupOverride(o fyne.CanvasObject) (*overrideScope, bool) {
	if currentOverrideObj == o {
		return currentOverrideScope, currentOverrideScope != nil
	}
	s, ok := overrides.Load(o)
	currentOverrideObj = o
	currentOverrideScope = s
	return s, ok
}

func storeOverride(o fyne.CanvasObject, s *overrideScope) {
	if currentOverrideObj == o {
		currentOverrideScope = s
	}
	overrides.Store(o, s)
}
