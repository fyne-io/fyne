package color_test

import (
	imagecolor "image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/internal/color"
)

func Test_Parse(t *testing.T) {
	for name, tt := range map[string]struct {
		input   string
		want    imagecolor.NRGBA
		wantErr string
	}{
		"ok: #rrggbbaa": {
			input: "#1a2b3c4d",
			want:  imagecolor.NRGBA{R: 0x1a, G: 0x2b, B: 0x3c, A: 0x4d},
		},
		"ok: rrggbbaa": {
			input: "1a2b3c4d",
			want:  imagecolor.NRGBA{R: 0x1a, G: 0x2b, B: 0x3c, A: 0x4d},
		},
		"ok: #rrggbb": {
			input: "#1a2b3c",
			want:  imagecolor.NRGBA{R: 0x1a, G: 0x2b, B: 0x3c, A: 0xff},
		},
		"ok: rrggbb": {
			input: "1a2b3c",
			want:  imagecolor.NRGBA{R: 0x1a, G: 0x2b, B: 0x3c, A: 0xff},
		},
		"ok: #rgba": {
			input: "#1a2b",
			want:  imagecolor.NRGBA{R: 0x11, G: 0xaa, B: 0x22, A: 0xbb},
		},
		"ok: rgba": {
			input: "1a2b",
			want:  imagecolor.NRGBA{R: 0x11, G: 0xaa, B: 0x22, A: 0xbb},
		},
		"ok: #rgb": {
			input: "#1a2",
			want:  imagecolor.NRGBA{R: 0x11, G: 0xaa, B: 0x22, A: 0xff},
		},
		"ok: rgb": {
			input: "1a2",
			want:  imagecolor.NRGBA{R: 0x11, G: 0xaa, B: 0x22, A: 0xff},
		},
		"err: #rg": {
			input:   "#12",
			wantErr: "invalid color format: #12",
		},
		"err: #rrggb": {
			input:   "#12345",
			wantErr: "invalid color format: #12345",
		},
		"err: #rrggbba": {
			input:   "#1234567",
			wantErr: "invalid color format: #1234567",
		},
		"err: #rrggbbaax": {
			input:   "#123456789",
			wantErr: "invalid color format: #123456789",
		},
		"err: xrrggbbaa": {
			input:   "x1a2b3c4d",
			wantErr: "invalid color format: x1a2b3c4d",
		},
		"err: #rrggbbax": {
			input:   "#1a2b3c4x",
			wantErr: "encoding/hex: invalid byte: U+0078 'x'",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := color.Parse(tt.input)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_ToNRGBA_unmultiplyAlpha(t *testing.T) {
	for name, tt := range map[string]struct {
		color imagecolor.Color
		wantR uint8
		wantG uint8
		wantB uint8
		wantA uint8
	}{
		"RGBA": {
			color: imagecolor.RGBA{R: 100, G: 100, B: 100, A: 100},
			wantR: 255,
			wantG: 255,
			wantB: 255,
			wantA: 100,
		},
		"RGBA opaque": {
			color: imagecolor.RGBA{R: 100, G: 100, B: 100, A: 255},
			wantR: 100,
			wantG: 100,
			wantB: 100,
			wantA: 255,
		},
		"RGBA64": {
			color: imagecolor.RGBA64{R: 100<<8 + 123, G: 100<<8 + 123, B: 100<<8 + 123, A: 100<<8 + 123},
			wantR: 255,
			wantG: 255,
			wantB: 255,
			wantA: 100,
		},
		"RGBA64 opaque": {
			color: imagecolor.RGBA64{R: 100<<8 + 123, G: 100<<8 + 123, B: 100<<8 + 123, A: 255 << 8},
			wantR: 100,
			wantG: 100,
			wantB: 100,
			wantA: 255,
		},
		"custom": {
			color: customColor{r: 100<<8 + 123, g: 100<<8 + 123, b: 100<<8 + 123, a: 100<<8 + 123},
			wantR: 255,
			wantG: 255,
			wantB: 255,
			wantA: 100,
		},
		"custom opaque": {
			color: customColor{r: 100<<8 + 123, g: 100<<8 + 123, b: 100<<8 + 123, a: 255 << 8},
			wantR: 100,
			wantG: 100,
			wantB: 100,
			wantA: 255,
		},
	} {
		t.Run(name, func(t *testing.T) {
			gotR, gotG, gotB, gotA := color.ToNRGBA(tt.color)
			assert.Equal(t, tt.wantR, gotR)
			assert.Equal(t, tt.wantG, gotG)
			assert.Equal(t, tt.wantB, gotB)
			assert.Equal(t, tt.wantA, gotA)
		})
	}
}

type customColor struct {
	r, g, b, a uint32
}

var _ imagecolor.Color = customColor{}

func (c customColor) RGBA() (r, g, b, a uint32) {
	return c.r, c.g, c.b, c.a
}
