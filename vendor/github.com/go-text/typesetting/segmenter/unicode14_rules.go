// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package segmenter

import (
	"unicode"

	ucd "github.com/go-text/typesetting/unicodedata"
)

// Apply the Line Breaking Rules and returns the computed break opportunity
// See https://unicode.org/reports/tr14/#BreakingRules
func (cr *cursor) applyLineBreakingRules() breakOpportunity {
	// start by attributing the break class for the current rune
	cr.ruleLB1()

	triggerNumSequence := cr.updateNumSequence()

	// add the line break rules in reverse order to override
	// the lower priority rules.
	breakOp := breakEmpty

	cr.ruleLB30(&breakOp)
	cr.ruleLB30ab(&breakOp)
	cr.ruleLB29To26(&breakOp)
	cr.ruleLB25(&breakOp, triggerNumSequence)
	cr.ruleLB24To22(&breakOp)
	cr.ruleLB21To9(&breakOp)
	cr.ruleLB8(&breakOp)
	cr.ruleLB7To4(&breakOp)

	return breakOp
}

// breakOpportunity is a convenient enum,
// mapped to the LineBreak and MandatoryBreak properties,
// avoiding too many bit operations
type breakOpportunity uint8

const (
	breakEmpty      breakOpportunity = iota // not specified
	breakProhibited                         // no break
	breakAllowed                            // direct break (can always break here)
	breakMandatory                          // break is mandatory (implies breakAllowed)
)

func (cr *cursor) ruleLB30(breakOp *breakOpportunity) {
	// (AL | HL | NU) × [OP-[\p{ea=F}\p{ea=W}\p{ea=H}]]
	if (cr.prevLine == ucd.BreakAL || cr.prevLine == ucd.BreakHL || cr.prevLine == ucd.BreakNU) &&
		cr.line == ucd.BreakOP && !unicode.Is(ucd.LargeEastAsian, cr.r) {
		*breakOp = breakProhibited
	}
	// [CP-[\p{ea=F}\p{ea=W}\p{ea=H}]] × (AL | HL | NU)
	if cr.prevLine == ucd.BreakCP && !unicode.Is(ucd.LargeEastAsian, cr.prev) &&
		(cr.line == ucd.BreakAL || cr.line == ucd.BreakHL || cr.line == ucd.BreakNU) {
		*breakOp = breakProhibited
	}
}

func (cr *cursor) ruleLB30ab(breakOp *breakOpportunity) {
	// (RI RI)* RI × RI
	if cr.isPrevLinebreakRIOdd && cr.line == ucd.BreakRI { // LB30a
		*breakOp = breakProhibited
	}

	// LB30b
	// EB × EM
	if cr.prevLine == ucd.BreakEB && cr.line == ucd.BreakEM {
		*breakOp = breakProhibited
	}
	// [\p{Extended_Pictographic}&\p{Cn}] × EM
	if unicode.Is(ucd.Extended_Pictographic, cr.prev) && ucd.LookupType(cr.prev) == nil &&
		cr.line == ucd.BreakEM {
		*breakOp = breakProhibited
	}
}

func (cr *cursor) ruleLB29To26(breakOp *breakOpportunity) {
	b0, b1 := cr.prevLine, cr.line
	// LB29 : IS × (AL | HL)
	if b0 == ucd.BreakIS && (b1 == ucd.BreakAL || b1 == ucd.BreakHL) {
		*breakOp = breakProhibited
	}
	// LB28 : (AL | HL) × (AL | HL)
	if (b0 == ucd.BreakAL || b0 == ucd.BreakHL) && (b1 == ucd.BreakAL || b1 == ucd.BreakHL) {
		*breakOp = breakProhibited
	}
	// LB27
	// (JL | JV | JT | H2 | H3) × PO
	if (b0 == ucd.BreakJL || b0 == ucd.BreakJV || b0 == ucd.BreakJT || b0 == ucd.BreakH2 || b0 == ucd.BreakH3) &&
		b1 == ucd.BreakPO {
		*breakOp = breakProhibited
	}
	// PR × (JL | JV | JT | H2 | H3)
	if b0 == ucd.BreakPR &&
		(b1 == ucd.BreakJL || b1 == ucd.BreakJV || b1 == ucd.BreakJT || b1 == ucd.BreakH2 || b1 == ucd.BreakH3) {
		*breakOp = breakProhibited
	}
	// LB26
	// JL × (JL | JV | H2 | H3)
	if b0 == ucd.BreakJL &&
		(b1 == ucd.BreakJL || b1 == ucd.BreakJV || b1 == ucd.BreakH2 || b1 == ucd.BreakH3) {
		*breakOp = breakProhibited
	}
	// (JV | H2) × (JV | JT)
	if (b0 == ucd.BreakJV || b0 == ucd.BreakH2) && (b1 == ucd.BreakJV || b1 == ucd.BreakJT) {
		*breakOp = breakProhibited
	}
	// (JT | H3) × JT
	if (b0 == ucd.BreakJT || b0 == ucd.BreakH3) && b1 == ucd.BreakJT {
		*breakOp = breakProhibited
	}
}

