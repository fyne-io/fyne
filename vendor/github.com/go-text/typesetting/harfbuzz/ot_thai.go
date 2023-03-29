package harfbuzz

import (
	"github.com/go-text/typesetting/language"
)

// ported from harfbuzz/src/hb-ot-shape-complex-thai.cc Copyright © 2010,2012  Google, Inc.  Behdad Esfahbod

/* Thai / Lao shaper */

var _ otComplexShaper = complexShaperThai{}

type complexShaperThai struct {
	complexShaperNil
}

/* PUA shaping */

// thai_consonant_type_t
const (
	tcNC = iota
	tcAC
	tcRC
	tcDC
	tcNOTCONSONANT
	numConsonantTypes = tcNOTCONSONANT
)

func getConsonantType(u rune) uint8 {
	switch u {
	case 0x0E1B, 0x0E1D, 0x0E1F /* , 0x0E2C*/ :
		return tcAC
	case 0x0E0D, 0x0E10:
		return tcRC
	case 0x0E0E, 0x0E0F:
		return tcDC
	}
	if 0x0E01 <= u && u <= 0x0E2E {
		return tcNC
	}
	return tcNOTCONSONANT
}

// thai_mark_type_t
const (
	tmAV = iota
	tmBV
	tmT
	tmNOTMARK
	numMarkTypes = tmNOTMARK
)

func getMarkType(u rune) uint8 {
	if u == 0x0E31 || (0x0E34 <= u && u <= 0x0E37) ||
		u == 0x0E47 || (0x0E4D <= u && u <= 0x0E4E) {
		return tmAV
	}
	if 0x0E38 <= u && u <= 0x0E3A {
		return tmBV
	}
	if 0x0E48 <= u && u <= 0x0E4C {
		return tmT
	}
	return tmNOTMARK
}

// thai_action_t
const (
	tcNOP = iota
	tcSD  /* Shift combining-mark down */
	tcSL  /* Shift combining-mark left */
	tcSDL /* Shift combining-mark down-left */
	tcRD  /* Remove descender from base */
)

type thaiPuaMapping struct {
	u, winPua, macPua rune
}

var (
	sdMappings = [...]thaiPuaMapping{
		{0x0E48, 0xF70A, 0xF88B}, /* MAI EK */
		{0x0E49, 0xF70B, 0xF88E}, /* MAI THO */
		{0x0E4A, 0xF70C, 0xF891}, /* MAI TRI */
		{0x0E4B, 0xF70D, 0xF894}, /* MAI CHATTAWA */
		{0x0E4C, 0xF70E, 0xF897}, /* THANTHAKHAT */
		{0x0E38, 0xF718, 0xF89B}, /* SARA U */
		{0x0E39, 0xF719, 0xF89C}, /* SARA UU */
		{0x0E3A, 0xF71A, 0xF89D}, /* PHINTHU */
		{0x0000, 0x0000, 0x0000},
	}
	sdlMappings = [...]thaiPuaMapping{
		{0x0E48, 0xF705, 0xF88C}, /* MAI EK */
		{0x0E49, 0xF706, 0xF88F}, /* MAI THO */
		{0x0E4A, 0xF707, 0xF892}, /* MAI TRI */
		{0x0E4B, 0xF708, 0xF895}, /* MAI CHATTAWA */
		{0x0E4C, 0xF709, 0xF898}, /* THANTHAKHAT */
		{0x0000, 0x0000, 0x0000},
	}
	slMappings = [...]thaiPuaMapping{
		{0x0E48, 0xF713, 0xF88A}, /* MAI EK */
		{0x0E49, 0xF714, 0xF88D}, /* MAI THO */
		{0x0E4A, 0xF715, 0xF890}, /* MAI TRI */
		{0x0E4B, 0xF716, 0xF893}, /* MAI CHATTAWA */
		{0x0E4C, 0xF717, 0xF896}, /* THANTHAKHAT */
		{0x0E31, 0xF710, 0xF884}, /* MAI HAN-AKAT */
		{0x0E34, 0xF701, 0xF885}, /* SARA I */
		{0x0E35, 0xF702, 0xF886}, /* SARA II */
		{0x0E36, 0xF703, 0xF887}, /* SARA UE */
		{0x0E37, 0xF704, 0xF888}, /* SARA UEE */
		{0x0E47, 0xF712, 0xF889}, /* MAITAIKHU */
		{0x0E4D, 0xF711, 0xF899}, /* NIKHAHIT */
		{0x0000, 0x0000, 0x0000},
	}
	rdMappings = [...]thaiPuaMapping{
		{0x0E0D, 0xF70F, 0xF89A}, /* YO YING */
		{0x0E10, 0xF700, 0xF89E}, /* THO THAN */
		{0x0000, 0x0000, 0x0000},
	}
)

