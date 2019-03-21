package widget

import (
	"testing"
	"time"

	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestButton_MinSize(t *testing.T) {
	button := NewButton("Hi", nil)
	min := button.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestButton_SetText(t *testing.T) {
	button := NewButton("Hi", nil)
	min1 := button.MinSize()

	button.SetText("Longer")
	min2 := button.MinSize()

	assert.True(t, min2.Width > min1.Width)
	assert.Equal(t, min2.Height, min1.Height)
}

func TestButton_MinSize_Icon(t *testing.T) {
	button := NewButton("Hi", nil)
	min1 := button.MinSize()

	button.SetIcon(theme.CancelIcon())
	min2 := button.MinSize()

	assert.True(t, min2.Width > min1.Width)
	assert.Equal(t, min2.Height, min1.Height)
}

func TestButton_Style(t *testing.T) {
	button := NewButton("Test", nil)
	bg := Renderer(button).BackgroundColor()

	button.Style = PrimaryButton
	assert.NotEqual(t, bg, Renderer(button).BackgroundColor())
}

func TestButton_Tapped(t *testing.T) {
	tapped := make(chan bool)
	button := NewButton("Hi", func() {
		tapped <- true
	})

	go test.Tap(button)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for button tap")
		}
	}()
}

func TestButtonRenderer_Layout(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := Renderer(button).(*buttonRenderer)

	assert.True(t, render.icon.Position().X < render.label.Position().X)
	assert.Equal(t, theme.Padding()*2, render.icon.Position().X)
	assert.Equal(t, theme.Padding()*2, render.MinSize().Width-render.label.Position().X-render.label.Size().Width)
}
