package widget

import (
	"testing"

	"fyne.io/fyne/canvas"
	"github.com/stretchr/testify/assert"
)

func TestNewTextGrid(t *testing.T) {
	grid := NewTextGrid("A")
	Renderer(grid).Refresh()

	assert.Equal(t, 1, grid.rows())
	assert.Equal(t, 1, grid.maxCols)
}

func TestTextGrid_Rows(t *testing.T) {
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

func TestTextGridRender_Cols(t *testing.T) {
	grid := NewTextGrid("Ab")
	grid.LineNumbers = true
	rend := Renderer(grid).(*textGridRender)
	rend.Refresh()

	assert.Equal(t, 4, rend.cols) // 1 for "1" and 1 space
}

func TestTextGridRender_ColsLong(t *testing.T) {
	grid := NewTextGrid("1\n2\n3\n4\n5\n6\n7\n8\n9\n10")
	grid.LineNumbers = true
	rend := Renderer(grid).(*textGridRender)
	rend.Refresh()

	assert.Equal(t, 5, rend.cols) // 2 for "10" and 1 space
}

func TestTextGridRender_Whitespace(t *testing.T) {
	grid := NewTextGrid("A b\nc")
	grid.Whitespace = true
	rend := Renderer(grid).(*textGridRender)
	rend.Refresh()

	assert.Equal(t, 4, rend.cols)                                                       // 1 for newline
	assert.Equal(t, rend.objects[1].(*canvas.Text).Text, string(textAreaSpaceSymbol))   // col 2 is space
	assert.Equal(t, rend.objects[3].(*canvas.Text).Text, string(textAreaNewLineSymbol)) // col 4 is newline
	assert.Equal(t, rend.objects[5].(*canvas.Text).Text, string(textAreaNewLineSymbol)) // col 2 on line 2
}
