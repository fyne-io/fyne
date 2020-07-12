package dialog

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type InputDialog struct {
	*dialog

	validator func(string) (bool, string)

	entry *widget.Entry
}

// SetText changes the current text value of the input dialog, this can
// be useful for setting a default value.
func (i *InputDialog) SetText(s string) {
	i.entry.SetText(s)
}

// newInput creates a dialog over the specified window for the user to enter
// a value, possibly with validation.
//
// validator is called when the entry text changes, and should return a bool
// indicating if the text is valid (valid => true), and a string message to the
// user. May be nil to allow any string.
func NewInput(title, message string, callback func(bool), validator func(string) (bool, string), parent fyne.Window) *InputDialog {

	entry := widget.NewEntry()
	icon := widget.NewIcon(theme.ConfirmIcon())
	response := widget.NewLabel("")

	content := widget.NewVBox(
		widget.NewHBox(
			widget.NewLabel(message),
			entry,
			icon,
		),
		response,
	)

	d := newDialog(title, message, theme.QuestionIcon(), callback, parent)
	d.content = content

	d.dismiss = widget.NewButton("Cancel", func() {
		d.Hide()
	})
	d.dismiss.Icon = theme.CancelIcon()

	confirm := widget.NewButton("Ok", func() {
		fmt.Printf("OK!\n")
		d.hideWithResponse(true)
	})

	d.setButtons(newButtonList(d.dismiss, confirm))

	i := &InputDialog{d, validator, entry}

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

	// set up the icon and such according to weather or not the validator
	// likes the empty string, which is our starting values
	entry.OnChanged("")

	return i
}

// ShowInput creates a new input dialog and shows it immediately.
func ShowInput(title, message string, callback func(bool), validator func(string) (bool, string), parent fyne.Window) {
	NewInput(title, message, callback, validator, parent).Show()
}
