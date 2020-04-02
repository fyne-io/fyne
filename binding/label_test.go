package binding_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"fyne.io/fyne/binding"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestBindLabelText(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	label := widget.NewLabel("label")
	data := &binding.StringBinding{}
	binding.BindLabelText(label, data)
	data.Set("foobar")
	time.Sleep(time.Second)
	assert.Equal(t, "foobar", label.Text)
}
