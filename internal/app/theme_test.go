package app

import (
	"testing"

	"fyne.io/fyne/test"
)

func TestApplySettings_BeforeContentSet(t *testing.T) {
	a := test.NewApp()
	_ = a.NewWindow("NoContent")

	ApplySettings(a.Settings(), a)
}
