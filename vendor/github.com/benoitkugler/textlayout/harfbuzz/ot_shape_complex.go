package harfbuzz

import (
	tt "github.com/benoitkugler/textlayout/fonts/truetype"
	"github.com/benoitkugler/textlayout/language"
)

type zeroWidthMarks uint8

const (
	zeroWidthMarksNone zeroWidthMarks = iota
	zeroWidthMarksByGdefEarly
	zeroWidthMarksByGdefLate
)

// implements the specialisation for a script
type otComplexShaper interface {
	marksBehavior() (zwm zeroWidthMarks, fallbackPosition bool)
	normalizationPreference() normalizationMode
	// If not 0, then must match found GPOS script tag for
	// GPOS to be applied. Otherwise, fallback positioning will be used.
	gposTag() tt.Tag

	// collectFeatures is alled during shape_plan().
	// Shapers should use plan.map to add their features and callbacks.
	collectFeatures(plan *otShapePlanner)

	// overrideFeatures is called during shape_plan().
	// Shapers should use plan.map to override features and add callbacks after
	// common features are added.
	overrideFeatures(plan *otShapePlanner)

	// dataCreate is called at the end of shape_plan().
	dataCreate(plan *otShapePlan)

	// called during shape(), shapers can use to modify text before shaping starts.
	preprocessText(plan *otShapePlan, buffer *Buffer, font *Font)

	// called during shape()'s normalization: may use decompose_unicode as fallback
	decompose(c *otNormalizeContext, ab rune) (a, b rune, ok bool)

	// called during shape()'s normalization: may use compose_unicode as fallback
	compose(c *otNormalizeContext, a, b rune) (ab rune, ok bool)

	// called during shape(), shapers should use map to get feature masks and set on buffer.
	// Shapers may NOT modify characters.
	setupMasks(plan *otShapePlan, buffer *Buffer, font *Font)

	// called during shape(), shapers can use to modify ordering of combining marks.
	reorderMarks(plan *otShapePlan, buffer *Buffer, start, end int)

	// called during shape(), shapers can use to modify glyphs after shaping ends.
	postprocessGlyphs(plan *otShapePlan, buffer *Buffer, font *Font)
}

/*
 * For lack of a better place, put Zawgyi script hack here.
 * https://github.com/harfbuzz/harfbuzz/issues/1162
 */
var scriptMyanmarZawgyi = language.Script(tt.NewTag('Q', 'a', 'a', 'g'))

func (planner *otShapePlanner) categorizeComplex() otComplexShaper {
	switch planner.props.Script {
	case language.Arabic, language.Syriac:
		/* For Arabic script, use the Arabic shaper even if no OT script tag was found.
		 * This is because we do fallback shaping for Arabic script (and not others).
		 * But note that Arabic shaping is applicable only to horizontal layout; for
		 * vertical text, just use the generic shaper instead. */
		if (planner.map_.chosenScript[0] != tagDefaultScript ||
			planner.props.Script == language.Arabic) &&
			planner.props.Direction.isHorizontal() {
			return &complexShaperArabic{}
		}
		return complexShaperDefault{}
	case language.Thai, language.Lao:
		return complexShaperThai{}
	case language.Hangul:
		return &complexShaperHangul{}
	case language.Hebrew:
		return complexShaperHebrew{}
	case language.Bengali, language.Devanagari, language.Gujarati, language.Gurmukhi, language.Kannada,
		language.Malayalam, language.Oriya, language.Tamil, language.Telugu, language.Sinhala:
		/* If the designer designed the font for the 'DFLT' script,
		 * (or we ended up arbitrarily pick 'latn'), use the default shaper.
		 * Otherwise, use the specific shaper.
		 *
		 * If it's indy3 tag, send to USE. */
		if planner.map_.chosenScript[0] == tt.NewTag('D', 'F', 'L', 'T') ||
			planner.map_.chosenScript[0] == tt.NewTag('l', 'a', 't', 'n') {
			return complexShaperDefault{}
		} else if (planner.map_.chosenScript[0] & 0x000000FF) == '3' {
			return &complexShaperUSE{}
		}
		return &complexShaperIndic{}
	case language.Khmer:
		return &complexShaperKhmer{}
	case language.Myanmar:
		/* If the designer designed the font for the 'DFLT' script,
		 * (or we ended up arbitrarily pick 'latn'), use the default shaper.
		 * Otherwise, use the specific shaper.
		 *
		 * If designer designed for 'mymr' tag, also send to default
		 * shaper.  That's tag used from before Myanmar shaping spec
		 * was developed.  The shaping spec uses 'mym2' tag. */
		if planner.map_.chosenScript[0] == tt.NewTag('D', 'F', 'L', 'T') ||
			planner.map_.chosenScript[0] == tt.NewTag('l', 'a', 't', 'n') ||
			planner.map_.chosenScript[0] == tt.NewTag('m', 'y', 'm', 'r') {
			return complexShaperDefault{}
		}
		return complexShaperMyanmar{}

	case scriptMyanmarZawgyi:
		/* Ugly Zawgyi encoding.
		 * Disable all auto processing.
		 * https://github.com/harfbuzz/harfbuzz/issues/1162 */
		return complexShaperDefault{dumb: true, disableNorm: true}
	case language.Tibetan,
		language.Mongolian,
		language.Buhid, language.Hanunoo, language.Tagalog, language.Tagbanwa,
		language.Limbu, language.Tai_Le,
		language.Buginese, language.Kharoshthi, language.Syloti_Nagri, language.Tifinagh,
		language.Balinese, language.Nko, language.Phags_Pa, language.Cham, language.Kayah_Li,
		language.Lepcha, language.Rejang, language.Saurashtra, language.Sundanese,
		language.Egyptian_Hieroglyphs, language.Javanese, language.Kaithi,
		language.Meetei_Mayek, language.Tai_Tham, language.Tai_Viet, language.Batak,
		language.Brahmi, language.Mandaic, language.Chakma, language.Miao, language.Sharada,
		language.Takri, language.Duployan, language.Grantha, language.Khojki, language.Khudawadi,
		language.Mahajani, language.Manichaean, language.Modi, language.Pahawh_Hmong,
		language.Psalter_Pahlavi, language.Siddham, language.Tirhuta, language.Ahom, language.Multani,
		language.Adlam, language.Bhaiksuki, language.Marchen, language.Newa, language.Masaram_Gondi,
		language.Soyombo, language.Zanabazar_Square, language.Dogra, language.Gunjala_Gondi,
		language.Hanifi_Rohingya, language.Makasar, language.Medefaidrin, language.Old_Sogdian,
		language.Sogdian, language.Elymaic, language.Nandinagari, language.Nyiakeng_Puachue_Hmong,
		language.Wancho,
		language.Chorasmian, language.Dives_Akuru, language.Khitan_Small_Script, language.Yezidi:

		/* If the designer designed the font for the 'DFLT' script,
		 * (or we ended up arbitrarily pick 'latn'), use the default shaper.
		 * Otherwise, use the specific shaper.
		 * Note that for some simple scripts, there may not be *any*
		 * GSUB/GPOS needed, so there may be no scripts found! */
		if planner.map_.chosenScript[0] == tt.NewTag('D', 'F', 'L', 'T') ||
			planner.map_.chosenScript[0] == tt.NewTag('l', 'a', 't', 'n') {
			return complexShaperDefault{}
		}
		return &complexShaperUSE{}
	default:
		return complexShaperDefault{}
	}
}

