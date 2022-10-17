// Package unicodedata provides additional lookup functions for unicode
// properties, not covered by the standard package unicode.
package unicodedata

import (
	"unicode"
)

var categories []*unicode.RangeTable

func init() {
	for cat, table := range unicode.Categories {
		if len(cat) == 2 {
			categories = append(categories, table)
		}
	}
}

// LookupType returns the unicode categorie of the rune,
// or nil if not found.
func LookupType(r rune) *unicode.RangeTable {
	for _, table := range categories {
		if unicode.Is(table, r) {
			return table
		}
	}
	return nil
}

// LookupCombiningClass returns the class used for the Canonical Ordering Algorithm in the Unicode Standard,
// defaulting to 0.
//
// From http://www.unicode.org/reports/tr44/#Canonical_Combining_Class:
// "This property could be considered either an enumerated property or a numeric property:
// the principal use of the property is in terms of the numeric values.
// For the property value names associated with different numeric values,
// see DerivedCombiningClass.txt and Canonical Combining Class Values."
func LookupCombiningClass(ch rune) uint8 {
	for i, t := range combiningClasses {
		if t == nil {
			continue
		}
		if unicode.Is(t, ch) {
			return uint8(i)
		}
	}
	return 0
}

// LookupLineBreakClass returns the break class for the rune (see the constants BreakXXX)
func LookupLineBreakClass(ch rune) *unicode.RangeTable {
	for _, class := range lineBreaks {
		if unicode.Is(class, ch) {
			return class
		}
	}
	return BreakXX
}

// LookupGraphemeBreakClass returns the grapheme break property for the rune (see the constants GraphemeBreakXXX),
// or nil
func LookupGraphemeBreakClass(ch rune) *unicode.RangeTable {
	// a lot of runes do not have a grapheme break property :
	// avoid testing all the graphemeBreaks classes for them
	if !unicode.Is(graphemeBreakAll, ch) {
		return nil
	}
	for _, class := range graphemeBreaks {
		if unicode.Is(class, ch) {
			return class
		}
	}
	return nil
}

// LookupMirrorChar finds the mirrored equivalent of a character as defined in
// the file BidiMirroring.txt of the Unicode Character Database available at
// http://www.unicode.org/Public/UNIDATA/BidiMirroring.txt.
//
// If the input character is declared as a mirroring character in the
// Unicode standard and has a mirrored equivalent, it is returned with `true`.
// Otherwise the input character itself returned with `false`.
func LookupMirrorChar(ch rune) (rune, bool) {
	m, ok := mirroring[ch]
	if !ok {
		m = ch
	}
	return m, ok
}

// Algorithmic hangul syllable [de]composition
const (
	HangulSBase  = 0xAC00
	HangulLBase  = 0x1100
	HangulVBase  = 0x1161
	HangulTBase  = 0x11A7
	HangulSCount = 11172
	HangulLCount = 19
	HangulVCount = 21
	HangulTCount = 28
	HangulNCount = HangulVCount * HangulTCount
)

func decomposeHangul(ab rune) (a, b rune, ok bool) {
	si := ab - HangulSBase

	if si < 0 || si >= HangulSCount {
		return 0, 0, false
	}

	if si%HangulTCount != 0 { /* LV,T */
		return HangulSBase + (si/HangulTCount)*HangulTCount, HangulTBase + (si % HangulTCount), true
	} /* L,V */
	return HangulLBase + (si / HangulNCount), HangulVBase + (si%HangulNCount)/HangulTCount, true
}

func composeHangul(a, b rune) (rune, bool) {
	if a >= HangulSBase && a < (HangulSBase+HangulSCount) && b > HangulTBase && b < (HangulTBase+HangulTCount) && (a-HangulSBase)%HangulTCount == 0 {
		/* LV,T */
		return a + (b - HangulTBase), true
	} else if a >= HangulLBase && a < (HangulLBase+HangulLCount) && b >= HangulVBase && b < (HangulVBase+HangulVCount) {
		/* L,V */
		li := a - HangulLBase
		vi := b - HangulVBase
		return HangulSBase + li*HangulNCount + vi*HangulTCount, true
	}
	return 0, false
}