// we follow other implementations by using the tailoring described
// in Example 7
func (cr *cursor) ruleLB25(breakOp *breakOpportunity, triggerNumSequence bool) {
	br0, br1 := cr.prevLine, cr.line
	// (PR | PO) × ( OP | HY )? NU
	if (br0 == ucd.BreakPR || br0 == ucd.BreakPO) && br1 == ucd.BreakNU {
		*breakOp = breakProhibited
	}
	if (br0 == ucd.BreakPR || br0 == ucd.BreakPO) &&
		(br1 == ucd.BreakOP || br1 == ucd.BreakHY) &&
		cr.nextLine == ucd.BreakNU {
		*breakOp = breakProhibited
	}
	// ( OP | HY ) × NU
	if (br0 == ucd.BreakOP || br0 == ucd.BreakHY) && br1 == ucd.BreakNU {
		*breakOp = breakProhibited
	}
	// NU × (NU | SY | IS)
	if br0 == ucd.BreakNU && (br1 == ucd.BreakNU || br1 == ucd.BreakSY || br1 == ucd.BreakIS) {
		*breakOp = breakProhibited
	}
	// NU (NU | SY | IS)* × (NU | SY | IS | CL | CP )
	if triggerNumSequence {
		*breakOp = breakProhibited
	}
}

func (cr *cursor) ruleLB24To22(breakOp *breakOpportunity) {
	br0, br1 := cr.prevLine, cr.line
	// LB24
	// (PR | PO) × (AL | HL)
	if (br0 == ucd.BreakPR || br0 == ucd.BreakPO) && (br1 == ucd.BreakAL || br1 == ucd.BreakHL) {
		*breakOp = breakProhibited
	}
	// (AL | HL) × (PR | PO)
	if (br0 == ucd.BreakAL || br0 == ucd.BreakHL) && (br1 == ucd.BreakPR || br1 == ucd.BreakPO) {
		*breakOp = breakProhibited
	}
	// LB23
	// (AL | HL) × NU
	if (br0 == ucd.BreakAL || br0 == ucd.BreakHL) && br1 == ucd.BreakNU {
		*breakOp = breakProhibited
	}
	// NU × (AL | HL)
	if br0 == ucd.BreakNU && (br1 == ucd.BreakAL || br1 == ucd.BreakHL) {
		*breakOp = breakProhibited
	}
	// LB23a
	// PR × (ID | EB | EM)
	if br0 == ucd.BreakPR && (br1 == ucd.BreakID || br1 == ucd.BreakEB || br1 == ucd.BreakEM) {
		*breakOp = breakProhibited
	}
	// (ID | EB | EM) × PO
	if (br0 == ucd.BreakID || br0 == ucd.BreakEB || br0 == ucd.BreakEM) && br1 == ucd.BreakPO {
		*breakOp = breakProhibited
	}

	// LB22 : × IN
	if br1 == ucd.BreakIN {
		*breakOp = breakProhibited
	}
}

