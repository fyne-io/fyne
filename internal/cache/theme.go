package cache

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async/migration"
)

var (
	overrides     = make(map[fyne.CanvasObject]*overrideScope)
	overrideCount uint64
	overridesLock migration.Mutex
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
	overridesLock.Lock()
	defer overridesLock.Unlock()

	id := overrideCount
	overrideCount++
	s := &overrideScope{th: th, cacheID: strconv.FormatUint(id, 10)}
	overrideTheme(o, s)
}

func OverrideThemeMatchingScope(o, parent fyne.CanvasObject) bool {
	overridesLock.Lock()
	defer overridesLock.Unlock()

	scope, ok := overrides[parent]
	if !ok { // not overridden in parent
		return false
	}

	overrideTheme(o, scope)
	return true
}

func WidgetScopeID(o fyne.CanvasObject) string {
	overridesLock.Lock()
	defer overridesLock.Unlock()

	scope, ok := overrides[o]
	if !ok {
		return ""
	}

	return scope.cacheID
}

func WidgetTheme(o fyne.CanvasObject) fyne.Theme {
	overridesLock.Lock()
	defer overridesLock.Unlock()

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
		overridesLock.Lock()
		overrides[c] = s
		overridesLock.Unlock()
	}
}

func overrideWidget(w fyne.Widget, s *overrideScope) {
	ResetThemeCaches()

	overridesLock.Lock()
	overrides[w] = s
	overridesLock.Unlock()

	r := Renderer(w)
	if r == nil {
		return
	}

	for _, o := range r.Objects() {
		overrideTheme(o, s)
	}
}
