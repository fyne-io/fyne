// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// https://learn.microsoft.com/en-us/typography/opentype/spec/hhea
type Hhea struct {
	majorVersion         uint16
	minorVersion         uint16
	Ascender             int16
	Descender            int16
	LineGap              int16
	AdvanceMax           uint16
	MinFirstSideBearing  int16
	MinSecondSideBearing int16
	MaxExtent            int16
	CaretSlopeRise       int16
	CaretSlopeRun        int16
	CaretOffset          int16
	reserved             [4]uint16
	metricDataformat     int16
	NumOfLongMetrics     uint16
}

// https://learn.microsoft.com/en-us/typography/opentype/spec/vhea
type Vhea = Hhea
