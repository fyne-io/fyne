package test

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func DarkTheme() fyne.Theme {
	return &forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark}
}

func LightTheme() fyne.Theme {
	return &forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight}
}

type forcedVariant struct {
	fyne.Theme

	variant fyne.ThemeVariant
}

func (f *forcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.variant)
}
