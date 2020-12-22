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
	//
	// Since: 1.4.0
	ColorRed = "red"
	// ColorOrange is the orange primary color name
	//
	// Since: 1.4.0
	ColorOrange = "orange"
	// ColorYellow is the yellow primary color name
	//
	// Since: 1.4.0
	ColorYellow = "yellow"
	// ColorGreen is the green primary color name
	//
	// Since: 1.4.0
	ColorGreen = "green"
	// ColorBlue is the blue primary color name
	//
	// Since: 1.4.0
	ColorBlue = "blue"
	// ColorPurple is the purple primary color name
	//
	// Since: 1.4.0
	ColorPurple = "purple"
	// ColorBrown is the brown primary color name
	//
	// Since: 1.4.0
	ColorBrown = "brown"
	// ColorGray is the gray primary color name
	//
	// Since: 1.4.0
	ColorGray = "gray"

	// ColorNameBackground is the name of theme lookup for background color.
	//
	// Since 2.0.0
	ColorNameBackground fyne.ThemeColorName = "background"

	// ColorNameButton is the name of theme lookup for button color.
	//
	// Since 2.0.0
	ColorNameButton fyne.ThemeColorName = "button"

	// ColorNameDisabledButton is the name of theme lookup for disabled button color.
	//
	// Since 2.0.0
	ColorNameDisabledButton fyne.ThemeColorName = "disabledButton"

	// ColorNameDisabled is the name of theme lookup for disabled foreground color.
	//
	// Since 2.0.0
	ColorNameDisabled fyne.ThemeColorName = "disabled"

	// ColorNameFocus is the name of theme lookup for focus color.
	//
	// Since 2.0.0
	ColorNameFocus fyne.ThemeColorName = "focus"

	// ColorNameForeground is the name of theme lookup for foreground color.
	//
	// Since 2.0.0
	ColorNameForeground fyne.ThemeColorName = "foreground"

	// ColorNameHover is the name of theme lookup for hover color.
	//
	// Since 2.0.0
	ColorNameHover fyne.ThemeColorName = "hover"

	// ColorNamePlaceHolder is the name of theme lookup for placeholder text color.
	//
	// Since 2.0.0
	ColorNamePlaceHolder fyne.ThemeColorName = "placeholder"

	// ColorNamePrimary is the name of theme lookup for primary color.
	//
	// Since 2.0.0
	ColorNamePrimary fyne.ThemeColorName = "primary"

	// ColorNameScrollBar is the name of theme lookup for scrollbar color.
	//
	// Since 2.0.0
	ColorNameScrollBar fyne.ThemeColorName = "scrollBar"

	// ColorNameShadow is the name of theme lookup for shadow color.
	//
	// Since 2.0.0
	ColorNameShadow fyne.ThemeColorName = "shadow"

	// SizeNameInlineIcon is the name of theme lookup for inline icons size.
	//
	// Since 2.0.0
	SizeNameInlineIcon fyne.ThemeSizeName = "iconInline"

	// SizeNamePadding is the name of theme lookup for padding size.
	//
	// Since 2.0.0
	SizeNamePadding fyne.ThemeSizeName = "padding"

	// SizeNameScrollBar is the name of theme lookup for the scrollbar size.
	//
	// Since 2.0.0
	SizeNameScrollBar fyne.ThemeSizeName = "scrollBar"

	// SizeNameScrollBarSmall is the name of theme lookup for the shrunk scrollbar size.
	//
	// Since 2.0.0
	SizeNameScrollBarSmall fyne.ThemeSizeName = "scrollBarSmall"

	// SizeNameText is the name of theme lookup for text size.
	//
	// Since 2.0.0
	SizeNameText fyne.ThemeSizeName = "text"

	// VariantNameDark is the version of a theme that satisfies a user preference for a light look.
	//
	// Since 2.0.0
	VariantNameDark fyne.ThemeVariant = 0

	// VariantNameLight is the version of a theme that satisfies a user preference for a dark look.
	//
	// Since 2.0.0
	VariantNameLight fyne.ThemeVariant = 1

	// potential for adding theme types such as high visibility or monochrome...
	variantNameUserPreference fyne.ThemeVariant = 2 // locally used in builtinTheme for backward compatibility
)

