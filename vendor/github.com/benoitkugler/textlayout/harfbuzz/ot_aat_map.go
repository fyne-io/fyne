package harfbuzz

import (
	"sort"

	tt "github.com/benoitkugler/textlayout/fonts/truetype"
)

// ported from harfbuzz/src/hb-aat-map.cc, hb-att-map.hh Copyright © 2018  Google, Inc. Behdad Esfahbod

type aatMap struct {
	chainFlags []GlyphMask
}

type aatFeatureInfo struct {
	type_       aatLayoutFeatureType
	setting     aatLayoutFeatureSelector
	isExclusive bool
}

func (fi aatFeatureInfo) key() uint32 {
	return uint32(fi.type_)<<16 | uint32(fi.setting)
}

const selMask = ^aatLayoutFeatureSelector(1)

func cmpAATFeatureInfo(a, b aatFeatureInfo) bool {
	if a.type_ != b.type_ {
		return a.type_ < b.type_
	}
	if !a.isExclusive && (a.setting&selMask) != (b.setting&selMask) {
		return a.setting < b.setting
	}
	return false
}

type aatMapBuilder struct {
	tables   *tt.LayoutTables
	features []aatFeatureInfo // sorted by (type_, setting) after compilation
}

// binary search into `features`, comparing type_ and setting only
func (mb *aatMapBuilder) hasFeature(info aatFeatureInfo) bool {
	key := info.key()
	for i, j := 0, len(mb.features); i < j; {
		h := i + (j-i)/2
		entry := mb.features[h].key()
		if key < entry {
			j = h
		} else if entry < key {
			i = h + 1
		} else {
			return true
		}
	}
	return false
}

func (mb *aatMapBuilder) compileMap(map_ *aatMap) {
	morx := mb.tables.Morx
	for _, chain := range morx {
		map_.chainFlags = append(map_.chainFlags, mb.compileMorxFlag(chain))
	}

	// TODO: for now we dont support deprecated mort table
	// mort := mapper.face.table.mort
	// if mort.has_data() {
	// 	mort.compile_flags(mapper, map_)
	// 	return
	// }
}

func (mb *aatMapBuilder) compileMorxFlag(chain tt.MorxChain) GlyphMask {
	flags := chain.DefaultFlags

	for _, feature := range chain.Features {
		type_, setting := feature.Type, feature.Setting

	retry:
		// Check whether this type_/setting pair was requested in the map, and if so, apply its flags.
		// (The search here only looks at the type_ and setting fields of feature_info_t.)
		info := aatFeatureInfo{type_, setting, false}
		if mb.hasFeature(info) {
			flags &= feature.DisableFlags
			flags |= feature.EnableFlags
		} else if type_ == aatLayoutFeatureTypeLetterCase && setting == aatLayoutFeatureSelectorSmallCaps {
			/* Deprecated. https://github.com/harfbuzz/harfbuzz/issues/1342 */
			type_ = aatLayoutFeatureTypeLowerCase
			setting = aatLayoutFeatureSelectorLowerCaseSmallCaps
			goto retry
		}
	}
	return flags
}

func (mb *aatMapBuilder) addFeature(tag tt.Tag, value uint32) {
	feat := mb.tables.Feat
	if len(feat) == 0 {
		return
	}

	if tag == tt.NewTag('a', 'a', 'l', 't') {
		if fn := feat.GetFeature(aatLayoutFeatureTypeCharacterAlternatives); fn == nil || len(fn.Settings) == 0 {
			return
		}
		info := aatFeatureInfo{
			type_:       aatLayoutFeatureTypeCharacterAlternatives,
			setting:     aatLayoutFeatureSelector(value),
			isExclusive: true,
		}
		mb.features = append(mb.features, info)
		return
	}

	mapping := aatLayoutFindFeatureMapping(tag)
	if mapping == nil {
		return
	}

	feature := feat.GetFeature(mapping.aatFeatureType)
	if feature == nil || len(feature.Settings) == 0 {
		/* Special case: compileMorxFlag() will fall back to the deprecated version of
		 * small-caps if necessary, so we need to check for that possibility.
		 * https://github.com/harfbuzz/harfbuzz/issues/2307 */
		if mapping.aatFeatureType == aatLayoutFeatureTypeLowerCase &&
			mapping.selectorToEnable == aatLayoutFeatureSelectorLowerCaseSmallCaps {
			feature = feat.GetFeature(aatLayoutFeatureTypeLetterCase)
			if feature == nil || len(feature.Settings) == 0 {
				return
			}
		} else {
			return
		}
	}

	var info aatFeatureInfo
	info.type_ = mapping.aatFeatureType
	if value != 0 {
		info.setting = mapping.selectorToEnable
	} else {
		info.setting = mapping.selectorToDisable
	}
	info.isExclusive = feature.IsExclusive()
	mb.features = append(mb.features, info)
}

func (mb *aatMapBuilder) compile(m *aatMap) {
	// sort features and merge duplicates
	if len(mb.features) != 0 {
		sort.SliceStable(mb.features, func(i, j int) bool {
			return cmpAATFeatureInfo(mb.features[i], mb.features[j])
		})
		j := 0
		for i := 1; i < len(mb.features); i++ {
			/* Nonexclusive feature selectors come in even/odd pairs to turn a setting on/off
			* respectively, so we mask out the low-order bit when checking for "duplicates"
			* (selectors referring to the same feature setting) here. */
			if mb.features[i].type_ != mb.features[j].type_ ||
				(!mb.features[i].isExclusive && ((mb.features[i].setting & selMask) != (mb.features[j].setting & selMask))) {
				j++
				mb.features[j] = mb.features[i]
			}
		}
		mb.features = mb.features[:j+1]
	}

	mb.compileMap(m)
}
