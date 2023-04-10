// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package shaping

import (
	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/harfbuzz"
	"golang.org/x/image/math/fixed"
)

// HarfbuzzShaper implements the Shaper interface using harfbuzz.
// Reusing this shaper type across multiple shaping operations is
// faster and more memory-efficient than creating a new shaper
// for each operation.
type HarfbuzzShaper struct {
	buf *harfbuzz.Buffer

	fonts fontLRU
}

// SetFontCacheSize adjusts the size of the font cache within the shaper.
// It is safe to adjust the size after using the shaper, though shrinking
// it may result in many evictions on the next shaping.
func (h *HarfbuzzShaper) SetFontCacheSize(size int) {
	h.fonts.maxSize = size
}

var _ Shaper = (*HarfbuzzShaper)(nil)

// Shaper describes the signature of a font shaping operation.
type Shaper interface {
	// Shape takes an Input and shapes it into the Output.
	Shape(Input) Output
}

const (
	// scaleShift is the power of 2 with which to automatically scale
	// up the input coordinate space of the shaper. This factor will
	// be removed prior to returning dimensions. This ensures that the
	// returned glyph dimensions take advantage of all of the precision
	// that a fixed.Int26_6 can provide.
	scaleShift = 6
)

// clamp ensures val is in the inclusive range [low,high].
func clamp(val, low, high int) int {
	if val < low {
		return low
	}
	if val > high {
		return high
	}
	return val
}

// Shape turns an input into an output.
func (t *HarfbuzzShaper) Shape(input Input) Output {
	// Prepare to shape the text.
	if t.buf == nil {
		t.buf = harfbuzz.NewBuffer()
	} else {
		t.buf.Clear()
	}

	runes, start, end := input.Text, input.RunStart, input.RunEnd
	if end < start {
		// Try to guess what the caller actually wanted.
		end, start = start, end
	}
	start = clamp(start, 0, len(runes))
	end = clamp(end, 0, len(runes))
	t.buf.AddRunes(runes, start, end-start)
	switch input.Direction {
	case di.DirectionRTL:
		t.buf.Props.Direction = harfbuzz.RightToLeft
	case di.DirectionBTT:
		t.buf.Props.Direction = harfbuzz.BottomToTop
	case di.DirectionTTB:
		t.buf.Props.Direction = harfbuzz.TopToBottom
	default:
		// Default to LTR.
		t.buf.Props.Direction = harfbuzz.LeftToRight
	}
	t.buf.Props.Language = input.Language
	t.buf.Props.Script = input.Script

	// reuse font when possible
	font, ok := t.fonts.Get(input.Face.Font)
	if !ok { // create a new font and cache it
		font = harfbuzz.NewFont(input.Face)
		t.fonts.Put(input.Face.Font, font)
	}
	// adjust the user provided fields
	font.XScale = int32(input.Size.Ceil()) << scaleShift
	font.YScale = font.XScale

	// Actually use harfbuzz to shape the text.
	t.buf.Shape(font, nil)

	// Convert the shaped text into an Output.
	glyphs := make([]Glyph, len(t.buf.Info))
	for i := range glyphs {
		g := t.buf.Info[i].Glyph
		glyphs[i] = Glyph{
			ClusterIndex: t.buf.Info[i].Cluster,
			GlyphID:      g,
			Mask:         t.buf.Info[i].Mask,
		}
		extents, ok := font.GlyphExtents(g)
		if !ok {
			// Leave the glyph having zero size if it isn't in the font. There
			// isn't really anything we can do to recover from such an error.
			continue
		}
		glyphs[i].Width = fixed.I(int(extents.Width)) >> scaleShift
		glyphs[i].Height = fixed.I(int(extents.Height)) >> scaleShift
		glyphs[i].XBearing = fixed.I(int(extents.XBearing)) >> scaleShift
		glyphs[i].YBearing = fixed.I(int(extents.YBearing)) >> scaleShift
		glyphs[i].XAdvance = fixed.I(int(t.buf.Pos[i].XAdvance)) >> scaleShift
		glyphs[i].YAdvance = fixed.I(int(t.buf.Pos[i].YAdvance)) >> scaleShift
		glyphs[i].XOffset = fixed.I(int(t.buf.Pos[i].XOffset)) >> scaleShift
		glyphs[i].YOffset = fixed.I(int(t.buf.Pos[i].YOffset)) >> scaleShift
	}
	countClusters(glyphs, input.RunEnd, input.Direction)
	out := Output{
		Glyphs:    glyphs,
		Direction: input.Direction,
		Face:      input.Face,
		Size:      input.Size,
	}
	fontExtents := font.ExtentsForDirection(t.buf.Props.Direction)
	out.LineBounds = Bounds{
		Ascent:  fixed.I(int(fontExtents.Ascender)) >> scaleShift,
		Descent: fixed.I(int(fontExtents.Descender)) >> scaleShift,
		Gap:     fixed.I(int(fontExtents.LineGap)) >> scaleShift,
	}
	out.Runes.Offset = input.RunStart
	out.Runes.Count = input.RunEnd - input.RunStart
	out.RecalculateAll()
	return out
}

// countClusters tallies the number of runes and glyphs in each cluster
// and updates the relevant fields on the provided glyph slice.
func countClusters(glyphs []Glyph, textLen int, dir di.Direction) {
	currentCluster := -1
	runesInCluster := 0
	glyphsInCluster := 0
	previousCluster := textLen
	for i := range glyphs {
		g := glyphs[i].ClusterIndex
		if g != currentCluster {
			// If we're processing a new cluster, count the runes and glyphs
			// that compose it.
			runesInCluster = 0
			glyphsInCluster = 1
			currentCluster = g
			nextCluster := -1
		glyphCountLoop:
			for k := i + 1; k < len(glyphs); k++ {
				if glyphs[k].ClusterIndex == g {
					glyphsInCluster++
				} else {
					nextCluster = glyphs[k].ClusterIndex
					break glyphCountLoop
				}
			}
			if nextCluster == -1 {
				nextCluster = textLen
			}
			switch dir {
			case di.DirectionLTR:
				runesInCluster = nextCluster - currentCluster
			case di.DirectionRTL:
				runesInCluster = previousCluster - currentCluster
			}
			previousCluster = g
		}
		glyphs[i].GlyphCount = glyphsInCluster
		glyphs[i].RuneCount = runesInCluster
	}
}
