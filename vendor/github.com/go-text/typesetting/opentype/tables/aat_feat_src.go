// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// Feat is the feature name table.
// See - https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6feat.html
type Feat struct {
	version          uint32        // 	Version number of the feature name table (0x00010000 for the current version).
	featureNameCount uint16        // 	The number of entries in the feature name array.
	none1            uint16        // 	Reserved (set to zero).
	none2            uint32        // 	Reserved (set to zero).
	Names            []FeatureName `arrayCount:"ComputedField-featureNameCount"` //	The feature name array.
}

type FeatureName struct {
	Feature      uint16               // Feature type.
	nSettings    uint16               // The number of records in the setting name array.
	SettingTable []FeatureSettingName `offsetSize:"Offset32" offsetRelativeTo:"Parent" arrayCount:"ComputedField-nSettings"` // Offset in bytes from the beginning of the 'feat' table to this feature's setting name array. The actual type of record this offset refers to will depend on the exclusivity value, as described below.
	FeatureFlags uint16               // Single-bit flags associated with the feature type.
	NameIndex    uint16               // The name table index for the feature's name. This index has values greater than 255 and less than 32768.
}

type FeatureSettingName struct {
	Setting   uint16 //	The setting.
	NameIndex uint16 //	The name table index for the setting's name. The nameIndex must be greater than 255 and less than 32768.
}
