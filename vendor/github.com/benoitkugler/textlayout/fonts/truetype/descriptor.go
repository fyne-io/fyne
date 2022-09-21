// Package truetype provides support for OpenType and TrueType font formats, used in PDF.
//
// It is largely influenced by github.com/ConradIrwin/font, golang.org/x/image/font/sfnt,
// and FreeType2.
package truetype

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/textlayout/fonts"
)

func (font *Font) PostscriptInfo() (fonts.PSInfo, bool) {
	return fonts.PSInfo{}, false
}

func (font *Font) Cmap() (fonts.Cmap, fonts.CmapEncoding) { return font.cmap, font.cmapEncoding }

// PoscriptName returns the optional PoscriptName of the font
func (font *Font) PoscriptName() string {
	// adapted from freetype

	// scan the name table to see whether we have a Postscript name here,
	// either in Macintosh or Windows platform encodings
	windows, mac := font.Names.getEntry(NamePostscript)

	// prefer Windows entries over Apple
	if windows != nil {
		return windows.String()
	}
	if mac != nil {
		return mac.String()
	}
	return ""
}

type fontSummary struct {
	head            *TableHead
	os2             *TableOS2
	names           TableName
	hasOutline      bool
	hasBitmap       bool
	hasColor        bool
	hasVerticalInfo bool
}

// loadSummary loads various tables to compute meta data about the font
func (pr *FontParser) loadSummary(font *Font) error {
	// adapted from freetype

	var out fontSummary
	out.names = font.Names
	if pr.HasTable(tagCBLC) || pr.HasTable(tagSbix) || pr.HasTable(tagCOLR) {
		out.hasColor = true
	}
	out.head = &font.Head

	// do we have outlines in there ?
	out.hasOutline = pr.HasTable(tagGlyf) || pr.HasTable(tagCFF) || pr.HasTable(tagCFF2)

	isAppleSbix := pr.HasTable(tagSbix)

	// Apple 'sbix' color bitmaps are rendered scaled and then the 'glyf'
	// outline rendered on top.  We don't support that yet, so just ignore
	// the 'glyf' outline and advertise it as a bitmap-only font.
	if isAppleSbix {
		out.hasOutline = false
	}

	isAppleSbit := pr.isBinary

	hasCblc := pr.HasTable(tagCBLC)
	hasCbdt := pr.HasTable(tagCBDT)

	// Ignore outlines for CBLC/CBDT fonts.
	if hasCblc || hasCbdt {
		out.hasOutline = false
	}

	out.hasBitmap = hasCblc && hasCbdt || pr.HasTable(tagEBDT) && pr.HasTable(tagEBLC) || isAppleSbit || isAppleSbix

	// OpenType 1.8.2 introduced limits to this value;
	// however, they make sense for older SFNT fonts also
	if out.head.UnitsPerEm < 16 || out.head.UnitsPerEm > 16384 {
		return fmt.Errorf("invalid UnitsPerEm value %d", out.head.UnitsPerEm)
	}

	// do not load the metrics headers and tables if this is an Apple
	// sbit font file
	if isAppleSbit {
		return nil
	}

	// load the `hhea' and `hmtx' tables
	if font.hhea != nil {
		_, err := pr.HtmxTable(font.NumGlyphs)
		if err != nil {
			return err
		}
	} else {
		// No `hhea' table necessary for SFNT Mac fonts.
		if font.Type == TypeAppleTrueType {
			out.hasOutline = false
		}
	}

	// try to load the `vhea' and `vmtx' tables
	if font.vhea != nil {
		_, err := pr.VtmxTable(font.NumGlyphs)
		out.hasVerticalInfo = err == nil
	}

	out.os2 = font.OS2 // we treat the table as missing if there are any errors

	font.fontSummary = out
	return nil
}

func (font *Font) LoadSummary() (fonts.FontSummary, error) {
	isItalic, isBold, familyName, styleName := font.fontSummary.getStyle()
	return fonts.FontSummary{
		IsItalic: isItalic,
		IsBold:   isBold,
		Familly:  familyName,
		Style:    styleName,
		// a font with no bitmaps and no outlines is scalable;
		// it has only empty glyphs then
		HasScalableGlyphs: !font.fontSummary.hasBitmap,
		HasBitmapGlyphs:   font.fontSummary.hasBitmap,
		HasColorGlyphs:    font.fontSummary.hasColor,
	}, nil
}