// zero byte struct providing no-ops, used to reduced boilerplate
type complexShaperNil struct{}

func (complexShaperNil) gposTag() tt.Tag { return 0 }

func (complexShaperNil) collectFeatures(plan *otShapePlanner)  {}
func (complexShaperNil) overrideFeatures(plan *otShapePlanner) {}
func (complexShaperNil) dataCreate(plan *otShapePlan)          {}
func (complexShaperNil) decompose(_ *otNormalizeContext, ab rune) (a, b rune, ok bool) {
	return uni.decompose(ab)
}

func (complexShaperNil) compose(_ *otNormalizeContext, a, b rune) (ab rune, ok bool) {
	return uni.compose(a, b)
}
func (complexShaperNil) preprocessText(*otShapePlan, *Buffer, *Font) {}
func (complexShaperNil) postprocessGlyphs(*otShapePlan, *Buffer, *Font) {
}
func (complexShaperNil) setupMasks(*otShapePlan, *Buffer, *Font)      {}
func (complexShaperNil) reorderMarks(*otShapePlan, *Buffer, int, int) {}

type complexShaperDefault struct {
	complexShaperNil

	/* if true, no mark advance zeroing / fallback positioning.
	 * Dumbest shaper ever, basically. */
	dumb        bool
	disableNorm bool
}

func (cs complexShaperDefault) marksBehavior() (zeroWidthMarks, bool) {
	if cs.dumb {
		return zeroWidthMarksNone, false
	}
	return zeroWidthMarksByGdefLate, true
}

func (cs complexShaperDefault) normalizationPreference() normalizationMode {
	if cs.disableNorm {
		return nmNone
	}
	return nmDefault
}

func syllabicInsertDottedCircles(font *Font, buffer *Buffer, brokenSyllableType,
	dottedcircleCategory uint8, rephaCategory, dottedCirclePosition int) {
	if (buffer.Flags & DoNotinsertDottedCircle) != 0 {
		return
	}

	hasBrokenSyllables := false
	info := buffer.Info
	for _, inf := range info {
		if (inf.syllable & 0x0F) == brokenSyllableType {
			hasBrokenSyllables = true
			break
		}
	}
	if !hasBrokenSyllables {
		return
	}

	dottedcircleGlyph, ok := font.face.NominalGlyph(0x25CC)
	if !ok {
		return
	}

	dottedcircle := GlyphInfo{
		Glyph:           dottedcircleGlyph,
		complexCategory: dottedcircleCategory,
	}

	if dottedCirclePosition != -1 {
		dottedcircle.complexAux = uint8(dottedCirclePosition)
	}

	buffer.clearOutput()

	buffer.idx = 0
	var lastSyllable uint8
	for buffer.idx < len(buffer.Info) {
		syllable := buffer.cur(0).syllable
		if lastSyllable != syllable && (syllable&0x0F) == brokenSyllableType {
			lastSyllable = syllable

			ginfo := dottedcircle
			ginfo.Cluster = buffer.cur(0).Cluster
			ginfo.Mask = buffer.cur(0).Mask
			ginfo.syllable = buffer.cur(0).syllable

			/* Insert dottedcircle after possible Repha. */
			if rephaCategory != -1 {
				for buffer.idx < len(buffer.Info) &&
					lastSyllable == buffer.cur(0).syllable &&
					buffer.cur(0).complexCategory == uint8(rephaCategory) {
					buffer.nextGlyph()
				}
			}
			buffer.outInfo = append(buffer.outInfo, ginfo)
		} else {
			buffer.nextGlyph()
		}
	}
	buffer.swapBuffers()
}
