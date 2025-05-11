package widget

import (
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewTextGrid(t *testing.T) {
	grid := NewTextGridFromString("A")

	assert.Len(t, grid.Rows, 1)
	assert.Len(t, grid.Rows[0].Cells, 1)
}

func TestTextGrid_Append(t *testing.T) {
	grid := NewTextGridFromString("Something\nElse")
	grid.Append("Newline")

	assert.Equal(t, "Something\nElse\nNewline", grid.Text())
}

func TestTextGrid_CursorLocationForPosition(t *testing.T) {
	grid := NewTextGridFromString("Position\nTest")

	row, col := grid.CursorLocationForPosition(fyne.NewPos(2, 2))
	assert.Equal(t, 0, row)
	assert.Equal(t, 0, col)

	row, col = grid.CursorLocationForPosition(fyne.NewPos(12, 12))
	assert.Equal(t, 0, row)
	assert.Equal(t, 1, col)

	row, col = grid.CursorLocationForPosition(fyne.NewPos(20, 20))
	assert.Equal(t, 1, row)
	assert.Equal(t, 2, col)
}

func TestTextGrid_PositionForCursorLocation(t *testing.T) {
	grid := NewTextGridFromString("Position\nTest")

	pos := grid.PositionForCursorLocation(0, 0)
	assert.True(t, pos.IsZero())

	col1Pos := grid.PositionForCursorLocation(0, 1)
	col2Pos := grid.PositionForCursorLocation(0, 2)
	assert.Equal(t, col1Pos.Y, col2Pos.Y)
	assert.Greater(t, col2Pos.X, col1Pos.X)

	col1Row1Pos := grid.PositionForCursorLocation(1, 1)
	assert.Equal(t, col1Row1Pos.X, col1Pos.X)
	assert.Greater(t, col1Row1Pos.Y, col1Pos.Y)
}

func TestTextGrid_Scroll(t *testing.T) {
	grid := NewTextGridFromString("Something\nElse")
	grid.Resize(fyne.NewSize(50, 20))
	test.AssertObjectRendersToMarkup(t, "textgrid/basic.xml", grid)

	scrolling := NewTextGridFromString("Something\nElse")
	scrolling.Scroll = widget.ScrollBoth
	scrolling.Resize(fyne.NewSize(50, 20))
	scrolling.Refresh()
	scrolling.scroll.ScrollToTop()
	test.AssertObjectRendersToMarkup(t, "textgrid/scroll.xml", scrolling)

	scrolling = NewTextGrid()
	scrolling.Scroll = widget.ScrollBoth
	scrolling.Resize(fyne.NewSize(50, 20))
	scrolling.SetText("Something\nElse")
	scrolling.scroll.ScrollToTop()
	test.AssertObjectRendersToMarkup(t, "textgrid/scroll.xml", scrolling)

	scrolling.Scroll = widget.ScrollNone
	scrolling.Resize(fyne.NewSize(50, 20))
	scrolling.Refresh()
	test.AssertObjectRendersToMarkup(t, "textgrid/basic.xml", grid)
}

func TestTextGrid_CreateRendererRows(t *testing.T) {
	grid := NewTextGrid()
	grid.Resize(fyne.NewSize(52, 22))
	wrap := test.TempWidgetRenderer(t, grid).(*textGridRenderer).text
	rend := test.TempWidgetRenderer(t, wrap).(*textGridContentRenderer)
	rend.Refresh()

	row := rend.visible[0].(fyne.Widget)
	rr := test.TempWidgetRenderer(t, row).(*textGridRowRenderer)
	assert.Len(t, rr.objects, 18)
}

func TestTextGrid_Row(t *testing.T) {
	grid := NewTextGridFromString("Ab\nC")
	test.TempWidgetRenderer(t, grid).Refresh()

	assert.NotNil(t, grid.Row(0))
	assert.Len(t, grid.Row(0).Cells, 2)
	assert.Equal(t, 'b', grid.Row(0).Cells[1].Rune)
}

func TestTextGrid_Rows(t *testing.T) {
	grid := NewTextGridFromString("Ab\nC")
	test.TempWidgetRenderer(t, grid).Refresh()

	assert.Len(t, grid.Rows, 2)
	assert.Len(t, grid.Rows[0].Cells, 2)
}

func TestTextGrid_RowText(t *testing.T) {
	grid := NewTextGridFromString("Ab\nC")
	test.TempWidgetRenderer(t, grid).Refresh()

	assert.Equal(t, "Ab", grid.RowText(0))
	assert.Equal(t, "C", grid.RowText(1))
}

func TestTextGrid_SetText(t *testing.T) {
	grid := NewTextGrid()
	grid.Resize(fyne.NewSize(20, 20))
	text := "\n\n\n\n\n\n\n\n\n\n\n\n"
	grid.SetText(text) // goes beyond the current view size - don't crash

	assert.Len(t, grid.Rows, 13)
	assert.Empty(t, grid.Rows[1].Cells)
}

func TestTextGrid_SetText_Overflow(t *testing.T) {
	grid := NewTextGrid()
	grid.Scroll = widget.ScrollNone
	grid.SetText("Hello Long\nthere")

	assert.Len(t, grid.Rows, 2)
	assert.Len(t, grid.Rows[1].Cells, 5)
	render := test.WidgetRenderer(test.WidgetRenderer(grid).(*textGridRenderer).text).(*textGridContentRenderer)
	assert.Equal(t, 3, len(render.visible))
	row0 := test.WidgetRenderer(render.visible[0].(*textGridRow)).(*textGridRowRenderer)
	row1 := test.WidgetRenderer(render.visible[1].(*textGridRow)).(*textGridRowRenderer)
	row2 := test.WidgetRenderer(render.visible[2].(*textGridRow)).(*textGridRowRenderer)
	assert.Equal(t, "H", row0.objects[1].(*canvas.Text).Text)
	assert.Equal(t, "g", row0.objects[28].(*canvas.Text).Text)
	assert.Equal(t, "t", row1.objects[1].(*canvas.Text).Text)
	assert.Equal(t, " ", row2.objects[1].(*canvas.Text).Text)

	grid.SetText("Replace")

	assert.Len(t, grid.Rows, 1)
	assert.Equal(t, 2, len(render.visible))
	assert.Len(t, grid.Rows[0].Cells, 7)

	assert.Equal(t, "R", row0.objects[1].(*canvas.Text).Text)
	assert.Equal(t, " ", row0.objects[28].(*canvas.Text).Text)
	assert.Equal(t, " ", row1.objects[1].(*canvas.Text).Text)
}

func TestTextGrid_SetRowStyle(t *testing.T) {
	grid := NewTextGridFromString("Abc")
	grid.SetRowStyle(0, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})

	assert.NotNil(t, grid.Rows[0].Style)
	assert.Equal(t, color.White, grid.Rows[0].Style.TextColor())
	assert.Equal(t, color.Black, grid.Rows[0].Style.BackgroundColor())
}

