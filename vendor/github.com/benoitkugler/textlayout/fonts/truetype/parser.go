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
	file   fonts.Resource        // source, needed to parse each table
	tables map[Tag]*tableSection // header only, contents is processed on demand

	// Cmaps is not empty after successful parsing
	Cmaps TableCmap

	font Font // target font to fill

	// True for fonts which include a 'hbed' table instead
	// of a 'head' table. Apple uses it as a flag that a font doesn't have
	// any glyph outlines but only embedded bitmaps
	isBinary bool
}

// tableSection represents a table within the font file.
type tableSection struct {
	offset  uint32 // Offset into the file this table starts.
	length  uint32 // Length of this table within the file.
	zLength uint32 // Uncompressed length of this table.
}

func (pr *FontParser) findTableBuffer(s *tableSection) ([]byte, error) {
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
func (pr *FontParser) loadHeadTable() error {
	s, hasbhed := pr.tables[tagBhed]
	if !hasbhed {
		var hasHead bool
		s, hasHead = pr.tables[tagHead]
		if !hasHead {
			return errors.New("missing required head (or bhed) table")
		}
	}
	pr.isBinary = hasbhed

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return err
	}

	pr.font.Head, err = parseTableHead(buf)
	return err
}

// loads the table corresponding to the 'name' tag.
// error only if the table is present and invalid
func (pr *FontParser) tryAndLoadNameTable() error {
	s, found := pr.tables[tagName]
	if !found {
		return nil
	}

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return err
	}

	pr.font.Names, err = parseTableName(buf)
	return err
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

