package binding_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"fyne.io/fyne/binding"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestBindButtonTapped_Binding(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	tapped := false
	button := widget.NewButton("button", func() {
		tapped = true
	})
	data := &binding.BoolBinding{}
	binding.BindButtonTapped(button, data)
	data.Set(true)
	time.Sleep(time.Second)
	assert.True(t, tapped)
}

func TestBindButtonTapped_Event(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	button := widget.NewButton("button", nil)
	data := &binding.BoolBinding{}
	binding.BindButtonTapped(button, data)
	tapped := false
	data.AddBoolListener(func(b bool) {
		tapped = b
	})
	test.Tap(button)
	time.Sleep(time.Second)
	assert.True(t, tapped)
}

func TestBindButtonText(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	button := widget.NewButton("button", nil)
	data := &binding.StringBinding{}
	binding.BindButtonText(button, data)
	data.Set("foobar")
	time.Sleep(time.Second)
	assert.Equal(t, "foobar", button.Text)
}

func TestBindButtonIcon(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	button := widget.NewButtonWithIcon("button", theme.WarningIcon(), nil)
	data := &binding.ResourceBinding{}
	binding.BindButtonIcon(button, data)
	data.Set(theme.InfoIcon())
	time.Sleep(time.Second)
	assert.Equal(t, theme.InfoIcon(), button.Icon)
}
