// Package theme defines how a Fyne app should look when rendered.
package theme // import "fyne.io/fyne/v2/theme"

import (
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
)

const (
	// ColorRed is the red primary color name.
	//
	// Since: 1.4
	ColorRed = "red"
	// ColorOrange is the orange primary color name.
	//
	// Since: 1.4
	ColorOrange = "orange"
	// ColorYellow is the yellow primary color name.
	//
	// Since: 1.4
	ColorYellow = "yellow"
	// ColorGreen is the green primary color name.
	//
	// Since: 1.4
	ColorGreen = "green"
	// ColorBlue is the blue primary color name.
	//
	// Since: 1.4
	ColorBlue = "blue"
	// ColorPurple is the purple primary color name.
	//
	// Since: 1.4
	ColorPurple = "purple"
	// ColorBrown is the brown primary color name.
	//
	// Since: 1.4
	ColorBrown = "brown"
	// ColorGray is the gray primary color name.
	//
	// Since: 1.4
	ColorGray = "gray"

	// ColorNameBackground is the name of theme lookup for background color.
	//
	// Since: 2.0
	ColorNameBackground fyne.ThemeColorName = "background"

	// ColorNameButton is the name of theme lookup for button color.
	//
	// Since: 2.0
	ColorNameButton fyne.ThemeColorName = "button"

	// ColorNameDisabledButton is the name of theme lookup for disabled button color.
	//
	// Since: 2.0
	ColorNameDisabledButton fyne.ThemeColorName = "disabledButton"

	// ColorNameDisabled is the name of theme lookup for disabled foreground color.
	//
	// Since: 2.0
	ColorNameDisabled fyne.ThemeColorName = "disabled"

	// ColorNameError is the name of theme lookup for foreground error color.
	//
	// Since: 2.0
	ColorNameError fyne.ThemeColorName = "error"

	// ColorNameFocus is the name of theme lookup for focus color.
	//
	// Since: 2.0
	ColorNameFocus fyne.ThemeColorName = "focus"

	// ColorNameForeground is the name of theme lookup for foreground color.
	//
	// Since: 2.0
	ColorNameForeground fyne.ThemeColorName = "foreground"

	// ColorNameHover is the name of theme lookup for hover color.
	//
	// Since: 2.0
	ColorNameHover fyne.ThemeColorName = "hover"

	// ColorNameInputBackground is the name of theme lookup for background color of an input field.
	//
	// Since: 2.0
	ColorNameInputBackground fyne.ThemeColorName = "inputBackground"

	// ColorNameInputBorder is the name of theme lookup for border color of an input field.
	//
	// Since: 2.3
	ColorNameInputBorder fyne.ThemeColorName = "inputBorder"

	// ColorNameMenuBackground is the name of theme lookup for background color of menus.
	//
	// Since: 2.3
	ColorNameMenuBackground fyne.ThemeColorName = "menuBackground"

	// ColorNameOverlayBackground is the name of theme lookup for background color of overlays like dialogs.
	//
	// Since: 2.3
	ColorNameOverlayBackground fyne.ThemeColorName = "overlayBackground"

	// ColorNamePlaceHolder is the name of theme lookup for placeholder text color.
	//
	// Since: 2.0
	ColorNamePlaceHolder fyne.ThemeColorName = "placeholder"

	// ColorNamePressed is the name of theme lookup for the tap overlay color.
	//
	// Since: 2.0
	ColorNamePressed fyne.ThemeColorName = "pressed"

	// ColorNamePrimary is the name of theme lookup for primary color.
	//
	// Since: 2.0
	ColorNamePrimary fyne.ThemeColorName = "primary"

	// ColorNameScrollBar is the name of theme lookup for scrollbar color.
	//
	// Since: 2.0
	ColorNameScrollBar fyne.ThemeColorName = "scrollBar"

	// ColorNameSelection is the name of theme lookup for selection color.
	//
	// Since: 2.1
	ColorNameSelection fyne.ThemeColorName = "selection"

	// ColorNameSeparator is the name of theme lookup for separator bars.
	//
	// Since: 2.3
	ColorNameSeparator fyne.ThemeColorName = "separator"

	// ColorNameShadow is the name of theme lookup for shadow color.
	//
	// Since: 2.0
	ColorNameShadow fyne.ThemeColorName = "shadow"

	// ColorNameSuccess is the name of theme lookup for foreground success color.
	//
	// Since: 2.3
	ColorNameSuccess fyne.ThemeColorName = "success"

	// ColorNameWarning is the name of theme lookup for foreground warning color.
	//
	// Since: 2.3
	ColorNameWarning fyne.ThemeColorName = "warning"

	// SizeNameCaptionText is the name of theme lookup for helper text size, normally smaller than regular text size.
	//
	// Since: 2.0
	SizeNameCaptionText fyne.ThemeSizeName = "helperText"

	// SizeNameInlineIcon is the name of theme lookup for inline icons size.
	//
	// Since: 2.0
	SizeNameInlineIcon fyne.ThemeSizeName = "iconInline"

	// SizeNameInnerPadding is the name of theme lookup for internal widget padding size.
	//
	// Since: 2.3
	SizeNameInnerPadding fyne.ThemeSizeName = "innerPadding"

	// SizeNameLineSpacing is the name of theme lookup for between text line spacing.
	//
	// Since: 2.3
	SizeNameLineSpacing fyne.ThemeSizeName = "lineSpacing"

	// SizeNamePadding is the name of theme lookup for padding size.
	//
	// Since: 2.0
	SizeNamePadding fyne.ThemeSizeName = "padding"

	// SizeNameScrollBar is the name of theme lookup for the scrollbar size.
	//
	// Since: 2.0
	SizeNameScrollBar fyne.ThemeSizeName = "scrollBar"

	// SizeNameScrollBarSmall is the name of theme lookup for the shrunk scrollbar size.
	//
	// Since: 2.0
	SizeNameScrollBarSmall fyne.ThemeSizeName = "scrollBarSmall"

	// SizeNameSeparatorThickness is the name of theme lookup for the thickness of a separator.
	//
	// Since: 2.0
	SizeNameSeparatorThickness fyne.ThemeSizeName = "separator"

	// SizeNameText is the name of theme lookup for text size.
	//
	// Since: 2.0
	SizeNameText fyne.ThemeSizeName = "text"

	// SizeNameHeadingText is the name of theme lookup for text size of a heading.
	//
	// Since: 2.1
	SizeNameHeadingText fyne.ThemeSizeName = "headingText"

	// SizeNameSubHeadingText is the name of theme lookup for text size of a sub-heading.
	//
	// Since: 2.1
	SizeNameSubHeadingText fyne.ThemeSizeName = "subHeadingText"

	// SizeNameInputBorder is the name of theme lookup for input border size.
	//
	// Since: 2.0
	SizeNameInputBorder fyne.ThemeSizeName = "inputBorder"

	// VariantDark is the version of a theme that satisfies a user preference for a light look.
	//
	// Since: 2.0
	VariantDark fyne.ThemeVariant = 0

	// VariantLight is the version of a theme that satisfies a user preference for a dark look.
	//
	// Since: 2.0
	VariantLight fyne.ThemeVariant = 1

	// potential for adding theme types such as high visibility or monochrome...
	variantNameUserPreference fyne.ThemeVariant = 2 // locally used in builtinTheme for backward compatibility
)

