package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// ProgressInfiniteDialog is a simple dialog window that displays text and a infinite progress bar.
type ProgressInfiniteDialog struct {
	*dialog

	bar *widget.ProgressBarInfinite
}

// NewProgressInfinite creates a infinite progress dialog and returns the handle.
// Using the returned type you should call Show().
func NewProgressInfinite(title, message string, parent fyne.Window) *ProgressInfiniteDialog {
	d := newDialog(title, message, theme.InfoIcon(), nil /*cancel?*/, parent)
	bar := widget.NewProgressBarInfinite()
	bar.Resize(fyne.NewSize(200, bar.MinSize().Height))

	d.setButtons(bar)
	return &ProgressInfiniteDialog{d, bar}
}

// Hide this dialog and stop the infinite progress goroutine
func (d *ProgressInfiniteDialog) Hide() {
	d.bar.Hide()
	d.dialog.Hide()
}
