package truetype

type fvarHeader struct {
	majorVersion    uint16 // Major version number of the font variations table — set to 1.
	minorVersion    uint16 // Minor version number of the font variations table — set to 0.
	axesArrayOffset uint16 // Offset in bytes from the beginning of the table to the start of the VariationAxisRecord array.
	reserved        uint16 // This field is permanently reserved. Set to 2.
	axisCount       uint16 // The number of variation axes in the font (the number of records in the axes array).
	axisSize        uint16 // The size in bytes of each VariationAxisRecord — set to 20 (0x0014) for this version.
	instanceCount   uint16 // The number of named instances defined in the font (the number of records in the instances array).
	instanceSize    uint16 // The size in bytes of each InstanceRecord — set to either axisCount * sizeof(Fixed) + 4, or to axisCount * sizeof(Fixed) + 6.
}

type VarAxis struct {
	Tag     Tag       // Tag identifying the design variation for the axis.
	Minimum Float1616 // mininum value on the variation axis that the font covers
	Default Float1616 // default position on the axis
	Maximum Float1616 // maximum value on the variation axis that the font covers
	flags   uint16    // Axis qualifiers — see details below.
	strid   NameID    // name entry in the font's ‘name’ table
}

// Float1616 is a float32, represented in
// fixed 16.16 format in font files.
type Float1616 = float32

func Float1616FromUint(v uint32) Float1616 {
	// value are actually signed integers
	return Float1616(int32(v)) / (1 << 16)
}

func Float1616ToUint(f Float1616) uint32 {
	return uint32(int32(f * (1 << 16)))
}