// BackgroundColor returns the theme's background color.
func BackgroundColor() color.Color {
	return safeColorLookup(ColorNameBackground, currentVariant())
}

// ButtonColor returns the theme's standard button color.
func ButtonColor() color.Color {
	return safeColorLookup(ColorNameButton, currentVariant())
}

// CaptionTextSize returns the size for caption text.
func CaptionTextSize() float32 {
	return current().Size(SizeNameCaptionText)
}

// DarkTheme defines the built-in dark theme colors and sizes.
//
// Deprecated: This method ignores user preference and should not be used, it will be removed in v3.0.
func DarkTheme() fyne.Theme {
	theme := &builtinTheme{variant: VariantDark}

	theme.initFonts()
	return theme
}

// DefaultTextBoldFont returns the font resource for the built-in bold font style.
func DefaultTextBoldFont() fyne.Resource {
	return bold
}

// DefaultTextBoldItalicFont returns the font resource for the built-in bold and italic font style.
func DefaultTextBoldItalicFont() fyne.Resource {
	return bolditalic
}

// DefaultTextFont returns the font resource for the built-in regular font style.
func DefaultTextFont() fyne.Resource {
	return regular
}

// DefaultTextItalicFont returns the font resource for the built-in italic font style.
func DefaultTextItalicFont() fyne.Resource {
	return italic
}

