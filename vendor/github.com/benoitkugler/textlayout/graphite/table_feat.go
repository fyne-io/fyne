package graphite

import (
	"errors"
	"fmt"
	"sort"

	"github.com/benoitkugler/textlayout/fonts/binaryreader"
	"github.com/benoitkugler/textlayout/fonts/truetype"
)

// FeatureValue specifies a value for a given feature.
type FeatureValue struct {
	ID    Tag   // ID of the feature
	Value int16 // Value to use
}

// FeaturesValue are sorted by Id
type FeaturesValue []FeatureValue

// FindFeature return the feature for the given tag, or nil.
func (feats FeaturesValue) FindFeature(id Tag) *FeatureValue {
	// binary search
	for i, j := 0, len(feats); i < j; {
		h := i + (j-i)/2
		entry := &feats[h]
		if id < entry.ID {
			j = h
		} else if entry.ID < id {
			i = h + 1
		} else {
			return entry
		}
	}
	return nil
}

// features are NOT sorted; they are accessed by (slice) index
// from the opcodes
type tableFeat []feature

type feature struct {
	settings []featureSetting
	id       Tag
	flags    uint16
	label    truetype.NameID
}

type featureSetting struct {
	Value int16
	Label truetype.NameID
}

// return the feature with their first setting selected (or 0)
func (tf tableFeat) defaultFeatures() FeaturesValue {
	out := make(FeaturesValue, len(tf))
	for i, f := range tf {
		out[i].ID = zeroToSpace(f.id)
		if len(f.settings) != 0 {
			out[i].Value = f.settings[0].Value
		}
	}

	// sort by id
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })

	return out
}

func (tf tableFeat) findFeature(id Tag) (feature, bool) {
	for _, feat := range tf {
		if feat.id == id {
			return feat, true
		}
	}
	return feature{}, false
}

func parseTableFeat(data []byte) (tableFeat, error) {
	const headerSize = 12
	if len(data) < headerSize {
		return nil, errors.New("invalid Feat table (EOF)")
	}
	r := binaryreader.NewReader(data)
	version_, _ := r.Uint32()
	version := version_ >> 16
	numFeat, _ := r.Uint16()
	r.Skip(6)

	recordSize := 12
	if version >= 2 {
		recordSize = 16
	}
	featSlice, err := r.FixedSizes(int(numFeat), recordSize)
	if err != nil {
		return nil, fmt.Errorf("invalid Feat table: %s", err)
	}

	rFeat := binaryreader.NewReader(featSlice)
	out := make(tableFeat, numFeat)
	tmpIndexes := make([][2]int, numFeat)
	var maxSettingsLength int
	for i := range out {
		if version >= 2 {
			id_, _ := rFeat.Uint32()
			out[i].id = Tag(id_)
		} else {
			id_, _ := rFeat.Uint16()
			out[i].id = Tag(id_)
		}
		numSettings, _ := rFeat.Uint16()
		if version >= 2 {
			rFeat.Skip(2)
		}
		offset, _ := rFeat.Uint32()
		out[i].flags, _ = rFeat.Uint16()
		label_, _ := rFeat.Uint16()
		out[i].label = truetype.NameID(label_)

		// convert from offset to index
		index := (int(offset) - headerSize - len(featSlice)) / 4
		end := index + int(numSettings)
		if numSettings != 0 && end > maxSettingsLength {
			maxSettingsLength = end
		}

		tmpIndexes[i] = [2]int{index, int(numSettings)}
	}

	// parse the settings array
	allSettings := make([]featureSetting, maxSettingsLength)
	err = r.ReadStruct(allSettings)
	if err != nil {
		return nil, fmt.Errorf("invalid Feat table: %s", err)
	}

	for i, indexes := range tmpIndexes {
		index, length := indexes[0], indexes[1]
		if length == 0 {
			continue
		}
		out[i].settings = allSettings[index : index+length]
	}

	return out, nil
}
