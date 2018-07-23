package widget

import "testing"
import "time"

import "github.com/stretchr/testify/assert"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/test"

var applied = make(chan bool)

type myWidget struct {
	baseWidget
}

func (m *myWidget) ApplyTheme() {
	applied <- true
}

func TestApplyThemeCalled(t *testing.T) {
	widget := &myWidget{}

	app := test.NewApp()
	window := app.NewWindow("Testing")
	window.SetContent(widget)

	fyne.GetSettings().SetTheme("light")

	func() {
		select {
		case <-applied:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for theme apply")
		}
	}()
}
