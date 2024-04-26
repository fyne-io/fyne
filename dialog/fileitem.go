package dialog

import (
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	fileIconSize       = 64
	fileInlineIconSize = 24
	fileIconCellWidth  = fileIconSize * 1.25
)

type fileDialogItem struct {
	widget.BaseWidget
	picker *fileDialog

	name     string
	id       int // id in the parent container
	choose   func(id int)
	open     func()
	location fyne.URI
	dir      bool

	lastClick time.Time
}

func (i *fileDialogItem) CreateRenderer() fyne.WidgetRenderer {
	text := widget.NewLabelWithStyle(i.name, fyne.TextAlignCenter, fyne.TextStyle{})
	text.Truncation = fyne.TextTruncateEllipsis
	text.Wrapping = fyne.TextWrapBreak
	icon := widget.NewFileIcon(i.location)

	return &fileItemRenderer{
		item:         i,
		icon:         icon,
		text:         text,
		objects:      []fyne.CanvasObject{icon, text},
		fileTextSize: widget.NewLabel("M\nM").MinSize().Height, // cache two-line label height,
	}
}

func (i *fileDialogItem) setLocation(l fyne.URI, dir, up bool) {
	i.dir = dir
	i.location = l
	i.name = l.Name()

	if i.picker.view == GridView {
		ext := filepath.Ext(i.name[1:])
		i.name = i.name[:len(i.name)-len(ext)]
	}

	if up {
		i.name = "(" + lang.X("file.parent", "Parent") + ")"
	}

	i.Refresh()
}

func (i *fileDialogItem) Tapped(*fyne.PointEvent) {
	if i.choose != nil {
		i.choose(i.id)
	}
	now := time.Now()
	if !i.dir && now.Sub(i.lastClick) < fyne.CurrentApp().Driver().DoubleTapDelay() && i.open != nil {
		// It is a double click, so we ask the dialog to open
		i.open()
	}
	i.lastClick = now
}

func (f *fileDialog) newFileItem(location fyne.URI, dir, up bool) *fileDialogItem {
	item := &fileDialogItem{
		picker:   f,
		location: location,
		name:     location.Name(),
		dir:      dir,
	}

	if f.view == GridView {
		ext := filepath.Ext(item.name[1:])
		item.name = item.name[:len(item.name)-len(ext)]
	}

	if up {
		item.name = "(" + lang.X("file.parent", "Parent") + ")"
	}

	item.ExtendBaseWidget(item)
	return item
}

type fileItemRenderer struct {
	item         *fileDialogItem
	fileTextSize float32

	icon    *widget.FileIcon
	text    *widget.Label
	objects []fyne.CanvasObject
}

func (s *fileItemRenderer) Layout(size fyne.Size) {
	if s.item.picker.view == GridView {
		s.icon.Resize(fyne.NewSize(fileIconSize, fileIconSize))
		s.icon.Move(fyne.NewPos((size.Width-fileIconSize)/2, 0))

		s.text.Alignment = fyne.TextAlignCenter
		s.text.Resize(fyne.NewSize(size.Width, s.fileTextSize))
		s.text.Move(fyne.NewPos(0, size.Height-s.fileTextSize))
	} else {
		s.icon.Resize(fyne.NewSize(fileInlineIconSize, fileInlineIconSize))
		s.icon.Move(fyne.NewPos(theme.Padding(), (size.Height-fileInlineIconSize)/2))

		s.text.Alignment = fyne.TextAlignLeading
		textMin := s.text.MinSize()
		s.text.Resize(fyne.NewSize(size.Width, textMin.Height))
		s.text.Move(fyne.NewPos(fileInlineIconSize, (size.Height-textMin.Height)/2))
	}
}

func (s *fileItemRenderer) MinSize() fyne.Size {
	if s.item.picker.view == GridView {
		return fyne.NewSize(fileIconCellWidth, fileIconSize+s.fileTextSize)
	}

	textMin := s.text.MinSize()
	return fyne.NewSize(fileInlineIconSize+textMin.Width+theme.Padding(), textMin.Height)
}

func (s *fileItemRenderer) Refresh() {
	s.fileTextSize = widget.NewLabel("M\nM").MinSize().Height // cache two-line label height

	s.text.SetText(s.item.name)
	s.icon.SetURI(s.item.location)
}

func (s *fileItemRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *fileItemRenderer) Destroy() {
}
