package truetype

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"

	"github.com/benoitkugler/textlayout/fonts"
	type1c "github.com/benoitkugler/textlayout/fonts/type1C"
)

// FontParser serves as an intermediate when reading font files.
// Most of the time, the Parse and Loader.Load functions are enough,
// but `FontParser` may be used on its own when more control over table loading is needed.
type FontParser struct {
	file   fonts.Resource       // source, needed to parse each table
	tables map[Tag]tableSection // header only, contents is processed on demand

	Type Tag

	// True for fonts which include a 'hbed' table instead
	// of a 'head' table. Apple uses it as a flag that a font doesn't have
	// any glyph outlines but only embedded bitmaps
	isBinary bool
}

// NewFontParser reads the `file` header and returns
// a parser.
// `file` will be used to parse tables, and should not be close.
func NewFontParser(file fonts.Resource) (*FontParser, error) {
	return parseOneFont(file, 0, false)
}

// NewFontParsers is the same as `NewFontParser`, but supports collections.
func NewFontParsers(file fonts.Resource) ([]*FontParser, error) {
	_, err := file.Seek(0, io.SeekStart) // file might have been used before
	if err != nil {
		return nil, err
	}

	var bytes [4]byte
	_, err = file.Read(bytes[:])
	if err != nil {
		return nil, err
	}
	magic := newTag(bytes[:])

	file.Seek(0, io.SeekStart)

	var (
		pr             *FontParser
		offsets        []uint32
		relativeOffset bool
	)
	switch magic {
	case SignatureWOFF, TypeTrueType, TypeOpenType, TypePostScript1, TypeAppleTrueType:
		pr, err = parseOneFont(file, 0, false)
	case ttcTag:
		offsets, err = parseTTCHeader(file)
	case dfontResourceDataOffset:
		offsets, err = parseDfont(file)
		relativeOffset = true
	default:
		return nil, fmt.Errorf("unsupported font format %v", bytes)
	}
	if err != nil {
		return nil, err
	}

	// only one font
	if pr != nil {
		return []*FontParser{pr}, nil
	}

	// collection
	out := make([]*FontParser, len(offsets))
	for i, o := range offsets {
		out[i], err = parseOneFont(file, o, relativeOffset)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// tableSection represents a table within the font file.
type tableSection struct {
	offset  uint32 // Offset into the file this table starts.
	length  uint32 // Length of this table within the file.
	zLength uint32 // Uncompressed length of this table.
}

func (pr *FontParser) findTableBuffer(s tableSection) ([]byte, error) {
	var buf []byte

	if s.length != 0 && s.length < s.zLength {
		zbuf := io.NewSectionReader(pr.file, int64(s.offset), int64(s.length))
		r, err := zlib.NewReader(zbuf)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		buf = make([]byte, s.zLength)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
	} else {
		buf = make([]byte, s.length)
		if _, err := pr.file.ReadAt(buf, int64(s.offset)); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

// HasTable returns `true` is the font has the given table.
func (pr *FontParser) HasTable(tag Tag) bool {
	_, has := pr.tables[tag]
	return has
}

// GetRawTable returns the binary content of the given table,
// or an error if not found.
// Note that many tables are already interpreted by this package,
// see the various XXXTable().
func (pr *FontParser) GetRawTable(tag Tag) ([]byte, error) {
	s, found := pr.tables[tag]
	if !found {
		return nil, errMissingTable
	}

	return pr.findTableBuffer(s)
}

// loads the table corresponding to the 'head' tag.
// if a 'bhed' Apple table is present, it replaces the 'head' one
func (pr *FontParser) loadHeadTable() (TableHead, error) {
	s, hasbhed := pr.tables[tagBhed]
	if !hasbhed {
		var hasHead bool
		s, hasHead = pr.tables[tagHead]
		if !hasHead {
			return TableHead{}, errors.New("missing required head (or bhed) table")
		}
	}
	pr.isBinary = hasbhed

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return TableHead{}, err
	}

	return parseTableHead(buf)
}

// loads the table corresponding to the 'name' tag.
// error only if the table is present and invalid
func (pr *FontParser) tryAndLoadNameTable() (TableName, error) {
	s, found := pr.tables[tagName]
	if !found {
		return nil, nil
	}

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return nil, err
	}

	return parseTableName(buf)
}

// GlyfTable parse the 'glyf' table.
// Note that glyphs may be defined in various format (like CFF or bitmaps), and stored
// in other tables.
// `locationIndexFormat` is found in the 'head' table.
func (pr *FontParser) GlyfTable(numGlyphs int, locationIndexFormat int16) (TableGlyf, error) {
	buf, err := pr.GetRawTable(tagLoca)
	if err != nil {
		return nil, err
	}

	loca, err := parseTableLoca(buf, numGlyphs, locationIndexFormat == 1)
	if err != nil {
		return nil, err
	}

	buf, err = pr.GetRawTable(tagGlyf)
	if err != nil {
		return nil, err
	}

	return parseTableGlyf(buf, loca)
}

func (pr *FontParser) cffTable(numGlyphs int) (*type1c.Font, error) {
	buf, err := pr.GetRawTable(tagCFF)
	if err != nil {
		return nil, err
	}

	out, err := type1c.Parse(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	if N := out.NumGlyphs(); N != numGlyphs {
		return nil, fmt.Errorf("invalid number of glyphs in CFF table (%d != %d)", N, numGlyphs)
	}

	return out, nil
}

func (pr *FontParser) sbixTable(numGlyphs int) (tableSbix, error) {
	buf, err := pr.GetRawTable(tagSbix)
	if err != nil {
		return tableSbix{}, err
	}

	return parseTableSbix(buf, numGlyphs)
}

// parse cblc and cbdt tables
func (pr *FontParser) colorBitmapTable() (bitmapTable, error) {
	buf, err := pr.GetRawTable(tagCBLC)
	if err != nil {
		return nil, err
	}

	rawImageData, err := pr.GetRawTable(tagCBDT)
	if err != nil {
		return nil, err
	}

	return parseTableBitmap(buf, rawImageData)
}

// parse eblc and ebdt tables
func (pr *FontParser) grayBitmapTable() (bitmapTable, error) {
	buf, err := pr.GetRawTable(tagEBLC)
	if err != nil {
		return nil, err
	}

	rawImageData, err := pr.GetRawTable(tagEBDT)
	if err != nil {
		return nil, err
	}

	return parseTableBitmap(buf, rawImageData)
}

// parse bloc and bdat tables
func (pr *FontParser) appleBitmapTable() (bitmapTable, error) {
	buf, err := pr.GetRawTable(tagBloc)
	if err != nil {
		return nil, err
	}

	rawImageData, err := pr.GetRawTable(tagBdat)
	if err != nil {
		return nil, err
	}

	return parseTableBitmap(buf, rawImageData)
}

func (pr *FontParser) HheaTable() (*TableHVhea, error) {
	buf, err := pr.GetRawTable(tagHhea)
	if err != nil {
		return nil, err
	}

	return parseTableHVhea(buf)
}

func (pr *FontParser) VheaTable() (*TableHVhea, error) {
	buf, err := pr.GetRawTable(tagVhea)
	if err != nil {
		return nil, err
	}

	return parseTableHVhea(buf)
}

func (pr *FontParser) OS2Table() (*TableOS2, error) {
	buf, err := pr.GetRawTable(tagOS2)
	if err != nil {
		return nil, err
	}

	return parseTableOS2(buf)
}

// GPOSTable returns the Glyph Positioning table identified with the 'GPOS' tag.
func (pr *FontParser) GPOSTable() (TableGPOS, error) {
	buf, err := pr.GetRawTable(TagGpos)
	if err != nil {
		return TableGPOS{}, err
	}

	return parseTableGPOS(buf)
}

// GSUBTable returns the Glyph Substitution table identified with the 'GSUB' tag.
func (pr *FontParser) GSUBTable() (TableGSUB, error) {
	buf, err := pr.GetRawTable(TagGsub)
	if err != nil {
		return TableGSUB{}, err
	}

	return parseTableGSUB(buf)
}

// GDEFTable returns the Glyph Definition table identified with the 'GDEF' tag.
func (pr *FontParser) GDEFTable(nbAxis int) (TableGDEF, error) {
	buf, err := pr.GetRawTable(TagGdef)
	if err != nil {
		return TableGDEF{}, err
	}

	return parseTableGdef(buf, nbAxis)
}

func (pr *FontParser) CmapTable() (TableCmap, error) {
	s, found := pr.tables[tagCmap]
	if !found {
		return TableCmap{}, errors.New("missing required 'cmap' table")
	}

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return TableCmap{}, fmt.Errorf("invalid required cmap table: %s", err)
	}

	return parseTableCmap(buf)
}

// PostTable returns the Post table names
func (pr *FontParser) PostTable(numGlyphs int) (TablePost, error) {
	buf, err := pr.GetRawTable(tagPost)
	if err != nil {
		return TablePost{}, err
	}

	return parseTablePost(buf, uint16(numGlyphs))
}

// svgTable returns the Post table names
func (pr *FontParser) svgTable() (tableSVG, error) {
	buf, err := pr.GetRawTable(tagSVG)
	if err != nil {
		return nil, err
	}

	return parseTableSVG(buf)
}

// NumGlyphs parses the 'maxp' table to find the number of glyphs in the font.
func (pr *FontParser) NumGlyphs() (int, error) {
	buf, err := pr.GetRawTable(tagMaxp)
	if err != nil {
		return -1, err
	}

	return parseTableMaxp(buf)
}

// HtmxTable returns the glyphs horizontal metrics (array of size numGlyphs),
// expressed in fonts units.
func (pr *FontParser) HtmxTable(numGlyphs int) (TableHVmtx, error) {
	hhea, err := pr.HheaTable()
	if err != nil {
		return nil, err
	}

	buf, err := pr.GetRawTable(tagHmtx)
	if err != nil {
		return nil, err
	}

	return parseHVmtxTable(buf, hhea.numOfLongMetrics, uint16(numGlyphs))
}

// VtmxTable returns the glyphs vertical metrics (array of size numGlyphs),
// expressed in fonts units.
func (pr *FontParser) VtmxTable(numGlyphs int) (TableHVmtx, error) {
	vhea, err := pr.VheaTable()
	if err != nil {
		return nil, err
	}

	buf, err := pr.GetRawTable(tagVmtx)
	if err != nil {
		return nil, err
	}

	return parseHVmtxTable(buf, vhea.numOfLongMetrics, uint16(numGlyphs))
}

// KernTable parses and returns the 'kern' table.
func (pr *FontParser) KernTable(numGlyphs int) (TableKernx, error) {
	buf, err := pr.GetRawTable(tagKern)
	if err != nil {
		return nil, err
	}

	return parseKernTable(buf, numGlyphs)
}

// MorxTable parse the AAT 'morx' table.
func (pr *FontParser) MorxTable(numGlyphs int) (TableMorx, error) {
	buf, err := pr.GetRawTable(tagMorx)
	if err != nil {
		return nil, err
	}

	return parseTableMorx(buf, numGlyphs)
}

// KerxTable parse the AAT 'kerx' table.
func (pr *FontParser) KerxTable(numGlyphs int) (TableKernx, error) {
	buf, err := pr.GetRawTable(tagKerx)
	if err != nil {
		return nil, err
	}

	return parseTableKerx(buf, numGlyphs)
}

// AnkrTable parse the AAT 'ankr' table.
func (pr *FontParser) AnkrTable(numGlyphs int) (TableAnkr, error) {
	buf, err := pr.GetRawTable(tagAnkr)
	if err != nil {
		return TableAnkr{}, err
	}

	return parseTableAnkr(buf, numGlyphs)
}

// TrakTable parse the AAT 'trak' table.
func (pr *FontParser) TrakTable() (TableTrak, error) {
	buf, err := pr.GetRawTable(tagTrak)
	if err != nil {
		return TableTrak{}, err
	}

	return parseTrakTable(buf)
}

// FeatTable parse the AAT 'feat' table.
func (pr *FontParser) FeatTable() (TableFeat, error) {
	buf, err := pr.GetRawTable(tagFeat)
	if err != nil {
		return nil, err
	}

	return parseTableFeat(buf)
}

// error only if the table is present and invalid
func (pr *FontParser) tryAndLoadFvarTable(names TableName) (TableFvar, error) {
	s, found := pr.tables[tagFvar]
	if !found {
		return TableFvar{}, nil
	}

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return TableFvar{}, err
	}

	return parseTableFvar(buf, names)
}

// error only if the table is present and invalid
func (pr *FontParser) tryAndLoadAvarTable(fvar TableFvar) (tableAvar, error) {
	s, found := pr.tables[tagAvar]
	if !found {
		return nil, nil
	}

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return nil, err
	}

	return parseTableAvar(buf, len(fvar.Axis))
}

