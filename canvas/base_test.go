package canvas

import (
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

// These tests group a few methods by function

func TestBase_Resize(t *testing.T) {
	b := &baseObject{}
	targetSize := fyne.NewSize(50, 50)
	b.Resize(targetSize)

	assert.Equal(t, b.Size(), targetSize)
}

func TestBase_Move(t *testing.T) {
	originPosition := fyne.NewPos(0, 0)
	targetPosition := fyne.NewPos(10, 15)
	b := &baseObject{
		position: originPosition,
	}
	assert.Equal(t, b.Position(), originPosition)
	b.Move(targetPosition)
	assert.Equal(t, b.Position(), targetPosition)
}

func TestBase_MinSize(t *testing.T) {
	zeroSizetest := fyne.NewSize(1, 1)
	minSize := fyne.NewSize(5, 5)
	b := &baseObject{}

	assert.Equal(t, b.MinSize(), zeroSizetest)

	b.SetMinSize(minSize)

	assert.Equal(t, b.MinSize(), minSize)

}

func TestBase_Visible(t *testing.T) {
	b := &baseObject{}
	assert.Equal(t, b.Visible(), true)

	b.Hide()

	assert.Equal(t, b.Visible(), false)

	b.Show()

	assert.Equal(t, b.Visible(), true)
}
