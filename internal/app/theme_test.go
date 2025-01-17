package app_test

import (
	"testing"

	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/test"
)

func TestApplySettings_BeforeContentSet(t *testing.T) {
	a := test.NewApp()
	w := a.NewWindow("NoContent")
	defer w.Close()

	app.ApplySettings(a.Settings(), a)
}
