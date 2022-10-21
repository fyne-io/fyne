package truetype

import (
	"errors"
	"fmt"

	"github.com/benoitkugler/textlayout/fonts"
)

var _ fonts.FaceRenderer = (*Font)(nil)

func (sbix tableSbix) glyphData(gid GID, xPpem, yPpem uint16) (fonts.GlyphBitmap, error) {
	st := sbix.chooseStrike(xPpem, yPpem)
	if st == nil {
		return fonts.GlyphBitmap{}, errors.New("empty 'sbix' table")
	}

	glyph := st.getGlyph(gid, 0)
	if glyph.isNil() {
		return fonts.GlyphBitmap{}, fmt.Errorf("no glyph %d in 'sbix' table for resolution (%d, %d)", gid, xPpem, yPpem)
	}

	out := fonts.GlyphBitmap{Data: glyph.data}
	var err error
	out.Width, out.Height, out.Format, err = glyph.decodeConfig()

	return out, err
}

func (colorBitmap bitmapTable) glyphData(gid GID, xPpem, yPpem uint16) (fonts.GlyphBitmap, error) {
	st := colorBitmap.chooseStrike(xPpem, yPpem)
	if st == nil || st.ppemX == 0 || st.ppemY == 0 {
		return fonts.GlyphBitmap{}, errors.New("empty bitmap table")
	}

	subtable := st.findTable(gid)
	if subtable == nil {
		return fonts.GlyphBitmap{}, fmt.Errorf("no glyph %d in bitmap table for resolution (%d, %d)", gid, xPpem, yPpem)
	}

	glyph := subtable.getImage(gid)
	if glyph == nil {
		return fonts.GlyphBitmap{}, fmt.Errorf("no glyph %d in bitmap table for resolution (%d, %d)", gid, xPpem, yPpem)
	}

	out := fonts.GlyphBitmap{
		Data:   glyph.image,
		Width:  int(glyph.metrics.width),
		Height: int(glyph.metrics.height),
	}
	switch subtable.imageFormat() {
	case 17, 18, 19: // PNG
		out.Format = fonts.PNG
	case 2, 5:
		out.Format = fonts.BlackAndWhite
	default:
		return fonts.GlyphBitmap{}, fmt.Errorf("unsupported format %d in bitmap table", subtable.imageFormat())
	}

	return out, nil
}

// look for data in 'glyf' and 'cff' tables
func (f *Font) outlineGlyphData(gid GID) (fonts.GlyphOutline, bool) {
	out, err := f.glyphDataFromCFF1(gid)
	if err == nil {
		return out, true
	}

	out, err = f.glyphDataFromGlyf(gid)
	if err == nil {
		return out, true
	}

	return fonts.GlyphOutline{}, false
}

func (f *Font) GlyphData(gid GID, xPpem, yPpem uint16) fonts.GlyphData {
	var out fonts.GlyphData

	// try every table
	out, err := f.sbix.glyphData(gid, xPpem, yPpem)
	if err == nil {
		return out
	}

	out, err = f.bitmap.glyphData(gid, xPpem, yPpem)
	if err == nil {
		return out
	}

	out_, ok := f.svg.glyphData(gid)
	if ok {
		// Spec :
		// For every SVG glyph description, there must be a corresponding TrueType,
		// CFF or CFF2 glyph description in the font.
		out_.Outline, _ = f.outlineGlyphData(gid)
		return out_
	}

	if out, ok := f.outlineGlyphData(gid); ok {
		return out
	}

	return nil
}
