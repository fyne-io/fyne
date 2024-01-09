package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// EntryDialog is a variation of a dialog which prompts the user to enter some text.
//
// Deprecated: Use dialog.NewForm() or dialog.ShowForm() with a widget.Entry inside instead.
type EntryDialog struct {
	*FormDialog

	entry *widget.Entry

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

// NewEntryDialog creates a dialog over the specified window for the user to enter a value.
//
// onConfirm is a callback that runs when the user enters a string of
// text and clicks the "confirm" button. May be nil.
//
// Deprecated: Use dialog.NewForm() with a widget.Entry inside instead.
func NewEntryDialog(title, message string, onConfirm func(string), parent fyne.Window) *EntryDialog {
	i := &EntryDialog{entry: widget.NewEntry()}
	items := []*widget.FormItem{widget.NewFormItem(message, i.entry)}
	i.FormDialog = NewForm(title, "Ok", "Cancel", items, func(ok bool) {
		// User has confirmed and entered an input
		if ok && onConfirm != nil {
			onConfirm(i.entry.Text)
		}

		if i.onClosed != nil {
			i.onClosed()
		}

		i.entry.Text = ""
		i.win.Hide() // Close directly without executing the callback. This is the callback.
	}, parent)

	return i
}

// ShowEntryDialog creates a new entry dialog and shows it immediately.
//
// Deprecated: Use dialog.ShowForm() with a widget.Entry inside instead.
func ShowEntryDialog(title, message string, onConfirm func(string), parent fyne.Window) {
	NewEntryDialog(title, message, onConfirm, parent).Show()
}
