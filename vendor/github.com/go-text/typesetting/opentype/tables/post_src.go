// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// PostScript table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/post
type Post struct {
	version     postVersion
	italicAngle uint32
	// UnderlinePosition is the suggested distance of the top of the
	// underline from the baseline (negative values indicate below baseline).
	UnderlinePosition int16
	// Suggested values for the underline thickness.
	UnderlineThickness int16
	// IsFixedPitch indicates that the font is not proportionally spaced
	// (i.e. monospaced).
	isFixedPitch uint32
	memoryUsage  [4]uint32
	Names        PostNames `unionField:"version"`
}

type PostNames interface {
	isPostNames()
}

func (PostNames10) isPostNames() {}
func (PostNames20) isPostNames() {}
func (PostNames30) isPostNames() {}

type postVersion uint32

const (
	postVersion10 postVersion = 0x00010000
	postVersion20 postVersion = 0x00020000
	postVersion30 postVersion = 0x00030000
)

type PostNames10 struct{}

type PostNames20 struct {
	GlyphNameIndexes []uint16 `arrayCount:"FirstUint16"` // size numGlyph
	StringData       []byte   `arrayCount:"ToEnd"`
}

type PostNames30 PostNames10
