package test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var defaultTheme fyne.Theme

// Try to keep these in sync with the existing color names at theme/color.go.
var knownColorNames = [...]fyne.ThemeColorName{
	theme.ColorNameBackground,
	theme.ColorNameButton,
	theme.ColorNameDisabled,
	theme.ColorNameDisabledButton,
	theme.ColorNameError,
	theme.ColorNameFocus,
	theme.ColorNameForeground,
	theme.ColorNameForegroundOnError,
	theme.ColorNameForegroundOnPrimary,
	theme.ColorNameForegroundOnSuccess,
	theme.ColorNameForegroundOnWarning,
	theme.ColorNameHeaderBackground,
	theme.ColorNameHover,
	theme.ColorNameHyperlink,
	theme.ColorNameInputBackground,
	theme.ColorNameInputBorder,
	theme.ColorNameMenuBackground,
	theme.ColorNameOverlayBackground,
	theme.ColorNamePlaceHolder,
	theme.ColorNamePressed,
	theme.ColorNamePrimary,
	theme.ColorNameScrollBar,
	theme.ColorNameScrollBarBackground,
	theme.ColorNameSelection,
	theme.ColorNameSeparator,
	theme.ColorNameShadow,
	theme.ColorNameSuccess,
	theme.ColorNameWarning,
}

// AssertAllColorNamesDefined asserts that all known color names are defined for the given theme.
func AssertAllColorNamesDefined(t *testing.T, th fyne.Theme, themeName string) {
	oldApp := fyne.CurrentApp()
	defer fyne.SetCurrentApp(oldApp)

	for _, primaryName := range theme.PrimaryColorNames() {
		testApp := NewTempApp(t)
		testApp.Settings().(*testSettings).primaryColor = primaryName
		for variantName, variant := range KnownThemeVariants() {
			for _, cn := range knownColorNames {
				assert.NotNil(t, th.Color(cn, variant), "undefined color %s variant %s in theme %s", cn, variantName, themeName)
				// Transparent is used by the default theme as fallback for unknown color names.
				// Built-in color names should have well-defined non-transparent values.
				assert.NotEqual(t, color.Transparent, th.Color(cn, variant), "undefined color %s variant %s in theme %s", cn, variantName, themeName)
			}
		}
	}
}

// KnownThemeVariants returns the known theme variants mapped by a descriptive key.
func KnownThemeVariants() map[string]fyne.ThemeVariant {
	// Try to keep this in sync with the existing variants at theme/theme.go
	return map[string]fyne.ThemeVariant{
		"dark":  theme.VariantDark,
		"light": theme.VariantLight,
	}
}