var (
	defaultTheme = setupDefaultTheme()

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

	darkPalette = map[fyne.ThemeColorName]color.Color{
		"background":     color.NRGBA{0x30, 0x30, 0x30, 0xff},
		"button":         color.Transparent,
		"disabled":       color.NRGBA{0xff, 0xff, 0xff, 0x42},
		"disabledButton": color.NRGBA{0x26, 0x26, 0x26, 0xff},
		"foreground":     color.NRGBA{0xff, 0xff, 0xff, 0xff},
		"hover":          color.NRGBA{0xff, 0xff, 0xff, 0x0f},
		"placeholder":    color.NRGBA{0xb2, 0xb2, 0xb2, 0xff},
		"scrollBar":      color.NRGBA{0x0, 0x0, 0x0, 0x99},
		"shadow":         color.NRGBA{0x0, 0x0, 0x0, 0x66},
	}

	lightPalette = map[fyne.ThemeColorName]color.Color{
		"background":     color.NRGBA{0xff, 0xff, 0xff, 0xff},
		"button":         color.Transparent,
		"disabled":       color.NRGBA{0x0, 0x0, 0x0, 0x42},
		"disabledButton": color.NRGBA{0xe5, 0xe5, 0xe5, 0xff},
		"foreground":     color.NRGBA{0x21, 0x21, 0x21, 0xff},
		"hover":          color.NRGBA{0x0, 0x0, 0x0, 0x0f},
		"placeholder":    color.NRGBA{0x88, 0x88, 0x88, 0xff},
		"scrollBar":      color.NRGBA{0x0, 0x0, 0x0, 0x99},
		"shadow":         color.NRGBA{0x0, 0x0, 0x0, 0x33}}
)

type builtinTheme struct {
	variant fyne.ThemeVariant

	regular, bold, italic, boldItalic, monospace fyne.Resource
}

// DefaultTheme returns a built-in theme that can adapt to the user preference of light or dark colors.
//
// Since 2.0.0
func DefaultTheme() fyne.Theme {
	return defaultTheme
}

// LightTheme defines the built in light theme colors and sizes
func LightTheme() fyne.Theme {
	theme := &builtinTheme{variant: VariantNameLight}

	theme.initFonts()
	return theme
}

// DarkTheme defines the built in dark theme colors and sizes
func DarkTheme() fyne.Theme {
	theme := &builtinTheme{variant: VariantNameDark}

	theme.initFonts()
	return theme
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

func (t *builtinTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	colors := darkPalette
	if v == VariantNameLight {
		colors = lightPalette
	}

	if n == ColorNamePrimary || n == ColorNameFocus { // TODO new focus color
		return PrimaryColorNamed(fyne.CurrentApp().Settings().PrimaryColor())
	}

	if c, ok := colors[n]; ok {
		return c
	}
	return color.Transparent
}

// TextFont returns the font resource for the regular font style
func (t *builtinTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return t.monospace
	}
	if style.Bold {
		if style.Italic {
			return t.boldItalic
		}
		return t.bold
	}
	if style.Italic {
		return t.italic
	}
	return t.regular
}

func (t *builtinTheme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case SizeNameInlineIcon:
		return 20
	case SizeNamePadding:
		return 4
	case SizeNameScrollBar:
		return 16
	case SizeNameScrollBarSmall:
		return 3
	case SizeNameText:
		return 14
	default:
		return 0
	}
}

func current() fyne.Theme {
	if fyne.CurrentApp() == nil || fyne.CurrentApp().Settings().Theme() == nil {
		return DarkTheme()
	}

	return fyne.CurrentApp().Settings().Theme()
}

// BackgroundColor returns the theme's background color
func BackgroundColor() color.Color {
	return current().Color(ColorNameBackground, currentVariant())
}

// ButtonColor returns the theme's standard button color.
func ButtonColor() color.Color {
	return current().Color(ColorNameButton, currentVariant())
}