func thaiPuaShape(u rune, action uint8, font *Font) rune {
	var puaMappings []thaiPuaMapping
	switch action {
	case tcNOP:
		return u
	case tcSD:
		puaMappings = sdMappings[:]
	case tcSDL:
		puaMappings = sdlMappings[:]
	case tcSL:
		puaMappings = slMappings[:]
	case tcRD:
		puaMappings = rdMappings[:]
	}
	for _, pua := range puaMappings {
		if pua.u == u {
			_, ok := font.face.NominalGlyph(pua.winPua)
			if ok {
				return pua.winPua
			}
			_, ok = font.face.NominalGlyph(pua.macPua)
			if ok {
				return pua.macPua
			}
			break
		}
	}
	return u
}

const (
	/* Cluster above looks like: */
	tcT0 = iota /*  ⣤                      */
	tcT1        /*     ⣼                   */
	tcT2        /*        ⣾                */
	tcT3        /*           ⣿             */
	numAboveStates
)

var thaiAboveStartState = [numConsonantTypes + 1] /* For NOT_CONSONANT */ uint8{
	tcT0, /* NC */
	tcT1, /* AC */
	tcT0, /* RC */
	tcT0, /* DC */
	tcT3, /* NOT_CONSONANT */
}

var thaiAboveStateMachine = [numAboveStates][numMarkTypes]struct {
	action    uint8
	nextState uint8
}{ /*AV*/ /*BV*/ /*T*/
	/*T0*/ {{tcNOP, tcT3}, {tcNOP, tcT0}, {tcSD, tcT3}},
	/*T1*/ {{tcSL, tcT2}, {tcNOP, tcT1}, {tcSDL, tcT2}},
	/*T2*/ {{tcNOP, tcT3}, {tcNOP, tcT2}, {tcSL, tcT3}},
	/*T3*/ {{tcNOP, tcT3}, {tcNOP, tcT3}, {tcNOP, tcT3}},
}

// thai_below_state_t
const (
	tbB0 = iota /* No descender */
	tbB1        /* Removable descender */
	tbB2        /* Strict descender */
	numBelowStates
)

var thaiBelowStartState = [numConsonantTypes + 1] /* For NOT_CONSONANT */ uint8{
	tbB0, /* NC */
	tbB0, /* AC */
	tbB1, /* RC */
	tbB2, /* DC */
	tbB2, /* NOT_CONSONANT */
}

var thaiBelowStateMachine = [numBelowStates][numMarkTypes]struct {
	action    uint8
	nextState uint8
}{ /*AV*/ /*BV*/ /*T*/
	/*B0*/ {{tcNOP, tbB0}, {tcNOP, tbB2}, {tcNOP, tbB0}},
	/*B1*/ {{tcNOP, tbB1}, {tcRD, tbB2}, {tcNOP, tbB1}},
	/*B2*/ {{tcNOP, tbB2}, {tcSD, tbB2}, {tcNOP, tbB2}},
}

func doThaiPuaShaping(buffer *Buffer, font *Font) {
	aboveState := thaiAboveStartState[tcNOTCONSONANT]
	belowState := thaiBelowStartState[tcNOTCONSONANT]
	base := 0

	info := buffer.Info
	//    unsigned int count = buffer.len;
	for i := range info {
		mt := getMarkType(info[i].codepoint)

		if mt == tmNOTMARK {
			ct := getConsonantType(info[i].codepoint)
			aboveState = thaiAboveStartState[ct]
			belowState = thaiBelowStartState[ct]
			base = i
			continue
		}

		aboveEdge := &thaiAboveStateMachine[aboveState][mt]
		belowEdge := &thaiBelowStateMachine[belowState][mt]
		aboveState = aboveEdge.nextState
		belowState = belowEdge.nextState

		// at least one of the above/below actions is NOP.
		action := belowEdge.action
		if aboveEdge.action != tcNOP {
			action = aboveEdge.action
		}

		buffer.unsafeToBreak(base, i)
		if action == tcRD {
			info[base].codepoint = thaiPuaShape(info[base].codepoint, action, font)
		} else {
			info[i].codepoint = thaiPuaShape(info[i].codepoint, action, font)
		}
	}
}

