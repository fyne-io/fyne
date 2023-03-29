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

// LineWrapper holds reusable state for a line wrapping operation. Reusing
// LineWrappers for multiple paragraphs should improve performance.
type LineWrapper struct {
	// seg is an internal storage used to initiate the breaker iterator
	seg segmenter.Segmenter

	// breaker provides line-breaking candidates.
	breaker *breaker

	// mappingValid indicates that the mapping field is populated.
	mappingValid bool
	// mappedRun is the index of the mapped run within glyphRuns.
	mappedRun int
	// mapping holds the rune->glyph mapping for the run at index mappedRun within
	// glyphRuns.
	mapping []glyphIndex
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
func (l *LineWrapper) Prepare(paragraph []rune, shapedRuns ...Output) {
	l.breaker = newBreaker(&l.seg, paragraph)
	l.glyphRuns = shapedRuns
	l.isUnused = false
	l.currentRun = 0
	l.lineStartRune = 0
	l.more = true
	l.mappingValid = false
}

// WrapParagraph wraps the paragraph's shaped glyphs to a constant maxWidth.
// It is equivalent to iteratively invoking WrapLine with a constant maxWidth.
func (l *LineWrapper) WrapParagraph(maxWidth int, paragraph []rune, shapedRuns ...Output) []Line {
	if len(shapedRuns) == 1 && shapedRuns[0].Advance.Ceil() < maxWidth {
		return []Line{shapedRuns}
	}
	l.Prepare(paragraph, shapedRuns...)
	var lines []Line
	var done bool
	for !done {
		var line Line
		line, done = l.WrapNextLine(maxWidth)
		lines = append(lines, line)
	}
	return lines
}

// WrapNextLine wraps the shaped glyphs of a paragraph to a particular max width.
// It is meant to be called iteratively to wrap each line, allowing lines to
// be wrapped to different widths within the same paragraph. When done is true,
// subsequent calls to WrapNextLine (without calling Prepare) will return a nil line.
func (l *LineWrapper) WrapNextLine(maxWidth int) (_ Line, done bool) {
	defer func() {
		if done {
			l.more = false
		}
	}()
	if !l.more {
		return nil, true
	} else if len(l.glyphRuns) == 0 {
		return nil, true
	} else if len(l.glyphRuns[0].Glyphs) == 0 {
		// Pass empty lines through as empty.
		l.glyphRuns[0].Runes = Range{Count: l.breaker.totalRunes}
		return Line([]Output{l.glyphRuns[0]}), true
	} else if len(l.glyphRuns) == 1 && l.glyphRuns[0].Advance.Ceil() < maxWidth {
		return Line(l.glyphRuns), true
	}

	lineCandidate, bestCandidate := []Output{}, []Output{}
	candidateWidth := fixed.I(0)

	// mapRun performs a rune->glyph mapping for the given run, using the provided
	// run index to skip the work if that run was already mapped.
	mapRun := func(runIdx int, run Output) {
		if l.mappedRun != runIdx || !l.mappingValid {
			l.mapping = mapRunesToClusterIndices3(run.Direction, run.Runes, run.Glyphs, l.mapping)
			l.mappedRun = runIdx
			l.mappingValid = true
		}
	}

	// candidateCurrentRun tracks the glyph run in use by the lineCandidate. It is
	// incremented separately so that the candidate search can run ahead of the
	// l.currentRun.
	candidateCurrentRun := l.currentRun
	incRun := func() bool {
		candidateCurrentRun++
		return candidateCurrentRun >= len(l.glyphRuns)
	}

	for {
		run := l.glyphRuns[candidateCurrentRun]
		var option breakOption
		if l.isUnused {
			option = l.unusedBreak
			l.isUnused = false
		} else {
			var breakOk bool
			option, breakOk = l.breaker.next()
			if !breakOk {
				return bestCandidate, true
			}
			l.unusedBreak = option
		}
		for option.breakAtRune >= run.Runes.Count+run.Runes.Offset {
			if l.lineStartRune >= run.Runes.Offset+run.Runes.Count {
				if incRun() {
					return bestCandidate, true
				}
				run = l.glyphRuns[candidateCurrentRun]
				continue
			} else if l.lineStartRune > run.Runes.Offset {
				// If part of this run has already been used on a previous line, trim
				// the runes corresponding to those glyphs off.
				mapRun(candidateCurrentRun, run)
				run = cutRun(run, l.mapping, l.lineStartRune, run.Runes.Count+run.Runes.Offset)
			}
			// While the run being processed doesn't contain the current line breaking
			// candidate, just append it to the candidate line.
			lineCandidate = append(lineCandidate, run)
			candidateWidth += run.Advance
			if incRun() {
				return lineCandidate, true
			}
			run = l.glyphRuns[candidateCurrentRun]
		}
		mapRun(candidateCurrentRun, run)
		if !option.isValid(l.mapping, run) {
			// Reject invalid line break candidate and acquire a new one.
			continue
		}
		candidateRun := cutRun(run, l.mapping, l.lineStartRune, option.breakAtRune)
		if (candidateRun.Advance + candidateWidth).Ceil() > maxWidth {
			// The run doesn't fit on the line.
			if len(bestCandidate) < 1 {
				// There is no existing candidate that fits, and we have just hit the
				// first line breaking canddiate. Commit this break position as the
				// best available, even though it doesn't fit.
				lineCandidate = append(lineCandidate, candidateRun)
				l.lineStartRune = candidateRun.Runes.Offset + candidateRun.Runes.Count
				l.currentRun = candidateCurrentRun
				return lineCandidate, l.lineStartRune >= l.breaker.totalRunes
			} else {
				// The line is a valid, shorter wrapping. Return it and mark that
				// we should reuse the current line break candidate on the next
				// line.
				l.isUnused = true
				finalRunRunes := bestCandidate[len(bestCandidate)-1].Runes
				l.lineStartRune = finalRunRunes.Count + finalRunRunes.Offset
				return bestCandidate, false
			}
		} else {
			// The run does fit on the line. Commit this line as the best known
			// line, but keep lineCandidate unmodified so that later break
			// options can be attempted to see if a more optimal solution is
			// available.
			if target := len(lineCandidate) + 1; cap(bestCandidate) < target {
				bestCandidate = make([]Output, target-1, target)
			} else if len(bestCandidate) < target {
				bestCandidate = bestCandidate[:target-1]
			}
			bestCandidate = bestCandidate[:copy(bestCandidate, lineCandidate)]
			bestCandidate = append(bestCandidate, candidateRun)
			l.currentRun = candidateCurrentRun
		}
	}
}
