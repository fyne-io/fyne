package dialog

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/widget"
)

// NewInformation creates a dialog over the specified window for user information.
// The title is used for the ialog window and message is the content.
// After creation you should call Show().
func NewInformation(title, message string, parent fyne.Window) Dialog {
	d := newDialog(title, message, theme.InfoIcon(), nil, parent)
	d.setButtons(newButtonList(&widget.Button{Text: "OK", Style: widget.PrimaryButton,
		OnTapped: func() {
			d.response <- false
		},
	}))

	return d
}

// ShowInformation shows a dialog over the specified window for user
// information. The title is used for the dialog window and message is the content.
func ShowInformation(title, message string, parent fyne.Window) {
	NewInformation(title, message, parent).Show()
}

// ShowError shows a dialog over the specified window for an application
// error. The title and message are extracted from the provided error.
func ShowError(err error, parent fyne.Window) {
	ShowInformation("Error", err.Error(), parent)
}