// DefaultTextMonospaceFont returns the font resource for the built-in monospace font face.
func DefaultTextMonospaceFont() fyne.Resource {
	return monospace
}

// DefaultSymbolFont returns the font resource for the built-in symbol font.
//
// Since: 2.2
func DefaultSymbolFont() fyne.Resource {
	return symbol
}

// DefaultTheme returns a built-in theme that can adapt to the user preference of light or dark colors.
//
// Since: 2.0
func DefaultTheme() fyne.Theme {
	if defaultTheme == nil {
		defaultTheme = setupDefaultTheme()
	}

	return defaultTheme
}

// DisabledButtonColor returns the theme's disabled button color.
func DisabledButtonColor() color.Color {
	return safeColorLookup(ColorNameDisabledButton, currentVariant())
}

// DisabledColor returns the foreground color for a disabled UI element.
//
// Since: 2.0
func DisabledColor() color.Color {
	return safeColorLookup(ColorNameDisabled, currentVariant())
}

// DisabledTextColor returns the theme's disabled text color - this is actually the disabled color since 1.4.
//
// Deprecated: Use theme.DisabledColor() colour instead.
func DisabledTextColor() color.Color {
	return DisabledColor()
}

// ErrorColor returns the theme's error foreground color.
//
// Since: 2.0
func ErrorColor() color.Color {
	return safeColorLookup(ColorNameError, currentVariant())
}

// FocusColor returns the color used to highlight a focused widget.
func FocusColor() color.Color {
	return safeColorLookup(ColorNameFocus, currentVariant())
}

// ForegroundColor returns the theme's standard foreground color for text and icons.
//
// Since: 2.0
func ForegroundColor() color.Color {
	return safeColorLookup(ColorNameForeground, currentVariant())
}

