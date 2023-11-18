// Package theme defines how a Fyne app should look when rendered.
package theme // import "fyne.io/fyne/v2/theme"

import (
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
)

const (
	// VariantDark is the version of a theme that satisfies a user preference for a dark look.
	//
	// Since: 2.0
	VariantDark fyne.ThemeVariant = 0

	// VariantLight is the version of a theme that satisfies a user preference for a light look.
	//
	// Since: 2.0
	VariantLight fyne.ThemeVariant = 1

	// potential for adding theme types such as high visibility or monochrome...
	variantNameUserPreference fyne.ThemeVariant = 2 // locally used in builtinTheme for backward compatibility
)

// DarkTheme defines the built-in dark theme colors and sizes.
//
// Deprecated: This method ignores user preference and should not be used, it will be removed in v3.0.
func DarkTheme() fyne.Theme {
	theme := &builtinTheme{variant: VariantDark}

	theme.initFonts()
	return theme
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

// LightTheme defines the built-in light theme colors and sizes.
//
// Deprecated: This method ignores user preference and should not be used, it will be removed in v3.0.
func LightTheme() fyne.Theme {
	theme := &builtinTheme{variant: VariantLight}

	theme.initFonts()
	return theme
}

var (
	defaultTheme fyne.Theme
)

type builtinTheme struct {
	variant fyne.ThemeVariant

	regular, bold, italic, boldItalic, monospace, symbol fyne.Resource
}

func (t *builtinTheme) initFonts() {
	t.regular = regular
	t.bold = bold
	t.italic = italic
	t.boldItalic = bolditalic
	t.monospace = monospace
	t.symbol = symbol

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
	font = os.Getenv("FYNE_FONT_SYMBOL")
	if font != "" {
		t.symbol = loadCustomFont(font, "Regular", symbol)
	}
}

func (t *builtinTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if t.variant != variantNameUserPreference {
		v = t.variant
	}

	primary := fyne.CurrentApp().Settings().PrimaryColor()
	if n == ColorNamePrimary || n == ColorNameHyperlink {
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
	if style.Symbol {
		return t.symbol
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
		return 4
	case SizeNameScrollBar:
		return 16
	case SizeNameScrollBarSmall:
		return 3
	case SizeNameText:
		return 14
	case SizeNameHeadingText:
		return 24
	case SizeNameSubHeadingText:
		return 18
	case SizeNameCaptionText:
		return 11
	case SizeNameInputBorder:
		return 1
	case SizeNameInputRadius:
		return 5
	case SizeNameSelectionRadius:
		return 3
	default:
		return 0
	}
}

func current() fyne.Theme {
	app := fyne.CurrentApp()
	if app == nil {
		return DarkTheme()
	}
	currentTheme := app.Settings().Theme()
	if currentTheme == nil {
		return DarkTheme()
	}

	return currentTheme
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
		return color.NRGBA{R: 0x17, G: 0x17, B: 0x18, A: 0xff}
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
	case ColorNameHeaderBackground:
		return color.NRGBA{R: 0x1b, G: 0x1b, B: 0x1b, A: 0xff}
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
	case ColorNameHeaderBackground:
		return color.NRGBA{R: 0xf9, G: 0xf9, B: 0xf9, A: 0xff}
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
		return color.NRGBA{R: 0xe3, G: 0xe3, B: 0xe3, A: 0xff}
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
