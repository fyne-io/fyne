package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Primary color names.
const (
	ColorBlue   = "blue"
	ColorBrown  = "brown"
	ColorGray   = "gray"
	ColorGreen  = "green"
	ColorOrange = "orange"
	ColorPurple = "purple"
	ColorRed    = "red"
	ColorYellow = "yellow"
)

// Theme variants; the public available ones are defined in /theme/theme.go.
const (
	VariantDark fyne.ThemeVariant = iota
	VariantLight
	VariantNameUserPreference // locally used in builtinTheme for backward compatibility
)

var (
	colorLightOnPrimaryBlue   = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	colorLightOnPrimaryBrown  = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	colorLightOnPrimaryGray   = color.NRGBA{R: 0x17, G: 0x17, B: 0x18, A: 0xff}
	colorLightOnPrimaryGreen  = color.NRGBA{R: 0x17, G: 0x17, B: 0x18, A: 0xff}
	colorLightOnPrimaryOrange = color.NRGBA{R: 0x17, G: 0x17, B: 0x18, A: 0xff}
	colorLightOnPrimaryPurple = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	colorLightOnPrimaryRed    = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	colorLightOnPrimaryYellow = color.NRGBA{R: 0x17, G: 0x17, B: 0x18, A: 0xff}
	colorLightPrimaryBlue     = color.NRGBA{R: 0x29, G: 0x6f, B: 0xf6, A: 0xff}
	colorLightPrimaryBrown    = color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0xff}
	colorLightPrimaryGray     = color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0xff}
	colorLightPrimaryGreen    = color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0xff}
	colorLightPrimaryOrange   = color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff}
	colorLightPrimaryPurple   = color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0xff}
	colorLightPrimaryRed      = color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
	colorLightPrimaryYellow   = color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0xff}
)

// ForegroundOnPrimaryColorNamed returns a theme specific color used for text and icons against the named primary color.
func ForegroundOnPrimaryColorNamed(name string) color.Color {
	switch name {
	case ColorRed:
		return colorLightOnPrimaryRed
	case ColorOrange:
		return colorLightOnPrimaryOrange
	case ColorYellow:
		return colorLightOnPrimaryYellow
	case ColorGreen:
		return colorLightOnPrimaryGreen
	case ColorPurple:
		return colorLightOnPrimaryPurple
	case ColorBrown:
		return colorLightOnPrimaryBrown
	case ColorGray:
		return colorLightOnPrimaryGray
	}

	// We return the “on” value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return colorLightOnPrimaryBlue
}

// PrimaryColorNamed returns a theme specific color value for a named primary color.
func PrimaryColorNamed(name string) color.Color {
	switch name {
	case ColorRed:
		return colorLightPrimaryRed
	case ColorOrange:
		return colorLightPrimaryOrange
	case ColorYellow:
		return colorLightPrimaryYellow
	case ColorGreen:
		return colorLightPrimaryGreen
	case ColorPurple:
		return colorLightPrimaryPurple
	case ColorBrown:
		return colorLightPrimaryBrown
	case ColorGray:
		return colorLightPrimaryGray
	}

	// We return the value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return colorLightPrimaryBlue
}
