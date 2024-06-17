package widget

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestActivity_Start(t *testing.T) {
	test.NewTempApp(t)
	test.ApplyTheme(t, test.NewTheme())

	a := NewActivity()
	w := test.NewWindow(a)
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
	test.NewTempApp(t)
	test.ApplyTheme(t, test.NewTheme())

	a := NewActivity()
	w := test.NewWindow(a)
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
