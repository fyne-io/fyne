package truetype

var (
	// tagHead represents the 'head' table, which contains the font header
	tagHead = MustNewTag("head")
	// tagMaxp represents the 'maxp' table, which contains the maximum profile
	tagMaxp = MustNewTag("maxp")
	// tagHmtx represents the 'hmtx' table, which contains the horizontal metrics
	tagHmtx = MustNewTag("hmtx")
	// tagVmtx represents the 'vmtx' table, which contains the vertical metrics
	tagVmtx = MustNewTag("vmtx")
	// tagHhea represents the 'hhea' table, which contains the horizonal header
	tagHhea = MustNewTag("hhea")
	// tagVhea represents the 'vhea' table, which contains the vertical header
	tagVhea = MustNewTag("vhea")
	// tagOS2 represents the 'OS/2' table, which contains windows-specific metadata
	tagOS2 = MustNewTag("OS/2")
	// tagName represents the 'name' table, which contains font name information
	tagName = MustNewTag("name")
	// TagGpos represents the 'GPOS' table, which contains Glyph Positioning features
	TagGpos = MustNewTag("GPOS")
	// TagGsub represents the 'GSUB' table, which contains Glyph Substitution features
	TagGsub = MustNewTag("GSUB")
	// TagGdef represents the 'GDEF' table, which contains various Glyph Definitions
	TagGdef = MustNewTag("GDEF")

	tagCmap = MustNewTag("cmap")
	tagKern = MustNewTag("kern")
	tagPost = MustNewTag("post")
	TagSilf = MustNewTag("Silf")
	TagPrep = MustNewTag("prep")
	tagLoca = MustNewTag("loca")
	tagGlyf = MustNewTag("glyf")
	tagCFF  = MustNewTag("CFF ")
	tagCFF2 = MustNewTag("CFF2")
	tagVorg = MustNewTag("VORG")
	tagSbix = MustNewTag("sbix")
	tagBhed = MustNewTag("bhed")
	tagCBLC = MustNewTag("CBLC")
	tagCBDT = MustNewTag("CBDT")
	tagEBLC = MustNewTag("EBLC")
	tagEBDT = MustNewTag("EBDT")
	tagBloc = MustNewTag("bloc")
	tagBdat = MustNewTag("bdat")
	tagCOLR = MustNewTag("COLR")
	tagFvar = MustNewTag("fvar")
	tagAvar = MustNewTag("avar")
	tagGvar = MustNewTag("gvar")
	tagMvar = MustNewTag("MVAR")
	tagHvar = MustNewTag("HVAR")
	tagVvar = MustNewTag("VVAR")

	tagFeat = MustNewTag("feat")
	tagMorx = MustNewTag("morx")
	tagKerx = MustNewTag("kerx")
	tagAnkr = MustNewTag("ankr")
	tagTrak = MustNewTag("trak")

	// TypeTrueType is the first four bytes of an OpenType file containing a TrueType font
	TypeTrueType = Tag(0x00010000)
	// TypeAppleTrueType is the first four bytes of an OpenType file containing a TrueType font
	// (specifically one designed for Apple products, it's recommended to use TypeTrueType instead)
	TypeAppleTrueType = MustNewTag("true")
	// TypePostScript1 is the first four bytes of an OpenType file containing a PostScript Type 1 font
	TypePostScript1 = MustNewTag("typ1")
	// TypeOpenType is the first four bytes of an OpenType file containing a PostScript Type 2 font
	// as specified by OpenType
	TypeOpenType = MustNewTag("OTTO")

	// SignatureWOFF is the magic number at the start of a WOFF file.
	SignatureWOFF = MustNewTag("wOFF")

	ttcTag = MustNewTag("ttcf")

	// // SignatureWOFF2 is the magic number at the start of a WOFF2 file.
	// SignatureWOFF2 = MustNewTag("wOF2")
)

// dfontResourceDataOffset is the assumed value of a dfont file's resource data
// offset.
//
// https://github.com/kreativekorp/ksfl/wiki/Macintosh-Resource-File-Format
// says that "A Mac OS resource file... [starts with an] offset from start of
// file to start of resource data section... [usually] 0x0100". In theory,
// 0x00000100 isn't always a magic number for identifying dfont files. In
// practice, it seems to work.
const dfontResourceDataOffset = 0x00000100
