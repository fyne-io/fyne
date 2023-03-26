// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// https://learn.microsoft.com/en-us/typography/opentype/spec/hmtx
type Hmtx struct {
	Metrics []LongHorMetric `arrayCount:""`
	// avances are padded with the last value
	// and side bearings are given
	LeftSideBearings []int16 `arrayCount:""`
}

type LongHorMetric struct {
	AdvanceWidth, LeftSideBearing int16
}

// https://learn.microsoft.com/en-us/typography/opentype/spec/vmtx
type Vmtx = Hmtx
