package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ProgressDialog is a simple dialog window that displays text and a progress bar.
//
// Deprecated: Create a new custom dialog with a widget.ProgressBar() inside.
type ProgressDialog struct {
	*dialog

	bar *widget.ProgressBar
}

// SetValue updates the value of the progress bar - this should be between 0.0 and 1.0.
func (p *ProgressDialog) SetValue(v float64) {
	p.bar.SetValue(v)
}

// NewProgress creates a progress dialog and returns the handle.
// Using the returned type you should call Show() and then set its value through SetValue().
//
// Deprecated: Create a new custom dialog with a widget.ProgressBar() inside.
func NewProgress(title, message string, parent fyne.Window) *ProgressDialog {
	d := newDialog(title, message, theme.InfoIcon(), nil /*cancel?*/, parent)
	bar := widget.NewProgressBar()
	bar.Resize(fyne.NewSize(200, bar.MinSize().Height))

	d.setButtons(bar)
	return &ProgressDialog{d, bar}
}
