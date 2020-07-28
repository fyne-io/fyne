package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// InputDialog is a variation of a dialog which prompts the user to enter some
// text.
type InputDialog struct {
	*dialog

	onConfirm func(string)

	entry *widget.Entry

	confirmButton *widget.Button
}

// SetText changes the current text value of the input dialog, this can
// be useful for setting a default value.
func (i *InputDialog) SetText(s string) {
	i.entry.SetText(s)
}

// GetText returns the current text value entered in the input dialog.
func (i *InputDialog) GetText() string {
	return i.entry.Text
}

// SetPlaceholder defines the placeholder text for the entry
func (i *InputDialog) SetPlaceHolder(s string) {
	i.entry.SetPlaceHolder(s)
}

// NewInput creates a dialog over the specified window for the user to enter
// a value.
//
// onConfirm is a callback that runs when the user enters a string of
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

	// create confirmation button
	confirm := widget.NewButton("Ok", func() {
		// User has confirmed and entered an input
		if onConfirm != nil {
			onConfirm(entry.Text)
		}

		// Also hide the dialog, and trigger it's callback
		d.hideWithResponse(true)
	})

	// attach response buttons to the dialog
	d.setButtons(newButtonList(d.dismiss, confirm))

	// and instantiate ourselves
	i := &InputDialog{d, onConfirm, entry, confirm}

	return i
}

// ShowInput creates a new input dialog and shows it immediately.
func ShowInput(title, message string, onConfirm func(string), parent fyne.Window) {
	NewInput(title, message, onConfirm, parent).Show()
}
