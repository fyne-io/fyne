package dialog

import (
	"image/color"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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

	icon     fyne.CanvasObject
	name     string
	location fyne.URI
	dir      bool
}

func (i *fileDialogItem) Tapped(_ *fyne.PointEvent) {
	i.picker.setSelected(i)
	i.Refresh()
}

func (i *fileDialogItem) TappedSecondary(_ *fyne.PointEvent) {
}

func (i *fileDialogItem) CreateRenderer() fyne.WidgetRenderer {
	text := widget.NewLabelWithStyle(i.name, fyne.TextAlignCenter, fyne.TextStyle{})
	text.Wrapping = fyne.TextTruncate

	return &fileItemRenderer{item: i,
		img: i.icon, text: text, objects: []fyne.CanvasObject{i.icon, text}}
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
	var icon fyne.CanvasObject
	var name string
	if dir {
		icon = canvas.NewImageFromResource(theme.FolderIcon())
		name = location.Name()
	} else {
		icon = NewFileIcon(location)
		name = fileName(location)
	}

	ret := &fileDialogItem{
		picker:   f,
		icon:     icon,
		name:     name,
		location: location,
		dir:      dir,
	}
	ret.ExtendBaseWidget(ret)
	return ret
}

type fileItemRenderer struct {
	item *fileDialogItem

	img     fyne.CanvasObject
	text    *widget.Label
	objects []fyne.CanvasObject
}

func (s fileItemRenderer) Layout(size fyne.Size) {
	iconAlign := (size.Width - fileIconSize) / 2
	s.img.Resize(fyne.NewSize(fileIconSize, fileIconSize))
	s.img.Move(fyne.NewPos(iconAlign, 0))

	s.text.Resize(fyne.NewSize(size.Width, fileTextSize))
	s.text.Move(fyne.NewPos(0, fileIconSize+theme.Padding()))
}

func (s fileItemRenderer) MinSize() fyne.Size {
	return fyne.NewSize(fileIconSize, fileIconSize+fileTextSize+theme.Padding())
}

func (s fileItemRenderer) Refresh() {
	canvas.Refresh(s.item)
}

func (s fileItemRenderer) BackgroundColor() color.Color {
	if s.item.isCurrent {
		return theme.PrimaryColor()
	}
	return theme.BackgroundColor()
}

func (s fileItemRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s fileItemRenderer) Destroy() {
}
