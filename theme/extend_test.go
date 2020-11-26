package theme

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
)

var wrappedTheme = &extendedTheme{}

func TestExtendDefaultTheme(t *testing.T) {
	newTheme := ExtendDefaultTheme(wrappedTheme)
	assert.NotNil(t, newTheme)
	assert.Equal(t, wrappedTheme, newTheme.(*themeWrapper).override)
}

func TestExtendTheme_Color(t *testing.T) {
	newTheme := ExtendDefaultTheme(wrappedTheme)
	assert.Equal(t, overrideColor, newTheme.Color(Colors.Background, Variants.Light))
	assert.Equal(t, ShadowColor(), newTheme.Color(Colors.Shadow, Variants.Light))
	assert.Equal(t, overrideColor, newTheme.Color(Colors.Text, Variants.Light))
}

func TestExtendTheme_Font(t *testing.T) {
	newTheme := ExtendDefaultTheme(wrappedTheme)
	assert.Equal(t, overrideFont, newTheme.Font(fyne.TextStyle{Bold: true}))
	assert.Equal(t, TextItalicFont(), newTheme.Font(fyne.TextStyle{Italic: true}))
	assert.Equal(t, TextMonospaceFont(), newTheme.Font(fyne.TextStyle{Monospace: true}))
}

func TestExtendTheme_Size(t *testing.T) {
	newTheme := ExtendDefaultTheme(wrappedTheme)
	assert.Equal(t, IconInlineSize(), newTheme.Size(Sizes.InlineIcon))
	assert.Equal(t, overrideSize, newTheme.Size(Sizes.Padding))
	assert.Equal(t, TextSize(), newTheme.Size(Sizes.Text))
}

var (
	overrideColor = &color.Gray{Y: 0x66}
	overrideFont  = DefaultTextMonospaceFont()
	overrideSize  = 100
)

var _ fyne.Theme = (*extendedTheme)(nil)

type extendedTheme struct {
}

func (e *extendedTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if n == Colors.Background || n == Colors.Text {
		return overrideColor
	}

	return nil
}

func (e *extendedTheme) Size(n fyne.ThemeSizeName) int {
	if n == Sizes.Padding {
		return overrideSize
	}

	return 0
}

func (e *extendedTheme) Font(t fyne.TextStyle) fyne.Resource {
	if t.Bold {
		return overrideFont
	}

	return nil
}
