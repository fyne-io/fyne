package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTextGrid(t *testing.T) {
	grid := NewTextGrid("A")
	Renderer(grid).Refresh()

	assert.Equal(t, 1, grid.rows())
	assert.Equal(t, 1, grid.maxCols)
}

func TestNewTextGrid_Rows(t *testing.T) {
	grid := NewTextGrid("Ab\nC")
	Renderer(grid).Refresh()

	assert.Equal(t, 2, grid.rows())
	assert.Equal(t, 2, grid.maxCols)
}

func TestTextGrid_CreateRendererRows(t *testing.T) {
	grid := NewTextGrid("Ab\nC")
	rend := Renderer(grid).(*textGridRender)
	rend.Refresh()

	assert.Equal(t, 4, len(rend.objects))
}
