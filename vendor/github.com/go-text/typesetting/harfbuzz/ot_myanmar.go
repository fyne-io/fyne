package harfbuzz

import (
	"fmt"

	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
)

// ported from harfbuzz/src/hb-ot-shape-complex-myanmar.cc, .hh Copyright © 2011,2012,2013  Google, Inc.  Behdad Esfahbod

// Myanmar shaper.
type complexShaperMyanmar struct {
	complexShaperNil
}

var _ otComplexShaper = complexShaperMyanmar{}

/*
 * Basic features.
 * These features are applied in order, one at a time, after reordering.
 */
var myanmarBasicFeatures = [...]tables.Tag{
	loader.NewTag('r', 'p', 'h', 'f'),
	loader.NewTag('p', 'r', 'e', 'f'),
	loader.NewTag('b', 'l', 'w', 'f'),
	loader.NewTag('p', 's', 't', 'f'),
}

/*
* Other features.
* These features are applied all at once, after clearing syllables.
 */
var myanmarOtherFeatures = [...]tables.Tag{
	loader.NewTag('p', 'r', 'e', 's'),
	loader.NewTag('a', 'b', 'v', 's'),
	loader.NewTag('b', 'l', 'w', 's'),
	loader.NewTag('p', 's', 't', 's'),
}

func (complexShaperMyanmar) collectFeatures(plan *otShapePlanner) {
	map_ := &plan.map_

	/* Do this before any lookups have been applied. */
	map_.addGSUBPause(setupSyllablesMyanmar)

	map_.enableFeature(loader.NewTag('l', 'o', 'c', 'l'))
	/* The Indic specs do not require ccmp, but we apply it here since if
	* there is a use of it, it's typically at the beginning. */
	map_.enableFeature(loader.NewTag('c', 'c', 'm', 'p'))

	map_.addGSUBPause(reorderMyanmar)

	for _, feat := range myanmarBasicFeatures {
		map_.enableFeatureExt(feat, ffManualZWJ, 1)
		map_.addGSUBPause(nil)
	}

	map_.addGSUBPause(clearSyllables)

	for _, feat := range myanmarOtherFeatures {
		map_.enableFeatureExt(feat, ffManualZWJ, 1)
	}
}

func (complexShaperMyanmar) setupMasks(_ *otShapePlan, buffer *Buffer, _ *Font) {
	/* We cannot setup masks here.  We save information about characters
	* and setup masks later on in a pause-callback. */

	info := buffer.Info
	for i := range info {
		setMyanmarProperties(&info[i])
	}
}

func foundSyllableMyanmar(syllableType uint8, ts, te int, info []GlyphInfo, syllableSerial *uint8) {
	for i := ts; i < te; i++ {
		info[i].syllable = (*syllableSerial << 4) | syllableType
	}
	*syllableSerial++
	if *syllableSerial == 16 {
		*syllableSerial = 1
	}
}

func setupSyllablesMyanmar(_ *otShapePlan, _ *Font, buffer *Buffer) {
	findSyllablesMyanmar(buffer)
	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		buffer.unsafeToBreak(start, end)
	}
}

/* Rules from:
 * https://docs.microsoft.com/en-us/typography/script-development/myanmar */
func initialReorderingConsonantSyllableMyanmar(buffer *Buffer, start, end int) {
	info := buffer.Info

	base := end
	hasReph := false

	limit := start
	if start+3 <= end &&
		info[start].complexCategory == otRa &&
		info[start+1].complexCategory == otAs &&
		info[start+2].complexCategory == otH {
		limit += 3
		base = start
		hasReph = true
	}

	if !hasReph {
		base = limit
	}

	for i := limit; i < end; i++ {
		if isConsonant(&info[i]) {
			base = i
			break
		}
	}

	/* Reorder! */
	i := start
	endLoop := start
	if hasReph {
		endLoop = start + 3
	}
	for ; i < endLoop; i++ {
		info[i].complexAux = posAfterMain
	}
	for ; i < base; i++ {
		info[i].complexAux = posPreC
	}
	if i < end {
		info[i].complexAux = posBaseC
		i++
	}
	var pos uint8 = posAfterMain
	/* The following loop may be ugly, but it implements all of
	 * Myanmar reordering! */
	for ; i < end; i++ {
		if info[i].complexCategory == otMR /* Pre-base reordering */ {
			info[i].complexAux = posPreC
			continue
		}
		if info[i].complexAux < posBaseC /* Left matra */ {
			continue
		}
		if info[i].complexCategory == otVS {
			info[i].complexAux = info[i-1].complexAux
			continue
		}

		if pos == posAfterMain && info[i].complexCategory == otVBlw {
			pos = posBelowC
			info[i].complexAux = pos
			continue
		}

		if pos == posBelowC && info[i].complexCategory == otA {
			info[i].complexAux = posBeforeSub
			continue
		}
		if pos == posBelowC && info[i].complexCategory == otVBlw {
			info[i].complexAux = pos
			continue
		}
		if pos == posBelowC && info[i].complexCategory != otA {
			pos = posAfterSub
			info[i].complexAux = pos
			continue
		}
		info[i].complexAux = pos
	}

	/* Sit tight, rock 'n roll! */
	buffer.sort(start, end, func(a, b *GlyphInfo) int { return int(a.complexAux) - int(b.complexAux) })
}

