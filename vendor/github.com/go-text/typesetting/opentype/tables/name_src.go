// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// Naming table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/name
type Name struct {
	version     uint16
	count       uint16
	stringData  []byte       `offsetSize:"Offset16" arrayCount:"ToEnd"`
	nameRecords []nameRecord `arrayCount:"ComputedField-count"`
}

type nameRecord struct {
	platformID   PlatformID
	encodingID   EncodingID
	languageID   LanguageID
	nameID       NameID
	length       uint16
	stringOffset uint16
}
