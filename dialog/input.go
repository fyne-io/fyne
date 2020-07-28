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
func NewInput(title, message string, onConfirm func(string), parent fyne.Window) *InputDialog {

	// create the widgets necessary for the dialog
	entry := widget.NewEntry()

	// content container for our widgets
	content := widget.NewHBox(
		widget.NewLabel(message),
		entry,
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

	// attach response buttons to the dialog
	d.setButtons(newButtonList(d.dismiss, confirm))

	// and instantiate ourselves
	i := &InputDialog{d, onConfirm, entry}

	// now we have everything we need for the confirmation button
	confirm.OnTapped = func() {

		// User has confirmed and entered a valid input
		if i.onConfirm != nil {
			i.onConfirm(entry.Text)
		}

		// Also hide the dialog, and trigger it's callback
		d.hideWithResponse(true)
	}

	return i
}

// ShowInput creates a new input dialog and shows it immediately.
func ShowInput(title, message string, onConfirm func(string), parent fyne.Window) {
	NewInput(title, message, onConfirm, parent).Show()
}
