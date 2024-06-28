// Package theme defines how a Fyne app should look when rendered.
package theme // import "fyne.io/fyne/v2/theme"

import (
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	internaltheme "fyne.io/fyne/v2/internal/theme"
)

// Keep in mind to add new constants to the tests at test/theme.go.
const (
	// VariantDark is the version of a theme that satisfies a user preference for a dark look.
	//
	// Since: 2.0
	VariantDark = internaltheme.VariantDark

	// VariantLight is the version of a theme that satisfies a user preference for a light look.
	//
	// Since: 2.0
	VariantLight = internaltheme.VariantLight
)

var defaultTheme fyne.Theme

// DarkTheme defines the built-in dark theme colors and sizes.
//
// Deprecated: This method ignores user preference and should not be used, it will be removed in v3.0.
// If developers want to ignore user preference for theme variant they can set a custom theme.
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
// If developers want to ignore user preference for theme variant they can set a custom theme.
func LightTheme() fyne.Theme {
	theme := &builtinTheme{variant: VariantLight}

	theme.initFonts()
	return theme
}

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
	if t.variant != internaltheme.VariantNameUserPreference {
		v = t.variant
	}

	primary := fyne.CurrentApp().Settings().PrimaryColor()
	if n == ColorNamePrimary || n == ColorNameHyperlink {
		return internaltheme.PrimaryColorNamed(primary)
	} else if n == ColorNameForegroundOnPrimary {
		return internaltheme.ForegroundOnPrimaryColorNamed(primary)
	} else if n == ColorNameFocus {
		return focusColorNamed(primary)
	} else if n == ColorNameSelection {
		return selectionColorNamed(primary)
	}

	if v == VariantLight {
		return lightPaletteColorNamed(n)
	}

	return darkPaletteColorNamed(n)
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
		return 12
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
	case SizeNameScrollBarRadius:
		return 3
	default:
		return 0
	}
}

// Current returns the theme that is currently used for the running application.
// It looks up based on user preferences and application configuration.
//
// Since: 2.5
func Current() fyne.Theme {
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

// CurrentForWidget returns the theme that is currently used for the specified widget.
// It looks for widget overrides and falls back to the application's current theme.
//
// Since: 2.5
func CurrentForWidget(w fyne.CanvasObject) fyne.Theme {
	if custom := cache.WidgetTheme(w); custom != nil {
		return custom
	}

	return Current()
}

func currentVariant() fyne.ThemeVariant {
	if std, ok := Current().(*builtinTheme); ok {
		if std.variant != internaltheme.VariantNameUserPreference {
			return std.variant // override if using the old LightTheme() or DarkTheme() constructor
		}
	}

	return fyne.CurrentApp().Settings().ThemeVariant()
}

func darkPaletteColorNamed(name fyne.ThemeColorName) color.Color {
	switch name {
	case ColorNameBackground:
		return colorDarkBackground
	case ColorNameButton:
		return colorDarkButton
	case ColorNameDisabled:
		return colorDarkDisabled
	case ColorNameDisabledButton:
		return colorDarkDisabledButton
	case ColorNameError:
		return colorDarkError
	case ColorNameForeground:
		return colorDarkForeground
	case ColorNameForegroundOnError:
		return colorDarkForegroundOnError
	case ColorNameForegroundOnSuccess:
		return colorDarkForegroundOnSuccess
	case ColorNameForegroundOnWarning:
		return colorDarkForegroundOnWarning
	case ColorNameHover:
		return colorDarkHover
	case ColorNameHeaderBackground:
		return colorDarkHeaderBackground
	case ColorNameInputBackground:
		return colorDarkInputBackground
	case ColorNameInputBorder:
		return colorDarkInputBorder
	case ColorNameMenuBackground:
		return colorDarkMenuBackground
	case ColorNameOverlayBackground:
		return colorDarkOverlayBackground
	case ColorNamePlaceHolder:
		return colorDarkPlaceholder
	case ColorNamePressed:
		return colorDarkPressed
	case ColorNameScrollBar:
		return colorDarkScrollBar
	case ColorNameSeparator:
		return colorDarkSeparator
	case ColorNameShadow:
		return colorDarkShadow
	case ColorNameSuccess:
		return colorDarkSuccess
	case ColorNameWarning:
		return colorDarkWarning
	}

	return color.Transparent
}

func focusColorNamed(name string) color.NRGBA {
	switch name {
	case ColorRed:
		return colorLightFocusRed
	case ColorOrange:
		return colorLightFocusOrange
	case ColorYellow:
		return colorLightFocusYellow
	case ColorGreen:
		return colorLightFocusGreen
	case ColorPurple:
		return colorLightFocusPurple
	case ColorBrown:
		return colorLightFocusBrown
	case ColorGray:
		return colorLightFocusGray
	}

	// We return the value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return colorLightFocusBlue
}

func lightPaletteColorNamed(name fyne.ThemeColorName) color.Color {
	switch name {
	case ColorNameBackground:
		return colorLightBackground
	case ColorNameButton:
		return colorLightButton
	case ColorNameDisabled:
		return colorLightDisabled
	case ColorNameDisabledButton:
		return colorLightDisabledButton
	case ColorNameError:
		return colorLightError
	case ColorNameForeground:
		return colorLightForeground
	case ColorNameForegroundOnError:
		return colorLightForegroundOnError
	case ColorNameForegroundOnSuccess:
		return colorLightForegroundOnSuccess
	case ColorNameForegroundOnWarning:
		return colorLightForegroundOnWarning
	case ColorNameHover:
		return colorLightHover
	case ColorNameHeaderBackground:
		return colorLightHeaderBackground
	case ColorNameInputBackground:
		return colorLightInputBackground
	case ColorNameInputBorder:
		return colorLightInputBorder
	case ColorNameMenuBackground:
		return colorLightMenuBackground
	case ColorNameOverlayBackground:
		return colorLightOverlayBackground
	case ColorNamePlaceHolder:
		return colorLightPlaceholder
	case ColorNamePressed:
		return colorLightPressed
	case ColorNameScrollBar:
		return colorLightScrollBar
	case ColorNameSeparator:
		return colorLightSeparator
	case ColorNameShadow:
		return colorLightShadow
	case ColorNameSuccess:
		return colorLightSuccess
	case ColorNameWarning:
		return colorLightWarning
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

func selectionColorNamed(name string) color.NRGBA {
	switch name {
	case ColorRed:
		return colorLightSelectionRed
	case ColorOrange:
		return colorLightSelectionOrange
	case ColorYellow:
		return colorLightSelectionYellow
	case ColorGreen:
		return colorLightSelectionGreen
	case ColorPurple:
		return colorLightSelectionPurple
	case ColorBrown:
		return colorLightSelectionBrown
	case ColorGray:
		return colorLightSelectionGray
	}

	// We return the value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return colorLightSelectionBlue
}

func setupDefaultTheme() fyne.Theme {
	theme := &builtinTheme{variant: internaltheme.VariantNameUserPreference}

	theme.initFonts()
	return theme
}
