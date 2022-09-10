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
func (pr *FontParser) loadSummary() error {
	// adapted from freetype

	var out fontSummary
	out.names = pr.font.Names
	if pr.HasTable(tagCBLC) || pr.HasTable(tagSbix) || pr.HasTable(tagCOLR) {
		out.hasColor = true
	}
	out.head = &pr.font.Head

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
	if pr.font.hhea != nil {
		_, err := pr.HtmxTable()
		if err != nil {
			return err
		}
	} else {
		// No `hhea' table necessary for SFNT Mac fonts.
		if pr.font.Type == TypeAppleTrueType {
			out.hasOutline = false
		}
	}

	// try to load the `vhea' and `vmtx' tables
	if pr.font.vhea != nil {
		_, err := pr.VtmxTable()
		out.hasVerticalInfo = err == nil
	}

	out.os2 = pr.font.OS2 // we treat the table as missing if there are any errors

	pr.font.fontSummary = out
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
