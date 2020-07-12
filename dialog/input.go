package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// InputDialog is a variation of a dialog which prompts the user to enter some
// text, with an optional validation function.
type InputDialog struct {
	*dialog

	validator func(string) (bool, string)
	onConfirm func(string)

	entry *widget.Entry
}

// SetText changes the current text value of the input dialog, this can
// be useful for setting a default value.
func (i *InputDialog) SetText(s string) {
	i.entry.SetText(s)
}

// NewInput creates a dialog over the specified window for the user to enter
// a value, possibly with validation.
//
// onConfirm is a callback that runs when the user enters a valid string of
// text and clicks the "confirm" button. May be nil.
//
// validator is called when the entry text changes, and should return a bool
// indicating if the text is valid (valid => true), and a string message to the
// user. May be nil to allow any string.
func NewInput(title, message string, onConfirm func(string), validator func(string) (bool, string), parent fyne.Window) *InputDialog {

	// create the widgets necessary for the dialog
	entry := widget.NewEntry()
	icon := widget.NewIcon(theme.ConfirmIcon())
	response := widget.NewLabel("") // response from the validator

	// content container for our widgets
	content := widget.NewVBox(
		widget.NewHBox(
			widget.NewLabel(message),
			entry,
			icon,
		),
		response,
	)

	// instantiate the dialog, and override the content
	d := newDialog(title, message, theme.QuestionIcon(), func(response bool) {}, parent)
	d.content = content

	// hide the dialog, and empty the text
	d.dismiss = widget.NewButton("Cancel", func() {
		entry.Text = ""
		d.Hide()
	})
	d.dismiss.Icon = theme.CancelIcon()

	// we will override the OnTapped function later once we've instanced
	// the variables we need
	confirm := widget.NewButton("Ok", func() {
	})
	confirm.Disable()

	// attach response buttons to the dialog
	d.setButtons(newButtonList(d.dismiss, confirm))

	// and instantiateourselves
	i := &InputDialog{d, validator, onConfirm, entry}

	// handle validation if the validator is non-nil, notice that we use
	// i.validator, in case the user has changed it later
	entry.OnChanged = func(s string) {
		if i.validator != nil {
			ok, msg := i.validator(s)
			if !ok {
				icon.Resource = theme.CancelIcon()
				confirm.Disable()
			} else {
				icon.Resource = theme.ConfirmIcon()
				confirm.Enable()
			}
			response.SetText(msg)
			icon.Refresh()
		}
	}

	// now we have everything we need for the confirmation button
	confirm.OnTapped = func() {

		// This shouldn't happen, since the button should be disabled
		// if validation has failed.
		ok, _ := i.validator(entry.Text)
		if !ok {
			return
		}

		// User has confirmed and entered a valid input
		if i.onConfirm != nil {
			i.onConfirm(entry.Text)
		}

		// Also hide the dialog, and trigger it's callback
		d.hideWithResponse(true)
	}

	// set up the icon and such according to weather or not the validator
	// likes the empty string, which is our starting values
	entry.OnChanged("")

	return i
}

// ShowInput creates a new input dialog and shows it immediately.
func ShowInput(title, message string, onConfirm func(string), validator func(string) (bool, string), parent fyne.Window) {
	NewInput(title, message, onConfirm, validator, parent).Show()
}
