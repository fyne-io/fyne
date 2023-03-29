// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package font

import (
	"encoding/binary"
	"errors"

	"github.com/go-text/typesetting/opentype/tables"
)

type os2 struct {
	version       uint16
	xAvgCharWidth uint16

	useTypoMetrics bool // true if the field sTypoAscender, sTypoDescender and sTypoLineGap are valid.

	ySubscriptXSize     float32
	ySubscriptYSize     float32
	ySubscriptXOffset   float32
	ySubscriptYOffset   float32
	ySuperscriptXSize   float32
	ySuperscriptYSize   float32
	ySuperscriptXOffset float32
	yStrikeoutSize      float32
	yStrikeoutPosition  float32
	sTypoAscender       float32
	sTypoDescender      float32
	sTypoLineGap        float32
	sxHeigh             float32
	sCapHeight          float32
}

func newOs2(os tables.Os2) (os2, error) {
	out := os2{
		version:             os.Version,
		xAvgCharWidth:       os.XAvgCharWidth,
		ySubscriptXSize:     float32(os.YSubscriptXSize),
		ySubscriptYSize:     float32(os.YSubscriptYSize),
		ySubscriptXOffset:   float32(os.YSubscriptXOffset),
		ySubscriptYOffset:   float32(os.YSubscriptYOffset),
		ySuperscriptXSize:   float32(os.YSuperscriptXSize),
		ySuperscriptYSize:   float32(os.YSuperscriptYSize),
		ySuperscriptXOffset: float32(os.YSuperscriptXOffset),
		yStrikeoutSize:      float32(os.YStrikeoutSize),
		yStrikeoutPosition:  float32(os.YStrikeoutPosition),
		sTypoAscender:       float32(os.STypoAscender),
		sTypoDescender:      float32(os.STypoDescender),
		sTypoLineGap:        float32(os.STypoLineGap),
	}
	// add addition info for version >= 2
	if os.Version >= 2 {
		if len(os.HigherVersionData) < 12 {
			return os2{}, errors.New("invalid table os2")
		}
		out.sxHeigh = float32(binary.BigEndian.Uint16(os.HigherVersionData[8:]))
		out.sCapHeight = float32(binary.BigEndian.Uint16(os.HigherVersionData[10:]))
	}

	const useTypoMetrics = 1 << 7
	use := os.FsSelection&useTypoMetrics != 0
	hasData := os.USWeightClass != 0 || os.USWidthClass != 0 || os.USFirstCharIndex != 0 || os.USLastCharIndex != 0
	out.useTypoMetrics = use && hasData

	return out, nil
}