// HoverColor returns the color used to highlight interactive elements currently under a cursor.
func HoverColor() color.Color {
	return safeColorLookup(ColorNameHover, currentVariant())
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func IconInlineSize() float32 {
	return current().Size(SizeNameInlineIcon)
}

// InnerPadding is the standard gap between element content and the outside edge of a widget.
//
// Since: 2.3
func InnerPadding() float32 {
	return current().Size(SizeNameInnerPadding)
}

// InputBackgroundColor returns the color used to draw underneath input elements.
func InputBackgroundColor() color.Color {
	return current().Color(ColorNameInputBackground, currentVariant())
}

// InputBorderColor returns the color used to draw underneath input elements.
//
// Since: 2.3
func InputBorderColor() color.Color {
	return current().Color(ColorNameInputBorder, currentVariant())
}

// InputBorderSize returns the input border size (or underline size for an entry).
//
// Since: 2.0
func InputBorderSize() float32 {
	return current().Size(SizeNameInputBorder)
}

// LightTheme defines the built-in light theme colors and sizes.
//
// Deprecated: This method ignores user preference and should not be used, it will be removed in v3.0.
func LightTheme() fyne.Theme {
	theme := &builtinTheme{variant: VariantLight}

	theme.initFonts()
	return theme
}

// LineSpacing is the default gap between multiple lines of text.
//
// Since: 2.3
func LineSpacing() float32 {
	return current().Size(SizeNameLineSpacing)
}

// MenuBackgroundColor returns the theme's background color for menus.
//
// Since: 2.3
func MenuBackgroundColor() color.Color {
	return safeColorLookup(ColorNameMenuBackground, currentVariant())
}

// OverlayBackgroundColor returns the theme's background color for overlays like dialogs.
//
// Since: 2.3
func OverlayBackgroundColor() color.Color {
	return safeColorLookup(ColorNameOverlayBackground, currentVariant())
}

// Padding is the standard gap between elements and the border around interface elements.
func Padding() float32 {
	return current().Size(SizeNamePadding)
}

// PlaceHolderColor returns the theme's standard text color.
func PlaceHolderColor() color.Color {
	return safeColorLookup(ColorNamePlaceHolder, currentVariant())
}

// PressedColor returns the color used to overlap tapped features.
//
// Since: 2.0
func PressedColor() color.Color {
	return safeColorLookup(ColorNamePressed, currentVariant())
}

// PrimaryColor returns the color used to highlight primary features.
func PrimaryColor() color.Color {
	return safeColorLookup(ColorNamePrimary, currentVariant())
}

// PrimaryColorNamed returns a theme specific color value for a named primary color.
//
// Since: 1.4
func PrimaryColorNamed(name string) color.Color {
	return primaryColorNamed(name)
}

// PrimaryColorNames returns a list of the standard primary color options.
//
// Since: 1.4
func PrimaryColorNames() []string {
	return []string{ColorRed, ColorOrange, ColorYellow, ColorGreen, ColorBlue, ColorPurple, ColorBrown, ColorGray}
}

// ScrollBarColor returns the color (and translucency) for a scrollBar.
func ScrollBarColor() color.Color {
	return safeColorLookup(ColorNameScrollBar, currentVariant())
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer.
func ScrollBarSize() float32 {
	return current().Size(SizeNameScrollBar)
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer.
func ScrollBarSmallSize() float32 {
	return current().Size(SizeNameScrollBarSmall)
}

// SelectionColor returns the color for a selected element.
//
// Since: 2.1
func SelectionColor() color.Color {
	return safeColorLookup(ColorNameSelection, currentVariant())
}

// SeparatorColor returns the color for the separator element.
//
// Since: 2.3
func SeparatorColor() color.Color {
	return safeColorLookup(ColorNameSeparator, currentVariant())
}

// SeparatorThicknessSize is the standard thickness of the separator widget.
//
// Since: 2.0
func SeparatorThicknessSize() float32 {
	return current().Size(SizeNameSeparatorThickness)
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation.
func ShadowColor() color.Color {
	return safeColorLookup(ColorNameShadow, currentVariant())
}

// SuccessColor returns the theme's success foreground color.
//
// Since: 2.3
func SuccessColor() color.Color {
	return safeColorLookup(ColorNameSuccess, currentVariant())
}

// TextBoldFont returns the font resource for the bold font style.
func TextBoldFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Bold: true})
}

// TextBoldItalicFont returns the font resource for the bold and italic font style.
func TextBoldItalicFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Bold: true, Italic: true})
}

// TextColor returns the theme's standard text color - this is actually the foreground color since 1.4.
//
// Deprecated: Use theme.ForegroundColor() colour instead.
func TextColor() color.Color {
	return safeColorLookup(ColorNameForeground, currentVariant())
}

// TextFont returns the font resource for the regular font style.
func TextFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{})
}

// TextHeadingSize returns the text size for header text.
//
// Since: 2.1
func TextHeadingSize() float32 {
	return current().Size(SizeNameHeadingText)
}

// TextItalicFont returns the font resource for the italic font style.
func TextItalicFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Italic: true})
}

// TextMonospaceFont returns the font resource for the monospace font face.
func TextMonospaceFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Monospace: true})
}

// TextSize returns the standard text size.
func TextSize() float32 {
	return current().Size(SizeNameText)
}

// TextSubHeadingSize returns the text size for sub-header text.
//
// Since: 2.1
func TextSubHeadingSize() float32 {
	return current().Size(SizeNameSubHeadingText)
}

// WarningColor returns the theme's warning foreground color.
//
// Since: 2.3
func WarningColor() color.Color {
	return safeColorLookup(ColorNameWarning, currentVariant())
}

var (
	defaultTheme fyne.Theme

	errorColor   = color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
	successColor = color.NRGBA{R: 0x43, G: 0xf4, B: 0x36, A: 0xff}
	warningColor = color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff}
)

