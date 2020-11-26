package tutorials

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var (
	purple = &color.NRGBA{R: 128, G: 0, B: 128, A: 255}
	orange = &color.NRGBA{R: 198, G: 123, B: 0, A: 255}
	grey   = &color.Gray{Y: 123}
)

// customTheme is a simple demonstration of a bespoke theme loaded by a Fyne app.
type customTheme struct {
}

func (customTheme) Color(c fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch c {
	case theme.Colors.Background:
		return purple
	case theme.Colors.Button, theme.Colors.DisabledText:
		return color.Black
	case theme.Colors.PlaceHolder, theme.Colors.ScrollBar:
		return grey
	case theme.Colors.Primary, theme.Colors.Hover, theme.Colors.Focus:
		return orange
	case theme.Colors.Shadow:
		return &color.RGBA{R: 0xcc, G: 0xcc, B: 0xcc, A: 0xcc}
	default:
		return color.White
	}
}

func (customTheme) Size(s fyne.ThemeSizeName) int {
	switch s {
	case theme.Sizes.Padding:
		return 10
	case theme.Sizes.InlineIcon:
		return 20
	case theme.Sizes.ScrollBar:
		return 10
	case theme.Sizes.ScrollBarSmall:
		return 5
	case theme.Sizes.Text:
		return 12
	default:
		return 0
	}
}

func (customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

func newCustomTheme() fyne.Theme {
	return &customTheme{}
}
