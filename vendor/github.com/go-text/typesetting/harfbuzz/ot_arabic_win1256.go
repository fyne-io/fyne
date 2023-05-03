package harfbuzz

import (
	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
)

// ported from harfbuzz/src/hb-ot-shape-complex-arabic-win1256.hh Copyright © 2014  Google, Inc. Behdad Esfahbod

type manifest struct {
	lookup *lookupGSUB
	tag    tables.Tag
}

var arabicWin1256GsubLookups = [...]manifest{
	{&rligLookup, loader.NewTag('r', 'l', 'i', 'g')},
	{&initLookup, loader.NewTag('i', 'n', 'i', 't')},
	{&mediLookup, loader.NewTag('m', 'e', 'd', 'i')},
	{&finaLookup, loader.NewTag('f', 'i', 'n', 'a')},
	{&rligMarksLookup, loader.NewTag('r', 'l', 'i', 'g')},
}

// Lookups
var (
	initLookup = lookupGSUB{
		LookupOptions: font.LookupOptions{Flag: otIgnoreMarks},
		Subtables: []tables.GSUBLookup{
			initmediSubLookup,
			initSubLookup,
		},
	}
	mediLookup = lookupGSUB{
		LookupOptions: font.LookupOptions{Flag: otIgnoreMarks},
		Subtables: []tables.GSUBLookup{
			initmediSubLookup,
			mediSubLookup,
			medifinaLamAlefSubLookup,
		},
	}
	finaLookup = lookupGSUB{
		LookupOptions: font.LookupOptions{Flag: otIgnoreMarks},
		Subtables: []tables.GSUBLookup{
			finaSubLookup,
			/* We don't need this one currently as the sequence inherits masks
			 * from the first item. Just in case we change that in the future
			 * to be smart about Arabic masks when ligating... */
			medifinaLamAlefSubLookup,
		},
	}
	rligLookup = lookupGSUB{
		LookupOptions: font.LookupOptions{Flag: otIgnoreMarks},
		Subtables:     []tables.GSUBLookup{lamAlefLigaturesSubLookup},
	}
	rligMarksLookup = lookupGSUB{
		Subtables: []tables.GSUBLookup{shaddaLigaturesSubLookup},
	}
)

// init/medi/fina forms
var (
	initmediSubLookup = tables.SingleSubs{Data: tables.SingleSubstData2{
		Coverage:           tables.Coverage1{Glyphs: []gID{198, 200, 201, 202, 203, 204, 205, 206, 211, 212, 213, 214, 223, 225, 227, 228, 236, 237}},
		SubstituteGlyphIDs: []gID{162, 4, 5, 5, 6, 7, 9, 11, 13, 14, 15, 26, 140, 141, 142, 143, 154, 154},
	}}
	initSubLookup = tables.SingleSubs{Data: tables.SingleSubstData2{
		Coverage:           tables.Coverage1{Glyphs: []gID{218, 219, 221, 222, 229}},
		SubstituteGlyphIDs: []gID{27, 30, 128, 131, 144},
	}}
	mediSubLookup = tables.SingleSubs{Data: tables.SingleSubstData2{
		Coverage:           tables.Coverage1{Glyphs: []gID{218, 219, 221, 222, 229}},
		SubstituteGlyphIDs: []gID{28, 31, 129, 138, 149},
	}}
	finaSubLookup = tables.SingleSubs{Data: tables.SingleSubstData2{
		Coverage:           tables.Coverage1{Glyphs: []gID{194, 195, 197, 198, 199, 201, 204, 205, 206, 218, 219, 229, 236, 237}},
		SubstituteGlyphIDs: []gID{2, 1, 3, 181, 0, 159, 8, 10, 12, 29, 127, 152, 160, 156},
	}}
	medifinaLamAlefSubLookup = tables.SingleSubs{Data: tables.SingleSubstData2{
		Coverage:           tables.Coverage1{Glyphs: []gID{165, 178, 180, 252}},
		SubstituteGlyphIDs: []gID{170, 179, 185, 255},
	}}
)

type ligs = []tables.Ligature

var (
	// Lam+Alef ligatures
	lamAlefLigaturesSubLookup = tables.LigatureSubs{
		Coverage:     tables.Coverage1{Glyphs: []gID{225}},
		LigatureSets: []tables.LigatureSet{{Ligatures: lamLigatureSet}},
	}
	lamLigatureSet = ligs{
		{
			LigatureGlyph:     199,
			ComponentGlyphIDs: []uint16{165},
		},
		{
			LigatureGlyph:     195,
			ComponentGlyphIDs: []uint16{178},
		},
		{
			LigatureGlyph:     194,
			ComponentGlyphIDs: []uint16{180},
		},
		{
			LigatureGlyph:     197,
			ComponentGlyphIDs: []uint16{252},
		},
	}

	// Shadda ligatures
	shaddaLigaturesSubLookup = tables.LigatureSubs{
		Coverage:     tables.Coverage1{Glyphs: []gID{248}},
		LigatureSets: []tables.LigatureSet{{Ligatures: shaddaLigatureSet}},
	}
	shaddaLigatureSet = ligs{
		{
			LigatureGlyph:     243,
			ComponentGlyphIDs: []uint16{172},
		},
		{
			LigatureGlyph:     245,
			ComponentGlyphIDs: []uint16{173},
		},
		{
			LigatureGlyph:     246,
			ComponentGlyphIDs: []uint16{175},
		},
	}
)
