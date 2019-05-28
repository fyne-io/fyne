package canvas

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This is a terribly designed test but I don't know
// a better way to determine the idempotency
// of the image generation

func TestNewRectangleGradient(t *testing.T) {
	var test = []struct {
		Name       string
		StartColor color.Color
		EndColor   color.Color
		Alignment  GradientDirection
		Result     [][]color.RGBA64
	}{
		{"Black to Transparent Horizontal", color.Black, color.Transparent, HORIZONTAL, [][]color.RGBA64{
			{
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 65535},
			},
			{
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 52428},
			},
			{
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 39321},
			},
			{
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 26214},
			},
			{
				{R: 0, G: 0, B: 0, A: 13107},
				{R: 0, G: 0, B: 0, A: 13107},
				{R: 0, G: 0, B: 0, A: 13107},
				{R: 0, G: 0, B: 0, A: 13107},
				{R: 0, G: 0, B: 0, A: 13107},
			}},
		},
		{"Black to Transparent Vertical", color.Black, color.Transparent, VERTICAL, [][]color.RGBA64{
			{
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 13107},
			},
			{
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 13107},
			},
			{
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 13107},
			},
			{
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 13107},
			},
			{
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 52428},
				{R: 0, G: 0, B: 0, A: 39321},
				{R: 0, G: 0, B: 0, A: 26214},
				{R: 0, G: 0, B: 0, A: 13107},
			}},
		},
		{"Black to Transparent Circular", color.Black, color.Transparent, CIRCULAR, [][]color.RGBA64{
			{
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65021},
			},
			{
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65278},
				{R: 0, G: 0, B: 0, A: 65278},
				{R: 0, G: 0, B: 0, A: 65278},
				{R: 0, G: 0, B: 0, A: 65021},
			},
			{
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65278},
				{R: 0, G: 0, B: 0, A: 65535},
				{R: 0, G: 0, B: 0, A: 65278},
				{R: 0, G: 0, B: 0, A: 65021},
			},
			{
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65278},
				{R: 0, G: 0, B: 0, A: 65278},
				{R: 0, G: 0, B: 0, A: 65278},
				{R: 0, G: 0, B: 0, A: 65021},
			},
			{
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65021},
				{R: 0, G: 0, B: 0, A: 65021},
			}},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			grad := NewRectangleGradient(
				GradientStartColor(tt.StartColor),
				GradientEndColor(tt.EndColor),
				GradientAlignment(tt.Alignment),
			)
			img := grad.Generator(5, 5)

			for x := 0; x < 5; x++ {
				for y := 0; y < 5; y++ {
					//t.Error(img.At(x, y).RGBA())
					//t.Error(tt.Result)
					r, g, b, a := img.At(x, y).RGBA()
					built := color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}
					assert.True(t, reflect.DeepEqual(built, tt.Result[x][y]))
				}
			}
		})
	}

}
