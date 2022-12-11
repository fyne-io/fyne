package tutorials

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type iconInfo struct {
	name string
	icon fyne.Resource
}

type browser struct {
	current int
	icons   []iconInfo

	name *widget.Select
	icon *widget.Icon
}

func (b *browser) setIcon(index int) {
	if index < 0 || index > len(b.icons)-1 {
		return
	}
	b.current = index

	b.name.SetSelected(b.icons[index].name)
	b.icon.SetResource(b.icons[index].icon)
}

// iconScreen loads a panel that shows the various icons available in Fyne
func iconScreen(_ fyne.Window) fyne.CanvasObject {
	b := &browser{}
	b.icons = loadIcons()

	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		b.setIcon(b.current - 1)
	})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		b.setIcon(b.current + 1)
	})
	b.name = widget.NewSelect(iconList(b.icons), func(name string) {
		for i, icon := range b.icons {
			if icon.name == name {
				if b.current != i {
					b.setIcon(i)
				}
				break
			}
		}
	})
	b.name.SetSelected(b.icons[b.current].name)
	buttons := container.NewHBox(prev, next)
	bar := container.NewBorder(nil, nil, buttons, nil, b.name)

	background := canvas.NewRasterWithPixels(checkerPattern)
	background.SetMinSize(fyne.NewSize(280, 280))
	b.icon = widget.NewIcon(b.icons[b.current].icon)

	return container.NewBorder(bar, nil, nil, nil, background, b.icon)
}

func checkerPattern(x, y, _, _ int) color.Color {
	x /= 20
	y /= 20

	if x%2 == y%2 {
		return theme.BackgroundColor()
	}

	return theme.ShadowColor()
}

func iconList(icons []iconInfo) []string {
	ret := make([]string, len(icons))
	for i, icon := range icons {
		ret[i] = icon.name
	}

	return ret
}

func loadIcons() []iconInfo {
	return []iconInfo{
		{"CancelIcon", theme.CancelIcon()},
		{"ConfirmIcon", theme.ConfirmIcon()},
		{"DeleteIcon", theme.DeleteIcon()},
		{"SearchIcon", theme.SearchIcon()},
		{"SearchReplaceIcon", theme.SearchReplaceIcon()},

		{"CheckButtonIcon", theme.CheckButtonIcon()},
		{"CheckButtonCheckedIcon", theme.CheckButtonCheckedIcon()},
		{"RadioButtonIcon", theme.RadioButtonIcon()},
		{"RadioButtonCheckedIcon", theme.RadioButtonCheckedIcon()},

		{"ColorAchromaticIcon", theme.ColorAchromaticIcon()},
		{"ColorChromaticIcon", theme.ColorChromaticIcon()},
		{"ColorPaletteIcon", theme.ColorPaletteIcon()},

		{"ContentAddIcon", theme.ContentAddIcon()},
		{"ContentRemoveIcon", theme.ContentRemoveIcon()},
		{"ContentClearIcon", theme.ContentClearIcon()},
		{"ContentCutIcon", theme.ContentCutIcon()},
		{"ContentCopyIcon", theme.ContentCopyIcon()},
		{"ContentPasteIcon", theme.ContentPasteIcon()},
		{"ContentRedoIcon", theme.ContentRedoIcon()},
		{"ContentUndoIcon", theme.ContentUndoIcon()},

		{"InfoIcon", theme.InfoIcon()},
		{"ErrorIcon", theme.ErrorIcon()},
		{"QuestionIcon", theme.QuestionIcon()},
		{"WarningIcon", theme.WarningIcon()},

		{"DocumentIcon", theme.DocumentIcon()},
		{"DocumentCreateIcon", theme.DocumentCreateIcon()},
		{"DocumentPrintIcon", theme.DocumentPrintIcon()},
		{"DocumentSaveIcon", theme.DocumentSaveIcon()},

		{"FileIcon", theme.FileIcon()},
		{"FileApplicationIcon", theme.FileApplicationIcon()},
		{"FileAudioIcon", theme.FileAudioIcon()},
		{"FileImageIcon", theme.FileImageIcon()},
		{"FileTextIcon", theme.FileTextIcon()},
		{"FileVideoIcon", theme.FileVideoIcon()},
		{"FolderIcon", theme.FolderIcon()},
		{"FolderNewIcon", theme.FolderNewIcon()},
		{"FolderOpenIcon", theme.FolderOpenIcon()},
		{"ComputerIcon", theme.ComputerIcon()},
		{"HomeIcon", theme.HomeIcon()},
		{"HelpIcon", theme.HelpIcon()},
		{"HistoryIcon", theme.HistoryIcon()},
		{"SettingsIcon", theme.SettingsIcon()},
		{"StorageIcon", theme.StorageIcon()},
		{"DownloadIcon", theme.DownloadIcon()},
		{"UploadIcon", theme.UploadIcon()},

		{"ViewFullScreenIcon", theme.ViewFullScreenIcon()},
		{"ViewRestoreIcon", theme.ViewRestoreIcon()},
		{"ViewRefreshIcon", theme.ViewRefreshIcon()},
		{"VisibilityIcon", theme.VisibilityIcon()},
		{"VisibilityOffIcon", theme.VisibilityOffIcon()},
		{"ZoomFitIcon", theme.ZoomFitIcon()},
		{"ZoomInIcon", theme.ZoomInIcon()},
		{"ZoomOutIcon", theme.ZoomOutIcon()},

		{"MoreHorizontalIcon", theme.MoreHorizontalIcon()},
		{"MoreVerticalIcon", theme.MoreVerticalIcon()},

		{"MoveDownIcon", theme.MoveDownIcon()},
		{"MoveUpIcon", theme.MoveUpIcon()},

		{"NavigateBackIcon", theme.NavigateBackIcon()},
		{"NavigateNextIcon", theme.NavigateNextIcon()},

		{"Menu", theme.MenuIcon()},
		{"MenuExpand", theme.MenuExpandIcon()},
		{"MenuDropDown", theme.MenuDropDownIcon()},
		{"MenuDropUp", theme.MenuDropUpIcon()},

		{"MailAttachmentIcon", theme.MailAttachmentIcon()},
		{"MailComposeIcon", theme.MailComposeIcon()},
		{"MailForwardIcon", theme.MailForwardIcon()},
		{"MailReplyIcon", theme.MailReplyIcon()},
		{"MailReplyAllIcon", theme.MailReplyAllIcon()},
		{"MailSendIcon", theme.MailSendIcon()},

		{"MediaFastForward", theme.MediaFastForwardIcon()},
		{"MediaFastRewind", theme.MediaFastRewindIcon()},
		{"MediaPause", theme.MediaPauseIcon()},
		{"MediaPlay", theme.MediaPlayIcon()},
		{"MediaStop", theme.MediaStopIcon()},
		{"MediaRecord", theme.MediaRecordIcon()},
		{"MediaReplay", theme.MediaReplayIcon()},
		{"MediaSkipNext", theme.MediaSkipNextIcon()},
		{"MediaSkipPrevious", theme.MediaSkipPreviousIcon()},

		{"VolumeDown", theme.VolumeDownIcon()},
		{"VolumeMute", theme.VolumeMuteIcon()},
		{"VolumeUp", theme.VolumeUpIcon()},

		{"AccountIcon", theme.AccountIcon()},
		{"LoginIcon", theme.LoginIcon()},
		{"LogoutIcon", theme.LogoutIcon()},

		{"ListIcon", theme.ListIcon()},
		{"GridIcon", theme.GridIcon()},
	}
}
