package harfbuzz

import (
	"fmt"
	"math"
	"math/bits"
	"sort"

	tt "github.com/benoitkugler/textlayout/fonts/truetype"
)

// ported from harfbuzz/src/hb-ot-map.cc, hb-ot-map.hh Copyright © 2009,2010  Red Hat, Inc. 2010,2011,2013  Google, Inc. Behdad Esfahbod

var (
	otTagGSUB = tt.TagGsub
	otTagGPOS = tt.TagGpos
	tableTags = [2]tt.Tag{otTagGSUB, otTagGPOS}
)

type otMapFeatureFlags uint8

const (
	ffGLOBAL              otMapFeatureFlags = 1 << iota /* Feature applies to all characters; results in no mask allocated for it. */
	ffHasFallback                                       /* Has fallback implementation, so include mask bit even if feature not found. */
	ffManualZWNJ                                        /* Don't skip over ZWNJ when matching **context**. */
	ffManualZWJ                                         /* Don't skip over ZWJ when matching **input**. */
	ffGlobalSearch                                      /* If feature not found in LangSys, look for it in global feature list and pick one. */
	ffRandom                                            /* Randomly select a glyph from an AlternateSubstFormat1 subtable. */
	ffNone                otMapFeatureFlags = 0
	ffManualJoiners                         = ffManualZWNJ | ffManualZWJ
	ffGlobalManualJoiners                   = ffGLOBAL | ffManualJoiners
	ffGlobalHasFallback                     = ffGLOBAL | ffHasFallback
)

const (
	otMapMaxBits  = 8
	otMapMaxValue = (1 << otMapMaxBits) - 1
)

type otMapFeature struct {
	tag   tt.Tag
	flags otMapFeatureFlags
}

type featureInfo struct {
	Tag tt.Tag
	// seq           int /* sequence#, used for stable sorting only */
	maxValue     uint32
	flags        otMapFeatureFlags
	defaultValue uint32 /* for non-global features, what should the unset glyphs take */
	stage        [2]int /* GSUB/GPOS */
}

type stageInfo struct {
	pauseFunc pauseFunc
	index     int
}

type otMapBuilder struct {
	tables        *tt.LayoutTables
	props         SegmentProperties
	stages        [2][]stageInfo
	featureInfos  []featureInfo
	scriptIndex   [2]int
	languageIndex [2]int
	currentStage  [2]int
	chosenScript  [2]tt.Tag
	foundScript   [2]bool
}

//  void hb_ot_map_t::collect_lookups (uint tableIndex, hb_set_t *lookups_out) const
//  {
//    for (uint i = 0; i < lookups[tableIndex].length; i++)
// 	 lookups_out.add (lookups[tableIndex][i].index);
//  }

func newOtMapBuilder(tables *tt.LayoutTables, props SegmentProperties) otMapBuilder {
	var out otMapBuilder

	out.tables = tables
	out.props = props

	/* Fetch script/language indices for GSUB/GPOS.  We need these later to skip
	* features not available in either table and not waste precious bits for them. */
	scriptTags, languageTags := NewOTTagsFromScriptAndLanguage(props.Script, props.Language)

	out.scriptIndex[0], out.chosenScript[0], out.foundScript[0] = SelectScript(&tables.GSUB.TableLayout, scriptTags)
	out.languageIndex[0], _ = SelectLanguage(&tables.GSUB.TableLayout, out.scriptIndex[0], languageTags)

	out.scriptIndex[1], out.chosenScript[1], out.foundScript[1] = SelectScript(&tables.GPOS.TableLayout, scriptTags)
	out.languageIndex[1], _ = SelectLanguage(&tables.GPOS.TableLayout, out.scriptIndex[1], languageTags)

	return out
}

func (mb *otMapBuilder) addFeatureExt(tag tt.Tag, flags otMapFeatureFlags, value uint32) {
	var info featureInfo
	info.Tag = tag
	info.maxValue = value
	info.flags = flags
	if (flags & ffGLOBAL) != 0 {
		info.defaultValue = value
	}
	info.stage = mb.currentStage

	mb.featureInfos = append(mb.featureInfos, info)
}

type pauseFunc func(plan *otShapePlan, font *Font, buffer *Buffer)

func (mb *otMapBuilder) addPause(tableIndex int, fn pauseFunc) {
	s := stageInfo{
		index:     mb.currentStage[tableIndex],
		pauseFunc: fn,
	}
	mb.stages[tableIndex] = append(mb.stages[tableIndex], s)
	mb.currentStage[tableIndex]++
}

func (mb *otMapBuilder) addGSUBPause(fn pauseFunc) { mb.addPause(0, fn) }
func (mb *otMapBuilder) addGPOSPause(fn pauseFunc) { mb.addPause(1, fn) }