// Decompose decompose an input Unicode code point,
// returning the two decomposed code points, if successful.
// It returns `false` otherwise.
func Decompose(ab rune) (a, b rune, ok bool) {
	if a, b, ok = decomposeHangul(ab); ok {
		return a, b, true
	}
	if m1, ok := decompose1[ab]; ok {
		return m1, 0, true
	}
	if m2, ok := decompose2[ab]; ok {
		return m2[0], m2[1], true
	}
	return ab, 0, false
}

// Compose composes a sequence of two input Unicode code
// points by canonical equivalence, returning the composed code, if successful.
// It returns `false` otherwise
func Compose(a, b rune) (rune, bool) {
	if ab, ok := composeHangul(a, b); ok {
		return ab, true
	}
	u := compose[[2]rune{a, b}]
	return u, u != 0
}

// ArabicJoining is a property used to shape Arabic runes.
// See the table ArabicJoinings.
type ArabicJoining byte

const (
	U          ArabicJoining = 'U' // Un-joining, e.g. Full Stop
	R          ArabicJoining = 'R' // Right-joining, e.g. Arabic Letter Dal
	Alaph      ArabicJoining = 'a' // Alaph group (included in kind R)
	DalathRish ArabicJoining = 'd' // Dalat Rish group (included in kind R)
	D          ArabicJoining = 'D' // Dual-joining, e.g. Arabic Letter Ain
	C          ArabicJoining = 'C' // Join-Causing, e.g. Tatweel, ZWJ
	L          ArabicJoining = 'L' // Left-joining, i.e. fictional
	T          ArabicJoining = 'T' // Transparent, e.g. Arabic Fatha
	G          ArabicJoining = 'G' // Ignored, e.g. LRE, RLE, ZWNBSP
)

// JamoType is an enum that works as the states of the Hangul syllables system.
type JamoType uint8

const (
	JAMO_LV  JamoType = iota /* break HANGUL_LV_SYLLABLE */
	JAMO_LVT                 /* break HANGUL_LVT_SYLLABLE */
	JAMO_L                   /* break HANGUL_L_JAMO */
	JAMO_V                   /* break HANGUL_V_JAMO */
	JAMO_T                   /* break HANGUL_T_JAMO */
	NO_JAMO                  /* Other */
)

// There are Hangul syllables encoded as characters, that act like a
// sequence of Jamos. For each character we define a JamoType
// that the character starts with, and one that it ends with.  This
// decomposes JAMO_LV and JAMO_LVT to simple other JAMOs.  So for
// example, a character with LineBreak type
// BreakH2 has start=JAMO_L and end=JAMO_V.
type charJamoProps struct {
	Start, End JamoType
}

// HangulJamoProps maps from JamoType to CharJamoProps that hold only simple
// JamoTypes (no LV or LVT) or none.
var HangulJamoProps = [6]charJamoProps{
	{Start: JAMO_L, End: JAMO_V},   /* JAMO_LV */
	{Start: JAMO_L, End: JAMO_T},   /* JAMO_LVT */
	{Start: JAMO_L, End: JAMO_L},   /* JAMO_L */
	{Start: JAMO_V, End: JAMO_V},   /* JAMO_V */
	{Start: JAMO_T, End: JAMO_T},   /* JAMO_T */
	{Start: NO_JAMO, End: NO_JAMO}, /* NO_JAMO */
}

// Jamo returns the Jamo Type of `btype` or NO_JAMO
func Jamo(btype *unicode.RangeTable) JamoType {
	switch btype {
	case BreakH2:
		return JAMO_LV
	case BreakH3:
		return JAMO_LVT
	case BreakJL:
		return JAMO_L
	case BreakJV:
		return JAMO_V
	case BreakJT:
		return JAMO_T
	default:
		return NO_JAMO
	}
}
