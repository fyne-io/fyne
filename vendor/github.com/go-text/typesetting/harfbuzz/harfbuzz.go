// Package harfbuzz provides advanced text layout for various scripts and languages,
// with font-aware substitutions and positionning.
//
// Given a font and an input specified as runes, the package shapes this input
// and returns a slice of positionned glyphs, identified by their index in the font.
// See `Buffer` and its methods for more details.
//
// This package is a direct port of the C/C++ library.
package harfbuzz

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"strconv"

	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/opentype/api"
	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
)

// based on upstream commit 449c4296a11a8b3d3882d4a4d6705281f5ff52e5
// pending the change for Unicode 14 that are not merged yet, starting at
// 195c05df9925c7c4a4982a286ef9c416b2cde3af

// debugMode is only used in test:
//   - 0 : nothing
//   - 1 : only the main steps are printed
//   - 2 : details informations are printed
const debugMode = 0

type (
	GID = api.GID
	gID = tables.GlyphID
)

// Direction is the text direction.
// The zero value is the initial, unset, invalid direction.
type Direction uint8

const (
	LeftToRight Direction = 4 + iota // Text is set horizontally from left to right.
	RightToLeft                      // Text is set horizontally from right to left.
	TopToBottom                      // Text is set vertically from top to bottom.
	BottomToTop                      // Text is set vertically from bottom to top.
)

// Fetches the `Direction` of a script when it is
// set horizontally. All right-to-left scripts will return
// `RightToLeft`. All left-to-right scripts will return
// `LeftToRight`.  Scripts that can be written either
// horizontally or vertically will return `Invalid`.
// Unknown scripts will return `LeftToRight`.
func getHorizontalDirection(script language.Script) Direction {
	/* https://docs.google.com/spreadsheets/d/1Y90M0Ie3MUJ6UVCRDOypOtijlMDLNNyyLk36T6iMu0o */
	switch script {
	case language.Arabic, language.Hebrew, language.Syriac, language.Thaana,
		language.Cypriot, language.Kharoshthi, language.Phoenician, language.Nko, language.Lydian,
		language.Avestan, language.Imperial_Aramaic, language.Inscriptional_Pahlavi, language.Inscriptional_Parthian, language.Old_South_Arabian, language.Old_Turkic,
		language.Samaritan, language.Mandaic, language.Meroitic_Cursive, language.Meroitic_Hieroglyphs, language.Manichaean, language.Mende_Kikakui,
		language.Nabataean, language.Old_North_Arabian, language.Palmyrene, language.Psalter_Pahlavi, language.Hatran, language.Adlam, language.Hanifi_Rohingya,
		language.Old_Sogdian, language.Sogdian, language.Elymaic, language.Chorasmian, language.Yezidi:

		return RightToLeft

	/* https://github.com/harfbuzz/harfbuzz/issues/1000 */
	case language.Old_Hungarian, language.Old_Italic, language.Runic:
		return 0
	}

	return LeftToRight
}

// Tests whether a text direction is horizontal. Requires
// that the direction be valid.
func (dir Direction) isHorizontal() bool { return dir & ^Direction(1) == 4 }

// Tests whether a text direction is vertical. Requires
// that the direction be valid.
func (dir Direction) isVertical() bool { return dir & ^Direction(1) == 6 }

// Tests whether a text direction moves backward (from right to left, or from
// bottom to top). Requires that the direction be valid.
func (dir Direction) isBackward() bool { return dir & ^Direction(2) == 5 }

// Tests whether a text direction moves forward (from left to right, or from
// top to bottom). Requires that the direction be valid.
func (dir Direction) isForward() bool { return dir & ^Direction(2) == 4 }

// Reverses a text direction. Requires that the direction
// is valid.
func (dir Direction) Reverse() Direction {
	return dir ^ 1
}

// SegmentProperties holds various text properties of a `Buffer`.
type SegmentProperties struct {
	// Languages are crucial for selecting which OpenType feature to apply to the
	// buffer which can result in applying language-specific behaviour. Languages
	// are orthogonal to the scripts, and though they are related, they are
	// different concepts and should not be confused with each other.
	Language language.Language

	// Script is crucial for choosing the proper shaping behaviour for scripts that
	// require it (e.g. Arabic) and the OpenType features defined in the font
	// to be applied.
	//
	// See the package language for predefined values.
	Script language.Script

	// Direction is the text flow direction of the buffer. No shaping can happen without
	// setting direction, and it controls the visual direction for the
	// output glyphs; for RTL direction the glyphs will be reversed. Many layout
	// features depend on the proper setting of the direction, for example,
	// reversing RTL text before shaping, then shaping with LTR direction is not
	// the same as keeping the text in logical order and shaping with RTL
	// direction.
	Direction Direction
}

// ShappingOptions controls some fine tunning of the shaping
// (see the constants).
type ShappingOptions uint16

