// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package font

import "github.com/go-text/typesetting/opentype/tables"

// shared between GSUB and GPOS
type Layout struct {
	Scripts           []Script
	Features          []Feature
	FeatureVariations []tables.FeatureVariationRecord
}

func newLayout(table tables.Layout) Layout {
	out := Layout{
		Scripts:  make([]Script, len(table.ScriptList.Scripts)),
		Features: make([]Feature, len(table.FeatureList.Features)),
	}
	for i, s := range table.ScriptList.Scripts {
		out.Scripts[i] = Script{
			Script: s,
			Tag:    table.ScriptList.Records[i].Tag,
		}
	}
	for i, f := range table.FeatureList.Features {
		out.Features[i] = Feature{
			Feature: f,
			Tag:     table.FeatureList.Records[i].Tag,
		}
	}
	if table.FeatureVariations != nil {
		out.FeatureVariations = table.FeatureVariations.FeatureVariationRecords
	}
	return out
}

type Script struct {
	tables.Script
	Tag Tag
}

type Feature struct {
	tables.Feature
	Tag Tag
}

// FindScript looks for [script] and return its index into the Scripts slice,
// or -1 if the tag is not found.
func (la *Layout) FindScript(script Tag) int {
	// Scripts is sorted: binary search
	low, high := 0, len(la.Scripts)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := la.Scripts[mid].Tag
		if script < p {
			high = mid
		} else if script > p {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

// FindVariationIndex returns the first feature variation matching
// the specified variation coordinates, as an index in the
// `FeatureVariations` field.
// It returns `-1` if not found.
func (la *Layout) FindVariationIndex(coords []float32) int {
	for i, record := range la.FeatureVariations {
		if evaluateVarRec(record, coords) {
			return i
		}
	}
	return -1
}

// returns `true` if the feature is concerned by the `coords`
func evaluateVarRec(fv tables.FeatureVariationRecord, coords []float32) bool {
	for _, c := range fv.ConditionSet.Conditions {
		if !evaluateCondition(c, coords) {
			return false
		}
	}
	return true
}

// returns `true` if `coords` match the condition `c`
func evaluateCondition(c tables.ConditionFormat1, coords []float32) bool {
	var coord float32
	if int(c.AxisIndex) < len(coords) {
		coord = coords[c.AxisIndex]
	}
	return c.FilterRangeMinValue <= coord && coord <= c.FilterRangeMaxValue
}

// FindFeatureIndex fetches the index for a given feature tag in the GSUB or GPOS table.
// Returns false if not found
func (la *Layout) FindFeatureIndex(featureTag Tag) (uint16, bool) {
	for i, feat := range la.Features { // i fits in uint16
		if featureTag == feat.Tag {
			return uint16(i), true
		}
	}
	return 0, false
}

// ---------------------------------- GSUB ----------------------------------

type GSUB struct {
	Layout
	Lookups []GSUBLookup
}

type LookupOptions struct {
	// Lookup qualifiers.
	Flag uint16
	// Index (base 0) into GDEF mark glyph sets structure,
	// meaningfull only if UseMarkFilteringSet is set.
	MarkFilteringSet uint16
}

const UseMarkFilteringSet = 1 << 4

// Props returns a 32-bit integer where the lower 16-bit is `Flag` and
// the higher 16-bit is `MarkFilteringSet` if the lookup uses one.
func (lo LookupOptions) Props() uint32 {
	flag := uint32(lo.Flag)
	if lo.Flag&UseMarkFilteringSet != 0 {
		flag |= uint32(lo.MarkFilteringSet) << 16
	}
	return flag
}

type GSUBLookup struct {
	LookupOptions
	Subtables []tables.GSUBLookup
}

func newGSUB(table tables.Layout) (GSUB, error) {
	out := GSUB{
		Layout:  newLayout(table),
		Lookups: make([]GSUBLookup, len(table.LookupList.Lookups)),
	}
	for i, lk := range table.LookupList.Lookups {
		subtables, err := lk.AsGSUBLookups()
		if err != nil {
			return GSUB{}, err
		}
		for j, subtable := range subtables {
			// start by resolving extension
			if ext, isExt := subtable.(tables.ExtensionSubs); isExt {
				subtables[j], err = ext.Resolve()
				if err != nil {
					return GSUB{}, err
				}
			}

			// sanitize each lookup
			switch subtable := subtable.(type) {
			case tables.MultipleSubs:
				err = subtable.Sanitize()
			case tables.LigatureSubs:
				err = subtable.Sanitize()
			case tables.ContextualSubs:
				err = subtable.Sanitize(uint16(len(out.Lookups)))
			case tables.ReverseChainSingleSubs:
				err = subtable.Sanitize()
			}
			if err != nil {
				return GSUB{}, err
			}
		}
		out.Lookups[i] = GSUBLookup{
			LookupOptions: LookupOptions{
				Flag:             lk.LookupFlag,
				MarkFilteringSet: lk.MarkFilteringSet,
			},
			Subtables: subtables,
		}
	}
	return out, nil
}

type GPOS struct {
	Layout
	Lookups []GPOSLookup
}

type GPOSLookup struct {
	LookupOptions
	Subtables []tables.GPOSLookup
}

func newGPOS(table tables.Layout) (GPOS, error) {
	out := GPOS{
		Layout:  newLayout(table),
		Lookups: make([]GPOSLookup, len(table.LookupList.Lookups)),
	}
	for i, lk := range table.LookupList.Lookups {
		subtables, err := lk.AsGPOSLookups()
		if err != nil {
			return GPOS{}, err
		}
		for j, subtable := range subtables {
			// start by resolving extension
			if ext, isExt := subtable.(tables.ExtensionPos); isExt {
				subtables[j], err = ext.Resolve()
				if err != nil {
					return GPOS{}, err
				}
			}

			// sanitize each lookup
			switch subtable := subtable.(type) {
			case tables.SinglePos:
				err = subtable.Sanitize()
			case tables.PairPos:
				err = subtable.Sanitize()
			case tables.MarkBasePos:
				err = subtable.Sanitize()
			case tables.MarkLigPos:
				err = subtable.Sanitize()
			case tables.ContextualPos:
				err = subtable.Sanitize(uint16(len(out.Lookups)))
			}
			if err != nil {
				return GPOS{}, err
			}
		}
		out.Lookups[i] = GPOSLookup{
			LookupOptions: LookupOptions{
				Flag:             lk.LookupFlag,
				MarkFilteringSet: lk.MarkFilteringSet,
			},
			Subtables: subtables,
		}
	}
	return out, nil
}
