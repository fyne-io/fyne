// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"sort"
)

func (c Coverage1) Index(gi GlyphID) (int, bool) {
	num := len(c.Glyphs)
	idx := sort.Search(num, func(i int) bool { return gi <= c.Glyphs[i] })
	if idx < num && c.Glyphs[idx] == gi {
		return idx, true
	}
	return 0, false
}

func (cl Coverage1) Len() int { return len(cl.Glyphs) }

func (c Coverage2) Index(gi GlyphID) (int, bool) {
	num := len(c.Ranges)
	if num == 0 {
		return 0, false
	}

	idx := sort.Search(num, func(i int) bool { return gi <= c.Ranges[i].StartGlyphID })
	// idx either points to a matching start, or to the next range (or idx==num)
	// e.g. with the range example from above: 130 points to 130-135 range, 133 points to 137-137 range

	// check if gi is the start of a range, but only if sort.Search returned a valid result
	if idx < num {
		if rang := c.Ranges[idx]; gi == rang.StartGlyphID {
			return int(rang.StartCoverageIndex), true
		}
	}
	// check if gi is in previous range
	if idx > 0 {
		idx--
		if rang := c.Ranges[idx]; gi >= rang.StartGlyphID && gi <= rang.EndGlyphID {
			return int(rang.StartCoverageIndex) + int(gi-rang.StartGlyphID), true
		}
	}

	return 0, false
}

func (cr Coverage2) Len() int {
	size := 0
	for _, r := range cr.Ranges {
		size += int(r.EndGlyphID - r.StartGlyphID + 1)
	}
	return size
}

func (cl ClassDef1) Class(gi GlyphID) (uint16, bool) {
	if gi < cl.StartGlyphID || gi >= cl.StartGlyphID+GlyphID(len(cl.ClassValueArray)) {
		return 0, false
	}
	return cl.ClassValueArray[gi-cl.StartGlyphID], true
}

func (cl ClassDef1) Extent() int {
	max := uint16(0)
	for _, cid := range cl.ClassValueArray {
		if cid >= max {
			max = cid
		}
	}
	return int(max) + 1
}

func (cl ClassDef2) Class(g GlyphID) (uint16, bool) {
	// 'adapted' from golang/x/image/font/sfnt
	c := cl.ClassRangeRecords
	num := len(c)
	if num == 0 {
		return 0, false
	}

	// classRange is an array of startGlyphID, endGlyphID and target class ID.
	// Ranges are non-overlapping.
	// E.g. 130, 135, 1   137, 137, 5   etc

	idx := sort.Search(num, func(i int) bool { return g <= c[i].StartGlyphID })
	// idx either points to a matching start, or to the next range (or idx==num)
	// e.g. with the range example from above: 130 points to 130-135 range, 133 points to 137-137 range

	// check if gi is the start of a range, but only if sort.Search returned a valid result
	if idx < num {
		if class := c[idx]; g == c[idx].StartGlyphID {
			return class.Class, true
		}
	}
	// check if gi is in previous range
	if idx > 0 {
		idx--
		if class := c[idx]; g >= class.StartGlyphID && g <= class.EndGlyphID {
			return class.Class, true
		}
	}

	return 0, false
}

func (cl ClassDef2) Extent() int {
	max := uint16(0)
	for _, r := range cl.ClassRangeRecords {
		if r.Class >= max {
			max = r.Class
		}
	}
	return int(max) + 1
}

// ------------------------------------ layout getters ------------------------------------