func (pr *FontParser) gvarTable(glyphs TableGlyf, fvar TableFvar) (tableGvar, error) {
	buf, err := pr.GetRawTable(tagGvar)
	if err != nil {
		return tableGvar{}, err
	}

	return parseTableGvar(buf, len(fvar.Axis), glyphs)
}

func (pr *FontParser) hvarTable(fvar TableFvar) (tableHVvar, error) {
	buf, err := pr.GetRawTable(tagHvar)
	if err != nil {
		return tableHVvar{}, err
	}

	return parseTableHVvar(buf, len(fvar.Axis))
}

func (pr *FontParser) vvarTable(fvar TableFvar) (tableHVvar, error) {
	buf, err := pr.GetRawTable(tagVvar)
	if err != nil {
		return tableHVvar{}, err
	}

	return parseTableHVvar(buf, len(fvar.Axis))
}

func (pr *FontParser) mvarTable(fvar TableFvar) (TableMvar, error) {
	buf, err := pr.GetRawTable(tagMvar)
	if err != nil {
		return TableMvar{}, err
	}

	return parseTableMvar(buf, len(fvar.Axis))
}

func (pr *FontParser) vorgTable() (tableVorg, error) {
	buf, err := pr.GetRawTable(tagVorg)
	if err != nil {
		return tableVorg{}, err
	}

	return parseTableVorg(buf)
}

