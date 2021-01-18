package app_test

import (
	"testing"

	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/test"
)

func TestApplySettings_BeforeContentSet(t *testing.T) {
	a := test.NewApp()
	_ = a.NewWindow("NoContent")

	app.ApplySettings(a.Settings(), a)
}
