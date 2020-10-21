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
		func(row, col int, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", row, col)
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
		func(row, col int, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", row, col)
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

func TestTable_ClearSelection(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(row, col int, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", row, col)
			c.(*Label).SetText(text)
		})
	table.selectedColumn = 1
	table.selectedRow = 1
	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	table.ClearSelection()
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
		func(row, col int, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", row, col)
			c.(*Label).SetText(text)
		})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	test.MoveMouse(w.Canvas(), fyne.NewPos(35, 100))

	assert.Equal(t, -1, table.hoveredColumn)
	assert.Equal(t, -1, table.hoveredRow)

	test.AssertImageMatches(t, "table/nohovered.png", w.Canvas().Capture())

	table = NewTable(
		func() (int, int) { return 3, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(row, col int, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", row, col)
			c.(*Label).SetText(text)
		})

	w.SetContent(table)
	w.Resize(fyne.NewSize(180, 180))
	test.MoveMouse(w.Canvas(), fyne.NewPos(35, 50))

	assert.Equal(t, 0, table.hoveredColumn)
	assert.Equal(t, 1, table.hoveredRow)

	test.AssertImageMatches(t, "table/hovered.png", w.Canvas().Capture())
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
		func(row, col int, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", row, col)
			c.(*Label).SetText(text)
		})
	assert.Equal(t, -1, table.selectedColumn)
	assert.Equal(t, -1, table.selectedRow)

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	selectedCol, selectedRow := 0, 0
	table.OnSelectionChanged = func(row, col int) {
		selectedCol = col
		selectedRow = row
	}
	test.TapCanvas(w.Canvas(), fyne.NewPos(35, 50))
	assert.Equal(t, 0, table.selectedColumn)
	assert.Equal(t, 1, table.selectedRow)
	assert.Equal(t, 0, selectedCol)
	assert.Equal(t, 1, selectedRow)

	test.AssertImageMatches(t, "table/selected.png", w.Canvas().Capture())
}

func TestTable_SetSelection(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 5, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(row, col int, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", row, col)
			c.(*Label).SetText(text)
		})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	selectedCol, selectedRow := 0, 0
	table.OnSelectionChanged = func(row, col int) {
		selectedCol = col
		selectedRow = row
	}
	table.SetSelection(1, 0)
	assert.Equal(t, 0, table.selectedColumn)
	assert.Equal(t, 1, table.selectedRow)
	assert.Equal(t, 0, selectedCol)
	assert.Equal(t, 1, selectedRow)
	test.AssertImageMatches(t, "table/selected.png", w.Canvas().Capture())

	table.SetSelection(3, 3)
	assert.Equal(t, 3, table.selectedColumn)
	assert.Equal(t, 3, table.selectedRow)
	assert.Equal(t, 3, selectedCol)
	assert.Equal(t, 3, selectedRow)
	test.AssertImageMatches(t, "table/selected_scrolled.png", w.Canvas().Capture())
}

func TestTable_ShowVisible(t *testing.T) {
	table := NewTable(
		func() (int, int) { return 50, 50 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(int, int, fyne.CanvasObject) {})
	table.Resize(fyne.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh()
	assert.Equal(t, 15, len(cellRenderer.Objects()))
}
