package harfbuzz

import (
	"math"

	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/opentype/tables"
)

/* ported from harfbuzz/src/hb-buffer.hh and hb-buffer.h
 * Copyright © 1998-2004  David Turner and Werner Lemberg
 * Copyright © 2004,2007,2009,2010  Red Hat, Inc.
 * Copyright © 2011,2012  Google, Inc.
 * Red Hat Author(s): Owen Taylor, Behdad Esfahbod
 * Google Author(s): Behdad Esfahbod */

const (
	// The following are used internally; not derived from GDEF.
	substituted tables.GlyphProps = 1 << (iota + 4)
	ligated
	multiplied

	preserve = substituted | ligated | multiplied
)

type bufferScratchFlags uint32

const (
	bsfHasNonASCII bufferScratchFlags = 1 << iota
	bsfHasDefaultIgnorables
	bsfHasSpaceFallback
	bsfHasGPOSAttachment
	bsfHasUnsafeToBreak
	bsfHasCGJ
	bsfDefault bufferScratchFlags = 0x00000000

	// reserved for complex shapers' internal use.
	bsfComplex0 bufferScratchFlags = 0x01000000
	bsfComplex1 bufferScratchFlags = 0x02000000
	bsfComplex2 bufferScratchFlags = 0x04000000
	bsfComplex3 bufferScratchFlags = 0x08000000
)

// maximum length of additional context added outside
// input text
const contextLength = 5

const (
	maxOpsDefault = 0x1FFFFFFF
	maxLenDefault = 0x3FFFFFFF
)

// Buffer is the main structure holding the input text segment and its properties before shaping,
// and output glyphs and their information after shaping.
type Buffer struct {
	// Info is used as internal storage during the shaping,
	// and also exposes the result: the glyph to display
	// and its original Cluster value.
	Info []GlyphInfo

	// Pos gives the position of the glyphs resulting from the shapping
	// It has the same length has `Info`.
	Pos []GlyphPosition

	// Text before / after the main buffer contents, ordered outward !
	// Index 0 is for "pre-Context", 1 for "post-Context".
	context [2][]rune

	// temporary storage, usully used the following way:
	// 	- truncate the slice
	//	- walk through Info and append glyphs to outInfo
	//	- swap back Into and outInfo with 'swapBuffers'
	outInfo []GlyphInfo

	// Props is required to correctly interpret the input runes.
	Props SegmentProperties
	// Glyph that replaces invisible characters in
	// the shaping result. If set to zero (default), the glyph for the
	// U+0020 SPACE character is used. Otherwise, this value is used
	// verbatim.
	Invisible GID

	// Glyph that replaces characters not found in the font during shaping.
	// The not-found glyph defaults to zero, sometimes knows as the
	// ".notdef" glyph.
	NotFound GID

	// Information about how the text in the buffer should be treated.
	Flags ShappingOptions
	// Precise the cluster handling behavior.
	ClusterLevel ClusterLevel

	// some pathological cases can be constructed
	// (for example with GSUB tables), where the size of the buffer
	// grows out of bounds
	// these bounds avoid such cases, which should never happen with
	// decent font files
	maxOps int // maximum operations allowed
	maxLen int // maximum length allowed

	serial       uint
	idx          int                // Cursor into `info` and `pos` arrays
	scratchFlags bufferScratchFlags /* Have space-fallback, etc. */

	haveOutput bool

	planCache map[Face][]*shapePlan
}

// NewBuffer allocate a storage with default options.
// It should then be populated with `AddRunes` and shapped with `Shape`.
func NewBuffer() *Buffer {
	return &Buffer{
		ClusterLevel:  MonotoneGraphemes,
		maxOps:        maxOpsDefault,
		planCache:     map[Face][]*shapePlan{},
	}
}

// AddRune appends a character with the Unicode value of `codepoint` to `b`, and
// gives it the initial cluster value of `cluster`. Clusters can be any thing
// the client wants, they are usually used to refer to the index of the
// character in the input text stream and are output in the
// `GlyphInfo.Cluster` field.
// This also clears the posterior context (see `AddRunes`).
func (b *Buffer) AddRune(codepoint rune, cluster int) {
	b.append(codepoint, cluster)
	b.clearContext(1)
}

