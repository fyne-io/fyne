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

func TestBindSelectChanged_Binding(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	selected := ""
	selec := widget.NewSelect([]string{"a", "b", "c"}, func(value string) {
		selected = value
	})
	data := &binding.StringBinding{}
	binding.BindSelectChanged(selec, data)
	data.Set("b")
	time.Sleep(time.Second)
	assert.Equal(t, "b", selected)
}

func TestBindSelectChanged_Event(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	selec := widget.NewSelect([]string{"a", "b", "c"}, nil)
	data := &binding.StringBinding{}
	binding.BindSelectChanged(selec, data)
	selected := ""
	data.AddStringListener(func(s string) {
		selected = s
	})

	test.Tap(selec)

	// Get Popup
	canvas := fyne.CurrentApp().Driver().CanvasForObject(selec)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	popup := canvas.Overlays().Top().(*widget.PopUp)
	box := popup.Content.(*widget.Box)
	test.Tap(box.Children[1].(fyne.Tappable))

	time.Sleep(time.Second)
	assert.Equal(t, "b", selected)
}
