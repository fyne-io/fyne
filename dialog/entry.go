package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// EntryDialog is a variation of a dialog which prompts the user to enter some
// text.
type EntryDialog struct {
	*dialog

	entry *widget.Entry

	confirmButton *widget.Button

	onClosed func()
}

// SetText changes the current text value of the entry dialog, this can
// be useful for setting a default value.
func (i *EntryDialog) SetText(s string) {
	i.entry.SetText(s)
}

// SetPlaceholder defines the placeholder text for the entry
func (i *EntryDialog) SetPlaceholder(s string) {
	i.entry.SetPlaceHolder(s)
}

// SetOnClosed changes the callback which is run when the dialog is closed,
// which is nil by default.
//
// The callback is called unconditionally whether the user confirms or cancels.
//
// Note that the callback will be called after onConfirm, if both are non-nil.
// This way onConfirm can potential modify state that this callback needs to
// get the user input when the user confirms, while also being able to handle
// the case where the user cancelled.
func (i *EntryDialog) SetOnClosed(callback func()) {
	i.onClosed = callback
}

// NewEntryDialog creates a dialog over the specified window for the user to
// enter a value.
//
// onConfirm is a callback that runs when the user enters a string of
// text and clicks the "confirm" button. May be nil.
//
func NewEntryDialog(title, message string, onConfirm func(string), parent fyne.Window) *EntryDialog {
	i := &EntryDialog{nil, nil, nil, nil}

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
		if i.onClosed != nil {
			i.onClosed()
		}
	})
	d.dismiss.Icon = theme.CancelIcon()

	// create confirmation button
	confirm := widget.NewButton("Ok", func() {
		// User has confirmed and entered an input
		if onConfirm != nil {
			onConfirm(entry.Text)
		}

		if i.onClosed != nil {
			i.onClosed()
		}

		// Also hide the dialog, and trigger it's callback
		d.hideWithResponse(true)
	})

	// attach response buttons to the dialog
	d.setButtons(newButtonList(d.dismiss, confirm))

	i.dialog = d
	i.entry = entry
	i.confirmButton = confirm
	i.onClosed = nil

	return i
}

// ShowEntryDialog creates a new entry dialog and shows it immediately.
func ShowEntryDialog(title, message string, onConfirm func(string), parent fyne.Window) {
	NewEntryDialog(title, message, onConfirm, parent).Show()
}
