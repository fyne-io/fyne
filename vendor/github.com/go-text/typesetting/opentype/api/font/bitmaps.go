// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package font

import (
	"errors"
	"fmt"
	"math"

	"github.com/go-text/typesetting/opentype/api"
	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
)

// sbix

type sbix []tables.Strike

func newSbix(table tables.Sbix) sbix { return table.Strikes }

// chooseStrike selects the best match for the given resolution.
// It returns nil only if the table is empty
func (sb sbix) chooseStrike(xPpem, yPpem uint16) *tables.Strike {
	if len(sb) == 0 {
		return nil
	}

	request := maxu16(xPpem, yPpem)
	if request == 0 {
		request = math.MaxUint16 // choose largest strike
	}

	var (
		bestIndex = 0
		bestPpem  = sb[0].Ppem
	)
	for i, s := range sb {
		ppem := s.Ppem
		if request <= ppem && ppem < bestPpem || request > bestPpem && ppem > bestPpem {
			bestIndex = i
			bestPpem = ppem
		}
	}
	return &sb[bestIndex]
}

func (sb sbix) availableSizes(horizontal *tables.Hhea, avgWidth, upem uint16) []api.BitmapSize {
	out := make([]api.BitmapSize, 0, len(sb))
	for _, size := range sb {
		v := strikeSizeMetrics(size, horizontal, avgWidth, upem)
		// only use strikes with valid PPEM values
		if v.XPpem == 0 || v.YPpem == 0 {
			continue
		}
		out = append(out, v)
	}
	return out
}

func strikeSizeMetrics(b tables.Strike, hori *tables.Hhea, avgWidth, upem uint16) (out api.BitmapSize) {
	out.XPpem, out.YPpem = b.Ppem, b.Ppem
	out.Height = mulDiv(uint16(hori.Ascender-hori.Descender+hori.LineGap), b.Ppem, upem)

	inferBitmapWidth(&out, avgWidth, upem)

	return out
}

// ---------------------------- bitmap ----------------------------

func loadBitmap(ld *loader.Loader, tagLoc, tagData loader.Tag) (bitmap, error) {
	raw, err := ld.RawTable(tagLoc)
	if err != nil {
		return nil, err
	}
	loc, _, err := tables.ParseCBLC(raw)
	if err != nil {
		return nil, err
	}
	imageTable, err := ld.RawTable(tagData)
	if err != nil {
		return nil, err
	}
	return newBitmap(loc, imageTable)
}

// CBLC/CBDT or EBLC/EBDT or BLOC/BDAT
type bitmap []bitmapStrike

func newBitmap(table tables.EBLC, imageTable []byte) (bitmap, error) {
	out := make(bitmap, len(table.BitmapSizes))
	for i, strike := range table.BitmapSizes {
		subtables := table.IndexSubTables[i]
		out[i] = bitmapStrike{
			subTables: make([]bitmapSubtable, len(subtables)),
			hori:      strike.Hori,
			vert:      strike.Vert,
			ppemX:     uint16(strike.PpemX),
			ppemY:     uint16(strike.PpemY),
		}
		for j, subtable := range subtables {
			var err error
			out[i].subTables[j], err = newBitmapSubtable(subtable, imageTable)
			if err != nil {
				return nil, err
			}
		}
	}
	return out, nil
}

func (t bitmap) availableSizes(avgWidth, upem uint16) []api.BitmapSize {
	out := make([]api.BitmapSize, 0, len(t))
	for _, size := range t {
		v := size.sizeMetrics(avgWidth, upem)
		// only use strikes with valid PPEM values
		if v.XPpem == 0 || v.YPpem == 0 {
			continue
		}
		out = append(out, v)
	}
	return out
}

type bitmapStrike struct {
	subTables    []bitmapSubtable
	hori, vert   tables.SbitLineMetrics
	ppemX, ppemY uint16
}

