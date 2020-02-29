package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"

	"github.com/stretchr/testify/assert"
)

func TestNewTextGrid(t *testing.T) {
	grid := NewTextGridFromString("A")
	Renderer(grid).Refresh()

	assert.Equal(t, 1, len(grid.Buffer))
	assert.Equal(t, 1, len(grid.Buffer[0]))
}

func TestTextGrid_SetText(t *testing.T) {
	grid := NewTextGrid()
	grid.SetText("Hello\nthere")

	assert.Equal(t, 2, len(grid.Buffer))
	assert.Equal(t, 5, len(grid.Buffer[1]))
}

func TestTextGrid_Rows(t *testing.T) {
	grid := NewTextGridFromString("Ab\nC")
	Renderer(grid).Refresh()

	assert.Equal(t, 2, len(grid.Buffer))
	assert.Equal(t, 2, len(grid.Buffer[0]))
}

func TestTextGrid_CreateRendererRows(t *testing.T) {
	grid := NewTextGrid()
	grid.Resize(fyne.NewSize(56, 22))
	rend := Renderer(grid).(*textGridRender)
	rend.Refresh()

	assert.Equal(t, 4, len(rend.objects))
}

func TestTextGridRender_Size(t *testing.T) {
	grid := NewTextGrid()
	grid.Resize(fyne.NewSize(32, 42))
	rend := Renderer(grid).(*textGridRender)
	rend.Refresh()

	assert.Equal(t, 2, rend.cols)
	assert.Equal(t, 2, rend.rows)
}

func TestTextGridRender_Whitespace(t *testing.T) {
	grid := NewTextGridFromString("A b\nc")
	grid.Resize(fyne.NewSize(56, 42))
	grid.Whitespace = true
	rend := Renderer(grid).(*textGridRender)
	rend.Refresh()

	assert.Equal(t, 4, rend.cols)
	assert.Equal(t, string(textAreaSpaceSymbol), rend.objects[1].(*canvas.Text).Text)   // col 2 is space
	assert.Equal(t, string(textAreaNewLineSymbol), rend.objects[3].(*canvas.Text).Text) // col 4 is newline
	assert.Equal(t, string(textAreaNewLineSymbol), rend.objects[5].(*canvas.Text).Text) // col 2 on line 2
}
