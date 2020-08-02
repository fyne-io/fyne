package dialog

import (
	"testing"

	"fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestEntryDialogConfirm(t *testing.T) {
	ch := make(chan string)
	i := NewEntryDialog("test1", "test2", func(s string) { ch <- s }, test.NewWindow(nil))
	i.Show()

	assert.False(t, i.win.Hidden)
	test.Type(i.entry, "test3")
	go test.Tap(i.confirmButton)
	assert.EqualValues(t, "test3", <-ch)
	assert.EqualValues(t, "test3", i.entry.Text)
}
