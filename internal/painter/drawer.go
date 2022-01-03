package painter

import (
	"image"
	"image/draw"
	"log"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// FontDrawer extends "golang.org/x/image/font" to add support for tabs
// FontDrawer draws text on a destination image.
//
// A FontDrawer is not safe for concurrent use by multiple goroutines, since its
// Face is not.
type FontDrawer struct {
	font.Drawer
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

// DrawString draws s at the dot and advances the dot's location.
// Tabs are translated into a dot location change.
func (d *FontDrawer) DrawString(s string, tabWidth int) {
	walkString(d.Face, s, tabWidth, &d.Dot.X, func(r rune) (fixed.Int26_6, bool) {
		dr, mask, maskp, advance, ok := d.Face.Glyph(d.Dot, r)
		if ok {
			draw.DrawMask(d.Dst, dr, d.Src, image.Point{}, mask, maskp, draw.Over)
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
