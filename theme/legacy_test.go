package theme

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

var oldTheme = &legacyTheme{}

func TestFromLegacy(t *testing.T) {
	newTheme := FromLegacy(oldTheme)
	assert.NotNil(t, newTheme)
	assert.Equal(t, oldTheme, newTheme.(*legacyWrapper).old)
}

func TestLegacyWrapper_Color(t *testing.T) {
	newTheme := FromLegacy(oldTheme)
	assert.Equal(t, oldTheme.BackgroundColor(), newTheme.Color(ColorNameBackground, VariantLight))
	assert.Equal(t, oldTheme.ShadowColor(), newTheme.Color(ColorNameShadow, VariantLight))
	assert.Equal(t, oldTheme.TextColor(), newTheme.Color(ColorNameForeground, VariantLight))
}

func TestLegacyWrapper_Font(t *testing.T) {
	newTheme := FromLegacy(oldTheme)
	assert.Equal(t, oldTheme.TextFont(), newTheme.Font(fyne.TextStyle{}))
	assert.Equal(t, oldTheme.TextBoldFont(), newTheme.Font(fyne.TextStyle{Bold: true}))
	assert.Equal(t, oldTheme.TextItalicFont(), newTheme.Font(fyne.TextStyle{Italic: true}))
	assert.Equal(t, oldTheme.TextMonospaceFont(), newTheme.Font(fyne.TextStyle{Monospace: true}))
}

func TestLegacyWrapper_Size(t *testing.T) {
	newTheme := FromLegacy(oldTheme)
	assert.Equal(t, oldTheme.IconInlineSize(), int(newTheme.Size(SizeNameInlineIcon)))
	assert.Equal(t, oldTheme.Padding(), int(newTheme.Size(SizeNamePadding)))
	assert.Equal(t, oldTheme.TextSize(), int(newTheme.Size(SizeNameText)))
}

var _ fyne.LegacyTheme = (*legacyTheme)(nil)

type legacyTheme struct{}

func (*legacyTheme) BackgroundColor() color.Color {
	return BackgroundColor()
}

func (*legacyTheme) ButtonColor() color.Color {
	return ButtonColor()
}

func (*legacyTheme) DisabledButtonColor() color.Color {
	return DisabledButtonColor()
}

func (*legacyTheme) DisabledTextColor() color.Color {
	return DisabledColor()
}

func (*legacyTheme) FocusColor() color.Color {
	return FocusColor()
}

func (*legacyTheme) HoverColor() color.Color {
	return HoverColor()
}

func (*legacyTheme) PlaceHolderColor() color.Color {
	return PlaceHolderColor()
}

func (*legacyTheme) PrimaryColor() color.Color {
	return PrimaryColor()
}

func (*legacyTheme) ScrollBarColor() color.Color {
	return ScrollBarColor()
}

func (*legacyTheme) ShadowColor() color.Color {
	return ShadowColor()
}

func (*legacyTheme) TextColor() color.Color {
	return ForegroundColor()
}

func (*legacyTheme) TextSize() int {
	return int(TextSize())
}

func (*legacyTheme) TextFont() fyne.Resource {
	return TextFont()
}

func (*legacyTheme) TextBoldFont() fyne.Resource {
	return TextBoldFont()
}

func (*legacyTheme) TextItalicFont() fyne.Resource {
	return TextItalicFont()
}

func (*legacyTheme) TextBoldItalicFont() fyne.Resource {
	return TextBoldItalicFont()
}

func (*legacyTheme) TextMonospaceFont() fyne.Resource {
	return TextMonospaceFont()
}

func (*legacyTheme) Padding() int {
	return int(Padding())
}

func (*legacyTheme) IconInlineSize() int {
	return int(IconInlineSize())
}

func (*legacyTheme) ScrollBarSize() int {
	return int(ScrollBarSize())
}

func (*legacyTheme) ScrollBarSmallSize() int {
	return int(ScrollBarSmallSize())
}
