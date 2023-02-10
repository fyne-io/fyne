package dialog

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
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
	picker    *fileDialog
	isCurrent bool

	name     string
	location fyne.URI
	dir      bool

	hovered bool
}

func (i *fileDialogItem) MouseIn(*desktop.MouseEvent) {
	i.hovered = true
	i.Refresh()
}

func (i *fileDialogItem) MouseMoved(*desktop.MouseEvent) {
}

func (i *fileDialogItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

func (i *fileDialogItem) Tapped(_ *fyne.PointEvent) {
	i.picker.setSelected(i)
	i.Refresh()
}

func (i *fileDialogItem) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(nil)
	text := widget.NewLabelWithStyle(i.name, fyne.TextAlignCenter, fyne.TextStyle{})
	text.Wrapping = fyne.TextTruncate
	icon := widget.NewFileIcon(i.location)

	return &fileItemRenderer{
		item:       i,
		background: background,
		icon:       icon,
		text:       text,
		objects:    []fyne.CanvasObject{background, icon, text},
	}
}

func (i *fileDialogItem) isDirectory() bool {
	return i.dir
}

func (f *fileDialog) newFileItem(location fyne.URI, dir bool) *fileDialogItem {
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

	item.ExtendBaseWidget(item)
	return item
}

type fileItemRenderer struct {
	item *fileDialogItem

	background *canvas.Rectangle
	icon       *widget.FileIcon
	text       *widget.Label
	objects    []fyne.CanvasObject
}

func (s fileItemRenderer) Layout(size fyne.Size) {
	s.background.Resize(size)

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
	if s.item.isCurrent {
		s.background.FillColor = theme.SelectionColor()
	} else if s.item.hovered {
		s.background.FillColor = theme.HoverColor()
	} else {
		s.background.FillColor = nil
	}

	s.background.Refresh()
	canvas.Refresh(s.item)
}

func (s fileItemRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s fileItemRenderer) Destroy() {
}
