package binding_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestBindRadioChanged_Binding(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	selected := ""
	radio := widget.NewRadio([]string{"a", "b", "c"}, func(value string) {
		selected = value
	})
	data := &binding.StringBinding{}
	binding.BindRadioChanged(radio, data)
	data.Set("b")
	time.Sleep(time.Second)
	assert.Equal(t, "b", selected)
}

func TestBindRadioChanged_Event(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	radio := widget.NewRadio([]string{"a", "b", "c"}, nil)
	data := &binding.StringBinding{}
	binding.BindRadioChanged(radio, data)
	selected := ""
	data.AddStringListener(func(s string) {
		selected = s
	})
	min := radio.MinSize()
	middleX := min.Width / 2
	middleY := min.Height / 2
	test.TapAt(radio, fyne.NewPos(middleX, middleY))
	time.Sleep(time.Second)
	assert.Equal(t, "b", selected)
}
