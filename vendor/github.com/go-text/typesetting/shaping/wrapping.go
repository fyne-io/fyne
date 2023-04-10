package shaping

import (
	"sort"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/segmenter"
	"golang.org/x/image/math/fixed"
)

// glyphIndex is the index in a Glyph slice
type glyphIndex = int

// mapRunesToClusterIndices
// returns a slice that maps rune indicies in the text to the index of the
// first glyph in the glyph cluster containing that rune in the shaped text.
// The indicies are relative to the region of runes covered by the input run.
// To translate an absolute rune index in text into a rune index into the returned
// mapping, subtract run.Runes.Offset first. If the provided buf is large enough to
// hold the return value, it will be used instead of allocating a new slice.
func mapRunesToClusterIndices(dir di.Direction, runes Range, glyphs []Glyph, buf []glyphIndex) []glyphIndex {
	if runes.Count <= 0 {
		return nil
	}
	var mapping []glyphIndex
	if cap(buf) >= runes.Count {
		mapping = buf[:runes.Count]
	} else {
		mapping = make([]glyphIndex, runes.Count)
	}
	glyphCursor := 0
	rtl := dir.Progression() == di.TowardTopLeft
	if rtl {
		glyphCursor = len(glyphs) - 1
	}
	// off tracks the offset position of the glyphs from the first rune of the
	// shaped text. This must be subtracted from all cluster indicies in order to
	// normalize them into the range [0,runes.Count).
	off := runes.Offset
	for i := 0; i < runes.Count; i++ {
		for glyphCursor >= 0 && glyphCursor < len(glyphs) &&
			((rtl && glyphs[glyphCursor].ClusterIndex-off <= i) ||
				(!rtl && glyphs[glyphCursor].ClusterIndex-off < i)) {
			if rtl {
				glyphCursor--
			} else {
				glyphCursor++
			}
		}
		if rtl {
			glyphCursor++
		} else if (glyphCursor >= 0 && glyphCursor < len(glyphs) &&
			glyphs[glyphCursor].ClusterIndex-off > i) ||
			(glyphCursor == len(glyphs) && len(glyphs) > 1) {
			glyphCursor--
			targetClusterIndex := glyphs[glyphCursor].ClusterIndex - off
			for glyphCursor-1 >= 0 && glyphs[glyphCursor-1].ClusterIndex-off == targetClusterIndex {
				glyphCursor--
			}
		}
		if glyphCursor < 0 {
			glyphCursor = 0
		} else if glyphCursor >= len(glyphs) {
			glyphCursor = len(glyphs) - 1
		}
		mapping[i] = glyphCursor
	}
	return mapping
}

// mapRuneToClusterIndex finds the lowest-index glyph for the glyph cluster contiaining the rune
// at runeIdx in the source text. It uses a binary search of the glyphs in order to achieve this.
// It is equivalent to using mapRunesToClusterIndices on only a single rune index, and is thus
// more efficient for single lookups while being less efficient for runs which require many
// lookups anyway.
func mapRuneToClusterIndex(dir di.Direction, runes Range, glyphs []Glyph, runeIdx int) glyphIndex {
	var index int
	rtl := dir.Progression() == di.TowardTopLeft
	if !rtl {
		index = sort.Search(len(glyphs), func(index int) bool {
			return glyphs[index].ClusterIndex-runes.Offset > runeIdx
		})
	} else {
		index = sort.Search(len(glyphs), func(index int) bool {
			return glyphs[index].ClusterIndex-runes.Offset < runeIdx
		})
	}
	if index < 1 {
		return 0
	}
	cluster := glyphs[index-1].ClusterIndex
	if rtl && cluster-runes.Offset > runeIdx {
		return index
	}
	for index-1 >= 0 && glyphs[index-1].ClusterIndex == cluster {
		index--
	}
	return index
}

