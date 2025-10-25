package widget_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func BenchmarkProgressbarCreateRenderer(b *testing.B) {
	widget := &widget.ProgressBar{}
	b.ReportAllocs()

	for b.Loop() {
		widget.CreateRenderer()
	}
}

func BenchmarkProgressBarLayout(b *testing.B) {
	b.ReportAllocs() // We should see zero allocations.

	bar := &widget.ProgressBar{}
	renderer := bar.CreateRenderer()

	for b.Loop() {
		renderer.Layout(fyne.NewSize(100, 100))
	}
}

func TestNewProgressBarWithData(t *testing.T) {
	val := binding.NewFloat()
	val.Set(0.4)

	label := widget.NewProgressBarWithData(val)
	waitForBinding()
	assert.Equal(t, 0.4, label.Value)
}

func TestProgressBar_Binding(t *testing.T) {
	bar := widget.NewProgressBar()
	assert.Equal(t, 0.0, bar.Value)

	val := binding.NewFloat()
	val.Set(0.1)
	bar.Bind(val)
	waitForBinding()
	assert.Equal(t, 0.1, bar.Value)

	val.Set(0.4)
	waitForBinding()
	assert.Equal(t, 0.4, bar.Value)

	bar.Unbind()
	waitForBinding()
	assert.Equal(t, 0.4, bar.Value)
}

func TestProgressBar_SetValue(t *testing.T) {
	bar := widget.NewProgressBar()

	assert.Equal(t, 0.0, bar.Min)
	assert.Equal(t, 1.0, bar.Max)
	assert.Equal(t, 0.0, bar.Value)

	bar.SetValue(.5)
	assert.Equal(t, .5, bar.Value)
}

func TestProgressBar_TextFormatter(t *testing.T) {
	bar := widget.NewProgressBar()
	formatted := false

	bar.SetValue(0.2)
	assert.False(t, formatted)

	formatter := func() string {
		formatted = true
		return fmt.Sprintf("%.2f out of %.2f", bar.Value, bar.Max)
	}
	bar.TextFormatter = formatter

	bar.SetValue(0.4)

	assert.True(t, formatted)
}

func TestProgressRenderer_Layout(t *testing.T) {
	bar, c := barOnCanvas()
	test.AssertRendersToMarkup(t, "progressbar/empty.xml", c)

	bar.SetValue(.5)
	test.AssertRendersToMarkup(t, "progressbar/half.xml", c)

	bar.SetValue(1)
	test.AssertRendersToMarkup(t, "progressbar/full.xml", c)
}

func TestProgressRenderer_Layout_Overflow(t *testing.T) {
	bar, c := barOnCanvas()
	bar.SetValue(1)
	test.AssertRendersToMarkup(t, "progressbar/full.xml", c)

	bar.SetValue(1.2)
	test.AssertRendersToMarkup(t, "progressbar/full.xml", c)
}

func TestProgressRenderer_ApplyTheme(t *testing.T) {
	test.WithTestTheme(t, func() {
		bar, c := barOnCanvas()
		bar.SetValue(.2)
		test.AssertRendersToMarkup(t, "progressbar/themed.xml", c)
	})
}

func TestProgressBar_FromStruct(t *testing.T) {
	bar := &widget.ProgressBar{Max: 10, Value: 5}
	w := test.NewTempWindow(t, bar)
	w.Resize(bar.MinSize().AddWidthHeight(100, 0))

	test.AssertRendersToImage(t, "progressbar/initial.png", w.Canvas())
}

func barOnCanvas() (*widget.ProgressBar, fyne.Canvas) {
	bar := widget.NewProgressBar()
	window := test.NewWindow(container.NewVBox(bar))
	window.SetPadded(false)
	window.Resize(fyne.NewSize(100, 100))
	return bar, window.Canvas()
}
