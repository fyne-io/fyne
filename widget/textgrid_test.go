package widget

import (
	"image/color"
	"strings"
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

func TestTextGrid_SetText(t *testing.T) {
	grid := NewTextGrid()
	grid.SetText("Hello\nthere")

	assert.Equal(t, 2, len(grid.Content))
	assert.Equal(t, 5, len(grid.Content[1]))
}

func TestTextGrid_Text(t *testing.T) {
	input := "Hello\nthere"
	grid := NewTextGrid()
	grid.SetText(input)
	assert.Equal(t, input, grid.Text())
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

func TestTextGrid_SetStyleRange_Overflow(t *testing.T) {
	grid := NewTextGridFromString("Ab\ncd")

	grid.SetStyleRange(-2, 0, -1, 2, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	grid.SetStyleRange(2, 2, 4, 2, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	assert.Nil(t, grid.Content[0][0].Style)
	assert.Nil(t, grid.Content[0][1].Style)
	assert.Nil(t, grid.Content[1][0].Style)
	assert.Nil(t, grid.Content[1][1].Style)

	grid.SetStyleRange(-2, 0, 0, 0, &CustomTextGridStyle{FGColor: color.Black, BGColor: color.White})
	grid.SetStyleRange(1, 1, 4, 0, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	assert.Equal(t, color.Black, grid.Content[0][0].Style.TextColor())
	assert.Equal(t, color.White, grid.Content[0][0].Style.BackgroundColor())
	assert.Nil(t, grid.Content[0][1].Style)
	assert.Nil(t, grid.Content[1][0].Style)
	assert.Equal(t, color.White, grid.Content[1][1].Style.TextColor())
	assert.Equal(t, color.Black, grid.Content[1][1].Style.BackgroundColor())
}

func TestTextGrid_CreateRendererRows(t *testing.T) {
	grid := NewTextGrid()
	grid.Resize(fyne.NewSize(56, 22))
	rend := test.WidgetRenderer(grid).(*textGridRenderer)
	rend.Refresh()

	assert.Equal(t, 8, len(rend.objects))
}

func TestTextGridRenderer_ShowLineNumbers(t *testing.T) {
	grid := NewTextGridFromString("1\n2\n3\n4\n5\n6\n7\n8\n9\n10")
	grid.ShowLineNumbers = true
	grid.Resize(fyne.NewSize(100, 250))

	assertGridContent(t, grid, ` 1|1
 2|
 3|3
 4|4
 5|5
 6|6
 7|7
 8|8
 9|9
10|10
`)
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

	assertGridContent(t, grid, `A·b↵
c`)
}

func TestTextGridRender_TextColor(t *testing.T) {
	grid := NewTextGridFromString("Ab ")
	customStyle := &CustomTextGridStyle{FGColor: color.Black}
	grid.Content[0][1].Style = customStyle
	grid.ShowWhitespace = true
	grid.Resize(fyne.NewSize(56, 22)) // causes refresh

	assertGridStyle(t, grid, " 12", map[string]TextGridStyle{"1": customStyle, "2": TextGridStyleWhitespace})
}

func assertGridContent(t *testing.T, g *TextGrid, expected string) {
	lines := strings.Split(expected, "\n")
	renderer := test.WidgetRenderer(g).(*textGridRenderer)

	for y, line := range lines {
		x := 0 // rune count - using index below would be offset into string bytes
		for _, r := range line {
			_, fg := rendererCell(renderer, y, x)
			assert.Equal(t, r, []rune(fg.Text)[0])
			x++
		}
	}
}

func assertGridStyle(t *testing.T, g *TextGrid, expected string, expectedStyles map[string]TextGridStyle) {
	lines := strings.Split(expected, "\n")
	renderer := test.WidgetRenderer(g).(*textGridRenderer)

	for y, line := range lines {
		x := 0 // rune count - using index below would be offset into string bytes
		for _, r := range line {
			expected := expectedStyles[string(r)]
			bg, fg := rendererCell(renderer, y, x)

			if r == ' ' {
				assert.Equal(t, theme.TextColor(), fg.Color)
				assert.Equal(t, color.Transparent, bg.FillColor)
			} else {
				if expected.TextColor() == nil {
					assert.Equal(t, theme.TextColor(), fg.Color)
				} else {
					assert.Equal(t, expected.TextColor(), fg.Color)
				}

				if expected.BackgroundColor() == nil {
					assert.Equal(t, color.Transparent, bg.FillColor)
				} else {
					assert.Equal(t, expected.BackgroundColor(), bg.FillColor)
				}
			}
			x++
		}
	}
}

func rendererCell(r *textGridRenderer, row, col int) (*canvas.Rectangle, *canvas.Text) {
	i := (row*r.cols + col) * 2
	return r.objects[i].(*canvas.Rectangle), r.objects[i+1].(*canvas.Text)
}
