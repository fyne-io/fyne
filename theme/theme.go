// Package theme defines how a Fyne app should look when rendered
package theme // import "fyne.io/fyne/theme"

import (
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne"
)

type builtinTheme struct {
	background color.Color

	button, text, icon, hyperlink, placeholder, primary, hover, scrollBar, shadow color.Color
	regular, bold, italic, bolditalic, monospace                                  fyne.Resource
	disabledButton, disabledIcon, disabledText                                    color.Color
}

// LightTheme defines the built in light theme colors and sizes
func LightTheme() fyne.Theme {
	theme := &builtinTheme{
		background:     color.NRGBA{0xf5, 0xf5, 0xf5, 0xff},
		button:         color.NRGBA{0xd9, 0xd9, 0xd9, 0xff},
		disabledButton: color.NRGBA{0xe7, 0xe7, 0xe7, 0xff},
		text:           color.NRGBA{0x21, 0x21, 0x21, 0xff},
		disabledText:   color.NRGBA{0x80, 0x80, 0x80, 0xff},
		icon:           color.NRGBA{0x21, 0x21, 0x21, 0xff},
		disabledIcon:   color.NRGBA{0x80, 0x80, 0x80, 0xff},
		hyperlink:      color.NRGBA{0x0, 0x0, 0xd9, 0xff},
		placeholder:    color.NRGBA{0x88, 0x88, 0x88, 0xff},
		primary:        color.NRGBA{0x9f, 0xa8, 0xda, 0xff},
		hover:          color.NRGBA{0xe7, 0xe7, 0xe7, 0xff},
		scrollBar:      color.NRGBA{0x0, 0x0, 0x0, 0x99},
		shadow:         color.NRGBA{0x0, 0x0, 0x0, 0x33},
	}

	theme.initFonts()
	return theme
}

// DarkTheme defines the built in dark theme colors and sizes
func DarkTheme() fyne.Theme {
	theme := &builtinTheme{
		background:     color.NRGBA{0x42, 0x42, 0x42, 0xff},
		button:         color.NRGBA{0x21, 0x21, 0x21, 0xff},
		disabledButton: color.NRGBA{0x31, 0x31, 0x31, 0xff},
		text:           color.NRGBA{0xff, 0xff, 0xff, 0xff},
		disabledText:   color.NRGBA{0x60, 0x60, 0x60, 0xff},
		icon:           color.NRGBA{0xff, 0xff, 0xff, 0xff},
		disabledIcon:   color.NRGBA{0x60, 0x60, 0x60, 0xff},
		hyperlink:      color.NRGBA{0x99, 0x99, 0xff, 0xff},
		placeholder:    color.NRGBA{0xb2, 0xb2, 0xb2, 0xff},
		primary:        color.NRGBA{0x1a, 0x23, 0x7e, 0xff},
		hover:          color.NRGBA{0x31, 0x31, 0x31, 0xff},
		scrollBar:      color.NRGBA{0x0, 0x0, 0x0, 0x99},
		shadow:         color.NRGBA{0x0, 0x0, 0x0, 0x66},
	}

	theme.initFonts()
	return theme
}

// Shade will darken a light color and lighten a light color by the given % value
// should use 4,8, 12 or 14 percent darken/lighten to implement google material design rules
// Note that the percent lightening is twice the percent given. See material.io/design/interaction/states.
func Shade(c color.Color, pct uint32) color.Color {
	if pct == 0 {
		return c
	}
	r, g, b, a := c.RGBA()
	if pct > 50 {
		pct = 50
	}
	if r+g+b < 3*0x8080 {
		// Lighten by twice the percent given.
		return color.NRGBA{
			R: uint8((r + (0x10000-r)*pct/50) >> 8),
			G: uint8((g + (0x10000-g)*pct/50) >> 8),
			B: uint8((b + (0x10000-b)*pct/50) >> 8),
			A: uint8(a >> 8),
		}
	}
	// Darken by given percent
	return color.NRGBA{
		R: uint8((r * (100 - pct) / 100) >> 8),
		G: uint8((g * (100 - pct) / 100) >> 8),
		B: uint8((b * (100 - pct) / 100) >> 8),
		A: uint8(a >> 8),
	}
}