func mapRunesToClusterIndices2(dir di.Direction, runes Range, glyphs []Glyph, buf []glyphIndex) []glyphIndex {
	if runes.Count <= 0 {
		return nil
	}
	var mapping []glyphIndex
	if cap(buf) >= runes.Count {
		mapping = buf[:runes.Count]
	} else {
		mapping = make([]glyphIndex, runes.Count)
	}

	rtl := dir.Progression() == di.TowardTopLeft
	if rtl {
		for gIdx := len(glyphs) - 1; gIdx >= 0; gIdx-- {
			cluster := glyphs[gIdx].ClusterIndex
			clusterEnd := gIdx
			for gIdx-1 >= 0 && glyphs[gIdx-1].ClusterIndex == cluster {
				gIdx--
				clusterEnd = gIdx
			}
			var nextCluster int
			if gIdx-1 >= 0 {
				nextCluster = glyphs[gIdx-1].ClusterIndex
			} else {
				nextCluster = runes.Count + runes.Offset
			}
			runesInCluster := nextCluster - cluster
			clusterOffset := cluster - runes.Offset
			for i := clusterOffset; i <= runesInCluster+clusterOffset && i < len(mapping); i++ {
				mapping[i] = clusterEnd
			}
		}
	} else {
		for gIdx := 0; gIdx < len(glyphs); gIdx++ {
			cluster := glyphs[gIdx].ClusterIndex
			clusterStart := gIdx
			for gIdx+1 < len(glyphs) && glyphs[gIdx+1].ClusterIndex == cluster {
				gIdx++
			}
			var nextCluster int
			if gIdx+1 < len(glyphs) {
				nextCluster = glyphs[gIdx+1].ClusterIndex
			} else {
				nextCluster = runes.Count + runes.Offset
			}
			runesInCluster := nextCluster - cluster
			clusterOffset := cluster - runes.Offset
			for i := clusterOffset; i <= runesInCluster+clusterOffset && i < len(mapping); i++ {
				mapping[i] = clusterStart
			}
		}
	}
	return mapping
}

func mapRunesToClusterIndices3(dir di.Direction, runes Range, glyphs []Glyph, buf []glyphIndex) []glyphIndex {
	if runes.Count <= 0 {
		return nil
	}
	var mapping []glyphIndex
	if cap(buf) >= runes.Count {
		mapping = buf[:runes.Count]
	} else {
		mapping = make([]glyphIndex, runes.Count)
	}

	rtl := dir.Progression() == di.TowardTopLeft
	if rtl {
		for gIdx := len(glyphs) - 1; gIdx >= 0; {
			glyph := &glyphs[gIdx]
			// go to the start of the cluster
			gIdx -= (glyph.GlyphCount - 1)
			clusterStart := glyph.ClusterIndex - runes.Offset // map back to [0;runes.Count[
			clusterEnd := glyph.RuneCount + clusterStart
			for i := clusterStart; i <= clusterEnd && i < len(mapping); i++ {
				mapping[i] = gIdx
			}
			// go to the next cluster
			gIdx--
		}
	} else {
		for gIdx := 0; gIdx < len(glyphs); {
			glyph := &glyphs[gIdx]
			clusterStart := glyph.ClusterIndex - runes.Offset // map back to [0;runes.Count[
			clusterEnd := glyph.RuneCount + clusterStart
			for i := clusterStart; i <= clusterEnd && i < len(mapping); i++ {
				mapping[i] = gIdx
			}
			// go to the next cluster
			gIdx += glyph.GlyphCount
		}
	}
	return mapping
}

// inclusiveGlyphRange returns the inclusive range of runes and glyphs matching
// the provided start and breakAfter rune positions.
// runeToGlyph must be a valid mapping from the rune representation to the
// glyph reprsentation produced by mapRunesToClusterIndices.
// numGlyphs is the number of glyphs in the output representing the runes
// under consideration.
func inclusiveGlyphRange(dir di.Direction, start, breakAfter int, runeToGlyph []int, numGlyphs int) (glyphStart, glyphEnd glyphIndex) {
	rtl := dir.Progression() == di.TowardTopLeft
	if rtl {
		glyphStart = runeToGlyph[breakAfter]
		if start-1 >= 0 {
			glyphEnd = runeToGlyph[start-1] - 1
		} else {
			glyphEnd = numGlyphs - 1
		}
	} else {
		glyphStart = runeToGlyph[start]
		if breakAfter+1 < len(runeToGlyph) {
			glyphEnd = runeToGlyph[breakAfter+1] - 1
		} else {
			glyphEnd = numGlyphs - 1
		}
	}
	return
}