func (b *Buffer) append(codepoint rune, cluster int) {
	b.Info = append(b.Info, GlyphInfo{codepoint: codepoint, Cluster: cluster})
	b.Pos = append(b.Pos, GlyphPosition{})
}

// AddRunes appends characters from `text` array to `b`. `itemOffset` is the
// position of the first character from `text` that will be appended, and
// `itemLength` is the number of character to add (-1 means the end of the slice).
// When shaping part of a larger text (e.g. a run of text from a paragraph),
// instead of passing just the substring
// corresponding to the run, it is preferable to pass the whole
// paragraph and specify the run start and length as `itemOffset` and
// `itemLength`, respectively, to give HarfBuzz the full context to be able,
// for example, to do cross-run Arabic shaping or properly handle combining
// marks at start of run.
// The cluster value attributed to each rune is the index in the `text` slice.
func (b *Buffer) AddRunes(text []rune, itemOffset, itemLength int) {
	/* If buffer is empty and pre-context provided, install it.
	* This check is written this way, to make sure people can
	* provide pre-context in one add_utf() call, then provide
	* text in a follow-up call.  See:
	*
	* https://bugzilla.mozilla.org/show_bug.cgi?id=801410#c13 */
	if len(b.Info) == 0 && itemOffset > 0 {
		// add pre-context
		b.clearContext(0)
		for prev := itemOffset - 1; prev >= 0 && len(b.context[0]) < contextLength; prev-- {
			b.context[0] = append(b.context[0], text[prev])
		}
	}

	if itemLength < 0 {
		itemLength = len(text) - itemOffset
	}

	for i, u := range text[itemOffset : itemOffset+itemLength] {
		b.append(u, itemOffset+i)
	}

	// add post-context
	s := itemOffset + itemLength + contextLength
	if s > len(text) {
		s = len(text)
	}
	b.context[1] = text[itemOffset+itemLength : s]
}

// GuessSegmentProperties fills unset buffer segment properties based on buffer Unicode
// contents and can be used when no other information is available.
//
// If buffer `Props.Script` is zero, it
// will be set to the Unicode script of the first character in
// the buffer that has a script other than Common,
// Inherited, and Unknown.
//
// Next, if buffer `Props.Direction` is zero,
// it will be set to the natural horizontal direction of the
// buffer script, defaulting to `LeftToRight`.
//
// Finally, if buffer Props.Language is empty,
// it will be set to the process's default language.
func (b *Buffer) GuessSegmentProperties() {
	/* If script is not set, guess from buffer contents */
	if b.Props.Script == 0 {
		for _, info := range b.Info {
			script := language.LookupScript(info.codepoint)
			if script != language.Common && script != language.Inherited && script != language.Unknown {
				b.Props.Script = script
				break
			}
		}
	}

	/* If direction is unset, guess from script */
	if b.Props.Direction == 0 {
		b.Props.Direction = getHorizontalDirection(b.Props.Script)
		if b.Props.Direction == 0 {
			b.Props.Direction = LeftToRight
		}
	}

	/* If language is not set, use default language from locale */
	if b.Props.Language == "" {
		b.Props.Language = language.DefaultLanguage()
	}
}

// Clear resets `b` to its initial empty state (including user settings).
// This method should be used to reuse the allocated memory.
func (b *Buffer) Clear() {
	b.Flags = 0
	b.Invisible = 0
	b.NotFound = 0

	b.Props = SegmentProperties{}
	b.scratchFlags = 0

	b.haveOutput = false

	b.idx = 0
	b.Info = b.Info[:0]
	b.outInfo = b.outInfo[:0]
	b.Pos = b.Pos[:0]
	b.clearContext(0)
	b.clearContext(1)

	b.serial = 0
}