// best effort to load all valid tables
func (pr *FontParser) loadLayoutTables(numGlyphs int, fvar TableFvar) (out LayoutTables) {
	if tb, err := pr.GDEFTable(len(fvar.Axis)); err == nil {
		out.GDEF = tb
	}
	if tb, err := pr.GSUBTable(); err == nil {
		out.GSUB = tb
	}
	if tb, err := pr.GPOSTable(); err == nil {
		out.GPOS = tb
	}

	if tb, err := pr.MorxTable(numGlyphs); err == nil {
		out.Morx = tb
	}
	if tb, err := pr.KernTable(numGlyphs); err == nil {
		out.Kern = tb
	}
	if tb, err := pr.KerxTable(numGlyphs); err == nil {
		out.Kerx = tb
	}
	if tb, err := pr.AnkrTable(numGlyphs); err == nil {
		out.Ankr = tb
	}
	if tb, err := pr.TrakTable(); err == nil {
		out.Trak = tb
	}
	if tb, err := pr.FeatTable(); err == nil {
		out.Feat = tb
	}

	return out
}

// graphite support

var (
	tagSilf         = MustNewTag("Silf")
	tagSill         = MustNewTag("Sill")
	tagGraphiteFeat = MustNewTag("Feat")
	tagGloc         = MustNewTag("Gloc")
	tagGlat         = MustNewTag("Glat")
)

