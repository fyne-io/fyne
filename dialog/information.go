package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func createTextDialog(title, message string, icon fyne.Resource, parent fyne.Window) Dialog {
	d := newDialog(title, message, icon, nil, parent)

	d.dismiss = &widget.Button{Text: "OK",
		OnTapped: d.Hide,
	}
	d.setButtons(newButtonList(d.dismiss))

	return d
}

// NewInformation creates a dialog over the specified window for user information.
// The title is used for the dialog window and message is the content.
// After creation you should call Show().
func NewInformation(title, message string, parent fyne.Window) Dialog {
	return createTextDialog(title, message, theme.InfoIcon(), parent)
}

// ShowInformation shows a dialog over the specified window for user
// information. The title is used for the dialog window and message is the content.
func ShowInformation(title, message string, parent fyne.Window) {
	NewInformation(title, message, parent).Show()
}

// NewError creates a dialog over the specified window for an application
// error. The title and message are extracted from the provided error.
// After creation you should call Show().
func NewError(err error, parent fyne.Window) Dialog {
	return createTextDialog("Error", err.Error(), theme.ErrorIcon(), parent)
}

// ShowError shows a dialog over the specified window for an application
// error. The title and message are extracted from the provided error.
func ShowError(err error, parent fyne.Window) {
	NewError(err, parent).Show()
}
