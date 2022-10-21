package truetype

var _ FaceVariable = (*Font)(nil)

// FaceVariable is an extension interface supporting Opentype variable fonts.
// See the `Variations` method to check if a font is actually variable.
type FaceVariable interface {
	// Variations returns the variations for the font,
	// or an empty table for non-variable fonts.
	Variations() TableFvar

	// SetVarCoordinates apply the normalized coordinates values.
	// Use `NormalizeVariations` to convert from design space units.
	// See also `SetVariations`.
	SetVarCoordinates(coords []float32)

	// VarCoordinates returns the current variable coordinates,
	// in normalized units.
	VarCoordinates() []float32

	// NormalizeVariations should normalize the given design-space coordinates. The minimum and maximum
	// values for the axis are mapped to the interval [-1,1], with the default
	// axis value mapped to 0.
	// This should be a no-op for non-variable fonts.
	NormalizeVariations(coords []float32) []float32
}

// SetVariations applies a list of font-variation settings to a font,
// defaulting to the values given in the `fvar` table.
// Note that passing an empty slice will instead remove the coordinates.
func SetVariations(face FaceVariable, variations []Variation) {
	if len(variations) == 0 {
		face.SetVarCoordinates(nil)
		return
	}

	fvar := face.Variations()
	if len(fvar.Axis) == 0 {
		face.SetVarCoordinates(nil)
		return
	}

	designCoords := fvar.GetDesignCoordsDefault(variations)

	face.SetVarCoordinates(face.NormalizeVariations(designCoords))
}

func (font *Font) SetVarCoordinates(coords []float32) {
	font.varCoords = coords
}

func (font *Font) VarCoordinates() []float32 { return font.varCoords }

// Variation defines a value for a wanted variation axis.
type Variation struct {
	Tag   Tag     // variation-axis identifier tag
	Value float32 // in design units
}

type VarInstance struct {
	Coords    []float32 // in design units; length: number of axis
	Subfamily NameID

	PSStringID NameID
}

type TableFvar struct {
	Axis      []VarAxis
	Instances []VarInstance // contains the default instance
}

// IsDefaultInstance returns `true` is `instance` has the same
// coordinates as the default instance.
func (fvar TableFvar) IsDefaultInstance(it VarInstance) bool {
	for i, c := range it.Coords {
		if c != fvar.Axis[i].Default {
			return false
		}
	}
	return true
}

// add the default instance if it not already explicitely present
func (fvar *TableFvar) checkDefaultInstance(names TableName) {
	for _, instance := range fvar.Instances {
		if fvar.IsDefaultInstance(instance) {
			return
		}
	}

	// add the default instance
	// choose the subfamily entry
	subFamily := NamePreferredSubfamily
	if v1, v2 := names.getEntry(subFamily); v1 == nil && v2 == nil {
		subFamily = NameFontSubfamily
	}
	defaultInstance := VarInstance{
		Coords:     make([]float32, len(fvar.Axis)),
		Subfamily:  subFamily,
		PSStringID: NamePostscript,
	}
	for i, axe := range fvar.Axis {
		defaultInstance.Coords[i] = axe.Default
	}
	fvar.Instances = append(fvar.Instances, defaultInstance)
}

// GetDesignCoordsDefault returns the design coordinates corresponding to the given pairs of axis/value.
// The default value of the axis is used when not specified in the variations.
func (fvar *TableFvar) GetDesignCoordsDefault(variations []Variation) []float32 {
	designCoords := make([]float32, len(fvar.Axis))
	// start with default values
	for i, axis := range fvar.Axis {
		designCoords[i] = axis.Default
	}

	fvar.GetDesignCoords(variations, designCoords)

	return designCoords
}

// GetDesignCoords updates the design coordinates, with the given pairs of axis/value.
// It will panic if `designCoords` has not the length expected by the table, that is the number of axis.
func (fvar *TableFvar) GetDesignCoords(variations []Variation, designCoords []float32) {
	for _, variation := range variations {
		// allow for multiple axis with the same tag
		for index, axis := range fvar.Axis {
			if axis.Tag == variation.Tag {
				designCoords[index] = variation.Value
			}
		}
	}
}

// normalize based on the [min,def,max] values for the axis to be [-1,0,1].
func (fvar *TableFvar) normalizeCoordinates(coords []float32) []float32 {
	normalized := make([]float32, len(coords))
	for i, a := range fvar.Axis {
		coord := coords[i]

		// out of range: clamping
		if coord > a.Maximum {
			coord = a.Maximum
		} else if coord < a.Minimum {
			coord = a.Minimum
		}

		if coord < a.Default {
			normalized[i] = -(coord - a.Default) / (a.Minimum - a.Default)
		} else if coord > a.Default {
			normalized[i] = (coord - a.Default) / (a.Maximum - a.Default)
		} else {
			normalized[i] = 0
		}
	}
	return normalized
}

func (f *Font) Variations() TableFvar { return f.fvar }

// Normalizes the given design-space coordinates. The minimum and maximum
// values for the axis are mapped to the interval [-1,1], with the default
// axis value mapped to 0.
// Any additional scaling defined in the face's `avar` table is also
// applied, as described at https://docs.microsoft.com/en-us/typography/opentype/spec/avar
func (f *Font) NormalizeVariations(coords []float32) []float32 {
	// ported from freetype2

	// Axis normalization is a two-stage process.  First we normalize
	// based on the [min,def,max] values for the axis to be [-1,0,1].
	// Then, if there's an `avar' table, we renormalize this range.
	normalized := f.fvar.normalizeCoordinates(coords)

	// now applying 'avar'
	for i, av := range f.avar {
		for j := 1; j < len(av); j++ {
			previous, pair := av[j-1], av[j]
			if normalized[i] < pair.from {
				normalized[i] =
					previous.to + (normalized[i]-previous.from)*
						(pair.to-previous.to)/(pair.from-previous.from)
				break
			}
		}
	}

	return normalized
}
