package dialog

import (
	"testing"

	"fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestDialog_ConfirmDoubleCallback(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", func(_ bool) {
		ch <- 42
	}, test.NewWindow(nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.SetOnClosed(func() {
		ch <- 43
	})
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go test.Tap(cnf.dismiss)
	assert.EqualValues(t, <-ch, 43)
	assert.EqualValues(t, <-ch, 42)
	assert.True(t, cnf.win.Hidden)
}

func TestDialog_ConfirmCallbackOnlyOnClosed(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", nil, test.NewWindow(nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.SetOnClosed(func() {
		ch <- 43
	})
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go test.Tap(cnf.dismiss)
	assert.EqualValues(t, <-ch, 43)
	assert.True(t, cnf.win.Hidden)
}

func TestDialog_ConfirmCallbackOnlyOnConfirm(t *testing.T) {
	ch := make(chan int)
	cnf := NewConfirm("Test", "Test", func(_ bool) {
		ch <- 42
	}, test.NewWindow(nil))
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.Show()

	assert.False(t, cnf.win.Hidden)
	go test.Tap(cnf.dismiss)
	assert.EqualValues(t, <-ch, 42)
	assert.True(t, cnf.win.Hidden)
}