type builtinTheme struct {
	variant fyne.ThemeVariant

	regular, bold, italic, boldItalic, monospace fyne.Resource
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
		if t.regular == regular { // failed to load
			t.bold = loadCustomFont(font, "Bold", bold)
			t.italic = loadCustomFont(font, "Italic", italic)
			t.boldItalic = loadCustomFont(font, "BoldItalic", bolditalic)
		} else { // first custom font loaded, fall back to that
			t.bold = loadCustomFont(font, "Bold", t.regular)
			t.italic = loadCustomFont(font, "Italic", t.regular)
			t.boldItalic = loadCustomFont(font, "BoldItalic", t.regular)
		}
	}
	font = os.Getenv("FYNE_FONT_MONOSPACE")
	if font != "" {
		t.monospace = loadCustomFont(font, "Regular", monospace)
	}
}

func (t *builtinTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if t.variant != variantNameUserPreference {
		v = t.variant
	}

	primary := fyne.CurrentApp().Settings().PrimaryColor()
	if n == ColorNamePrimary {
		return primaryColorNamed(primary)
	} else if n == ColorNameFocus {
		return focusColorNamed(primary)
	} else if n == ColorNameSelection {
		return selectionColorNamed(primary)
	}

	if v == VariantLight {
		return lightPaletColorNamed(n)
	}

	return darkPaletColorNamed(n)
}

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
	case SizeNameSeparatorThickness:
		return 1
	case SizeNameInlineIcon:
		return 20
	case SizeNameInnerPadding:
		return 8
	case SizeNameLineSpacing:
		return 4
	case SizeNamePadding:
		return 6
	case SizeNameScrollBar:
		return 16
	case SizeNameScrollBarSmall:
		return 3
	case SizeNameText:
		return 13
	case SizeNameHeadingText:
		return 24
	case SizeNameSubHeadingText:
		return 18
	case SizeNameCaptionText:
		return 11
	case SizeNameInputBorder:
		return 1
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

func currentVariant() fyne.ThemeVariant {
	if std, ok := current().(*builtinTheme); ok {
		if std.variant != variantNameUserPreference {
			return std.variant // override if using the old LightTheme() or DarkTheme() constructor
		}
	}

	return fyne.CurrentApp().Settings().ThemeVariant()
}

func darkPaletColorNamed(name fyne.ThemeColorName) color.Color {
	switch name {
	case ColorNameBackground:
		return color.NRGBA{R: 0x14, G: 0x14, B: 0x15, A: 0xff}
	case ColorNameButton:
		return color.NRGBA{R: 0x28, G: 0x29, B: 0x2e, A: 0xff}
	case ColorNameDisabled:
		return color.NRGBA{R: 0x39, G: 0x39, B: 0x3a, A: 0xff}
	case ColorNameDisabledButton:
		return color.NRGBA{R: 0x28, G: 0x29, B: 0x2e, A: 0xff}
	case ColorNameError:
		return errorColor
	case ColorNameForeground:
		return color.NRGBA{R: 0xf3, G: 0xf3, B: 0xf3, A: 0xff}
	case ColorNameHover:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0f}
	case ColorNameInputBackground:
		return color.NRGBA{R: 0x20, G: 0x20, B: 0x23, A: 0xff}
	case ColorNameInputBorder:
		return color.NRGBA{R: 0x39, G: 0x39, B: 0x3a, A: 0xff}
	case ColorNameMenuBackground:
		return color.NRGBA{R: 0x28, G: 0x29, B: 0x2e, A: 0xff}
	case ColorNameOverlayBackground:
		return color.NRGBA{R: 0x18, G: 0x1d, B: 0x25, A: 0xff}
	case ColorNamePlaceHolder:
		return color.NRGBA{R: 0xb2, G: 0xb2, B: 0xb2, A: 0xff}
	case ColorNamePressed:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x66}
	case ColorNameScrollBar:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x99}
	case ColorNameSeparator:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff}
	case ColorNameShadow:
		return color.NRGBA{A: 0x66}
	case ColorNameSuccess:
		return successColor
	case ColorNameWarning:
		return warningColor
	}

	return color.Transparent
}

