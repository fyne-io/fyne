package theme

import (
	"fmt"
	"image/color"
	"strings"
)

// Brighten takes a color and brightens / dims it by the given
// amount. Each value in the amount represents a 10% increase
// in the intensity of the color.  Negative values allow for
// dimming the color. This is used to compute a color such as
// a hightlight on hover based on a base color.
func Brighten(clr color.Color, amount float64) color.Color {
	factor := 1.0 + (amount / 10.0)
	c, ok := clr.(color.NRGBA)
	if !ok {
		return clr
	}
	r := limit(int(float64(c.R) * factor))
	g := limit(int(float64(c.G) * factor))
	b := limit(int(float64(c.B) * factor))
	return color.NRGBA{r, g, b, c.A}
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

func brighten(c color.NRGBA, amount float64) color.Color {
	factor := 1.0 + (amount / 10.0)
	r := limit(int(float64(c.R) * factor))
	g := limit(int(float64(c.G) * factor))
	b := limit(int(float64(c.B) * factor))
	return color.NRGBA{r, g, b, c.A}
}
