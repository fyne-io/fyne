package main

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/layout"
	"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/widget"
	"image/color"
)

type browser struct {
	canvas  fyne.Canvas
	current int

	name *widget.Label
	icon *widget.Icon
}

func (b *browser) setIcon(index int) {
	if index < 0 || index > len(icons)-1 {
		return
	}
	b.current = index

	b.name.SetText(icons[index].name)
	b.icon.SetResource(icons[index].icon)
}

// Icons loads a window that shows the various icons available in Fyne
func Icons(app fyne.App) {
	win := app.NewWindow("Icons")
	b := &browser{canvas: win.Canvas()}

	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		b.setIcon(b.current - 1)
	})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		b.setIcon(b.current + 1)
	})
	b.name = widget.NewLabel(icons[b.current].name)
	bar := widget.NewHBox(prev, next, b.name, layout.NewSpacer())

	background := canvas.NewRaster(checkerPattern)
	background.SetMinSize(fyne.NewSize(280, 280))
	b.icon = widget.NewIcon(icons[b.current].icon)

	win.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(
		bar, nil, nil, nil), bar, background, b.icon))
	win.Show()
}

func checkerPattern(x, y, _, _ int) color.Color {
	x /= 20
	y /= 20

	if x%2 == y%2 {
		return theme.BackgroundColor()
	}

	return theme.ButtonColor()
}

var icons = []struct {
	name string
	icon fyne.Resource
}{
	{"CancelIcon", theme.CancelIcon()},
	{"ConfirmIcon", theme.ConfirmIcon()},
	{"DeleteIcon", theme.DeleteIcon()},
	{"CheckedIcon", theme.CheckedIcon()},
	{"UnCheckedIcon", theme.UncheckedIcon()},

	{"CutIcon", theme.CutIcon()},
	{"CopyIcon", theme.CopyIcon()},
	{"PasteIcon", theme.PasteIcon()},

	{"InfoIcon", theme.InfoIcon()},
	{"QuestionIcon", theme.QuestionIcon()},
	{"WarningIcon", theme.WarningIcon()},

	{"MailComposeIcon", theme.MailComposeIcon()},
	{"MailForwardIcon", theme.MailForwardIcon()},
	{"MailReplyIcon", theme.MailReplyIcon()},
	{"MailReplyAllIcon", theme.MailReplyAllIcon()},

	{"NavigateBackIcon", theme.NavigateBackIcon()},
	{"NavigateNextIcon", theme.NavigateNextIcon()},
}
