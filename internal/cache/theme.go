package cache

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/svg"
)

var overrides = make(map[fyne.Widget]*overrideScope)

type overrideScope struct {
	th      fyne.Theme
	cacheID string
}

var nextID = 1

// OverrideTheme allows an app to specify that a single object should use a different theme to the app.
// This should be used sparingly to avoid a jarring user experience.
// If the object is a container it will theme the children, if it is a canvas primitive it will do nothing.
//
// Since: 2.5
func OverrideTheme(o fyne.CanvasObject, th fyne.Theme) {
	s := &overrideScope{th: th, cacheID: strconv.Itoa(nextID)}
	overrideTheme(o, s, nextID)
	nextID++
}

func WidgetTheme(o fyne.Widget) fyne.Theme {
	data, ok := overrides[o]
	if !ok {
		return nil
	}

	return data.th
}

func OverrideResourceTheme(res fyne.Resource, w fyne.Widget) fyne.Resource {
	if th, ok := res.(fyne.ThemedResource); ok {
		return &WidgetResource{ThemedResource: th, Owner: w}
	}

	return res
}

func themeForResource(res fyne.Resource) fyne.Theme {
	if th, ok := res.(*WidgetResource); ok {
		if over, ok := overrides[th.Owner]; ok {
			return over.th
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
	name := res.ThemeColorName()
	if name == "" {
		name = "foreground"
	}

	th := themeForResource(res)
	return svg.Colorize(res.ThemedResource.Content(), th.Color(name, fyne.CurrentApp().Settings().ThemeVariant()))
}

func (res *WidgetResource) Name() string {
	cacheID := ""
	if over, ok := overrides[res.Owner]; ok {
		cacheID = over.cacheID
	}
	return cacheID + res.ThemedResource.Name()
}

func overrideContainer(c *fyne.Container, s *overrideScope, id int) {
	for _, o := range c.Objects {
		overrideTheme(o, s, id)
	}
}

func overrideTheme(o fyne.CanvasObject, s *overrideScope, id int) {
	switch c := o.(type) {
	case fyne.Widget:
		overrideWidget(c, s, id)
	case *fyne.Container:
		overrideContainer(c, s, id)
	}
}

func overrideWidget(w fyne.Widget, s *overrideScope, id int) {
	ResetThemeCaches()
	overrides[w] = s

	r := Renderer(w)
	if r == nil {
		return
	}

	for _, o := range r.Objects() {
		overrideTheme(o, s, id)
	}
}