// chooseStrike selects the best match for the given resolution.
// It returns nil only if the table is empty
func (bt bitmap) chooseStrike(xPpem, yPpem uint16) *bitmapStrike {
	if len(bt) == 0 {
		return nil
	}
	request := maxu16(xPpem, yPpem)
	if request == 0 {
		request = math.MaxUint16 // choose largest strike
	}
	var (
		bestIndex = 0
		bestPpem  = maxu16(bt[0].ppemX, bt[0].ppemY)
	)
	for i, s := range bt {
		ppem := maxu16(s.ppemX, s.ppemY)
		if request <= ppem && ppem < bestPpem || request > bestPpem && ppem > bestPpem {
			bestIndex = i
			bestPpem = ppem
		}
	}
	return &bt[bestIndex]
}

func (b *bitmapStrike) sizeMetrics(avgWidth, upem uint16) (out api.BitmapSize) {
	out.XPpem, out.YPpem = b.ppemX, b.ppemY
	ascender := int16(b.hori.Ascender)
	descender := int16(b.hori.Descender)

	maxBeforeBl := b.hori.MaxBeforeBL
	minAfterBl := b.hori.MinAfterBL

	/* Due to fuzzy wording in the EBLC documentation, we find both */
	/* positive and negative values for `descender'.  Additionally, */
	/* many fonts have both `ascender' and `descender' set to zero  */
	/* (which is definitely wrong).  MS Windows simply ignores all  */
	/* those values...  For these reasons we apply some heuristics  */
	/* to get a reasonable, non-zero value for the height.          */

	if descender > 0 {
		if minAfterBl < 0 {
			descender = -descender
		}
	} else if descender == 0 {
		if ascender == 0 {
			/* sanitize buggy ascender and descender values */
			if maxBeforeBl != 0 || minAfterBl != 0 {
				ascender = int16(maxBeforeBl)
				descender = int16(minAfterBl)
			} else {
				ascender = int16(out.YPpem)
				descender = 0
			}
		}
	}

	if h := ascender - descender; h > 0 {
		out.Height = uint16(h)
	} else {
		out.Height = out.YPpem
	}

	inferBitmapWidth(&out, avgWidth, upem)

	return out
}

func inferBitmapWidth(size *api.BitmapSize, avgWidth, upem uint16) {
	size.Width = uint16((uint32(avgWidth)*uint32(size.XPpem) + uint32(upem/2)) / uint32(upem))
}

// return nil when not found
func (b *bitmapStrike) findTable(glyph gID) *bitmapSubtable {
	for i, subtable := range b.subTables {
		if subtable.first <= glyph && glyph <= subtable.last {
			return &b.subTables[i]
		}
	}
	return nil
}

type bitmapSubtable struct {
	first       gID // First glyph ID of this range.
	last        gID // Last glyph ID of this range (inclusive).
	imageFormat uint16
	index       bitmapIndex
}

func newBitmapSubtable(header tables.BitmapSubtable, dataTable []byte) (bitmapSubtable, error) {
	out := bitmapSubtable{
		first:       header.FirstGlyph,
		last:        header.LastGlyph,
		imageFormat: header.ImageFormat,
	}
	if L, E := len(dataTable), int(header.ImageDataOffset); L < E {
		return bitmapSubtable{}, errors.New("invalid bitmap table (EOF)")
	}
	imageData := dataTable[header.ImageDataOffset:]

	var err error
	switch index := header.IndexData.(type) {
	case tables.IndexData1:
		out.index, err = parseIndexSubTable1(header, index, imageData)
	case tables.IndexData2:
		out.index, err = parseIndexSubTable2(header, index, imageData)
	case tables.IndexData3:
		out.index, err = parseIndexSubTable3(header, index, imageData)
	case tables.IndexData4:
		out.index, err = parseIndexSubTable4(header, index, imageData)
	case tables.IndexData5:
		out.index, err = parseIndexSubTable5(header, index, imageData)
	}
	return out, err
}

func (subT *bitmapSubtable) image(glyph gID) *bitmapImage {
	return subT.index.imageFor(glyph, subT.first, subT.last)
}

