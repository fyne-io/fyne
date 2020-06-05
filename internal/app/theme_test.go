package app_test

import (
	"testing"

	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/test"
)

func TestApplySettings_BeforeContentSet(t *testing.T) {
	a := test.NewApp()
	_ = a.NewWindow("NoContent")

	app.ApplySettings(a.Settings(), a)
}
