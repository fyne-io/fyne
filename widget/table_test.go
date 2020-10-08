package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestTable_ChangeTheme(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		}, func(row, col int, c fyne.CanvasObject) {
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

func TestTable_Selected(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	table := NewTable(
		func() (int, int) { return 3, 5 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		}, func(row, col int, c fyne.CanvasObject) {
			text := fmt.Sprintf("Cell %d, %d", row, col)
			c.(*Label).SetText(text)
		})
	assert.Equal(t, -1, table.SelectedColumn)
	assert.Equal(t, -1, table.SelectedRow)

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	selectedCol, selectedRow := 0, 0
	table.OnCellSelected = func(row, col int) {
		selectedCol = col
		selectedRow = row
	}
	test.TapCanvas(w.Canvas(), fyne.NewPos(35, 50))
	assert.Equal(t, 0, table.SelectedColumn)
	assert.Equal(t, 1, table.SelectedRow)
	assert.Equal(t, 0, selectedCol)
	assert.Equal(t, 1, selectedRow)

	test.AssertImageMatches(t, "table/selected.png", w.Canvas().Capture())
}

func TestTable_ShowVisible(t *testing.T) {
	table := NewTable(
		func() (int, int) { return 50, 50 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		}, func(int, int, fyne.CanvasObject) {})
	table.Resize(fyne.NewSize(120, 120))

	renderer := test.WidgetRenderer(table).(*tableRenderer)
	cellRenderer := test.WidgetRenderer(renderer.scroll.Content.(*tableCells))
	cellRenderer.Refresh()
	assert.Equal(t, 15, len(cellRenderer.Objects()))
}