// cur returns the glyph at the cursor, optionaly shifted by `i`.
// Its simply a syntactic sugar for `&b.Info[b.idx+i] `
func (b *Buffer) cur(i int) *GlyphInfo { return &b.Info[b.idx+i] }

// cur returns the position at the cursor, optionaly shifted by `i`.
// Its simply a syntactic sugar for `&b.Pos[b.idx+i]
func (b *Buffer) curPos(i int) *GlyphPosition { return &b.Pos[b.idx+i] }

// returns the last glyph of `outInfo`
func (b Buffer) prev() *GlyphInfo {
	return &b.outInfo[len(b.outInfo)-1]
}

// func (b Buffer) has_separate_output() bool { return info != b.outInfo }

func (b *Buffer) backtrackLen() int {
	if b.haveOutput {
		return len(b.outInfo)
	}
	return b.idx
}

func (b *Buffer) lookaheadLen() int { return len(b.Info) - b.idx }

// func (b *Buffer) nextSerial() uint {
// 	out := b.serial
// 	b.serial++
// 	return out
// }

// Copies glyph at `idx` to `outInfo` before replacing its codepoint by `u`
// Advances `idx`
func (b *Buffer) replaceGlyph(u rune) {
	b.replaceGlyphs(1, []rune{u}, nil)
}

// Copies glyph at `idx` to `outInfo` before replacing its codepoint by `u`
// Advances `idx`
func (b *Buffer) replaceGlyphIndex(g GID) {
	b.outInfo = append(b.outInfo, b.Info[b.idx])
	b.outInfo[len(b.outInfo)-1].Glyph = g
	b.idx++
}

// Merges clusters in [idx:idx+numIn], then duplicate `Info[idx]` len(codepoints) times to `outInfo`.
// Advances `idx` by `numIn`. Assume that idx + numIn <= len(Info)
// Also replaces their codepoint by `codepoints` and their glyph by `glyphs` if non nil
func (b *Buffer) replaceGlyphs(numIn int, codepoints []rune, glyphs []GID) {
	b.mergeClusters(b.idx, b.idx+numIn)

	var origInfo *GlyphInfo
	if b.idx < len(b.Info) {
		origInfo = b.cur(0)
	} else {
		origInfo = b.prev()
	}
	replaceCodepoints := codepoints != nil
	replaceGlyphs := glyphs != nil
	L := len(b.outInfo)

	Lplus := max(len(codepoints), len(glyphs))
	b.outInfo = append(b.outInfo, make([]GlyphInfo, Lplus)...)
	for i := 0; i < Lplus; i++ {
		b.outInfo[L+i] = *origInfo
		if replaceCodepoints {
			b.outInfo[L+i].codepoint = codepoints[i]
		}
		if replaceGlyphs {
			b.outInfo[L+i].Glyph = glyphs[i]
		}
	}

	b.idx += numIn
}

// makes a copy of the glyph at idx to output and replace in output `codepoint`
// by `r`. Does NOT adavance `idx`
func (b *Buffer) outputRune(r rune) {
	b.replaceGlyphs(0, []rune{r}, nil)
}

// same as outputRune
func (b *Buffer) outputGlyphIndex(g GID) {
	b.replaceGlyphs(0, nil, []GID{g})
}

// Copies glyph at idx to output but doesn't advance idx
func (b *Buffer) copyGlyph() {
	b.outInfo = append(b.outInfo, b.Info[b.idx])
}

// Copies glyph at `idx` to `outInfo` and advance `idx`.
// If there's no output, just advance `idx`.
func (b *Buffer) nextGlyph() {
	if b.haveOutput {
		b.outInfo = append(b.outInfo, b.Info[b.idx])
	}

	b.idx++
}

// Copies `n` glyphs from `idx` to `outInfo` and advances `idx`.
// If there's no output, just advance idx.
func (b *Buffer) nextGlyphs(n int) {
	if b.haveOutput {
		b.outInfo = append(b.outInfo, b.Info[b.idx:b.idx+n]...)
	}
	b.idx += n
}