// getStyle sum up the style of the font
func (summary fontSummary) getStyle() (isItalic, isBold bool, familyName, styleName string) {
	// Bit 8 of the `fsSelection' field in the `OS/2' table denotes
	// a WWS-only font face.  `WWS' stands for `weight', width', and
	// `slope', a term used by Microsoft's Windows Presentation
	// Foundation (WPF).  This flag has been introduced in version
	// 1.5 of the OpenType specification (May 2008).

	if summary.os2 != nil && summary.os2.FsSelection&256 != 0 {
		familyName = summary.names.getName(NamePreferredFamily)
		if familyName == "" {
			familyName = summary.names.getName(NameFontFamily)
		}

		styleName = summary.names.getName(NamePreferredSubfamily)
		if styleName == "" {
			styleName = summary.names.getName(NameFontSubfamily)
		}
	} else {
		familyName = summary.names.getName(NameWWSFamily)
		if familyName == "" {
			familyName = summary.names.getName(NamePreferredFamily)
		}
		if familyName == "" {
			familyName = summary.names.getName(NameFontFamily)
		}

		styleName = summary.names.getName(NameWWSSubfamily)
		if styleName == "" {
			styleName = summary.names.getName(NamePreferredSubfamily)
		}
		if styleName == "" {
			styleName = summary.names.getName(NameFontSubfamily)
		}
	}

	styleName = strings.TrimSpace(styleName)
	if styleName == "" { // assume `Regular' style because we don't know better
		styleName = "Regular"
	}

	// Compute style flags.
	if summary.hasOutline && summary.os2 != nil {
		// We have an OS/2 table; use the `fsSelection' field.  Bit 9
		// indicates an oblique font face.  This flag has been
		// introduced in version 1.5 of the OpenType specification.
		isItalic = summary.os2.FsSelection&(1<<9) != 0 || summary.os2.FsSelection&1 != 0
		isBold = summary.os2.FsSelection&(1<<5) != 0
	} else {
		// this is an old Mac font, use the header field
		isBold = summary.head.MacStyle&1 != 0
		isItalic = summary.head.MacStyle&2 != 0
	}

	return
}

// ScanFont lazily parse `file` to extract a summary of the font(s).
// Collections are supported.
func ScanFont(file fonts.Resource) ([]fonts.FontDescriptor, error) {
	parsers, err := NewFontParsers(file)
	if err != nil {
		return nil, err
	}

	out := make([]fonts.FontDescriptor, len(parsers))
	for i, p := range parsers {
		out[i] = newFontDescriptor(p)
	}

	return out, nil
}

var _ fonts.FontDescriptor = (*fontDescriptor)(nil)

type fontDescriptor struct {
	FontParser

	// these tables are required both in Family
	// and Aspect
	os2   *TableOS2
	names TableName
	head  TableHead
}

func newFontDescriptor(pr *FontParser) *fontDescriptor {
	// load required table
	out := fontDescriptor{FontParser: *pr}
	out.os2, _ = pr.OS2Table()
	out.names, _ = pr.tryAndLoadNameTable()
	out.head, _ = pr.loadHeadTable()
	return &out
}

func (fd *fontDescriptor) Family() string {
	var family string
	if fd.os2 != nil && fd.os2.FsSelection&256 != 0 {
		family = fd.names.getName(NamePreferredFamily)
		if family == "" {
			family = fd.names.getName(NameFontFamily)
		}
	} else {
		family = fd.names.getName(NameWWSFamily)
		if family == "" {
			family = fd.names.getName(NamePreferredFamily)
		}
		if family == "" {
			family = fd.names.getName(NameFontFamily)
		}
	}
	return family
}

func (fd *fontDescriptor) AdditionalStyle() string {
	var style string
	if fd.os2 != nil && fd.os2.FsSelection&256 != 0 {
		style = fd.names.getName(NamePreferredSubfamily)
		if style == "" {
			style = fd.names.getName(NameFontSubfamily)
		}
	} else {
		style = fd.names.getName(NameWWSSubfamily)
		if style == "" {
			style = fd.names.getName(NamePreferredSubfamily)
		}
		if style == "" {
			style = fd.names.getName(NameFontSubfamily)
		}
	}
	style = strings.TrimSpace(style)
	return style
}

func (fd *fontDescriptor) Aspect() (style fonts.Style, weight fonts.Weight, stretch fonts.Stretch) {
	if fd.os2 != nil {
		// We have an OS/2 table; use the `fsSelection' field.  Bit 9
		// indicates an oblique font face.  This flag has been
		// introduced in version 1.5 of the OpenType specification.
		if fd.os2.FsSelection&(1<<9) != 0 || fd.os2.FsSelection&1 != 0 {
			style = fonts.StyleItalic
		}

		weight = fonts.Weight(fd.os2.USWeightClass)

		switch fd.os2.USWidthClass {
		case 1:
			stretch = fonts.StretchUltraCondensed
		case 2:
			stretch = fonts.StretchExtraCondensed
		case 3:
			stretch = fonts.StretchCondensed
		case 4:
			stretch = fonts.StretchSemiCondensed
		case 5:
			stretch = fonts.StretchNormal
		case 6:
			stretch = fonts.StretchSemiExpanded
		case 7:
			stretch = fonts.StretchExpanded
		case 8:
			stretch = fonts.StretchExtraExpanded
		case 9:
			stretch = fonts.StretchUltraExpanded
		}

	} else {
		// this is an old Mac font, use the header field
		if isItalic := fd.head.MacStyle&2 != 0; isItalic {
			style = fonts.StyleItalic
		}
		if isBold := fd.head.MacStyle&1 != 0; isBold {
			weight = fonts.WeightBold
		}
	}

	return
}

func (fd *fontDescriptor) LoadCmap() (Cmap, error) {
	cmap, err := fd.FontParser.CmapTable()
	if err != nil {
		return nil, err
	}
	out, _ := cmap.BestEncoding()
	return out, nil
}
