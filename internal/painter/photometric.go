package painter

import (
	"image/color"
	"math"
)

// srgbToLinearLUT maps 8-bit non-linear sRGB channel values (0-255) to 16-bit linear photometric intensities (0-65535).
var srgbToLinearLUT [256]uint16

// linearToSrgbLUT maps 16-bit linear photometric intensities (0-65535) to 8-bit non-linear sRGB channel values (0-255).
var linearToSrgbLUT [65536]uint8

func init() {
	// Standard IEC 61966-2-1 transfer function
	for i := 0; i < 256; i++ {
		c := float64(i) / 255.0
		var lin float64
		if c <= 0.04045 {
			lin = c / 12.92
		} else {
			lin = math.Pow((c+0.055)/1.055, 2.4)
		}
		srgbToLinearLUT[i] = uint16(math.Round(lin * 65535.0))
	}

	for i := 0; i < 65536; i++ {
		lin := float64(i) / 65535.0
		var srgb float64
		if lin <= 0.0031308 {
			srgb = lin * 12.92
		} else {
			srgb = 1.055*math.Pow(lin, 1.0/2.4) - 0.055
		}
		if srgb < 0.0 {
			srgb = 0.0
		} else if srgb > 1.0 {
			srgb = 1.0
		}
		linearToSrgbLUT[i] = uint8(math.Round(srgb * 255.0))
	}
}

// BlendPhotometric blends a source color over a destination color in linear photometric sRGB space.
// 50% white over black yields sRGB 188 rather than 128. This primitive is not yet wired into draw.go.
func BlendPhotometric(dst, src color.NRGBA) color.NRGBA {
	if src.A == 0 {
		return dst
	}
	if src.A == 255 || dst.A == 0 {
		return src
	}

	sa := uint32(src.A)
	da := uint32(dst.A)
	invA := 255 - sa

	// Out alpha calculation
	outA := sa + (da*invA+127)/255
	if outA == 0 {
		return color.NRGBA{}
	}

	// Linearize RGB components
	srLin := uint32(srgbToLinearLUT[src.R])
	sgLin := uint32(srgbToLinearLUT[src.G])
	sbLin := uint32(srgbToLinearLUT[src.B])

	drLin := uint32(srgbToLinearLUT[dst.R])
	dgLin := uint32(srgbToLinearLUT[dst.G])
	dbLin := uint32(srgbToLinearLUT[dst.B])

	// Blend in physical linear space
	rLin := (srLin*sa + (drLin*da*invA+127)/255 + outA/2) / outA
	gLin := (sgLin*sa + (dgLin*da*invA+127)/255 + outA/2) / outA
	bLin := (sbLin*sa + (dbLin*da*invA+127)/255 + outA/2) / outA

	if rLin > 65535 {
		rLin = 65535
	}
	if gLin > 65535 {
		gLin = 65535
	}
	if bLin > 65535 {
		bLin = 65535
	}

	return color.NRGBA{
		R: linearToSrgbLUT[rLin],
		G: linearToSrgbLUT[gLin],
		B: linearToSrgbLUT[bLin],
		A: uint8(outA),
	}
}
