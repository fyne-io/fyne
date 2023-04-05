package harfbuzz

import (
	"unicode"

	"github.com/go-text/typesetting/unicodedata"
)

// uni exposes some lookup functions for Unicode properties.
var uni = unicodeFuncs{}

// generalCategory is an enum value to allow compact storage (see generalCategories)
type generalCategory uint8

const (
	control generalCategory = iota
	format
	unassigned
	privateUse
	surrogate
	lowercaseLetter
	modifierLetter
	otherLetter
	titlecaseLetter
	uppercaseLetter
	spacingMark
	enclosingMark
	nonSpacingMark
	decimalNumber
	letterNumber
	otherNumber
	connectPunctuation
	dashPunctuation
	closePunctuation
	finalPunctuation
	initialPunctuation
	otherPunctuation
	openPunctuation
	currencySymbol
	modifierSymbol
	mathSymbol
	otherSymbol
	lineSeparator
	paragraphSeparator
	spaceSeparator
)

// correspondance with *unicode.RangeTable classes
var generalCategories = [...]*unicode.RangeTable{
	control:            unicode.Cc,
	format:             unicode.Cf,
	unassigned:         nil,
	privateUse:         unicode.Co,
	surrogate:          unicode.Cs,
	lowercaseLetter:    unicode.Ll,
	modifierLetter:     unicode.Lm,
	otherLetter:        unicode.Lo,
	titlecaseLetter:    unicode.Lt,
	uppercaseLetter:    unicode.Lu,
	spacingMark:        unicode.Mc,
	enclosingMark:      unicode.Me,
	nonSpacingMark:     unicode.Mn,
	decimalNumber:      unicode.Nd,
	letterNumber:       unicode.Nl,
	otherNumber:        unicode.No,
	connectPunctuation: unicode.Pc,
	dashPunctuation:    unicode.Pd,
	closePunctuation:   unicode.Pe,
	finalPunctuation:   unicode.Pf,
	initialPunctuation: unicode.Pi,
	otherPunctuation:   unicode.Po,
	openPunctuation:    unicode.Ps,
	currencySymbol:     unicode.Sc,
	modifierSymbol:     unicode.Sk,
	mathSymbol:         unicode.Sm,
	otherSymbol:        unicode.So,
	lineSeparator:      unicode.Zl,
	paragraphSeparator: unicode.Zp,
	spaceSeparator:     unicode.Zs,
}

func (g generalCategory) isMark() bool {
	return g == spacingMark || g == enclosingMark || g == nonSpacingMark
}

func (g generalCategory) isLetter() bool {
	return g == lowercaseLetter || g == modifierLetter || g == otherLetter ||
		g == titlecaseLetter || g == uppercaseLetter
}

// Modified combining marks
const (
	/* Hebrew
	 *
	 * We permute the "fixed-position" classes 10-26 into the order
	 * described in the SBL Hebrew manual:
	 *
	 * https://www.sbl-site.org/Fonts/SBLHebrewUserManual1.5x.pdf
	 *
	 * (as recommended by:
	 *  https://forum.fontlab.com/archive-old-microsoft-volt-group/vista-and-diacritic-ordering/msg22823/)
	 *
	 * More details here:
	 * https://bugzilla.mozilla.org/show_bug.cgi?id=662055
	 */
	mcc10 uint8 = 22 /* sheva */
	mcc11 uint8 = 15 /* hataf segol */
	mcc12 uint8 = 16 /* hataf patah */
	mcc13 uint8 = 17 /* hataf qamats */
	mcc14 uint8 = 23 /* hiriq */
	mcc15 uint8 = 18 /* tsere */
	mcc16 uint8 = 19 /* segol */
	mcc17 uint8 = 20 /* patah */
	mcc18 uint8 = 21 /* qamats & qamats qatan */
	mcc19 uint8 = 14 /* holam & holam haser for vav*/
	mcc20 uint8 = 24 /* qubuts */
	mcc21 uint8 = 12 /* dagesh */
	mcc22 uint8 = 25 /* meteg */
	mcc23 uint8 = 13 /* rafe */
	mcc24 uint8 = 10 /* shin dot */
	mcc25 uint8 = 11 /* sin dot */
	mcc26 uint8 = 26 /* point varika */

	/*
	 * Arabic
	 *
	 * Modify to move Shadda (ccc=33) before other marks.  See:
	 * https://unicode.org/faq/normalization.html#8
	 * https://unicode.org/faq/normalization.html#9
	 */
	mcc27 uint8 = 28 /* fathatan */
	mcc28 uint8 = 29 /* dammatan */
	mcc29 uint8 = 30 /* kasratan */
	mcc30 uint8 = 31 /* fatha */
	mcc31 uint8 = 32 /* damma */
	mcc32 uint8 = 33 /* kasra */
	mcc33 uint8 = 27 /* shadda */
	mcc34 uint8 = 34 /* sukun */
	mcc35 uint8 = 35 /* superscript alef */

	/* Syriac */
	mcc36 uint8 = 36 /* superscript alaph */

	/* Telugu
	 *
	 * Modify Telugu length marks (ccc=84, ccc=91).
	 * These are the only matras in the main Indic scripts range that have
	 * a non-zero ccc.  That makes them reorder with the Halant (ccc=9).
	 * Assign 4 and 5, which are otherwise unassigned.
	 */
	mcc84 uint8 = 4 /* length mark */
	mcc91 uint8 = 5 /* ai length mark */

	/* Thai
	 *
	 * Modify U+0E38 and U+0E39 (ccc=103) to be reordered before U+0E3A (ccc=9).
	 * Assign 3, which is unassigned otherwise.
	 * Uniscribe does this reordering too.
	 */
	mcc103 uint8 = 3   /* sara u / sara uu */
	mcc107 uint8 = 107 /* mai * */

	/* Lao */
	mcc118 uint8 = 118 /* sign u / sign uu */
	mcc122 uint8 = 122 /* mai * */

	/* Tibetan
	 *
	 * In case of multiple vowel-signs, use u first (but after achung)
	 * this allows Dzongkha multi-vowel shortcuts to render correctly
	 */
	mcc129 = 129 /* sign aa */
	mcc130 = 132 /* sign i */
	mcc132 = 131 /* sign u */
)

