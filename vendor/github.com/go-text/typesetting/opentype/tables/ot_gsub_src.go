// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

type SingleSubs struct {
	Data SingleSubstData
}

type SingleSubstData interface {
	isSingleSubstData()

	Cov() Coverage
}

func (SingleSubstData1) isSingleSubstData() {}
func (SingleSubstData2) isSingleSubstData() {}

type SingleSubstData1 struct {
	format       uint16   `unionTag:"1"`
	Coverage     Coverage `offsetSize:"Offset16"` // Offset to Coverage table, from beginning of substitution subtable
	DeltaGlyphID int16    // Add to original glyph ID to get substitute glyph ID
}

type SingleSubstData2 struct {
	format             uint16    `unionTag:"2"`
	Coverage           Coverage  `offsetSize:"Offset16"`    // Offset to Coverage table, from beginning of substitution subtable
	SubstituteGlyphIDs []GlyphID `arrayCount:"FirstUint16"` //[glyphCount]	Array of substitute glyph IDs — ordered by Coverage index
}

type MultipleSubs struct {
	substFormat uint16     // Format identifier: format = 1
	Coverage    Coverage   `offsetSize:"Offset16"` // Offset to Coverage table, from beginning of substitution subtable
	Sequences   []Sequence `arrayCount:"FirstUint16"  offsetsArray:"Offset16"`
	//[sequenceCount]	Array of offsets to Sequence tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
}

type Sequence struct {
	SubstituteGlyphIDs []GlyphID `arrayCount:"FirstUint16"` // [glyphCount]	String of glyph IDs to substitute
}

type AlternateSubs struct {
	substFormat   uint16         //	Format identifier: format = 1
	Coverage      Coverage       `offsetSize:"Offset16"` //	Offset to Coverage table, from beginning of substitution subtable
	AlternateSets []AlternateSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"`
}

type AlternateSet struct {
	AlternateGlyphIDs []GlyphID `arrayCount:"FirstUint16"` // Array of alternate glyph IDs, in arbitrary order
}

type LigatureSubs struct {
	substFormat  uint16        // Format identifier: format = 1
	Coverage     Coverage      `offsetSize:"Offset16"`                             // Offset to Coverage table, from beginning of substitution subtable
	LigatureSets []LigatureSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[ligatureSetCount]	Array of offsets to LigatureSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
}

// All ligatures beginning with the same glyph
type LigatureSet struct {
	Ligatures []Ligature `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` // [LigatureCount]	Array of offsets to Ligature tables. Offsets are from beginning of LigatureSet table, ordered by preference.
}

// Glyph components for one ligature
type Ligature struct {
	LigatureGlyph     GlyphID   //	glyph ID of ligature to substitute
	componentCount    uint16    //	Number of components in the ligature
	ComponentGlyphIDs []GlyphID `arrayCount:"ComputedField-componentCount-1"` //  [componentCount - 1]	Array of component glyph IDs — start with the second component, ordered in writing direction
}

type ContextualSubs struct {
	Data ContextualSubsITF
}

type ContextualSubsITF interface {
	isContextualSubsITF()

	Cov() Coverage
}

type (
	ContextualSubs1 SequenceContextFormat1
	ContextualSubs2 SequenceContextFormat2
	ContextualSubs3 SequenceContextFormat3
)

func (ContextualSubs1) isContextualSubsITF() {}
func (ContextualSubs2) isContextualSubsITF() {}
func (ContextualSubs3) isContextualSubsITF() {}

type ChainedContextualSubs struct {
	Data ChainedContextualSubsITF
}

type ChainedContextualSubsITF interface {
	isChainedContextualSubsITF()

	Cov() Coverage
}

type (
	ChainedContextualSubs1 ChainedSequenceContextFormat1
	ChainedContextualSubs2 ChainedSequenceContextFormat2
	ChainedContextualSubs3 ChainedSequenceContextFormat3
)

func (ChainedContextualSubs1) isChainedContextualSubsITF() {}
func (ChainedContextualSubs2) isChainedContextualSubsITF() {}
func (ChainedContextualSubs3) isChainedContextualSubsITF() {}

type ExtensionSubs Extension

type ReverseChainSingleSubs struct {
	substFormat        uint16     // Format identifier: format = 1
	coverage           Coverage   `offsetSize:"Offset16"`                             // Offset to Coverage table, from beginning of substitution subtable.
	BacktrackCoverages []Coverage `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[backtrackGlyphCount]	Array of offsets to coverage tables in backtrack sequence, in glyph sequence order.
	LookaheadCoverages []Coverage `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[lookaheadGlyphCount]	Array of offsets to coverage tables in lookahead sequence, in glyph sequence order.
	SubstituteGlyphIDs []GlyphID  `arrayCount:"FirstUint16"`                          //[glyphCount]	Array of substitute glyph IDs — ordered by Coverage index.
}
