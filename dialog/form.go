package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// formDialog is a simple dialog window for displaying FormItems inside a form.
type formDialog struct {
	*dialog
	items   []*widget.FormItem
	confirm *widget.Button
	cancel  *widget.Button
}

// validateItems acts as a validation edge state handler that will respond to an individual widget's validation
// state before checking all others to determine the net validation state. If the error passed is not nil, then the
// confirm button will be disabled. If the error parameter is nil, then all other Validatable widgets in items are
// checked as well to determine whether the confirm button should be disabled.
// This method is passed to each Validatable widget's SetOnValidationChanged method in items by NewFormDialog.
func (d *formDialog) validateItems(err error) {
	if err != nil {
		d.confirm.Disable()
		return
	}
	for _, item := range d.items {
		if validatable, ok := item.Widget.(fyne.Validatable); ok {
			if err := validatable.Validate(); err != nil {
				d.confirm.Disable()
				return
			}
		}
	}
	d.confirm.Enable()
}

// NewFormDialog creates and returns a dialog over the specified application using
// the provided FormItems. The cancel button will have the dismiss text set and the confirm button will
// use the confirm text. The response callback is called on user action after validation passes.
// If any Validatable widget reports that validation has failed, then the confirm
// button will be disabled. The initial state of the confirm button will reflect the initial
// validation state of the items added to the form dialog.
func NewFormDialog(title, confirm, dismiss string, items []*widget.FormItem, callback func(bool),
	parent fyne.Window) Dialog {
	var itemObjects = make([]fyne.CanvasObject, len(items)*2)
	for i, item := range items {
		itemObjects[i*2] = widget.NewLabel(item.Text)
		itemObjects[i*2+1] = item.Widget
	}
	content := fyne.NewContainerWithLayout(layout.NewFormLayout(), itemObjects...)
	d := &dialog{content: content, title: title, icon: nil, parent: parent}
	d.callback = callback
	// TODO: Copied from NewCustomConfirm above.
	// This is still a problem because commenting out the `.Show()` call below will still result in the
	// dialog being shown.
	d.sendResponse = true

	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	confirmBtn := &widget.Button{Text: confirm, Icon: nil, Importance: widget.HighImportance,
		OnTapped: func() {
			d.hideWithResponse(true)
		}}
	// Mitigation for issue #1553
	confirmBtn.SetIcon(theme.ConfirmIcon())
	formDialog := &formDialog{
		dialog:  d,
		items:   items,
		confirm: confirmBtn,
		cancel:  d.dismiss,
	}
	formDialog.validateItems(nil)

	for _, item := range items {
		if validatable, ok := item.Widget.(fyne.Validatable); ok {
			validatable.SetOnValidationChanged(formDialog.validateItems)
		}
	}
	d.setButtons(container.NewHBox(layout.NewSpacer(), d.dismiss, confirmBtn, layout.NewSpacer()))
	return formDialog
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
