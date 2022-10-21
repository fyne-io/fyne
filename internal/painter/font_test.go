package painter_test

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/test"
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
			got, _ := painter.CachedFontFace(tt.style, 14, 1)
			for _, r := range tt.wantGlyphs {
				_, ok := got.GlyphAdvance(r)
				assert.True(t, ok, "symbol font should include: %c", r)
			}
		})
	}
}

func TestDrawString(t *testing.T) {
	for name, tt := range map[string]struct {
		color    color.Color
		style    fyne.TextStyle
		size     float32
		height   int
		string   string
		tabWidth int
		want     string
	}{
		"regular": {
			color:    color.Black,
			style:    fyne.TextStyle{},
			size:     40,
			height:   50,
			string:   "Hello\tworld!",
			tabWidth: 7,
			want:     "hello_TAB_world_regular_size_40_height_50_tab_width_7.png",
		},
		"bold italic": {
			color:    color.NRGBA{R: 255, A: 255},
			style:    fyne.TextStyle{Bold: true, Italic: true},
			size:     27.42,
			height:   42,
			string:   "Hello\tworld!",
			tabWidth: 3,
			want:     "hello_TAB_world_bold_italic_size_27.42_height_42_tab_width_3.png",
		},
		"missing glyphs": {
			color:    color.Black,
			style:    fyne.TextStyle{},
			size:     40,
			height:   50,
			string:   "Missing: ↩",
			tabWidth: 4,
			want:     "missing_glyph.png",
		},
	} {
		t.Run(name, func(t *testing.T) {
			img := image.NewNRGBA(image.Rect(0, 0, 300, 100))
			f, face := painter.CachedFontFace(tt.style, tt.size, 1)
			painter.DrawString(img, tt.string, tt.color, f, face, tt.size, 1, tt.height, tt.tabWidth)
			test.AssertImageMatches(t, "font/"+tt.want, img)
		})
	}
}

func TestMeasureString(t *testing.T) {
	for name, tt := range map[string]struct {
		style    fyne.TextStyle
		size     float32
		string   string
		tabWidth int
		want     float32
	}{
		"regular": {
			style:    fyne.TextStyle{},
			size:     40,
			string:   "Hello\tworld!",
			tabWidth: 7,
			want:     257.82812,
		},
		"bold italic": {
			style:    fyne.TextStyle{Bold: true, Italic: true},
			size:     27.42,
			string:   "Hello\tworld!",
			tabWidth: 3,
			want:     173.17188,
		},
		"missing glyph": {
			style:    fyne.TextStyle{},
			size:     40,
			string:   "Missing: ↩",
			tabWidth: 4,
			want:     189.65625,
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, face := painter.CachedFontFace(tt.style, tt.size, 1)
			got, _ := painter.MeasureString(face, tt.string, tt.size, tt.tabWidth)
			assert.Equal(t, tt.want, got.Width)
		})
	}
}

func TestRenderedTextSize(t *testing.T) {
	size1, baseline1 := painter.RenderedTextSize("Hello World!", 20, fyne.TextStyle{})
	size2, baseline2 := painter.RenderedTextSize("\rH\re\rl\rl\ro\r \rW\ro\rr\rl\rd\r!\r", 20, fyne.TextStyle{})
	assert.Equal(t, int(size1.Width), int(size2.Width))
	assert.Equal(t, size1.Height, size2.Height)
	assert.Equal(t, baseline1, baseline2)
}