// PressedShade is the shade used for pressed buttons, in %
const PressedShade = 14

// HoveredShade is the shade used for hovered buttons, in %
const HoveredShade = 4

// FocusedShade is the shade used for focused buttons, in %
const FocusedShade = 8

func (t *builtinTheme) BackgroundColor() color.Color {
	return t.background
}

// ButtonColor returns the theme's standard button color
func (t *builtinTheme) ButtonColor() color.Color {
	return t.button
}

// DisabledButtonColor returns the theme's disabled button color
func (t *builtinTheme) DisabledButtonColor() color.Color {
	return t.disabledButton
}

// HyperlinkColor returns the theme's standard hyperlink color
func (t *builtinTheme) HyperlinkColor() color.Color {
	return t.hyperlink
}

// TextColor returns the theme's standard text color
func (t *builtinTheme) TextColor() color.Color {
	return t.text
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *builtinTheme) DisabledTextColor() color.Color {
	return t.disabledText
}

// IconColor returns the theme's standard text color
func (t *builtinTheme) IconColor() color.Color {
	return t.icon
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *builtinTheme) DisabledIconColor() color.Color {
	return t.disabledIcon
}

// PlaceHolderColor returns the theme's placeholder text color
func (t *builtinTheme) PlaceHolderColor() color.Color {
	return t.placeholder
}

// PrimaryColor returns the color used to highlight primary features
func (t *builtinTheme) PrimaryColor() color.Color {
	return t.primary
}

// HoverColor returns the color used to highlight interactive elements currently under a cursor
func (t *builtinTheme) HoverColor() color.Color {
	return t.hover
}

