package dialog

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	fileIconSize      = 64
	fileTextSize      = 24
	fileIconCellWidth = fileIconSize * 1.25
)

type fileDialogItem struct {
	widget.BaseWidget
	picker    *fileDialog
	isCurrent bool

	name     string
	location fyne.URI
	dir      bool
}

func (i *fileDialogItem) Tapped(_ *fyne.PointEvent) {
	i.picker.setSelected(i)
	i.Refresh()
}

func (i *fileDialogItem) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(theme.PrimaryColor())
	background.Hide()
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

func fileName(path fyne.URI) (name string) {
	pathstr := path.String()[len(path.Scheme())+3:]
	name = filepath.Base(pathstr)
	ext := filepath.Ext(name[1:])
	name = name[:len(name)-len(ext)]

	return
}

func (i *fileDialogItem) isDirectory() bool {
	return i.dir
}

func (f *fileDialog) newFileItem(location fyne.URI, dir bool) *fileDialogItem {
	item := &fileDialogItem{
		picker:   f,
		location: location,
		dir:      dir,
	}

	if dir {
		item.name = location.Name()
	} else {
		item.name = fileName(location)
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

	iconAlign := (size.Width - fileIconSize) / 2
	s.icon.Resize(fyne.NewSize(fileIconSize, fileIconSize))
	s.icon.Move(fyne.NewPos(iconAlign, 0))

	s.text.Resize(fyne.NewSize(size.Width, fileTextSize))
	s.text.Move(fyne.NewPos(0, size.Height-fileTextSize-theme.Padding()*2))
}

func (s fileItemRenderer) MinSize() fyne.Size {
	return fyne.NewSize(fileIconSize, fileIconSize+fileTextSize+theme.Padding())
}

func (s fileItemRenderer) Refresh() {
	if s.item.isCurrent {
		s.background.FillColor = theme.PrimaryColor()
		s.background.Show()
	} else {
		s.background.Hide()
	}
	s.background.Refresh()
	s.icon.SetSelected(s.item.isCurrent)
	canvas.Refresh(s.item)
}

func (s fileItemRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s fileItemRenderer) Destroy() {
}
