// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// https://learn.microsoft.com/en-us/typography/opentype/spec/Maxp
type Maxp struct {
	version   maxpVersion
	NumGlyphs uint16
	data      maxpData `unionField:"version"`
}

type maxpVersion uint32

const (
	maxpVersion05 maxpVersion = 0x00005000
	maxpVersion1  maxpVersion = 0x00010000
)

type maxpData interface {
	isMaxpVersion()
}

func (maxpData05) isMaxpVersion() {}
func (maxpData1) isMaxpVersion()  {}

type maxpData05 struct{}

type maxpData1 struct {
	rawData [13]uint16
}