// FindLanguage looks for [language] and return its index into the [LangSys] slice,
// or -1 if the tag is not found.
func (sc Script) FindLanguage(language Tag) int {
	// LangSys is sorted: binary search
	low, high := 0, len(sc.LangSysRecords)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := sc.LangSysRecords[mid].Tag
		if language < p {
			high = mid
		} else if language > p {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

// GetLangSys return the language at [index]. It [index] is out of range (for example with 0xFFFF),
// it returns [DefaultLangSys] (which may be empty)
func (sc Script) GetLangSys(index uint16) LangSys {
	if int(index) >= len(sc.LangSys) {
		if sc.DefaultLangSys != nil {
			return *sc.DefaultLangSys
		}
		return LangSys{}
	}
	return sc.LangSys[index]
}

// --------------------------------------- gsub ---------------------------------------

func (d SingleSubstData1) Cov() Coverage { return d.Coverage }
func (d SingleSubstData2) Cov() Coverage { return d.Coverage }

func (cs ContextualSubs1) Cov() Coverage { return cs.coverage }
func (cs ContextualSubs2) Cov() Coverage { return cs.coverage }
func (cs ContextualSubs3) Cov() Coverage {
	if len(cs.Coverages) == 0 { // return an empty, valid Coverage
		return Coverage1{}
	}
	return cs.Coverages[0]
}

func (cc ChainedContextualSubs1) Cov() Coverage { return cc.coverage }
func (cc ChainedContextualSubs2) Cov() Coverage { return cc.coverage }
func (cc ChainedContextualSubs3) Cov() Coverage {
	if len(cc.InputCoverages) == 0 { // return an empty, valid Coverage
		return Coverage1{}
	}
	return cc.InputCoverages[0]
}

func (lk SingleSubs) Cov() Coverage             { return lk.Data.Cov() }
func (lk MultipleSubs) Cov() Coverage           { return lk.Coverage }
func (lk AlternateSubs) Cov() Coverage          { return lk.Coverage }
func (lk LigatureSubs) Cov() Coverage           { return lk.Coverage }
func (lk ContextualSubs) Cov() Coverage         { return lk.Data.Cov() }
func (lk ChainedContextualSubs) Cov() Coverage  { return lk.Data.Cov() }
func (lk ExtensionSubs) Cov() Coverage          { return nil } // not used anyway
func (lk ReverseChainSingleSubs) Cov() Coverage { return lk.coverage }

// --------------------------------------- gpos ---------------------------------------

func (d SinglePosData1) Cov() Coverage { return d.coverage }
func (d SinglePosData2) Cov() Coverage { return d.coverage }

func (d PairPosData1) Cov() Coverage { return d.coverage }
func (d PairPosData2) Cov() Coverage { return d.coverage }

func (cs ContextualPos1) Cov() Coverage { return cs.coverage }
func (cs ContextualPos2) Cov() Coverage { return cs.coverage }
func (cs ContextualPos3) Cov() Coverage {
	if len(cs.Coverages) == 0 { // return an empty, valid Coverage
		return Coverage1{}
	}
	return cs.Coverages[0]
}

func (cc ChainedContextualPos1) Cov() Coverage { return cc.coverage }
func (cc ChainedContextualPos2) Cov() Coverage { return cc.coverage }
func (cc ChainedContextualPos3) Cov() Coverage {
	if len(cc.InputCoverages) == 0 { // return an empty, valid Coverage
		return Coverage1{}
	}
	return cc.InputCoverages[0]
}

func (lk SinglePos) Cov() Coverage            { return lk.Data.Cov() }
func (lk PairPos) Cov() Coverage              { return lk.Data.Cov() }
func (lk CursivePos) Cov() Coverage           { return lk.coverage }
func (lk MarkBasePos) Cov() Coverage          { return lk.markCoverage }
func (lk MarkLigPos) Cov() Coverage           { return lk.MarkCoverage }
func (lk MarkMarkPos) Cov() Coverage          { return lk.Mark1Coverage }
func (lk ContextualPos) Cov() Coverage        { return lk.Data.Cov() }
func (lk ChainedContextualPos) Cov() Coverage { return lk.Data.Cov() }
func (lk ExtensionPos) Cov() Coverage         { return nil } // not used anyway

// FindGlyph performs a binary search in the list, returning the record for `secondGlyph`,
// or `nil` if not found.
func (ps PairSet) FindGlyph(secondGlyph GlyphID) *PairValueRecord {
	low, high := 0, len(ps.PairValueRecords)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := ps.PairValueRecords[mid].SecondGlyph
		if secondGlyph < p {
			high = mid
		} else if secondGlyph > p {
			low = mid + 1
		} else {
			return &ps.PairValueRecords[mid]
		}
	}
	return nil
}

// GetDelta returns the hint for the given `ppem`, scaled by `scale`.
// It returns 0 for out of range `ppem` values.
func (dev DeviceHinting) GetDelta(ppem uint16, scale int32) int32 {
	if ppem == 0 {
		return 0
	}

	if ppem < dev.StartSize || ppem > dev.EndSize {
		return 0
	}

	pixels := dev.Values[ppem-dev.StartSize]

	return int32(pixels) * (scale / int32(ppem))
}

// -------------------------------------- gdef --------------------------------------

// GlyphProps is a 16-bit integer where the lower 8-bit have bits representing
// glyph class, and high 8-bit the mark attachment type (if any).
type GlyphProps = uint16

const (
	GPBaseGlyph GlyphProps = 1 << (iota + 1)
	GPLigature
	GPMark
)

// GlyphProps return a summary of the glyph properties.
func (gd *GDEF) GlyphProps(glyph GlyphID) GlyphProps {
	klass, _ := gd.GlyphClassDef.Class(glyph)
	switch klass {
	case 1:
		return GPBaseGlyph
	case 2:
		return GPLigature
	case 3:
		var klass uint16 // it is actually a byte
		if gd.MarkAttachClass != nil {
			klass, _ = gd.MarkAttachClass.Class(glyph)
		}
		return GPMark | GlyphProps(klass)<<8
	default:
		return 0
	}
}

// -------------------------------------- var --------------------------------------

// GetDelta uses the variation [store] and the selected instance coordinates [coords]
// to compute the value at [index].
func (store ItemVarStore) GetDelta(index VariationStoreIndex, coords []float32) float32 {
	if int(index.DeltaSetOuter) >= len(store.ItemVariationDatas) {
		return 0
	}
	varData := store.ItemVariationDatas[index.DeltaSetOuter]
	if int(index.DeltaSetInner) >= len(varData.DeltaSets) {
		return 0
	}
	deltaSet := varData.DeltaSets[index.DeltaSetInner]
	var delta float32
	for i, regionIndex := range varData.RegionIndexes {
		region := store.VariationRegionList.VariationRegions[regionIndex].RegionAxes
		v := float32(1)
		for axis, coord := range coords {
			factor := region[axis].evaluate(coord)
			v *= factor
		}
		delta += float32(deltaSet[i]) * v
	}
	return delta
}
