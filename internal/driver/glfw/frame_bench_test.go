package glfw

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/widget"
)

// buildTree makes a tree roughly the shape of a real dashboard: tabs, a form of
// labels/entries and a grid of buttons.
func buildTree() fyne.CanvasObject {
	rows := make([]fyne.CanvasObject, 0, 60)
	for i := 0; i < 60; i++ {
		rows = append(rows, container.NewHBox(
			widget.NewLabel(fmt.Sprintf("symbol %d", i)),
			widget.NewLabel(fmt.Sprintf("%d.%d", i, i)),
			widget.NewButton("go", func() {}),
		))
	}
	grid := container.NewGridWithColumns(3, rows...)
	return container.NewAppTabs(
		container.NewTabItem("one", grid),
		container.NewTabItem("two", widget.NewLabel("second")),
	)
}

func BenchmarkFrameEnsureMinSize(b *testing.B) {
	c := newCanvas()
	c.SetContent(buildTree())
	c.Resize(fyne.NewSize(1200, 900))
	c.EnsureMinSize()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.EnsureMinSize()
	}
}

func BenchmarkFrameWalkTrees(b *testing.B) {
	c := newCanvas()
	c.SetContent(buildTree())
	c.Resize(fyne.NewSize(1200, 900))
	c.EnsureMinSize()

	count := 0
	visit := func(*common.RenderCacheNode, fyne.Position) { count++ }

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.WalkTrees(visit, visit)
	}
	b.StopTimer()
	b.Logf("visited %d nodes per pass", count/b.N/2)
}