// cutRun returns the sub-run of run containing glyphs corresponding to the provided
// _inclusive_ rune range.
func cutRun(run Output, mapping []glyphIndex, startRune, endRune int) Output {
	// Convert the rune range of interest into an inclusive range within the
	// current run's runes.
	runeStart := startRune - run.Runes.Offset
	runeEnd := endRune - run.Runes.Offset
	if runeStart < 0 {
		// If the start location is prior to the run of shaped text under consideration,
		// just work from the beginning of this run.
		runeStart = 0
	}
	if runeEnd >= len(mapping) {
		// If the break location is after the entire run of shaped text,
		// keep through the end of the run.
		runeEnd = len(mapping) - 1
	}
	glyphStart, glyphEnd := inclusiveGlyphRange(run.Direction, runeStart, runeEnd, mapping, len(run.Glyphs))

	// Construct a run out of the inclusive glyph range.
	run.Glyphs = run.Glyphs[glyphStart : glyphEnd+1]
	run.RecomputeAdvance()
	run.Runes.Offset = run.Runes.Offset + runeStart
	run.Runes.Count = runeEnd - runeStart + 1
	return run
}

// breakOption represets a location within the rune slice at which
// it may be safe to break a line of text.
type breakOption struct {
	// breakAtRune is the index at which it is safe to break.
	breakAtRune int
}

// isValid returns whether a given option violates shaping rules (like breaking
// a shaped text cluster).
func (option breakOption) isValid(runeToGlyph []int, out Output) bool {
	breakAfter := option.breakAtRune - out.Runes.Offset
	nextRune := breakAfter + 1
	if nextRune < len(runeToGlyph) && breakAfter >= 0 {
		// Check if this break is valid.
		gIdx := runeToGlyph[breakAfter]
		g2Idx := runeToGlyph[nextRune]
		cIdx := out.Glyphs[gIdx].ClusterIndex
		c2Idx := out.Glyphs[g2Idx].ClusterIndex
		if cIdx == c2Idx {
			// This break is within a harfbuzz cluster, and is
			// therefore invalid.
			return false
		}
	}
	return true
}

// breaker generates line breaking candidates for a text.
type breaker struct {
	segmenter  *segmenter.LineIterator
	totalRunes int
}

// newBreaker returns a breaker initialized to break the provided text.
func newBreaker(seg *segmenter.Segmenter, text []rune) *breaker {
	seg.Init(text)
	br := &breaker{
		segmenter:  seg.LineIterator(),
		totalRunes: len(text),
	}
	return br
}

// next returns a naive break candidate which may be invalid.
func (b *breaker) next() (option breakOption, ok bool) {
	if b.segmenter.Next() {
		currentSegment := b.segmenter.Line()
		// Note : we dont use penalties for Mandatory Breaks so far,
		// we could add it with currentSegment.IsMandatoryBreak
		option := breakOption{
			breakAtRune: currentSegment.Offset + len(currentSegment.Text) - 1,
		}
		return option, true
	}
	// Unicode rules impose to always break at the end
	return breakOption{}, false
}

// Range indicates the location of a sequence of elements within a longer slice.
type Range struct {
	Offset int
	Count  int
}

// Line holds runs of shaped text wrapped onto a single line. All the contained
// Output should be displayed sequentially on one line.
type Line []Output

// WrapConfig provides line-wrapper settings.
type WrapConfig struct {
	// TruncateAfterLines is the number of lines of text to allow before truncating
	// the text. A value of zero means no limit.
	TruncateAfterLines int
	// Truncator, if provided, will be inserted at the end of a truncated line. This
	// feature is only active if TruncateAfterLines is nonzero.
	Truncator Output
	// TextContinues indicates that the paragraph wrapped by this config is not the
	// final paragraph in the text. This alters text truncation when filling the
	// final line permitted by TruncateAfterLines. If the text of this paragraph
	// does fit entirely on TruncateAfterLines, normally the truncator symbol would
	// not be inserted. However, if the overall body of text continues beyond this
	// paragraph (indicated by TextContinues), the truncator should still be inserted
	// to indicate that further paragraphs of text were truncated. This field has
	// no effect if TruncateAfterLines is zero.
	TextContinues bool
}

