// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

// Package segmenter implements Unicode rules used
// to segment a paragraph of text according to several criteria.
// In particular, it provides a way of delimiting line break opportunities.
//
// The API of the package follows the very nice iterator pattern proposed
// in github.com/npillmayer/uax,
// but use a somewhat simpler internal implementation, inspired by Pango.
//
// The reference documentation is at https://unicode.org/reports/tr14
// and https://unicode.org/reports/tr29.
package segmenter

import (
	"unicode"

	ucd "github.com/go-text/typesetting/unicodedata"
)

// runeAttr is a flag storing the break properties between two runes of
// the input text, as computed by computeAttributes
type runeAttr uint8

const (
	aLineBreak runeAttr = 1 << iota
	aMandatoryBreak

	// this flag is on if the cursor can appear in front of a character.
	// i.e. if we are at a grapheme boundary.
	aGraphemeBoundary
)

const paragraphSeparator rune = 0x2029

// lineBreakClass stores the Line Break Property
// See https://unicode.org/reports/tr14/#Properties
type lineBreakClass = *unicode.RangeTable

// graphemeBreakClass stores the Unicode Grapheme Cluster Break Property
// See https://unicode.org/reports/tr29/#Grapheme_Cluster_Break_Property_Values
type graphemeBreakClass = *unicode.RangeTable

// cursor holds the information for the current index
// processed by `computeAttributes`, that is
// the context provided by previous and next runes in the text
type cursor struct {
	prev rune // the rune at index i-1
	r    rune // the rune at index i
	next rune // the rune at index i+1

	// is r included in `ucd.Extended_Pictographic`,
	// cached for efficiency
	isExtentedPic bool

	// the following fields persists across iterations

	prevGrapheme graphemeBreakClass // the Grapheme Break property at index i-1
	grapheme     graphemeBreakClass // the Grapheme Break property at index i

	// true if the `prev` rune was an odd Regional_Indicator, false if it was even or not an RI
	// used for rules GB12 and GB13
	// see [updateGraphemeRIOdd]
	isPrevGraphemeRIOdd bool

	prevPrevLine lineBreakClass // the Line Break Class at index i-2 (see rules LB9 and LB10 for edge cases)
	prevLine     lineBreakClass // the Line Break Class at index i-1 (see rules LB9 and LB10 for edge cases)
	line         lineBreakClass // the Line Break Class at index i
	nextLine     lineBreakClass // the Line Break Class at index i+1

	// the last rune after spaces, used in rules LB14,LB15,LB16,LB17
	// to match ... SP* ...
	beforeSpaces lineBreakClass

	// true if the `prev` rune was an odd Regional_Indicator, false if it was even or not an RI
	// used for rules LB30a
	// see [updateGrRIOdd]
	isPrevLinebreakRIOdd bool

	// are we in a numeric sequence, as defined in Example 7 of customisation for LB25
	numSequence numSequenceState

	// are we in an emoji sequence, as defined in rule GB11
	// see [updatePictoSequence]
	pictoSequence pictoSequenceState
}

// initialise the cursor properties
// some of them are set in [startIteration]
func newCursor(text []rune) *cursor {
	cr := cursor{
		prevPrevLine: ucd.BreakXX,
	}

	// `startIteration` set `breakCl` from `nextBreakCl`
	// so we need to init this field before the first iteration
	cr.nextLine = ucd.BreakXX
	if len(text) != 0 {
		cr.nextLine = ucd.LookupLineBreakClass(text[0])
	}
	return &cr
}

// computeAttributes does the heavy lifting of the segmentation,
// by computing a break attribute for each rune
//
// More precisely, `attributes` is a slice of length len(text)+1,
// which will be filled at index i by the attribute describing the
// break between rune at index i-1 and index i
//
// Unicode defines a lot of properties; for now we only handle
// graphemes and line breaking
//
// The rules are somewhat complex, but the general logic is pretty simple:
// iterate through the input slice, fetch context information
// from previous and following runes required by the rules,
// and finaly apply them.
// Some rules require variable length lookup, which we handle by keeping
// a state in a [cursor] object.
func computeAttributes(text []rune, attributes []runeAttr) {
	// initialise the cursor properties
	cr := newCursor(text)

	for i := 0; i <= len(text); i++ { // note that we accept i == len(text) to fill the last attribute
		cr.startIteration(text, i)

		var attr runeAttr

		// UAX#29 Grapheme Boundaries

		isGraphemeBoundary := cr.applyGraphemeBoundaryRules()
		if isGraphemeBoundary {
			attr |= aGraphemeBoundary
		}

		// UAX#14 Line Breaking

		bo := cr.applyLineBreakingRules()
		switch bo {
		case breakEmpty:
			// rule LB31 : default to allow line break
			attr |= aLineBreak
		case breakProhibited:
			attr &^= aLineBreak
		case breakAllowed:
			attr |= aLineBreak
		case breakMandatory:
			attr |= aLineBreak
			attr |= aMandatoryBreak
		}

		cr.endIteration(i == 0)

		attributes[i] = attr
	}

	// start and end of the paragraph are always
	// grapheme boundaries
	attributes[0] |= aGraphemeBoundary         // Rule GB1
	attributes[len(text)] |= aGraphemeBoundary // Rule GB2

	// never break before the first char,
	// but always break after the last
	attributes[0] &^= aLineBreak             // Rule LB2
	attributes[len(text)] |= aLineBreak      // Rule LB3
	attributes[len(text)] |= aMandatoryBreak // Rule LB3
}

