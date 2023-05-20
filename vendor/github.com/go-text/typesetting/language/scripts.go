package language

import (
	"encoding/binary"
	"fmt"
)

// Script identifies different writing systems.
// It is represented as the binary encoding of a script tag of 4 (case sensitive) letters,
// as specified by ISO 15924.
// Note that the default value is usually the Unknown script, not the 0 value (which is invalid)
type Script uint32

// ParseScript simply converts a 4 bytes string into its binary encoding.
func ParseScript(script string) (Script, error) {
	if len(script) != 4 {
		return 0, fmt.Errorf("invalid script string: %s", script)
	}
	return Script(binary.BigEndian.Uint32([]byte(script))), nil
}

// LookupScript looks up the script for a particular character (as defined by
// Unicode Standard Annex #24), and returns Unknown if not found.
func LookupScript(r rune) Script {
	// binary search
	for i, j := 0, len(scriptRanges); i < j; {
		h := i + (j-i)/2
		entry := scriptRanges[h]
		if r < entry.start {
			j = h
		} else if entry.end < r {
			i = h + 1
		} else {
			return entry.script
		}
	}
	return Unknown
}

// String returns the ISO 4 lower letters code of the script
func (s Script) String() string {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(s))
	return string(buf[:])
}