type GraphiteTables struct {
	Sill, Feat, Gloc, Glat, Silf []byte
}

// LoadGraphiteTables returns the raw tables required for
// Graphite engine support.
// See the package graphite for how to interpret these tables.
func (pr *FontParser) LoadGraphiteTables() (gr GraphiteTables, err error) {
	gr.Sill, err = pr.GetRawTable(tagSill)
	if err != nil {
		return gr, fmt.Errorf("loading table Sill: %s", err)
	}

	gr.Feat, err = pr.GetRawTable(tagGraphiteFeat)
	if err != nil {
		return gr, fmt.Errorf("loading table Feat: %s", err)
	}

	gr.Gloc, err = pr.GetRawTable(tagGloc)
	if err != nil {
		return gr, fmt.Errorf("loading table Gloc: %s", err)
	}

	gr.Glat, err = pr.GetRawTable(tagGlat)
	if err != nil {
		return gr, fmt.Errorf("loading table Glat: %s", err)
	}

	gr.Silf, err = pr.GetRawTable(tagSilf)
	if err != nil {
		return gr, fmt.Errorf("loading table Silf: %s", err)
	}

	return gr, nil
}

func parseOneFont(file fonts.Resource, offset uint32, relativeOffset bool) (parser *FontParser, err error) {
	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("invalid offset: %s", err)
	}

	var bytes [4]byte
	_, err = file.Read(bytes[:])
	if err != nil {
		return nil, err
	}
	magic := newTag(bytes[:])

	switch magic {
	case SignatureWOFF:
		parser, err = parseWOFF(file, offset, relativeOffset)
	case TypeTrueType, TypeOpenType, TypePostScript1, TypeAppleTrueType:
		parser, err = parseOTF(file, offset, relativeOffset)
	default:
		// no more collections allowed here
		return nil, errUnsupportedFormat
	}

	if err != nil {
		return nil, err
	}

	return parser, nil
}