func reorderSyllableMyanmar(buffer *Buffer, start, end int) {
	syllableType := buffer.Info[start].syllable & 0x0F
	switch syllableType {
	/* We already inserted dotted-circles, so just call the consonant_syllable. */
	case myanmarBrokenCluster, myanmarConsonantSyllable:
		initialReorderingConsonantSyllableMyanmar(buffer, start, end)
	}
}

func reorderMyanmar(_ *otShapePlan, font *Font, buffer *Buffer) {
	if debugMode >= 1 {
		fmt.Println("MYANMAR - start reordering myanmar")
	}

	syllabicInsertDottedCircles(font, buffer, myanmarBrokenCluster, otGB, -1, -1)

	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		reorderSyllableMyanmar(buffer, start, end)
	}

	if debugMode >= 1 {
		fmt.Println("MYANMAR - end reordering myanmar")
	}
}

/* Note: This enum is duplicated in the -machine.rl source file.
 * Not sure how to avoid duplication. */
const (
	otAs = 18  /* Asat */
	otD0 = 20  /* Digit zero */
	otDB = otN /* Dot below */
	otGB = otPLACEHOLDER
	otMH = 21 /* Various consonant medial types */
	otMR = 22 /* Various consonant medial types */
	otMW = 23 /* Various consonant medial types */
	otMY = 24 /* Various consonant medial types */
	otPT = 25 /* Pwo and other tones */
	// otVAbv = 26
	// otVBlw = 27
	// otVPre = 28
	// otVPst = 29
	otVS = 30 /* Variation selectors */
	otP  = 31 /* Punctuation */
	otD  = 32 /* Digits except zero */
	otML = 33 /* Various consonant medial types */
)

func computeMyanmarProperties(u rune) (cat, pos uint8) {
	type_ := indicGetCategories(u)
	cat = uint8(type_ & 0xFF)
	pos = uint8(type_ >> 8)

	/* Myanmar
	* https://docs.microsoft.com/en-us/typography/script-development/myanmar#analyze */
	if 0xFE00 <= u && u <= 0xFE0F {
		cat = otVS
	}

	switch u {
	case 0x104E:
		cat = otC /* The spec says C, IndicSyllableCategory doesn't have. */
	case 0x002D, 0x00A0, 0x00D7, 0x2012, 0x2013, 0x2014, 0x2015, 0x2022,
		0x25CC, 0x25FB, 0x25FC, 0x25FD, 0x25FE:
		cat = otGB
	case 0x1004, 0x101B, 0x105A:
		cat = otRa
	case 0x1032, 0x1036:
		cat = otA
	case 0x1039:
		cat = otH
	case 0x103A:
		cat = otAs
	case 0x1041, 0x1042, 0x1043, 0x1044, 0x1045, 0x1046, 0x1047, 0x1048,
		0x1049, 0x1090, 0x1091, 0x1092, 0x1093, 0x1094, 0x1095, 0x1096, 0x1097, 0x1098, 0x1099:
		cat = otD
	case 0x1040:
		cat = otD /* The spec says D0, but Uniscribe doesn't seem to do. */
	case 0x103E:
		cat = otMH
	case 0x1060:
		cat = otML
	case 0x103C:
		cat = otMR
	case 0x103D, 0x1082:
		cat = otMW
	case 0x103B, 0x105E, 0x105F:
		cat = otMY
	case 0x1063, 0x1064, 0x1069, 0x106A, 0x106B, 0x106C, 0x106D, 0xAA7B:
		cat = otPT
	case 0x1038, 0x1087, 0x1088, 0x1089, 0x108A, 0x108B, 0x108C, 0x108D,
		0x108F, 0x109A, 0x109B, 0x109C:
		cat = otSM
	case 0x104A, 0x104B:
		cat = otP
	case 0xAA74, 0xAA75, 0xAA76:
		/* https://github.com/harfbuzz/harfbuzz/issues/218 */
		cat = otC
	}

	if cat == otM {
		switch pos {
		case posPreC:
			cat = otVPre
			pos = posPreM
		case posAboveC:
			cat = otVAbv
		case posBelowC:
			cat = otVBlw
		case posPostC:
			cat = otVPst
		}
	}

	return cat, pos
}

func setMyanmarProperties(info *GlyphInfo) {
	u := info.codepoint
	cat, pos := computeMyanmarProperties(u)
	info.complexCategory = cat
	info.complexAux = pos
}

func (complexShaperMyanmar) marksBehavior() (zeroWidthMarks, bool) {
	return zeroWidthMarksByGdefEarly, false
}

func (complexShaperMyanmar) normalizationPreference() normalizationMode {
	return nmComposedDiacriticsNoShortCircuit
}
