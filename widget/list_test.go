package widget_test

import (
	"fmt"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestList_ThemeChange(t *testing.T) {
	list, w := setupList(t)

	test.AssertImageMatches(t, "list/list_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		time.Sleep(100 * time.Millisecond)
		list.Refresh()
		test.AssertImageMatches(t, "list/list_theme_changed.png", w.Canvas().Capture())
	})
}

func setupList(t *testing.T) (*widget.List, fyne.Window) {
	test.NewTempApp(t)
	list := widget.NewList(
		func() int {
			return 25
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Test Item 55")
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("Test Item %d", id))
		})
	w := test.NewTempWindow(t, list)
	w.SetPadded(false)
	w.Resize(fyne.NewSize(200, 200))

	return list, w
}