// loadTables calls all the functions loading the
// various font tables,
// and return the loaded font
func (pr *FontParser) loadTables() (*Font, error) {
	var (
		out Font
		err error
	)
	out.Type = pr.Type

	out.NumGlyphs, err = pr.NumGlyphs()
	if err != nil {
		return nil, err
	}
	cmaps, err := pr.CmapTable()
	if err != nil {
		return nil, err
	}
	out.Head, err = pr.loadHeadTable()
	if err != nil {
		return nil, err
	}
	out.Names, err = pr.tryAndLoadNameTable()
	if err != nil {
		return nil, err
	}
	out.fvar, err = pr.tryAndLoadFvarTable(out.Names)
	if err != nil {
		return nil, err
	}
	out.avar, err = pr.tryAndLoadAvarTable(out.fvar)
	if err != nil {
		return nil, err
	}

	out.upem = out.Head.Upem()

	out.OS2, _ = pr.OS2Table()

	out.Glyf, _ = pr.GlyfTable(out.NumGlyphs, out.Head.indexToLocFormat)

	out.bitmap = pr.selectBitmapTable()

	out.sbix, _ = pr.sbixTable(out.NumGlyphs)
	out.cff, _ = pr.cffTable(out.NumGlyphs)
	out.post, _ = pr.PostTable(out.NumGlyphs)
	out.svg, _ = pr.svgTable()

	out.hhea, _ = pr.HheaTable()
	out.vhea, _ = pr.VheaTable()
	out.Hmtx, _ = pr.HtmxTable(out.NumGlyphs)
	out.vmtx, _ = pr.VtmxTable(out.NumGlyphs)

	if len(out.fvar.Axis) != 0 {
		out.mvar, _ = pr.mvarTable(out.fvar)
		out.gvar, _ = pr.gvarTable(out.Glyf, out.fvar)
		if v, err := pr.hvarTable(out.fvar); err == nil {
			out.hvar = &v
		}
		if v, err := pr.vvarTable(out.fvar); err == nil {
			out.vvar = &v
		}
	}

	out.cmap, out.cmapEncoding = cmaps.BestEncoding()
	out.cmapVar = cmaps.unicodeVariation

	if vorg, err := pr.vorgTable(); err == nil {
		out.vorg = &vorg
	}

	out.layoutTables = pr.loadLayoutTables(out.NumGlyphs, out.fvar)

	if pr.HasTable(TagSilf) {
		var gr GraphiteTables
		gr, err = pr.LoadGraphiteTables()
		if err != nil {
			return nil, err
		}
		out.Graphite = &gr
	}

	if pr.HasTable(TagPrep) {
		out.HasHint = true
	}

	err = pr.loadSummary(&out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Parse parses an OpenType or TrueType file and returns a Font.
// See Loader for support for collections, and FontParser for
// more control over table loading.
func Parse(file fonts.Resource) (*Font, error) {
	pr, err := NewFontParser(file)
	if err != nil {
		return nil, err
	}

	return pr.loadTables()
}

// Load implements fonts.FontLoader. For collection font files (.ttc, .otc),
// multiple fonts may be returned.
func Load(file fonts.Resource) (fonts.Faces, error) {
	prs, err := NewFontParsers(file)
	if err != nil {
		return nil, err
	}
	out := make(fonts.Faces, len(prs))
	for i, pr := range prs {
		out[i], err = pr.loadTables()
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
