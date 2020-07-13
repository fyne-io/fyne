package theme

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne"
)

// TappedShade is the shade used for pressed buttons, in %
var TappedShade = 14

// HoveredShade is the shade used for hovered buttons, in %
var HoveredShade = 8

// FocusedShade is the shade used for focused buttons, in %.
var FocusedShade = 12

// Palette is the main color palette acording to Google recomendations.
type Palette struct {
	Primary      color.NRGBA
	Secondary    color.NRGBA
	Background   color.NRGBA
	Surface      color.NRGBA
	Error        color.NRGBA
	OnPrimary    color.NRGBA
	OnSecondary  color.NRGBA
	OnBackground color.NRGBA
	OnSurface    color.NRGBA
	OnError      color.NRGBA
	Shadow       color.NRGBA
	Hyperlink    color.NRGBA
	DisabledText color.NRGBA
	Disabled     color.NRGBA
}

// paletteDescription is a struct used to parse json files
// The fields will receive data with corresponding names
type paletteDescription struct {
	Primary      string
	Secondary    string
	Background   string
	Surface      string
	Error        string
	OnPrimary    string
	OnSecondary  string
	OnBackground string
	OnSurface    string
	OnError      string
	Shadow       string
	DisabledText string
	Disabled     string
	// Fyne special colors
	Hyperlink string
	Button    string
}

// AlternativePalette is the default light Google Material design Palette
// see: https://material.io/design/color/the-color-system.html#color-theme-creation
const AlternativePalette string = `
{
	"Palette": {
		"Background": "#ffffff",
		"OnBackground": "#000000",
		"Primary": "#6200ee",
		"OnPrimary": "#ffffff",
		"Secondary": "#03dac6",
		"OnSecondary": "#000000",
		"Surface": "#dddddd",
		"OnSurface": "#000000",
		"Error": "#b00020",
        "OnError": "#ffffff",
		"Shadow": "rgba(0, 0, 0, 0.2}",
		"Disabled": "#e7e7e7",
		"DisabledText": "#808080",
		"Hyperlink": "#0000d9",
		"Button": "#222222"
	}
}
`

// YellowAlternative is a strikingly different theme
const YellowAlternative string = `
{
	"Palette": {
		"Background": "#ffde30",
		"OnBackground": "#000000",
		"Primary": "#ffff22",
		"OnPrimary": "#000044",
		"Secondary": "#aaaa00",
		"OnSecondary": "#ffffff",
		"Surface": "#ffee77",
		"OnSurface": "#000000",
		"Error": "#b00020",
        "OnError": "#ffffff",
		"Shadow": "rgba(0, 0, 0, 0.2}",
		"Disabled": "#e7e7e7",
		"DisabledText": "#808080",
		"Hyperlink": "#0000d9",
		"Button": "#eeee20"
	}
}
`

// ParsePallete converts string to a builtin theme
func ParsePallete(s string) fyne.Theme {
	t, _ := newPalette(s)
	t.initFonts()
	return t
}

// newPalette will parse the string (containing json text)
// and generate a pallete struct
func newPalette(str string) (*builtinTheme, error) {
	type exportedPalette struct {
		Palette paletteDescription
	}
	ep := &exportedPalette{}
	np := &builtinTheme{}
	if err := json.Unmarshal([]byte(str), ep); err != nil {
		return nil, err
	}
	p := ep.Palette
	np.Background = rgb(np.Background, p.Background)
	np.OnBackground = rgb(np.OnBackground, p.OnBackground)
	np.Primary = rgb(np.Primary, p.Primary)
	np.OnPrimary = rgb(np.OnBackground, p.OnPrimary)
	np.Secondary = rgb(np.Secondary, p.Secondary)
	np.OnSecondary = rgb(np.OnBackground, p.OnSecondary)
	np.Surface = rgb(np.Background, p.Surface)
	np.OnSurface = rgb(np.OnBackground, p.OnSurface)
	np.Error = rgb(np.Error, p.Error)
	np.OnError = rgb(np.OnBackground, p.OnError)
	np.Disabled = rgb(np.Disabled, p.Disabled)
	np.DisabledText = rgb(np.DisabledText, p.DisabledText)
	np.Shadow = rgb(np.Shadow, p.Shadow)
	np.Hyperlink = rgb(np.Hyperlink, p.Hyperlink)
	np.placeholder = np.DisabledText
	np.button = rgb(np.Surface, p.Button)
	np.text = np.OnBackground
	np.hover = Blend(np.ButtonColor(), np.TextColor(), 0.1)
	np.icon = np.OnBackground
	np.hyperlink = np.Hyperlink
	np.primary = np.Primary
	np.background = np.Background
	np.disabledButton = np.Disabled
	np.disabledText = np.DisabledText
	return np, nil
}

func rgb(c color.NRGBA, s string) color.NRGBA {
	if strings.HasPrefix(s, "#") && len(s) == 7 {
		fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
		c.A = 255
	} else if strings.HasPrefix(s, "rgba(") {
		s = strings.Replace(s, " ", "", -1)
		alpha := 1.0
		fmt.Sscanf(s, "rgba(%d,%d,%d,%f)", &c.R, &c.G, &c.B, &alpha)
		c.A = uint8(255.0 * alpha)
	} else if s != "" {
		panic("Palette colors has wrong format, should start with # or rgba( : " + s)
	}
	return c
}

// Blend will return a new color that is a alpha blend of the two given colors
// Alpha=1.0 wil return the second color, alpha=0.0 will return the first color
// I.e. the second color is blended into the first color by the given alpha factor
func Blend(c1 color.Color, c2 color.Color, alpha float64) color.Color {
	m := uint32(alpha * 255 * 255)
	c1r, c1g, c1b, _ := c1.RGBA()
	c2r, c2g, c2b, _ := c2.RGBA()
	const M = (1<<16 - 1) // = 65535 * 256
	return color.RGBA{
		uint8(((c1r*(M-m) + c2r*m) / M) >> 8),
		uint8(((c1g*(M-m) + c2g*m) / M) >> 8),
		uint8(((c1b*(M-m) + c2b*m) / M) >> 8),
		255}
}

// This init is needed because the tests depend on a default theme.
func init() {
	UpdateDefaultStyling(DarkTheme())
}
