package cache

import "fyne.io/fyne/v2"

var overrides = make(map[fyne.Widget]fyne.Theme)

func overrideWidget(w fyne.Widget, th fyne.Theme) {
	ResetThemeCaches()
	overrides[w] = th

	r := Renderer(w)
	if r == nil {
		return
	}

	for _, o := range r.Objects() {
		OverrideTheme(o, th)
	}
}

func overrideContainer(c *fyne.Container, th fyne.Theme) {
	for _, o := range c.Objects {
		OverrideTheme(o, th)
	}
}

// OverrideTheme allows an app to specify that a single object should use a different theme to the app.
// This should be used sparingly to avoid a jarring user experience.
// If the object is a container it will theme the children, if it is a canvas primitive it will do nothing.
//
// Since: 2.5
func OverrideTheme(o fyne.CanvasObject, th fyne.Theme) {
	switch c := o.(type) {
	case fyne.Widget:
		overrideWidget(c, th)
	case *fyne.Container:
		overrideContainer(c, th)
	}
}

func WidgetTheme(o fyne.Widget) fyne.Theme {
	th, ok := overrides[o]
	if !ok {
		return nil
	}

	return th
}