// NewTheme returns a new test theme using quiet ugly colors.
func NewTheme() fyne.Theme {
	blue := func(alpha uint8) color.Color {
		return &color.NRGBA{R: 0, G: 0, B: 255, A: alpha}
	}
	gray := func(level uint8) color.Color {
		return &color.Gray{Y: level}
	}
	green := func(alpha uint8) color.Color {
		return &color.NRGBA{R: 0, G: 255, B: 0, A: alpha}
	}
	red := func(alpha uint8) color.Color {
		return &color.NRGBA{R: 200, G: 0, B: 0, A: alpha}
	}

	return &configurableTheme{
		colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:          red(255),
			theme.ColorNameButton:              gray(100),
			theme.ColorNameDisabled:            gray(20),
			theme.ColorNameDisabledButton:      gray(230),
			theme.ColorNameError:               blue(255),
			theme.ColorNameFocus:               red(66),
			theme.ColorNameForeground:          gray(255),
			theme.ColorNameForegroundOnError:   red(210),
			theme.ColorNameForegroundOnPrimary: red(200),
			theme.ColorNameForegroundOnSuccess: blue(201),
			theme.ColorNameForegroundOnWarning: blue(202),
			theme.ColorNameHeaderBackground:    red(22),
			theme.ColorNameHover:               green(200),
			theme.ColorNameHyperlink:           blue(240),
			theme.ColorNameInputBackground:     red(30),
			theme.ColorNameInputBorder:         gray(10),
			theme.ColorNameMenuBackground:      red(50),
			theme.ColorNameOverlayBackground:   red(44),
			theme.ColorNamePlaceHolder:         blue(200),
			theme.ColorNamePressed:             blue(250),
			theme.ColorNamePrimary:             green(255),
			theme.ColorNameScrollBar:           blue(220),
			theme.ColorNameScrollBarBackground: red(20),
			theme.ColorNameSelection:           red(55),
			theme.ColorNameSeparator:           gray(30),
			theme.ColorNameShadow:              blue(150),
			theme.ColorNameSuccess:             green(150),
			theme.ColorNameWarning:             red(100),
		},
		fonts: map[fyne.TextStyle]fyne.Resource{
			{}:                         theme.DefaultTextBoldFont(),
			{Bold: true}:               theme.DefaultTextItalicFont(),
			{Bold: true, Italic: true}: theme.DefaultTextMonospaceFont(),
			{Italic: true}:             theme.DefaultTextBoldItalicFont(),
			{Monospace: true}:          theme.DefaultTextFont(),
			{Symbol: true}:             theme.DefaultSymbolFont(),
		},
		name: "Ugly Test Theme",
		sizes: map[fyne.ThemeSizeName]float32{
			theme.SizeNameInlineIcon:         float32(24),
			theme.SizeNameInnerPadding:       float32(20),
			theme.SizeNameLineSpacing:        float32(6),
			theme.SizeNamePadding:            float32(10),
			theme.SizeNameScrollBar:          float32(10),
			theme.SizeNameScrollBarSmall:     float32(2),
			theme.SizeNameSeparatorThickness: float32(1),
			theme.SizeNameText:               float32(18),
			theme.SizeNameHeadingText:        float32(30.6),
			theme.SizeNameSubHeadingText:     float32(24),
			theme.SizeNameCaptionText:        float32(15),
			theme.SizeNameInputBorder:        float32(5),
			theme.SizeNameInputRadius:        float32(2),
			theme.SizeNameSelectionRadius:    float32(6),
			theme.SizeNameScrollBarRadius:    float32(2),
		},
	}
}

