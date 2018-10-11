package dialog

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/widget"
)

// NewConfirm creates a dialog over the specified window for user confirmation.
// The title is used for the ialog window and message is the content.
// The callback is executed when the user decides. After creation you should call Show().
func NewConfirm(title, message string, callback func(bool), parent fyne.Window) Dialog {
	d := newDialog(title, message, theme.QuestionIcon(), callback, parent)
	d.setButtons(newButtonList(
		&widget.Button{Text: "No",
			OnTapped: func() {
				d.response <- false
			},
		},
		&widget.Button{Text: "Yes", Style: widget.PrimaryButton,
			OnTapped: func() {
				d.response <- true
			},
		},
	))

	return d
}

// ShowConfirm shows a dialog over the specified window for a user
// confirmation. The title is used for the dialog window and message is the content.
// The callback is executed when the user decides.
func ShowConfirm(title, message string, callback func(bool), parent fyne.Window) {
	NewConfirm(title, message, callback, parent).Show()
}
