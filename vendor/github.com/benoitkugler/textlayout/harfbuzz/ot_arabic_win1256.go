package harfbuzz

import (
	tt "github.com/benoitkugler/textlayout/fonts/truetype"
)

// ported from harfbuzz/src/hb-ot-shape-complex-arabic-win1256.hh Copyright © 2014  Google, Inc. Behdad Esfahbod

type manifest struct {
	lookup *lookupGSUB
	tag    tt.Tag
}

var arabicWin1256GsubLookups = [...]manifest{
	{&rligLookup, tt.NewTag('r', 'l', 'i', 'g')},
	{&initLookup, tt.NewTag('i', 'n', 'i', 't')},
	{&mediLookup, tt.NewTag('m', 'e', 'd', 'i')},
	{&finaLookup, tt.NewTag('f', 'i', 'n', 'a')},
	{&rligMarksLookup, tt.NewTag('r', 'l', 'i', 'g')},
}

// Lookups
var (
	initLookup = lookupGSUB{
		LookupOptions: tt.LookupOptions{Flag: tt.IgnoreMarks},
		Subtables: []tt.GSUBSubtable{
			initmediSubLookup,
			initSubLookup,
		},
	}
	mediLookup = lookupGSUB{
		LookupOptions: tt.LookupOptions{Flag: tt.IgnoreMarks},
		Subtables: []tt.GSUBSubtable{
			initmediSubLookup,
			mediSubLookup,
			medifinaLamAlefSubLookup,
		},
	}
	finaLookup = lookupGSUB{
		LookupOptions: tt.LookupOptions{Flag: tt.IgnoreMarks},
		Subtables: []tt.GSUBSubtable{
			finaSubLookup,
			/* We don't need this one currently as the sequence inherits masks
			 * from the first item. Just in case we change that in the future
			 * to be smart about Arabic masks when ligating... */
			medifinaLamAlefSubLookup,
		},
	}
	rligLookup = lookupGSUB{
		LookupOptions: tt.LookupOptions{Flag: tt.IgnoreMarks},
		Subtables:     []tt.GSUBSubtable{lamAlefLigaturesSubLookup},
	}
	rligMarksLookup = lookupGSUB{
		Subtables: []tt.GSUBSubtable{shaddaLigaturesSubLookup},
	}
)

// init/medi/fina forms
var (
	initmediSubLookup = tt.GSUBSubtable{
		Coverage: tt.CoverageList{198, 200, 201, 202, 203, 204, 205, 206, 211, 212, 213, 214, 223, 225, 227, 228, 236, 237},
		Data:     tt.GSUBSingle2{162, 4, 5, 5, 6, 7, 9, 11, 13, 14, 15, 26, 140, 141, 142, 143, 154, 154},
	}
	initSubLookup = tt.GSUBSubtable{
		Coverage: tt.CoverageList{218, 219, 221, 222, 229},
		Data:     tt.GSUBSingle2{27, 30, 128, 131, 144},
	}
	mediSubLookup = tt.GSUBSubtable{
		Coverage: tt.CoverageList{218, 219, 221, 222, 229},
		Data:     tt.GSUBSingle2{28, 31, 129, 138, 149},
	}
	finaSubLookup = tt.GSUBSubtable{
		Coverage: tt.CoverageList{194, 195, 197, 198, 199, 201, 204, 205, 206, 218, 219, 229, 236, 237},
		Data:     tt.GSUBSingle2{2, 1, 3, 181, 0, 159, 8, 10, 12, 29, 127, 152, 160, 156},
	}
	medifinaLamAlefSubLookup = tt.GSUBSubtable{
		Coverage: tt.CoverageList{165, 178, 180, 252},
		Data:     tt.GSUBSingle2{170, 179, 185, 255},
	}
)

type ligs = []tt.LigatureGlyph

var (
	// Lam+Alef ligatures
	lamAlefLigaturesSubLookup = tt.GSUBSubtable{
		Coverage: tt.CoverageList{225},
		Data:     tt.GSUBLigature1{lamLigatureSet},
	}
	lamLigatureSet = ligs{
		tt.LigatureGlyph{
			Glyph:      199,
			Components: []uint16{165},
		},
		tt.LigatureGlyph{
			Glyph:      195,
			Components: []uint16{178},
		},
		tt.LigatureGlyph{
			Glyph:      194,
			Components: []uint16{180},
		},
		tt.LigatureGlyph{
			Glyph:      197,
			Components: []uint16{252},
		},
	}

	// Shadda ligatures
	shaddaLigaturesSubLookup = tt.GSUBSubtable{
		Coverage: tt.CoverageList{248},
		Data:     tt.GSUBLigature1{shaddaLigatureSet},
	}
	shaddaLigatureSet = ligs{
		tt.LigatureGlyph{
			Glyph:      243,
			Components: []uint16{172},
		},
		tt.LigatureGlyph{
			Glyph:      245,
			Components: []uint16{173},
		},
		tt.LigatureGlyph{
			Glyph:      246,
			Components: []uint16{175},
		},
	}
)
