package widget_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestEntry_ActionItemBorder(t *testing.T) {
	test.NewTempApp(t)
	entry := widget.NewEntry()
	tapped := false
	entry.ActionItem = widget.NewButtonWithIcon("", theme.QuestionIcon(), func() {
		tapped = true
	})
	window := test.NewTempWindow(t, entry)
	window.SetPadded(false)
	window.Resize(fyne.NewSize(240, entry.MinSize().Height))
	window.Canvas().Focus(entry)

	test.AssertImageMatches(t, "entry/action_item_border.png", window.Canvas().Capture())
	test.TapCanvas(window.Canvas(), entry.ActionItem.Position().Add(fyne.NewPos(10, 10)))
	assert.True(t, tapped)
}
