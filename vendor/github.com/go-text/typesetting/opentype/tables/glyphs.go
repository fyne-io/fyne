// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

type BitmapSubtable struct {
	FirstGlyph GlyphID //	First glyph ID of this range.
	LastGlyph  GlyphID //	Last glyph ID of this range (inclusive).
	IndexSubHeader
}

// EBLC is the Embedded Bitmap Location Table
// See - https://learn.microsoft.com/fr-fr/typography/opentype/spec/eblc
type EBLC = CBLC

// Bloc is the bitmap location table
// See - https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6bloc.html
type Bloc = CBLC