var modifiedCombiningClass = [256]uint8{
	0, /* HB_UNICODE_COMBINING_CLASS_NOT_REORDERED */
	1, /* HB_UNICODE_COMBINING_CLASS_OVERLAY */
	2, 3, 4, 5, 6,
	7, /* HB_UNICODE_COMBINING_CLASS_NUKTA */
	8, /* HB_UNICODE_COMBINING_CLASS_KANA_VOICING */
	9, /* HB_UNICODE_COMBINING_CLASS_VIRAMA */

	/* Hebrew */
	mcc10,
	mcc11,
	mcc12,
	mcc13,
	mcc14,
	mcc15,
	mcc16,
	mcc17,
	mcc18,
	mcc19,
	mcc20,
	mcc21,
	mcc22,
	mcc23,
	mcc24,
	mcc25,
	mcc26,

	/* Arabic */
	mcc27,
	mcc28,
	mcc29,
	mcc30,
	mcc31,
	mcc32,
	mcc33,
	mcc34,
	mcc35,

	/* Syriac */
	mcc36,

	37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
	80, 81, 82, 83,

	/* Telugu */
	mcc84,
	85, 86, 87, 88, 89, 90,
	mcc91,
	92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102,

	/* Thai */
	mcc103,
	104, 105, 106,
	mcc107,
	108, 109, 110, 111, 112, 113, 114, 115, 116, 117,

	/* Lao */
	mcc118,
	119, 120, 121,
	mcc122,
	123, 124, 125, 126, 127, 128,

	/* Tibetan */
	mcc129,
	mcc130,
	131,
	mcc132,
	133, 134, 135, 136, 137, 138, 139,

	140, 141, 142, 143, 144, 145, 146, 147, 148, 149,
	150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
	160, 161, 162, 163, 164, 165, 166, 167, 168, 169,
	170, 171, 172, 173, 174, 175, 176, 177, 178, 179,
	180, 181, 182, 183, 184, 185, 186, 187, 188, 189,
	190, 191, 192, 193, 194, 195, 196, 197, 198, 199,

	200, /* HB_UNICODE_COMBINING_CLASS_ATTACHED_BELOW_LEFT */
	201,
	202, /* HB_UNICODE_COMBINING_CLASS_ATTACHED_BELOW */
	203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213,
	214, /* HB_UNICODE_COMBINING_CLASS_ATTACHED_ABOVE */
	215,
	216, /* HB_UNICODE_COMBINING_CLASS_ATTACHED_ABOVE_RIGHT */
	217,
	218, /* HB_UNICODE_COMBINING_CLASS_BELOW_LEFT */
	219,
	220, /* HB_UNICODE_COMBINING_CLASS_BELOW */
	221,
	222, /* HB_UNICODE_COMBINING_CLASS_BELOW_RIGHT */
	223,
	224, /* HB_UNICODE_COMBINING_CLASS_LEFT */
	225,
	226, /* HB_UNICODE_COMBINING_CLASS_RIGHT */
	227,
	228, /* HB_UNICODE_COMBINING_CLASS_ABOVE_LEFT */
	229,
	230, /* HB_UNICODE_COMBINING_CLASS_ABOVE */
	231,
	232, /* HB_UNICODE_COMBINING_CLASS_ABOVE_RIGHT */
	233, /* HB_UNICODE_COMBINING_CLASS_DOUBLE_BELOW */
	234, /* HB_UNICODE_COMBINING_CLASS_DOUBLE_ABOVE */
	235, 236, 237, 238, 239,
	240, /* HB_UNICODE_COMBINING_CLASS_IOTA_SUBSCRIPT */
	241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254,
	255, /* HB_UNICODE_COMBINING_CLASS_INVALID */
}

