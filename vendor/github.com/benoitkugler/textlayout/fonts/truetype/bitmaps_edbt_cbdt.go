package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/benoitkugler/textlayout/fonts"
)

// group the location (cblc/eblc/bloc) and the data (cbdt/cbdt/bdat)
type bitmapTable []bitmapSize

// return nil if the table is empty
func (t bitmapTable) chooseStrike(xPpem, yPpem uint16) *bitmapSize {
	if len(t) == 0 {
		return nil
	}
	request := maxu16(xPpem, yPpem)
	if request == 0 {
		request = math.MaxUint16 // choose largest strike
	}
	var (
		bestIndex = 0
		bestPpem  = maxu16(t[0].ppemX, t[0].ppemY)
	)
	for i, s := range t {
		ppem := maxu16(s.ppemX, s.ppemY)
		if request <= ppem && ppem < bestPpem || request > bestPpem && ppem > bestPpem {
			bestIndex = i
			bestPpem = ppem
		}
	}
	return &t[bestIndex]
}

func parseTableBitmap(locationTable, rawDataTable []byte) (bitmapTable, error) {
	if len(locationTable) < 8 {
		return nil, errors.New("invalid bitmap location table (EOF)")
	}
	numSizes := int(binary.BigEndian.Uint32(locationTable[4:]))
	if len(locationTable) < 8+numSizes*bitmapSizeLength {
		return nil, errors.New("invalid bitmap location table (EOF)")
	}
	out := make(bitmapTable, numSizes) // guarded by the check above
	var err error
	for i := range out {
		out[i], err = parseBitmapSize(locationTable, 8+i*bitmapSizeLength, rawDataTable)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

type bitmapSize struct {
	subTables            []indexSubTable
	hori, vert           sbitLineMetrics
	startGlyph, endGlyph gid
	ppemX, ppemY         uint16
	bitDepth             uint8
	flags                uint8
}

func (b *bitmapSize) sizeMetrics(avgWidth, upem uint16) (out fonts.BitmapSize) {
	out.XPpem, out.YPpem = b.ppemX, b.ppemY
	ascender := int16(b.hori.ascender)
	descender := int16(b.hori.descender)

	maxBeforeBl := b.hori.maxBeforeBL
	minAfterBl := b.hori.minAfterBL

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

// return nil when not found
func (b *bitmapSize) findTable(glyph GID) indexSubTable {
	for i, subtable := range b.subTables {
		if f, l := subtable.glyphRange(); f <= glyph && glyph <= l {
			return b.subTables[i]
		}
	}
	return nil
}

const (
	sbitLineMetricsLength = 12
	bitmapSizeLength      = 24 + 2*sbitLineMetricsLength
)

// length as been checked
func parseBitmapSize(data []byte, offset int, rawImageData []byte) (out bitmapSize, err error) {
	strikeData := data[offset:]
	subtableArrayOffset := int(binary.BigEndian.Uint32(strikeData))
	// tablesSize := binary.BigEndian.Uint32(strikeData[4:])
	numberSubtables := int(binary.BigEndian.Uint32(strikeData[8:]))
	// color ref
	out.hori = parseSbitLineMetrics(strikeData[16:])
	out.vert = parseSbitLineMetrics(strikeData[16+sbitLineMetricsLength:])

	strikeData = strikeData[16+2*sbitLineMetricsLength:]

	out.startGlyph = gid(binary.BigEndian.Uint16(strikeData))
	out.endGlyph = gid(binary.BigEndian.Uint16(strikeData[2:]))
	out.ppemX = uint16(strikeData[4])
	out.ppemY = uint16(strikeData[5])
	out.bitDepth = strikeData[6]
	out.flags = strikeData[7]

	if len(data) < subtableArrayOffset+numberSubtables*8 {
		return out, errors.New("invalid bitmap strike subtable (EOF)")
	}

	out.subTables = make([]indexSubTable, numberSubtables)
	for i := range out.subTables {
		firstGlyph := GID(binary.BigEndian.Uint16(data[subtableArrayOffset+8*i:]))
		lastGlyph := GID(binary.BigEndian.Uint16(data[subtableArrayOffset+8*i+2:]))
		additionalOffset := int(binary.BigEndian.Uint32(data[subtableArrayOffset+8*i+4:]))

		out.subTables[i], err = parseIndexSubTableData(data, subtableArrayOffset+additionalOffset, firstGlyph, lastGlyph, rawImageData)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

type sbitLineMetrics struct {
	ascender, descender                        int8
	widthMax                                   uint8
	caretSlopeNumerator, caretSlopeDenominator int8
	caretOffset                                int8
	minOriginSB                                int8
	minAdvanceSB                               int8
	maxBeforeBL                                int8
	minAfterBL                                 int8
}

// data must have suffisant length
func parseSbitLineMetrics(data []byte) (out sbitLineMetrics) {
	out.ascender = int8(data[0])
	out.descender = int8(data[1])
	out.widthMax = data[2]
	out.caretSlopeNumerator = int8(data[3])
	out.caretSlopeDenominator = int8(data[4])
	out.caretOffset = int8(data[5])
	out.minOriginSB = int8(data[6])
	out.minAdvanceSB = int8(data[7])
	out.maxBeforeBL = int8(data[8])
	out.minAfterBL = int8(data[9])
	return out
}

type indexSubTable interface {
	// returns the metrics and content of the glyph image,
	// or nil for out of range glyph indexes
	getImage(GID) *bitmapDataMetrics

	// return the format for this sequence of glyphs
	imageFormat() uint16
	glyphRange() (first, last GID)
}

func (idx indexSubTable1And3) imageFormat() uint16    { return idx.format }
func (idx indexSubTable2) imageFormat() uint16        { return idx.format }
func (idx indexSubTable4) imageFormat() uint16        { return idx.format }
func (idx indexSubTable5) imageFormat() uint16        { return idx.format }
func (idx indexSubTable1And3) glyphRange() (f, l GID) { return idx.firstGlyph, idx.lastGlyph }
func (idx indexSubTable2) glyphRange() (f, l GID)     { return idx.firstGlyph, idx.lastGlyph }
func (idx indexSubTable4) glyphRange() (f, l GID)     { return idx.firstGlyph, idx.lastGlyph }
func (idx indexSubTable5) glyphRange() (f, l GID)     { return idx.firstGlyph, idx.lastGlyph }

func parseIndexSubTableData(data []byte, offset int, firstGlyph, lastGlyph GID, rawData []byte) (indexSubTable, error) {
	if len(data) < offset+8 {
		return nil, errors.New("invalid bitmap index subtable (EOF)")
	}
	data = data[offset:]
	indexFormat := binary.BigEndian.Uint16(data)
	imageFormat := binary.BigEndian.Uint16(data[2:])
	imageDataOffset := int(binary.BigEndian.Uint32(data[4:]))

	if len(rawData) < imageDataOffset {
		return nil, errors.New("invalid bitmap data table (EOF)")
	}
	imageData := rawData[imageDataOffset:]

	switch indexFormat {
	case 1:
		return parseIndexSubTable1(firstGlyph, lastGlyph, imageFormat, imageData, data[8:])
	case 2:
		return parseIndexSubTable2(firstGlyph, lastGlyph, imageFormat, imageData, data[8:])
	case 3:
		return parseIndexSubTable3(firstGlyph, lastGlyph, imageFormat, imageData, data[8:])
	case 4:
		return parseIndexSubTable4(firstGlyph, lastGlyph, imageFormat, imageData, data[8:])
	case 5:
		return parseIndexSubTable5(firstGlyph, lastGlyph, imageFormat, imageData, data[8:])
	default:
		return nil, fmt.Errorf("unsupported bitmap index subtable format: %d", indexFormat)
	}
}

type indexSubTable1And3 struct {
	// length lastGlyph - firstGlyph + 1, elements may be nil
	glyphs                []*bitmapDataMetrics
	firstGlyph, lastGlyph GID
	format                uint16
}

func (idx indexSubTable1And3) getImage(gid GID) *bitmapDataMetrics {
	if gid < idx.firstGlyph || gid > idx.lastGlyph {
		return nil
	}
	return idx.glyphs[gid-idx.firstGlyph]
}

// data starts after the header, imageData at the image
func parseIndexSubTable1(firstGlyph, lastGlyph GID, imageFormat uint16, imageData, data []byte) (out indexSubTable1And3, err error) {
	numGlyphs := int(lastGlyph-firstGlyph) + 1
	if len(data) < (numGlyphs+1)*4 {
		return out, errors.New("invalid bitmap index subtable format 1 (EOF)")
	}
	offsets := parseUint32s(data, numGlyphs+1)
	out.firstGlyph, out.lastGlyph = firstGlyph, lastGlyph
	out.format = imageFormat
	out.glyphs = make([]*bitmapDataMetrics, numGlyphs)
	for i := range out.glyphs {
		if offsets[i] == offsets[i+1] {
			continue
		}
		out.glyphs[i], err = parseBitmapDataMetrics(imageData, offsets[i], offsets[i+1], imageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 1: %s", err)
		}
	}
	return out, nil
}

type indexSubTable2 struct {
	glyphs                []bitmapDataStandalone
	firstGlyph, lastGlyph GID
	format                uint16
	metrics               bigGlyphMetrics
}

func (idx indexSubTable2) getImage(gid GID) *bitmapDataMetrics {
	if gid < idx.firstGlyph || gid > idx.lastGlyph {
		return nil
	}
	return &bitmapDataMetrics{image: idx.glyphs[gid-idx.firstGlyph], metrics: idx.metrics.smallGlyphMetrics}
}

func parseIndexSubTable2(firstGlyph, lastGlyph GID, imageFormat uint16, imageData, data []byte) (out indexSubTable2, err error) {
	numGlyphs := int(lastGlyph) - int(firstGlyph) + 1
	if len(data) < 4+bigGlyphMetricsSize {
		return out, errors.New("invalid bitmap index subtable format 2 (EOF)")
	}
	imageSize := binary.BigEndian.Uint32(data)
	out.firstGlyph, out.lastGlyph = firstGlyph, lastGlyph
	out.format = imageFormat
	out.metrics = parseBigGlyphMetrics(data[4:])
	out.glyphs = make([]bitmapDataStandalone, numGlyphs)
	for i := range out.glyphs {
		out.glyphs[i], err = parseBitmapDataStandalone(imageData, imageSize*uint32(i), imageSize*uint32(i+1), imageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 2: %s", err)
		}
	}
	return out, nil
}

func parseIndexSubTable3(firstGlyph, lastGlyph GID, imageFormat uint16, imageData, data []byte) (out indexSubTable1And3, err error) {
	numGlyphs := int(lastGlyph-firstGlyph) + 1
	offsets, err := parseUint16s(data, numGlyphs+1)
	if err != nil {
		return out, err
	}
	out.firstGlyph, out.lastGlyph = firstGlyph, lastGlyph
	out.glyphs = make([]*bitmapDataMetrics, numGlyphs)
	for i := range out.glyphs {
		if offsets[i] == offsets[i+1] {
			continue
		}
		out.glyphs[i], err = parseBitmapDataMetrics(imageData, uint32(offsets[i]), uint32(offsets[i+1]), imageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 3: %s", err)
		}
	}
	return out, err
}

type indexedBitmapGlyph struct {
	data  bitmapDataMetrics
	glyph GID
}

type indexSubTable4 struct {
	glyphs                []indexedBitmapGlyph
	firstGlyph, lastGlyph GID
	format                uint16
}

func (idx indexSubTable4) getImage(gid GID) *bitmapDataMetrics {
	if gid < idx.firstGlyph || gid > idx.lastGlyph {
		return nil
	}
	for i, g := range idx.glyphs {
		if g.glyph == gid {
			return &idx.glyphs[i].data
		}
	}
	return nil
}

func parseIndexSubTable4(firstGlyph, lastGlyph GID, imageFormat uint16, imageData, data []byte) (out indexSubTable4, err error) {
	if len(data) < 4 {
		return out, errors.New("invalid bitmap index subtable format 4 (EOF)")
	}
	numGlyphs := int(binary.BigEndian.Uint32(data))
	if len(data) < 4+(numGlyphs+1)*4 {
		return out, errors.New("invalid bitmap index subtable format 4 (EOF)")
	}
	out.firstGlyph, out.lastGlyph = firstGlyph, lastGlyph
	out.format = imageFormat
	out.glyphs = make([]indexedBitmapGlyph, numGlyphs)
	var currentOffset, nextOffset uint32
	nextOffset = uint32(binary.BigEndian.Uint16(data[4+2:]))
	for i := range out.glyphs {
		out.glyphs[i].glyph = GID(binary.BigEndian.Uint16(data[4+4*i:]))
		currentOffset = nextOffset
		nextOffset = uint32(binary.BigEndian.Uint16(data[4+4*(i+1)+2:]))
		data, err := parseBitmapDataMetrics(imageData, currentOffset, nextOffset, imageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 4: %s", err)
		}
		out.glyphs[i].data = *data
	}
	return out, nil
}

type indexSubTable5 struct {
	glyphIndexes          []GID                  // sorted by glyph index
	glyphs                []bitmapDataStandalone // corresponding to glyphIndexes
	firstGlyph, lastGlyph GID
	format                uint16
	metrics               bigGlyphMetrics
}

func (idx indexSubTable5) getImage(gid GID) *bitmapDataMetrics {
	if gid < idx.firstGlyph || gid > idx.lastGlyph {
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
			return &bitmapDataMetrics{image: idx.glyphs[h], metrics: idx.metrics.smallGlyphMetrics}
		}
	}
	return nil
}

func parseIndexSubTable5(firstGlyph, lastGlyph GID, imageFormat uint16, imageData, data []byte) (out indexSubTable5, err error) {
	if len(data) < 8+bigGlyphMetricsSize {
		return out, errors.New("invalid bitmap index subtable format 5 (EOF)")
	}
	imageSize := binary.BigEndian.Uint32(data)
	out.firstGlyph, out.lastGlyph = firstGlyph, lastGlyph
	out.format = imageFormat
	out.metrics = parseBigGlyphMetrics(data[4:])
	numGlyphs := int(binary.BigEndian.Uint32(data[4+bigGlyphMetricsSize:]))
	data = data[8+bigGlyphMetricsSize:]
	if len(data) < 2*numGlyphs {
		return out, errors.New("invalid bitmap index subtable format 5 (EOF)")
	}
	out.glyphIndexes = make([]GID, numGlyphs)
	out.glyphs = make([]bitmapDataStandalone, numGlyphs)
	for i := range out.glyphs {
		out.glyphIndexes[i] = GID(binary.BigEndian.Uint16(data[2*i:]))
		out.glyphs[i], err = parseBitmapDataStandalone(imageData, imageSize*uint32(i), (imageSize+1)*uint32(i), imageFormat)
		if err != nil {
			return out, fmt.Errorf("invalid bitmap index format 5: %s", err)
		}
	}
	return out, nil
}

type smallGlyphMetrics struct {
	height       uint8 // Number of rows of data.
	width        uint8 // Number of columns of data.
	horiBearingX int8  // Distance in pixels from the horizontal origin to the left edge of the bitmap.
	horiBearingY int8  // Distance in pixels from the horizontal origin to the top edge of the bitmap.
	horiAdvance  uint8 // Horizontal advance width in pixels.
}

func (s smallGlyphMetrics) glyphExtents() (out fonts.GlyphExtents) {
	out.XBearing = float32(s.horiBearingX)
	out.YBearing = float32(s.horiBearingY)
	out.Width = float32(s.width)
	out.Height = -float32(s.height)
	return out
}

type bigGlyphMetrics struct {
	smallGlyphMetrics

	vertBearingX int8  // Distance in pixels from the vertical origin to the left edge of the bitmap.
	vertBearingY int8  // Distance in pixels from the vertical origin to the top edge of the bitmap.
	vertAdvance  uint8 // Vertical advance width in pixels.
}

const (
	smallGlyphMetricsSize = 5
	bigGlyphMetricsSize   = smallGlyphMetricsSize + 3
)

// data must have a sufficient length
func parseSmallGlyphMetrics(data []byte) (out smallGlyphMetrics) {
	out.height = data[0]
	out.width = data[1]
	out.horiBearingX = int8(data[2])
	out.horiBearingY = int8(data[3])
	out.horiAdvance = data[4]
	return out
}

// data must have a sufficient length
func parseBigGlyphMetrics(data []byte) (out bigGlyphMetrics) {
	out.smallGlyphMetrics = parseSmallGlyphMetrics(data)
	out.vertBearingX = int8(data[5])
	out.vertBearingY = int8(data[6])
	out.vertAdvance = data[7]
	return out
}

// --------------------- actual bitmap data ---------------------
// for now, we simplify the implementation to two cases:
//	- data, metrics (small)
//  - data only

type bitmapDataMetrics struct {
	image   []byte
	metrics smallGlyphMetrics
}

type bitmapDataStandalone []byte

func parseBitmapDataMetrics(imageData []byte, start, end uint32, format uint16) (*bitmapDataMetrics, error) {
	if len(imageData) < int(end) || start > end {
		return nil, errors.New("invalid bitmap data table (EOF)")
	}
	imageData = imageData[start:end]
	switch format {
	case 1, 6, 7, 8, 9:
		return nil, fmt.Errorf("valid but currently not implemented bitmap image format: %d", format)
	case 2:
		return parseBitmapDataFormat2(imageData)
	case 17:
		return parseBitmapDataFormat17(imageData)
	case 18:
		return parseBitmapDataFormat18(imageData)
	default:
		return nil, fmt.Errorf("unsupported bitmap image format: %d", format)
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
		return parseBitmapDataFormat5(imageData)
	case 19:
		return parseBitmapDataFormat19(imageData)
	default:
		return nil, fmt.Errorf("unsupported bitmap image format: %d", format)
	}
}

// small metrics, bit-aligned data
// data start at the image data
func parseBitmapDataFormat2(data []byte) (*bitmapDataMetrics, error) {
	if len(data) < smallGlyphMetricsSize {
		return nil, errors.New("invalid bitmap data format 2 (EOF)")
	}
	return &bitmapDataMetrics{
		metrics: parseSmallGlyphMetrics(data),
		image:   data[smallGlyphMetricsSize:],
	}, nil
}

// Format 5: metrics in CBLC table, bit-aligned image data only
// data start at the image data
func parseBitmapDataFormat5(data []byte) (out bitmapDataStandalone, err error) {
	return data, nil
}

// small metrics, PNG image data
// data start at the image data
func parseBitmapDataFormat17(data []byte) (*bitmapDataMetrics, error) {
	if len(data) < smallGlyphMetricsSize+4 {
		return nil, errors.New("invalid bitmap data format 17 (EOF)")
	}
	var out bitmapDataMetrics
	out.metrics = parseSmallGlyphMetrics(data)
	length := int(binary.BigEndian.Uint32(data[smallGlyphMetricsSize:]))
	if len(data) < smallGlyphMetricsSize+4+length {
		return nil, errors.New("invalid bitmap data format 17 (EOF)")
	}
	out.image = data[smallGlyphMetricsSize+4 : smallGlyphMetricsSize+4+length]
	return &out, nil
}

// big metrics, PNG image data
// data start at the image data
func parseBitmapDataFormat18(data []byte) (*bitmapDataMetrics, error) {
	if len(data) < bigGlyphMetricsSize+4 {
		return nil, errors.New("invalid bitmap data format 18 (EOF)")
	}
	var out bitmapDataMetrics

	// for now, we only use the first metrics
	out.metrics = parseBigGlyphMetrics(data).smallGlyphMetrics
	length := int(binary.BigEndian.Uint32(data[bigGlyphMetricsSize:]))
	if len(data) < bigGlyphMetricsSize+4+length {
		return nil, errors.New("invalid bitmap data format 18 (EOF)")
	}
	out.image = data[bigGlyphMetricsSize+4 : bigGlyphMetricsSize+4+length]
	return &out, nil
}

// Format 19: metrics in CBLC table, PNG image data
// data start at the image data
func parseBitmapDataFormat19(data []byte) (out bitmapDataStandalone, err error) {
	if len(data) < 4 {
		return out, errors.New("invalid bitmap data format 19 (EOF)")
	}
	length := int(binary.BigEndian.Uint32(data))
	if len(data) < 4+length {
		return out, errors.New("invalid bitmap data format 19 (EOF)")
	}
	return data[4 : 4+length], nil
}
