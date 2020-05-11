package theme

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strings"
)

func defaultPalette() *palette {
	p, _ := newPalette(nil, `
{
	"palette": {
		"textColor": "#000000",
		"canvasColor": "#ffffff",
		"primary1Color": "#9fa8da",
		"primary2Color": "#b39ddb",
		"primary2Color": "#757575",
		"accent1Color": "#00bfa5",
		"accent2Color": "#b9f6ca",
		"accent3Color": "#ccff90"
	}
}
`)
	return p
}

type palette struct {
	primary1Color      color.NRGBA
	primary2Color      color.NRGBA
	primary3Color      color.NRGBA
	accent1Color       color.NRGBA
	accent2Color       color.NRGBA
	accent3Color       color.NRGBA
	textColor          color.NRGBA
	secondaryTextColor color.NRGBA
	alternateTextColor color.NRGBA
	canvasColor        color.NRGBA
	borderColor        color.NRGBA
	disabledColor      color.NRGBA
	pickerHeaderColor  color.NRGBA
	clockCircleColor   color.NRGBA
	shadowColor        color.NRGBA
}

type paletteDescription struct {
	Primary1Color      string
	Primary2Color      string
	Primary3Color      string
	Accent1Color       string
	Accent2Color       string
	Accent3Color       string
	TextColor          string
	SecondaryTextColor string
	AlternateTextColor string
	CanvasColor        string
	BorderColor        string
	DisabledColor      string
	PickerHeaderColor  string
	ClockCircleColor   string
	ShadowColor        string
}

func newPalette(base *palette, str string) (*palette, error) {
	type exportedPalette struct {
		Palette paletteDescription
	}
	ep := &exportedPalette{}
	np := &palette{}
	if base != nil {
		*np = *base
	}
	if err := json.Unmarshal([]byte(str), ep); err != nil {
		return nil, err
	}
	p := ep.Palette
	np.primary1Color = rgb(np.primary1Color, p.Primary1Color)
	np.primary2Color = rgb(np.primary2Color, p.Primary2Color)
	np.primary3Color = rgb(np.primary3Color, p.Primary3Color)
	np.accent1Color = rgb(np.accent1Color, p.Accent1Color)
	np.accent2Color = rgb(np.accent2Color, p.Accent2Color)
	np.accent3Color = rgb(np.accent3Color, p.Accent3Color)
	np.textColor = rgb(np.textColor, p.TextColor)
	np.secondaryTextColor = rgb(np.secondaryTextColor, p.SecondaryTextColor)
	np.alternateTextColor = rgb(np.alternateTextColor, p.AlternateTextColor)
	np.canvasColor = rgb(np.canvasColor, p.CanvasColor)
	np.borderColor = rgb(np.borderColor, p.BorderColor)
	np.disabledColor = rgb(np.disabledColor, p.DisabledColor)
	np.pickerHeaderColor = rgb(np.pickerHeaderColor, p.PickerHeaderColor)
	np.clockCircleColor = rgb(np.clockCircleColor, p.ClockCircleColor)
	np.shadowColor = rgb(np.shadowColor, p.ShadowColor)
	return np, nil
}

func rgb(c color.NRGBA, s string) color.NRGBA {
	if strings.HasPrefix(s, "#") && len(s) == 7 {
		fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
		c.A = 255
	}
	if strings.HasPrefix(s, "rgba(") {
		s = strings.Replace(s, " ", "", -1)
		alpha := 1.0
		fmt.Sscanf(s, "rgba(%d,%d,%d,%f)", &c.R, &c.G, &c.B, &alpha)
		c.A = uint8(255.0 * alpha)
	}
	return c
}

func limit(i int) uint8 {
	if i < 0 {
		return 0
	}
	if i > 255 {
		return 255
	}
	return uint8(i)
}

func brighten(c color.NRGBA, amount int) color.Color {
	factor := 1.0 + (float64(amount) / 10.0)
	r := limit(int(float64(c.R) * factor))
	g := limit(int(float64(c.G) * factor))
	b := limit(int(float64(c.B) * factor))
	return color.NRGBA{r, g, b, c.A}
}
