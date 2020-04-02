package binding_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"fyne.io/fyne/binding"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestBindCheckChanged_Binding(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	checked := false
	check := widget.NewCheck("check", func(value bool) {
		checked = value
	})
	data := &binding.BoolBinding{}
	binding.BindCheckChanged(check, data)
	data.Set(true)
	time.Sleep(time.Second)
	assert.True(t, checked)
}

func TestBindCheckChangedEvent(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	check := widget.NewCheck("check", nil)
	data := &binding.BoolBinding{}
	binding.BindCheckChanged(check, data)
	checked := false
	data.AddBoolListener(func(b bool) {
		checked = b
	})
	test.Tap(check)
	time.Sleep(time.Second)
	assert.True(t, checked)
}

func TestBindCheckText(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	check := widget.NewCheck("check", nil)
	data := &binding.StringBinding{}
	binding.BindCheckText(check, data)
	data.Set("foobar")
	time.Sleep(time.Second)
	assert.Equal(t, "foobar", check.Text)
}
