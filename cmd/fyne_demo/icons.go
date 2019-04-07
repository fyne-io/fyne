package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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
	bar := widget.NewHBox(prev, next, b.name)

	background := canvas.NewRasterWithPixels(checkerPattern)
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
	{"SearchIcon", theme.SearchIcon()},
	{"SearchReplaceIcon", theme.SearchReplaceIcon()},

	{"CheckButtonIcon", theme.CheckButtonIcon()},
	{"CheckButtonCheckedIcon", theme.CheckButtonCheckedIcon()},
	{"RadioButtonIcon", theme.RadioButtonIcon()},
	{"RadioButtonCheckedIcon", theme.RadioButtonCheckedIcon()},

	{"ContentAddIcon", theme.ContentAddIcon()},
	{"ContentRemoveIcon", theme.ContentRemoveIcon()},
	{"ContentClearIcon", theme.ContentClearIcon()},
	{"ContentCutIcon", theme.ContentCutIcon()},
	{"ContentCopyIcon", theme.ContentCopyIcon()},
	{"ContentPasteIcon", theme.ContentPasteIcon()},
	{"ContentRedoIcon", theme.ContentRedoIcon()},
	{"ContentUndoIcon", theme.ContentUndoIcon()},

	{"InfoIcon", theme.InfoIcon()},
	{"QuestionIcon", theme.QuestionIcon()},
	{"WarningIcon", theme.WarningIcon()},

	{"DocumentCreateIcon", theme.DocumentCreateIcon()},
	{"DocumentPrintIcon", theme.DocumentPrintIcon()},
	{"DocumentSaveIcon", theme.DocumentSaveIcon()},

	{"FolderIcon", theme.FolderIcon()},
	{"FolderNewIcon", theme.FolderNewIcon()},
	{"FolderOpenIcon", theme.FolderOpenIcon()},
	{"HomeIcon", theme.HomeIcon()},
	{"HelpIcon", theme.HelpIcon()},

	{"ViewFullScreenIcon", theme.ViewFullScreenIcon()},
	{"ViewRestoreIcon", theme.ViewRestoreIcon()},
	{"ZoomFitIcon", theme.ZoomFitIcon()},
	{"ZoomInIcon", theme.ZoomInIcon()},
	{"ZoomOutIcon", theme.ZoomOutIcon()},

	{"MoveDownIcon", theme.MoveDownIcon()},
	{"MoveUpIcon", theme.MoveUpIcon()},

	{"NavigateBackIcon", theme.NavigateBackIcon()},
	{"NavigateNextIcon", theme.NavigateNextIcon()},

	{"MailAttachmentIcon", theme.MailAttachmentIcon()},
	{"MailComposeIcon", theme.MailComposeIcon()},
	{"MailForwardIcon", theme.MailForwardIcon()},
	{"MailReplyIcon", theme.MailReplyIcon()},
	{"MailReplyAllIcon", theme.MailReplyAllIcon()},
	{"MailSendIcon", theme.MailSendIcon()},
}
