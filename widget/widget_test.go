package widget

import (
	"testing"
	"time"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/test"
	"github.com/fyne-io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

type myWidget struct {
	baseWidget

	applied chan bool
}

func (m *myWidget) Resize(size fyne.Size) {
	m.resize(size, m)
}

func (m *myWidget) Move(pos fyne.Position) {
	m.move(pos, m)
}

func (m *myWidget) MinSize() fyne.Size {
	return m.minSize(m)
}

func (m *myWidget) Show() {
	m.show(m)
}

func (m *myWidget) Hide() {
	m.hide(m)
}

func (m *myWidget) ApplyTheme() {
	m.applied <- true
}

func (m *myWidget) CreateRenderer() fyne.WidgetRenderer {
	return (&Box{}).CreateRenderer()
}

func TestApplyThemeCalled(t *testing.T) {
	widget := &myWidget{applied: make(chan bool)}

	window := test.NewTestWindow(widget)
	fyne.GlobalSettings().SetTheme(theme.LightTheme())

	func() {
		select {
		case <-widget.applied:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for theme apply")
		}
	}()

	close(widget.applied)
	window.Close()
}

func TestApplyThemeCalledChild(t *testing.T) {
	child := &myWidget{applied: make(chan bool)}
	parent := NewVBox(child)

	window := test.NewTestWindow(parent)
	fyne.GlobalSettings().SetTheme(theme.LightTheme())
	func() {
		select {
		case <-child.applied:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for child theme apply")
		}
	}()

	close(child.applied)
	window.Close()
}
