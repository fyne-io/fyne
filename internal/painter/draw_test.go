package painter_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
	"github.com/stretchr/testify/assert"
)

func TestPaint_GetCornerRadius(t *testing.T) {
	cases := []struct {
		perCorner, base, want float32
	}{
		{0, 5, 5},
		{3, 5, 3},
		{0, 0, 0},
		{7, 0, 7},
	}
	for _, c := range cases {
		got := painter.GetCornerRadius(c.perCorner, c.base)
		if got != c.want {
			t.Errorf("GetCornerRadius(%v, %v) = %v, want %v", c.perCorner, c.base, got, c.want)
		}
	}
}

func TestPaint_GetMaximumRadius(t *testing.T) {
	cases := []struct {
		size fyne.Size
		want float32
	}{
		{fyne.NewSize(10, 20), 5},
		{fyne.NewSize(30, 30), 15},
		{fyne.NewSize(0, 0), 0},
		{fyne.NewSize(100, 50), 25},
	}
	for _, c := range cases {
		got := painter.GetMaximumRadius(c.size)
		if got != c.want {
			t.Errorf("GetMaximumRadius(%v) = %v, want %v", c.size, got, c.want)
		}
	}
}

func TestPaint_GetMaximumCornerRadius(t *testing.T) {
	cases := []struct {
		radius, adjW, adjH float32
		size               fyne.Size
		want               float32
	}{
		{5, 5, 5, fyne.NewSquareSize(10), 5},
		{10, 2, 2, fyne.NewSquareSize(10), 8},
		{20, 0, 0, fyne.NewSquareSize(20), 20},
		{15, 5, 10, fyne.NewSquareSize(20), 10},
		{30, 2, 10, fyne.NewSquareSize(20), 10},
		{8, 4, 4, fyne.NewSquareSize(16), 8},
		{12, 6, 6, fyne.NewSize(24, 16), 10},
		{16, 8, 8, fyne.NewSize(32, 16), 8},
		{20, 10, 5, fyne.NewSize(40, 20), 15},
		{10, 5, 5, fyne.NewSize(40, 20), 10},
		{10, 20, 5, fyne.NewSize(20, 40), 10},
		{13, 5, 2, fyne.NewSize(17, 20), 12},
		{13, 2, 5, fyne.NewSize(17, 20), 13},
		{canvas.RadiusMaximum, 3, 5, fyne.NewSize(36, 16), 11},
		{canvas.RadiusMaximum, canvas.RadiusMaximum, canvas.RadiusMaximum, fyne.NewSize(36, 16), 8},
		{canvas.RadiusMaximum, 0, 0, fyne.NewSize(10, 14), 10},
		{canvas.RadiusMaximum, 1, 0, fyne.NewSize(10, 14), 9},
		{canvas.RadiusMaximum, canvas.RadiusMaximum, 4, fyne.NewSize(30, 20), 15},
	}
	for _, c := range cases {
		got := painter.GetMaximumCornerRadius(c.radius, c.adjW, c.adjH, c.size)
		if got != c.want {
			t.Errorf("GetMaximumCornerRadius(%v, %v, %v, %v) = %v, want %v", c.radius, c.adjW, c.adjH, c.size, got, c.want)
		}
	}
}

func TestPaint_GetMaximumRadiusArc(t *testing.T) {
	cases := []struct {
		outer, inner, sweep, want float32
	}{
		{10, 0, 180, 3.75},
		{10, 8, 180, 1},
		{20, 10, 90, 5},
		{30, 20, 360, 5},
		{50, 45, 45, 2.5},
	}
	for _, c := range cases {
		got := painter.GetMaximumRadiusArc(c.outer, c.inner, c.sweep)
		if got != c.want {
			t.Errorf("GetMaximumRadiusArc(%v, %v, %v) = %v, want %v", c.outer, c.inner, c.sweep, got, c.want)
		}
	}
}

func TestPaint_NormalizeArcAngles(t *testing.T) {
	cases := []struct {
		start, end, wantStart, wantEnd float32
	}{
		{0, 90, 90, 0},
		{180, 270, -90, -180},
		{360, 0, -270, 90},
		{-90, -180, 180, 270},
		{45, 45, 45, 45},
		{-360, 0, 450, 90},
		{-50, -220, 140, 310},
		{-180, 0, 270, 90},
		{90, -90, 0, 180},
		{270, 90, -180, 0},
		{0, -90, 90, 180},
		{110, 150, -20, -60},
	}
	for _, c := range cases {
		gotStart, gotEnd := painter.NormalizeArcAngles(c.start, c.end)
		if gotStart != c.wantStart || gotEnd != c.wantEnd {
			t.Errorf("NormalizeArcAngles(%v, %v) = (%v, %v), want (%v, %v)", c.start, c.end, gotStart, gotEnd, c.wantStart, c.wantEnd)
		}
	}
}

func TestShadow_ShadowPaddings_Empty(t *testing.T) {
	b := canvas.Shadow{}
	expected := [4]float32{0, 0, 0, 0}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetOnlyLeft(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 3,
		ShadowOffset:   fyne.NewPos(5, 0),
		ShadowType:     canvas.BoxShadow,
	}
	expected := [4]float32{8, 3, 0, 3}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetOnlyRight(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 2,
		ShadowOffset:   fyne.NewPos(-6, 0),
		ShadowType:     canvas.DropShadow,
	}
	expected := [4]float32{0, 2, 8, 2}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetOnlyTop(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 5,
		ShadowOffset:   fyne.NewPos(0, -4),
		ShadowType:     canvas.BoxShadow,
	}
	expected := [4]float32{5, 9, 5, 1}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetOnlyBottom(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 1,
		ShadowOffset:   fyne.NewPos(0, 7),
		ShadowType:     canvas.DropShadow,
	}
	expected := [4]float32{1, 0, 1, 8}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetLeftTop(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 4,
		ShadowOffset:   fyne.NewPos(3, -2),
		ShadowType:     canvas.DropShadow,
	}
	expected := [4]float32{7, 6, 1, 2}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetLeftBottom(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 2,
		ShadowOffset:   fyne.NewPos(4, 5),
		ShadowType:     canvas.BoxShadow,
	}
	expected := [4]float32{6, 0, 0, 7}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetRightTop(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 3,
		ShadowOffset:   fyne.NewPos(-3, -2),
	}
	expected := [4]float32{0, 5, 6, 1}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetRightBottom(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 4,
		ShadowOffset:   fyne.NewPos(2, 3),
		ShadowType:     canvas.DropShadow,
	}
	expected := [4]float32{6, 1, 2, 7}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetRightBottom2(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 2,
		ShadowOffset:   fyne.NewPos(-4, 5),
		ShadowType:     canvas.DropShadow,
	}
	expected := [4]float32{0, 0, 6, 7}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_NoOffsetOnlySoftness(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 5,
		ShadowOffset:   fyne.NewPos(0, 0),
		ShadowType:     canvas.BoxShadow,
	}
	expected := [4]float32{5, 5, 5, 5}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_NoOffsetOnlySoftness2(t *testing.T) {
	b := canvas.Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 9,
		ShadowOffset:   fyne.NewPos(0, 0),
		ShadowType:     canvas.BoxShadow,
	}
	expected := [4]float32{9, 9, 9, 9}
	pads := painter.GetShadowPaddings(b)
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}
