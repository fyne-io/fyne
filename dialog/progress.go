package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type ProgressDialog struct {
	*dialog

	bar *widget.ProgressBar
}

func (p *ProgressDialog) SetValue(v float64) {
	p.bar.SetValue(v)
}

func NewProgress(title, message string, parent fyne.Window) *ProgressDialog {
	d := newDialog(title, message, theme.InfoIcon(), nil /*cancel?*/, parent)
	bar := widget.NewProgressBar()
	bar.Resize(fyne.NewSize(200, bar.MinSize().Height))

	d.setButtons(bar)
	return &ProgressDialog{d, bar}
}
