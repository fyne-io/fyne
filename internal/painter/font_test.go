package painter_test

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/test"
	"github.com/goki/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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

func TestDrawString(t *testing.T) {
	for name, tt := range map[string]struct {
		color    color.Color
		face     font.Face
		height   int
		string   string
		tabWidth int
		want     string
	}{
		"regular": {
			color:    color.Black,
			face:     regular,
			height:   50,
			string:   "Hello\tworld!",
			tabWidth: 7,
			want:     "hello_TAB_world_regular_size_40_height_50_tab_width_7.png",
		},
		"bold italic": {
			color:    color.NRGBA{R: 255, A: 255},
			face:     boldItalic,
			height:   42,
			string:   "Hello\tworld!",
			tabWidth: 3,
			want:     "hello_TAB_world_bold_italic_size_27.42_height_42_tab_width_3.png",
		},
		"missing glyphs": {
			color:    color.Black,
			face:     regular,
			height:   50,
			string:   "Missing: ↩",
			tabWidth: 4,
			want:     "missing_glyph.png",
		},
	} {
		t.Run(name, func(t *testing.T) {
			img := image.NewNRGBA(image.Rect(0, 0, 300, 100))
			painter.DrawString(img, tt.string, tt.color, tt.face, tt.height, tt.tabWidth)
			test.AssertImageMatches(t, "font/"+tt.want, img)
		})
	}
}

func TestMeasureString(t *testing.T) {
	for name, tt := range map[string]struct {
		face     font.Face
		string   string
		tabWidth int
		want     fixed.Int26_6
	}{
		"regular": {
			face:     regular,
			string:   "Hello\tworld!",
			tabWidth: 7,
			want:     18263, // 285.359375
		},
		"bold italic": {
			face:     boldItalic,
			string:   "Hello\tworld!",
			tabWidth: 3,
			want:     11576, // 180.875
		},
		"missing glyph": {
			face:     regular,
			string:   "Missing: ↩",
			tabWidth: 4,
			want:     14257, // 222.765625
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := painter.MeasureString(tt.face, tt.string, tt.tabWidth)
			assert.Equal(t, tt.want, got)
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

var regular = painter.CachedFontFace(
	fyne.TextStyle{},
	&truetype.Options{
		Size: 40.0,
		DPI:  painter.TextDPI,
	},
)

var boldItalic = painter.CachedFontFace(
	fyne.TextStyle{
		Bold:   true,
		Italic: true,
	},
	&truetype.Options{
		Size: 27.42,
		DPI:  painter.TextDPI,
	},
)
