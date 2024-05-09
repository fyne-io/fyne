package widget

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/software"
	"fyne.io/fyne/v2/test"
)

func TestActivity_Start(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, test.NewTheme())

	a := NewActivity()
	w := software.NewWindow(a)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50))

	img1 := w.Canvas().Capture()
	a.Start()
	time.Sleep(time.Millisecond * 50)
	img2 := w.Canvas().Capture()
	a.Stop()

	assert.NotEqual(t, img1, img2)
}

func TestActivity_Stop(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, test.NewTheme())

	a := NewActivity()
	w := software.NewWindow(a)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50))

	a.Start()
	time.Sleep(time.Millisecond * 50)
	a.Stop()

	img1 := w.Canvas().Capture()
	time.Sleep(time.Millisecond * 50)
	img2 := w.Canvas().Capture()
	assert.Equal(t, img1, img2)
}