// FocusColor returns the color used to highlight a focused widget
func (t *builtinTheme) FocusColor() color.Color {
	return Shade(t.primary, FocusedShade)
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func (t *builtinTheme) ScrollBarColor() color.Color {
	return t.scrollBar
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation
func (t *builtinTheme) ShadowColor() color.Color {
	return t.shadow
}

// TextSize returns the standard text size
func (t *builtinTheme) TextSize() int {
	return 14
}

func loadCustomFont(env, variant string, fallback fyne.Resource) fyne.Resource {
	variantPath := strings.Replace(env, "Regular", variant, 0)

	res, err := fyne.LoadResourceFromPath(variantPath)
	if err != nil {
		fyne.LogError("Error loading specified font", err)
		return fallback
	}

	return res
}

func (t *builtinTheme) initFonts() {
	t.regular = regular
	t.bold = bold
	t.italic = italic
	t.bolditalic = bolditalic
	t.monospace = monospace

	font := os.Getenv("FYNE_FONT")
	if font != "" {
		t.regular = loadCustomFont(font, "Regular", regular)
		t.bold = loadCustomFont(font, "Bold", bold)
		t.italic = loadCustomFont(font, "Italic", italic)
		t.bolditalic = loadCustomFont(font, "BoldItalic", bolditalic)
	}
	font = os.Getenv("FYNE_FONT_MONOSPACE")
	if font != "" {
		t.monospace = loadCustomFont(font, "Regular", monospace)
	}
}

// TextFont returns the font resource for the regular font style
func (t *builtinTheme) TextFont() fyne.Resource {
	return t.regular
}

// TextBoldFont retutns the font resource for the bold font style
func (t *builtinTheme) TextBoldFont() fyne.Resource {
	return t.bold
}

// TextItalicFont returns the font resource for the italic font style
func (t *builtinTheme) TextItalicFont() fyne.Resource {
	return t.italic
}

// TextBoldItalicFont returns the font resource for the bold and italic font style
func (t *builtinTheme) TextBoldItalicFont() fyne.Resource {
	return t.bolditalic
}

// TextMonospaceFont retutns the font resource for the monospace font face
func (t *builtinTheme) TextMonospaceFont() fyne.Resource {
	return t.monospace
}

// Padding is the standard gap between elements and the border around interface
// elements
func (t *builtinTheme) Padding() int {
	return 4
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func (t *builtinTheme) IconInlineSize() int {
	return 20
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer
func (t *builtinTheme) ScrollBarSize() int {
	return 16
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer
func (t *builtinTheme) ScrollBarSmallSize() int {
	return 3
}

func current() fyne.Theme {
	if fyne.CurrentApp() == nil || fyne.CurrentApp().Settings().Theme() == nil {
		return DarkTheme()
	}

	return fyne.CurrentApp().Settings().Theme()
}

// BackgroundColor returns the theme's background color
func BackgroundColor() color.Color {
	return current().BackgroundColor()
}

// ButtonColor returns the theme's standard button color
func ButtonColor() color.Color {
	return current().ButtonColor()
}

// DisabledButtonColor returns the theme's disabled button color
func DisabledButtonColor() color.Color {
	return current().DisabledButtonColor()
}

// HyperlinkColor returns the theme's standard hyperlink color
func HyperlinkColor() color.Color {
	return current().HyperlinkColor()
}

// TextColor returns the theme's standard text color
func TextColor() color.Color {
	return current().TextColor()
}

// DisabledTextColor returns the color for a disabledIcon UI element
func DisabledTextColor() color.Color {
	return current().DisabledTextColor()
}

// IconColor returns the theme's standard text color
func IconColor() color.Color {
	return current().IconColor()
}

// DisabledIconColor returns the color for a disabledIcon UI element
func DisabledIconColor() color.Color {
	return current().DisabledIconColor()
}

// PlaceHolderColor returns the theme's standard text color
func PlaceHolderColor() color.Color {
	return current().PlaceHolderColor()
}

// PrimaryColor returns the color used to highlight primary features
func PrimaryColor() color.Color {
	return current().PrimaryColor()
}

// HoverColor returns the color used to highlight interactive elements currently under a cursor
func HoverColor() color.Color {
	return current().HoverColor()
}

// PressedColor returns the colour used for a pressed button
func PressedColor() color.Color {
	return Shade(FocusColor(), PressedShade)
}

// HoverFocusedColor returns the colour used for a focused/primary hovered button
func HoverFocusedColor() color.Color {
	return Shade(FocusColor(), HoveredShade)
}

// FocusColor returns the color used to highlight a focussed widget
func FocusColor() color.Color {
	return current().FocusColor()
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func ScrollBarColor() color.Color {
	return current().ScrollBarColor()
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation
func ShadowColor() color.Color {
	return current().ShadowColor()
}

// TextSize returns the standard text size
func TextSize() int {
	return current().TextSize()
}

// TextFont returns the font resource for the regular font style
func TextFont() fyne.Resource {
	return current().TextFont()
}

// TextBoldFont retutns the font resource for the bold font style
func TextBoldFont() fyne.Resource {
	return current().TextBoldFont()
}

// TextItalicFont returns the font resource for the italic font style
func TextItalicFont() fyne.Resource {
	return current().TextItalicFont()
}

// TextBoldItalicFont returns the font resource for the bold and italic font style
func TextBoldItalicFont() fyne.Resource {
	return current().TextBoldItalicFont()
}

// TextMonospaceFont retutns the font resource for the monospace font face
func TextMonospaceFont() fyne.Resource {
	return current().TextMonospaceFont()
}

// Padding is the standard gap between elements and the border around interface
// elements
func Padding() int {
	return current().Padding()
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func IconInlineSize() int {
	return current().IconInlineSize()
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer
func ScrollBarSize() int {
	return current().ScrollBarSize()
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer
func ScrollBarSmallSize() int {
	return current().ScrollBarSmallSize()
}

// DefaultTextFont returns the font resource for the built-in regular font style
func DefaultTextFont() fyne.Resource {
	return regular
}

// DefaultTextBoldFont retutns the font resource for the built-in bold font style
func DefaultTextBoldFont() fyne.Resource {
	return bold
}

// DefaultTextItalicFont returns the font resource for the built-in italic font style
func DefaultTextItalicFont() fyne.Resource {
	return italic
}

// DefaultTextBoldItalicFont returns the font resource for the built-in bold and italic font style
func DefaultTextBoldItalicFont() fyne.Resource {
	return bolditalic
}

// DefaultTextMonospaceFont retutns the font resource for the built-in monospace font face
func DefaultTextMonospaceFont() fyne.Resource {
	return monospace
}
