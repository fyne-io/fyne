package canvas

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestShadow_ShadowPaddings_OffsetOnlyLeft(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 3,
		ShadowOffset:   fyne.NewPos(5, 0),
		ShadowType:     BoxShadow,
	}
	expected := [4]float32{8, 3, 0, 3}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetOnlyRight(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 2,
		ShadowOffset:   fyne.NewPos(-6, 0),
		ShadowType:     DropShadow,
	}
	expected := [4]float32{0, 2, 8, 2}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetOnlyTop(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 5,
		ShadowOffset:   fyne.NewPos(0, -4),
		ShadowType:     BoxShadow,
	}
	expected := [4]float32{5, 9, 5, 1}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetOnlyBottom(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 1,
		ShadowOffset:   fyne.NewPos(0, 7),
		ShadowType:     DropShadow,
	}
	expected := [4]float32{1, 0, 1, 8}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetLeftTop(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 4,
		ShadowOffset:   fyne.NewPos(3, -2),
		ShadowType:     DropShadow,
	}
	expected := [4]float32{7, 6, 1, 2}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetLeftBottom(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 2,
		ShadowOffset:   fyne.NewPos(4, 5),
		ShadowType:     BoxShadow,
	}
	expected := [4]float32{6, 0, 0, 7}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetRightTop(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 3,
		ShadowOffset:   fyne.NewPos(-3, -2),
	}
	expected := [4]float32{0, 5, 6, 1}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetRightBottom(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 4,
		ShadowOffset:   fyne.NewPos(2, 3),
		ShadowType:     DropShadow,
	}
	expected := [4]float32{6, 1, 2, 7}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_OffsetRightBottom2(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 2,
		ShadowOffset:   fyne.NewPos(-4, 5),
		ShadowType:     DropShadow,
	}
	expected := [4]float32{0, 0, 6, 7}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_NoOffsetOnlySoftness(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 5,
		ShadowOffset:   fyne.NewPos(0, 0),
		ShadowType:     BoxShadow,
	}
	expected := [4]float32{5, 5, 5, 5}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}

func TestShadow_ShadowPaddings_NoOffsetOnlySoftness2(t *testing.T) {
	b := &Shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 9,
		ShadowOffset:   fyne.NewPos(0, 0),
		ShadowType:     BoxShadow,
	}
	expected := [4]float32{9, 9, 9, 9}
	pads := b.ShadowPaddings()
	assert.Equal(t, expected[0], pads[0], "left")
	assert.Equal(t, expected[1], pads[1], "top")
	assert.Equal(t, expected[2], pads[2], "right")
	assert.Equal(t, expected[3], pads[3], "bottom")
}
