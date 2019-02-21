package gl

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

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

func TestDrawImage_Ratio(t *testing.T) {
	d := NewGLDriver()
	win := d.CreateWindow("Test")
	c := win.Canvas().(*glCanvas)

	img := canvas.NewImageFromResource(theme.FyneLogo())
	img.Resize(fyne.NewSize(10, 10))
	c.newGlImageTexture(img)
	assert.Equal(t, float32(1.0), c.aspects[img])
}

func TestDrawImage_Ratio2(t *testing.T) {
	d := NewGLDriver()
	win := d.CreateWindow("Test")
	c := win.Canvas().(*glCanvas)

	// make sure we haven't used the visual ratio
	img := canvas.NewImageFromResource(theme.FyneLogo())
	img.Resize(fyne.NewSize(20, 10))
	c.newGlImageTexture(img)
	assert.Equal(t, float32(1.0), c.aspects[img])
}
