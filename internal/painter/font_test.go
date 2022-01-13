package painter_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/painter"
	"github.com/goki/freetype/truetype"
	"github.com/stretchr/testify/assert"
)

func TestCachedFontFace(t *testing.T) {
	for name, tt := range map[string]struct {
		style      fyne.TextStyle
		wantGlyphs string
	}{
		"symbol font": {
			fyne.TextStyle{
				Symbol: true,
			},
			"←↑→↓↖↘↵↵⇞⇟⇥⇧⌃⌘⌥⌦⌫⎋␣⌃⌥⇧⌘",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := painter.CachedFontFace(tt.style, &truetype.Options{})
			for _, r := range tt.wantGlyphs {
				_, ok := got.GlyphAdvance(r)
				assert.True(t, ok, "symbol font should include: %c", r)
			}
		})
	}
}

func TestRenderedTextSize(t *testing.T) {
	size1, baseline1 := painter.RenderedTextSize("Hello World!", 20, fyne.TextStyle{})
	size2, baseline2 := painter.RenderedTextSize("\rH\re\rl\rl\ro\r \rW\ro\rr\rl\rd\r!\r", 20, fyne.TextStyle{})
	assert.Equal(t, size1.Width, size2.Width)
	assert.Equal(t, size1.Height, size2.Height)
	assert.Equal(t, baseline1, baseline2)
}