type unicodeFuncs struct{}

func (unicodeFuncs) modifiedCombiningClass(u rune) uint8 {
	/* This hack belongs to the USE shaper (for Tai Tham):
	 * Reorder SAKOT to ensure it comes after any tone marks. */
	if u == 0x1A60 {
		return 254
	}

	/* This hack belongs to the Tibetan shaper:
	 * Reorder PADMA to ensure it comes after any vowel marks. */
	if u == 0x0FC6 {
		return 254
	}
	/* Reorder TSA -PHRU to reorder before U+0F74 */
	if u == 0x0F39 {
		return 127
	}
	return modifiedCombiningClass[unicodedata.LookupCombiningClass(u)]
}

// IsDefaultIgnorable returns `true` for
// codepoints with the Default_Ignorable property
// (as defined in unicode data DerivedCoreProperties.txt)
func IsDefaultIgnorable(ch rune) bool {
	// Note: While U+115F, U+1160, U+3164 and U+FFA0 are Default_Ignorable,
	// we do NOT want to hide them, as the way Uniscribe has implemented them
	// is with regular spacing glyphs, and that's the way fonts are made to work.
	// As such, we make exceptions for those four.
	// Also ignoring U+1BCA0..1BCA3. https://github.com/harfbuzz/harfbuzz/issues/503
	plane := ch >> 16
	if plane == 0 {
		/* BMP */
		page := ch >> 8
		switch page {
		case 0x00:
			return (ch == 0x00AD)
		case 0x03:
			return (ch == 0x034F)
		case 0x06:
			return (ch == 0x061C)
		case 0x17:
			return 0x17B4 <= ch && ch <= 0x17B5
		case 0x18:
			return 0x180B <= ch && ch <= 0x180E
		case 0x20:
			return 0x200B <= ch && ch <= 0x200F ||
				0x202A <= ch && ch <= 0x202E ||
				0x2060 <= ch && ch <= 0x206F
		case 0xFE:
			return 0xFE00 <= ch && ch <= 0xFE0F || ch == 0xFEFF
		case 0xFF:
			return 0xFFF0 <= ch && ch <= 0xFFF8
		default:
			return false
		}
	} else {
		/* Other planes */
		switch plane {
		case 0x01:
			return 0x1D173 <= ch && ch <= 0x1D17A
		case 0x0E:
			return 0xE0000 <= ch && ch <= 0xE0FFF
		default:
			return false
		}
	}
}

func (unicodeFuncs) isDefaultIgnorable(ch rune) bool {
	return IsDefaultIgnorable(ch)
}

// retrieves the General Category property for
// a specified Unicode code point, expressed as enumeration value.
func (unicodeFuncs) generalCategory(ch rune) generalCategory {
	for i, cat := range generalCategories {
		if cat != nil && unicode.Is(cat, ch) {
			return generalCategory(i)
		}
	}
	return unassigned
}

func (unicodeFuncs) isExtendedPictographic(ch rune) bool {
	return unicode.Is(unicodedata.Extended_Pictographic, ch)
}

// returns the mirroring Glyph code point (for bi-directional
// replacement) of a code point, or itself
func (unicodeFuncs) mirroring(ch rune) rune {
	out, _ := unicodedata.LookupMirrorChar(ch)
	return out
}

/* Space estimates based on:
 * https://unicode.org/charts/PDF/U2000.pdf
 * https://docs.microsoft.com/en-us/typography/develop/character-design-standards/whitespace
 */
const (
	spaceEM16  = 16 + iota
	space4EM18 // 4/18th of an EM!
	space
	spaceFigure
	spacePunctuation
	spaceNarrow
	notSpace = 0
	spaceEM  = 1
	spaceEM2 = 2
	spaceEM3 = 3
	spaceEM4 = 4
	spaceEM5 = 5
	spaceEM6 = 6
)