const (
	// Flag indicating that special handling of the beginning
	// of text paragraph can be applied to this buffer. Should usually
	// be set, unless you are passing to the buffer only part
	// of the text without the full context.
	Bot ShappingOptions = 1 << iota
	// Flag indicating that special handling of the end of text
	// paragraph can be applied to this buffer, similar to
	// `Bot`.
	Eot
	// Flag indication that character with Default_Ignorable
	// Unicode property should use the corresponding glyph
	// from the font, instead of hiding them (done by
	// replacing them with the space glyph and zeroing the
	// advance width.)  This flag takes precedence over
	// `RemoveDefaultIgnorables`.
	PreserveDefaultIgnorables
	// Flag indication that character with Default_Ignorable
	// Unicode property should be removed from glyph string
	// instead of hiding them (done by replacing them with the
	// space glyph and zeroing the advance width.)
	// `PreserveDefaultIgnorables` takes
	// precedence over this flag.
	RemoveDefaultIgnorables
	// Flag indicating that a dotted circle should
	// not be inserted in the rendering of incorrect
	// character sequences (such at <0905 093E>).
	DoNotinsertDottedCircle
)

// ClusterLevel allows selecting more fine-grained Cluster handling.
// It defaults to `MonotoneGraphemes`.
type ClusterLevel uint8

const (
	// Return cluster values grouped by graphemes into monotone order.
	MonotoneGraphemes ClusterLevel = iota
	//  Return cluster values grouped into monotone order.
	MonotoneCharacters
	// Don't group cluster values.
	Characters
)

func (cl ClusterLevel) String() string {
	switch cl {
	case MonotoneCharacters:
		return "MonotoneCharacters"
	case MonotoneGraphemes:
		return "MonotoneGraphemes"
	case Characters:
		return "Characters"
	default:
		return fmt.Sprintf("<unknown cluster level: %d>", cl)
	}
}

// Feature holds information about requested
// feature application. The feature will be applied with the given value to all
// glyphs which are in clusters between `start` (inclusive) and `end` (exclusive).
// Setting start to `FeatureGlobalStart` and end to `FeatureGlobalEnd`
// specifies that the feature always applies to the entire buffer.
type Feature struct {
	Tag loader.Tag
	// Value of the feature: 0 disables the feature, non-zero (usually
	// 1) enables the feature. For features implemented as lookup type 3 (like
	// 'salt') `Value` is a one-based index into the alternates.
	Value uint32
	// The cluster to Start applying this feature setting (inclusive)
	Start int
	// The cluster to End applying this feature setting (exclusive)
	End int
}

const (
	// Special setting for `Feature.Start` to apply the feature from the start
	// of the buffer.
	FeatureGlobalStart = 0
	// Special setting for `Feature.End` to apply the feature from to the end
	// of the buffer.
	FeatureGlobalEnd = maxInt
)

// ParseVariation parse the string representation of a variation
// of the form tag=value
func ParseVariation(s string) (font.Variation, error) {
	pr := parser{data: []byte(s)}
	return pr.parseOneVariation()
}

type parser struct {
	data []byte
	pos  int
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\f' || c == '\n' || c == '\r' || c == '\t' || c == '\v'
}

func (p *parser) skipSpaces() {
	for p.pos < len(p.data) && isSpace(p.data[p.pos]) {
		p.pos++
	}
}

// return true if `c` was found
func (p *parser) parseChar(c byte) bool {
	p.skipSpaces()

	if p.pos == len(p.data) || p.data[p.pos] != c {
		return false
	}
	p.pos++
	return true
}

func (p *parser) parseUint32() (uint32, bool) {
	start := p.pos
	// go to the next space
	for p.pos < len(p.data) && isAlnum(p.data[p.pos]) {
		p.pos++
	}
	out, err := strconv.Atoi(string(p.data[start:p.pos]))
	return uint32(out), err == nil
}

func (p *parser) parseBool() (uint32, bool) {
	p.skipSpaces()

	startPos := p.pos
	for p.pos < len(p.data) && isAlpha(p.data[p.pos]) {
		p.pos++
	}
	data := string(p.data[startPos:p.pos])

	/* CSS allows on/off as aliases 1/0. */
	if data == "on" {
		return 1, true
	} else if data == "off" {
		return 0, true
	} else {
		return 0, false
	}
}

func (p *parser) parseTag() (loader.Tag, error) {
	p.skipSpaces()

	var quote byte

	if p.pos < len(p.data) && (p.data[p.pos] == '\'' || p.data[p.pos] == '"') {
		quote = p.data[p.pos]
		p.pos++
	}

	start := p.pos
	for p.pos < len(p.data) && (isAlnum(p.data[p.pos]) || p.data[p.pos] == '_') {
		p.pos++
	}

	if p.pos == start || p.pos > start+4 {
		return 0, errors.New("invalid tag length")
	}

	// padd with space if necessary, since MustNewTag requires 4 bytes
	tagBytes := [4]byte{' ', ' ', ' ', ' '}
	copy(tagBytes[:], p.data[start:p.pos])
	tag := loader.MustNewTag(string(tagBytes[:]))

	if quote != 0 {
		/* CSS expects exactly four bytes.  And we only allow quotations for
		 * CSS compatibility.  So, enforce the length. */
		if p.pos != start+4 {
			return 0, errors.New("tag must have 4 bytes")
		}
		if p.pos == len(p.data) || p.data[p.pos] != quote {
			return 0, errors.New("tag is missing end quote")
		}
		p.pos++
	}

	return tag, nil
}

