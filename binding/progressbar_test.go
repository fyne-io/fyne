package binding_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"fyne.io/fyne/binding"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestBindProgressBarValue(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	progressBar := widget.NewProgressBar()
	data := &binding.Float64Binding{}
	binding.BindProgressBarValue(progressBar, data)
	data.Set(0.75)
	time.Sleep(time.Second)
	assert.Equal(t, 0.75, progressBar.Value)
}