// WithTruncator returns a copy of WrapConfig with the Truncator field set to the
// result of shaping input with shaper.
func (w WrapConfig) WithTruncator(shaper Shaper, input Input) WrapConfig {
	w.Truncator = shaper.Shape(input)
	return w
}

// runMapper efficiently maps a run to glyph clusters.
type runMapper struct {
	// valid indicates that the mapping field is populated.
	valid bool
	// runIdx is the index of the mapped run within glyphRuns.
	runIdx int
	// mapping holds the rune->glyph mapping for the run at index mappedRun within
	// glyphRuns.
	mapping []glyphIndex
}

// mapRun updates the mapping field to be valid for the given run. It will skip the mapping
// operation if the provided runIdx is equal to the runIdx of the previous call, as the
// current mapping value is already correct.
func (r *runMapper) mapRun(runIdx int, run Output) {
	if r.runIdx != runIdx || !r.valid {
		r.mapping = mapRunesToClusterIndices3(run.Direction, run.Runes, run.Glyphs, r.mapping)
		r.runIdx = runIdx
		r.valid = true
	}
}

// LineWrapper holds reusable state for a line wrapping operation. Reusing
// LineWrappers for multiple paragraphs should improve performance.
type LineWrapper struct {
	// config holds the current line wrapping settings.
	config WrapConfig
	// truncating tracks whether the wrapper should be performing truncation.
	truncating bool
	// seg is an internal storage used to initiate the breaker iterator.
	seg segmenter.Segmenter

	// breaker provides line-breaking candidates.
	breaker *breaker

	// mapper tracks rune->glyphCluster mappings.
	mapper runMapper
	// unusedBreak is a break requested from the breaker in a previous iteration
	// but which was not chosen as the line ending. Subsequent invocations of
	// WrapLine should start with this break.
	unusedBreak breakOption
	// isUnused indicates that the unusedBreak field is valid.
	isUnused bool
	// glyphRuns holds the runs of shaped text being wrapped.
	glyphRuns []Output
	// currentRun holds the index in use within glyphRuns.
	currentRun int
	// lineStartRune is the rune index of the first rune on the next line to
	// be shaped.
	lineStartRune int
	// more indicates that the iteration API has more data to return.
	more bool
}

// Prepare initializes the LineWrapper for the given paragraph and shaped text.
// It must be called prior to invoking WrapNextLine.
func (l *LineWrapper) Prepare(config WrapConfig, paragraph []rune, shapedRuns ...Output) {
	l.config = config
	l.truncating = l.config.TruncateAfterLines > 0
	l.breaker = newBreaker(&l.seg, paragraph)
	l.glyphRuns = shapedRuns
	l.isUnused = false
	l.currentRun = 0
	l.lineStartRune = 0
	l.more = true
	l.mapper.valid = false
}

// WrapParagraph wraps the paragraph's shaped glyphs to a constant maxWidth.
// It is equivalent to iteratively invoking WrapLine with a constant maxWidth.
// If the config has a non-zero TruncateAfterLines, WrapParagraph will return at most
// that many lines. The truncated return value is the count of runes truncated from
// the end of the text.
func (l *LineWrapper) WrapParagraph(config WrapConfig, maxWidth int, paragraph []rune, shapedRuns ...Output) (_ []Line, truncated int) {
	if len(shapedRuns) == 1 && shapedRuns[0].Advance.Ceil() < maxWidth && !(config.TextContinues && config.TruncateAfterLines == 1) {
		return []Line{shapedRuns}, 0
	}
	l.Prepare(config, paragraph, shapedRuns...)
	var lines []Line
	var done bool
	for !done {
		var line Line
		line, truncated, done = l.WrapNextLine(maxWidth)
		lines = append(lines, line)
	}
	return lines, truncated
}

// nextBreakOption returns the next rune offset at which the line can be broken,
// if any. If it returns false, there are no more candidates.
func (l *LineWrapper) nextBreakOption() (breakOption, bool) {
	var option breakOption
	if l.isUnused {
		option = l.unusedBreak
		l.isUnused = false
	} else {
		var breakOk bool
		option, breakOk = l.breaker.next()
		if !breakOk {
			return option, false
		}
		l.unusedBreak = option
	}
	return option, true
}

