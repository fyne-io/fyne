// Package theme defines how a Fyne app should look when rendered
package theme // import "fyne.io/fyne/theme"

import (
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne"
)

const (
	// ColorRed is the red primary color name
	ColorRed = "red"
	// ColorOrange is the orange primary color name
	ColorOrange = "orange"
	// ColorYellow is the yellow primary color name
	ColorYellow = "yellow"
	// ColorGreen is the green primary color name
	ColorGreen = "green"
	// ColorBlue is the blue primary color name
	ColorBlue = "blue"
	// ColorPurple is the purple primary color name
	ColorPurple = "purple"
	// ColorBrown is the brown primary color name
	ColorBrown = "brown"
	// ColorGray is the gray primary color name
	ColorGray = "gray"
)

type builtinTheme struct {
	background color.Color

	button, disabledButton, text, placeholder, hover, shadow, disabled, scrollBar color.Color
	regular, bold, italic, boldItalic, monospace                                  fyne.Resource
}

var (
	primaryColors = map[string]color.Color{
		ColorRed:    color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff},
		ColorOrange: color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff},
		ColorYellow: color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0xff},
		ColorGreen:  color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0xff},
		ColorBlue:   color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0xff},
		ColorPurple: color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0xff},
		ColorBrown:  color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0xff},
		ColorGray:   color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0xff},
	}

	//	themeSecondaryColor = color.NRGBA{R: 0xff, G: 0x40, B: 0x81, A: 0xff}
)

// LightTheme defines the built in light theme colors and sizes
func LightTheme() fyne.Theme {
	theme := &builtinTheme{
		background:     color.NRGBA{0xff, 0xff, 0xff, 0xff},
		button:         color.Transparent,
		disabled:       color.NRGBA{0x0, 0x0, 0x0, 0x42},
		disabledButton: color.NRGBA{0xe5, 0xe5, 0xe5, 0xff},
		text:           color.NRGBA{0x21, 0x21, 0x21, 0xff},
		placeholder:    color.NRGBA{0x88, 0x88, 0x88, 0xff},
		hover:          color.NRGBA{0x0, 0x0, 0x0, 0x0f},
		scrollBar:      color.NRGBA{0x0, 0x0, 0x0, 0x99},
		shadow:         color.NRGBA{0x0, 0x0, 0x0, 0x33},
	}

	theme.initFonts()
	return theme
}

// DarkTheme defines the built in dark theme colors and sizes
func DarkTheme() fyne.Theme {
	theme := &builtinTheme{
		background:     color.NRGBA{0x30, 0x30, 0x30, 0xff},
		button:         color.Transparent,
		disabled:       color.NRGBA{0xff, 0xff, 0xff, 0x42},
		disabledButton: color.NRGBA{0x26, 0x26, 0x26, 0xff},
		text:           color.NRGBA{0xff, 0xff, 0xff, 0xff},
		placeholder:    color.NRGBA{0xb2, 0xb2, 0xb2, 0xff},
		hover:          color.NRGBA{0xff, 0xff, 0xff, 0x0f},
		scrollBar:      color.NRGBA{0x0, 0x0, 0x0, 0x99},
		shadow:         color.NRGBA{0x0, 0x0, 0x0, 0x66},
	}

	theme.initFonts()
	return theme
}

func (t *builtinTheme) BackgroundColor() color.Color {
	return t.background
}

// ButtonColor returns the theme's standard button color.
func (t *builtinTheme) ButtonColor() color.Color {
	return t.button
}

// DisabledButtonColor returns the theme's disabled button color.
func (t *builtinTheme) DisabledButtonColor() color.Color {
	return t.disabledButton
}

// HyperlinkColor returns the theme's standard hyperlink color.
//
// Deprecated: Hyperlinks now use the primary color for consistency.
func (t *builtinTheme) HyperlinkColor() color.Color {
	return t.PrimaryColor()
}

