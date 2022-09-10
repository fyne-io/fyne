package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type TableFeat []AATFeatureName

// GetFeature performs a binary seach into the names, using `Feature` as key,
// returning `nil` if not found.
func (t TableFeat) GetFeature(feature uint16) *AATFeatureName {
	for i, j := 0, len(t); i < j; {
		h := i + (j-i)/2
		entry := t[h].Feature
		if feature < entry {
			j = h
		} else if entry < feature {
			i = h + 1
		} else {
			return &t[h]
		}
	}
	return nil
}

func parseTableFeat(data []byte) (TableFeat, error) {
	if len(data) < 12 {
		return nil, errors.New("invalid feat table (EOF)")
	}
	featureNameCount := binary.BigEndian.Uint16(data[4:])
	if len(data) < 12+12*int(featureNameCount) {
		return nil, errors.New("invalid feat table (EOF)")
	}
	out := make(TableFeat, featureNameCount)
	var err error
	for i := range out {
		out[i].Feature = binary.BigEndian.Uint16(data[12+12*i:])
		nSettings := binary.BigEndian.Uint16(data[12+12*i+2:])
		offsetSetting := binary.BigEndian.Uint32(data[12+12*i+4:])
		out[i].Flags = binary.BigEndian.Uint16(data[12+12*i+8:])
		out[i].NameIndex = NameID(binary.BigEndian.Uint16(data[12+12*i+10:]))
		out[i].Settings, err = parseAATSettingNames(data, offsetSetting, nSettings)
		if err != nil {
			return nil, err
		}

		// sanitize the index
		if di := out[i].defaultIndex(); di >= nSettings {
			return nil, fmt.Errorf("invalid feat table setting index: %d (for %d)", di, nSettings)
		}
	}

	return out, nil
}

type AATFeatureName struct {
	Settings  []AATSettingName
	Feature   uint16
	Flags     uint16
	NameIndex NameID
}

// IsExclusive returns true if the feature settings are mutually exclusive.
func (feature *AATFeatureName) IsExclusive() bool {
	const Exclusive = 0x8000
	return feature.Flags&Exclusive != 0
}

func (feature *AATFeatureName) defaultIndex() uint16 {
	const (
		aatFeatureNotDefault = 0x4000
		aatFeatureIndexMask  = 0x00FF
	)
	var defaultIndex uint16
	if feature.Flags&aatFeatureNotDefault != 0 {
		defaultIndex = feature.Flags & aatFeatureIndexMask
	}
	return defaultIndex
}

// GetSelectorInfos fetches a list of the selectors available for the feature,
// and the default index. It the later equals 0xFFFF, then
// the feature type is non-exclusive.  Otherwise, it is the index of
// the selector that is selected by default.
func (feature *AATFeatureName) GetSelectorInfos() ([]AATFeatureSelector, uint16) {
	settingsTable := feature.Settings

	defaultSelector := uint16(0xFFFF)
	defaultIndex := uint16(0xFFFF)
	if feature.IsExclusive() {
		defaultIndex = feature.defaultIndex()
		defaultSelector = settingsTable[defaultIndex].Setting
	}

	out := make([]AATFeatureSelector, len(settingsTable))
	for i, setting := range settingsTable {
		out[i] = setting.getSelector(defaultSelector)
	}

	return out, defaultIndex
}

// AATFeatureSelector represents a setting for an AAT feature type.
type AATFeatureSelector struct {
	Name    NameID // selector's name identifier
	Enable  uint16 // value to turn the selector on
	Disable uint16 // value to turn the selector off
}

type AATSettingName struct {
	Setting uint16
	Name    NameID
}

func (s AATSettingName) getSelector(defaultSelector uint16) AATFeatureSelector {
	// AATFeatureSelectorUnset is the initial, unset feature selector
	const AATFeatureSelectorUnset = 0xFFFF
	out := AATFeatureSelector{Name: s.Name, Enable: s.Setting}
	if defaultSelector == AATFeatureSelectorUnset {
		out.Disable = s.Setting + 1
	} else {
		out.Disable = defaultSelector
	}
	return out
}

func parseAATSettingNames(data []byte, offset uint32, count uint16) ([]AATSettingName, error) {
	if len(data) < int(offset)+4*int(count) {
		return nil, errors.New("invalid feat table settings names (EOF)")
	}

	out := make([]AATSettingName, count)
	data = data[offset:]
	for i := range out {
		out[i].Setting = binary.BigEndian.Uint16(data[4*i:])
		out[i].Name = NameID(binary.BigEndian.Uint16(data[4*i+2:]))
	}
	return out, nil
}
