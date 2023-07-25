package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ConfirmDialog is like the standard Dialog but with an additional confirmation button
type ConfirmDialog struct {
	*dialog

	confirm *widget.Button
}

// SetConfirmText allows custom text to be set in the confirmation button
func (d *ConfirmDialog) SetConfirmText(label string) {
	d.confirm.SetText(label)
	d.win.Refresh()
}

// SetConfirmImportance sets the importance level of the confirm button.
//
// Since 2.4
func (d *ConfirmDialog) SetConfirmImportance(importance widget.Importance) {
	d.confirm.Importance = importance
}

// NewConfirm creates a dialog over the specified window for user confirmation.
// The title is used for the dialog window and message is the content.
// The callback is executed when the user decides. After creation you should call Show().
func NewConfirm(title, message string, callback func(bool), parent fyne.Window) *ConfirmDialog {
	d := newDialog(title, message, theme.QuestionIcon(), callback, parent)

	d.dismiss = &widget.Button{Text: "No", Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	confirm := &widget.Button{Text: "Yes", Icon: theme.ConfirmIcon(), Importance: widget.HighImportance,
		OnTapped: func() {
			d.hideWithResponse(true)
		},
	}
	d.create(container.NewGridWithColumns(2, d.dismiss, confirm))

	return &ConfirmDialog{dialog: d, confirm: confirm}
}

// ShowConfirm shows a dialog over the specified window for a user
// confirmation. The title is used for the dialog window and message is the content.
// The callback is executed when the user decides.
func ShowConfirm(title, message string, callback func(bool), parent fyne.Window) {
	NewConfirm(title, message, callback, parent).Show()
}