// Theme returns a test theme useful for image based tests.
func Theme() fyne.Theme {
	if defaultTheme == nil {
		defaultTheme = &configurableTheme{
			colors: map[fyne.ThemeColorName]color.Color{
				theme.ColorNameBackground:          color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff},
				theme.ColorNameButton:              color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
				theme.ColorNameDisabled:            color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
				theme.ColorNameDisabledButton:      color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
				theme.ColorNameError:               color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff},
				theme.ColorNameFocus:               color.NRGBA{R: 0x78, G: 0x3a, B: 0x3a, A: 0xff},
				theme.ColorNameForeground:          color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
				theme.ColorNameForegroundOnError:   color.NRGBA{R: 0x08, G: 0x0a, B: 0x0f, A: 0xff},
				theme.ColorNameForegroundOnPrimary: color.NRGBA{R: 0x08, G: 0x0c, B: 0x0f, A: 0xff},
				theme.ColorNameForegroundOnSuccess: color.NRGBA{R: 0x0a, G: 0x0c, B: 0x0f, A: 0xff},
				theme.ColorNameForegroundOnWarning: color.NRGBA{R: 0x08, G: 0x0c, B: 0x0a, A: 0xff},
				theme.ColorNameHeaderBackground:    color.NRGBA{R: 0x25, G: 0x25, B: 0x25, A: 0xff},
				theme.ColorNameHover:               color.NRGBA{R: 0x88, G: 0xff, B: 0xff, A: 0x22},
				theme.ColorNameHyperlink:           color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
				theme.ColorNameInputBackground:     color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xff},
				theme.ColorNameInputBorder:         color.NRGBA{R: 0x86, G: 0x86, B: 0x86, A: 0xff},
				theme.ColorNameMenuBackground:      color.NRGBA{R: 0x56, G: 0x56, B: 0x56, A: 0xff},
				theme.ColorNameOverlayBackground:   color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0xff},
				theme.ColorNamePlaceHolder:         color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff},
				theme.ColorNamePressed:             color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33},
				theme.ColorNamePrimary:             color.NRGBA{R: 0xff, G: 0xc0, B: 0x80, A: 0xff},
				theme.ColorNameScrollBar:           color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa},
				theme.ColorNameScrollBarBackground: color.NRGBA{R: 0x67, G: 0x66, B: 0x66, A: 0xff},
				theme.ColorNameSelection:           color.NRGBA{R: 0x78, G: 0x3a, B: 0x3a, A: 0x99},
				theme.ColorNameSeparator:           color.NRGBA{R: 0x90, G: 0x90, B: 0x90, A: 0xff},
				theme.ColorNameShadow:              color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x88},
				theme.ColorNameSuccess:             color.NRGBA{R: 0x00, G: 0x99, B: 0x00, A: 0xff},
				theme.ColorNameWarning:             color.NRGBA{R: 0xee, G: 0xee, B: 0x00, A: 0xff},
			},
			fonts: map[fyne.TextStyle]fyne.Resource{
				{}:                         theme.DefaultTextFont(),
				{Bold: true}:               theme.DefaultTextBoldFont(),
				{Bold: true, Italic: true}: theme.DefaultTextBoldItalicFont(),
				{Italic: true}:             theme.DefaultTextItalicFont(),
				{Monospace: true}:          theme.DefaultTextMonospaceFont(),
				{Symbol: true}:             theme.DefaultSymbolFont(),
			},
			name: "Default Test Theme",
			sizes: map[fyne.ThemeSizeName]float32{
				theme.SizeNameInlineIcon:           float32(20),
				theme.SizeNameInnerPadding:         float32(8),
				theme.SizeNameLineSpacing:          float32(4),
				theme.SizeNamePadding:              float32(4),
				theme.SizeNameScrollBar:            float32(16),
				theme.SizeNameScrollBarSmall:       float32(3),
				theme.SizeNameSeparatorThickness:   float32(1),
				theme.SizeNameText:                 float32(14),
				theme.SizeNameHeadingText:          float32(23.8),
				theme.SizeNameSubHeadingText:       float32(18),
				theme.SizeNameCaptionText:          float32(11),
				theme.SizeNameInputBorder:          float32(2),
				theme.SizeNameInputRadius:          float32(4),
				theme.SizeNameSelectionRadius:      float32(4),
				theme.SizeNameScrollBarRadius:      float32(3),
				theme.SizeNameWindowTitleBarHeight: float32(20),
				theme.SizeNameWindowButtonHeight:   float32(10),
				theme.SizeNameWindowButtonIcon:     float32(8),
				theme.SizeNameWindowButtonRadius:   float32(5),
			},
		}
	}
	return defaultTheme
}

type configurableTheme struct {
	colors map[fyne.ThemeColorName]color.Color
	fonts  map[fyne.TextStyle]fyne.Resource
	name   string
	sizes  map[fyne.ThemeSizeName]float32
}

var _ fyne.Theme = (*configurableTheme)(nil)

func (t *configurableTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if t.colors[n] == nil {
		fyne.LogError(fmt.Sprintf("color %s not defined in theme %s", n, t.name), nil)
	}

	return t.colors[n]
}

func (t *configurableTheme) Font(style fyne.TextStyle) fyne.Resource {
	if t.fonts[style] == nil {
		fyne.LogError(fmt.Sprintf("font for style %#v not defined in theme %s", style, t.name), nil)
	}

	return t.fonts[style]
}

func (t *configurableTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *configurableTheme) Size(s fyne.ThemeSizeName) float32 {
	if _, ok := t.sizes[s]; !ok {
		fyne.LogError(fmt.Sprintf("size %s not defined in theme %s", s, t.name), nil)
		return 0
	}

	return t.sizes[s]
}
