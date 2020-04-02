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

func TestBindIconResource(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	icon := widget.NewIcon(theme.WarningIcon())
	data := &binding.ResourceBinding{}
	binding.BindIconResource(icon, data)
	data.Set(theme.InfoIcon())
	time.Sleep(time.Second)
	assert.Equal(t, theme.InfoIcon(), icon.Resource)
}
