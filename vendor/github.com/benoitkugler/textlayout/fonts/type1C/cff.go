// Package type1c provides a parser for the CFF font format
// defined at https://www.adobe.com/content/dam/acom/en/devnet/font/pdfs/5176.CFF.pdf.
// It can be used to read standalone CFF font files, but is mainly used
// through the truetype package to read embedded CFF glyph tables.
package type1c

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/benoitkugler/textlayout/fonts"
	"github.com/benoitkugler/textlayout/fonts/glyphsnames"
	"github.com/benoitkugler/textlayout/fonts/simpleencodings"
)

// Load reads standalone .cff font files and may
// return multiple fonts.
func Load(file fonts.Resource) ([]Font, error) {
	return parse(file)
}

// Font represents a parsed Font font.
type Font struct {
	userStrings userStrings
	fdSelect    fdSelect // only valid for CIDFonts
	charset     []uint16 // indexed by glyph ID
	Encoding    *simpleencodings.Encoding

	cmap fonts.CmapSimple // see synthetizeCmap

	cidFontName string
	charstrings [][]byte // indexed by glyph ID
	fontName    []byte   // name from the Name INDEX
	globalSubrs [][]byte
	// array of length 1 for non CIDFonts
	// For CIDFonts, it can be safely indexed by `fdSelect` output
	localSubrs [][][]byte
	fonts.PSInfo
}

// Parse parse a .cff font file.
// Although CFF enables multiple font or CIDFont programs to be bundled together in a
// single file, embedded CFF font file in PDF or in TrueType/OpenType fonts
// shall consist of exactly one font or CIDFont. Thus, this function
// returns an error if the file contains more than one font.
// See Loader to read standalone .cff files
func Parse(file fonts.Resource) (*Font, error) {
	fonts, err := parse(file)
	if err != nil {
		return nil, err
	}
	if len(fonts) != 1 {
		return nil, errors.New("only one CFF font is allowed in embedded files")
	}
	return &fonts[0], nil
}

func parse(file fonts.Resource) ([]Font, error) {
	_, err := file.Seek(0, io.SeekStart) // file might have been used before
	if err != nil {
		return nil, err
	}
	// read 4 bytes to check if its a supported CFF file
	var buf [4]byte
	file.Read(buf[:])
	if buf[0] != 1 || buf[1] != 0 || buf[2] != 4 {
		return nil, errUnsupportedCFFVersion
	}
	file.Seek(0, io.SeekStart)

	// if this is really needed, we can modify the parser to directly use `file`
	// without reading all in memory
	input, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	p := cffParser{src: input}
	p.skip(4)
	out, err := p.parse()
	if err != nil {
		return nil, err
	}

	// for standalone use, generate a cmap
	for _, ft := range out {
		ft.synthetizeCmap()
	}

	return out, nil
}

// Type1 fonts have no natural notion of Unicode code points
// We use a glyph names table to identify the most commonly used runes
func (f *Font) synthetizeCmap() {
	f.cmap = make(map[rune]fonts.GID)
	for gid := range f.charstrings {
		glyphName := f.GlyphName(fonts.GID(gid))
		r, _ := glyphsnames.GlyphToRune(glyphName)
		f.cmap[r] = fonts.GID(gid)
	}
}

func (f *Font) Cmap() (fonts.Cmap, fonts.CmapEncoding) {
	return f.cmap, fonts.EncUnicode
}

// GlyphName returns the name of the glyph or an empty string if not found.
func (f *Font) GlyphName(glyph fonts.GID) string {
	if f.fdSelect != nil || int(glyph) >= len(f.charset) {
		return ""
	}
	out, _ := f.userStrings.getString(f.charset[glyph])
	return out
}

// NumGlyphs returns the number of glyphs in this font.
// It is also the maximum glyph index + 1.
func (f *Font) NumGlyphs() int { return len(f.charstrings) }

func (f *Font) PostscriptInfo() (fonts.PSInfo, bool) { return f.PSInfo, true }

func (f *Font) PoscriptName() string { return f.PSInfo.FontName }

// Strip all subset prefixes of the form `ABCDEF+'.  Usually, there
// is only one, but font names like `APCOOG+JFABTD+FuturaBQ-Bold'
// have been seen in the wild.
func removeSubsetPrefix(name []byte) []byte {
	for keep := true; keep; {
		if len(name) >= 7 && name[6] == '+' {
			for idx := 0; idx < 6; idx++ {
				/* ASCII uppercase letters */
				if !('A' <= name[idx] && name[idx] <= 'Z') {
					keep = false
				}
			}
			if keep {
				name = name[7:]
			}
		} else {
			keep = false
		}
	}
	return name
}

// remove the style part from the family name (if present).
func removeStyle(familyName, styleName string) string {
	if lF, lS := len(familyName), len(styleName); lF > lS {
		idx := 1
		for ; idx <= len(styleName); idx++ {
			if familyName[lF-idx] != styleName[lS-idx] {
				break
			}
		}

		if idx > lS {
			// familyName ends with styleName; remove it
			idx = lF - lS - 1

			// also remove special characters
			// between real family name and style
			for idx > 0 &&
				(familyName[idx] == '-' || familyName[idx] == ' ' ||
					familyName[idx] == '_' || familyName[idx] == '+') {
				idx--
			}

			if idx > 0 {
				familyName = familyName[:idx+1]
			}
		}
	}
	return familyName
}

func (f *Font) getStyle() (isItalic, isBold bool, familyName, styleName string) {
	// adapted from freetype/src/cff/cffobjs.c

	// retrieve font family & style name
	familyName = f.PSInfo.FamilyName
	if familyName == "" {
		familyName = string(removeSubsetPrefix(f.fontName))
	}
	if familyName != "" {
		full := f.PSInfo.FullName

		// We try to extract the style name from the full name.
		// We need to ignore spaces and dashes during the search.
		for i, j := 0, 0; i < len(full); {
			// skip common characters at the start of both strings
			if full[i] == familyName[j] {
				i++
				j++
				continue
			}

			// ignore spaces and dashes in full name during comparison
			if full[i] == ' ' || full[i] == '-' {
				i++
				continue
			}

			// ignore spaces and dashes in family name during comparison
			if familyName[j] == ' ' || familyName[j] == '-' {
				j++
				continue
			}

			if j == len(familyName) && i < len(full) {
				/* The full name begins with the same characters as the  */
				/* family name, with spaces and dashes removed.  In this */
				/* case, the remaining string in `full' will be used as */
				/* the style name.                                       */
				styleName = full[i:]

				/* remove the style part from the family name (if present) */
				familyName = removeStyle(familyName, styleName)
			}
			break
		}
	} else {
		// do we have a `/FontName' for a CID-keyed font?
		familyName = f.cidFontName
	}

	styleName = strings.TrimSpace(styleName)
	if styleName == "" {
		// assume "Regular" style if we don't know better
		styleName = "Regular"
	}

	isItalic = f.PSInfo.ItalicAngle != 0
	isBold = f.PSInfo.Weight == "Bold" || f.PSInfo.Weight == "Black"

	// double check
	if !isBold {
		isBold = strings.HasPrefix(styleName, "Bold") || strings.HasPrefix(styleName, "Black")
	}
	return
}

func (f *Font) LoadSummary() (fonts.FontSummary, error) {
	isItalic, isBold, familyName, styleName := f.getStyle()
	return fonts.FontSummary{
		IsItalic:          isItalic,
		IsBold:            isBold,
		Familly:           familyName,
		Style:             styleName,
		HasScalableGlyphs: true,
		HasBitmapGlyphs:   false,
		HasColorGlyphs:    false,
	}, nil
}
