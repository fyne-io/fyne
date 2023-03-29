package font

import (
	"github.com/go-text/typesetting/harfbuzz"
	"github.com/go-text/typesetting/opentype/api"
	"github.com/go-text/typesetting/opentype/api/font"
)

type (
	Face      = *font.Face
	GID       = api.GID
	GlyphMask = harfbuzz.GlyphMask
)