func TestTextGrid_SetStyle(t *testing.T) {
	grid := NewTextGridFromString("Abc")
	grid.SetStyle(0, 1, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})

	assert.Nil(t, grid.Rows[0].Cells[0].Style)
	assert.Equal(t, color.White, grid.Rows[0].Cells[1].Style.TextColor())
	assert.Equal(t, color.Black, grid.Rows[0].Cells[1].Style.BackgroundColor())
}

func TestTextGrid_SetStyleRange(t *testing.T) {
	grid := NewTextGridFromString("Ab\ncd\nef")
	grid.SetStyleRange(0, 1, 2, 0, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})

	assert.Nil(t, grid.Rows[0].Cells[0].Style)
	assert.Equal(t, color.White, grid.Rows[0].Cells[1].Style.TextColor())
	assert.Equal(t, color.Black, grid.Rows[0].Cells[1].Style.BackgroundColor())
	assert.Equal(t, color.White, grid.Rows[1].Cells[0].Style.TextColor())
	assert.Equal(t, color.Black, grid.Rows[1].Cells[0].Style.BackgroundColor())
	assert.Equal(t, color.White, grid.Rows[1].Cells[1].Style.TextColor())
	assert.Equal(t, color.Black, grid.Rows[1].Cells[1].Style.BackgroundColor())
	assert.Equal(t, color.White, grid.Rows[2].Cells[0].Style.TextColor())
	assert.Equal(t, color.Black, grid.Rows[2].Cells[0].Style.BackgroundColor())
	assert.Nil(t, grid.Rows[2].Cells[1].Style)
}

