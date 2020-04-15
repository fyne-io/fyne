package widget

import (
	"fmt"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/driver/desktop"
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

func TestButton_Cursor(t *testing.T) {
	button := NewButton("Test", nil)
	assert.Equal(t, desktop.DefaultCursor, button.Cursor())
}

func TestButton_Style(t *testing.T) {
	button := NewButton("Test", nil)
	bg := test.WidgetRenderer(button).BackgroundColor()

	button.Style = PrimaryButton
	assert.NotEqual(t, bg, test.WidgetRenderer(button).BackgroundColor())
}

func TestButton_DisabledColor(t *testing.T) {
	button := NewButton("Test", nil)
	bg := test.WidgetRenderer(button).BackgroundColor()
	button.Style = DefaultButton
	assert.Equal(t, bg, theme.ButtonColor())

	button.Disable()
	bg = test.WidgetRenderer(button).BackgroundColor()
	assert.Equal(t, bg, theme.DisabledButtonColor())
}

func TestButton_DisabledIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	button.Enable()
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())
}

func TestButton_DisabledIconChangeUsingSetIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
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
	render := test.WidgetRenderer(button).(*buttonRenderer)
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

func TestButton_BindText(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	button := NewButton("button", nil)
	data := binding.NewString("foo")
	button.BindText(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, "foo", button.Text)

	data.Set("foobar")
	timedWait(t, done)
	assert.Equal(t, "foobar", button.Text)
}

func TestButton_BindIcon(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	done := make(chan bool)
	button := NewButtonWithIcon("button", theme.WarningIcon(), nil)
	data := binding.NewResource(theme.QuestionIcon())
	button.BindIcon(data)
	data.AddListener(binding.NewNotifyFunction(func(binding.Binding) {
		done <- true
	}))
	timedWait(t, done)
	assert.Equal(t, theme.QuestionIcon(), button.Icon)

	data.Set(theme.InfoIcon())
	timedWait(t, done)
	assert.Equal(t, theme.InfoIcon(), button.Icon)
}

func TestButtonRenderer_Layout(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	render.Layout(render.MinSize())

	assert.True(t, render.icon.Position().X < render.label.Position().X)
	assert.Equal(t, theme.Padding()*2, render.icon.Position().X)
	assert.Equal(t, theme.Padding()*2, render.MinSize().Width-render.label.Position().X-render.label.Size().Width)
}

func TestButtonRenderer_Layout_Stretch(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	button.Resize(button.MinSize().Add(fyne.NewSize(100, 100)))
	render := test.WidgetRenderer(button).(*buttonRenderer)

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
	render := test.WidgetRenderer(button).(*buttonRenderer)

	button.Resize(fyne.NewSize(100, 100))

	assert.Equal(t, 50-theme.IconInlineSize()/2, render.icon.Position().X)
	assert.Equal(t, 50-theme.IconInlineSize()/2, render.icon.Position().Y)
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

func TestButton_Shadow(t *testing.T) {
	{
		button := NewButton("Test", func() {})
		shadowFound := false
		for _, o := range test.LaidOutObjects(button) {
			if s, ok := o.(*shadow); ok {
				shadowFound = true
				assert.Equal(t, elevationLevel(2), s.level)
			}
		}
		if !shadowFound {
			assert.Fail(t, "button should cast a shadow")
		}
	}
	{
		button := NewButton("Test", func() {})
		button.HideShadow = true
		for _, o := range test.LaidOutObjects(button) {
			if _, ok := o.(*shadow); ok {
				assert.Fail(t, "button with HideShadow == true should not create a shadow")
			}
		}
	}
}

func TestButtonRenderer_ApplyTheme(t *testing.T) {
	button := &Button{}
	render := test.WidgetRenderer(button).(*buttonRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	withTestTheme(func() {
		render.applyTheme()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}
