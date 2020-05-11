// Package theme defines how a Fyne app should look when rendered
package theme // import "fyne.io/fyne/theme"

import (
	"image/color"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne"
)

type builtinTheme struct {
	palette                                      *palette
	regular, bold, italic, bolditalic, monospace fyne.Resource
	direction float64
}

// LightTheme defines the built in light theme colours and sizes
func LightTheme() fyne.Theme {
	theme := &builtinTheme{direction: 1.0}
	theme.palette = defaultPalette()
	theme.initFonts()
	return theme
}

// DarkTheme defines the built in dark theme colours and sizes
// See https://www.material.io/design/color/dark-theme.html
func DarkTheme() fyne.Theme {
	theme := &builtinTheme{direction: -1.0}
	var err error
	theme.palette, err = newPalette(defaultPalette(), `
{
	"palette": {
		"canvasColor": "#242424",
		"textColor": "#aaaaaa",
		"secondaryTextColor": "rgba(255, 255, 255, 0.7)",
		"textColor": "rgba(255, 255, 255, 0.9)"
	}
}
	`)
	if err != nil {
		println(err.Error())
	}
	theme.initFonts()
	return theme
}

func (t *builtinTheme) BackgroundColor() color.Color {
	return t.palette.canvasColor
}

// ButtonColor returns the theme's standard button colour
func (t *builtinTheme) ButtonColor() color.Color {
	return brighten(t.palette.canvasColor, -2.4 * t.direction)
}

// DisabledButtonColor returns the theme's disabled button colour
func (t *builtinTheme) DisabledButtonColor() color.Color {
	return brighten(t.palette.canvasColor, -1.2 * t.direction)
}

// HyperlinkColor returns the theme's standard hyperlink colour
func (t *builtinTheme) HyperlinkColor() color.Color {
	return color.NRGBA{0x99, 0x99, 0xff, 0xff}
	//return t.hyperlink
}

// TextColor returns the theme's standard text colour
func (t *builtinTheme) TextColor() color.Color {
	return t.palette.textColor
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *builtinTheme) DisabledTextColor() color.Color {
	return t.palette.disabledColor
}

// PrimaryTextColor returns the theme's standard text colour on primary background
func (t *builtinTheme) PrimaryTextColor() color.Color {
	return t.palette.canvasColor
}

// IconColor returns the theme's standard text colour
func (t *builtinTheme) IconColor() color.Color {
	return brighten(t.palette.textColor, 0.8 * t.direction)
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *builtinTheme) DisabledIconColor() color.Color {
	return brighten(t.palette.disabledColor, 0.8 * t.direction)
}

// PlaceHolderColor returns the theme's placeholder text colour
func (t *builtinTheme) PlaceHolderColor() color.Color {
	return t.palette.disabledColor
}

// PrimaryColor returns the colour used to highlight primary features
func (t *builtinTheme) PrimaryColor() color.Color {
	return t.palette.primary1Color
}

// PrimaryHoverColor returns the colour used to highlight primary features on hover
func (t *builtinTheme) PrimaryHoverColor() color.Color {
	return brighten(t.palette.primary1Color, 0.8 * t.direction)
}

// HoverColor returns the colour used to highlight interactive elements currently under a cursor
func (t *builtinTheme) HoverColor() color.Color {
	return brighten(t.palette.canvasColor, -1.0 * t.direction)
}

// FocusColor returns the colour used to highlight a focused widget
func (t *builtinTheme) FocusColor() color.Color {
	return t.palette.primary1Color
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func (t *builtinTheme) ScrollBarColor() color.Color {
	return t.palette.primary1Color
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation
func (t *builtinTheme) ShadowColor() color.Color {
	return t.palette.shadowColor
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

// Extend will apply the new palette from the input JSON, or return error
func Extend(str string) error { 
	c := current()
	t,ok := c.(*builtinTheme)
	if !ok {
		fmt.Printf("Unknown type %T for current theme", c)
		fyne.CurrentApp().Settings().SetTheme(DarkTheme())
		t = current().(*builtinTheme)
	}
	newPalette,err := newPalette(t.palette, str)
	if err != nil {
		return err
	}
	theme := &builtinTheme{}
	theme.palette = newPalette
	theme.initFonts()
	fyne.CurrentApp().Settings().SetTheme(theme)
	return nil
}

// BackgroundColor returns the theme's background colour

func BackgroundColor() color.Color {
	return current().BackgroundColor()
}

// ButtonColor returns the theme's standard button colour
func ButtonColor() color.Color {
	return current().ButtonColor()
}

// DisabledButtonColor returns the theme's disabled button colour
func DisabledButtonColor() color.Color {
	return current().DisabledButtonColor()
}

// HyperlinkColor returns the theme's standard hyperlink colour
func HyperlinkColor() color.Color {
	return current().HyperlinkColor()
}

// TextColor returns the theme's standard text colour
func TextColor() color.Color {
	return current().TextColor()
}

// DisabledTextColor returns the color for a disabledIcon UI element
func DisabledTextColor() color.Color {
	return current().DisabledTextColor()
}

// PrimaryTextColor returns the color for text on a primary background
func PrimaryTextColor() color.Color {
	return current().PrimaryTextColor()
}

// IconColor returns the theme's standard icon colour
func IconColor() color.Color {
	return current().IconColor()
}

// DisabledIconColor returns the color for a disabledIcon UI element
func DisabledIconColor() color.Color {
	return current().DisabledIconColor()
}

// PlaceHolderColor returns the theme's standard text colour
func PlaceHolderColor() color.Color {
	return current().PlaceHolderColor()
}

// PrimaryColor returns the colour used to highlight primary features
func PrimaryColor() color.Color {
	return current().PrimaryColor()
}

// PrimaryHoverColor returns the colour used to highlight primary features on hover
func PrimaryHoverColor() color.Color {
	return current().PrimaryHoverColor()
}

// HoverColor returns the colour used to highlight interactive elements currently under a cursor
func HoverColor() color.Color {
	return current().HoverColor()
}

// FocusColor returns the colour used to highlight a focussed widget
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

// DefaultTextMonospaceFont returns the font resource for the built-in monospace font face
func DefaultTextMonospaceFont() fyne.Resource {
	return monospace
}
