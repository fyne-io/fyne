package main

import (
	"image/color"

	"fyne.io/fyne/v2"
)

type forceVariantTheme struct {
	fyne.Theme

	variant fyne.ThemeVariant
}

func (f *forceVariantTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.variant)
}
