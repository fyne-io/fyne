package widget

import (
	"testing"
	"time"

	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestButtonSize(t *testing.T) {
	button := NewButton("Hi", nil)
	min := button.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestButtonNotify(t *testing.T) {
	tapped := make(chan bool)
	button := NewButton("Hi", func() {
		tapped <- true
	})

	go button.OnMouseDown(nil)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for button tap")
		}
	}()
}
