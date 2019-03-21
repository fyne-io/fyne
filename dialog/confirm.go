package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// ConfirmDialog is like the standard Dialog but with an additional confirmation button
type ConfirmDialog struct {
	*dialog

	confirm *widget.Button
}

// SetConfirmText allows custom text to be set in the confirmation button
func (d *ConfirmDialog) SetConfirmText(label string) {
	d.confirm.SetText(label)
	d.Layout(d.win.Content().(*fyne.Container).Objects, d.win.Content().MinSize())
}

// NewConfirm creates a dialog over the specified window for user confirmation.
// The title is used for the ialog window and message is the content.
// The callback is executed when the user decides. After creation you should call Show().
func NewConfirm(title, message string, callback func(bool), parent fyne.Window) *ConfirmDialog {
	d := newDialog(title, message, theme.QuestionIcon(), callback, parent)

	d.dismiss = &widget.Button{Text: "No", Icon: theme.CancelIcon(),
		OnTapped: func() {
			d.response <- false
		},
	}
	confirm := &widget.Button{Text: "Yes", Icon: theme.ConfirmIcon(), Style: widget.PrimaryButton,
		OnTapped: func() {
			d.response <- true
		},
	}
	d.setButtons(newButtonList(d.dismiss, confirm))

	return &ConfirmDialog{d, confirm}
}

// ShowConfirm shows a dialog over the specified window for a user
// confirmation. The title is used for the dialog window and message is the content.
// The callback is executed when the user decides.
func ShowConfirm(title, message string, callback func(bool), parent fyne.Window) {
	NewConfirm(title, message, callback, parent).Show()
}
