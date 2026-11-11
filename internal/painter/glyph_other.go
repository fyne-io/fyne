//go:build !darwin

package painter

import (
	"image"

	"github.com/go-text/typesetting/font"
)

// glyphMargin is the room the atlas leaves around a glyph's ink, for its
// antialiased edge.
// nativeGlyphs is true where rasteriseGlyphNative draws, so the whole-run
// path draws a glyph at a time to match the atlas.
const nativeGlyphs = false

const glyphMargin = 1

func rememberFontSource(*font.Font, []byte) {}

func forgetNativeFonts() {}

// rasteriseGlyphNative is darwin's; elsewhere go-text draws every glyph.
func rasteriseGlyphNative(*image.RGBA, PlacedGlyph, float32, float32, int, int, float32) bool {
	return false
}
