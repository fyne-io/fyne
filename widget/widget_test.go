package widget

import "testing"
import "time"

import "github.com/stretchr/testify/assert"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/test"

type myWidget struct {
	baseWidget

	applied chan bool
}

func (m *myWidget) ApplyTheme() {
	m.applied <- true
	close(m.applied)
}

func TestApplyThemeCalled(t *testing.T) {
	widget := &myWidget{applied: make(chan bool)}

	window := test.NewTestWindow(widget)
	fyne.GetSettings().SetTheme("light")

	func() {
		select {
		case <-widget.applied:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for theme apply")
		}
	}()

	window.Close()
}

func TestAppplyThemeCalledChild(t *testing.T) {
	child := &myWidget{applied: make(chan bool)}
	parent := NewList(child)

	window := test.NewTestWindow(parent)
	fyne.GetSettings().SetTheme("light")

	func() {
		select {
		case <-child.applied:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for child theme apply")
		}
	}()

	window.Close()
}
