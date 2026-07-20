//go:build !windows || !ci

package gl

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
)

func TestGetFragmentColor(t *testing.T) {
	var c color.Color

	r, g, b, a := getFragmentColor(c)
	assert.Equal(t, float32(0), r)
	assert.Equal(t, float32(0), g)
	assert.Equal(t, float32(0), b)
	assert.Equal(t, float32(0), a)

	c = color.NRGBA{R: 0x0, G: 0x66, B: 0x99, A: 0xff}
	r, g, b, a = getFragmentColor(c)
	assert.Equal(t, float32(0), r)
	assert.Equal(t, float32(0.4), g)
	assert.Equal(t, float32(0.6), b)
	assert.Equal(t, float32(1), a)

	c = color.NRGBA{R: 0x0, G: 0x66, B: 0x99, A: 0x99}
	r, g, b, a = getFragmentColor(c)
	assert.Equal(t, float32(0), r)
	assert.Equal(t, float32(0.3999898), g)
	assert.Equal(t, float32(0.59998477), b)
	assert.Equal(t, float32(0.6), a)
}

func TestInnerRect_Stretch(t *testing.T) {
	pos := fyne.NewPos(10, 10)
	size := fyne.NewSize(40, 40)

	innerSize, innerPos := rectInnerCoords(size, pos, canvas.ImageFillStretch, 0.0)

	assert.Equal(t, size, innerSize)
	assert.Equal(t, pos, innerPos)
}

func TestInnerRect_StretchIgnoreRatio(t *testing.T) {
	pos := fyne.NewPos(10, 10)
	size := fyne.NewSize(40, 40)

	innerSize, innerPos := rectInnerCoords(size, pos, canvas.ImageFillStretch, 2.0)

	assert.Equal(t, size, innerSize)
	assert.Equal(t, pos, innerPos)
}

func TestInnerRect_ContainScale(t *testing.T) {
	pos := fyne.NewPos(10, 10)
	size := fyne.NewSize(40, 40)

	innerSize, innerPos := rectInnerCoords(size, pos, canvas.ImageFillContain, 1.0)

	assert.Equal(t, size, innerSize)
	assert.Equal(t, pos, innerPos)
}

func TestInnerRect_ContainPillarbox(t *testing.T) {
	pos := fyne.NewPos(10, 10)
	size := fyne.NewSize(40, 40)

	innerSize, innerPos := rectInnerCoords(size, pos, canvas.ImageFillContain, 0.5)

	assert.Equal(t, fyne.NewSize(20, 40), innerSize)
	assert.Equal(t, fyne.NewPos(20, 10), innerPos)
}

func TestInnerRect_Original(t *testing.T) {
	// TODO add check for minsize somehow?
	pos := fyne.NewPos(10, 10)
	size := fyne.NewSize(40, 40)

	innerSize1, innerPos1 := rectInnerCoords(size, pos, canvas.ImageFillOriginal, 0.5)
	innerSize2, innerPos2 := rectInnerCoords(size, pos, canvas.ImageFillContain, 0.5)

	assert.Equal(t, innerSize2, innerSize1)
	assert.Equal(t, innerPos2, innerPos1)
}

func TestVisibleTextPixels(t *testing.T) {
	frame := fyne.NewSize(400, 100)

	offset, width := visibleTextPixels(fyne.NewPos(10, 0), fyne.NewSize(100, 20), frame, nil, 1)
	assert.Equal(t, 0, offset)
	assert.Equal(t, 100, width)

	stack := &internal.ClipStack{}
	clip := stack.Push(fyne.NewPos(100, 0), fyne.NewSize(200, 100))
	offset, width = visibleTextPixels(fyne.NewPos(50, 0), fyne.NewSize(1000, 20), frame, clip, 1)
	assert.Equal(t, 50, offset)
	assert.Equal(t, 200, width)

	offset, width = visibleTextPixels(fyne.NewPos(-500, 0), fyne.NewSize(1000, 20), frame, nil, 2)
	assert.Equal(t, 1000, offset)
	assert.Equal(t, 800, width)
}

func TestTextTextureWindow(t *testing.T) {
	offset, width := textTextureWindow(0, 400, 20000, 4096)
	assert.Equal(t, 0, offset)
	assert.Equal(t, 4096, width)

	offset, width = textTextureWindow(5000, 400, 20000, 4096)
	assert.Equal(t, 3152, offset)
	assert.Equal(t, 4096, width)

	offset, width = textTextureWindow(19800, 200, 20000, 4096)
	assert.Equal(t, 15904, offset)
	assert.Equal(t, 4096, width)

	cached := clippedTextTexture{offset: 3152, width: 4096, height: 40, scale: 1}
	assert.True(t, cached.covers(5000, 400, 40, 1))
	assert.False(t, cached.covers(7200, 400, 40, 1))
}

func TestVecRectCoordsWithPad_Shadow(t *testing.T) {
	p := &painter{pixScale: 1.0}
	rect := &canvas.Rectangle{}
	pos := fyne.NewPos(5, 5)
	frame := fyne.NewSize(100, 100)

	coords, bounds := p.vecRectCoordsWithPad(pos, rect, frame, 0, 0, rect.Shadow)
	assert.Equal(t, [4]float32{5, 5, 5, 5}, bounds)
	assert.Equal(t, [8]float32{
		-0.92, 0.92,
		-0.88, 0.92,
		-0.92, 0.88,
		-0.88, 0.88,
	}, coords)

	rect.Shadow = canvas.Shadow{
		Color:      color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Offset:     fyne.NewPos(30, -20),
		BlurRadius: 80,
	}

	coords, bounds = p.vecRectCoordsWithPad(pos, rect, frame, 0, 0, rect.Shadow)
	// Check that shadow paddings affect the normalized coordinates
	assert.Equal(t, [4]float32{5, 5, 5, 5}, bounds)
	assert.Equal(t, [8]float32{
		-1.9200001, 2.92,
		1.3199999, 2.92,
		-1.9200001, -0.32000005,
		1.3199999, -0.32000005,
	}, coords)
}