func (pr *FontParser) sbixTable() (tableSbix, error) {
	buf, err := pr.GetRawTable(tagSbix)
	if err != nil {
		return tableSbix{}, err
	}

	return parseTableSbix(buf, pr.font.NumGlyphs)
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

func (pr *FontParser) loadCmapTable() error {
	s, found := pr.tables[tagCmap]
	if !found {
		return errors.New("missing required 'cmap' table")
	}

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return fmt.Errorf("invalid required cmap table: %s", err)
	}

	pr.Cmaps, err = parseTableCmap(buf)
	return err
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
func (pr *FontParser) HtmxTable() (TableHVmtx, error) {
	hhea, err := pr.HheaTable()
	if err != nil {
		return nil, err
	}

	buf, err := pr.GetRawTable(tagHmtx)
	if err != nil {
		return nil, err
	}

	return parseHVmtxTable(buf, hhea.numOfLongMetrics, uint16(pr.font.NumGlyphs))
}

// VtmxTable returns the glyphs vertical metrics (array of size numGlyphs),
// expressed in fonts units.
func (pr *FontParser) VtmxTable() (TableHVmtx, error) {
	vhea, err := pr.VheaTable()
	if err != nil {
		return nil, err
	}

	buf, err := pr.GetRawTable(tagVmtx)
	if err != nil {
		return nil, err
	}

	return parseHVmtxTable(buf, vhea.numOfLongMetrics, uint16(pr.font.NumGlyphs))
}

// KernTable parses and returns the 'kern' table.
func (pr *FontParser) KernTable() (TableKernx, error) {
	buf, err := pr.GetRawTable(tagKern)
	if err != nil {
		return nil, err
	}

	return parseKernTable(buf, pr.font.NumGlyphs)
}

// MorxTable parse the AAT 'morx' table.
func (pr *FontParser) MorxTable() (TableMorx, error) {
	buf, err := pr.GetRawTable(tagMorx)
	if err != nil {
		return nil, err
	}

	return parseTableMorx(buf, pr.font.NumGlyphs)
}

// KerxTable parse the AAT 'kerx' table.
func (pr *FontParser) KerxTable() (TableKernx, error) {
	buf, err := pr.GetRawTable(tagKerx)
	if err != nil {
		return nil, err
	}

	return parseTableKerx(buf, pr.font.NumGlyphs)
}

// AnkrTable parse the AAT 'ankr' table.
func (pr *FontParser) AnkrTable() (TableAnkr, error) {
	buf, err := pr.GetRawTable(tagAnkr)
	if err != nil {
		return TableAnkr{}, err
	}

	return parseTableAnkr(buf, pr.font.NumGlyphs)
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
func (pr *FontParser) tryAndLoadFvarTable() error {
	s, found := pr.tables[tagFvar]
	if !found {
		return nil
	}

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return err
	}

	pr.font.fvar, err = parseTableFvar(buf, pr.font.Names)
	return err
}

// error only if the table is present and invalid
func (pr *FontParser) tryAndLoadAvarTable() error {
	s, found := pr.tables[tagAvar]
	if !found {
		return nil
	}

	buf, err := pr.findTableBuffer(s)
	if err != nil {
		return err
	}

	pr.font.avar, err = parseTableAvar(buf, len(pr.font.fvar.Axis))
	return err
}

func (pr *FontParser) gvarTable(glyphs TableGlyf) (tableGvar, error) {
	buf, err := pr.GetRawTable(tagGvar)
	if err != nil {
		return tableGvar{}, err
	}

	return parseTableGvar(buf, len(pr.font.fvar.Axis), glyphs)
}

func (pr *FontParser) hvarTable() (tableHVvar, error) {
	buf, err := pr.GetRawTable(tagHvar)
	if err != nil {
		return tableHVvar{}, err
	}

	return parseTableHVvar(buf, len(pr.font.fvar.Axis))
}

func (pr *FontParser) vvarTable() (tableHVvar, error) {
	buf, err := pr.GetRawTable(tagVvar)
	if err != nil {
		return tableHVvar{}, err
	}

	return parseTableHVvar(buf, len(pr.font.fvar.Axis))
}

func (pr *FontParser) mvarTable() (TableMvar, error) {
	buf, err := pr.GetRawTable(tagMvar)
	if err != nil {
		return TableMvar{}, err
	}

	return parseTableMvar(buf, len(pr.font.fvar.Axis))
}

func (pr *FontParser) vorgTable() (tableVorg, error) {
	buf, err := pr.GetRawTable(tagVorg)
	if err != nil {
		return tableVorg{}, err
	}

	return parseTableVorg(buf)
}

func (pr *FontParser) loadLayoutTables() {
	if tb, err := pr.GDEFTable(len(pr.font.fvar.Axis)); err == nil {
		pr.font.layoutTables.GDEF = tb
	}
	if tb, err := pr.GSUBTable(); err == nil {
		pr.font.layoutTables.GSUB = tb
	}
	if tb, err := pr.GPOSTable(); err == nil {
		pr.font.layoutTables.GPOS = tb
	}

	if tb, err := pr.MorxTable(); err == nil {
		pr.font.layoutTables.Morx = tb
	}
	if tb, err := pr.KernTable(); err == nil {
		pr.font.layoutTables.Kern = tb
	}
	if tb, err := pr.KerxTable(); err == nil {
		pr.font.layoutTables.Kerx = tb
	}
	if tb, err := pr.AnkrTable(); err == nil {
		pr.font.layoutTables.Ankr = tb
	}
	if tb, err := pr.TrakTable(); err == nil {
		pr.font.layoutTables.Trak = tb
	}
	if tb, err := pr.FeatTable(); err == nil {
		pr.font.layoutTables.Feat = tb
	}
}

func (pr *FontParser) loadMainTables() {
	if pr.font.Head.UnitsPerEm < 16 || pr.font.Head.UnitsPerEm > 16384 {
		pr.font.upem = 1000
	} else {
		pr.font.upem = pr.font.Head.UnitsPerEm
	}

	pr.font.OS2, _ = pr.OS2Table()

	pr.font.Glyf, _ = pr.GlyfTable(pr.font.NumGlyphs, pr.font.Head.indexToLocFormat)

	pr.font.bitmap = pr.selectBitmapTable()

	pr.font.sbix, _ = pr.sbixTable()
	pr.font.cff, _ = pr.cffTable(pr.font.NumGlyphs)
	pr.font.post, _ = pr.PostTable(pr.font.NumGlyphs)
	pr.font.svg, _ = pr.svgTable()

	pr.font.hhea, _ = pr.HheaTable()
	pr.font.vhea, _ = pr.VheaTable()
	pr.font.Hmtx, _ = pr.HtmxTable()
	pr.font.vmtx, _ = pr.VtmxTable()

	if len(pr.font.fvar.Axis) != 0 {
		pr.font.mvar, _ = pr.mvarTable()
		pr.font.gvar, _ = pr.gvarTable(pr.font.Glyf)
		if v, err := pr.hvarTable(); err == nil {
			pr.font.hvar = &v
		}
		if v, err := pr.vvarTable(); err == nil {
			pr.font.vvar = &v
		}
	}

	pr.font.cmap, pr.font.cmapEncoding = pr.Cmaps.BestEncoding()
	pr.font.cmapVar = pr.Cmaps.unicodeVariation

	if vorg, err := pr.vorgTable(); err == nil {
		pr.font.vorg = &vorg
	}
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
	var err error
	pr.font.NumGlyphs, err = pr.NumGlyphs()
	if err != nil {
		return nil, err
	}
	err = pr.loadCmapTable()
	if err != nil {
		return nil, err
	}
	err = pr.loadHeadTable()
	if err != nil {
		return nil, err
	}
	err = pr.tryAndLoadNameTable()
	if err != nil {
		return nil, err
	}
	err = pr.tryAndLoadFvarTable()
	if err != nil {
		return nil, err
	}
	err = pr.tryAndLoadAvarTable()
	if err != nil {
		return nil, err
	}

	pr.loadMainTables()

	pr.loadLayoutTables()

	if pr.HasTable(TagSilf) {
		var gr GraphiteTables
		gr, err = pr.LoadGraphiteTables()
		if err != nil {
			return nil, err
		}
		pr.font.Graphite = &gr
	}

	if pr.HasTable(TagPrep) {
		// TODO: load the table
		pr.font.HasHint = true
	}

	err = pr.loadSummary()
	if err != nil {
		return nil, err
	}

	return &pr.font, nil
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
func (loader) Load(file fonts.Resource) (fonts.Faces, error) {
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
