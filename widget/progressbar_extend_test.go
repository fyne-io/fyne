package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

type extendedProgressBar struct {
	ProgressBar
}

func newExtendedProgressBar() *extendedProgressBar {
	ret := &extendedProgressBar{}
	ret.ExtendBaseWidget(ret)

	return ret
}

func TestProgressBarRenderer_Extended_Layout(t *testing.T) {
	bar := newExtendedProgressBar()
	bar.Resize(fyne.NewSize(100, 100))
	r := test.WidgetRenderer(bar).(*progressRenderer)

	assert.Equal(t, 0.0, bar.Value)
	assert.Equal(t, 0, r.bar.Size().Width)

	bar.SetValue(.5)
	assert.Equal(t, .5, bar.Value)
	assert.Equal(t, 50, r.bar.Size().Width)
}
