// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// OS/2 and Windows Metrics Table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/os2
type Os2 struct {
	Version             uint16
	XAvgCharWidth       uint16
	USWeightClass       uint16
	USWidthClass        uint16
	fSType              uint16
	YSubscriptXSize     int16
	YSubscriptYSize     int16
	YSubscriptXOffset   int16
	YSubscriptYOffset   int16
	YSuperscriptXSize   int16
	YSuperscriptYSize   int16
	YSuperscriptXOffset int16
	ySuperscriptYOffset int16
	YStrikeoutSize      int16
	YStrikeoutPosition  int16
	sFamilyClass        int16
	panose              [10]byte
	ulCharRange         [4]uint32
	achVendID           Tag
	FsSelection         uint16
	USFirstCharIndex    uint16
	USLastCharIndex     uint16
	STypoAscender       int16
	STypoDescender      int16
	STypoLineGap        int16
	usWinAscent         uint16
	usWinDescent        uint16
	HigherVersionData   []byte `arrayCount:"ToEnd"`
}
