package truetype

import (
	"bytes"
	"encoding/binary"
)

// TableHVhea exposes global metrics for horizontal or vertical writting.
type TableHVhea struct {
	Ascent               int16 // typoAscent in version 1.1
	Descent              int16 // typoDescent in version 1.1
	LineGap              int16 // typoLineGap in version 1.1
	AdvanceMax           uint16
	MinFirstSideBearing  int16 // left/top
	MinSecondSideBearing int16 // right/bottom
	MaxExtent            int16
	CaretSlopeRise       int16
	CaretSlopeRun        int16
	CaretOffset          int16
	numOfLongMetrics     uint16
}

func parseTableHVhea(buf []byte) (*TableHVhea, error) {
	var fields struct {
		Version              fixed
		Ascent               int16
		Descent              int16
		LineGap              int16
		AdvanceMax           uint16
		MinFirstSideBearing  int16
		MinSecondSideBearing int16
		MaxExtent            int16
		CaretSlopeRise       int16
		CaretSlopeRun        int16
		CaretOffset          int16
		Reserved1            int16
		Reserved2            int16
		Reserved3            int16
		Reserved4            int16
		MetricDataformat     int16
		NumOfLongMetrics     uint16
	}
	err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &fields)
	if err != nil {
		return nil, err
	}
	out := TableHVhea{
		Ascent:               fields.Ascent,
		Descent:              fields.Descent,
		LineGap:              fields.LineGap,
		AdvanceMax:           fields.AdvanceMax,
		MinFirstSideBearing:  fields.MinFirstSideBearing,
		MinSecondSideBearing: fields.MinSecondSideBearing,
		MaxExtent:            fields.MaxExtent,
		CaretSlopeRise:       fields.CaretSlopeRise,
		CaretSlopeRun:        fields.CaretSlopeRun,
		CaretOffset:          fields.CaretOffset,
		numOfLongMetrics:     fields.NumOfLongMetrics,
	}
	return &out, nil
}
