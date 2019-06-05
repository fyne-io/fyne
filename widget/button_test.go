package widget

import (
	"fmt"
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
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

func TestButton_DisabledColor(t *testing.T) {
	button := NewButton("Test", nil)
	bg := Renderer(button).BackgroundColor()
	button.Style = DefaultButton
	assert.Equal(t, bg, theme.ButtonColor())

	button.Disable()
	bg = Renderer(button).BackgroundColor()
	assert.Equal(t, bg, theme.DisabledButtonColor())
}

func TestButton_DisabledIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := Renderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	button.Enable()
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())
}

func TestButton_DisabledIconChangeUsingSetIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := Renderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	// assert we are using the disabled original icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	// re-enable, then change the icon
	button.Enable()
	button.SetIcon(theme.SearchIcon())
	assert.Equal(t, render.icon.Resource.Name(), theme.SearchIcon().Name())

	// assert we are using the disabled new icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.SearchIcon().Name()))

}

func TestButton_DisabledIconChangedDirectly(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := Renderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	// assert we are using the disabled original icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	// re-enable, then change the icon
	button.Enable()
	button.Icon = theme.SearchIcon()
	render.Refresh()
	assert.Equal(t, render.icon.Resource.Name(), theme.SearchIcon().Name())

	// assert we are using the disabled new icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.SearchIcon().Name()))

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

func TestButtonRenderer_Layout_Stretch(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	button.Resize(button.MinSize().Add(fyne.NewSize(100, 100)))
	render := Renderer(button).(*buttonRenderer)

	iconYOffset, labelYOffset := 0, 0
	textHeight := render.label.MinSize().Height
	if theme.IconInlineSize() > textHeight {
		labelYOffset = (theme.IconInlineSize() - textHeight) / 2
	} else {
		iconYOffset = (textHeight - theme.IconInlineSize()) / 2
	}
	assert.Equal(t, fyne.NewPos(50+theme.Padding()*2, 50+theme.Padding()+iconYOffset), render.icon.Position(), "icon position")
	assert.Equal(t, fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()), render.icon.Size(), "icon size")
	assert.Equal(t, fyne.NewPos(50+theme.Padding()*3+theme.IconInlineSize(), 50+theme.Padding()+labelYOffset), render.label.Position(), "label position")
	assert.Equal(t, render.label.MinSize(), render.label.Size(), "label size")
}

func TestButtonRenderer_Layout_NoText(t *testing.T) {
	button := NewButtonWithIcon("", theme.CancelIcon(), nil)
	render := Renderer(button).(*buttonRenderer)

	width := button.MinSize().Width + 20
	height := button.MinSize().Height + 20
	button.Resize(fyne.NewSize(width, height))

	iconOffset := 0
	if theme.IconInlineSize() < theme.TextSize() {
		iconOffset = (canvas.NewText("", color.White).MinSize().Height - theme.IconInlineSize()) / 2
	}

	assert.Equal(t, theme.Padding()+10, render.icon.Position().X)
	assert.Equal(t, theme.Padding()+10+iconOffset, render.icon.Position().Y)
}

func TestButton_Disable(t *testing.T) {
	tapped := make(chan bool)
	button := NewButton("Test", func() {
		tapped <- true
	})

	button.Disable()
	go test.Tap(button)
	func() {
		select {
		case <-tapped:
			assert.Fail(t, "Button should have been disabled")
		case <-time.After(1 * time.Second):
		}
	}()
}

func TestButton_Enable(t *testing.T) {
	tapped := make(chan bool)
	button := NewButton("Test", func() {
		tapped <- true
	})

	button.Disable()
	go test.Tap(button)
	func() {
		select {
		case <-tapped:
			assert.Fail(t, "Button should have been disabled")
		case <-time.After(1 * time.Second):
		}
	}()

	button.Enable()
	go test.Tap(button)
	func() {
		select {
		case <-tapped:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Button should have been re-enabled")
		}
	}()
}

func TestButton_Disabled(t *testing.T) {
	button := NewButton("Test", func() {})
	assert.False(t, button.Disabled())
	button.Disable()
	assert.True(t, button.Disabled())
	button.Enable()
	assert.False(t, button.Disabled())
}
