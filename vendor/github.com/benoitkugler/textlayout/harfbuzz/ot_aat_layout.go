package harfbuzz

import (
	"fmt"

	"github.com/benoitkugler/textlayout/fonts"
	tt "github.com/benoitkugler/textlayout/fonts/truetype"
)

// ported from harfbuzz/src/hb-aat-layout.h  Copyright © 2018 Ebrahim Byagowi, Behdad Esfahbod

// The possible feature types defined for AAT shaping,
// from https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html
type aatLayoutFeatureType = uint16

const (
	// Initial, unset feature type
	// aatLayoutFeatureTypeInvalid = 0xFFFF
	// [All Typographic Features](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type0)
	// aatLayoutFeatureTypeAllTypographic = 0
	// [Ligatures](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type1)
	aatLayoutFeatureTypeLigatures = 1
	// [Cursive Connection](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type2)
	// aatLayoutFeatureTypeCurisveConnection = 2
	// [Letter Case](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type3)
	aatLayoutFeatureTypeLetterCase = 3
	// [Vertical Substitution](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type4)
	aatLayoutFeatureTypeVerticalSubstitution = 4
	// [Linguistic Rearrangement](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type5)
	// aatLayoutFeatureTypeLinguisticRearrangement = 5
	// [Number Spacing](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type6)
	aatLayoutFeatureTypeNumberSpacing = 6
	// [Smart Swash](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type8)
	// aatLayoutFeatureTypeSmartSwashType = 8
	// [Diacritics](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type9)
	// aatLayoutFeatureTypeDiacriticsType = 9
	// [Vertical Position](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type10)
	aatLayoutFeatureTypeVerticalPosition = 10
	// [Fractions](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type11)
	aatLayoutFeatureTypeFractions = 11
	// [Overlapping Characters](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type13)
	// aatLayoutFeatureTypeOverlappingCharactersType = 13
	// [Typographic Extras](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type14)
	aatLayoutFeatureTypeTypographicExtras = 14
	// [Mathematical Extras](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type15)
	aatLayoutFeatureTypeMathematicalExtras = 15
	// [Ornament Sets](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type16)
	// aatLayoutFeatureTypeOrnamentSetsType = 16
	// [Character Alternatives](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type17)
	aatLayoutFeatureTypeCharacterAlternatives = 17
	// [Design Complexity](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type18)
	aatLayoutFeatureTypeDesignComplexityType = 18
	// [Style Options](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type19)
	aatLayoutFeatureTypeStyleOptions = 19
	// [Character Shape](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type20)
	aatLayoutFeatureTypeCharacterShape = 20
	// [Number Case](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type21)
	aatLayoutFeatureTypeNumberCase = 21
	// [Text Spacing](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type22)
	aatLayoutFeatureTypeTextSpacing = 22
	// [Transliteration](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type23)
	aatLayoutFeatureTypeTransliteration = 23

	// [Ruby Kana](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type28)
	aatLayoutFeatureTypeRubyKana = 28

	// [Italic CJK Roman](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type32)
	aatLayoutFeatureTypeItalicCjkRoman = 32
	// [Case Sensitive Layout](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type33)
	aatLayoutFeatureTypeCaseSensitiveLayout = 33
	// [Alternate Kana](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type34)
	aatLayoutFeatureTypeAlternateKana = 34
	// [Stylistic Alternatives](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type35)
	aatLayoutFeatureTypeStylisticAlternatives = 35
	// [Contextual Alternatives](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type36)
	aatLayoutFeatureTypeContextualAlternatives = 36
	// [Lower Case](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type37)
	aatLayoutFeatureTypeLowerCase = 37
	// [Upper Case](https://developer.apple.com/fonts/TrueType-Reference-Manual/RM09/AppendixF.html#Type38)
	aatLayoutFeatureTypeUpperCase = 38
)

// The selectors defined for specifying AAT feature settings.
type aatLayoutFeatureSelector = uint16