func (p *parser) parseVariationValue() (float32, error) {
	p.parseChar('=') // Optional.
	start := p.pos
	// go to the next space
	for p.pos < len(p.data) && !isSpace(p.data[p.pos]) {
		p.pos++
	}
	v, err := strconv.ParseFloat(string(p.data[start:p.pos]), 32)
	return float32(v), err
}

func (p *parser) parseOneVariation() (vari font.Variation, err error) {
	vari.Tag, err = p.parseTag()
	if err != nil {
		return
	}
	vari.Value, err = p.parseVariationValue()
	if err != nil {
		return
	}
	p.skipSpaces()
	return
}

func (p *parser) parseFeatureIndices() (start, end int, err error) {
	p.skipSpaces()

	start, end = FeatureGlobalStart, FeatureGlobalEnd

	if !p.parseChar('[') {
		return start, end, nil
	}

	startU, hasStart := p.parseUint32()
	start = int(startU)

	if p.parseChar(':') || p.parseChar(';') {
		if endU, ok := p.parseUint32(); ok {
			end = int(endU)
		}
	} else {
		if hasStart {
			end = start + 1
		}
	}

	if !p.parseChar(']') {
		return 0, 0, errors.New("expecting closing bracked after feature indices")
	}

	return start, end, nil
}

// return true if a value was specified
func (p *parser) parseFeatureValuePostfix() (uint32, bool) {
	/* CSS doesn't use equal-sign between tag and value.
	 * If there was an equal-sign, then there *must* be a value.
	 * A value without an equal-sign is ok, but not required. */
	p.parseChar('=')

	val, hadValue := p.parseUint32()
	if !hadValue {
		val, hadValue = p.parseBool()
	}
	return val, hadValue
}

func (p *parser) parseFeatureValuePrefix() uint32 {
	if p.parseChar('-') {
		return 0
	} else {
		_ = p.parseChar('+')
		return 1
	}
}

func (p *parser) parseOneFeature() (feature Feature, err error) {
	feature.Value = p.parseFeatureValuePrefix()
	feature.Tag, err = p.parseTag()
	if err != nil {
		return feature, err
	}
	feature.Start, feature.End, err = p.parseFeatureIndices()
	if err != nil {
		return feature, err
	}
	if val, ok := p.parseFeatureValuePostfix(); ok {
		feature.Value = val
	}
	p.skipSpaces()
	return feature, nil
}

// ParseFeature parses one feature string (usually coming from a comma-separated list of font features).
//
//	Features can be enabled or disabled, either globally or limited to
//	specific character ranges.  The format for specifying feature settings
//	follows.  All valid CSS font-feature-settings values other than 'normal'
//	and the global values are also accepted, though not documented below.
//	CSS string escapes are not supported.
//
//	The range indices refer to the positions between Unicode characters,
//	unless the --utf8-clusters is provided, in which case range indices
//	refer to UTF-8 byte indices. The position before the first character
//	is always 0.
//
//	The format is Python-esque.  Here is how it all works:
//
//	  Syntax:       Value:    Start:    End:
//
//	Setting value:
//	  "kern"        1         0         ∞         // Turn feature on
//	  "+kern"       1         0         ∞         // Turn feature on
//	  "-kern"       0         0         ∞         // Turn feature off
//	  "kern=0"      0         0         ∞         // Turn feature off
//	  "kern=1"      1         0         ∞         // Turn feature on
//	  "aalt=2"      2         0         ∞         // Choose 2nd alternate
//
//	Setting index:
//	  "kern[]"      1         0         ∞         // Turn feature on
//	  "kern[:]"     1         0         ∞         // Turn feature on
//	  "kern[5:]"    1         5         ∞         // Turn feature on, partial
//	  "kern[:5]"    1         0         5         // Turn feature on, partial
//	  "kern[3:5]"   1         3         5         // Turn feature on, range
//	  "kern[3]"     1         3         3+1       // Turn feature on, single char
//
//	Mixing it all:
//
//	  "aalt[3:5]=2" 2         3         5         // Turn 2nd alternate on for range
func ParseFeature(feature string) (Feature, error) {
	pr := parser{data: []byte(feature)}
	return pr.parseOneFeature()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func isAlpha(c byte) bool { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isAlnum(c byte) bool { return isAlpha(c) || (c >= '0' && c <= '9') }
func toUpper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 'A'
	}
	return c
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c - 'A' + 'a'
	}
	return c
}

const maxInt = int(^uint(0) >> 1)

// bitStorage returns the number of bits needed to store the number.
func bitStorage(v uint32) int { return 32 - bits.LeadingZeros32(v) }

func roundf(f float32) Position {
	return Position(math.Round(float64(f)))
}
