package widget_test

import (
	"fmt"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestList_ThemeChange(t *testing.T) {
	list, w, _ := setupList(t)

	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		time.Sleep(100 * time.Millisecond)
		list.Refresh()
		test.AssertImageMatches(t, "list/list_theme_changed.png", w.Canvas().Capture())
	})
}

func TestList_ThemeOverride(t *testing.T) {
	list, w, _ := setupList(t)

	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "list/list_theme_changed.png", w.Canvas().Capture())

	normal := test.Theme()
	bg := canvas.NewRectangle(normal.Color(theme.ColorNameBackground, theme.VariantDark))
	w.SetContent(container.NewStack(bg, container.NewThemeOverride(list, normal)))
	w.Resize(fyne.NewSize(200, 200))
	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())
}

func TestList_Resize(t *testing.T) {
	l, _, rows := setupList(t)

	var refreshCounts []int
	var resizeCounts []int
	for _, row := range rows {
		refreshCounts = append(refreshCounts, row.refreshCount)
		resizeCounts = append(resizeCounts, row.resizeCount)
	}

	l.Resize(fyne.NewSize(250, 200)) // changing width only
	var newRefreshCounts []int
	var newResizeCounts []int
	for _, row := range rows {
		newRefreshCounts = append(newRefreshCounts, row.refreshCount)
		newResizeCounts = append(newResizeCounts, row.resizeCount)
	}

	// resizing width only (no new visible rows) calls Resize on the rows, not Refresh
	// some rows may not be visible, so their resize does not need to be called
	assert.Equal(t, refreshCounts, newRefreshCounts)
	atLeastOneGreaterResizeCount := false
	resizeCountAllGreaterOrEqual := true
	for i := range resizeCounts {
		if newResizeCounts[i] < resizeCounts[i] {
			resizeCountAllGreaterOrEqual = false
			break
		}
		if newResizeCounts[i] > resizeCounts[i] {
			atLeastOneGreaterResizeCount = true
		}
	}
	assert.True(t, atLeastOneGreaterResizeCount)
	assert.True(t, resizeCountAllGreaterOrEqual)
}

func setupList(t *testing.T) (*widget.List, fyne.Window, []*resizeRefreshCountingLabel) {
	var rows []*resizeRefreshCountingLabel
	test.NewTempApp(t)
	list := widget.NewList(
		func() int {
			return 25
		},
		func() fyne.CanvasObject {
			row := newResizeRefreshCountingLabel("Test Item 55")
			rows = append(rows, row)
			return row
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*resizeRefreshCountingLabel).SetText(fmt.Sprintf("Test Item %d", id))
		})
	w := test.NewTempWindow(t, list)
	w.SetPadded(false)
	w.Resize(fyne.NewSize(200, 200))

	return list, w, rows
}

type resizeRefreshCountingLabel struct {
	widget.Label
	resizeCount  int
	refreshCount int
}

func newResizeRefreshCountingLabel(text string) *resizeRefreshCountingLabel {
	r := &resizeRefreshCountingLabel{}
	r.Text = text
	r.ExtendBaseWidget(r)
	return r
}

func (r *resizeRefreshCountingLabel) Refresh() {
	r.refreshCount++
	r.Label.Refresh()
}

func (r *resizeRefreshCountingLabel) Resize(s fyne.Size) {
	r.resizeCount++
	r.Label.Resize(s)
}
