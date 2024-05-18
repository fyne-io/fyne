package dialog

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestEntryDialog_Confirm(t *testing.T) {
	value := ""
	ed := NewEntryDialog("Test", "message", func(v string) {
		value = v
	}, test.NewTempWindow(t, nil))
	ed.Show()
	test.Type(ed.entry, "123")
	test.Tap(ed.confirm)

	assert.Equal(t, value, "123", "Control form should be confirmed with no validation")
}

func TestEntryDialog_Dismiss(t *testing.T) {
	value := "123"
	ed := NewEntryDialog("Test", "message", func(v string) {
		value = v
	}, test.NewTempWindow(t, nil))
	ed.Show()
	test.Type(ed.entry, "XYZ")
	test.Tap(ed.cancel)

	assert.Equal(t, value, "123", "Control form should not change value on dismiss")
}
