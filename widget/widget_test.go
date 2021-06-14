package widget

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

type myWidget struct {
	DisableableWidget

	refreshed chan bool
}

func (m *myWidget) Refresh() {
	m.refreshed <- true
}

func (m *myWidget) CreateRenderer() fyne.WidgetRenderer {
	m.ExtendBaseWidget(m)
	return internalWidget.NewSimpleRenderer(&fyne.Container{})
}

func TestApplyThemeCalled(t *testing.T) {
	widget := &myWidget{refreshed: make(chan bool)}

	window := test.NewWindow(widget)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())

	func() {
		select {
		case <-widget.refreshed:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for theme apply")
		}
	}()

	window.Close()
}

func TestApplyThemeCalledChild(t *testing.T) {
	child := &myWidget{refreshed: make(chan bool)}
	parent := &fyne.Container{Layout: layout.NewVBoxLayout(), Objects: []fyne.CanvasObject{child}}

	window := test.NewWindow(parent)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	func() {
		select {
		case <-child.refreshed:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for child theme apply")
		}
	}()

	window.Close()
}

func TestSimpleRenderer(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	c := &fyne.Container{Layout: layout.NewMaxLayout(), Objects: []fyne.CanvasObject{
		newTestWidget(canvas.NewRectangle(color.Gray{Y: 0x79})),
		newTestWidget(canvas.NewText("Hi", color.Black))}}

	window := test.NewWindow(c)
	defer window.Close()

	test.AssertImageMatches(t, "simple_renderer.png", window.Canvas().Capture())
}

type testWidget struct {
	BaseWidget
	obj fyne.CanvasObject
}

func newTestWidget(o fyne.CanvasObject) fyne.Widget {
	t := &testWidget{obj: o}
	t.ExtendBaseWidget(t)
	return t
}

func (t *testWidget) CreateRenderer() fyne.WidgetRenderer {
	return NewSimpleRenderer(t.obj)
}

func waitForBinding() {
	time.Sleep(time.Millisecond * 100) // data resolves on background thread
}
