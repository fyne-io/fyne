package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewTextGrid(t *testing.T) {
	grid := NewTextGridFromString("A")
	test.WidgetRenderer(grid).Refresh()

	assert.Equal(t, 1, len(grid.Content))
	assert.Equal(t, 1, len(grid.Content[0]))
}

func TestTextGrid_LineNumbers(t *testing.T) {
	grid := NewTextGridFromString("1\n2\n3\n4\n5\n6\n7\n8\n9\n10")
	grid.ShowLineNumbers = true

	grid.Resize(fyne.NewSize(100, 250))
	r := test.WidgetRenderer(grid).(*textGridRenderer)

	assert.Equal(t, ' ', rendererCellRune(r, 0, 0))
	assert.Equal(t, '1', rendererCellRune(r, 0, 1))
	assert.Equal(t, '|', rendererCellRune(r, 0, 2))
	assert.Equal(t, '1', rendererCellRune(r, 0, 3))
	assert.Equal(t, '2', rendererCellRune(r, 1, 3))

	assert.Equal(t, '1', rendererCellRune(r, 9, 0))
	assert.Equal(t, '0', rendererCellRune(r, 9, 1))
	assert.Equal(t, '|', rendererCellRune(r, 9, 2))
	assert.Equal(t, '1', rendererCellRune(r, 9, 3))
	assert.Equal(t, '0', rendererCellRune(r, 9, 4))
}

func TestTextGrid_SetText(t *testing.T) {
	grid := NewTextGrid()
	grid.SetText("Hello\nthere")

	assert.Equal(t, 2, len(grid.Content))
	assert.Equal(t, 5, len(grid.Content[1]))
}

func TestTextGrid_Rows(t *testing.T) {
	grid := NewTextGridFromString("Ab\nC")
	test.WidgetRenderer(grid).Refresh()

	assert.Equal(t, 2, len(grid.Content))
	assert.Equal(t, 2, len(grid.Content[0]))
}

func TestTextGrid_SetStyle(t *testing.T) {
	grid := NewTextGridFromString("Abc")
	grid.SetStyle(0, 1, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})

	assert.Nil(t, grid.Content[0][0].Style)
	assert.Equal(t, color.White, grid.Content[0][1].Style.TextColor())
	assert.Equal(t, color.Black, grid.Content[0][1].Style.BackgroundColor())
}

func TestTextGrid_SetStyleRange(t *testing.T) {
	grid := NewTextGridFromString("Ab\ncd")
	grid.SetStyleRange(0, 1, 1, 0, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})

	assert.Nil(t, grid.Content[0][0].Style)
	assert.Equal(t, color.White, grid.Content[0][1].Style.TextColor())
	assert.Equal(t, color.Black, grid.Content[0][1].Style.BackgroundColor())
	assert.Equal(t, color.White, grid.Content[1][0].Style.TextColor())
	assert.Equal(t, color.Black, grid.Content[1][0].Style.BackgroundColor())
	assert.Nil(t, grid.Content[1][1].Style)
}

func TestTextGrid_CreateRendererRows(t *testing.T) {
	grid := NewTextGrid()
	grid.Resize(fyne.NewSize(56, 22))
	rend := test.WidgetRenderer(grid).(*textGridRenderer)
	rend.Refresh()

	assert.Equal(t, 8, len(rend.objects))
}

func TestTextGridRender_Size(t *testing.T) {
	grid := NewTextGrid()
	grid.Resize(fyne.NewSize(32, 42)) // causes refresh
	rend := test.WidgetRenderer(grid).(*textGridRenderer)

	assert.Equal(t, 2, rend.cols)
	assert.Equal(t, 2, rend.rows)
}

func TestTextGridRender_Whitespace(t *testing.T) {
	grid := NewTextGridFromString("A b\nc")
	grid.ShowWhitespace = true
	grid.Resize(fyne.NewSize(56, 42)) // causes refresh
	rend := test.WidgetRenderer(grid).(*textGridRenderer)

	assert.Equal(t, 4, rend.cols)
	assert.Equal(t, 2, rend.rows)
	// indexes of text are at n*2+1 due to bg rects appearing before letter objects
	assert.Equal(t, string(textAreaSpaceSymbol), rend.objects[3].(*canvas.Text).Text)       // col 1 is space
	assert.Equal(t, string(textAreaNewLineSymbol), rend.objects[7].(*canvas.Text).Text)     // col 3 is newline
	assert.NotEqual(t, string(textAreaNewLineSymbol), rend.objects[11].(*canvas.Text).Text) // no newline on end of content
}

func TestTextGridRender_TextColor(t *testing.T) {
	grid := NewTextGridFromString("Ab ")
	grid.Content[0][1].Style = &CustomTextGridStyle{FGColor: color.Black}
	grid.ShowWhitespace = true
	grid.Resize(fyne.NewSize(56, 22)) // causes refresh
	rend := test.WidgetRenderer(grid).(*textGridRenderer)

	assert.Equal(t, 4, rend.cols)
	assert.Equal(t, 1, rend.rows)
	assert.Equal(t, theme.TextColor(), rend.objects[1].(*canvas.Text).Color)
	assert.Equal(t, color.Black, rend.objects[3].(*canvas.Text).Color)
	assert.Equal(t, TextGridStyleWhitespace.TextColor(), rend.objects[5].(*canvas.Text).Color)
}

func rendererCell(r *textGridRenderer, row, col int) (*canvas.Rectangle, *canvas.Text) {
	i := (row*r.cols + col) * 2
	return r.objects[i].(*canvas.Rectangle), r.objects[i+1].(*canvas.Text)
}

func rendererCellRune(r *textGridRenderer, row, col int) rune {
	_, text := rendererCell(r, row, col)
	return []rune(text.Text)[0]
}