/* We only get one script at a time, so a script-agnostic implementation
* is adequate here. */
func isSaraAm(x rune) bool           { return x & ^0x0080 == 0x0E33 }
func nikhahitFromSaraAm(x rune) rune { return x - 0x0E33 + 0x0E4D }
func saraAaFromSaraAm(x rune) rune   { return x - 1 }
func isToneMark(x rune) bool {
	u := x & ^0x0080
	return 0x0E34 <= u && u <= 0x0E37 ||
		0x0E47 <= u && u <= 0x0E4E ||
		0x0E31 <= u && u <= 0x0E31
}

/* This function implements the shaping logic documented here:
 *
 *   https://linux.thai.net/~thep/th-otf/shaping.html
 *
 * The first shaping rule listed there is needed even if the font has Thai
 * OpenType tables.  The rest do fallback positioning based on PUA codepoints.
 * We implement that only if there exist no Thai GSUB in the font.
 */
func (complexShaperThai) preprocessText(plan *otShapePlan, buffer *Buffer, font *Font) {
	/* The following is NOT specified in the MS OT Thai spec, however, it seems
	* to be what Uniscribe and other engines implement.  According to Eric Muller:
	*
	* When you have a SARA AM, decompose it in NIKHAHIT + SARA AA, *and* move the
	* NIKHAHIT backwards over any tone mark (0E48-0E4B).
	*
	* <0E14, 0E4B, 0E33> . <0E14, 0E4D, 0E4B, 0E32>
	*
	* This reordering is legit only when the NIKHAHIT comes from a SARA AM, not
	* when it's there to start with. The string <0E14, 0E4B, 0E4D> is probably
	* not what a user wanted, but the rendering is nevertheless nikhahit above
	* chattawa.
	*
	* Same for Lao.
	*
	* Note:
	*
	* Uniscribe also does some below-marks reordering.  Namely, it positions U+0E3A
	* after U+0E38 and U+0E39.  We do that by modifying the ccc for U+0E3A.
	* See unicode.modified_combining_class ().  Lao does NOT have a U+0E3A
	* equivalent.
	 */

	/*
	* Here are the characters of significance:
	*
	*			Thai	Lao
	* SARA AM:		U+0E33	U+0EB3
	* SARA AA:		U+0E32	U+0EB2
	* Nikhahit:		U+0E4D	U+0ECD
	*
	* Testing shows that Uniscribe reorder the following marks:
	* Thai:	<0E31,0E34..0E37,0E47..0E4E>
	* Lao:	<0EB1,0EB4..0EB7,0EC7..0ECE>
	*
	* Note how the Lao versions are the same as Thai + 0x80.
	 */

	buffer.clearOutput()
	count := len(buffer.Info)
	for buffer.idx = 0; buffer.idx < count; {
		u := buffer.cur(0).codepoint
		if !isSaraAm(u) {
			buffer.nextGlyph()
			continue
		}

		/* Is SARA AM. Decompose and reorder. */
		buffer.outputRune(nikhahitFromSaraAm(u))
		buffer.prev().setContinuation()
		buffer.replaceGlyph(saraAaFromSaraAm(u))

		/* Make Nikhahit be recognized as a ccc=0 mark when zeroing widths. */
		end := len(buffer.outInfo)
		buffer.outInfo[end-2].setGeneralCategory(nonSpacingMark)

		/* Ok, let's see... */
		start := end - 2
		for start > 0 && isToneMark(buffer.outInfo[start-1].codepoint) {
			start--
		}

		if start+2 < end {
			/* Move Nikhahit (end-2) to the beginning */
			buffer.mergeOutClusters(start, end)
			t := buffer.outInfo[end-2]
			copy(buffer.outInfo[start+1:], buffer.outInfo[start:end-2])
			buffer.outInfo[start] = t
		} else {
			/* Since we decomposed, and NIKHAHIT is combining, merge clusters with the
			* previous cluster. */
			if start != 0 && buffer.ClusterLevel == MonotoneGraphemes {
				buffer.mergeOutClusters(start-1, end)
			}
		}
	}
	buffer.swapBuffers()

	/* If font has Thai GSUB, we are done. */
	if plan.props.Script == language.Thai && !plan.map_.foundScript[0] {
		doThaiPuaShaping(buffer, font)
	}
}

func (complexShaperThai) marksBehavior() (zeroWidthMarks, bool) {
	return zeroWidthMarksByGdefLate, false
}

func (complexShaperThai) normalizationPreference() normalizationMode {
	return nmDefault
}
