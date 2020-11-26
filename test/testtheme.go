package test

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var (
	red   = &color.RGBA{R: 200, G: 0, B: 0, A: 255}
	green = &color.RGBA{R: 0, G: 255, B: 0, A: 255}
	blue  = &color.RGBA{R: 0, G: 0, B: 255, A: 255}
)

const testTextSize = 18

var _ fyne.Theme = testTheme{}

// testTheme is a simple theme variation used for testing that widgets adapt correctly
type testTheme struct {
}

// NewTheme returns a new testTheme.
func NewTheme() fyne.Theme {
	return &testTheme{}
}

func (testTheme) Color(c fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch c {
	case theme.Colors.Background:
		return red
	case theme.Colors.Button, theme.Colors.DisabledText:
		return color.Black
	case theme.Colors.PlaceHolder, theme.Colors.ScrollBar, theme.Colors.Shadow:
		return blue
	case theme.Colors.Primary, theme.Colors.Hover, theme.Colors.Focus:
		return green
	default:
		return color.White
	}
}

func (testTheme) Size(s fyne.ThemeSizeName) int {
	switch s {
	case theme.Sizes.Padding:
		return 10
	case theme.Sizes.InlineIcon:
		return 24
	case theme.Sizes.ScrollBar:
		return 10
	case theme.Sizes.ScrollBarSmall:
		return 2
	case theme.Sizes.Text:
		return testTextSize
	default:
		return 0
	}
}

func (testTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return theme.DefaultTextFont()
	}

	if style.Bold {
		if style.Italic {
			return theme.DefaultTextMonospaceFont()
		}
		return theme.DefaultTextItalicFont()
	}
	if style.Italic {
		return theme.DefaultTextBoldItalicFont()
	}

	return theme.DefaultTextBoldFont()
}
