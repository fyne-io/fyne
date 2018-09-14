package dialog

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
)
import "github.com/fyne-io/fyne/widget"

// ShowInformationDialog shows a dialog over the specified application for user
// information. The title is used for the dialog window and message is the content.
func ShowInformationDialog(title, message string, parent fyne.App) {
	dialog := newDialog(title, message, theme.InfoIcon(), nil, parent)
	dialog.setButtons(newButtonList(&widget.Button{Text: "OK",
		OnTapped: func() {
			dialog.response <- false
		},
	}))

	go dialog.wait()
	dialog.win.Show()
}

// ShowConfirmDialog shows a dialog over the specified application for a user
// confirmation. The title is used for the dialog window and message is the content.
// The callback is executed when the user decides.
func ShowConfirmDialog(title, message string, callback func(bool), parent fyne.App) {
	dialog := newDialog(title, message, theme.QuestionIcon(), callback, parent)
	dialog.setButtons(newButtonList(
		&widget.Button{Text: "No",
			OnTapped: func() {
				dialog.response <- false
			},
		},
		&widget.Button{Text: "Yes", Style: widget.PrimaryButton,
			OnTapped: func() {
				dialog.response <- true
			},
		},
	))

	go dialog.wait()
	dialog.win.Show()
}
