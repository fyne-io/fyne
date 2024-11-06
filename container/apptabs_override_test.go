//go:build !mobile

package container_test

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestAppTabs_OverrideMobile(t *testing.T) {
	test.NewTempApp(t)

	item1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	item2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Text 2")}
	item3 := &container.TabItem{Text: "Test3", Content: widget.NewLabel("Text 3")}
	tabs := container.NewAppTabs(item1, item2, item3)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	min := tabs.MinSize()
	w.Resize(min)

	test.AssertRendersToMarkup(t, "apptabs/desktop/tab_location_top.xml", c)

	override := container.NewThemeOverride(tabs, test.Theme())
	override.SetDeviceIsMobile(true)
	w.Resize(min.AddWidthHeight(-4, -0))

	test.AssertRendersToMarkup(t, "apptabs/mobile/tab_location_top.xml", c)
}