func (mb *otMapBuilder) enableFeatureExt(tag tt.Tag, flags otMapFeatureFlags, value uint32) {
	mb.addFeatureExt(tag, ffGLOBAL|flags, value)
}
func (mb *otMapBuilder) enableFeature(tag tt.Tag)  { mb.enableFeatureExt(tag, ffNone, 1) }
func (mb *otMapBuilder) addFeature(tag tt.Tag)     { mb.addFeatureExt(tag, ffNone, 1) }
func (mb *otMapBuilder) disableFeature(tag tt.Tag) { mb.addFeatureExt(tag, ffGLOBAL, 0) }

func (mb *otMapBuilder) compile(m *otMap, key otShapePlanKey) {
	const globalBitShift = 8*4 - 1
	const globalBitMask = 1 << globalBitShift

	m.globalMask = globalBitMask

	var (
		requiredFeatureIndex [2]uint16 // HB_OT_LAYOUT_NO_FEATURE_INDEX for empty
		requiredFeatureTag   [2]tt.Tag
		/* We default to applying required feature in stage 0. If the required
		* feature has a tag that is known to the shaper, we apply the required feature
		* in the stage for that tag. */
		requiredFeatureStage [2]int
	)

	gsub, gpos := mb.tables.GSUB, mb.tables.GPOS
	tables := [2]*tt.TableLayout{&gsub.TableLayout, &gpos.TableLayout}

	m.chosenScript = mb.chosenScript
	m.foundScript = mb.foundScript
	requiredFeatureIndex[0], requiredFeatureTag[0] = getRequiredFeature(tables[0], mb.scriptIndex[0], mb.languageIndex[0])
	requiredFeatureIndex[1], requiredFeatureTag[1] = getRequiredFeature(tables[1], mb.scriptIndex[1], mb.languageIndex[1])

	// sort features and merge duplicates
	if len(mb.featureInfos) != 0 {
		sort.SliceStable(mb.featureInfos, func(i, j int) bool {
			return mb.featureInfos[i].Tag < mb.featureInfos[j].Tag
		})
		j := 0
		for i, feat := range mb.featureInfos {
			if i == 0 {
				continue
			}
			if feat.Tag != mb.featureInfos[j].Tag {
				j++
				mb.featureInfos[j] = feat
				continue
			}
			if feat.flags&ffGLOBAL != 0 {
				mb.featureInfos[j].flags |= ffGLOBAL
				mb.featureInfos[j].maxValue = feat.maxValue
				mb.featureInfos[j].defaultValue = feat.defaultValue
			} else {
				if mb.featureInfos[j].flags&ffGLOBAL != 0 {
					mb.featureInfos[j].flags ^= ffGLOBAL
				}
				mb.featureInfos[j].maxValue = max32(mb.featureInfos[j].maxValue, feat.maxValue)
				// inherit default_value from j
			}
			mb.featureInfos[j].flags |= (feat.flags & ffHasFallback)
			mb.featureInfos[j].stage[0] = min(mb.featureInfos[j].stage[0], feat.stage[0])
			mb.featureInfos[j].stage[1] = min(mb.featureInfos[j].stage[1], feat.stage[1])
		}
		mb.featureInfos = mb.featureInfos[0 : j+1]
	}

	// allocate bits now
	nextBit := bits.OnesCount32(glyphFlagDefined) + 1

	for _, info := range mb.featureInfos {

		bitsNeeded := 0

		if (info.flags&ffGLOBAL) != 0 && info.maxValue == 1 {
			// uses the global bit
			bitsNeeded = 0
		} else {
			// limit bits per feature.
			bitsNeeded = min(otMapMaxBits, bitStorage(info.maxValue))
		}

		if info.maxValue == 0 || nextBit+bitsNeeded >= globalBitShift {
			continue // feature disabled, or not enough bits.
		}

		var (
			found        = false
			featureIndex [2]uint16
		)
		for tableIndex, table := range tables {
			if requiredFeatureTag[tableIndex] == info.Tag {
				requiredFeatureStage[tableIndex] = info.stage[tableIndex]
			}
			featureIndex[tableIndex] = FindFeatureForLang(table, mb.scriptIndex[tableIndex], mb.languageIndex[tableIndex], info.Tag)
			found = found || featureIndex[tableIndex] != NoFeatureIndex
		}
		if !found && (info.flags&ffGlobalSearch) != 0 {
			for tableIndex, table := range tables {
				featureIndex[tableIndex] = findFeature(table, info.Tag)
				found = found || featureIndex[tableIndex] != NoFeatureIndex
			}
		}
		if !found && info.flags&ffHasFallback == 0 {
			continue
		}

		var map_ featureMap
		map_.tag = info.Tag
		map_.index = featureIndex
		map_.stage = info.stage
		map_.autoZWNJ = info.flags&ffManualZWNJ == 0
		map_.autoZWJ = info.flags&ffManualZWJ == 0
		map_.random = info.flags&ffRandom != 0
		if (info.flags&ffGLOBAL) != 0 && info.maxValue == 1 {
			// uses the global bit
			map_.shift = globalBitShift
			map_.mask = globalBitMask
		} else {
			map_.shift = nextBit
			map_.mask = (1 << (nextBit + bitsNeeded)) - (1 << nextBit)
			nextBit += bitsNeeded
			m.globalMask |= (info.defaultValue << map_.shift) & map_.mask
		}
		map_.mask1 = (1 << map_.shift) & map_.mask
		map_.needsFallback = !found

		if debugMode >= 1 {
			fmt.Printf("\tMAP - adding feature %s (%d) for stage %v\n", info.Tag, info.Tag, info.stage)
		}

		m.features = append(m.features, map_)
	}
	mb.featureInfos = mb.featureInfos[:0] // done with these

	mb.addGSUBPause(nil)
	mb.addGPOSPause(nil)

	// collect lookup indices for features
	for tableIndex, table := range tables {
		stageIndex := 0
		lastNumLookups := 0
		for stage := 0; stage < mb.currentStage[tableIndex]; stage++ {
			if requiredFeatureIndex[tableIndex] != NoFeatureIndex &&
				requiredFeatureStage[tableIndex] == stage {
				m.addLookups(table, tableIndex, requiredFeatureIndex[tableIndex],
					key[tableIndex], globalBitMask, true, true, false)
			}

			for _, feat := range m.features {
				if feat.stage[tableIndex] == stage {
					m.addLookups(table, tableIndex,
						feat.index[tableIndex],
						key[tableIndex],
						feat.mask,
						feat.autoZWNJ,
						feat.autoZWJ,
						feat.random)
				}
			}
			// sort lookups and merge duplicates

			if ls := m.lookups[tableIndex]; lastNumLookups < len(ls) {
				view := ls[lastNumLookups:]
				sort.Slice(view, func(i, j int) bool { return view[i].index < view[j].index })

				j := lastNumLookups
				for i := j + 1; i < len(ls); i++ {
					if ls[i].index != ls[j].index {
						j++
						ls[j] = ls[i]
					} else {
						ls[j].mask |= ls[i].mask
						ls[j].autoZWNJ = ls[j].autoZWNJ && ls[i].autoZWNJ
						ls[j].autoZWJ = ls[j].autoZWJ && ls[i].autoZWJ
					}
				}
				m.lookups[tableIndex] = m.lookups[tableIndex][:j+1]
			}

			lastNumLookups = len(m.lookups[tableIndex])

			if stageIndex < len(mb.stages[tableIndex]) && mb.stages[tableIndex][stageIndex].index == stage {
				sm := stageMap{
					lastLookup: lastNumLookups,
					pauseFunc:  mb.stages[tableIndex][stageIndex].pauseFunc,
				}
				m.stages[tableIndex] = append(m.stages[tableIndex], sm)
				stageIndex++
			}
		}
	}
}