const (
	/* Selectors for #aatLayoutFeatureTypeLigatures */
	// for #aatLayoutFeatureTypeLigatures
	aatLayoutFeatureSelectorCommonLigaturesOn = 2
	// for #aatLayoutFeatureTypeLigatures
	aatLayoutFeatureSelectorCommonLigaturesOff = 3
	// for #aatLayoutFeatureTypeLigatures
	aatLayoutFeatureSelectorRareLigaturesOn = 4
	// for #aatLayoutFeatureTypeLigatures
	aatLayoutFeatureSelectorRareLigaturesOff = 5

	// for #aatLayoutFeatureTypeLigatures
	aatLayoutFeatureSelectorContextualLigaturesOn = 18
	// for #aatLayoutFeatureTypeLigatures
	aatLayoutFeatureSelectorContextualLigaturesOff = 19
	// for #aatLayoutFeatureTypeLigatures
	aatLayoutFeatureSelectorHistoricalLigaturesOn = 20
	// for #aatLayoutFeatureTypeLigatures
	aatLayoutFeatureSelectorHistoricalLigaturesOff = 21

	/* Selectors for #aatLayoutFeatureTypeLetterCase */

	// Deprecated
	aatLayoutFeatureSelectorSmallCaps = 3 /* deprecated */

	/* Selectors for #aatLayoutFeatureTypeVerticalSubstitution */
	// for #aatLayoutFeatureTypeVerticalSubstitution
	aatLayoutFeatureSelectorSubstituteVerticalFormsOn = 0
	// for #aatLayoutFeatureTypeVerticalSubstitution
	aatLayoutFeatureSelectorSubstituteVerticalFormsOff = 1

	/* Selectors for #aatLayoutFeatureTypeNumberSpacing */
	// for #aatLayoutFeatureTypeNumberSpacing
	aatLayoutFeatureSelectorMonospacedNumbers = 0
	// for #aatLayoutFeatureTypeNumberSpacing
	aatLayoutFeatureSelectorProportionalNumbers = 1

	/* Selectors for #aatLayoutFeatureTypeVerticalPosition */
	// for #aatLayoutFeatureTypeVerticalPosition
	aatLayoutFeatureSelectorNormalPosition = 0
	// for #aatLayoutFeatureTypeVerticalPosition
	aatLayoutFeatureSelectorSuperiors = 1
	// for #aatLayoutFeatureTypeVerticalPosition
	aatLayoutFeatureSelectorInferiors = 2
	// for #aatLayoutFeatureTypeVerticalPosition
	aatLayoutFeatureSelectorOrdinals = 3
	// for #aatLayoutFeatureTypeVerticalPosition
	aatLayoutFeatureSelectorScientificInferiors = 4

	/* Selectors for #aatLayoutFeatureTypeFractions */
	// for #aatLayoutFeatureTypeFractions
	aatLayoutFeatureSelectorNoFractions = 0
	// for #aatLayoutFeatureTypeFractions
	aatLayoutFeatureSelectorVerticalFractions = 1
	// for #aatLayoutFeatureTypeFractions
	aatLayoutFeatureSelectorDiagonalFractions = 2

	// for #aatLayoutFeatureTypeTypographicExtras
	aatLayoutFeatureSelectorSlashedZeroOn = 4
	// for #aatLayoutFeatureTypeTypographicExtras
	aatLayoutFeatureSelectorSlashedZeroOff = 5

	/* Selectors for #aatLayoutFeatureTypeMathematicalExtras */
	// for #aatLayoutFeatureTypeMathematicalExtras
	aatLayoutFeatureSelectorMathematicalGreekOn = 10
	// for #aatLayoutFeatureTypeMathematicalExtras
	aatLayoutFeatureSelectorMathematicalGreekOff = 11

	/* Selectors for #aatLayoutFeatureTypeStyleOptions */
	// for #aatLayoutFeatureTypeStyleOptions
	aatLayoutFeatureSelectorNoStyleOptions = 0

	// for #aatLayoutFeatureTypeStyleOptions
	aatLayoutFeatureSelectorTitlingCaps = 4

	/* Selectors for #aatLayoutFeatureTypeCharacterShape */
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorTraditionalCharacters = 0
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorSimplifiedCharacters = 1
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorJis1978Characters = 2
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorJis1983Characters = 3
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorJis1990Characters = 4

	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorExpertCharacters = 10
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorJis2004Characters = 11
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorHojoCharacters = 12
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorNlccharacters = 13
	// for #aatLayoutFeatureTypeCharacterShape
	aatLayoutFeatureSelectorTraditionalNamesCharacters = 14

	/* Selectors for #aatLayoutFeatureTypeNumberCase */
	// for #aatLayoutFeatureTypeNumberCase
	aatLayoutFeatureSelectorLowerCaseNumbers = 0
	// for #aatLayoutFeatureTypeNumberCase
	aatLayoutFeatureSelectorUpperCaseNumbers = 1

	/* Selectors for #aatLayoutFeatureTypeTextSpacing */
	// for #aatLayoutFeatureTypeTextSpacing
	aatLayoutFeatureSelectorProportionalText = 0
	// for #aatLayoutFeatureTypeTextSpacing
	aatLayoutFeatureSelectorMonospacedText = 1
	// for #aatLayoutFeatureTypeTextSpacing
	aatLayoutFeatureSelectorHalfWidthText = 2
	// for #aatLayoutFeatureTypeTextSpacing
	aatLayoutFeatureSelectorThirdWidthText = 3
	// for #aatLayoutFeatureTypeTextSpacing
	aatLayoutFeatureSelectorQuarterWidthText = 4
	// for #aatLayoutFeatureTypeTextSpacing
	aatLayoutFeatureSelectorAltProportionalText = 5
	// for #aatLayoutFeatureTypeTextSpacing
	aatLayoutFeatureSelectorAltHalfWidthText = 6

	/* Selectors for #aatLayoutFeatureTypeTransliteration */
	// for #aatLayoutFeatureTypeTransliteration
	aatLayoutFeatureSelectorNoTransliteration = 0
	// for #aatLayoutFeatureTypeTransliteration
	aatLayoutFeatureSelectorHanjaToHangul = 1

	/* Selectors for #aatLayoutFeatureTypeRubyKana */
	// for #aatLayoutFeatureTypeRubyKana
	aatLayoutFeatureSelectorRubyKanaOn = 2
	// for #aatLayoutFeatureTypeRubyKana
	aatLayoutFeatureSelectorRubyKanaOff = 3

	/* Selectors for #aatLayoutFeatureTypeItalicCjkRoman */
	// for #aatLayoutFeatureTypeItalicCjkRoman
	aatLayoutFeatureSelectorCjkItalicRomanOn = 2
	// for #aatLayoutFeatureTypeItalicCjkRoman
	aatLayoutFeatureSelectorCjkItalicRomanOff = 3

	/* Selectors for #aatLayoutFeatureTypeCaseSensitiveLayout */
	// for #aatLayoutFeatureTypeCaseSensitiveLayout
	aatLayoutFeatureSelectorCaseSensitiveLayoutOn = 0
	// for #aatLayoutFeatureTypeCaseSensitiveLayout
	aatLayoutFeatureSelectorCaseSensitiveLayoutOff = 1
	// for #aatLayoutFeatureTypeCaseSensitiveLayout
	aatLayoutFeatureSelectorCaseSensitiveSpacingOn = 2
	// for #aatLayoutFeatureTypeCaseSensitiveLayout
	aatLayoutFeatureSelectorCaseSensitiveSpacingOff = 3

	/* Selectors for #aatLayoutFeatureTypeAlternateKana */
	// for #aatLayoutFeatureTypeAlternateKana
	aatLayoutFeatureSelectorAlternateHorizKanaOn = 0
	// for #aatLayoutFeatureTypeAlternateKana
	aatLayoutFeatureSelectorAlternateHorizKanaOff = 1
	// for #aatLayoutFeatureTypeAlternateKana
	aatLayoutFeatureSelectorAlternateVertKanaOn = 2
	// for #aatLayoutFeatureTypeAlternateKana
	aatLayoutFeatureSelectorAlternateVertKanaOff = 3

	/* Selectors for #aatLayoutFeatureTypeStylisticAlternatives */

	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltOneOn = 2
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltOneOff = 3
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltTwoOn = 4
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltTwoOff = 5
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltThreeOn = 6
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltThreeOff = 7
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltFourOn = 8
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltFourOff = 9
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltFiveOn = 10
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltFiveOff = 11
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltSixOn = 12
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltSixOff = 13
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltSevenOn = 14
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltSevenOff = 15
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltEightOn = 16
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltEightOff = 17
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltNineOn = 18
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltNineOff = 19
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltTenOn = 20
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltTenOff = 21
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltElevenOn = 22
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltElevenOff = 23
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltTwelveOn = 24
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltTwelveOff = 25
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltThirteenOn = 26
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltThirteenOff = 27
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltFourteenOn = 28
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltFourteenOff = 29
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltFifteenOn = 30
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltFifteenOff = 31
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltSixteenOn = 32
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltSixteenOff = 33
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltSeventeenOn = 34
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltSeventeenOff = 35
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltEighteenOn = 36
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltEighteenOff = 37
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltNineteenOn = 38
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltNineteenOff = 39
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltTwentyOn = 40
	// for #aatLayoutFeatureTypeStylisticAlternatives
	aatLayoutFeatureSelectorStylisticAltTwentyOff = 41

	/* Selectors for #aatLayoutFeatureTypeContextualAlternatives */
	// for #aatLayoutFeatureTypeContextualAlternatives
	aatLayoutFeatureSelectorContextualAlternatesOn = 0
	// for #aatLayoutFeatureTypeContextualAlternatives
	aatLayoutFeatureSelectorContextualAlternatesOff = 1
	// for #aatLayoutFeatureTypeContextualAlternatives
	aatLayoutFeatureSelectorSwashAlternatesOn = 2
	// for #aatLayoutFeatureTypeContextualAlternatives
	aatLayoutFeatureSelectorSwashAlternatesOff = 3
	// for #aatLayoutFeatureTypeContextualAlternatives
	aatLayoutFeatureSelectorContextualSwashAlternatesOn = 4
	// for #aatLayoutFeatureTypeContextualAlternatives
	aatLayoutFeatureSelectorContextualSwashAlternatesOff = 5

	/* Selectors for #aatLayoutFeatureTypeLowerCase */
	// for #aatLayoutFeatureTypeLowerCase
	aatLayoutFeatureSelectorDefaultLowerCase = 0
	// for #aatLayoutFeatureTypeLowerCase
	aatLayoutFeatureSelectorLowerCaseSmallCaps = 1
	// for #aatLayoutFeatureTypeLowerCase
	aatLayoutFeatureSelectorLowerCasePetiteCaps = 2

	/* Selectors for #aatLayoutFeatureTypeUpperCase */
	// for #aatLayoutFeatureTypeUpperCase
	aatLayoutFeatureSelectorDefaultUpperCase = 0
	// for #aatLayoutFeatureTypeUpperCase
	aatLayoutFeatureSelectorUpperCaseSmallCaps = 1
	// for #aatLayoutFeatureTypeUpperCase
	aatLayoutFeatureSelectorUpperCasePetiteCaps = 2
)

/* Mapping from OpenType feature tags to AAT feature names and selectors.
 *
 * Table data courtesy of Apple.  Converted from mnemonics to integers
 * when moving to this file. */
