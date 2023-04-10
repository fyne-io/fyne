// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// TableHead contains critical information about the rest of the font.
// https://learn.microsoft.com/en-us/typography/opentype/spec/head
// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6head.html
type Head struct {
	majorVersion       uint16
	minorVersion       uint16
	fontRevision       uint32
	checksumAdjustment uint32
	magicNumber        uint32
	flags              uint16
	UnitsPerEm         uint16
	created            longdatetime
	modified           longdatetime
	XMin               int16
	YMin               int16
	XMax               int16
	YMax               int16
	MacStyle           uint16
	lowestRecPPEM      uint16
	fontDirectionHint  int16
	IndexToLocFormat   int16
	glyphDataFormat    int16
}

// Upem returns a sanitize version of the 'UnitsPerEm' field.
func (head *Head) Upem() uint16 {
	if head.UnitsPerEm < 16 || head.UnitsPerEm > 16384 {
		return 1000
	}
	return head.UnitsPerEm
}
