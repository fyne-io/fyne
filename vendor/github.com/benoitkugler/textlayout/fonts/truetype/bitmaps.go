package truetype

import (
	"github.com/benoitkugler/textlayout/fonts"
)

func (t bitmapTable) availableSizes(avgWidth, upem uint16) []fonts.BitmapSize {
	out := make([]fonts.BitmapSize, 0, len(t))
	for _, size := range t {
		v := size.sizeMetrics(avgWidth, upem)
		/* only use strikes with valid PPEM values */
		if v.XPpem == 0 || v.YPpem == 0 {
			continue
		}
		out = append(out, v)
	}
	return out
}

func (t tableSbix) availableSizes(horizontal *TableHVhea, avgWidth, upem uint16) []fonts.BitmapSize {
	out := make([]fonts.BitmapSize, 0, len(t.strikes))
	for _, size := range t.strikes {
		v := size.sizeMetrics(horizontal, avgWidth, upem)
		/* only use strikes with valid PPEM values */
		if v.XPpem == 0 || v.YPpem == 0 {
			continue
		}
		out = append(out, v)
	}
	return out
}

func inferBitmapWidth(size *fonts.BitmapSize, avgWidth, upem uint16) {
	size.Width = uint16((uint32(avgWidth)*uint32(size.XPpem) + uint32(upem/2)) / uint32(upem))
}

// return nil if no table is valid (or present)
func (pr *FontParser) selectBitmapTable() bitmapTable {
	color, err := pr.colorBitmapTable()
	if err == nil {
		return color
	}

	gray, err := pr.grayBitmapTable()
	if err == nil {
		return gray
	}

	apple, err := pr.appleBitmapTable()
	if err == nil {
		return apple
	}

	return nil
}

// LoadBitmaps checks for the various bitmaps table and returns
// the first valid
func (font *Font) LoadBitmaps() []fonts.BitmapSize {
	upem := font.Head.UnitsPerEm

	avgWidth := font.OS2.XAvgCharWidth

	if upem == 0 || font.OS2.Version == 0xFFFF {
		avgWidth = 1
		upem = 1
	}

	// adapted from freetype tt_face_load_sbit
	if font.bitmap != nil {
		return font.bitmap.availableSizes(avgWidth, upem)
	}

	if hori := font.hhea; hori != nil {
		return font.sbix.availableSizes(hori, avgWidth, upem)
	}

	return nil
}
