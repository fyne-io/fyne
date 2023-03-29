// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

// Package opentype provides the low level routines
// required to read Opentype font files, including collections.
//
// This package is designed to provide an efficient, lazy, reading API.
//
// For the parsing of the various tables, see package [tables].
package loader

// Tag represents an open-type name.
// These are technically uint32's, but are usually
// displayed in ASCII as they are all acronyms.
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6.html#Overview
type Tag uint32

// NewTag returns the tag for <abcd>.
func NewTag(a, b, c, d byte) Tag {
	return Tag(uint32(d) | uint32(c)<<8 | uint32(b)<<16 | uint32(a)<<24)
}

// MustNewTag gives you the Tag corresponding to the acronym.
// This function will panic if the string passed in is not 4 bytes long.
func MustNewTag(str string) Tag {
	if len(str) != 4 {
		panic("invalid tag: must be exactly 4 bytes")
	}
	_ = str[3]
	return NewTag(str[0], str[1], str[2], str[3])
}

// String return the ASCII form of the tag.
func (t Tag) String() string {
	return string([]byte{
		byte(t >> 24),
		byte(t >> 16),
		byte(t >> 8),
		byte(t),
	})
}
