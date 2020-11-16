package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type FormDialog struct {
	*dialog
	confirm *widget.Button
	cancel  *widget.Button
}

// NewFormDialog creates and returns a dialog over the specified application using
// the provided FormItems. The cancel button will have the dismiss text set and the confirm button will
// use the confirm text. The response callback is called on user action after validation passes.
// If any Validatable widget reports that validation has failed, then the confirm
// button will be disabled. The initial state of the confirm button will reflect the initial
// validation state of the items added to the form dialog.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewFormDialog(title, confirm, dismiss string, items []*widget.FormItem, callback func(bool),
	parent fyne.Window) Dialog {
	formDialog, _, _ := testableNewFormDialog(title, confirm, dismiss, items, callback, parent)
	return formDialog
}

func testableNewFormDialog(title string, confirm string, dismiss string, items []*widget.FormItem, callback func(bool),
	parent fyne.Window) (d *dialog, confirmBtn *widget.Button, dismissBtn *widget.Button) {
	var itemObjects = make([]fyne.CanvasObject, 0, len(items)*2)
	for _, ii := range items {
		itemObjects = append(itemObjects, widget.NewLabel(ii.Text), ii.Widget)
	}
	content := fyne.NewContainerWithLayout(layout.NewFormLayout(), itemObjects...)
	d = &dialog{content: content, title: title, icon: nil, parent: parent}
	d.callback = callback
	// TODO: Copied from NewCustomConfirm above.
	// This is still a problem because commenting out the `.Show()` call below will still result in the
	// dialog being shown.
	d.sendResponse = true

	dismissBtn = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	d.dismiss = dismissBtn
	confirmBtn = &widget.Button{Text: confirm, Icon: nil, Importance: widget.HighImportance,
		OnTapped: func() {
			d.hideWithResponse(true)
		}}
	// Mitigation for issue #1553
	confirmBtn.SetIcon(theme.ConfirmIcon())
	validateItems := func() {
		for _, item := range items {
			if validatable, canValidate := item.Widget.(fyne.Validatable); canValidate {
				if err := validatable.Validate(); err != nil {
					confirmBtn.Disable()
					return
				}
			}
		}
		confirmBtn.Enable()
	}
	validateItems()

	for _, item := range items {
		if validatable, canValidate := item.Widget.(fyne.Validatable); canValidate {
			validatable.SetOnValidationChanged(func(error) {
				validateItems()
			})
		}
	}
	d.setButtons(container.NewHBox(layout.NewSpacer(), d.dismiss, confirmBtn, layout.NewSpacer()))
	return
}

// ShowFormDialog shows a dialog over the specified application using
// the provided FormItems. The cancel button will have the dismiss text set and the confirm button will
// use the confirm text. The response callback is called on user action after validation passes.
// If any Validatable widget reports that validation has failed, then the confirm
// button will be disabled. The initial state of the confirm button will reflect the initial
// validation state of the items added to the form dialog.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowFormDialog(title, confirm, dismiss string, content []*widget.FormItem,
	callback func(bool), parent fyne.Window) {
	NewFormDialog(title, confirm, dismiss, content, callback, parent).Show()
}
