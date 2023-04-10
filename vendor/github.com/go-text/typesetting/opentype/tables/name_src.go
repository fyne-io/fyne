// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"encoding/binary"
	"unicode/utf16"
)

const (
	PlatformUnicode PlatformID = iota
	PlatformMac
	PlatformIso // deprecated
	PlatformMicrosoft
	PlatformCustom
	_
	_
	PlatformAdobe // artificial
)

const (
	PEUnicodeDefault     = EncodingID(0)
	PEUnicodeBMP         = EncodingID(3)
	PEUnicodeFull        = EncodingID(4)
	PEUnicodeFull13      = EncodingID(6)
	PEMacRoman           = PEUnicodeDefault
	PEMicrosoftSymbolCs  = EncodingID(0)
	PEMicrosoftUnicodeCs = EncodingID(1)
	PEMicrosoftUcs4      = EncodingID(10)
)

const (
	plMacEnglish       = LanguageID(0)
	plUnicodeDefault   = LanguageID(0)
	plMicrosoftEnglish = LanguageID(0x0409)
)

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

// selectRecord return the entry for `name` or nil if not found.
func (names Name) selectRecord(name NameID) *nameRecord {
	var (
		foundAppleRoman   = -1
		foundAppleEnglish = -1
		foundWin          = -1
		foundUnicode      = -1
		isEnglish         = false
	)

	for n, rec := range names.nameRecords {
		// According to the OpenType 1.3 specification, only Microsoft or
		// Apple platform IDs might be used in the `name' table.  The
		// `Unicode' platform is reserved for the `cmap' table, and the
		// `ISO' one is deprecated.
		//
		// However, the Apple TrueType specification doesn't say the same
		// thing and goes to suggest that all Unicode `name' table entries
		// should be coded in UTF-16.
		if rec.nameID == name && rec.length > 0 {
			switch rec.platformID {
			case PlatformUnicode, PlatformIso:
				// there is `languageID' to check there.  We should use this
				// field only as a last solution when nothing else is
				// available.
				foundUnicode = n
			case PlatformMac:
				// This is a bit special because some fonts will use either
				// an English language id, or a Roman encoding id, to indicate
				// the English version of its font name.
				if rec.languageID == plMacEnglish {
					foundAppleEnglish = n
				} else if rec.encodingID == PEMacRoman {
					foundAppleRoman = n
				}
			case PlatformMicrosoft:
				// we only take a non-English name when there is nothing
				// else available in the font
				if foundWin == -1 || (rec.languageID&0x3FF) == 0x009 {
					switch rec.encodingID {
					case PEMicrosoftSymbolCs, PEMicrosoftUnicodeCs, PEMicrosoftUcs4:
						isEnglish = (rec.languageID & 0x3FF) == 0x009
						foundWin = n
					}
				}
			}
		}
	}

	foundApple := foundAppleRoman
	if foundAppleEnglish >= 0 {
		foundApple = foundAppleEnglish
	}

	// some fonts contain invalid Unicode or Macintosh formatted entries;
	// we will thus favor names encoded in Windows formats if available
	// (provided it is an English name)
	if foundWin >= 0 && !(foundApple >= 0 && !isEnglish) {
		return &names.nameRecords[foundWin]
	} else if foundApple >= 0 {
		return &names.nameRecords[foundApple]
	} else if foundUnicode >= 0 {
		return &names.nameRecords[foundUnicode]
	}
	return nil
}

// Name returns the entry at [name], encoded in UTF-8 when possible,
// or an empty string if not found
func (names Name) Name(name NameID) string {
	if record := names.selectRecord(name); record != nil {
		return names.decodeRecord(*record)
	}
	return ""
}

// decode is a best-effort attempt to get an UTF-8 encoded version of
// Value. Only MicrosoftUnicode (3,1 ,X), MacRomain (1,0,X) and Unicode platform
// strings are supported.
func (names Name) decodeRecord(n nameRecord) string {
	end := int(n.stringOffset) + int(n.length)
	if end > len(names.stringData) {
		// invalid record
		return ""
	}
	value := names.stringData[n.stringOffset:end]

	if n.platformID == PlatformUnicode || (n.platformID == PlatformMicrosoft &&
		n.encodingID == PEMicrosoftUnicodeCs) {
		return decodeUtf16(value)
	}

	if n.platformID == PlatformMac && n.encodingID == PEMacRoman {
		return DecodeMacintosh(value)
	}

	// no encoding detected, hope for utf8
	return string(value)
}

// decode a big ending, no BOM utf16 string
func decodeUtf16(b []byte) string {
	ints := make([]uint16, len(b)/2)
	for i := range ints {
		ints[i] = binary.BigEndian.Uint16(b[2*i:])
	}
	return string(utf16.Decode(ints))
}

// Support for the old macintosh encoding

var macintoshEncoding = [256]rune{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 196, 197, 199, 201, 209, 214, 220, 225, 224, 226, 228, 227, 229, 231, 233, 232, 234, 235, 237, 236, 238, 239, 241, 243, 242, 244, 246, 245, 250, 249, 251, 252, 8224, 176, 162, 163, 167, 8226, 182, 223, 174, 169, 8482, 180, 168, 8800, 198, 216, 8734, 177, 8804, 8805, 165, 181, 8706, 8721, 8719, 960, 8747, 170, 186, 937, 230, 248, 191, 161, 172, 8730, 402, 8776, 8710, 171, 187, 8230, 160, 192, 195, 213, 338, 339, 8211, 8212, 8220, 8221, 8216, 8217, 247, 9674, 255, 376, 8260, 8364,
	8249, 8250, 64257, 64258, 8225, 183, 8218, 8222, 8240, 194, 202, 193, 203, 200, 205, 206, 207, 204, 211, 212, 63743, 210, 218, 219, 217, 305, 710, 732, 175, 728, 729, 730, 184, 733, 731, 711,
}

// DecodeMacintoshByte returns the rune for the given byte
func DecodeMacintoshByte(b byte) rune { return macintoshEncoding[b] }

// DecodeMacintosh decode a Macintosh encoded string
func DecodeMacintosh(encoded []byte) string {
	out := make([]rune, len(encoded))
	for i, b := range encoded {
		out[i] = macintoshEncoding[b]
	}
	return string(out)
}
