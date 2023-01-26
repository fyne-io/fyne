package harfbuzz

// ported from harfbuzz/src/hb-fallback-shape.cc Copyright © 2011  Google, Inc. Behdad Esfahbod

var _ shaper = shaperFallback{}

// shaperFallback implements a naive shaper, which does the minimum,
// without requiring advanced Opentype font features.
type shaperFallback struct{}

func (shaperFallback) kind() shaperKind { return skFallback }

func (shaperFallback) compile(props SegmentProperties, userFeatures []Feature) {
}

func (shaperFallback) shape(font *Font, buffer *Buffer, _ []Feature) {
	space, hasSpace := font.face.NominalGlyph(' ')

	buffer.clearPositions()

	direction := buffer.Props.Direction
	info := buffer.Info
	pos := buffer.Pos
	for i := range info {
		if hasSpace && uni.isDefaultIgnorable(info[i].codepoint) {
			info[i].Glyph = space
			pos[i].XAdvance = 0
			pos[i].YAdvance = 0
		} else {
			info[i].Glyph, _ = font.face.NominalGlyph(info[i].codepoint)
			pos[i].XAdvance, pos[i].YAdvance = font.GlyphAdvanceForDirection(info[i].Glyph, direction)
			pos[i].XOffset, pos[i].YOffset = font.subtractGlyphOriginForDirection(info[i].Glyph, direction,
				pos[i].XOffset, pos[i].YOffset)
		}
	}

	if direction.isBackward() {
		buffer.Reverse()
	}

	buffer.clearGlyphFlags(0)
}