type featureMap struct {
	tag           tt.Tag    /* should be first for our bsearch to work */
	index         [2]uint16 /* GSUB/GPOS */
	stage         [2]int    /* GSUB/GPOS */
	shift         int
	mask          GlyphMask
	mask1         GlyphMask /* mask for value=1, for quick access */
	needsFallback bool      // = 1;
	autoZWNJ      bool      // = 1;
	autoZWJ       bool      // = 1;
	random        bool      // = 1;

	// int cmp (const hb_tag_t tag_) const
	// { return tag_ < tag ? -1 : tag_ > tag ? 1 : 0; }
}

func bsearchFeature(features []featureMap, tag tt.Tag) *featureMap {
	low, high := 0, len(features)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := features[mid].tag
		if tag < p {
			high = mid
		} else if tag > p {
			low = mid + 1
		} else {
			return &features[mid]
		}
	}
	return nil
}

type lookupMap struct {
	index    uint16
	autoZWNJ bool // = 1;
	autoZWJ  bool // = 1;
	random   bool // = 1;
	mask     GlyphMask

	// HB_INTERNAL static int cmp (const void *pa, const void *pb)
	// {
	//   const lookup_map_t *a = (const lookup_map_t *) pa;
	//   const lookup_map_t *b = (const lookup_map_t *) pb;
	//   return a.index < b.index ? -1 : a.index > b.index ? 1 : 0;
	// }
}

type stageMap struct {
	pauseFunc  pauseFunc
	lastLookup int
}

type otMap struct {
	lookups      [2][]lookupMap
	stages       [2][]stageMap
	features     []featureMap // sorted
	chosenScript [2]tt.Tag
	globalMask   GlyphMask
	foundScript  [2]bool
}