func (cr *cursor) ruleLB21To9(breakOp *breakOpportunity) {
	br0, br1 := cr.prevLine, cr.line
	// LB21
	// × BA
	// × HY
	// × NS
	// BB ×
	if br1 == ucd.BreakBA || br1 == ucd.BreakHY || br1 == ucd.BreakNS || br0 == ucd.BreakBB {
		*breakOp = breakProhibited
	}
	// LB21a : HL (HY | BA) ×
	if cr.prevPrevLine == ucd.BreakHL &&
		(br0 == ucd.BreakHY || br0 == ucd.BreakBA) {
		*breakOp = breakProhibited
	}
	// LB21b : SY × HL
	if br0 == ucd.BreakSY && br1 == ucd.BreakHL {
		*breakOp = breakProhibited
	}
	// LB20
	// ÷ CB
	// CB ÷
	if br0 == ucd.BreakCB || br1 == ucd.BreakCB {
		*breakOp = breakAllowed
	}
	// LB19
	// × QU
	// QU ×
	if br0 == ucd.BreakQU || br1 == ucd.BreakQU {
		*breakOp = breakProhibited
	}
	// LB18 : SP ÷
	if br0 == ucd.BreakSP {
		*breakOp = breakAllowed
	}
	// LB17 : B2 SP* × B2
	if cr.beforeSpaces == ucd.BreakB2 && br1 == ucd.BreakB2 {
		*breakOp = breakProhibited
	}
	// LB16 : (CL | CP) SP* × NS
	if (cr.beforeSpaces == ucd.BreakCL || cr.beforeSpaces == ucd.BreakCP) && br1 == ucd.BreakNS {
		*breakOp = breakProhibited
	}
	// LB15 : QU SP* × OP
	if cr.beforeSpaces == ucd.BreakQU && br1 == ucd.BreakOP {
		*breakOp = breakProhibited
	}
	// LB14 : OP SP* ×
	if cr.beforeSpaces == ucd.BreakOP {
		*breakOp = breakProhibited
	}

	// rule LB13, with the tailoring described in Example 7
	// × EX
	if br1 == ucd.BreakEX {
		*breakOp = breakProhibited
	}
	// [^NU] × CL
	// [^NU] × CP
	// [^NU] × IS
	// [^NU] × SY
	if br0 != ucd.BreakNU &&
		(br1 == ucd.BreakCL || br1 == ucd.BreakCP || br1 == ucd.BreakIS || br1 == ucd.BreakSY) {
		*breakOp = breakProhibited
	}
	// LB12 : GL ×
	if br0 == ucd.BreakGL {
		*breakOp = breakProhibited
	}
	// LB12a : [^SP BA HY] × GL
	if (br0 != ucd.BreakSP && br0 != ucd.BreakBA && br0 != ucd.BreakHY) &&
		br1 == ucd.BreakGL {
		*breakOp = breakProhibited
	}
	// LB11
	// × WJ
	// WJ ×
	if br0 == ucd.BreakWJ || br1 == ucd.BreakWJ {
		*breakOp = breakProhibited
	}

	// rule LB9 : "Do not break a combining character sequence"
	// where X is any line break class except BK, CR, LF, NL, SP, or ZW.
	// see also [endIteration]
	if br1 == ucd.BreakCM || br1 == ucd.BreakZWJ {
		if !(br0 == ucd.BreakBK || br0 == ucd.BreakCR || br0 == ucd.BreakLF ||
			br0 == ucd.BreakNL || br0 == ucd.BreakSP || br0 == ucd.BreakZW) {
			*breakOp = breakProhibited
		}
	}
}

func (cr *cursor) ruleLB8(breakOp *breakOpportunity) {
	// rule LB8 : ZW SP* ÷
	if cr.beforeSpaces == ucd.BreakZW {
		*breakOp = breakAllowed
	}
	// rule LB8a : ZWJ ×
	// there is a catch here : prevLine is not always exactly
	// the class at index i-1, because of rules LB9 and LB10
	// however, rule LB8a applies before LB9 and LB10, meaning
	// we need to use the real class
	if unicode.Is(ucd.BreakZWJ, cr.prev) {
		*breakOp = breakProhibited
	}
}

func (cr *cursor) ruleLB7To4(breakOp *breakOpportunity) {
	// LB7
	// × SP
	// × ZW
	if cr.line == ucd.BreakSP || cr.line == ucd.BreakZW {
		*breakOp = breakProhibited
	}
	// LB6 : × ( BK | CR | LF | NL )
	if cr.line == ucd.BreakBK || cr.line == ucd.BreakCR || cr.line == ucd.BreakLF || cr.line == ucd.BreakNL {
		*breakOp = breakProhibited
	}

	// LB4 and LB5
	// BK !
	// CR !
	// LF !
	// NL !
	// (CR × LF is actually handled in rule LB6)
	if cr.prevLine == ucd.BreakBK || (cr.prevLine == ucd.BreakCR && cr.r != '\n') ||
		cr.prevLine == ucd.BreakLF || cr.prevLine == ucd.BreakNL {
		*breakOp = breakMandatory
	}
}

// apply rule LB1 to resolve break classses AI, SG, XX, SA and CJ.
// We use the default values specified in https://unicode.org/reports/tr14/#BreakingRules.
func (cr *cursor) ruleLB1() {
	switch cr.line {
	case ucd.BreakAI, ucd.BreakSG, ucd.BreakXX:
		cr.line = ucd.BreakAL
	case ucd.BreakSA:
		generalCategory := ucd.LookupType(cr.r)
		if generalCategory == unicode.Mn || generalCategory == unicode.Mc {
			cr.line = ucd.BreakCM
		} else {
			cr.line = ucd.BreakAL
		}
	case ucd.BreakCJ:
		cr.line = ucd.BreakNS
	}
}

