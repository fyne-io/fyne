package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	// absolute max width of text dialogs
	// (prevent them from looking unnaturally large on desktop)
	maxTextDialogAbsoluteWidth float32 = 600

	// max width of text dialogs as a percentage of the current window width
	maxTextDialogWinPcntWidth float32 = .9
)

func newTextDialog(title, message string, icon fyne.Resource, parent fyne.Window) *dialog {
	d := &dialog{
		title:   title,
		icon:    icon,
		parent:  parent,
		content: newCenterWrappedLabel(message),
	}
	d.beforeShowHook = createBeforeShowHook(d, message)

	return d
}

func createBeforeShowHook(d *dialog, message string) func() {
	return func() {
		// set the desired width of the dialog to the min of:
		// - width needed to show message without wrapping
		// - maxTextDialogAbsoluteWidth
		// - current window width * maxTextDialogWinPcntWidth
		if d.desiredSize.IsZero() {
			noWrapWidth := fyne.MeasureText(message, theme.TextSize(), fyne.TextStyle{}).Width + padWidth*2
			maxWinWitth := d.parent.Canvas().Size().Width * maxTextDialogWinPcntWidth
			w := fyne.Min(fyne.Min(noWrapWidth, maxTextDialogAbsoluteWidth), maxWinWitth)
			d.desiredSize = fyne.NewSize(w, d.MinSize().Height)
		}
	}
}

func newCenterWrappedLabel(message string) fyne.CanvasObject {
	return &widget.Label{Text: message, Alignment: fyne.TextAlignCenter, Wrapping: fyne.TextWrapWord}
}
