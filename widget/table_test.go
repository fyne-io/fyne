package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne"
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
	assert.Equal(t, 9, len(cellRenderer.Objects()))
	assert.Equal(t, "Cell 0, 0", cellRenderer.Objects()[0].(*Label).Text)
	objRef := cellRenderer.Objects()[0].(*Label)

	test.Scroll(c, fyne.NewPos(10, 10), -150, -150)
	assert.Equal(t, 0, renderer.scroll.Offset.Y) // we didn't scroll as data shorter
	assert.Equal(t, 150, renderer.scroll.Offset.X)
	assert.Equal(t, 9, len(cellRenderer.Objects()))
	assert.Equal(t, "Cell 0, 1", cellRenderer.Objects()[0].(*Label).Text)
	assert.NotEqual(t, objRef, cellRenderer.Objects()[0].(*Label)) // we want to re-use visible cells without rewriting them
}

func TestTable_ChangeTheme(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 3, 5 },
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
	test.AssertImageMatches(t, "table/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		test.WidgetRenderer(table).Refresh()
		test.AssertImageMatches(t, "table/theme_changed.png", w.Canvas().Capture())
	})
}

func TestTable_Unselect(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

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

func TestTable_Hovered(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 2, 2 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c fyne.CanvasObject) {
			c.(*Label).SetText(fmt.Sprintf("Cell %d, %d", id.Row, id.Col))
		})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	test.MoveMouse(w.Canvas(), fyne.NewPos(35, 50))
	test.MoveMouse(w.Canvas(), fyne.NewPos(35, 100))

	assert.Nil(t, table.hoveredCell)

	test.AssertImageMatches(t, "table/hovered_out.png", w.Canvas().Capture())

	table.Length = func() (int, int) { return 3, 5 }
	table.Refresh()

	w.SetContent(table)
	w.Resize(fyne.NewSize(180, 180))
	test.MoveMouse(w.Canvas(), fyne.NewPos(35, 50))

	assert.Equal(t, 0, table.hoveredCell.Col)
	assert.Equal(t, 1, table.hoveredCell.Row)

	test.AssertImageMatches(t, "table/hovered.png", w.Canvas().Capture())
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
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

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

	test.AssertImageMatches(t, "table/selected.png", w.Canvas().Capture())
}

func TestTable_Select(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

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
	test.AssertImageMatches(t, "table/selected.png", w.Canvas().Capture())

	table.Select(TableCellID{4, 3})
	assert.Equal(t, 3, table.selectedCell.Col)
	assert.Equal(t, 4, table.selectedCell.Row)
	assert.Equal(t, 3, selectedCol)
	assert.Equal(t, 4, selectedRow)
	test.AssertImageMatches(t, "table/selected_scrolled.png", w.Canvas().Capture())
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
	assert.Equal(t, 15, len(cellRenderer.Objects()))
}
