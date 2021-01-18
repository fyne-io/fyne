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

type legacyTheme struct {
}

func (t *legacyTheme) BackgroundColor() color.Color {
	return BackgroundColor()
}

func (t *legacyTheme) ButtonColor() color.Color {
	return ButtonColor()
}

func (t *legacyTheme) DisabledButtonColor() color.Color {
	return DisabledButtonColor()
}

func (t *legacyTheme) DisabledTextColor() color.Color {
	return DisabledColor()
}

func (t *legacyTheme) FocusColor() color.Color {
	return FocusColor()
}

func (t *legacyTheme) HoverColor() color.Color {
	return HoverColor()
}

func (t *legacyTheme) PlaceHolderColor() color.Color {
	return PlaceHolderColor()
}

func (t *legacyTheme) PrimaryColor() color.Color {
	return PrimaryColor()
}

func (t *legacyTheme) ScrollBarColor() color.Color {
	return ScrollBarColor()
}

func (t *legacyTheme) ShadowColor() color.Color {
	return ShadowColor()
}

func (t *legacyTheme) TextColor() color.Color {
	return ForegroundColor()
}

func (t *legacyTheme) TextSize() int {
	return int(TextSize())
}

func (t *legacyTheme) TextFont() fyne.Resource {
	return TextFont()
}

func (t *legacyTheme) TextBoldFont() fyne.Resource {
	return TextBoldFont()
}

func (t *legacyTheme) TextItalicFont() fyne.Resource {
	return TextItalicFont()
}

func (t *legacyTheme) TextBoldItalicFont() fyne.Resource {
	return TextBoldItalicFont()
}

func (t *legacyTheme) TextMonospaceFont() fyne.Resource {
	return TextMonospaceFont()
}

func (t *legacyTheme) Padding() int {
	return int(Padding())
}

func (t *legacyTheme) IconInlineSize() int {
	return int(IconInlineSize())
}

func (t *legacyTheme) ScrollBarSize() int {
	return int(ScrollBarSize())
}

func (t *legacyTheme) ScrollBarSmallSize() int {
	return int(ScrollBarSmallSize())
}
