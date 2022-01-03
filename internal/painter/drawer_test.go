package painter_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/test"
)

func TestDrawString(t *testing.T) {
	regular := painter.CachedFontFace(
		fyne.TextStyle{},
		&truetype.Options{
			Size: 40.0,
			DPI:  painter.TextDPI,
		},
	)
	boldItalic := painter.CachedFontFace(
		fyne.TextStyle{
			Bold:   true,
			Italic: true,
		},
		&truetype.Options{
			Size: 27.42,
			DPI:  painter.TextDPI,
		},
	)

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
	} {
		t.Run(name, func(t *testing.T) {
			img := image.NewNRGBA(image.Rect(0, 0, 300, 100))
			painter.DrawString(img, tt.string, tt.color, tt.face, tt.height, tt.tabWidth)
			test.AssertImageMatches(t, "font/"+tt.want, img)
		})
	}
}
