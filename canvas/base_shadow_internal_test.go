package canvas

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestBaseShadow_ShadowPaddings(t *testing.T) {
	b := &baseShadow{
		baseObject:     baseObject{},
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 4,
		ShadowOffset:   fyne.NewPos(2, 3),
		ShadowType:     DropShadow,
	}

	pads := b.ShadowPaddings()
	assert.Equal(t, float32(6), pads[0])
	assert.Equal(t, float32(1), pads[1])
	assert.Equal(t, float32(2), pads[2])
	assert.Equal(t, float32(7), pads[3])
}

func TestBaseShadow_SizeAndPositionWithShadow(t *testing.T) {
	b := &baseShadow{
		baseObject:     baseObject{},
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 4,
		ShadowOffset:   fyne.NewPos(2, 3),
		ShadowType:     DropShadow,
	}
	size := fyne.NewSize(10, 20)
	totalSize, pos := b.SizeAndPositionWithShadow(size)
	pads := b.ShadowPaddings()

	assert.Equal(t, size.Width+pads[0]+pads[2], totalSize.Width)
	assert.Equal(t, size.Height+pads[1]+pads[3], totalSize.Height)
	assert.Equal(t, -pads[0], pos.X)
	assert.Equal(t, -pads[1], pos.Y)
}

func TestBaseShadow_ContentSizeAndPos(t *testing.T) {
	b := &baseShadow{
		baseObject:     baseObject{},
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 4,
		ShadowOffset:   fyne.NewPos(2, 3),
		ShadowType:     BoxShadow,
	}
	b.baseObject.Resize(fyne.NewSize(30, 40))
	b.baseObject.Move(fyne.NewPos(5, 6))

	pads := b.ShadowPaddings()
	contentSize := b.ContentSize()
	contentPos := b.ContentPos()

	assert.Equal(t, 30-pads[0]-pads[2], contentSize.Width)
	assert.Equal(t, 40-pads[1]-pads[3], contentSize.Height)
	assert.Equal(t, 5+pads[0], contentPos.X)
	assert.Equal(t, 6+pads[1], contentPos.Y)
}