func (unicodeFuncs) spaceFallbackType(u rune) uint8 {
	switch u {
	// all GC=Zs chars that can use a fallback.
	case 0x0020:
		return space /* U+0020 SPACE */
	case 0x00A0:
		return space /* U+00A0 NO-BREAK SPACE */
	case 0x2000:
		return spaceEM2 /* U+2000 EN QUAD */
	case 0x2001:
		return spaceEM /* U+2001 EM QUAD */
	case 0x2002:
		return spaceEM2 /* U+2002 EN SPACE */
	case 0x2003:
		return spaceEM /* U+2003 EM SPACE */
	case 0x2004:
		return spaceEM3 /* U+2004 THREE-PER-EM SPACE */
	case 0x2005:
		return spaceEM4 /* U+2005 FOUR-PER-EM SPACE */
	case 0x2006:
		return spaceEM6 /* U+2006 SIX-PER-EM SPACE */
	case 0x2007:
		return spaceFigure /* U+2007 FIGURE SPACE */
	case 0x2008:
		return spacePunctuation /* U+2008 PUNCTUATION SPACE */
	case 0x2009:
		return spaceEM5 /* U+2009 THIN SPACE */
	case 0x200A:
		return spaceEM16 /* U+200A HAIR SPACE */
	case 0x202F:
		return spaceNarrow /* U+202F NARROW NO-BREAK SPACE */
	case 0x205F:
		return space4EM18 /* U+205F MEDIUM MATHEMATICAL SPACE */
	case 0x3000:
		return spaceEM /* U+3000 IDEOGRAPHIC SPACE */
	default:
		return notSpace /* U+1680 OGHAM SPACE MARK */
	}
}

func (unicodeFuncs) isVariationSelector(r rune) bool {
	/* U+180B..180D, U+180F MONGOLIAN FREE VARIATION SELECTORs are handled in the
	 * Arabic shaper.  No need to match them here. */
	/* VARIATION SELECTOR-1..16 */
	/* VARIATION SELECTOR-17..256 */
	return (0xFE00 <= r && r <= 0xFE0F) || (0xE0100 <= r && r <= 0xE01EF)
}

func (unicodeFuncs) decompose(ab rune) (a, b rune, ok bool) { return unicodedata.Decompose(ab) }
func (unicodeFuncs) compose(a, b rune) (rune, bool)         { return unicodedata.Compose(a, b) }

/* Prepare */

/* Implement enough of Unicode Graphemes here that shaping
 * in reverse-direction wouldn't break graphemes.  Namely,
 * we mark all marks and ZWJ and ZWJ,Extended_Pictographic
 * sequences as continuations.  The foreach_grapheme()
 * macro uses this bit.
 *
 * https://www.unicode.org/reports/tr29/#Regex_Definitions
 */
func (b *Buffer) setUnicodeProps() {
	info := b.Info
	for i := 0; i < len(info); i++ {
		info[i].setUnicodeProps(b)

		/* Marks are already set as continuation by the above line.
		 * Handle Emoji_Modifier and ZWJ-continuation. */
		if info[i].unicode.generalCategory() == modifierSymbol && (0x1F3FB <= info[i].codepoint && info[i].codepoint <= 0x1F3FF) {
			info[i].setContinuation()
		} else if i != 0 && 0x1F1E6 <= info[i].codepoint && info[i].codepoint <= 0x1F1FF {
			/* Regional_Indicators are hairy as hell...
			* https://github.com/harfbuzz/harfbuzz/issues/2265 */
			if 0x1F1E6 <= info[i-1].codepoint && info[i-1].codepoint <= 0x1F1FF && !info[i-1].isContinuation() {
				info[i].setContinuation()
			}
		} else if info[i].isZwj() {
			info[i].setContinuation()
			if i+1 < len(b.Info) && uni.isExtendedPictographic(info[i+1].codepoint) {
				i++
				info[i].setUnicodeProps(b)
				info[i].setContinuation()
			}
		} else if 0xE0020 <= info[i].codepoint && info[i].codepoint <= 0xE007F {
			/* Or part of the Other_Grapheme_Extend that is not marks.
			 * As of Unicode 11 that is just:
			 *
			 * 200C          ; Other_Grapheme_Extend # Cf       ZERO WIDTH NON-JOINER
			 * FF9E..FF9F    ; Other_Grapheme_Extend # Lm   [2] HALFWIDTH KATAKANA VOICED SOUND MARK..HALFWIDTH KATAKANA SEMI-VOICED SOUND MARK
			 * E0020..E007F  ; Other_Grapheme_Extend # Cf  [96] TAG SPACE..CANCEL TAG
			 *
			 * ZWNJ is special, we don't want to merge it as there's no need, and keeping
			 * it separate results in more granular clusters.  Ignore Katakana for now.
			 * Tags are used for Emoji sub-region flag sequences:
			 * https://github.com/harfbuzz/harfbuzz/issues/1556
			 */
			info[i].setContinuation()
		}
	}
}

