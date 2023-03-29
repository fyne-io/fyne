package harfbuzz

import fontP "github.com/go-text/typesetting/opentype/api/font"

func simpleKern(kernTable fontP.Kernx) fontP.SimpleKerns {
	for _, subtable := range kernTable {
		if simple, ok := subtable.Data.(fontP.SimpleKerns); ok {
			return simple
		}
	}
	return nil
}

func kern(driver fontP.SimpleKerns, crossStream bool, font *Font, buffer *Buffer, kernMask GlyphMask, scale bool) {
	c := newOtApplyContext(1, font, buffer)
	c.setLookupMask(kernMask)
	c.setLookupProps(uint32(otIgnoreMarks))
	skippyIter := &c.iterInput
	horizontal := buffer.Props.Direction.isHorizontal()
	info := buffer.Info
	pos := buffer.Pos
	for idx := 0; idx < len(pos); {
		if info[idx].Mask&kernMask == 0 {
			idx++
			continue
		}

		skippyIter.reset(idx, 1)
		if !skippyIter.next() {
			idx++
			continue
		}

		i := idx
		j := skippyIter.idx

		rawKern := driver.KernPair(info[i].Glyph, info[j].Glyph)
		kern := Position(rawKern)

		if rawKern == 0 {
			goto skip
		}

		if horizontal {
			if scale {
				kern = font.emScaleX(rawKern)
			}
			if crossStream {
				pos[j].YOffset = kern
				buffer.scratchFlags |= bsfHasGPOSAttachment
			} else {
				kern1 := kern >> 1
				kern2 := kern - kern1
				pos[i].XAdvance += kern1
				pos[j].XAdvance += kern2
				pos[j].XOffset += kern2
			}
		} else {
			if scale {
				kern = font.emScaleY(rawKern)
			}
			if crossStream {
				pos[j].XOffset = kern
				buffer.scratchFlags |= bsfHasGPOSAttachment
			} else {
				kern1 := kern >> 1
				kern2 := kern - kern1
				pos[i].YAdvance += kern1
				pos[j].YAdvance += kern2
				pos[j].YOffset += kern2
			}
		}

		buffer.unsafeToBreak(i, j+1)

	skip:
		idx = skippyIter.idx
	}
}

func (sp *otShapePlan) otApplyFallbackKern(font *Font, buffer *Buffer) {
	reverse := buffer.Props.Direction.isBackward()

	if reverse {
		buffer.Reverse()
	}

	if driver := simpleKern(font.face.Kern); driver != nil {
		kern(driver, false, font, buffer, sp.kernMask, false)
	}

	if reverse {
		buffer.Reverse()
	}
}
