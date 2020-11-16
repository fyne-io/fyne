package dialog

import (
	"errors"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
	testify "github.com/stretchr/testify/assert"
)

type formDialogtestContext struct {
	subject *dialog
	confirm *widget.Button
	dismiss *widget.Button
	items   []*widget.FormItem
}

func TestFormDialog_Control(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	tc := controlFormDialog(ch, test.NewWindow(nil))
	go test.Tap(tc.confirm)

	select {
	case result := <-ch:
		assert.True(result, "Control should allow confirmation with no validation constraints")
	case <-time.After(500 * time.Millisecond):
		assert.Fail("Should have received a confirmation by now")
	}
}

func TestFormDialog_InvalidCannotSubmit(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	tc := validatingFormDialog(ch, test.NewWindow(nil))
	tc.subject.Show()

	assert.False(tc.subject.win.Hidden)
	assert.True(tc.confirm.Disabled(), "Confirm button should be disabled due to validation state")
	go test.Tap(tc.confirm)

	select {
	case <-ch:
		assert.Fail("Callback should not have ran with an invalid form")
	case <-time.After(500 * time.Millisecond):
	}
}

func TestFormDialog_ValidCanSubmit(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	tc := validatingFormDialog(ch, test.NewWindow(nil))
	tc.subject.Show()

	assert.False(tc.subject.win.Hidden)
	assert.True(tc.confirm.Disabled(), "Confirm button should be disabled due to validation state")

	if validatingEntry, ok := tc.items[0].Widget.(*widget.Entry); ok {
		validatingEntry.SetText("abc")
		assert.False(tc.confirm.Disabled())
		go test.Tap(tc.confirm)

		select {
		case result := <-ch:
			assert.True(result, "Confirm should return true result")
		case <-time.After(500 * time.Millisecond):
			assert.Fail("Callback should have ran with a valid form")
		}
	} else {
		assert.Fail("First item's widget should be an Entry (check validatingFormDialog)")
	}
}

func TestFormDialog_CanCancelInvalid(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	tc := validatingFormDialog(ch, test.NewWindow(nil))
	tc.subject.Show()
	assert.False(tc.subject.win.Hidden)

	go test.Tap(tc.dismiss)

	select {
	case result := <-ch:
		assert.False(result, "Result should be false with cancellation")
	case <-time.After(500 * time.Millisecond):
		assert.Fail("Should have received a cancellation by now")
	}
}

func TestFormDialog_CanCancelNoValidation(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	tc := controlFormDialog(ch, test.NewWindow(nil))
	tc.subject.Show()
	assert.False(tc.subject.win.Hidden)

	go test.Tap(tc.dismiss)

	select {
	case result := <-ch:
		assert.False(result, "Result should be false with cancellation")
	case <-time.After(500 * time.Millisecond):
		assert.Fail("Should have received a cancellation by now")
	}
}

func validatingFormDialog(ch chan bool, parent fyne.Window) formDialogtestContext {
	validatingEntry := widget.NewEntry()
	validatingEntry.Validator = func(input string) error {
		if input != "abc" {
			return errors.New("only accepts 'abc'")
		}
		return nil
	}
	validatingItem := &widget.FormItem{
		Text:   "Only accepts 'abc'",
		Widget: validatingEntry,
	}
	controlEntry := widget.NewEntry()
	controlItem := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry,
	}

	items := []*widget.FormItem{validatingItem, controlItem}
	formDialog, confirm, dismiss := testableNewFormDialog("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		ch <- confirm
	}, parent)
	tc := formDialogtestContext{
		subject: formDialog,
		confirm: confirm,
		dismiss: dismiss,
		items:   items,
	}
	return tc
}

func controlFormDialog(ch chan bool, parent fyne.Window) formDialogtestContext {
	controlEntry := widget.NewEntry()
	controlItem := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry,
	}
	controlEntry2 := widget.NewEntry()
	controlItem2 := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry2,
	}
	items := []*widget.FormItem{controlItem, controlItem2}
	formDialog, confirm, dismiss := testableNewFormDialog("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		ch <- confirm
	}, parent)
	tc := formDialogtestContext{
		subject: formDialog,
		confirm: confirm,
		dismiss: dismiss,
		items:   items,
	}
	return tc
}
