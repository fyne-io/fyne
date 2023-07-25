package dialog

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// formDialogResult is the result of the test form dialog callback.
type formDialogResult int

const (
	formDialogNoAction formDialogResult = iota
	formDialogConfirm
	formDialogCancel
)

func TestFormDialog_Control(t *testing.T) {
	var result formDialogResult
	fd := controlFormDialog(&result, test.NewWindow(nil))
	fd.Show()
	test.Tap(fd.confirm)

	assert.Equal(t, formDialogConfirm, result, "Control form should be confirmed with no validation")
}

func TestFormDialog_InvalidCannotSubmit(t *testing.T) {
	var result formDialogResult
	fd := validatingFormDialog(&result, test.NewWindow(nil))
	fd.Show()

	assert.False(t, fd.win.Hidden)
	assert.True(t, fd.confirm.Disabled(), "Confirm button should be disabled due to validation state")
	test.Tap(fd.confirm)

	assert.Equal(t, formDialogNoAction, result, "Callback should not have ran with invalid form")
}

func TestFormDialog_ValidCanSubmit(t *testing.T) {
	var result formDialogResult
	fd := validatingFormDialog(&result, test.NewWindow(nil))
	fd.Show()

	assert.False(t, fd.win.Hidden)
	assert.True(t, fd.confirm.Disabled(), "Confirm button should be disabled due to validation state")

	if validatingEntry, ok := fd.items[0].Widget.(*widget.Entry); ok {
		validatingEntry.SetText("abc")
		assert.False(t, fd.confirm.Disabled())
		test.Tap(fd.confirm)

		assert.Equal(t, formDialogConfirm, result, "Valid form should be able to be confirmed")
	} else {
		assert.Fail(t, "First item's widget should be an Entry (check validatingFormDialog)")
	}
}

func TestFormDialog_CanCancelInvalid(t *testing.T) {
	var result formDialogResult
	fd := validatingFormDialog(&result, test.NewWindow(nil))
	fd.Show()
	assert.False(t, fd.win.Hidden)

	test.Tap(fd.dismiss)

	assert.Equal(t, formDialogCancel, result, "Expected cancel result")
}

func TestFormDialog_CanCancelNoValidation(t *testing.T) {
	var result formDialogResult
	fd := controlFormDialog(&result, test.NewWindow(nil))
	fd.Show()
	assert.False(t, fd.win.Hidden)

	test.Tap(fd.dismiss)

	assert.Equal(t, formDialogCancel, result, "Expected cancel result")
}

func TestFormDialog_Hints(t *testing.T) {

	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())
	w := test.NewWindow(nil)
	w.SetFullScreen(true)

	var result formDialogResult
	fd := hintsFormDialog(&result, w)
	w.Resize(fd.MinSize())
	fd.Show()
	assert.False(t, fd.win.Hidden)

	assert.True(t, fd.confirm.Disabled(), "Confirm button should be disabled due to validation state")

	test.AssertImageMatches(t, "form/hint_initial.png", w.Canvas().Capture())

	validatingEntry, ok := fd.items[0].Widget.(*widget.Entry)
	require.True(t, ok, "First item's widget should be an Entry (check hintsFormDialog)")

	validatingEntry.SetText("n")
	test.AssertImageMatches(t, "form/hint_invalid.png", w.Canvas().Capture())
	assert.True(t, fd.confirm.Disabled())

	validatingEntry.SetText("abc")
	test.AssertImageMatches(t, "form/hint_valid.png", w.Canvas().Capture())
	assert.False(t, fd.confirm.Disabled())

	test.Tap(fd.confirm)
	assert.Equal(t, formDialogConfirm, result, "Valid form should be able to be confirmed")

	test.Tap(fd.dismiss)
	assert.Equal(t, formDialogCancel, result, "Expected cancel result")
}

func TestFormDialog_Submit(t *testing.T) {
	validatingEntry := widget.NewEntry()
	validatingEntry.Validator = func(input string) error {
		if input != "abc" {
			return errors.New("only accepts 'abc'")
		}
		return nil
	}
	validatingItem := &widget.FormItem{Widget: validatingEntry}

	confirmed := false

	items := []*widget.FormItem{validatingItem}
	form := NewForm("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		confirmed = confirm
	}, test.NewWindow(nil))

	form.Show()
	validatingEntry.SetText("cba")

	form.Submit()
	assert.Equal(t, false, confirmed)
	assert.Equal(t, false, form.win.Hidden)

	validatingEntry.SetText("abc")

	form.Submit()
	assert.Equal(t, true, confirmed)
	assert.Equal(t, true, form.win.Hidden)
}

func validatingFormDialog(result *formDialogResult, parent fyne.Window) *FormDialog {
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
	controlEntry := widget.NewPasswordEntry()
	controlItem := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry,
	}

	items := []*widget.FormItem{validatingItem, controlItem}
	return NewForm("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		if confirm {
			*result = formDialogConfirm
		} else {
			*result = formDialogCancel
		}
	}, parent)
}

func controlFormDialog(result *formDialogResult, parent fyne.Window) *FormDialog {
	controlEntry := widget.NewEntry()
	controlItem := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry,
	}
	controlEntry2 := widget.NewPasswordEntry()
	controlItem2 := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry2,
	}
	items := []*widget.FormItem{controlItem, controlItem2}
	return NewForm("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		if confirm {
			*result = formDialogConfirm
		} else {
			*result = formDialogCancel
		}
	}, parent)
}

func hintsFormDialog(result *formDialogResult, parent fyne.Window) *FormDialog {
	validatingEntry := widget.NewEntry()
	validatingEntry.Validator = func(input string) error {
		if input != "abc" {
			return errors.New("only accepts 'abc'")
		}
		return nil
	}
	validatingItem := &widget.FormItem{
		Text:     "Only accepts 'abc'",
		Widget:   validatingEntry,
		HintText: "Type abc",
	}

	items := []*widget.FormItem{validatingItem}
	return NewForm("Validating Form Dialog With Hints", "Submit", "Cancel", items, func(confirm bool) {
		if confirm {
			*result = formDialogConfirm
		} else {
			*result = formDialogCancel
		}
	}, parent)
}
