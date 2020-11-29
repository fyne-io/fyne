package widget

import (
	"fmt"
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestTable_Empty(t *testing.T) {
	table := &Table{}
	table.Resize(fyne.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh() // let's not crash :)
}

func TestTable_Cache(t *testing.T) {
	c := test.NewCanvas()
	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})
	c.SetContent(table)
	c.SetPadded(false)
	c.Resize(fyne.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh()
	assert.Equal(t, 6, len(cellRenderer.Objects()))
	assert.Equal(t, "Cell 0, 0", cellRenderer.Objects()[0].(*Label).Text)
	objRef := cellRenderer.Objects()[0].(*Label)

	test.Scroll(c, fyne.NewPos(10, 10), -150, -150)
	assert.Equal(t, 0, renderer.scroll.Offset.Y) // we didn't scroll as data shorter
	assert.Equal(t, 150, renderer.scroll.Offset.X)
	assert.Equal(t, 6, len(cellRenderer.Objects()))
	assert.Equal(t, "Cell 0, 1", cellRenderer.Objects()[0].(*Label).Text)
	assert.NotEqual(t, objRef, cellRenderer.Objects()[0].(*Label)) // we want to re-use visible cells without rewriting them
}

func TestTable_ChangeTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})
	content := test.WidgetRenderer(test.WidgetRenderer(table).(*tableRenderer).scroll.Content.(*tableCells)).(*tableCellsRenderer)
	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))
	test.AssertImageMatches(t, "table/theme_initial.png", w.Canvas().Capture())
	assert.Equal(t, "Cell 0, 0", content.Objects()[0].(*Label).Text)
	assert.Equal(t, NewLabel("placeholder").MinSize(), content.Objects()[0].(*Label).Size())

	test.WithTestTheme(t, func() {
		test.WidgetRenderer(table).Refresh()
		test.AssertImageMatches(t, "table/theme_changed.png", w.Canvas().Capture())
	})
	assert.Equal(t, NewLabel("placeholder").MinSize(), content.Objects()[0].(*Label).Size())
}

func TestTable_Filled(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() fyne.CanvasObject {
			r := canvas.NewRectangle(color.Black)
			r.SetMinSize(fyne.NewSize(30, 20))
			r.Resize(fyne.NewSize(30, 20))
			return r
		},
		func(TableCellID, fyne.CanvasObject) {})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))
	test.AssertImageMatches(t, "table/filled.png", w.Canvas().Capture())
}

func TestTable_Unselect(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})
	unselectedRow, unselectedColumn := 0, 0
	table.OnUnselected = func(id TableCellID) {
		unselectedRow = id.Row
		unselectedColumn = id.Col
	}
	table.selectedCell = &TableCellID{1, 1}
	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	table.Unselect(*table.selectedCell)
	assert.Equal(t, 1, unselectedRow)
	assert.Equal(t, 1, unselectedColumn)
	test.AssertImageMatches(t, "table/theme_initial.png", w.Canvas().Capture())
}

func TestTable_Refresh(t *testing.T) {
	displayText := "placeholder"
	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() fyne.CanvasObject {
			return NewLabel("template")
		},
		func(_ TableCellID, obj fyne.CanvasObject) {
			obj.(*Label).SetText(displayText)
		})
	table.Resize(fyne.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	assert.Equal(t, "placeholder", cellRenderer.(*tableCellsRenderer).Objects()[7].(*Label).Text)

	displayText = "replaced"
	table.Refresh()
	assert.Equal(t, "replaced", cellRenderer.(*tableCellsRenderer).Objects()[7].(*Label).Text)
}

