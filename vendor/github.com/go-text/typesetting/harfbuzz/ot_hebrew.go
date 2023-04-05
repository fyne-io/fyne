package harfbuzz

import (
	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
)

// ported from harfbuzz/src/hb-ot-shape-complex-hebrew.cc Copyright © 2010,2012  Google, Inc.  Behdad Esfahbod

var _ otComplexShaper = complexShaperHebrew{}

type complexShaperHebrew struct {
	complexShaperNil
}

/* Hebrew presentation-form shaping.
* https://bugzilla.mozilla.org/show_bug.cgi?id=728866
* Hebrew presentation forms with dagesh, for characters U+05D0..05EA;
* Note that some letters do not have a dagesh presForm encoded. */
var sDageshForms = [0x05EA - 0x05D0 + 1]rune{
	0xFB30, /* ALEF */
	0xFB31, /* BET */
	0xFB32, /* GIMEL */
	0xFB33, /* DALET */
	0xFB34, /* HE */
	0xFB35, /* VAV */
	0xFB36, /* ZAYIN */
	0x0000, /* HET */
	0xFB38, /* TET */
	0xFB39, /* YOD */
	0xFB3A, /* FINAL KAF */
	0xFB3B, /* KAF */
	0xFB3C, /* LAMED */
	0x0000, /* FINAL MEM */
	0xFB3E, /* MEM */
	0x0000, /* FINAL NUN */
	0xFB40, /* NUN */
	0xFB41, /* SAMEKH */
	0x0000, /* AYIN */
	0xFB43, /* FINAL PE */
	0xFB44, /* PE */
	0x0000, /* FINAL TSADI */
	0xFB46, /* TSADI */
	0xFB47, /* QOF */
	0xFB48, /* RESH */
	0xFB49, /* SHIN */
	0xFB4A, /* TAV */
}

func (complexShaperHebrew) compose(c *otNormalizeContext, a, b rune) (rune, bool) {
	ab, found := uni.compose(a, b)

	if !found && !c.plan.hasGposMark {
		/* Special-case Hebrew presentation forms that are excluded from
		* standard normalization, but wanted for old fonts. */
		switch b {
		case 0x05B4: /* HIRIQ */
			if a == 0x05D9 { /* YOD */
				return 0xFB1D, true
			}
		case 0x05B7: /* patah */
			if a == 0x05F2 { /* YIDDISH YOD YOD */
				return 0xFB1F, true
			} else if a == 0x05D0 { /* ALEF */
				return 0xFB2E, true
			}
		case 0x05B8: /* QAMATS */
			if a == 0x05D0 { /* ALEF */
				return 0xFB2F, true
			}
		case 0x05B9: /* HOLAM */
			if a == 0x05D5 { /* VAV */
				return 0xFB4B, true
			}
		case 0x05BC: /* DAGESH */
			if a >= 0x05D0 && a <= 0x05EA {
				ab = sDageshForms[a-0x05D0]
				return ab, ab != 0
			} else if a == 0xFB2A { /* SHIN WITH SHIN DOT */
				return 0xFB2C, true
			} else if a == 0xFB2B { /* SHIN WITH SIN DOT */
				return 0xFB2D, true
			}
		case 0x05BF: /* RAFE */
			switch a {
			case 0x05D1: /* BET */
				return 0xFB4C, true
			case 0x05DB: /* KAF */
				return 0xFB4D, true
			case 0x05E4: /* PE */
				return 0xFB4E, true
			}
		case 0x05C1: /* SHIN DOT */
			if a == 0x05E9 { /* SHIN */
				return 0xFB2A, true
			} else if a == 0xFB49 { /* SHIN WITH DAGESH */
				return 0xFB2C, true
			}
		case 0x05C2: /* SIN DOT */
			if a == 0x05E9 { /* SHIN */
				return 0xFB2B, true
			} else if a == 0xFB49 { /* SHIN WITH DAGESH */
				return 0xFB2D, true
			}
		}
	}

	return ab, found
}

func (complexShaperHebrew) marksBehavior() (zeroWidthMarks, bool) {
	return zeroWidthMarksByGdefLate, true
}

func (complexShaperHebrew) normalizationPreference() normalizationMode {
	return nmDefault
}

func (complexShaperHebrew) gposTag() tables.Tag {
	// https://github.com/harfbuzz/harfbuzz/issues/347#issuecomment-267838368
	return loader.NewTag('h', 'e', 'b', 'r')
}
