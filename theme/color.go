package theme

import (
	"image/color"
)

// Brighten takes a color and brightens / dims it by the given
// amount. Each value in the amount represents a 10% increase
// in the intensity of the color.  Negative values allow for
// dimming the color. This is used to compute a color such as
// a hightlight on hover based on a base color.
func Brighten(c color.Color, amount float64) color.Color {
	r, g, b, a := c.RGBA()
	factor := 1.0 + (amount / 10.0)
	fn := func(i uint32, f float64) uint8 {
		return uint8(float64(i) * f / 255)
	}
	return color.RGBA{fn(r, factor), fn(g, factor), fn(b, factor), uint8(a)}
}

// PrimaryTextColor returns either {black, white} depending on
// the brightness of the primary color
// add some code here
func PrimaryTextColor() color.Color {
	r, g, b, _ := PrimaryColor().RGBA()
	if (r+g+b)/3 < 24000 {
		// is a dark primary
		return color.White
	}
	return color.Black
}

func limit(i uint8) uint8 {
	if i < 0 {
		return 0
	}
	if i > 255 {
		return 255
	}
	return uint8(i)
}