// skipGlyph advances idx without copying to output
func (b *Buffer) skipGlyph() { b.idx++ }

func (b *Buffer) resetMasks(mask GlyphMask) {
	for j := range b.Info {
		b.Info[j].Mask = mask
	}
}

func (b *Buffer) setMasks(value, mask GlyphMask, clusterStart, clusterEnd int) {
	notMask := ^mask
	value &= mask

	if mask == 0 {
		return
	}

	for i, info := range b.Info {
		if clusterStart <= info.Cluster && info.Cluster < clusterEnd {
			b.Info[i].Mask = (info.Mask & notMask) | value
		}
	}
}

func (b *Buffer) mergeClusters(start, end int) {
	if end-start < 2 {
		return
	}

	if b.ClusterLevel == Characters {
		b.unsafeToBreak(start, end)
		return
	}

	cluster := b.Info[start].Cluster

	for i := start + 1; i < end; i++ {
		cluster = min(cluster, b.Info[i].Cluster)
	}

	/* Extend end */
	for end < len(b.Info) && b.Info[end-1].Cluster == b.Info[end].Cluster {
		end++
	}

	/* Extend start */
	for b.idx < start && b.Info[start-1].Cluster == b.Info[start].Cluster {
		start--
	}

	/* If we hit the start of buffer, continue in out-buffer. */
	if b.idx == start {
		startC := b.Info[start].Cluster
		for i := len(b.outInfo); i != 0 && b.outInfo[i-1].Cluster == startC; i-- {
			b.outInfo[i-1].setCluster(cluster, 0)
		}
	}

	for i := start; i < end; i++ {
		b.Info[i].setCluster(cluster, 0)
	}
}

// merge clusters for deleting current glyph, and skip it.
func (b *Buffer) deleteGlyph() {
	/* The logic here is duplicated in hb_ot_hide_default_ignorables(). */

	cluster := b.Info[b.idx].Cluster
	if b.idx+1 < len(b.Info) && cluster == b.Info[b.idx+1].Cluster {
		/* Cluster survives; do nothing. */
		goto done
	}

	if len(b.outInfo) != 0 {
		/* Merge cluster backward. */
		if cluster < b.outInfo[len(b.outInfo)-1].Cluster {
			mask := b.Info[b.idx].Mask
			oldCluster := b.outInfo[len(b.outInfo)-1].Cluster
			for i := len(b.outInfo); i != 0 && b.outInfo[i-1].Cluster == oldCluster; i-- {
				b.outInfo[i-1].setCluster(cluster, mask)
			}
		}
		goto done
	}

	if b.idx+1 < len(b.Info) {
		/* Merge cluster forward. */
		b.mergeClusters(b.idx, b.idx+2)
		goto done
	}

done:
	b.skipGlyph()
}

// unsafeToBreak adds the flag `GlyphFlagUnsafeToBreak`
// when needed, between `start` and `end`.
func (b *Buffer) unsafeToBreak(start, end int) {
	if end-start < 2 {
		return
	}
	b.unsafeToBreakImpl(start, end)
}

func (b *Buffer) unsafeToBreakImpl(start, end int) {
	cluster := findMinCluster(b.Info, start, end, maxInt)
	b.unsafeToBreakSetMask(b.Info, start, end, cluster)
}

// return the smallest cluster between `cluster` and  infos[start:end]
func findMinCluster(infos []GlyphInfo, start, end, cluster int) int {
	for i := start; i < end; i++ {
		cluster = min(cluster, infos[i].Cluster)
	}
	return cluster
}

func (b *Buffer) unsafeToBreakSetMask(infos []GlyphInfo,
	start, end, cluster int,
) {
	for i := start; i < end; i++ {
		if cluster != infos[i].Cluster {
			b.scratchFlags |= bsfHasUnsafeToBreak
			infos[i].Mask |= GlyphUnsafeToBreak
		}
	}
}