func (b *Buffer) insertDottedCircle(font *Font) {
	if b.Flags&DoNotinsertDottedCircle != 0 {
		return
	}

	if b.Flags&Bot == 0 || len(b.context[0]) != 0 ||
		len(b.Info) == 0 || !b.Info[0].isUnicodeMark() {
		return
	}

	if !font.hasGlyph(0x25CC) {
		return
	}

	dottedcircle := GlyphInfo{codepoint: 0x25CC}
	dottedcircle.setUnicodeProps(b)

	b.clearOutput()

	b.idx = 0
	dottedcircle.Cluster = b.cur(0).Cluster
	dottedcircle.Mask = b.cur(0).Mask
	b.outInfo = append(b.outInfo, dottedcircle)
	b.swapBuffers()
}

func (b *Buffer) formClusters() {
	if b.scratchFlags&bsfHasNonASCII == 0 {
		return
	}

	iter, count := b.graphemesIterator()

	if b.ClusterLevel == MonotoneGraphemes {
		for start, end := iter.next(); start < count; start, end = iter.next() {
			b.mergeClusters(start, end)
		}
	} else {
		for start, end := iter.next(); start < count; start, end = iter.next() {
			b.unsafeToBreak(start, end)
		}
	}
}

func (b *Buffer) ensureNativeDirection() {
	direction := b.Props.Direction
	horizDir := getHorizontalDirection(b.Props.Script)

	/* Numeric runs in natively-RTL scripts are actually native-LTR, so we reset
	 * the horiz_dir if the run contains at least one decimal-number char, and no
	 * letter chars (ideally we should be checking for chars with strong
	 * directionality but hb-unicode currently lacks bidi categories).
	 *
	 * This allows digit sequences in Arabic etc to be shaped in "native"
	 * direction, so that features like ligatures will work as intended.
	 *
	 * https://github.com/harfbuzz/harfbuzz/issues/501
	 */

	if horizDir == RightToLeft && direction == LeftToRight {
		var foundNumber, foundLetter bool
		for _, info := range b.Info {
			gc := info.unicode.generalCategory()
			if gc == decimalNumber {
				foundNumber = true
			} else if gc.isLetter() {
				foundLetter = true
				break
			}
		}
		if foundNumber && !foundLetter {
			horizDir = LeftToRight
		}
	}

	if (direction.isHorizontal() && direction != horizDir && horizDir != 0) ||
		(direction.isVertical() && direction != TopToBottom) {

		reverseGraphemes(b)

		b.Props.Direction = b.Props.Direction.Reverse()
	}
}

// the returned flag must be ORed with the current
func computeUnicodeProps(u rune) (unicodeProp, bufferScratchFlags) {
	genCat := uni.generalCategory(u)
	props := unicodeProp(genCat)
	var flags bufferScratchFlags
	if u >= 0x80 {
		flags |= bsfHasNonASCII

		if uni.isDefaultIgnorable(u) {
			flags |= bsfHasDefaultIgnorables
			props |= upropsMaskIgnorable
			if u == 0x200C {
				props |= upropsMaskCfZwnj
			} else if u == 0x200D {
				props |= upropsMaskCfZwj
			} else if (0x180B <= u && u <= 0x180D) || u == 0x180F {
				/* Mongolian Free Variation Selectors need to be remembered
				 * because although we need to hide them like default-ignorables,
				 * they need to non-ignorable during shaping.  This is similar to
				 * what we do for joiners in Indic-like shapers, but since the
				 * FVSes are GC=Mn, we have use a separate bit to remember them.
				 * Fixes:
				 * https://github.com/harfbuzz/harfbuzz/issues/234 */
				props |= upropsMaskHidden
			} else if 0xE0020 <= u && u <= 0xE007F {
				/* TAG characters need similar treatment. Fixes:
				 * https://github.com/harfbuzz/harfbuzz/issues/463 */
				props |= upropsMaskHidden
			} else if u == 0x034F {
				/* COMBINING GRAPHEME JOINER should not be skipped; at least some times.
				 * https://github.com/harfbuzz/harfbuzz/issues/554 */
				flags |= bsfHasCGJ
				props |= upropsMaskHidden
			}
		}

		if genCat.isMark() {
			props |= upropsMaskContinuation
			props |= unicodeProp(uni.modifiedCombiningClass(u)) << 8
		}
	}

	return props, flags
}
