package dialog

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	fileIconSize       = 64
	fileInlineIconSize = 24
	fileTextSize       = 24
	fileIconCellWidth  = fileIconSize * 1.25
)

type fileDialogItem struct {
	widget.BaseWidget
	picker *fileDialog

	name     string
	location fyne.URI
	dir      bool
}

func (i *fileDialogItem) CreateRenderer() fyne.WidgetRenderer {
	text := widget.NewLabelWithStyle(i.name, fyne.TextAlignCenter, fyne.TextStyle{})
	text.Wrapping = fyne.TextTruncate
	icon := widget.NewFileIcon(i.location)

	return &fileItemRenderer{
		item:    i,
		icon:    icon,
		text:    text,
		objects: []fyne.CanvasObject{icon, text},
	}
}

func (i *fileDialogItem) setLocation(l fyne.URI, dir, up bool) {
	i.dir = dir
	i.location = l
	i.name = l.Name()

	if i.picker.view == gridView {
		ext := filepath.Ext(i.name[1:])
		i.name = i.name[:len(i.name)-len(ext)]
	}

	if up {
		i.name = "(Parent)"
	}

	i.Refresh()
}

func (f *fileDialog) newFileItem(location fyne.URI, dir, up bool) *fileDialogItem {
	item := &fileDialogItem{
		picker:   f,
		location: location,
		name:     location.Name(),
		dir:      dir,
	}

	if f.view == gridView {
		ext := filepath.Ext(item.name[1:])
		item.name = item.name[:len(item.name)-len(ext)]
	}

	if up {
		item.name = "(Parent)"
	}

	item.ExtendBaseWidget(item)
	return item
}

type fileItemRenderer struct {
	item *fileDialogItem

	icon    *widget.FileIcon
	text    *widget.Label
	objects []fyne.CanvasObject
}

func (s fileItemRenderer) Layout(size fyne.Size) {
	if s.item.picker.view == gridView {
		s.icon.Resize(fyne.NewSize(fileIconSize, fileIconSize))
		s.icon.Move(fyne.NewPos((size.Width-fileIconSize)/2, 0))

		s.text.Alignment = fyne.TextAlignCenter
		s.text.Resize(fyne.NewSize(size.Width, fileTextSize))
		s.text.Move(fyne.NewPos(0, size.Height-s.text.MinSize().Height))
	} else {
		s.icon.Resize(fyne.NewSize(fileInlineIconSize, fileInlineIconSize))
		s.icon.Move(fyne.NewPos(theme.Padding(), (size.Height-fileInlineIconSize)/2))

		s.text.Alignment = fyne.TextAlignLeading
		s.text.Resize(fyne.NewSize(size.Width, fileTextSize))
		s.text.Move(fyne.NewPos(fileInlineIconSize, (size.Height-s.text.MinSize().Height)/2))
	}
	s.text.Refresh()
}

func (s fileItemRenderer) MinSize() fyne.Size {
	var padding fyne.Size

	if s.item.picker.view == gridView {
		padding = fyne.NewSize(fileIconCellWidth-fileIconSize, theme.Padding())
		return fyne.NewSize(fileIconSize, fileIconSize+fileTextSize).Add(padding)
	}

	padding = fyne.NewSize(theme.Padding(), theme.Padding()*4)
	return fyne.NewSize(fileInlineIconSize+s.text.MinSize().Width, fileTextSize).Add(padding)
}

func (s fileItemRenderer) Refresh() {
	s.text.SetText(s.item.name)
	s.icon.SetURI(s.item.location)
}

func (s fileItemRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s fileItemRenderer) Destroy() {
}