func (b *Buffer) unsafeToBreakFromOutbuffer(start, end int) {
	if !b.haveOutput {
		b.unsafeToBreakImpl(start, end)
		return
	}

	//   assert (start <= out_len);
	//   assert (idx <= end);

	cluster := math.MaxInt32
	cluster = findMinCluster(b.outInfo, start, len(b.outInfo), cluster)
	cluster = findMinCluster(b.Info, b.idx, end, cluster)
	b.unsafeToBreakSetMask(b.outInfo, start, len(b.outInfo), cluster)
	b.unsafeToBreakSetMask(b.Info, b.idx, end, cluster)
}

// reset `b.outInfo`, and adjust `pos` to have
// same length as `Info` (without zeroing its values)
func (b *Buffer) clearPositions() {
	b.haveOutput = false
	// b.have_positions = true

	b.outInfo = b.outInfo[:0]

	L := len(b.Info)
	if cap(b.Pos) >= L {
		b.Pos = b.Pos[:L]
	} else {
		b.Pos = make([]GlyphPosition, L)
	}
}

// truncate `outInfo` and set `haveOutput`
func (b *Buffer) removeOutput(setOutput bool) {
	b.haveOutput = setOutput
	// b.have_positions = false

	b.outInfo = b.outInfo[:0]
}

// truncate `outInfo` and set `haveOutput` to true
func (b *Buffer) clearOutput() {
	b.removeOutput(true)
	b.idx = 0
}

func (b *Buffer) clearContext(side uint) { b.context[side] = b.context[side][:0] }

// reverses the subslice [start:end] of the buffer contents
func (b *Buffer) reverseRange(start, end int) {
	if end-start < 2 {
		return
	}
	info := b.Info[start:end]
	pos := b.Pos[start:end]
	L := len(info)
	_ = pos[L-1] // BCE
	for i := L/2 - 1; i >= 0; i-- {
		opp := L - 1 - i
		info[i], info[opp] = info[opp], info[i]
		pos[i], pos[opp] = pos[opp], pos[i] // same length
	}
}

// Reverse reverses buffer contents, that is the `Info` and `Pos` slices.
func (b *Buffer) Reverse() { b.reverseRange(0, len(b.Info)) }

func (b *Buffer) reverseClusters() {
	b.reverseGroups(func(gi1, gi2 *GlyphInfo) bool {
		return gi1.Cluster == gi2.Cluster
	}, false)
}

// mergeClusters = false
func (b *Buffer) reverseGroups(groupFunc func(*GlyphInfo, *GlyphInfo) bool, mergeClusters bool) {
	if len(b.Info) == 0 {
		return
	}

	count := len(b.Info)
	start := 0
	var i int
	for i = 1; i < count; i++ {
		if !groupFunc(&b.Info[i-1], &b.Info[i]) {
			if mergeClusters {
				b.mergeClusters(start, i)
			}
			b.reverseRange(start, i)
			start = i
		}
	}

	if mergeClusters {
		b.mergeClusters(start, i)
	}
	b.reverseRange(start, i)

	b.Reverse()
}

// swap back the temporary outInfo buffer to `Info`
// and resets the cursor `idx`.
// Assume that haveOutput is true, and toogle it.
func (b *Buffer) swapBuffers() {
	b.nextGlyphs(len(b.Info) - b.idx)
	b.haveOutput = false
	b.Info, b.outInfo = b.outInfo, b.Info
	b.idx = 0
}

// returns an unique id
func (b *Buffer) allocateLigID() uint8 {
	ligID := uint8(b.serial & 0x07)
	b.serial++
	if ligID == 0 { // in case of overflow
		ligID = b.allocateLigID()
	}
	return ligID
}

func (b *Buffer) shiftForward(count int) {
	//   assert (have_output);
	L := len(b.Info)
	b.Info = append(b.Info, make([]GlyphInfo, count)...)
	copy(b.Info[b.idx+count:], b.Info[b.idx:L])
	b.idx += count
}

