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
}

// SetText changes the current text value of the entry dialog, this can
// be useful for setting a default value.
func (i *EntryDialog) SetText(s string) {
	i.entry.SetText(s)
}

// GetText returns the current text value entered in the entry dialog.
func (i *EntryDialog) GetText() string {
	return i.entry.Text
}

// SetPlaceholder defines the placeholder text for the entry
func (i *EntryDialog) SetPlaceholder(s string) {
	i.entry.SetPlaceHolder(s)
}

// NewEntryDialog creates a dialog over the specified window for the user to
// enter a value.
//
// onConfirm is a callback that runs when the user enters a string of
// text and clicks the "confirm" button. May be nil.
//
// onClosed is called unconditionally weather the user confirms or cancels. The
// argument will be true if the user confirmed, and false otherwise. May be
// nil.
//
// Note that onClosed will be called after onConfirm, if both are non-nil. This
// way onConfirm can potential modify state that onClosed needs to get the user
// input when the user confirms, while also being able to handle the case where
// the user cancelled.
func NewEntryDialog(title, message string, onConfirm func(string), onClosed func(bool), parent fyne.Window) *EntryDialog {

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
		if onClosed != nil {
			onClosed(false)
		}
	})
	d.dismiss.Icon = theme.CancelIcon()

	// create confirmation button
	confirm := widget.NewButton("Ok", func() {
		// User has confirmed and entered an input
		if onConfirm != nil {
			onConfirm(entry.Text)
		}

		if onClosed != nil {
			onClosed(true)
		}

		// Also hide the dialog, and trigger it's callback
		d.hideWithResponse(true)
	})

	// attach response buttons to the dialog
	d.setButtons(newButtonList(d.dismiss, confirm))

	// and instantiate ourselves
	i := &EntryDialog{d, entry, confirm}

	return i
}

// ShowEntryDialog creates a new entry dialog and shows it immediately.
func ShowEntryDialog(title, message string, onConfirm func(string), onClosed func(bool), parent fyne.Window) {
	NewEntryDialog(title, message, onConfirm, onClosed, parent).Show()
}