type fillResult uint8

const (
	// noCandidate indicates that it is not possible to compose a new line candidate using the provided
	// breakOption, so the best known line should be used instead.
	noCandidate fillResult = iota
	// noRunWithBreak indicates that none of the runs available to the line wrapper contain the break
	// option, so the returned candidate is the best option.
	noRunWithBreak
	// newCandidate indicates that the returned line candidate is valid.
	newCandidate
)

// fillUntil tries to fill the provided line candidate slice with runs until it reaches a run containing the
// provided break option. It returns the index of the run containing the option, the new width of the candidate
// line, the contents of the new candidate line, and a result indicating how to proceed.
func (l *LineWrapper) fillUntil(option breakOption, startRunIdx int, startWidth fixed.Int26_6, lineCandidate []Output) (newRunIdx int, newWidth fixed.Int26_6, newLineCandidate []Output, status fillResult) {
	run := l.glyphRuns[startRunIdx]
	for option.breakAtRune >= run.Runes.Count+run.Runes.Offset {
		if l.lineStartRune >= run.Runes.Offset+run.Runes.Count {
			startRunIdx++
			if startRunIdx >= len(l.glyphRuns) {
				return startRunIdx, startWidth, lineCandidate, noCandidate
			}
			run = l.glyphRuns[startRunIdx]
			continue
		} else if l.lineStartRune > run.Runes.Offset {
			// If part of this run has already been used on a previous line, trim
			// the runes corresponding to those glyphs off.
			l.mapper.mapRun(startRunIdx, run)
			run = cutRun(run, l.mapper.mapping, l.lineStartRune, run.Runes.Count+run.Runes.Offset)
		}
		// While the run being processed doesn't contain the current line breaking
		// candidate, just append it to the candidate line.
		lineCandidate = append(lineCandidate, run)
		startWidth += run.Advance
		startRunIdx++
		if startRunIdx >= len(l.glyphRuns) {
			return startRunIdx, startWidth, lineCandidate, noRunWithBreak
		}
		run = l.glyphRuns[startRunIdx]
	}
	return startRunIdx, startWidth, lineCandidate, newCandidate
}