func (b *Buffer) moveTo(i int) {
	if !b.haveOutput {
		// assert(i <= len)
		b.idx = i
		return
	}

	// assert(i <= out_len+(len-idx))
	outL := len(b.outInfo)
	if outL < i {
		count := i - outL
		b.outInfo = append(b.outInfo, b.Info[b.idx:count+b.idx]...)
		b.idx += count
	} else if outL > i {
		/* Tricky part: rewinding... */
		count := outL - i

		if b.idx < count {
			b.shiftForward(count - b.idx)
		}

		// assert(idx >= count)

		b.idx -= count
		copy(b.Info[b.idx:], b.outInfo[outL-count:outL])
		b.outInfo = b.outInfo[:outL-count]
	}
}

// iterator over the grapheme of a buffer
type graphemesIterator struct {
	buffer *Buffer
	start  int
}

// at the end of the buffer, start >= len(info)
func (g *graphemesIterator) next() (start, end int) {
	info := g.buffer.Info
	count := len(info)
	start = g.start
	if start >= count {
		return
	}
	for end = g.start + 1; end < count && info[end].isContinuation(); end++ {
	}
	g.start = end
	return start, end
}

func (b *Buffer) graphemesIterator() (*graphemesIterator, int) {
	return &graphemesIterator{buffer: b}, len(b.Info)
}

// iterator over clusters of a buffer with the loop
// for start, end := iter.Next(); start < count; start, end = iter.Next() {}
type clusterIterator struct {
	buffer *Buffer
	start  int
}

// returns the next cluster delimited by [start, end[
func (c *clusterIterator) next() (start, end int) {
	info := c.buffer.Info
	count := len(info)
	start = c.start
	if start >= count {
		return
	}
	cluster := info[start].Cluster
	for end = start + 1; end < count && cluster == info[end].Cluster; end++ {
	}
	c.start = end

	return start, end
}

func (b *Buffer) clusterIterator() (*clusterIterator, int) {
	return &clusterIterator{buffer: b}, len(b.Info)
}

// iterator over syllables of a buffer
type syllableIterator struct {
	buffer *Buffer
	start  int
}

func (c *syllableIterator) next() (start, end int) {
	info := c.buffer.Info
	count := len(c.buffer.Info)
	start = c.start
	if start >= count {
		return
	}
	syllable := info[start].syllable
	for end = start + 1; end < count && syllable == info[end].syllable; end++ {
	}
	c.start = end
	return start, end
}

func (b *Buffer) syllableIterator() (*syllableIterator, int) {
	return &syllableIterator{buffer: b}, len(b.Info)
}

// only modifies Info, thus assume Pos is not used yet
func (b *Buffer) sort(start, end int, compar func(a, b *GlyphInfo) int) {
	for i := start + 1; i < end; i++ {
		j := i
		for j > start && compar(&b.Info[j-1], &b.Info[i]) > 0 {
			j--
		}
		if i == j {
			continue
		}
		// move item i to occupy place for item j, shift what's in between.
		b.mergeClusters(j, i+1)

		t := b.Info[i]
		copy(b.Info[j+1:], b.Info[j:i])
		b.Info[j] = t
	}
}

func (b *Buffer) mergeOutClusters(start, end int) {
	if b.ClusterLevel == Characters {
		return
	}

	if end-start < 2 {
		return
	}

	cluster := b.outInfo[start].Cluster

	for i := start + 1; i < end; i++ {
		cluster = min(cluster, b.outInfo[i].Cluster)
	}

	/* Extend start */
	for start != 0 && b.outInfo[start-1].Cluster == b.outInfo[start].Cluster {
		start--
	}

	/* Extend end */
	for end < len(b.outInfo) && b.outInfo[end-1].Cluster == b.outInfo[end].Cluster {
		end++
	}

	/* If we hit the end of out-buffer, continue in buffer. */
	if end == len(b.outInfo) {
		endC := b.outInfo[end-1].Cluster
		for i := b.idx; i < len(b.Info) && b.Info[i].Cluster == endC; i++ {
			b.Info[i].setCluster(cluster, 0)
		}
	}

	for i := start; i < end; i++ {
		b.outInfo[i].setCluster(cluster, 0)
	}
}