type numSequenceState uint8

const (
	noNumSequence numSequenceState = iota // we are not in a sequence
	inNumSequence                         // we are in NU (NU | SY | IS)*
	seenCloseNum                          // we are at NU (NU | SY | IS)* (CL | CP)?
)

// update the `numSequence` state used for rule LB25
// and returns true if we matched one
func (cr *cursor) updateNumSequence() bool {
	// note that rule LB9 also apply : (CM|ZWJ) do not change
	// the flag
	if cr.line == ucd.BreakCM || cr.line == ucd.BreakZWJ {
		return false
	}

	switch cr.numSequence {
	case noNumSequence:
		if cr.line == ucd.BreakNU { // start a sequence
			cr.numSequence = inNumSequence
		}
		return false
	case inNumSequence:
		switch cr.line {
		case ucd.BreakNU, ucd.BreakSY, ucd.BreakIS:
			// NU (NU | SY | IS)* × (NU | SY | IS) : the sequence continue
			return true
		case ucd.BreakCL, ucd.BreakCP:
			// NU (NU | SY | IS)* × (CL | CP)
			cr.numSequence = seenCloseNum
			return true
		case ucd.BreakPO, ucd.BreakPR:
			// NU (NU | SY | IS)* × (PO | PR) : close the sequence
			cr.numSequence = noNumSequence
			return true
		default:
			cr.numSequence = noNumSequence
			return false
		}
	case seenCloseNum:
		cr.numSequence = noNumSequence // close the sequence anyway
		if cr.line == ucd.BreakPO || cr.line == ucd.BreakPR {
			// NU (NU | SY | IS)* (CL | CP) × (PO | PR)
			return true
		}
		return false
	default:
		panic("exhaustive switch")
	}
}

// startIteration updates the cursor properties, setting the current
// rune to text[i].
// Some properties depending on the context are rather
// updated in the previous `endIteration` call.
func (cr *cursor) startIteration(text []rune, i int) {
	cr.prev = cr.r
	if i < len(text) {
		cr.r = text[i]
	} else {
		cr.r = paragraphSeparator
	}
	if i == len(text) {
		cr.next = 0
	} else if i == len(text)-1 {
		// we fill in the last element of `attrs` by assuming
		// there's a paragraph separators off the end of text
		cr.next = paragraphSeparator
	} else {
		cr.next = text[i+1]
	}

	// query general unicode properties for the current rune
	cr.isExtentedPic = unicode.Is(ucd.Extended_Pictographic, cr.r)

	cr.prevGrapheme = cr.grapheme
	cr.grapheme = ucd.LookupGraphemeBreakClass(cr.r)

	// prevPrevLine and prevLine are handled in endIteration
	cr.line = cr.nextLine // avoid calling LookupBreakClass twice
	cr.nextLine = ucd.LookupLineBreakClass(cr.next)
}

// end the current iteration, computing some of the properties
// required for the next rune and respecting rule LB9 and LB10
func (cr *cursor) endIteration(isStart bool) {
	// start by handling rule LB9 and LB10
	if cr.line == ucd.BreakCM || cr.line == ucd.BreakZWJ {
		isLB10 := cr.prevLine == ucd.BreakBK ||
			cr.prevLine == ucd.BreakCR ||
			cr.prevLine == ucd.BreakLF ||
			cr.prevLine == ucd.BreakNL ||
			cr.prevLine == ucd.BreakSP ||
			cr.prevLine == ucd.BreakZW
		if isStart || isLB10 { // Rule LB10
			cr.prevLine = ucd.BreakAL
		} // else rule LB9 : ignore the rune for prevLine and prevPrevLine

	} else { // regular update
		cr.prevPrevLine = cr.prevLine
		cr.prevLine = cr.line
	}

	// keep track of the rune before the spaces
	if cr.prevLine != ucd.BreakSP {
		cr.beforeSpaces = cr.prevLine
	}

	// update RegionalIndicator parity used for LB30a
	if cr.line == ucd.BreakRI {
		cr.isPrevLinebreakRIOdd = !cr.isPrevLinebreakRIOdd
	} else if !(cr.line == ucd.BreakCM || cr.line == ucd.BreakZWJ) { // beware of the rule LB9: (CM|ZWJ) ignore the update
		cr.isPrevLinebreakRIOdd = false
	}
}
