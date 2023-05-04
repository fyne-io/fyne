// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package shaping

import (
	"unicode"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/harfbuzz"
	"github.com/go-text/typesetting/language"
	"golang.org/x/image/math/fixed"
)

type Input struct {
	// Text is the body of text being shaped. Only the range Text[RunStart:RunEnd] is considered
	// for shaping, with the rest provided as context for the shaper. This helps with, for example,
	// cross-run Arabic shaping or handling combining marks at the start of a run.
	Text []rune
	// RunStart and RunEnd indicate the subslice of Text being shaped.
	RunStart, RunEnd int
	// Direction is the directionality of the text.
	Direction di.Direction
	// Face is the font face to render the text in.
	Face font.Face

	// Size is the requested size of the font.
	// More generally, it is a scale factor applied to the resulting metrics.
	// For instance, given a device resolution (in dpi) and a point size (like 14), the `Size` to
	// get result in pixels is given by : pointSize * dpi / 72
	Size fixed.Int26_6

	// Script is an identifier for the writing system used in the text.
	Script language.Script

	// Language is an identifier for the language of the text.
	Language language.Language
}

// Fontmap provides a general mechanism to select
// a face to use when shaping text.
type Fontmap interface {
	// ResolveFace is called by `SplitByFace` for each input rune potentially
	// triggering a face change.
	// It must always return a valid (non nil) font.Face value.
	ResolveFace(r rune) font.Face
}

var _ Fontmap = fixedFontmap(nil)

type fixedFontmap []font.Face

// ResolveFace panics if the slice is empty
func (ff fixedFontmap) ResolveFace(r rune) font.Face {
	for _, f := range ff {
		if _, has := f.NominalGlyph(r); has {
			return f
		}
	}
	return ff[0]
}

// SplitByFontGlyphs split the runes from 'input' to several items, sharing the same
// characteristics as 'input', expected for the `Face` which is set to
// the first font among 'availableFonts' providing support for all the runes
// in the item.
// Runes supported by no fonts are mapped to the first element of 'availableFonts', which
// must not be empty.
// The 'Face' field of 'input' is ignored: only 'availableFaces' are consulted.
// Rune coverage is obtained by calling the NominalGlyph() method of each font.
// See also SplitByFace for a more general approach of font selection.
func SplitByFontGlyphs(input Input, availableFaces []font.Face) []Input {
	return SplitByFace(input, fixedFontmap(availableFaces))
}

// SplitByFace split the runes from 'input' to several items, sharing the same
// characteristics as 'input', expected for the `Face` which is set to
// the return value of the `Fontmap.ResolveFace` call.
// The 'Face' field of 'input' is ignored: only 'availableFaces' is used to select the face.
func SplitByFace(input Input, availableFaces Fontmap) []Input {
	var splitInputs []Input
	currentInput := input
	for i := input.RunStart; i < input.RunEnd; i++ {
		r := input.Text[i]
		if currentInput.Face != nil && ignoreFaceChange(r) {
			// add the rune to the current input
			continue
		}

		// select the first font supporting r
		selectedFace := availableFaces.ResolveFace(r)

		if currentInput.Face == selectedFace {
			// add the rune to the current input
			continue
		}

		// new face needed

		if i != input.RunStart {
			// close the current input ...
			currentInput.RunEnd = i
			// ... add it to the output ...
			splitInputs = append(splitInputs, currentInput)
		}

		// ... and create a new one
		currentInput = input
		currentInput.RunStart = i
		currentInput.Face = selectedFace
	}

	// close and add the last input
	currentInput.RunEnd = input.RunEnd
	splitInputs = append(splitInputs, currentInput)
	return splitInputs
}

// ignoreFaceChange returns `true` is the given rune should not trigger
// a change of font.
//
// We don't want space characters to affect font selection; in general,
// it's always wrong to select a font just to render a space.
// We assume that all fonts have the ASCII space, and for other space
// characters if they don't, HarfBuzz will compatibility-decompose them
// to ASCII space...
//
// We don't want to change fonts for line or paragraph separators.
//
// Finaly, we also don't change fonts for what Harfbuzz consider
// as ignorable (however, some Control Format runes like 06DD are not ignored).
//
// The rationale is taken from pango : see bugs
// https://bugzilla.gnome.org/show_bug.cgi?id=355987
// https://bugzilla.gnome.org/show_bug.cgi?id=701652
// https://bugzilla.gnome.org/show_bug.cgi?id=781123
// for more details.
func ignoreFaceChange(r rune) bool {
	return unicode.Is(unicode.Cc, r) || // control
		unicode.Is(unicode.Cs, r) || // surrogate
		unicode.Is(unicode.Zl, r) || // line separator
		unicode.Is(unicode.Zp, r) || // paragraph separator
		(unicode.Is(unicode.Zs, r) && r != '\u1680') || // space separator != OGHAM SPACE MARK
		harfbuzz.IsDefaultIgnorable(r)
}