// DisabledButtonColor returns the theme's disabled button color.
func DisabledButtonColor() color.Color {
	return current().Color(ColorNameDisabledButton, currentVariant())
}

// TextColor returns the theme's standard text color
//
// Deprecated: Use theme.Foreground colour instead
func TextColor() color.Color {
	return current().Color(ColorNameForeground, currentVariant())
}

// DisabledColor returns the foreground color for a disabled UI element
//
// Since: 2.0.0
func DisabledColor() color.Color {
	return current().Color(ColorNameDisabled, currentVariant())
}

// DisabledTextColor returns the color for a disabledIcon UI element
//
// Deprecated: Use DisabledColor() instead
func DisabledTextColor() color.Color {
	return current().Color(ColorNameDisabled, currentVariant())
}

// PlaceHolderColor returns the theme's standard text color
func PlaceHolderColor() color.Color {
	return current().Color(ColorNamePlaceHolder, currentVariant())
}

// PrimaryColor returns the color used to highlight primary features
func PrimaryColor() color.Color {
	return current().Color(ColorNamePrimary, currentVariant())
}

// HoverColor returns the color used to highlight interactive elements currently under a cursor
func HoverColor() color.Color {
	return current().Color(ColorNameHover, currentVariant())
}

// FocusColor returns the color used to highlight a focused widget
func FocusColor() color.Color {
	return current().Color(ColorNameFocus, currentVariant())
}

// ForegroundColor returns the theme's standard foreground color for text and icons
//
// Since: 2.0.0
func ForegroundColor() color.Color {
	return current().Color(ColorNameForeground, currentVariant())
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func ScrollBarColor() color.Color {
	return current().Color(ColorNameScrollBar, currentVariant())
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation
func ShadowColor() color.Color {
	return current().Color(ColorNameShadow, currentVariant())
}

// TextSize returns the standard text size
func TextSize() float32 {
	return current().Size(SizeNameText)
}

// TextFont returns the font resource for the regular font style
func TextFont() fyne.Resource {
	return current().Font(fyne.TextStyle{})
}

// TextBoldFont returns the font resource for the bold font style
func TextBoldFont() fyne.Resource {
	return current().Font(fyne.TextStyle{Bold: true})
}

// TextItalicFont returns the font resource for the italic font style
func TextItalicFont() fyne.Resource {
	return current().Font(fyne.TextStyle{Italic: true})
}

// TextBoldItalicFont returns the font resource for the bold and italic font style
func TextBoldItalicFont() fyne.Resource {
	return current().Font(fyne.TextStyle{Bold: true, Italic: true})
}

// TextMonospaceFont returns the font resource for the monospace font face
func TextMonospaceFont() fyne.Resource {
	return current().Font(fyne.TextStyle{Monospace: true})
}

// Padding is the standard gap between elements and the border around interface
// elements
func Padding() float32 {
	return current().Size(SizeNamePadding)
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func IconInlineSize() float32 {
	return current().Size(SizeNameInlineIcon)
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer
func ScrollBarSize() float32 {
	return current().Size(SizeNameScrollBar)
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer
func ScrollBarSmallSize() float32 {
	return current().Size(SizeNameScrollBarSmall)
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
//
// Since: 1.4.0
func PrimaryColorNames() []string {
	return []string{ColorRed, ColorOrange, ColorYellow, ColorGreen, ColorBlue, ColorPurple, ColorBrown, ColorGray}
}

// PrimaryColorNamed returns a theme specific color value for a named primary color.
//
// Since 1.4.0
func PrimaryColorNamed(name string) color.Color {
	col, ok := primaryColors[name]
	if !ok {
		return primaryColors[ColorBlue]
	}
	return col
}

func currentVariant() fyne.ThemeVariant {
	if std, ok := current().(*builtinTheme); ok {
		if std.variant != variantNameUserPreference {
			return std.variant // override if using the old LightTheme() or DarkTheme() constructor
		}
	}

	return fyne.CurrentApp().Settings().ThemeVariant()
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

func setupDefaultTheme() fyne.Theme {
	theme := &builtinTheme{}

	theme.initFonts()
	return theme
}
