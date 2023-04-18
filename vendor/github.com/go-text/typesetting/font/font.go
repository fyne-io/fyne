package font

import (
	"fmt"

	"github.com/go-text/typesetting/harfbuzz"
	"github.com/go-text/typesetting/opentype/api"
	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/loader"
)

type (
	// Font is a readonly view of a font file, safe for concurrent use.
	Font = *font.Font
	// Face is a [Font] with user settings, not safe for concurrent use.
	Face = *font.Face

	GID       = api.GID
	GlyphMask = harfbuzz.GlyphMask
)

type Resource = loader.Resource

// ParseTTF parse an Opentype font file (.otf, .ttf).
// See ParseTTC for support for collections.
func ParseTTF(file Resource) (Face, error) {
	ld, err := loader.NewLoader(file)
	if err != nil {
		return nil, err
	}
	ft, err := font.NewFont(ld)
	if err != nil {
		return nil, err
	}
	return &font.Face{Font: ft}, nil
}

// ParseTTC parse an Opentype font file, with support for collections.
// Single font files are supported, returning a slice with length 1.
func ParseTTC(file Resource) ([]Face, error) {
	lds, err := loader.NewLoaders(file)
	if err != nil {
		return nil, err
	}
	out := make([]Face, len(lds))
	for i, ld := range lds {
		ft, err := font.NewFont(ld)
		if err != nil {
			return nil, fmt.Errorf("reading font %d of collection: %s", i, err)
		}
		out[i] = &font.Face{Font: ft}
	}

	return out, nil
}