var featureMappings = [...]aatFeatureMapping{
	{tt.NewTag('a', 'f', 'r', 'c'), aatLayoutFeatureTypeFractions, aatLayoutFeatureSelectorVerticalFractions, aatLayoutFeatureSelectorNoFractions},
	{tt.NewTag('c', '2', 'p', 'c'), aatLayoutFeatureTypeUpperCase, aatLayoutFeatureSelectorUpperCasePetiteCaps, aatLayoutFeatureSelectorDefaultUpperCase},
	{tt.NewTag('c', '2', 's', 'c'), aatLayoutFeatureTypeUpperCase, aatLayoutFeatureSelectorUpperCaseSmallCaps, aatLayoutFeatureSelectorDefaultUpperCase},
	{tt.NewTag('c', 'a', 'l', 't'), aatLayoutFeatureTypeContextualAlternatives, aatLayoutFeatureSelectorContextualAlternatesOn, aatLayoutFeatureSelectorContextualAlternatesOff},
	{tt.NewTag('c', 'a', 's', 'e'), aatLayoutFeatureTypeCaseSensitiveLayout, aatLayoutFeatureSelectorCaseSensitiveLayoutOn, aatLayoutFeatureSelectorCaseSensitiveLayoutOff},
	{tt.NewTag('c', 'l', 'i', 'g'), aatLayoutFeatureTypeLigatures, aatLayoutFeatureSelectorContextualLigaturesOn, aatLayoutFeatureSelectorContextualLigaturesOff},
	{tt.NewTag('c', 'p', 's', 'p'), aatLayoutFeatureTypeCaseSensitiveLayout, aatLayoutFeatureSelectorCaseSensitiveSpacingOn, aatLayoutFeatureSelectorCaseSensitiveSpacingOff},
	{tt.NewTag('c', 's', 'w', 'h'), aatLayoutFeatureTypeContextualAlternatives, aatLayoutFeatureSelectorContextualSwashAlternatesOn, aatLayoutFeatureSelectorContextualSwashAlternatesOff},
	{tt.NewTag('d', 'l', 'i', 'g'), aatLayoutFeatureTypeLigatures, aatLayoutFeatureSelectorRareLigaturesOn, aatLayoutFeatureSelectorRareLigaturesOff},
	{tt.NewTag('e', 'x', 'p', 't'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorExpertCharacters, 16},
	{tt.NewTag('f', 'r', 'a', 'c'), aatLayoutFeatureTypeFractions, aatLayoutFeatureSelectorDiagonalFractions, aatLayoutFeatureSelectorNoFractions},
	{tt.NewTag('f', 'w', 'i', 'd'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorMonospacedText, 7},
	{tt.NewTag('h', 'a', 'l', 't'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorAltHalfWidthText, 7},
	{tt.NewTag('h', 'i', 's', 't'), aatLayoutFeatureTypeLigatures, aatLayoutFeatureSelectorHistoricalLigaturesOn, aatLayoutFeatureSelectorHistoricalLigaturesOff},
	{tt.NewTag('h', 'k', 'n', 'a'), aatLayoutFeatureTypeAlternateKana, aatLayoutFeatureSelectorAlternateHorizKanaOn, aatLayoutFeatureSelectorAlternateHorizKanaOff},
	{tt.NewTag('h', 'l', 'i', 'g'), aatLayoutFeatureTypeLigatures, aatLayoutFeatureSelectorHistoricalLigaturesOn, aatLayoutFeatureSelectorHistoricalLigaturesOff},
	{tt.NewTag('h', 'n', 'g', 'l'), aatLayoutFeatureTypeTransliteration, aatLayoutFeatureSelectorHanjaToHangul, aatLayoutFeatureSelectorNoTransliteration},
	{tt.NewTag('h', 'o', 'j', 'o'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorHojoCharacters, 16},
	{tt.NewTag('h', 'w', 'i', 'd'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorHalfWidthText, 7},
	{tt.NewTag('i', 't', 'a', 'l'), aatLayoutFeatureTypeItalicCjkRoman, aatLayoutFeatureSelectorCjkItalicRomanOn, aatLayoutFeatureSelectorCjkItalicRomanOff},
	{tt.NewTag('j', 'p', '0', '4'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorJis2004Characters, 16},
	{tt.NewTag('j', 'p', '7', '8'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorJis1978Characters, 16},
	{tt.NewTag('j', 'p', '8', '3'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorJis1983Characters, 16},
	{tt.NewTag('j', 'p', '9', '0'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorJis1990Characters, 16},
	{tt.NewTag('l', 'i', 'g', 'a'), aatLayoutFeatureTypeLigatures, aatLayoutFeatureSelectorCommonLigaturesOn, aatLayoutFeatureSelectorCommonLigaturesOff},
	{tt.NewTag('l', 'n', 'u', 'm'), aatLayoutFeatureTypeNumberCase, aatLayoutFeatureSelectorUpperCaseNumbers, 2},
	{tt.NewTag('m', 'g', 'r', 'k'), aatLayoutFeatureTypeMathematicalExtras, aatLayoutFeatureSelectorMathematicalGreekOn, aatLayoutFeatureSelectorMathematicalGreekOff},
	{tt.NewTag('n', 'l', 'c', 'k'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorNlccharacters, 16},
	{tt.NewTag('o', 'n', 'u', 'm'), aatLayoutFeatureTypeNumberCase, aatLayoutFeatureSelectorLowerCaseNumbers, 2},
	{tt.NewTag('o', 'r', 'd', 'n'), aatLayoutFeatureTypeVerticalPosition, aatLayoutFeatureSelectorOrdinals, aatLayoutFeatureSelectorNormalPosition},
	{tt.NewTag('p', 'a', 'l', 't'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorAltProportionalText, 7},
	{tt.NewTag('p', 'c', 'a', 'p'), aatLayoutFeatureTypeLowerCase, aatLayoutFeatureSelectorLowerCasePetiteCaps, aatLayoutFeatureSelectorDefaultLowerCase},
	{tt.NewTag('p', 'k', 'n', 'a'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorProportionalText, 7},
	{tt.NewTag('p', 'n', 'u', 'm'), aatLayoutFeatureTypeNumberSpacing, aatLayoutFeatureSelectorProportionalNumbers, 4},
	{tt.NewTag('p', 'w', 'i', 'd'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorProportionalText, 7},
	{tt.NewTag('q', 'w', 'i', 'd'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorQuarterWidthText, 7},
	{tt.NewTag('r', 'u', 'b', 'y'), aatLayoutFeatureTypeRubyKana, aatLayoutFeatureSelectorRubyKanaOn, aatLayoutFeatureSelectorRubyKanaOff},
	{tt.NewTag('s', 'i', 'n', 'f'), aatLayoutFeatureTypeVerticalPosition, aatLayoutFeatureSelectorScientificInferiors, aatLayoutFeatureSelectorNormalPosition},
	{tt.NewTag('s', 'm', 'c', 'p'), aatLayoutFeatureTypeLowerCase, aatLayoutFeatureSelectorLowerCaseSmallCaps, aatLayoutFeatureSelectorDefaultLowerCase},
	{tt.NewTag('s', 'm', 'p', 'l'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorSimplifiedCharacters, 16},
	{tt.NewTag('s', 's', '0', '1'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltOneOn, aatLayoutFeatureSelectorStylisticAltOneOff},
	{tt.NewTag('s', 's', '0', '2'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltTwoOn, aatLayoutFeatureSelectorStylisticAltTwoOff},
	{tt.NewTag('s', 's', '0', '3'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltThreeOn, aatLayoutFeatureSelectorStylisticAltThreeOff},
	{tt.NewTag('s', 's', '0', '4'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltFourOn, aatLayoutFeatureSelectorStylisticAltFourOff},
	{tt.NewTag('s', 's', '0', '5'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltFiveOn, aatLayoutFeatureSelectorStylisticAltFiveOff},
	{tt.NewTag('s', 's', '0', '6'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltSixOn, aatLayoutFeatureSelectorStylisticAltSixOff},
	{tt.NewTag('s', 's', '0', '7'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltSevenOn, aatLayoutFeatureSelectorStylisticAltSevenOff},
	{tt.NewTag('s', 's', '0', '8'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltEightOn, aatLayoutFeatureSelectorStylisticAltEightOff},
	{tt.NewTag('s', 's', '0', '9'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltNineOn, aatLayoutFeatureSelectorStylisticAltNineOff},
	{tt.NewTag('s', 's', '1', '0'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltTenOn, aatLayoutFeatureSelectorStylisticAltTenOff},
	{tt.NewTag('s', 's', '1', '1'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltElevenOn, aatLayoutFeatureSelectorStylisticAltElevenOff},
	{tt.NewTag('s', 's', '1', '2'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltTwelveOn, aatLayoutFeatureSelectorStylisticAltTwelveOff},
	{tt.NewTag('s', 's', '1', '3'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltThirteenOn, aatLayoutFeatureSelectorStylisticAltThirteenOff},
	{tt.NewTag('s', 's', '1', '4'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltFourteenOn, aatLayoutFeatureSelectorStylisticAltFourteenOff},
	{tt.NewTag('s', 's', '1', '5'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltFifteenOn, aatLayoutFeatureSelectorStylisticAltFifteenOff},
	{tt.NewTag('s', 's', '1', '6'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltSixteenOn, aatLayoutFeatureSelectorStylisticAltSixteenOff},
	{tt.NewTag('s', 's', '1', '7'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltSeventeenOn, aatLayoutFeatureSelectorStylisticAltSeventeenOff},
	{tt.NewTag('s', 's', '1', '8'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltEighteenOn, aatLayoutFeatureSelectorStylisticAltEighteenOff},
	{tt.NewTag('s', 's', '1', '9'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltNineteenOn, aatLayoutFeatureSelectorStylisticAltNineteenOff},
	{tt.NewTag('s', 's', '2', '0'), aatLayoutFeatureTypeStylisticAlternatives, aatLayoutFeatureSelectorStylisticAltTwentyOn, aatLayoutFeatureSelectorStylisticAltTwentyOff},
	{tt.NewTag('s', 'u', 'b', 's'), aatLayoutFeatureTypeVerticalPosition, aatLayoutFeatureSelectorInferiors, aatLayoutFeatureSelectorNormalPosition},
	{tt.NewTag('s', 'u', 'p', 's'), aatLayoutFeatureTypeVerticalPosition, aatLayoutFeatureSelectorSuperiors, aatLayoutFeatureSelectorNormalPosition},
	{tt.NewTag('s', 'w', 's', 'h'), aatLayoutFeatureTypeContextualAlternatives, aatLayoutFeatureSelectorSwashAlternatesOn, aatLayoutFeatureSelectorSwashAlternatesOff},
	{tt.NewTag('t', 'i', 't', 'l'), aatLayoutFeatureTypeStyleOptions, aatLayoutFeatureSelectorTitlingCaps, aatLayoutFeatureSelectorNoStyleOptions},
	{tt.NewTag('t', 'n', 'a', 'm'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorTraditionalNamesCharacters, 16},
	{tt.NewTag('t', 'n', 'u', 'm'), aatLayoutFeatureTypeNumberSpacing, aatLayoutFeatureSelectorMonospacedNumbers, 4},
	{tt.NewTag('t', 'r', 'a', 'd'), aatLayoutFeatureTypeCharacterShape, aatLayoutFeatureSelectorTraditionalCharacters, 16},
	{tt.NewTag('t', 'w', 'i', 'd'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorThirdWidthText, 7},
	{tt.NewTag('u', 'n', 'i', 'c'), aatLayoutFeatureTypeLetterCase, 14, 15},
	{tt.NewTag('v', 'a', 'l', 't'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorAltProportionalText, 7},
	{tt.NewTag('v', 'e', 'r', 't'), aatLayoutFeatureTypeVerticalSubstitution, aatLayoutFeatureSelectorSubstituteVerticalFormsOn, aatLayoutFeatureSelectorSubstituteVerticalFormsOff},
	{tt.NewTag('v', 'h', 'a', 'l'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorAltHalfWidthText, 7},
	{tt.NewTag('v', 'k', 'n', 'a'), aatLayoutFeatureTypeAlternateKana, aatLayoutFeatureSelectorAlternateVertKanaOn, aatLayoutFeatureSelectorAlternateVertKanaOff},
	{tt.NewTag('v', 'p', 'a', 'l'), aatLayoutFeatureTypeTextSpacing, aatLayoutFeatureSelectorAltProportionalText, 7},
	{tt.NewTag('v', 'r', 't', '2'), aatLayoutFeatureTypeVerticalSubstitution, aatLayoutFeatureSelectorSubstituteVerticalFormsOn, aatLayoutFeatureSelectorSubstituteVerticalFormsOff},
	{tt.NewTag('z', 'e', 'r', 'o'), aatLayoutFeatureTypeTypographicExtras, aatLayoutFeatureSelectorSlashedZeroOn, aatLayoutFeatureSelectorSlashedZeroOff},
}

// Used when getting or setting AAT feature selectors. Indicates that
// there is no selector index corresponding to the selector of interest.
const aatLayoutNoSelectorIndex = 0xFFFF

/* Note: This context is used for kerning, even without AAT, hence the condition. */

/**
 * SECTION:hb-aat-layout
 * @title: hb-aat-layout
 * @short_description: Apple Advanced Typography Layout
 * @include: hb-aat.h
 *
 * Functions for querying AAT Layout features in the font face.
 *
 * HarfBuzz supports all of the AAT tables (in their modern version) used to implement shaping. Other
 * AAT tables and their associated features are not supported.
 **/

// execute the state machine in AAT tables
type stateTableDriver struct {
	buffer  *Buffer
	machine tt.AATStateTable
}

func newStateTableDriver(machine tt.AATStateTable, buffer *Buffer, face fonts.FaceMetrics) stateTableDriver {
	return stateTableDriver{
		machine: machine,
		buffer:  buffer,
	}
}

// implemented by the subtables
type driverContext interface {
	inPlace() bool
	isActionable(s stateTableDriver, entry tt.AATStateEntry) bool
	transition(s stateTableDriver, entry tt.AATStateEntry)
}

func (s stateTableDriver) drive(c driverContext) {
	const (
		stateStartOfText = uint16(0)

		classEndOfText = uint16(0)

		DontAdvance = 0x4000
	)
	if !c.inPlace() {
		s.buffer.clearOutput()
	}

	state := stateStartOfText
	for s.buffer.idx = 0; ; {
		class := classEndOfText
		if s.buffer.idx < len(s.buffer.Info) {
			class = s.machine.GetClass(s.buffer.Info[s.buffer.idx].Glyph)
		}

		if debugMode >= 2 {
			fmt.Printf("\t\tState machine - state %d, class %d at index %d\n", state, class, s.buffer.idx)
		}

		entry := s.machine.GetEntry(state, class)
		nextState := entry.NewState // we only supported extended table

		/* Conditions under which it's guaranteed safe-to-break before current glyph:
		 *
		 * 1. There was no action in this transition; and
		 *
		 * 2. If we break before current glyph, the results will be the same. That
		 *    is guaranteed if:
		 *
		 *    2a. We were already in start-of-text state; or
		 *
		 *    2b. We are epsilon-transitioning to start-of-text state; or
		 *
		 *    2c. Starting from start-of-text state seeing current glyph:
		 *
		 *        2c'. There won't be any actions; and
		 *
		 *        2c". We would end up in the same state that we were going to end up
		 *             in now, including whether epsilon-transitioning.
		 *
		 *    and
		 *
		 * 3. If we break before current glyph, there won't be any end-of-text action
		 *    after previous glyph.
		 *
		 * This triples the transitions we need to look up, but is worth returning
		 * granular unsafe-to-break results. See eg.:
		 *
		 *   https://github.com/harfbuzz/harfbuzz/issues/2860
		 */

		wouldbeEntry := s.machine.GetEntry(stateStartOfText, class)
		safeToBreak := /* 1. */ !c.isActionable(s, entry) &&
			/* 2. */
			(
			/* 2a. */
			state == stateStartOfText ||
				/* 2b. */
				((entry.Flags&DontAdvance != 0) && nextState == stateStartOfText) ||
				/* 2c. */
				(
				/* 2c'. */
				!c.isActionable(s, wouldbeEntry) &&
					/* 2c". */
					(nextState == wouldbeEntry.NewState) &&
					(entry.Flags&DontAdvance) == (wouldbeEntry.Flags&DontAdvance))) &&
			/* 3. */
			!c.isActionable(s, s.machine.GetEntry(state, classEndOfText))

		if !safeToBreak && s.buffer.backtrackLen() != 0 && s.buffer.idx < len(s.buffer.Info) {
			s.buffer.unsafeToBreakFromOutbuffer(s.buffer.backtrackLen()-1, s.buffer.idx+1)
		}

		c.transition(s, entry)

		state = nextState

		if debugMode >= 2 {
			fmt.Printf("\t\tState machine - new state %d\n", state)
		}

		if s.buffer.idx == len(s.buffer.Info) {
			break
		}

		if entry.Flags&DontAdvance == 0 {
			s.buffer.nextGlyph()
		} else {
			if s.buffer.maxOps <= 0 {
				s.buffer.maxOps--
				s.buffer.nextGlyph()
			}
			s.buffer.maxOps--
		}
	}

	if !c.inPlace() {
		s.buffer.swapBuffers()
	}
}

type aatApplyContext struct {
	plan      *otShapePlan
	font      *Font
	face      fonts.FaceMetrics
	buffer    *Buffer
	gdefTable *tt.TableGDEF
	ankrTable tt.TableAnkr
}

func newAatApplyContext(plan *otShapePlan, font *Font, buffer *Buffer) *aatApplyContext {
	var out aatApplyContext
	out.plan = plan
	out.font = font
	out.face = font.face
	out.buffer = buffer
	out.gdefTable = &font.otTables.GDEF
	return &out
}

func (c *aatApplyContext) applyMorx(chain tt.MorxChain, flags GlyphMask) {
	//  Coverage, see https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6morx.html
	const (
		Vertical      = 0x80
		Backwards     = 0x40
		AllDirections = 0x20
		Logical       = 0x10
	)

	for i, subtable := range chain.Subtables {

		if subtable.Flags&flags == 0 {
			continue
		}

		if subtable.Coverage&AllDirections == 0 && c.buffer.Props.Direction.isVertical() !=
			(subtable.Coverage&Vertical != 0) {
			continue
		}

		/* Buffer contents is always in logical direction.  Determine if
		we need to reverse before applying this subtable.  We reverse
		back after if we did reverse indeed.

		Quoting the spec:
		"""
		Bits 28 and 30 of the coverage field control the order in which
		glyphs are processed when the subtable is run by the layout engine.
		Bit 28 is used to indicate if the glyph processing direction is
		the same as logical order or layout order. Bit 30 is used to
		indicate whether glyphs are processed forwards or backwards within
		that order.

		Bit 30	Bit 28	Interpretation for Horizontal Text
		0	0	The subtable is processed in layout order 	(the same order as the glyphs, which is
			always left-to-right).
		1	0	The subtable is processed in reverse layout order (the order opposite that of the glyphs, which is
			always right-to-left).
		0	1	The subtable is processed in logical order (the same order as the characters, which may be
			left-to-right or right-to-left).
		1	1	The subtable is processed in reverse logical order 	(the order opposite that of the characters, which
			may be right-to-left or left-to-right).
		"""
		*/
		var reverse bool
		if subtable.Coverage&Logical != 0 {
			reverse = subtable.Coverage&Backwards != 0
		} else {
			reverse = subtable.Coverage&Backwards != 0 != c.buffer.Props.Direction.isBackward()
		}

		if debugMode >= 2 {
			fmt.Printf("MORX - start chainsubtable %d\n", i)
		}

		if reverse {
			reverseGraphemes(c.buffer)
		}

		c.applyMorxSubtable(subtable)

		if reverse {
			reverseGraphemes(c.buffer)
		}

		if debugMode >= 2 {
			fmt.Printf("MORX - end chainsubtable %d\n", i)
			fmt.Println(c.buffer.Info)
		}

	}
}

func (c *aatApplyContext) applyMorxSubtable(subtable tt.MortxSubtable) bool {
	if debugMode >= 2 {
		fmt.Printf("\tMORX subtable %T\n", subtable.Data)
	}
	switch data := subtable.Data.(type) {
	case tt.MorxRearrangementSubtable:
		var dc driverContextRearrangement
		driver := newStateTableDriver(tt.AATStateTable(data), c.buffer, c.face)
		driver.drive(&dc)
	case tt.MorxContextualSubtable:
		dc := driverContextContextual{table: data, gdef: c.gdefTable, hasGlyphClass: c.gdefTable.Class != nil}
		driver := newStateTableDriver(data.Machine, c.buffer, c.face)
		driver.drive(&dc)
		return dc.ret
	case tt.MorxLigatureSubtable:
		dc := driverContextLigature{table: data}
		driver := newStateTableDriver(data.Machine, c.buffer, c.face)
		driver.drive(&dc)
	case tt.MorxInsertionSubtable:
		dc := driverContextInsertion{insertionAction: data.Insertions}
		driver := newStateTableDriver(data.Machine, c.buffer, c.face)
		driver.drive(&dc)
	case tt.MorxNonContextualSubtable:
		var ret bool
		gdef := c.gdefTable
		hasGlyphClass := gdef.Class != nil
		info := c.buffer.Info
		for i := range c.buffer.Info {
			replacement, has := data.ClassID(info[i].Glyph)
			if has {
				info[i].Glyph = fonts.GID(replacement)
				if hasGlyphClass {
					info[i].glyphProps = gdef.GetGlyphProps(fonts.GID(replacement))
				}
				ret = true
			}
		}
		return ret
	}
	return false
}

type driverContextRearrangement struct {
	start int
	end   int
}

func (driverContextRearrangement) inPlace() bool { return true }

func (d driverContextRearrangement) isActionable(_ stateTableDriver, entry tt.AATStateEntry) bool {
	return (entry.Flags&tt.MRVerb) != 0 && d.start < d.end
}

/* The following map has two nibbles, for start-side
 * and end-side. Values of 0,1,2 mean move that many
 * to the other side. Value of 3 means move 2 and
 * flip them. */
var mapRearrangement = [16]int{
	0x00, /* 0	no change */
	0x10, /* 1	Ax => xA */
	0x01, /* 2	xD => Dx */
	0x11, /* 3	AxD => DxA */
	0x20, /* 4	ABx => xAB */
	0x30, /* 5	ABx => xBA */
	0x02, /* 6	xCD => CDx */
	0x03, /* 7	xCD => DCx */
	0x12, /* 8	AxCD => CDxA */
	0x13, /* 9	AxCD => DCxA */
	0x21, /* 10	ABxD => DxAB */
	0x31, /* 11	ABxD => DxBA */
	0x22, /* 12	ABxCD => CDxAB */
	0x32, /* 13	ABxCD => CDxBA */
	0x23, /* 14	ABxCD => DCxAB */
	0x33, /* 15	ABxCD => DCxBA */
}

func (d *driverContextRearrangement) transition(driver stateTableDriver, entry tt.AATStateEntry) {
	buffer := driver.buffer
	flags := entry.Flags

	if flags&tt.MRMarkFirst != 0 {
		d.start = buffer.idx
	}

	if flags&tt.MRMarkLast != 0 {
		d.end = min(buffer.idx+1, len(buffer.Info))
	}

	if (flags&tt.MRVerb) != 0 && d.start < d.end {

		m := mapRearrangement[flags&tt.MRVerb]
		l := min(2, m>>4)
		r := min(2, m&0x0F)
		reverseL := m>>4 == 3
		reverseR := m&0x0F == 3

		if d.end-d.start >= l+r {
			buffer.mergeClusters(d.start, min(buffer.idx+1, len(buffer.Info)))
			buffer.mergeClusters(d.start, d.end)

			info := buffer.Info
			var buf [4]GlyphInfo

			copy(buf[:], info[d.start:d.start+l])
			copy(buf[2:], info[d.end-r:d.end])

			if l != r {
				copy(info[d.start+r:], info[d.start+l:d.end-r])
			}

			copy(info[d.start:d.start+r], buf[2:])
			copy(info[d.end-l:d.end], buf[:])
			if reverseL {
				buf[0] = info[d.end-1]
				info[d.end-1] = info[d.end-2]
				info[d.end-2] = buf[0]
			}
			if reverseR {
				buf[0] = info[d.start]
				info[d.start] = info[d.start+1]
				info[d.start+1] = buf[0]
			}
		}
	}
}

type driverContextContextual struct {
	gdef          *tt.TableGDEF
	table         tt.MorxContextualSubtable
	mark          int
	markSet       bool
	ret           bool
	hasGlyphClass bool // cached version from gdef
}

func (driverContextContextual) inPlace() bool { return true }

func (dc driverContextContextual) isActionable(driver stateTableDriver, entry tt.AATStateEntry) bool {
	buffer := driver.buffer

	if buffer.idx == len(buffer.Info) && !dc.markSet {
		return false
	}
	markIndex, currentIndex := entry.AsMorxContextual()
	return markIndex != 0xFFFF || currentIndex != 0xFFFF
}

func (dc *driverContextContextual) transition(driver stateTableDriver, entry tt.AATStateEntry) {
	buffer := driver.buffer

	/* Looks like CoreText applies neither mark nor current substitution for
	 * end-of-text if mark was not explicitly set. */
	if buffer.idx == len(buffer.Info) && !dc.markSet {
		return
	}

	var (
		replacement             uint32 // intepreted as GlyphIndex
		hasRep                  bool
		markIndex, currentIndex = entry.AsMorxContextual()
	)
	if markIndex != 0xFFFF {
		lookup := dc.table.Substitutions[markIndex]
		replacement, hasRep = lookup.ClassID(buffer.Info[dc.mark].Glyph)
	}
	if hasRep {
		buffer.unsafeToBreak(dc.mark, min(buffer.idx+1, len(buffer.Info)))
		buffer.Info[dc.mark].Glyph = fonts.GID(replacement)
		if dc.hasGlyphClass {
			buffer.Info[dc.mark].glyphProps = dc.gdef.GetGlyphProps(fonts.GID(replacement))
		}
		dc.ret = true
	}

	hasRep = false
	idx := min(buffer.idx, len(buffer.Info)-1)
	if currentIndex != 0xFFFF {
		lookup := dc.table.Substitutions[currentIndex]
		replacement, hasRep = lookup.ClassID(buffer.Info[idx].Glyph)
	}

	if hasRep {
		buffer.Info[idx].Glyph = fonts.GID(replacement)
		if dc.hasGlyphClass {
			buffer.Info[idx].glyphProps = dc.gdef.GetGlyphProps(fonts.GID(replacement))
		}
		dc.ret = true
	}

	if entry.Flags&tt.MCSetMark != 0 {
		dc.markSet = true
		dc.mark = buffer.idx
	}
}

type driverContextLigature struct {
	table          tt.MorxLigatureSubtable
	matchLength    int
	matchPositions [maxContextLength]int
}

func (driverContextLigature) inPlace() bool { return false }

func (driverContextLigature) isActionable(_ stateTableDriver, entry tt.AATStateEntry) bool {
	return entry.Flags&tt.MLOffset != 0
}

func (dc *driverContextLigature) transition(driver stateTableDriver, entry tt.AATStateEntry) {
	buffer := driver.buffer

	if debugMode >= 2 {
		fmt.Printf("\tLigature - Ligature transition at %d\n", buffer.idx)
	}

	if entry.Flags&tt.MLSetComponent != 0 {
		/* Never mark same index twice, in case DontAdvance was used... */
		if dc.matchLength != 0 && dc.matchPositions[(dc.matchLength-1)%len(dc.matchPositions)] == len(buffer.outInfo) {
			dc.matchLength--
		}

		dc.matchPositions[dc.matchLength%len(dc.matchPositions)] = len(buffer.outInfo)
		dc.matchLength++

		if debugMode >= 2 {
			fmt.Printf("\tLigature - Set component at %d\n", len(buffer.outInfo))
		}

	}

	if dc.isActionable(driver, entry) {

		if debugMode >= 2 {
			fmt.Printf("\tLigature - Perform action with %d\n", dc.matchLength)
		}

		end := len(buffer.outInfo)

		if dc.matchLength == 0 {
			return
		}

		if buffer.idx >= len(buffer.Info) {
			return
		}
		cursor := dc.matchLength

		actionIdx := entry.AsMorxLigature()
		actionData := dc.table.LigatureAction[actionIdx:]

		ligatureIdx := 0
		var action uint32
		for do := true; do; do = action&tt.MLActionLast == 0 {
			if cursor == 0 {
				/* Stack underflow.  Clear the stack. */
				if debugMode >= 2 {
					fmt.Println("\tLigature - Stack underflow")
				}
				dc.matchLength = 0
				break
			}

			if debugMode >= 2 {
				fmt.Printf("\tLigature - Moving to stack position %d\n", cursor-1)
			}

			cursor--
			buffer.moveTo(dc.matchPositions[cursor%len(dc.matchPositions)])

			if len(actionData) == 0 {
				break
			}
			action = actionData[0]

			uoffset := action & tt.MLActionOffset
			if uoffset&0x20000000 != 0 {
				uoffset |= 0xC0000000 /* Sign-extend. */
			}
			offset := int32(uoffset)
			componentIdx := int32(buffer.cur(0).Glyph) + offset
			if int(componentIdx) >= len(dc.table.Component) {
				break
			}
			componentData := dc.table.Component[componentIdx]
			ligatureIdx += int(componentData)

			if debugMode >= 2 {
				fmt.Printf("\tLigature - Action store %d last %d\n", action&tt.MLActionStore, action&tt.MLActionLast)
			}

			if action&(tt.MLActionStore|tt.MLActionLast) != 0 {
				if ligatureIdx >= len(dc.table.Ligatures) {
					break
				}
				lig := dc.table.Ligatures[ligatureIdx]

				if debugMode >= 2 {
					fmt.Printf("\tLigature - Produced ligature %d\n", lig)
				}

				buffer.replaceGlyphIndex(lig)

				ligEnd := dc.matchPositions[(dc.matchLength-1)%len(dc.matchPositions)] + 1
				/* Now go and delete all subsequent components. */
				for dc.matchLength-1 > cursor {

					if debugMode >= 2 {
						fmt.Println("\tLigature - Skipping ligature component")
					}

					dc.matchLength--
					buffer.moveTo(dc.matchPositions[dc.matchLength%len(dc.matchPositions)])
					buffer.replaceGlyphIndex(0xFFFF)
				}

				buffer.moveTo(ligEnd)
				buffer.mergeOutClusters(dc.matchPositions[cursor%len(dc.matchPositions)], len(buffer.outInfo))
			}

			actionData = actionData[1:]
		}
		buffer.moveTo(end)
	}
}

type driverContextInsertion struct {
	insertionAction []fonts.GID
	mark            int
}

func (driverContextInsertion) inPlace() bool { return false }

func (driverContextInsertion) isActionable(_ stateTableDriver, entry tt.AATStateEntry) bool {
	current, marked := entry.AsMorxInsertion()
	return entry.Flags&(tt.MICurrentInsertCount|tt.MIMarkedInsertCount) != 0 && (current != 0xFFFF || marked != 0xFFFF)
}

func (dc *driverContextInsertion) transition(driver stateTableDriver, entry tt.AATStateEntry) {
	buffer := driver.buffer
	flags := entry.Flags

	markLoc := len(buffer.outInfo)
	currentInsertIndex, markedInsertIndex := entry.AsMorxInsertion()
	if markedInsertIndex != 0xFFFF {
		count := int(flags & tt.MIMarkedInsertCount)
		buffer.maxOps -= count
		if buffer.maxOps <= 0 {
			return
		}
		start := markedInsertIndex
		glyphs := dc.insertionAction[start:]

		before := flags&tt.MIMarkedInsertBefore != 0

		end := len(buffer.outInfo)
		buffer.moveTo(dc.mark)

		if buffer.idx < len(buffer.Info) && !before {
			buffer.copyGlyph()
		}
		/* TODO We ignore KashidaLike setting. */
		buffer.replaceGlyphs(0, nil, glyphs[:count])

		if buffer.idx < len(buffer.Info) && !before {
			buffer.skipGlyph()
		}

		buffer.moveTo(end + count)

		buffer.unsafeToBreakFromOutbuffer(dc.mark, min(buffer.idx+1, len(buffer.Info)))
	}

	if flags&tt.MISetMark != 0 {
		dc.mark = markLoc
	}

	if currentInsertIndex != 0xFFFF {
		count := int(flags&tt.MICurrentInsertCount) >> 5
		if buffer.maxOps <= 0 {
			buffer.maxOps -= count
			return
		}
		buffer.maxOps -= count
		start := currentInsertIndex
		glyphs := dc.insertionAction[start:]

		before := flags&tt.MICurrentInsertBefore != 0

		end := len(buffer.outInfo)

		if buffer.idx < len(buffer.Info) && !before {
			buffer.copyGlyph()
		}

		/* TODO We ignore KashidaLike setting. */
		buffer.replaceGlyphs(0, nil, glyphs[:count])

		if buffer.idx < len(buffer.Info) && !before {
			buffer.skipGlyph()
		}

		/* Humm. Not sure where to move to.  There's this wording under
		 * DontAdvance flag:
		 *
		 * "If set, don't update the glyph index before going to the new state.
		 * This does not mean that the glyph pointed to is the same one as
		 * before. If you've made insertions immediately downstream of the
		 * current glyph, the next glyph processed would in fact be the first
		 * one inserted."
		 *
		 * This suggests that if DontAdvance is NOT set, we should move to
		 * end+count.  If it *was*, then move to end, such that newly inserted
		 * glyphs are now visible.
		 *
		 * https://github.com/harfbuzz/harfbuzz/issues/1224#issuecomment-427691417
		 */
		moveTo := end
		if flags&tt.MIDontAdvance == 0 {
			moveTo = end + count
		}
		buffer.moveTo(moveTo)
	}
}

///////

type aatFeatureMapping struct {
	otFeatureTag      tt.Tag
	aatFeatureType    aatLayoutFeatureType
	selectorToEnable  aatLayoutFeatureSelector
	selectorToDisable aatLayoutFeatureSelector
}

// FaatLayoutFindFeatureMapping fetches the AAT feature-and-selector combination that corresponds
// to a given OpenType feature tag, or `nil` if not found.
func aatLayoutFindFeatureMapping(tag tt.Tag) *aatFeatureMapping {
	low, high := 0, len(featureMappings)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := featureMappings[mid].otFeatureTag
		if tag < p {
			high = mid
		} else if tag > p {
			low = mid + 1
		} else {
			return &featureMappings[mid]
		}
	}
	return nil
}

func (sp *otShapePlan) aatLayoutSubstitute(font *Font, buffer *Buffer) {
	morx := font.otTables.Morx
	c := newAatApplyContext(sp, font, buffer)
	for i, chain := range morx {
		c.applyMorx(chain, c.plan.aatMap.chainFlags[i])
	}
	// TODO: we dont support obsolete 'mort' table
}

func aatLayoutZeroWidthDeletedGlyphs(buffer *Buffer) {
	pos := buffer.Pos
	for i, inf := range buffer.Info {
		if inf.Glyph == 0xFFFF {
			pos[i].XAdvance, pos[i].YAdvance, pos[i].XOffset, pos[i].YOffset = 0, 0, 0, 0
		}
	}
}

func aatLayoutRemoveDeletedGlyphsInplace(buffer *Buffer) {
	otLayoutDeleteGlyphsInplace(buffer, func(info *GlyphInfo) bool {
		return info.Glyph == 0xFFFF
	})
}

func (sp *otShapePlan) aatLayoutPosition(font *Font, buffer *Buffer) {
	kerx := font.otTables.Kerx

	c := newAatApplyContext(sp, font, buffer)
	c.ankrTable = font.otTables.Ankr
	c.applyKernx(kerx)
}

func (c *aatApplyContext) applyKernx(kerx tt.TableKernx) {
	var ret, seenCrossStream bool

	for i, st := range kerx {
		var reverse bool

		if !st.IsExtended && st.IsVariation() {
			continue
		}

		if c.buffer.Props.Direction.isHorizontal() != st.IsHorizontal() {
			continue
		}
		reverse = st.IsBackwards() != c.buffer.Props.Direction.isBackward()

		if debugMode >= 2 {
			fmt.Printf("AAT kerx : start subtable %d\n", i)
		}

		if !seenCrossStream && st.IsCrossStream() {
			/* Attach all glyphs into a chain. */
			seenCrossStream = true
			pos := c.buffer.Pos
			for i := range pos {
				pos[i].attachType = attachTypeCursive
				if c.buffer.Props.Direction.isForward() {
					pos[i].attachChain = -1
				} else {
					pos[i].attachChain = +1
				}
				/* We intentionally don't set HB_BUFFER_SCRATCH_FLAG_HAS_GPOS_ATTACHMENT,
				 * since there needs to be a non-zero attachment for post-positioning to
				 * be needed. */
			}
		}

		if reverse {
			c.buffer.Reverse()
		}

		applied := c.applyKerxSubtable(st)
		ret = ret || applied

		if reverse {
			c.buffer.Reverse()
		}

		if debugMode >= 2 {
			fmt.Printf("AAT kerx : end subtable %d\n", i)
			fmt.Println(c.buffer.Pos)
		}

	}
}

func (c *aatApplyContext) applyKerxSubtable(st tt.KernSubtable) bool {
	if debugMode >= 2 {
		fmt.Printf("\tKERNX table %T\n", st.Data)
	}
	switch data := st.Data.(type) {
	case tt.Kern0:
		if !c.plan.requestedKerning {
			return false
		}
		if st.IsBackwards() {
			return false
		}
		kern(data, st.IsCrossStream(), c.font, c.buffer, c.plan.kernMask, true)
	case tt.Kern1:
		crossStream := st.IsCrossStream()
		if !c.plan.requestedKerning && !crossStream {
			return false
		}
		dc := driverContextKerx1{c: c, table: data, crossStream: crossStream}
		driver := newStateTableDriver(data.Machine, c.buffer, c.face)
		driver.drive(&dc)
	case tt.Kern2:
		if !c.plan.requestedKerning {
			return false
		}
		if st.IsBackwards() {
			return false
		}
		kern(data, st.IsCrossStream(), c.font, c.buffer, c.plan.kernMask, true)
	case tt.Kern3:
		if !c.plan.requestedKerning {
			return false
		}
		if st.IsBackwards() {
			return false
		}
		kern(data, st.IsCrossStream(), c.font, c.buffer, c.plan.kernMask, true)
	case tt.Kerx4:
		crossStream := st.IsCrossStream()
		if !c.plan.requestedKerning && !crossStream {
			return false
		}
		dc := driverContextKerx4{c: c, table: data, actionType: data.ActionType()}
		driver := newStateTableDriver(data.Machine, c.buffer, c.face)
		driver.drive(&dc)
	case tt.Kerx6:
		if !c.plan.requestedKerning {
			return false
		}
		if st.IsBackwards() {
			return false
		}
		kern(data, st.IsCrossStream(), c.font, c.buffer, c.plan.kernMask, true)
	}
	return true
}

type driverContextKerx1 struct {
	c           *aatApplyContext
	table       tt.Kern1
	stack       [8]int
	depth       int
	crossStream bool
}

func (driverContextKerx1) inPlace() bool { return true }

func (dc driverContextKerx1) isActionable(_ stateTableDriver, entry tt.AATStateEntry) bool {
	return entry.AsKernxIndex() != 0xFFFF
}

func (dc *driverContextKerx1) transition(driver stateTableDriver, entry tt.AATStateEntry) {
	buffer := driver.buffer
	flags := entry.Flags

	if flags&tt.Kerx1Reset != 0 {
		dc.depth = 0
	}

	if flags&tt.Kerx1Push != 0 {
		if dc.depth < len(dc.stack) {
			dc.stack[dc.depth] = buffer.idx
			dc.depth++
		} else {
			dc.depth = 0 /* Probably not what CoreText does, but better? */
		}
	}

	if dc.isActionable(driver, entry) && dc.depth != 0 {
		tupleCount := 1 // we do not support tupleCount > 0

		kernIdx := entry.AsKernxIndex()

		actions := dc.table.Values[kernIdx:]
		if len(actions) < tupleCount*dc.depth {
			dc.depth = 0
			return
		}

		kernMask := dc.c.plan.kernMask

		/* From Apple 'kern' spec:
		 * "Each pops one glyph from the kerning stack and applies the kerning value to it.
		 * The end of the list is marked by an odd value... */
		var last bool
		for !last && dc.depth != 0 {
			dc.depth--
			idx := dc.stack[dc.depth]
			v := actions[0]
			actions = actions[tupleCount:]
			if idx >= len(buffer.Pos) {
				continue
			}

			/* "The end of the list is marked by an odd value..." */
			last = v&1 != 0
			v &= ^1

			o := &buffer.Pos[idx]
			if buffer.Props.Direction.isHorizontal() {
				if dc.crossStream {
					/* The following flag is undocumented in the spec, but described
					 * in the 'kern' table example. */
					if v == -0x8000 {
						o.attachType = attachTypeNone
						o.attachChain = 0
						o.YOffset = 0
					} else if o.attachType != 0 {
						o.YOffset += dc.c.font.emScaleY(v)
						buffer.scratchFlags |= bsfHasGPOSAttachment
					}
				} else if buffer.Info[idx].Mask&kernMask != 0 {
					o.XAdvance += dc.c.font.emScaleX(v)
					o.XOffset += dc.c.font.emScaleX(v)
				}
			} else {
				if dc.crossStream {
					/* CoreText doesn't do crossStream kerning in vertical.  We do. */
					if v == -0x8000 {
						o.attachType = attachTypeNone
						o.attachChain = 0
						o.XOffset = 0
					} else if o.attachType != 0 {
						o.XOffset += dc.c.font.emScaleX(v)
						buffer.scratchFlags |= bsfHasGPOSAttachment
					}
				} else if buffer.Info[idx].Mask&kernMask != 0 {
					o.YAdvance += dc.c.font.emScaleY(v)
					o.YOffset += dc.c.font.emScaleY(v)
				}
			}
		}
	}
}

type driverContextKerx4 struct {
	c          *aatApplyContext
	table      tt.Kerx4
	mark       int
	markSet    bool
	actionType uint8
}

func (driverContextKerx4) inPlace() bool { return true }

func (driverContextKerx4) isActionable(_ stateTableDriver, entry tt.AATStateEntry) bool {
	return entry.AsKernxIndex() != 0xFFFF
}

func (dc *driverContextKerx4) transition(driver stateTableDriver, entry tt.AATStateEntry) {
	buffer := driver.buffer

	ankrActionIndex := entry.AsKernxIndex()
	if dc.markSet && ankrActionIndex != 0xFFFF && buffer.idx < len(buffer.Pos) {
		o := buffer.curPos(0)
		switch dc.actionType {
		case 0: /* Control Point Actions.*/
			/* Indexed into glyph outline. */
			action := dc.table.Anchors[ankrActionIndex].(tt.KerxAnchorControl)

			markX, markY, okMark := dc.c.font.getGlyphContourPointForOrigin(dc.c.buffer.Info[dc.mark].Glyph,
				action.Mark, LeftToRight)
			currX, currY, okCurr := dc.c.font.getGlyphContourPointForOrigin(dc.c.buffer.cur(0).Glyph,
				action.Current, LeftToRight)
			if !okMark || !okCurr {
				return
			}

			o.XOffset = markX - currX
			o.YOffset = markY - currY

		case 1: /* Anchor Point Actions. */
			/* Indexed into 'ankr' table. */
			action := dc.table.Anchors[ankrActionIndex].(tt.KerxAnchorAnchor)

			markAnchor := dc.c.ankrTable.GetAnchor(dc.c.buffer.Info[dc.mark].Glyph, int(action.Mark))
			currAnchor := dc.c.ankrTable.GetAnchor(dc.c.buffer.cur(0).Glyph, int(action.Current))

			o.XOffset = dc.c.font.emScaleX(markAnchor.X) - dc.c.font.emScaleX(currAnchor.X)
			o.YOffset = dc.c.font.emScaleY(markAnchor.Y) - dc.c.font.emScaleY(currAnchor.Y)

		case 2: /* Control Point Coordinate Actions. */
			action := dc.table.Anchors[ankrActionIndex].(tt.KerxAnchorCoordinates)
			o.XOffset = dc.c.font.emScaleX(action.MarkX) - dc.c.font.emScaleX(action.CurrentX)
			o.YOffset = dc.c.font.emScaleY(action.MarkY) - dc.c.font.emScaleY(action.CurrentY)
		}
		o.attachType = attachTypeMark
		o.attachChain = int16(dc.mark - buffer.idx)
		buffer.scratchFlags |= bsfHasGPOSAttachment
	}

	const Mark = 0x8000 /* If set, remember this glyph as the marked glyph. */
	if entry.Flags&Mark != 0 {
		dc.markSet = true
		dc.mark = buffer.idx
	}
}

func (sp *otShapePlan) aatLayoutTrack(font *Font, buffer *Buffer) {
	trak := font.otTables.Trak

	c := newAatApplyContext(sp, font, buffer)
	c.applyTrak(trak)
}

func (c *aatApplyContext) applyTrak(trak tt.TableTrak) {
	trakMask := c.plan.trakMask

	ptem := c.font.Ptem
	if ptem <= 0. {
		return
	}

	buffer := c.buffer
	if buffer.Props.Direction.isHorizontal() {
		trackData := trak.Horizontal
		tracking := int(trackData.GetTracking(ptem, 0))
		advanceToAdd := c.font.emScalefX(float32(tracking))
		offsetToAdd := c.font.emScalefX(float32(tracking / 2))

		iter, count := buffer.graphemesIterator()
		for start, _ := iter.next(); start < count; start, _ = iter.next() {
			if buffer.Info[start].Mask&trakMask == 0 {
				continue
			}
			buffer.Pos[start].XAdvance += advanceToAdd
			buffer.Pos[start].XOffset += offsetToAdd
		}

	} else {
		trackData := trak.Vertical
		tracking := int(trackData.GetTracking(ptem, 0))
		advanceToAdd := c.font.emScalefY(float32(tracking))
		offsetToAdd := c.font.emScalefY(float32(tracking / 2))
		iter, count := buffer.graphemesIterator()
		for start, _ := iter.next(); start < count; start, _ = iter.next() {
			if buffer.Info[start].Mask&trakMask == 0 {
				continue
			}
			buffer.Pos[start].YAdvance += advanceToAdd
			buffer.Pos[start].YOffset += offsetToAdd
		}

	}
}
