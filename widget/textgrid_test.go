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
	Renderer(grid).Refresh()

	assert.Equal(t, 1, len(grid.Content))
	assert.Equal(t, 1, len(grid.Content[0]))
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
	rend := test.WidgetRenderer(grid).(*textGridRender)
	rend.Refresh()

	assert.Equal(t, 8, len(rend.objects))
}

func TestTextGridRender_Size(t *testing.T) {
	grid := NewTextGrid()
	grid.Resize(fyne.NewSize(32, 42)) // causes refresh
	rend := test.WidgetRenderer(grid).(*textGridRender)

	assert.Equal(t, 2, rend.cols)
	assert.Equal(t, 2, rend.rows)
}

func TestTextGridRender_Whitespace(t *testing.T) {
	grid := NewTextGridFromString("A b\nc")
	grid.Whitespace = true
	grid.Resize(fyne.NewSize(56, 42)) // causes refresh
	rend := test.WidgetRenderer(grid).(*textGridRender)

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
	grid.Whitespace = true
	grid.Resize(fyne.NewSize(56, 22)) // causes refresh
	rend := test.WidgetRenderer(grid).(*textGridRender)

	assert.Equal(t, 4, rend.cols)
	assert.Equal(t, 1, rend.rows)
	assert.Equal(t, theme.TextColor(), rend.objects[1].(*canvas.Text).Color)
	assert.Equal(t, color.Black, rend.objects[3].(*canvas.Text).Color)
	assert.Equal(t, TextGridStyleWhitespace.TextColor(), rend.objects[5].(*canvas.Text).Color)
}