// WrapNextLine wraps the shaped glyphs of a paragraph to a particular max width.
// It is meant to be called iteratively to wrap each line, allowing lines to
// be wrapped to different widths within the same paragraph. When done is true,
// subsequent calls to WrapNextLine (without calling Prepare) will return a nil line.
// The truncated return value is the count of runes truncated from the end of the line,
// if this line was truncated.
func (l *LineWrapper) WrapNextLine(maxWidth int) (finalLine Line, truncated int, done bool) {
	defer func() {
		if len(finalLine) > 0 {
			finalRun := finalLine[len(finalLine)-1]
			l.lineStartRune = finalRun.Runes.Count + finalRun.Runes.Offset
		}
		done = done || l.lineStartRune >= l.breaker.totalRunes
		if l.truncating {
			l.config.TruncateAfterLines--
			insertTruncator := false
			if l.config.TruncateAfterLines == 0 {
				done = true
				truncated = l.breaker.totalRunes - l.lineStartRune
				insertTruncator = truncated > 0 || l.config.TextContinues
			}
			if insertTruncator {
				finalLine = append(finalLine, l.config.Truncator)
			}
		}
		if done {
			l.more = false
		}
	}()
	if !l.more {
		return nil, truncated, true
	} else if len(l.glyphRuns) == 0 {
		return nil, truncated, true
	} else if len(l.glyphRuns[0].Glyphs) == 0 {
		// Pass empty lines through as empty.
		l.glyphRuns[0].Runes = Range{Count: l.breaker.totalRunes}
		return Line([]Output{l.glyphRuns[0]}), truncated, true
	} else if len(l.glyphRuns) == 1 && l.glyphRuns[0].Advance.Ceil() < maxWidth && !(l.config.TextContinues && l.config.TruncateAfterLines == 1) {
		return Line(l.glyphRuns), truncated, true
	}

	// lineCandidate is filled with runs as we search for valid line breaks. When we find a valid
	// option, we commit it into bestCandidate and keep looking.
	var lineCandidate, bestCandidate []Output
	// lineWidth tracks the width of the lineCandidate.
	lineWidth := fixed.I(0)
	var result fillResult

	// lineRun tracks the glyph run in use by the lineCandidate. It is
	// incremented separately so that the candidate search can run ahead of the
	// l.currentRun.
	lineRun := l.currentRun

	// truncating tracks whether this line should consider truncation options.
	truncating := l.config.TruncateAfterLines == 1
	// truncatedMaxWidth holds the maximum width of the line available for text if the truncator
	// is occupying part of the line.
	truncatedMaxWidth := maxWidth - l.config.Truncator.Advance.Ceil()

	for {
		option, ok := l.nextBreakOption()
		if !ok {
			return bestCandidate, truncated, true
		}
		lineRun, lineWidth, lineCandidate, result = l.fillUntil(
			option,
			lineRun,
			lineWidth,
			lineCandidate,
		)
		if result == noCandidate {
			return bestCandidate, truncated, true
		} else if result == noRunWithBreak {
			return lineCandidate, truncated, true
		}
		run := l.glyphRuns[lineRun]
		l.mapper.mapRun(lineRun, run)
		if !option.isValid(l.mapper.mapping, run) {
			// Reject invalid line break candidate and acquire a new one.
			continue
		}
		candidateRun := cutRun(run, l.mapper.mapping, l.lineStartRune, option.breakAtRune)
		candidateLineWidth := (candidateRun.Advance + lineWidth).Ceil()
		if candidateLineWidth > maxWidth {
			// The run doesn't fit on the line.
			if len(bestCandidate) < 1 {
				if truncating {
					return bestCandidate, truncated, true
				}
				// There is no existing candidate that fits, and we have just hit the
				// first line breaking candidate. Commit this break position as the
				// best available, even though it doesn't fit.
				lineCandidate = append(lineCandidate, candidateRun)
				l.currentRun = lineRun
				return lineCandidate, truncated, false
			} else {
				// The line is a valid, shorter wrapping. Return it and mark that
				// we should reuse the current line break candidate on the next
				// line.
				l.isUnused = true
				return bestCandidate, truncated, false
			}
		} else if truncating && candidateLineWidth > truncatedMaxWidth {
			// The run would not fit if truncated.
			finalRunRune := candidateRun.Runes.Count + candidateRun.Runes.Offset
			if finalRunRune == l.breaker.totalRunes && !l.config.TextContinues {
				// The run contains the entire end of the text, so no truncation is
				// necessary.
				bestCandidate = commitCandidate(bestCandidate, lineCandidate, candidateRun)
				l.currentRun = lineRun
				return bestCandidate, truncated, true
			}
			// We must truncate the line in order to show it.
			return bestCandidate, truncated, true
		} else {
			// The run does fit on the line. Commit this line as the best known
			// line, but keep lineCandidate unmodified so that later break
			// options can be attempted to see if a more optimal solution is
			// available.
			bestCandidate = commitCandidate(bestCandidate, lineCandidate, candidateRun)
			l.currentRun = lineRun
		}
	}
}

// commitCandidate efficiently updates destination to contain append(source, newRuns...),
// returning the resulting slice. This operation only makes sense when destination
// is not known to contain the elements of source already.
func commitCandidate(destination, source []Output, newRuns ...Output) []Output {
	destination = resize(destination, len(source), len(source)+1)
	destination = destination[:copy(destination, source)]
	destination = append(destination, newRuns...)
	return destination
}

// resize returns input resized to have the provided length and at least the provided
// capacity. It may copy the data if the provided capacity is greater than the capacity
// of in. If the provided length is greater than the provided capacity, the capacity will
// be used as the length.
func resize(input []Output, length, capacity int) []Output {
	if length > capacity {
		length = capacity
	}
	out := input
	if cap(input) < capacity {
		out = make([]Output, capacity)
		copy(out, input)
	}
	if len(out) != length {
		out = out[:length]
	}
	return out
}
