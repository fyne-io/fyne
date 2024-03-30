package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct {
	fyne.Theme
}

func NewThemeDark() fyne.Theme {
	return &Theme{
		Theme: theme.DefaultTheme(),
	}
}

func (t *Theme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, theme.VariantDark)
}

func (t *Theme) Size(name fyne.ThemeSizeName) float32 {
	// if name == theme.SizeNamePadding {
	// 	return float32(10)
	// }
	if name == theme.SizeNameInlineIcon {
		return float32(30)
	}
	if name == theme.SizeNameInnerPadding {
		return float32(8)
	}
	if name == theme.SizeNameText {
		return float32(30)
	}
	return t.Theme.Size(name)
}
