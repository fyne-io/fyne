package cache

import (
	"strconv"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/svg"
)

var (
	overrides     = &sync.Map{} // map[fyne.Widget]*overrideScope
	overrideCount = atomic.Uint32{}
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
	overrideTheme(o, s, id)
}

func WidgetTheme(o fyne.Widget) fyne.Theme {
	data, ok := overrides.Load(o)
	if !ok {
		return nil
	}

	return data.(*overrideScope).th
}

func OverrideResourceTheme(res fyne.Resource, w fyne.Widget) fyne.Resource {
	if th, ok := res.(fyne.ThemedResource); ok {
		return &WidgetResource{ThemedResource: th, Owner: w}
	}

	return res
}

func themeForResource(res fyne.Resource) fyne.Theme {
	if th, ok := res.(*WidgetResource); ok {
		if over, ok := overrides.Load(th.Owner); ok {
			return over.(*overrideScope).th
		}
	}

	return fyne.CurrentApp().Settings().Theme()
}

type WidgetResource struct {
	fyne.ThemedResource
	Owner fyne.Widget
}

// Content returns the underlying content of the resource adapted to the current text color.
func (res *WidgetResource) Content() []byte {
	th := themeForResource(res)
	return svg.Colorize(res.ThemedResource.Content(), th.Color(res.ThemeColorName(), fyne.CurrentApp().Settings().ThemeVariant()))
}

func (res *WidgetResource) Name() string {
	cacheID := ""
	if over, ok := overrides.Load(res.Owner); ok {
		cacheID = over.(*overrideScope).cacheID
	}
	return cacheID + res.ThemedResource.Name()
}

func overrideContainer(c *fyne.Container, s *overrideScope, id uint32) {
	for _, o := range c.Objects {
		overrideTheme(o, s, id)
	}
}

func overrideTheme(o fyne.CanvasObject, s *overrideScope, id uint32) {
	switch c := o.(type) {
	case fyne.Widget:
		overrideWidget(c, s, id)
	case *fyne.Container:
		overrideContainer(c, s, id)
	}
}

func overrideWidget(w fyne.Widget, s *overrideScope, id uint32) {
	ResetThemeCaches()
	overrides.Store(w, s)

	r := Renderer(w)
	if r == nil {
		return
	}

	for _, o := range r.Objects() {
		overrideTheme(o, s, id)
	}
}
