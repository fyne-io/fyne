package widget_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

var shadowLevel = widget.ElevationLevel(5)

func TestShadow_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowAround, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50))

	s.Resize(fyne.NewSize(30, 30))
	s.Move(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "shadow/theme_default.png", w.Canvas().Capture())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "shadow/theme_ugly.png", w.Canvas().Capture())
}

func TestShadow_AroundShadow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowAround, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50))

	s.Resize(fyne.NewSize(30, 30))
	s.Move(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "shadow/around.png", w.Canvas().Capture())
}

func TestShadow_Transparency(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowAround, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50))

	s.Resize(fyne.NewSize(40, 20))
	s.Move(fyne.NewPos(5, 15))
	s2 := widget.NewShadow(widget.ShadowAround, shadowLevel)
	w.Canvas().Overlays().Add(s2)
	s2.Resize(fyne.NewSize(20, 40))
	s2.Move(fyne.NewPos(15, 5))
	test.AssertImageMatches(t, "shadow/transparency.png", w.Canvas().Capture())
}

func TestShadow_BottomShadow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowBottom, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50))

	s.Resize(fyne.NewSize(30, 30))
	s.Move(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "shadow/bottom.png", w.Canvas().Capture())
}

func TestShadow_MinSize(t *testing.T) {
	assert.Equal(t, fyne.NewSize(0, 0), widget.NewShadow(widget.ShadowAround, 1).MinSize())
}

func TestShadow_TopShadow(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	s := widget.NewShadow(widget.ShadowTop, shadowLevel)
	w := test.NewWindow(s)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50))

	s.Resize(fyne.NewSize(30, 30))
	s.Move(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "shadow/top.png", w.Canvas().Capture())
}