func focusColorNamed(name string) color.NRGBA {
	switch name {
	case ColorRed:
		return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0x7f}
	case ColorOrange:
		return color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0x7f}
	case ColorYellow:
		return color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0x7f}
	case ColorGreen:
		return color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0x7f}
	case ColorPurple:
		return color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0x7f}
	case ColorBrown:
		return color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0x7f}
	case ColorGray:
		return color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0x7f}
	}

	// We return the value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return color.NRGBA{R: 0x00, G: 0x6C, B: 0xff, A: 0x2a}
}

func lightPaletColorNamed(name fyne.ThemeColorName) color.Color {
	switch name {
	case ColorNameBackground:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	case ColorNameButton:
		return color.NRGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff}
	case ColorNameDisabled:
		return color.NRGBA{R: 0xe3, G: 0xe3, B: 0xe3, A: 0xff}
	case ColorNameDisabledButton:
		return color.NRGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff}
	case ColorNameError:
		return errorColor
	case ColorNameForeground:
		return color.NRGBA{R: 0x56, G: 0x56, B: 0x56, A: 0xff}
	case ColorNameHover:
		return color.NRGBA{A: 0x0f}
	case ColorNameInputBackground:
		return color.NRGBA{R: 0xf3, G: 0xf3, B: 0xf3, A: 0xff}
	case ColorNameInputBorder:
		return color.NRGBA{R: 0xe3, G: 0xe3, B: 0xe3, A: 0xff}
	case ColorNameMenuBackground:
		return color.NRGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff}
	case ColorNameOverlayBackground:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	case ColorNamePlaceHolder:
		return color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff}
	case ColorNamePressed:
		return color.NRGBA{A: 0x19}
	case ColorNameScrollBar:
		return color.NRGBA{A: 0x99}
	case ColorNameSeparator:
		return color.NRGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff}
	case ColorNameShadow:
		return color.NRGBA{A: 0x33}
	case ColorNameSuccess:
		return successColor
	case ColorNameWarning:
		return warningColor
	}

	return color.Transparent
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

func primaryColorNamed(name string) color.NRGBA {
	switch name {
	case ColorRed:
		return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
	case ColorOrange:
		return color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff}
	case ColorYellow:
		return color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0xff}
	case ColorGreen:
		return color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0xff}
	case ColorPurple:
		return color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0xff}
	case ColorBrown:
		return color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0xff}
	case ColorGray:
		return color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0xff}
	}

	// We return the value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return color.NRGBA{R: 0x29, G: 0x6f, B: 0xf6, A: 0xff}
}

func safeColorLookup(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	col := current().Color(n, v)
	if col == nil {
		fyne.LogError("Loaded theme returned nil color", nil)
		return fallbackColor
	}
	return col
}

func safeFontLookup(s fyne.TextStyle) fyne.Resource {
	font := current().Font(s)
	if font != nil {
		return font
	}
	fyne.LogError("Loaded theme returned nil font", nil)

	if s.Monospace {
		return DefaultTextMonospaceFont()
	}
	if s.Bold {
		if s.Italic {
			return DefaultTextBoldItalicFont()
		}
		return DefaultTextBoldFont()
	}
	if s.Italic {
		return DefaultTextItalicFont()
	}

	return DefaultTextFont()
}

func selectionColorNamed(name string) color.NRGBA {
	switch name {
	case ColorRed:
		return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0x3f}
	case ColorOrange:
		return color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0x3f}
	case ColorYellow:
		return color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0x3f}
	case ColorGreen:
		return color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0x3f}
	case ColorPurple:
		return color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0x3f}
	case ColorBrown:
		return color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0x3f}
	case ColorGray:
		return color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0x3f}
	}

	// We return the value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return color.NRGBA{R: 0x00, G: 0x6C, B: 0xff, A: 0x40}
}

func setupDefaultTheme() fyne.Theme {
	theme := &builtinTheme{variant: variantNameUserPreference}

	theme.initFonts()
	return theme
}
