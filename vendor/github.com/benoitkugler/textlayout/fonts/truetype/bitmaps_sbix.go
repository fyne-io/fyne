package truetype

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"

	"github.com/benoitkugler/textlayout/fonts"
	"golang.org/x/image/tiff"
)

var (
	// TagPNG identifies bitmap glyph with png format
	TagPNG = MustNewTag("png ")
	// TagTIFF identifies bitmap glyph with tiff format
	TagTIFF = MustNewTag("tiff")
	// TagJPG identifies bitmap glyph with jpg format
	TagJPG = MustNewTag("jpg ")
)

// ---------------------------------------- sbix ----------------------------------------

type tableSbix struct {
	strikes      []bitmapStrike
	drawOutlines bool
}

// return nil only if the table is empty
func (t tableSbix) chooseStrike(xPpem, yPpem uint16) *bitmapStrike {
	if len(t.strikes) == 0 {
		return nil
	}

	request := maxu16(xPpem, yPpem)
	if request == 0 {
		request = math.MaxUint16 // choose largest strike
	}

	/* TODO Add DPI sensitivity as well? */
	var (
		bestIndex = 0
		bestPpem  = t.strikes[0].ppem
	)
	for i, s := range t.strikes {
		ppem := s.ppem
		if request <= ppem && ppem < bestPpem || request > bestPpem && ppem > bestPpem {
			bestIndex = i
			bestPpem = ppem
		}
	}
	return &t.strikes[bestIndex]
}

func parseTableSbix(data []byte, numGlyphs int) (out tableSbix, err error) {
	if len(data) < 8 {
		return out, errors.New("invalid 'sbix' table (EOF)")
	}
	flag := binary.BigEndian.Uint16(data[2:])
	numStrikes := int(binary.BigEndian.Uint32(data[4:]))

	out.drawOutlines = flag&0x02 != 0

	if len(data) < 8+8*numStrikes {
		return out, errors.New("invalid 'sbix' table (EOF)")
	}
	out.strikes = make([]bitmapStrike, numStrikes)
	for i := range out.strikes {
		offset := binary.BigEndian.Uint32(data[8+4*i:])
		out.strikes[i], err = parseBitmapStrike(data, offset, numGlyphs)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

type bitmapStrike struct {
	// length numGlyph; items may be empty (see isNil)
	glyphs    []bitmapGlyphData
	ppem, ppi uint16
}

func mulDiv(a, b, c uint16) uint16 {
	return uint16(uint32(a) * uint32(b) / uint32(c))
}

func (b *bitmapStrike) sizeMetrics(hori *TableHVhea, avgWidth, upem uint16) (out fonts.BitmapSize) {
	out.XPpem, out.YPpem = b.ppem, b.ppem
	out.Height = mulDiv(uint16(hori.Ascent-hori.Descent+hori.LineGap), b.ppem, upem)

	inferBitmapWidth(&out, avgWidth, upem)

	return out
}

// may return a zero value
func (b *bitmapStrike) getGlyph(glyph GID, recursionLevel int) bitmapGlyphData {
	const maxRecursionLevel = 8

	if int(glyph) >= len(b.glyphs) {
		return bitmapGlyphData{}
	}
	out := b.glyphs[glyph]
	if out.graphicType == MustNewTag("dupe") {
		if len(out.data) < 2 || recursionLevel > maxRecursionLevel {
			return bitmapGlyphData{}
		}
		glyph = GID(binary.BigEndian.Uint16(out.data))
		return b.getGlyph(glyph, recursionLevel+1)
	}
	return out
}

func parseBitmapStrike(data []byte, offset uint32, numGlyphs int) (out bitmapStrike, err error) {
	if len(data) < int(offset)+4+4*(numGlyphs+1) {
		return out, errors.New("invalud sbix bitmap strike (EOF)")
	}
	data = data[offset:]
	out.ppem = binary.BigEndian.Uint16(data)
	out.ppi = binary.BigEndian.Uint16(data[2:])

	offsets, _ := parseTableLoca(data[4:], numGlyphs, true)
	out.glyphs = make([]bitmapGlyphData, numGlyphs)
	for i := range out.glyphs {
		if offsets[i] == offsets[i+1] { // no data
			continue
		}

		out.glyphs[i], err = parseBitmapGlyphData(data, offsets[i], offsets[i+1])
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

type bitmapGlyphData struct {
	data                         []byte
	originOffsetX, originOffsetY int16 // in font units
	graphicType                  Tag
}

func (b bitmapGlyphData) isNil() bool { return b.graphicType == 0 }

// decodeConfig parse the data to find the width and height
func (b bitmapGlyphData) decodeConfig() (width, height int, format fonts.BitmapFormat, err error) {
	var config image.Config
	switch b.graphicType {
	case TagPNG:
		format = fonts.PNG
		config, err = png.DecodeConfig(bytes.NewReader(b.data))
	case TagTIFF:
		format = fonts.TIFF
		config, err = tiff.DecodeConfig(bytes.NewReader(b.data))
	case TagJPG:
		format = fonts.JPG
		config, err = jpeg.DecodeConfig(bytes.NewReader(b.data))
	default:
		err = fmt.Errorf("unsupported graphic type in sbix table: %s", b.graphicType)
	}
	if err != nil {
		return 0, 0, 0, err
	}
	return config.Width, config.Height, format, nil
}

// return the extents computed from the data
// should only be called on valid, non nil glyph data
func (b bitmapGlyphData) glyphExtents() (out fonts.GlyphExtents, ok bool) {
	width, height, _, err := b.decodeConfig()
	if err != nil {
		return out, false
	}
	out.XBearing = float32(b.originOffsetX)
	out.YBearing = float32(height) + float32(b.originOffsetY)
	out.Width = float32(width)
	out.Height = -float32(height)
	return out, true
}

func parseBitmapGlyphData(data []byte, offsetStart, offsetNext uint32) (out bitmapGlyphData, err error) {
	if len(data) < int(offsetStart)+8 || offsetStart+8 > offsetNext {
		return out, errors.New("invalid 'sbix' bitmap glyph data (EOF)")
	}
	data = data[offsetStart:]
	out.originOffsetX = int16(binary.BigEndian.Uint16(data))
	out.originOffsetY = int16(binary.BigEndian.Uint16(data[2:]))
	out.graphicType = Tag(binary.BigEndian.Uint32(data[4:]))
	out.data = data[8 : offsetNext-offsetStart]

	if out.graphicType == 0 {
		return out, errors.New("invalid 'sbix' zero bitmap type")
	}
	return out, nil
}