//   friend struct hb_ot_map_builder_t;

func (m *otMap) needsFallback(featureTag tt.Tag) bool {
	if ma := bsearchFeature(m.features, featureTag); ma != nil {
		return ma.needsFallback
	}
	return false
}

func (m *otMap) getMask(featureTag tt.Tag) (GlyphMask, int) {
	if ma := bsearchFeature(m.features, featureTag); ma != nil {
		return ma.mask, ma.shift
	}
	return 0, 0
}

func (m *otMap) getMask1(featureTag tt.Tag) GlyphMask {
	if ma := bsearchFeature(m.features, featureTag); ma != nil {
		return ma.mask1
	}
	return 0
}

func (m *otMap) getFeatureIndex(tableIndex int, featureTag tt.Tag) uint16 {
	if ma := bsearchFeature(m.features, featureTag); ma != nil {
		return ma.index[tableIndex]
	}
	return NoFeatureIndex
}

func (m *otMap) getFeatureStage(tableIndex int, featureTag tt.Tag) int {
	if ma := bsearchFeature(m.features, featureTag); ma != nil {
		return ma.stage[tableIndex]
	}
	return math.MaxInt32
}

func (m *otMap) getStageLookups(tableIndex, stage int) []lookupMap {
	if stage > len(m.stages[tableIndex]) {
		return nil
	}
	start, end := 0, len(m.lookups[tableIndex])
	if stage != 0 {
		start = m.stages[tableIndex][stage-1].lastLookup
	}
	if stage < len(m.stages[tableIndex]) {
		end = m.stages[tableIndex][stage].lastLookup
	}
	return m.lookups[tableIndex][start:end]
}

func (m *otMap) addLookups(table *tt.TableLayout, tableIndex int, featureIndex uint16, variationsIndex int,
	mask GlyphMask, autoZwnj, autoZwj, random bool) {
	lookupIndices := getFeatureLookupsWithVar(table, featureIndex, variationsIndex)
	for _, lookupInd := range lookupIndices {
		lookup := lookupMap{
			mask:     mask,
			index:    lookupInd,
			autoZWNJ: autoZwnj,
			autoZWJ:  autoZwj,
			random:   random,
		}
		m.lookups[tableIndex] = append(m.lookups[tableIndex], lookup)
	}
}

// apply the GSUB table
func (m *otMap) substitute(plan *otShapePlan, font *Font, buffer *Buffer) {
	if debugMode >= 1 {
		fmt.Println("SUBSTITUTE - start table GSUB")
	}

	proxy := otProxy{otProxyMeta: proxyGSUB, accels: font.gsubAccels}
	m.apply(proxy, plan, font, buffer)

	if debugMode >= 1 {
		fmt.Println("SUBSTITUTE - end table GSUB")
	}
}

// apply the GPOS table
func (m *otMap) position(plan *otShapePlan, font *Font, buffer *Buffer) {
	if debugMode >= 1 {
		fmt.Println("POSITION - start table GPOS")
	}

	proxy := otProxy{otProxyMeta: proxyGPOS, accels: font.gposAccels}
	m.apply(proxy, plan, font, buffer)

	if debugMode >= 1 {
		fmt.Println("POSITION - end table GPOS")
	}
}

func (m *otMap) apply(proxy otProxy, plan *otShapePlan, font *Font, buffer *Buffer) {
	tableIndex := proxy.tableIndex
	i := 0
	c := newOtApplyContext(tableIndex, font, buffer)
	c.recurseFunc = proxy.recurseFunc

	for stageI, stage := range m.stages[tableIndex] {

		if debugMode >= 2 {
			fmt.Printf("\tAPPLY - stage %d\n", stageI)
		}

		for ; i < stage.lastLookup; i++ {
			lookupIndex := m.lookups[tableIndex][i].index

			if debugMode >= 1 {
				fmt.Printf("\t\tLookup %d start\n", lookupIndex)
			}

			c.lookupIndex = lookupIndex
			c.setLookupMask(m.lookups[tableIndex][i].mask)
			c.setAutoZWJ(m.lookups[tableIndex][i].autoZWJ)
			c.setAutoZWNJ(m.lookups[tableIndex][i].autoZWNJ)
			c.random = m.lookups[tableIndex][i].random

			// pathological cases
			if len(c.buffer.Info) > c.buffer.maxLen {
				return
			}
			c.applyString(proxy.otProxyMeta, &proxy.accels[lookupIndex])

			if debugMode >= 1 {
				fmt.Println("\t\tLookup end")
				fmt.Println(c.buffer.Info)
			}

		}

		if stage.pauseFunc != nil {
			if debugMode >= 1 {
				fmt.Println("\t\tExecuting pause function")
			}

			stage.pauseFunc(plan, font, buffer)
		}
	}
}
