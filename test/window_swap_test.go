package test

import (
	"testing"

	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func Test_WindowSwap(t *testing.T) {
	test := NewApp()
	window := test.NewWindow("Window Swap Test")
	zort := widget.NewLabel("Zort")
	narf := widget.NewLabel("Narf")
	window.SetContent(zort)
	window.SetContent(narf)
	window.SetContent(zort)
	assert.False(t, zort.Hidden)
	window.ShowAndRun()
}
