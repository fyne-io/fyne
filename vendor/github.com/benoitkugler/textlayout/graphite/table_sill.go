package graphite

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sort"

	"github.com/benoitkugler/textlayout/fonts/binaryreader"
)

type tableSill []languageRecord

// replace the trailing space by zero
func spaceToZero(x Tag) Tag {
	switch {
	case x == 0x20202020:
		return 0
	case (x & 0x00FFFFFF) == 0x00202020:
		return x & 0xFF000000
	case (x & 0x0000FFFF) == 0x00002020:
		return x & 0xFFFF0000
	case (x & 0x000000FF) == 0x00000020:
		return x & 0xFFFFFF00
	default:
		return x
	}
}

func zeroToSpace(x Tag) Tag {
	switch {
	case x == 0:
		return 0x20202020
	case (x & 0x00FFFFFF) == 0:
		return x & 0xFF202020
	case (x & 0x0000FFFF) == 0:
		return x & 0xFFFF2020
	case (x & 0x000000FF) == 0:
		return x & 0xFFFFFF20
	default:
		return x
	}
}

// getFeatures selects the features and values for the given language, or
// the default ones if the language is not found.
func (si tableSill) getFeatures(langname Tag, features tableFeat) FeaturesValue {
	langname = spaceToZero(langname)

	for _, rec := range si {
		if rec.langcode == langname {
			return rec.applyValues(features)
		}
	}

	return features.defaultFeatures()
}

type languageRecord struct {
	settings []languageSetting
	langcode Tag
}

// resolve the feature
func (lr languageRecord) applyValues(features tableFeat) FeaturesValue {
	var out FeaturesValue
	for _, set := range lr.settings {
		if _, ok := features.findFeature(set.FeatureId); ok {
			out = append(out, FeatureValue{
				ID:    zeroToSpace(set.FeatureId), // from the internal convention to the external
				Value: set.Value,
			})
		}
	}

	// sort by Id
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })

	return out
}

type languageSetting struct {
	FeatureId Tag
	Value     int16
	_         [2]byte
}

func parseTableSill(data []byte) (tableSill, error) {
	r := binaryreader.NewReader(data)
	if len(data) < 12 {
		return nil, errors.New("invalid Sill table (EOF)")
	}
	_, _ = r.Uint32()
	numLangs, _ := r.Uint16()
	r.Skip(6)

	type languageEntry struct {
		Langcode    [4]byte
		NumSettings uint16
		Offset      uint16
	}

	entries := make([]languageEntry, numLangs)
	err := r.ReadStruct(entries)
	if err != nil {
		return nil, fmt.Errorf("invalid Sill table: %s", err)
	}

	out := make(tableSill, numLangs)
	for i, entry := range entries {
		out[i].langcode = Tag(binary.BigEndian.Uint32(entry.Langcode[:]))
		out[i].settings = make([]languageSetting, entry.NumSettings)
		r.SetPos(int(entry.Offset))
		err := r.ReadStruct(out[i].settings)
		if err != nil {
			return nil, fmt.Errorf("invalid Sill table: %s", err)
		}
	}

	return out, nil
}
