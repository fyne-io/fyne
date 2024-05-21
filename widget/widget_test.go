package widget

import (
	"image/color"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	internalTest "fyne.io/fyne/v2/internal/test"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
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
	fyne.CurrentApp().Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))

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
	fyne.CurrentApp().Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))
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

	c := &fyne.Container{Layout: layout.NewStackLayout(), Objects: []fyne.CanvasObject{
		newTestWidget(canvas.NewRectangle(color.Gray{Y: 0x79})),
		newTestWidget(canvas.NewText("Hi", color.Black))}}

	window := test.NewWindow(c)
	defer window.Close()

	test.AssertImageMatches(t, "simple_renderer.png", window.Canvas().Capture())
}

func TestMinSizeCache(t *testing.T) {
	label := NewLabel("a")

	wid := newTestWidget(label)
	assert.Equal(t, 0, wid.minSizeCalls)

	minSize := wid.MinSize()
	assert.NotEqual(t, 0, minSize)
	assert.Equal(t, 1, wid.minSizeCalls)

	wid.MinSize()
	assert.Equal(t, 1, wid.minSizeCalls)
}

type testWidget struct {
	BaseWidget
	obj          fyne.CanvasObject
	minSizeCalls int
}

func newTestWidget(o fyne.CanvasObject) *testWidget {
	t := &testWidget{obj: o}
	t.ExtendBaseWidget(t)
	return t
}

func (t *testWidget) CreateRenderer() fyne.WidgetRenderer {
	return &testRenderer{
		WidgetRenderer: NewSimpleRenderer(t.obj),
		widget:         t,
	}
}

func waitForBinding() {
	time.Sleep(time.Millisecond * 100) // data resolves on background thread
}

type testRenderer struct {
	widget *testWidget
	fyne.WidgetRenderer
}

func (r *testRenderer) MinSize() fyne.Size {
	r.widget.minSizeCalls++
	return r.WidgetRenderer.MinSize()
}