// Segmenter is the entry point of the package.
//
// Usage :
//
//	var seg Segmenter
//	seg.Init(...)
//	iter := seg.LineIterator()
//	for iter.Next() {
//	  ... // do something with iter.Line()
//	}
type Segmenter struct {
	text []rune
	// with length len(text) + 1 :
	// the attribute at indice i is about the
	// rune at i-1 and i
	// See also `computeAttributes`
	// Example :
	// 	text : 			[b, 		u, 	l, 	l]
	// 	attributes :	[<start> b, b u, u l, l l, l <end>]
	attributes []runeAttr
}

// Init resets the segmenter storage with the given input,
// and computes the attributes required to segment the text.
func (seg *Segmenter) Init(paragraph []rune) {
	seg.text = append(seg.text[:0], paragraph...)
	seg.attributes = append(seg.attributes[:0], make([]runeAttr, len(paragraph)+1)...)
	computeAttributes(seg.text, seg.attributes)
}

// attributeIterator is an helper type used to
// handle iterating over a slice of runeAttr
type attributeIterator struct {
	src       *Segmenter
	pos       int      // the current position in the input slice
	lastBreak int      // the start of the current segment
	flag      runeAttr // break where this flag is on
}

// next returns true if there is still a segment to process,
// and advances the iterator; or return false.
// if returning true, the segment it at li.lastBreak:li.pos
func (iter *attributeIterator) next() bool {
	iter.lastBreak = iter.pos // remember the start of the next line
	iter.pos++
	for iter.pos <= len(iter.src.text) {
		// can we break before i ?
		if iter.src.attributes[iter.pos]&iter.flag != 0 {
			return true
		}
		iter.pos++
	}
	return false
}

// LineIterator provides a convenient way of
// iterating over the lines delimited by a `Segmenter`.
type LineIterator struct {
	attributeIterator
}

// Next returns true if there is still a line to process,
// and advances the iterator; or return false.
func (li *LineIterator) Next() bool { return li.next() }

// Line returns the current `Line`
func (li *LineIterator) Line() Line {
	return Line{
		Offset:           li.lastBreak,
		Text:             li.src.text[li.lastBreak:li.pos], // pos is not included since we break right before
		IsMandatoryBreak: li.src.attributes[li.pos]&aMandatoryBreak != 0,
	}
}

// Line is the content of a line delimited by the segmenter.
type Line struct {
	// Text is a subslice of the original input slice, containing the delimited line
	Text []rune
	// Offset is the start of the line in the input rune slice
	Offset int
	// IsMandatoryBreak is true if breaking (at the end of the line)
	// is mandatory
	IsMandatoryBreak bool
}

// LineIterator returns an iterator on the lines
// delimited in [Init].
func (sg *Segmenter) LineIterator() *LineIterator {
	return &LineIterator{attributeIterator: attributeIterator{src: sg, flag: aLineBreak}}
}

// GraphemeIterator provides a convenient way of
// iterating over the graphemes delimited by a `Segmenter`.
type GraphemeIterator struct {
	attributeIterator
}

// Next returns true if there is still a grapheme to process,
// and advances the iterator; or return false.
func (gr *GraphemeIterator) Next() bool { return gr.next() }

// Grapheme returns the current `Grapheme`
func (gr *GraphemeIterator) Grapheme() Grapheme {
	return Grapheme{
		Offset: gr.lastBreak,
		Text:   gr.src.text[gr.lastBreak:gr.pos],
	}
}

// Line is the content of a grapheme delimited by the segmenter.
type Grapheme struct {
	// Text is a subslice of the original input slice, containing the delimited grapheme
	Text []rune
	// Offset is the start of the grapheme in the input rune slice
	Offset int
}

// GraphemeIterator returns an iterator over the graphemes
// delimited in [Init].
func (sg *Segmenter) GraphemeIterator() *GraphemeIterator {
	return &GraphemeIterator{attributeIterator: attributeIterator{src: sg, flag: aGraphemeBoundary}}
}
