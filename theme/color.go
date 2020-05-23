package theme

import (
	"image/color"
)

// Brighten takes a color and brightens / dims it by the given
// amount.
// A value of 1.0 == the same
// 1.1 = 10% brighter
// 0.9 = 10% dimmer, etc
func Brighten(c color.Color, amount float64) color.Color {
	r, g, b, a := c.RGBA()
	if IsDark() {
		amount = amount - 2.0*(amount - 1.0)
	}
	fn := func(i uint32) uint8 {
		x := uint32(float64(i) * amount)
		if amount > 1.0 && x > 0xffff {
			x = 0xffff
		}
		if amount < 1.0 && x > i {
			x = 0
		}
		return uint8(x >> 8)
	}
	return color.RGBA{fn(r), fn(g), fn(b), uint8(a)}
}

// PrimaryTextColor returns either {black, white} depending on
// the brightness of the primary color
func PrimaryTextColor() color.Color {
	r, g, b, _ := PrimaryColor().RGBA()
	if (r+g+b)/3 < 24000 {
		// is a dark primary
		return color.White
	}
	return color.Black
}

func IsDark() bool {
	c := BackgroundColor()
	r,g,b,_ := c.RGBA()
	if r+g+b > (0x8888*3) {
		return false
	}
	return true
}

func limit(i uint8) uint8 {
	if i > 255 {
		return 255
	}
	return uint8(i)
}
