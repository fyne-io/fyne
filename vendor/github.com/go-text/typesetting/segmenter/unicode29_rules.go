// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package segmenter

import ucd "github.com/go-text/typesetting/unicodedata"

// Apply the Grapheme_Cluster_Boundary_Rules and returns a true if we are
// at a grapheme break.
// See https://unicode.org/reports/tr29/#Grapheme_Cluster_Boundary_Rules
func (cr *cursor) applyGraphemeBoundaryRules() bool {
	triggerGB11 := cr.updatePictoSequence()    // apply rule GB11
	triggerGB12_13 := cr.updateGraphemeRIOdd() // apply rule GB12 and GB13

	br0, br1 := cr.prevGrapheme, cr.grapheme
	if cr.r == '\n' && cr.prev == '\r' {
		return false // Rule GB3
	} else if br0 == ucd.GraphemeBreakControl || br0 == ucd.GraphemeBreakCR || br0 == ucd.GraphemeBreakLF ||
		br1 == ucd.GraphemeBreakControl || br1 == ucd.GraphemeBreakCR || br1 == ucd.GraphemeBreakLF {
		return true // Rules GB4 && GB5
	} else if br0 == ucd.GraphemeBreakL &&
		(br1 == ucd.GraphemeBreakL || br1 == ucd.GraphemeBreakV || br1 == ucd.GraphemeBreakLV || br1 == ucd.GraphemeBreakLVT) { // rule GB6
		return false
	} else if (br0 == ucd.GraphemeBreakLV || br0 == ucd.GraphemeBreakV) && (br1 == ucd.GraphemeBreakV || br1 == ucd.GraphemeBreakT) {
		return false // rule GB7
	} else if (br0 == ucd.GraphemeBreakLVT || br0 == ucd.GraphemeBreakT) && br1 == ucd.GraphemeBreakT {
		return false // rule GB8
	} else if br1 == ucd.GraphemeBreakExtend || br1 == ucd.GraphemeBreakZWJ {
		return false // Rule GB9
	} else if br1 == ucd.GraphemeBreakSpacingMark {
		return false // Rule GB9a
	} else if br0 == ucd.GraphemeBreakPrepend {
		return false // Rule GB9b
	} else if triggerGB11 { // Rule GB11
		return false
	} else if triggerGB12_13 {
		return false // Rule GB12 && GB13
	}

	return true // Rule GB999
}

// update `isPrevGrRIOdd` used for the rules GB12 and GB13
// and returns `true` if one of them triggered
func (cr *cursor) updateGraphemeRIOdd() (trigger bool) {
	if cr.grapheme == ucd.GraphemeBreakRegional_Indicator {
		trigger = cr.isPrevGraphemeRIOdd
		cr.isPrevGraphemeRIOdd = !cr.isPrevGraphemeRIOdd // switch the parity
	} else {
		cr.isPrevGraphemeRIOdd = false
	}
	return trigger
}

// see rule GB11
type pictoSequenceState uint8

const (
	noPictoSequence pictoSequenceState = iota // we are not in a sequence
	inPictoExtend                             // we are in (ExtendedPic)(Extend*) pattern
	seenPictoZWJ                              // we have seen (ExtendedPic)(Extend*)(ZWJ)
)

// update the `pictoSequence` state used for rule GB11 pattern :
// (ExtendedPic)(Extend*)(ZWJ)(ExtendedPic)
// and returns true if we matched one
func (cr *cursor) updatePictoSequence() bool {
	switch cr.pictoSequence {
	case noPictoSequence:
		// we are not in a sequence yet, start it if we have an ExtendedPic
		if cr.isExtentedPic {
			cr.pictoSequence = inPictoExtend
		}
		return false
	case inPictoExtend:
		if cr.grapheme == ucd.GraphemeBreakExtend {
			// continue the sequence with an Extend rune
		} else if cr.grapheme == ucd.GraphemeBreakZWJ {
			// close the variable part of the sequence with (ZWJ)
			cr.pictoSequence = seenPictoZWJ
		} else {
			// stop the sequence
			cr.pictoSequence = noPictoSequence
		}
		return false
	case seenPictoZWJ:
		// trigger GB11 if we have an ExtendedPic,
		// and reset the sequence
		if cr.isExtentedPic {
			cr.pictoSequence = inPictoExtend
			return true
		}
		cr.pictoSequence = noPictoSequence
		return false
	default:
		panic("exhaustive switch")
	}
}
