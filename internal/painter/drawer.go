package painter

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"

	"github.com/goki/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// DrawString draws a string into an image.
func DrawString(dst draw.Image, s string, color color.Color, face font.Face, height int, tabWidth int) {
	src := &image.Uniform{C: color}
	dot := freetype.Pt(0, height-face.Metrics().Descent.Ceil())
	walkString(face, s, tabWidth, &dot.X, func(r rune) (fixed.Int26_6, bool) {
		dr, mask, maskp, advance, ok := face.Glyph(dot, r)
		if ok {
			draw.DrawMask(dst, dr, src, image.Point{}, mask, maskp, draw.Over)
		}
		return advance, ok
	})
}

// MeasureString returns how far dot would advance by drawing s with f.
// Tabs are translated into a dot location change.
func MeasureString(f font.Face, s string, tabWidth int) (advance fixed.Int26_6) {
	walkString(f, s, tabWidth, &advance, f.GlyphAdvance)
	return
}

func tabStop(f font.Face, x fixed.Int26_6, tabWidth int) fixed.Int26_6 {
	if tabWidth <= 0 {
		tabWidth = DefaultTabWidth
	}
	spacew, ok := f.GlyphAdvance(' ')
	if !ok {
		log.Print("Failed to find space width for tab")
		return x
	}
	tabw := spacew * fixed.Int26_6(tabWidth)
	tabs, _ := math.Modf(float64((x + tabw) / tabw))
	return tabw * fixed.Int26_6(tabs)
}

func walkString(f font.Face, s string, tabWidth int, advance *fixed.Int26_6, cb func(r rune) (fixed.Int26_6, bool)) {
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			*advance += f.Kern(prevC, c)
		}
		if c == '\t' {
			*advance = tabStop(f, *advance, tabWidth)
		} else {
			a, ok := cb(c)
			if !ok {
				// TODO: is falling back on the U+FFFD glyph the responsibility of
				// the Drawer or the Face?
				// TODO: set prevC = '\ufffd'?
				continue
			}
			*advance += a
		}

		prevC = c
	}
}