type bitmapIndex interface {
	// first, last is the range of the subtable
	imageFor(glyph gID, first, last gID) *bitmapImage
}

type bitmapImage struct {
	image   []byte
	metrics tables.SmallGlyphMetrics
}

type indexSubTable1And3 struct {
	// length lastGlyph - firstGlyph + 1, elements may be nil
	glyphs []bitmapImage
	format uint16
}

func (idx indexSubTable1And3) imageFor(gid gID, first, last gID) *bitmapImage {
	if gid < first || gid > last {
		return nil
	}
	return &idx.glyphs[gid-first]
}

// imageData starts at the image (table[imageDataOffset:])
func parseIndexSubTable1(header tables.BitmapSubtable, index tables.IndexData1, imageData []byte) (indexSubTable1And3, error) {
	out := indexSubTable1And3{
		format: header.ImageFormat,
		glyphs: make([]bitmapImage, len(index.SbitOffsets)-1),
	}
	for i := range out.glyphs {
		if index.SbitOffsets[i] == index.SbitOffsets[i+1] {
			continue
		}
		var err error
		out.glyphs[i], err = parseBitmapDataMetrics(imageData, index.SbitOffsets[i], index.SbitOffsets[i+1], header.ImageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 1: %s", err)
		}
	}
	return out, nil
}

func parseIndexSubTable3(header tables.BitmapSubtable, index tables.IndexData3, imageData []byte) (indexSubTable1And3, error) {
	out := indexSubTable1And3{
		format: header.ImageFormat,
		glyphs: make([]bitmapImage, len(index.SbitOffsets)-1),
	}
	for i := range out.glyphs {
		if index.SbitOffsets[i] == index.SbitOffsets[i+1] {
			continue
		}
		var err error
		out.glyphs[i], err = parseBitmapDataMetrics(imageData, tables.Offset32(index.SbitOffsets[i]), tables.Offset32(index.SbitOffsets[i+1]), header.ImageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 1: %s", err)
		}
	}
	return out, nil
}

type bitmapDataStandalone []byte

type indexSubTable2 struct {
	glyphs  []bitmapDataStandalone
	format  uint16
	metrics tables.BigGlyphMetrics
}

func (idx indexSubTable2) imageFor(gid gID, first, last gID) *bitmapImage {
	if gid < first || gid > last {
		return nil
	}
	return &bitmapImage{image: idx.glyphs[gid-first], metrics: idx.metrics.SmallGlyphMetrics}
}

// imageData starts at the image (table[imageDataOffset:])
func parseIndexSubTable2(header tables.BitmapSubtable, index tables.IndexData2, imageData []byte) (indexSubTable2, error) {
	out := indexSubTable2{
		format:  header.ImageFormat,
		metrics: index.BigMetrics,
		glyphs:  make([]bitmapDataStandalone, int(header.LastGlyph)-int(header.FirstGlyph)+1),
	}
	for i := range out.glyphs {
		var err error
		out.glyphs[i], err = parseBitmapDataStandalone(imageData, index.ImageSize*uint32(i), index.ImageSize*uint32(i+1), header.ImageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 2: %s", err)
		}
	}
	return out, nil
}

type indexedBitmapGlyph struct {
	data  bitmapImage
	glyph gID
}

type indexSubTable4 struct {
	glyphs []indexedBitmapGlyph
	format uint16
}

func (idx indexSubTable4) imageFor(gid gID, first, last gID) *bitmapImage {
	if gid < first || gid > last {
		return nil
	}
	for i, g := range idx.glyphs {
		if g.glyph == gid {
			return &idx.glyphs[i].data
		}
	}
	return nil
}

