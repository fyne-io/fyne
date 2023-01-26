package language

import (
	"encoding/binary"
	"fmt"
)

// Script identifies different writing systems.
// It is represented as the binary encoding of a script tag of 4 letters,
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

func (s Script) String() string {
	for k, v := range scriptToTag {
		if v == s {
			return k
		}
	}
	return fmt.Sprintf("<script unknown: %d>", s)
}

// IsRealScript return `true` if `s` if valid,
// and neither common or inherited.
func (s Script) IsRealScript() bool {
	switch s {
	case 0, Unknown, Common, Inherited:
		return false
	default:
		return true
	}
}

// IsSameScript compares two scripts: if one them
// is not 'real' (see IsRealScript), they are compared equal.
func (s1 Script) IsSameScript(s2 Script) bool {
	return s1 == s2 || !s1.IsRealScript() || !s2.IsRealScript()
}
