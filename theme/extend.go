package theme

import (
	"image/color"

	"fyne.io/fyne"
)

// ExtendDefaultTheme allows a theme to specify alterations to the default theme.
// It will delegate all colour, theme or font lookups ignored to the built-in theme.
//
// Since: 2.0.0
func ExtendDefaultTheme(t fyne.Theme) fyne.Theme {
	return &themeWrapper{override: t, fallback: DefaultTheme()}
}

var _ fyne.Theme = (*themeWrapper)(nil)

type themeWrapper struct {
	override, fallback fyne.Theme
}

func (t *themeWrapper) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	c := t.override.Color(n, v)
	if c != nil {
		return c
	}

	return t.fallback.Color(n, v)
}

func (t *themeWrapper) Size(n fyne.ThemeSizeName) int {
	s := t.override.Size(n)
	if s != 0 {
		return s
	}

	return t.fallback.Size(n)
}

func (t *themeWrapper) Font(s fyne.TextStyle) fyne.Resource {
	f := t.override.Font(s)
	if f != nil {
		return f
	}

	return t.fallback.Font(s)
}
