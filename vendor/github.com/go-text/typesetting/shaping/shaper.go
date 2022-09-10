// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package shaping

import (
	"fmt"

	"github.com/benoitkugler/textlayout/harfbuzz"
	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"golang.org/x/image/math/fixed"
)

type Shaper interface {
	// Shape takes an Input and shapes it into the Output.
	Shape(Input) Output
}

// MissingGlyphError indicates that the font used in shaping did not
// have a glyph needed to complete the shaping.
type MissingGlyphError struct {
	font.GID
}

func (m MissingGlyphError) Error() string {
	return fmt.Sprintf("missing glyph with id %d", m.GID)
}

// InvalidRunError represents an invalid run of text, either because
// the end is before the start or because start or end is greater
// than the length.
type InvalidRunError struct {
	RunStart, RunEnd, TextLength int
}

func (i InvalidRunError) Error() string {
	return fmt.Sprintf("run from %d to %d is not valid for text len %d", i.RunStart, i.RunEnd, i.TextLength)
}

const (
	// scaleShift is the power of 2 with which to automatically scale
	// up the input coordinate space of the shaper. This factor will
	// be removed prior to returning dimensions. This ensures that the
	// returned glyph dimensions take advantage of all of the precision
	// that a fixed.Int26_6 can provide.
	scaleShift = 6
)

// Shape turns an input into an output.
func Shape(input Input) (Output, error) {
	// Prepare to shape the text.
	// TODO: maybe reuse these buffers for performance?
	buf := harfbuzz.NewBuffer()
	runes, start, end := input.Text, input.RunStart, input.RunEnd
	if end < start {
		return Output{}, InvalidRunError{RunStart: start, RunEnd: end, TextLength: len(input.Text)}
	}
	buf.AddRunes(runes, start, end-start)
	// TODO: handle vertical text?
	switch input.Direction {
	case di.DirectionLTR:
		buf.Props.Direction = harfbuzz.LeftToRight
	case di.DirectionRTL:
		buf.Props.Direction = harfbuzz.RightToLeft
	default:
		return Output{}, UnimplementedDirectionError{
			Direction: input.Direction,
		}
	}
	buf.Props.Language = input.Language
	buf.Props.Script = input.Script
	// TODO: figure out what (if anything) to do if this type assertion fails.
	font := harfbuzz.NewFont(input.Face.(harfbuzz.Face))
	font.XScale = int32(input.Size.Ceil()) << scaleShift
	font.YScale = font.XScale

	// Actually use harfbuzz to shape the text.
	buf.Shape(font, nil)

	// Convert the shaped text into an Output.
	glyphs := make([]Glyph, len(buf.Info))
	for i := range glyphs {
		g := buf.Info[i].Glyph
		extents, ok := font.GlyphExtents(g)
		if !ok {
			// TODO: can this error happen? Will harfbuzz return a
			// GID for a glyph that isn't in the font?
			return Output{}, MissingGlyphError{GID: g}
		}
		glyphs[i] = Glyph{
			Width:        fixed.I(int(extents.Width)) >> scaleShift,
			Height:       fixed.I(int(extents.Height)) >> scaleShift,
			XBearing:     fixed.I(int(extents.XBearing)) >> scaleShift,
			YBearing:     fixed.I(int(extents.YBearing)) >> scaleShift,
			XAdvance:     fixed.I(int(buf.Pos[i].XAdvance)) >> scaleShift,
			YAdvance:     fixed.I(int(buf.Pos[i].YAdvance)) >> scaleShift,
			XOffset:      fixed.I(int(buf.Pos[i].XOffset)) >> scaleShift,
			YOffset:      fixed.I(int(buf.Pos[i].YOffset)) >> scaleShift,
			ClusterIndex: buf.Info[i].Cluster,
			GlyphID:      g,
			Mask:         buf.Info[i].Mask,
		}
	}
	countClusters(glyphs, input.RunEnd-input.RunStart, input.Direction)
	out := Output{
		Glyphs:    glyphs,
		Direction: input.Direction,
	}
	fontExtents := font.ExtentsForDirection(buf.Props.Direction)
	out.LineBounds = Bounds{
		Ascent:  fixed.I(int(fontExtents.Ascender)) >> scaleShift,
		Descent: fixed.I(int(fontExtents.Descender)) >> scaleShift,
		Gap:     fixed.I(int(fontExtents.LineGap)) >> scaleShift,
	}
	return out, out.RecalculateAll()
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