// TextColor returns the theme's standard text color
func (t *builtinTheme) TextColor() color.Color {
	return t.text
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *builtinTheme) DisabledTextColor() color.Color {
	return t.disabled
}

// IconColor returns the theme's standard text color.
//
// Deprecated: Icons now use the text colour for consistency.
func (t *builtinTheme) IconColor() color.Color {
	return t.text
}

// DisabledIconColor returns the color for a disabledIcon UI element.
//
// Deprecated: Disabled icons match disabled text color for consistency.
func (t *builtinTheme) DisabledIconColor() color.Color {
	return t.disabled
}

// PlaceHolderColor returns the theme's placeholder text color
func (t *builtinTheme) PlaceHolderColor() color.Color {
	return t.placeholder
}

// PrimaryColor returns the color used to highlight primary features
func (t *builtinTheme) PrimaryColor() color.Color {
	return PrimaryColorNamed(fyne.CurrentApp().Settings().PrimaryColor())
}

// HoverColor returns the color used to highlight interactive elements currently under a cursor
func (t *builtinTheme) HoverColor() color.Color {
	return t.hover
}

// FocusColor returns the color used to highlight a focused widget
func (t *builtinTheme) FocusColor() color.Color {
	return t.PrimaryColor()
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
	variantPath := strings.Replace(env, "Regular", variant, -1)

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
	t.boldItalic = bolditalic
	t.monospace = monospace

	font := os.Getenv("FYNE_FONT")
	if font != "" {
		t.regular = loadCustomFont(font, "Regular", regular)
		t.bold = loadCustomFont(font, "Bold", bold)
		t.italic = loadCustomFont(font, "Italic", italic)
		t.boldItalic = loadCustomFont(font, "BoldItalic", bolditalic)
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

// TextBoldFont returns the font resource for the bold font style
func (t *builtinTheme) TextBoldFont() fyne.Resource {
	return t.bold
}

// TextItalicFont returns the font resource for the italic font style
func (t *builtinTheme) TextItalicFont() fyne.Resource {
	return t.italic
}

// TextBoldItalicFont returns the font resource for the bold and italic font style
func (t *builtinTheme) TextBoldItalicFont() fyne.Resource {
	return t.boldItalic
}

// TextMonospaceFont returns the font resource for the monospace font face
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

// ButtonColor returns the theme's standard button color.
func ButtonColor() color.Color {
	return current().ButtonColor()
}

// DisabledButtonColor returns the theme's disabled button color.
func DisabledButtonColor() color.Color {
	return current().DisabledButtonColor()
}

// HyperlinkColor returns the theme's standard hyperlink color.
//
// Deprecated: Hyperlinks now use the primary color for consistency.
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

// IconColor returns the theme's standard text color.
//
// Deprecated: Icons now use the text colour for consistency.
func IconColor() color.Color {
	return current().IconColor()
}

// DisabledIconColor returns the color for a disabledIcon UI element.
//
// Deprecated: Disabled icons match disabled text color for consistency.
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

// FocusColor returns the color used to highlight a focused widget
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

// TextBoldFont returns the font resource for the bold font style
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

// TextMonospaceFont returns the font resource for the monospace font face
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

// DefaultTextBoldFont returns the font resource for the built-in bold font style
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

// DefaultTextMonospaceFont returns the font resource for the built-in monospace font face
func DefaultTextMonospaceFont() fyne.Resource {
	return monospace
}

// PrimaryColorNames returns a list of the standard primary color options.
func PrimaryColorNames() []string {
	return []string{ColorRed, ColorOrange, ColorYellow, ColorGreen, ColorBlue, ColorPurple, ColorBrown, ColorGray}
}

// PrimaryColorNamed returns a theme specific color value for a named primary color.
func PrimaryColorNamed(name string) color.Color {
	col, ok := primaryColors[name]
	if !ok {
		return primaryColors[ColorBlue]
	}
	return col
}
