//go:build mobile && flakey

package container_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestAppTabs_ApplyTheme(t *testing.T) {
	test.NewTempApp(t)

	w := test.NewWindow(
		container.NewAppTabs(&container.TabItem{Text: "Test", Content: widget.NewLabel("Text")}),
	)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToImage(t, "apptabs/mobile/theme_default.png", c)

	test.ApplyTheme(t, test.NewTheme())
	test.AssertRendersToImage(t, "apptabs/mobile/theme_ugly.png", c)
}

func TestAppTabs_ChangeItemContent(t *testing.T) {
	test.NewTempApp(t)

	item1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := container.NewAppTabs(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "apptabs/mobile/change_content_initial.xml", c)

	item1.Content = widget.NewLabel("Text3")
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "apptabs/mobile/change_content_change_visible.xml", c)

	item2.Content = widget.NewLabel("Text4")
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "apptabs/mobile/change_content_change_hidden.xml", c)
}
