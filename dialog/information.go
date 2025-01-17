package dialog

import (
	"unicode"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createInformationDialog(title, message string, icon fyne.Resource, parent fyne.Window) Dialog {
	d := newTextDialog(title, message, icon, parent)
	d.dismiss = &widget.Button{
		Text:     lang.L("OK"),
		OnTapped: d.Hide,
	}
	d.create(container.NewGridWithColumns(1, d.dismiss))
	return d
}

// NewInformation creates a dialog over the specified window for user information.
// The title is used for the dialog window and message is the content.
// After creation you should call Show().
func NewInformation(title, message string, parent fyne.Window) Dialog {
	return createInformationDialog(title, message, theme.InfoIcon(), parent)
}

// ShowInformation shows a dialog over the specified window for user information.
// The title is used for the dialog window and message is the content.
func ShowInformation(title, message string, parent fyne.Window) {
	NewInformation(title, message, parent).Show()
}

// NewError creates a dialog over the specified window for an application error.
// The message is extracted from the provided error (should not be nil).
// After creation you should call Show().
func NewError(err error, parent fyne.Window) Dialog {
	dialogText := err.Error()
	r, size := utf8.DecodeRuneInString(dialogText)
	if r != utf8.RuneError {
		dialogText = string(unicode.ToUpper(r)) + dialogText[size:]
	}
	return createInformationDialog(lang.L("Error"), dialogText, theme.ErrorIcon(), parent)
}

// ShowError shows a dialog over the specified window for an application error.
// The message is extracted from the provided error (should not be nil).
func ShowError(err error, parent fyne.Window) {
	NewError(err, parent).Show()
}
