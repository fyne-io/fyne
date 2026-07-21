package painter_test

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
	intTest "fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/test"
)

func TestCachedFontFace(t *testing.T) {
	for name, tt := range map[string]struct {
		style fyne.TextStyle
		runes string
	}{
		"symbol font": {
			fyne.TextStyle{
				Symbol: true,
			},
			"←↑→↓↖↘↵↵⇞⇟⇥⇧⌃⌘⌥⌦⌫⎋␣⌃⌥⇧⌘",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := painter.CachedFontFace(tt.style, nil, nil)
			for _, r := range tt.runes {
				f := got.Fonts.ResolveFace(r)
				assert.NotNil(t, f, "symbol Font should include: %c", r)
			}
		})
	}

	// check the wide symbol rune
	symbol := canvas.NewText("⌘", color.Black)
	symbol.TextStyle.Symbol = true
	assert.True(t, symbol.MinSize().Width > 10)
}

func TestDrawString(t *testing.T) {
	for name, tt := range map[string]struct {
		color    color.Color
		style    fyne.TextStyle
		size     float32
		string   string
		tabWidth int
		want     string
	}{
		"regular": {
			color:    color.Black,
			style:    fyne.TextStyle{},
			size:     40,
			string:   "Hello\tworld!",
			tabWidth: 7,
			want:     "hello_TAB_world_regular_size_40_height_50_tab_width_7.png",
		},
		"bold italic": {
			color:    color.NRGBA{R: 255, A: 255},
			style:    fyne.TextStyle{Bold: true, Italic: true},
			size:     27.42,
			string:   "Hello\tworld!",
			tabWidth: 3,
			want:     "hello_TAB_world_bold_italic_size_27.42_height_42_tab_width_3.png",
		},
		"missing glyphs": {
			color:    color.Black,
			style:    fyne.TextStyle{},
			size:     40,
			string:   "Missing: स",
			tabWidth: 4,
			want:     "missing_glyph.png",
		},
	} {
		t.Run(name, func(t *testing.T) {
			img := image.NewNRGBA(image.Rect(0, 0, 300, 100))
			f := painter.CachedFontFace(tt.style, nil, nil)

			fontMap := &intTest.FontMap{f.Fonts.ResolveFace(' ')} // first (ascii) font
			painter.DrawString(img, tt.string, tt.color, fontMap, tt.size, 1, fyne.TextStyle{TabWidth: tt.tabWidth})
			test.AssertImageMatches(t, "font/"+tt.want, img)
		})
	}
}

func TestDrawStringOffset(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
	}{
		{name: "long ascii", text: "0123456789" + strings.Repeat(" abcdefghijklmnopqrstuvwxyz", 60)},
		{name: "ligatures", text: strings.Repeat("office affine official ", 40)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			style := fyne.TextStyle{}
			size := float32(24)
			color := color.Black
			faces := painter.CachedFontFace(style, nil, nil)
			fontMap := &intTest.FontMap{faces.Fonts.ResolveFace(' ')}
			textSize, _ := painter.MeasureString(fontMap, tc.text, size, style)
			width := int(textSize.Width) + 32
			height := int(textSize.Height) + 32
			offset := width / 3
			cropWidth := width / 4

			full := image.NewNRGBA(image.Rect(0, 0, width, height))
			painter.DrawString(full, tc.text, color, fontMap, size, 1, style)

			cropped := image.NewNRGBA(image.Rect(0, 0, cropWidth, height))
			painter.DrawStringOffset(cropped, tc.text, color, fontMap, size, 1, style, offset)

			maxDiff := 0
			for y := 0; y < height; y++ {
				fullStart := full.PixOffset(offset, y)
				croppedStart := cropped.PixOffset(0, y)
				for x := 0; x < cropWidth*4; x++ {
					diff := int(full.Pix[fullStart+x]) - int(cropped.Pix[croppedStart+x])
					if diff < 0 {
						diff = -diff
					}
					if diff > maxDiff {
						maxDiff = diff
					}
				}
			}
			assert.LessOrEqual(t, maxDiff, 8)
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
			string:   "Missing: स",
			tabWidth: 4,
			want:     213.65625,
		},
	} {
		t.Run(name, func(t *testing.T) {
			faces := painter.CachedFontFace(tt.style, nil, nil)
			fontMap := &intTest.FontMap{faces.Fonts.ResolveFace(' ')} // first (ascii) font
			got, _ := painter.MeasureString(fontMap, tt.string, tt.size, fyne.TextStyle{TabWidth: tt.tabWidth})
			assert.Equal(t, tt.want, got.Width)
		})
	}
}

func TestRenderedTextSize(t *testing.T) {
	size1, baseline1 := painter.RenderedTextSize("Hello World!", 20, fyne.TextStyle{}, nil)
	size2, baseline2 := painter.RenderedTextSize("\rH\re\rl\rl\ro\r \rW\ro\rr\rl\rd\r!\r", 20, fyne.TextStyle{}, nil)
	assert.Equal(t, int(size1.Width), int(size2.Width))
	assert.Equal(t, size1.Height, size2.Height)
	assert.Equal(t, baseline1, baseline2)
}