func TestTable_Selection(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})
	assert.Nil(t, table.selectedCell)

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	selectedCol, selectedRow := 0, 0
	table.OnSelected = func(id TableCellID) {
		selectedCol = id.Col
		selectedRow = id.Row
	}
	test.TapCanvas(w.Canvas(), fyne.NewPos(35, 50))
	assert.Equal(t, 0, table.selectedCell.Col)
	assert.Equal(t, 1, table.selectedCell.Row)
	assert.Equal(t, 0, selectedCol)
	assert.Equal(t, 1, selectedRow)

	test.AssertRendersToMarkup(t, `
		<canvas padded size="180x180">
			<content>
				<widget pos="4,4" size="172x172" type="*widget.Table">
					<widget pos="4,4" size="168x168" type="*widget.ScrollContainer">
						<widget size="509x189" type="*widget.tableCells">
							<widget pos="4,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 0</text>
							</widget>
							<widget pos="106,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 1</text>
							</widget>
							<widget pos="4,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 0</text>
							</widget>
							<widget pos="106,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 1</text>
							</widget>
							<widget pos="4,80" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 2, 0</text>
							</widget>
							<widget pos="106,80" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 2, 1</text>
							</widget>
							<widget pos="4,118" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 3, 0</text>
							</widget>
							<widget pos="106,118" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 3, 1</text>
							</widget>
							<widget pos="4,156" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 4, 0</text>
							</widget>
							<widget pos="106,156" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 4, 1</text>
							</widget>
						</widget>
						<widget pos="162,0" size="6x168" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="3,0" size="3x149" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="0,168" size="168x0" type="*widget.Shadow">
							<linearGradient endColor="shadow" pos="0,-8" size="168x8"/>
						</widget>
						<widget pos="0,162" size="168x6" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="0,3" size="55x3" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="168,0" size="0x168" type="*widget.Shadow">
							<linearGradient angle="270" endColor="shadow" pos="-8,0" size="8x168"/>
						</widget>
					</widget>
					<rectangle fillColor="primary" pos="4,0" size="101x4"/>
					<rectangle fillColor="primary" pos="0,42" size="4x37"/>
					<widget pos="105,4" size="1x168" type="*widget.Separator">
						<rectangle fillColor="disabled" size="1x168"/>
					</widget>
					<widget pos="4,41" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,79" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,117" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,155" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
}

func TestTable_Select(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", id.Row, id.Col)
			c.(*Label).SetText(text)
		})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	selectedCol, selectedRow := 0, 0
	table.OnSelected = func(id TableCellID) {
		selectedCol = id.Col
		selectedRow = id.Row
	}
	table.Select(TableCellID{1, 0})
	assert.Equal(t, 0, table.selectedCell.Col)
	assert.Equal(t, 1, table.selectedCell.Row)
	assert.Equal(t, 0, selectedCol)
	assert.Equal(t, 1, selectedRow)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="180x180">
			<content>
				<widget pos="4,4" size="172x172" type="*widget.Table">
					<widget pos="4,4" size="168x168" type="*widget.ScrollContainer">
						<widget size="509x189" type="*widget.tableCells">
							<widget pos="4,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 0</text>
							</widget>
							<widget pos="106,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 1</text>
							</widget>
							<widget pos="4,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 0</text>
							</widget>
							<widget pos="106,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 1</text>
							</widget>
							<widget pos="4,80" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 2, 0</text>
							</widget>
							<widget pos="106,80" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 2, 1</text>
							</widget>
							<widget pos="4,118" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 3, 0</text>
							</widget>
							<widget pos="106,118" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 3, 1</text>
							</widget>
							<widget pos="4,156" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 4, 0</text>
							</widget>
							<widget pos="106,156" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 4, 1</text>
							</widget>
						</widget>
						<widget pos="162,0" size="6x168" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="3,0" size="3x149" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="0,168" size="168x0" type="*widget.Shadow">
							<linearGradient endColor="shadow" pos="0,-8" size="168x8"/>
						</widget>
						<widget pos="0,162" size="168x6" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="0,3" size="55x3" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="168,0" size="0x168" type="*widget.Shadow">
							<linearGradient angle="270" endColor="shadow" pos="-8,0" size="8x168"/>
						</widget>
					</widget>
					<rectangle fillColor="primary" pos="4,0" size="101x4"/>
					<rectangle fillColor="primary" pos="0,42" size="4x37"/>
					<widget pos="105,4" size="1x168" type="*widget.Separator">
						<rectangle fillColor="disabled" size="1x168"/>
					</widget>
					<widget pos="4,41" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,79" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,117" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,155" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())

	table.Select(TableCellID{4, 3})
	assert.Equal(t, 3, table.selectedCell.Col)
	assert.Equal(t, 4, table.selectedCell.Row)
	assert.Equal(t, 3, selectedCol)
	assert.Equal(t, 4, selectedRow)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="180x180">
			<content>
				<widget pos="4,4" size="172x172" type="*widget.Table">
					<widget pos="4,4" size="168x168" type="*widget.ScrollContainer">
						<widget pos="-239,-21" size="509x189" type="*widget.tableCells">
							<widget pos="208,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 2</text>
							</widget>
							<widget pos="310,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 3</text>
							</widget>
							<widget pos="208,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 2</text>
							</widget>
							<widget pos="310,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 3</text>
							</widget>
							<widget pos="208,80" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 2, 2</text>
							</widget>
							<widget pos="310,80" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 2, 3</text>
							</widget>
							<widget pos="208,118" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 3, 2</text>
							</widget>
							<widget pos="310,118" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 3, 3</text>
							</widget>
							<widget pos="208,156" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 4, 2</text>
							</widget>
							<widget pos="310,156" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 4, 3</text>
							</widget>
						</widget>
						<widget pos="162,0" size="6x168" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="3,19" size="3x149" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget size="168x0" type="*widget.Shadow">
							<linearGradient size="168x8" startColor="shadow"/>
						</widget>
						<widget pos="0,162" size="168x6" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="79,3" size="55x3" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget size="0x168" type="*widget.Shadow">
							<linearGradient angle="270" size="8x168" startColor="shadow"/>
						</widget>
						<widget pos="168,0" size="0x168" type="*widget.Shadow">
							<linearGradient angle="270" endColor="shadow" pos="-8,0" size="8x168"/>
						</widget>
					</widget>
					<rectangle fillColor="primary" pos="71,0" size="101x4"/>
					<rectangle fillColor="primary" pos="0,135" size="4x37"/>
					<widget pos="70,4" size="1x168" type="*widget.Separator">
						<rectangle fillColor="disabled" size="1x168"/>
					</widget>
					<widget pos="4,20" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,58" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,96" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,134" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
}

func TestTable_SetColumnWidth(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, obj fyne.CanvasObject) {
			if id.Col == 0 {
				obj.(*Label).Text = "p"
			} else {
				obj.(*Label).Text = "placeholder"
			}
		})
	table.SetColumnWidth(0, 16)
	table.Resize(fyne.NewSize(120, 120))
	table.Select(TableCellID{1, 0})

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh()
	assert.Equal(t, 10, len(cellRenderer.Objects()))
	assert.Equal(t, 16, cellRenderer.(*tableCellsRenderer).Objects()[0].Size().Width)

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(120+(2*theme.Padding()), 120+(2*theme.Padding())))
	test.AssertImageMatches(t, "table/col_size.png", w.Canvas().Capture())
}

func TestTable_ShowVisible(t *testing.T) {
	table := NewTable(
		func() (int, int) { return 50, 50 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(TableCellID, fyne.CanvasObject) {})
	table.Resize(fyne.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh()
	assert.Equal(t, 10, len(cellRenderer.Objects()))
}