// imageData starts at the image (table[imageDataOffset:])
func parseIndexSubTable4(header tables.BitmapSubtable, index tables.IndexData4, imageData []byte) (indexSubTable4, error) {
	out := indexSubTable4{
		format: header.ImageFormat,
		glyphs: make([]indexedBitmapGlyph, len(index.GlyphArray)-1),
	}
	for i := range out.glyphs {
		current, next := index.GlyphArray[i], index.GlyphArray[i+1]
		out.glyphs[i].glyph = current.GlyphID
		var err error
		out.glyphs[i].data, err = parseBitmapDataMetrics(imageData, tables.Offset32(current.SbitOffset), tables.Offset32(next.SbitOffset), header.ImageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 4: %s", err)
		}
	}
	return out, nil
}

type indexSubTable5 struct {
	glyphIndexes []gID                  // sorted by glyph index
	glyphs       []bitmapDataStandalone // corresponding to glyphIndexes
	format       uint16
	metrics      tables.BigGlyphMetrics
}

func (idx indexSubTable5) imageFor(gid gID, first, last gID) *bitmapImage {
	if gid < first || gid > last {
		return nil
	}
	// binary search
	for i, j := 0, len(idx.glyphIndexes); i < j; {
		h := i + (j-i)/2
		entry := idx.glyphIndexes[h]
		if gid < entry {
			j = h
		} else if entry < gid {
			i = h + 1
		} else {
			return &bitmapImage{image: idx.glyphs[h], metrics: idx.metrics.SmallGlyphMetrics}
		}
	}
	return nil
}

// imageData starts at the image (table[imageDataOffset:])
func parseIndexSubTable5(header tables.BitmapSubtable, index tables.IndexData5, imageData []byte) (indexSubTable5, error) {
	out := indexSubTable5{
		format:       header.ImageFormat,
		metrics:      index.BigMetrics,
		glyphIndexes: index.GlyphIdArray,
		glyphs:       make([]bitmapDataStandalone, len(index.GlyphIdArray)),
	}

	for i := range out.glyphs {
		var err error
		out.glyphs[i], err = parseBitmapDataStandalone(imageData, index.ImageSize*uint32(i), (index.ImageSize+1)*uint32(i), header.ImageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 5: %s", err)
		}
	}
	return out, nil
}

func parseBitmapDataMetrics(imageData []byte, start, end tables.Offset32, imageFormat uint16) (bitmapImage, error) {
	if len(imageData) < int(end) || start > end {
		return bitmapImage{}, errors.New("invalid bitmap data table (EOF)")
	}
	imageData = imageData[start:end]
	switch imageFormat {
	case 1, 6, 7, 8, 9:
		return bitmapImage{}, fmt.Errorf("valid but currently not implemented bitmap image format: %d", imageFormat)
	case 2:
		data, _, err := tables.ParseBitmapData2(imageData)
		return bitmapImage{metrics: data.SmallGlyphMetrics, image: data.Image}, err
	case 17:
		data, _, err := tables.ParseBitmapData17(imageData)
		return bitmapImage{metrics: data.SmallGlyphMetrics, image: data.Image}, err
	case 18:
		data, _, err := tables.ParseBitmapData18(imageData)
		return bitmapImage{metrics: data.SmallGlyphMetrics, image: data.Image}, err
	default:
		return bitmapImage{}, fmt.Errorf("unsupported bitmap image format: %d", imageFormat)
	}
}

func parseBitmapDataStandalone(imageData []byte, start, end uint32, format uint16) (bitmapDataStandalone, error) {
	if len(imageData) < int(end) || start > end {
		return nil, fmt.Errorf("invalid bitmap data table (EOF for [%d,%d])", start, end)
	}
	imageData = imageData[start:end]
	switch format {
	case 4:
		return nil, fmt.Errorf("valid but currently not implemented bitmap image format: %d", format)
	case 5:
		data, _, err := tables.ParseBitmapData5(imageData)
		return data.Image, err
	case 19:
		data, _, err := tables.ParseBitmapData19(imageData)
		return data.Image, err
	default:
		return nil, fmt.Errorf("unsupported bitmap image format: %d", format)
	}
}

func maxu16(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

func mulDiv(a, b, c uint16) uint16 {
	return uint16(uint32(a) * uint32(b) / uint32(c))
}