func TestTextGrid_SetStyleRange_Overflow(t *testing.T) {
	grid := NewTextGridFromString("Ab\ncd")

	grid.SetStyleRange(-2, 0, -1, 2, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	grid.SetStyleRange(2, 2, 4, 2, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	assert.Nil(t, grid.Rows[0].Cells[0].Style)
	assert.Nil(t, grid.Rows[0].Cells[1].Style)
	assert.Nil(t, grid.Rows[1].Cells[0].Style)
	assert.Nil(t, grid.Rows[1].Cells[1].Style)

	grid.SetStyleRange(-2, 0, 0, 0, &CustomTextGridStyle{FGColor: color.Black, BGColor: color.White})
	grid.SetStyleRange(1, 1, 4, 0, &CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	assert.Equal(t, color.Black, grid.Rows[0].Cells[0].Style.TextColor())
	assert.Equal(t, color.White, grid.Rows[0].Cells[0].Style.BackgroundColor())
	assert.Nil(t, grid.Rows[0].Cells[1].Style)
	assert.Nil(t, grid.Rows[1].Cells[0].Style)
	assert.Equal(t, color.White, grid.Rows[1].Cells[1].Style.TextColor())
	assert.Equal(t, color.Black, grid.Rows[1].Cells[1].Style.BackgroundColor())
}

func TestTextGrid_Text(t *testing.T) {
	grid := NewTextGrid()
	assert.Equal(t, "", grid.Text())

	input := "Hello\nthere"
	grid.SetText(input)
	assert.Equal(t, input, grid.Text())
}

func TestTextGridRenderer_Resize(t *testing.T) {
	grid := NewTextGridFromString("1\n2")
	grid.ShowLineNumbers = true

	renderer := test.TempWidgetRenderer(t, grid)
	min := renderer.MinSize()

	grid.Resize(fyne.NewSize(100, 250))
	assert.Equal(t, min, renderer.MinSize())
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
	grid.Resize(fyne.NewSize(30, 42)) // causes refresh
	wrap := test.TempWidgetRenderer(t, grid).(*textGridRenderer).text
	rend := test.TempWidgetRenderer(t, wrap).(*textGridContentRenderer)

	assert.Equal(t, 2, rend.text.rows)

	row := rend.visible[0].(fyne.Widget)
	rend2 := test.TempWidgetRenderer(t, row).(*textGridRowRenderer)
	assert.Equal(t, 3, rend2.cols)
}

func TestTextGridRender_Whitespace(t *testing.T) {
	grid := NewTextGridFromString("A b\nc")
	grid.ShowWhitespace = true
	grid.Resize(fyne.NewSize(56, 42)) // causes refresh

	assertGridContent(t, grid, `A·b↵
c`)
}

func TestTextGridRender_WhitespaceTab(t *testing.T) {
	grid := NewTextGridFromString("A\n\tb")
	grid.ShowWhitespace = true
	grid.Resize(fyne.NewSize(56, 42)) // causes refresh

	assertGridContent(t, grid, `A↵
→···b`)
	assert.Equal(t, "A\n\tb", grid.Text())
}

func TestTextGridRender_RowColor(t *testing.T) {
	grid := NewTextGridFromString("Ab ")
	customStyle := &CustomTextGridStyle{FGColor: color.Black}
	grid.Rows[0].Style = customStyle
	grid.ShowWhitespace = true
	grid.Resize(fyne.NewSize(56, 22)) // causes refresh

	assertGridStyle(t, grid, "112", map[string]TextGridStyle{"1": customStyle, "2": TextGridStyleWhitespace})
}

func TestTextGridRender_Style(t *testing.T) {
	grid := NewTextGridFromString("Abcd ")
	boldStyle := &CustomTextGridStyle{TextStyle: fyne.TextStyle{Bold: true}}
	italicStyle := &CustomTextGridStyle{TextStyle: fyne.TextStyle{Italic: true}}
	boldItalicStyle := &CustomTextGridStyle{TextStyle: fyne.TextStyle{Bold: true, Italic: true}}
	grid.Rows[0].Cells[1].Style = boldStyle
	grid.Rows[0].Cells[2].Style = italicStyle
	grid.Rows[0].Cells[3].Style = boldItalicStyle
	grid.ShowWhitespace = true
	grid.Resize(fyne.NewSize(56, 22)) // causes refresh

	assertGridStyle(t, grid, "0123", map[string]TextGridStyle{"1": boldStyle, "2": italicStyle, "3": boldItalicStyle})
}

func TestTextGridRender_TextColor(t *testing.T) {
	grid := NewTextGridFromString("Ab ")
	customStyle := &CustomTextGridStyle{FGColor: color.Black}
	grid.Rows[0].Cells[1].Style = customStyle
	grid.ShowWhitespace = true
	grid.Resize(fyne.NewSize(56, 22)) // causes refresh

	currentTextColor := TextGridStyleWhitespace.TextColor()
	assertGridStyle(t, grid, " 12", map[string]TextGridStyle{"1": customStyle, "2": TextGridStyleWhitespace})

	test.WithTestTheme(t, func() {
		grid.Refresh()
		assert.NotEqual(t, TextGridStyleWhitespace.TextColor(), currentTextColor)
		assertGridStyle(t, grid, " 12", map[string]TextGridStyle{"1": customStyle, "2": TextGridStyleWhitespace})
	})
}

func assertGridContent(t *testing.T, g *TextGrid, expected string) {
	lines := strings.Split(expected, "\n")
	wrap := test.TempWidgetRenderer(t, g).(*textGridRenderer).text
	renderer := test.TempWidgetRenderer(t, wrap).(*textGridContentRenderer)

	for y, line := range lines {
		x := 0 // rune count - using index below would be offset into string bytes
		for _, r := range line {
			row := renderer.visible[y].(fyne.Widget)
			rend2 := test.TempWidgetRenderer(t, row).(*textGridRowRenderer)

			_, fg := rendererCell(rend2, x)
			assert.Equal(t, string(r), string([]rune(fg.Text)[0]))
			x++
		}
	}
}

func assertGridStyle(t *testing.T, g *TextGrid, content string, expectedStyles map[string]TextGridStyle) {
	lines := strings.Split(content, "\n")
	wrap := test.TempWidgetRenderer(t, g).(*textGridRenderer).text
	renderer := test.TempWidgetRenderer(t, wrap).(*textGridContentRenderer)

	for y, line := range lines {
		x := 0 // rune count - using index below would be offset into string bytes

		row := renderer.visible[y].(fyne.Widget)
		rend2 := test.TempWidgetRenderer(t, row).(*textGridRowRenderer)

		for _, r := range line {
			expected := expectedStyles[string(r)]
			bg, fg := rendererCell(rend2, x)

			if r == ' ' {
				assert.Equal(t, theme.Color(theme.ColorNameForeground), fg.Color)
				assert.Equal(t, color.Transparent, bg.FillColor)
			} else if expected != nil {
				if expected.TextColor() == nil {
					assert.Equal(t, theme.Color(theme.ColorNameForeground), fg.Color)
				} else {
					assert.Equal(t, expected.TextColor(), fg.Color)
				}

				if expected.BackgroundColor() == nil {
					assert.Equal(t, color.Transparent, bg.FillColor)
				} else {
					assert.Equal(t, expected.BackgroundColor(), bg.FillColor)
				}
			}

			style := fyne.TextStyle{}
			if expected != nil {
				style = expected.Style()
			}
			style.Monospace = true
			assert.Equal(t, style, fg.TextStyle)
			x++
		}
	}
}

func rendererCell(r *textGridRowRenderer, col int) (*canvas.Rectangle, *canvas.Text) {
	i := col * 3
	return r.objects[i].(*canvas.Rectangle), r.objects[i+1].(*canvas.Text)
}
