package test

import "github.com/go-text/typesetting/font"

type FontMap []*font.Face

func (f FontMap) ResolveFace(r rune) *font.Face {
	if len(f) == 1 {
		return f[0]
	}

	face := f.ResolveFace(r)
	if face != nil {
		return face
	}

	return f[0]
}
