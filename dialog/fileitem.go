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

	path      string
	icon      fyne.Resource
	name, ext string
}

func (i *fileDialogItem) Tapped(_ *fyne.PointEvent) {
	i.picker.setSelected(i)
	i.Refresh()
}

func (i *fileDialogItem) TappedSecondary(_ *fyne.PointEvent) {
}

func (i *fileDialogItem) CreateRenderer() fyne.WidgetRenderer {
	img := canvas.NewImageFromResource(i.icon)
	text := widget.NewLabelWithStyle(i.name, fyne.TextAlignCenter, fyne.TextStyle{})
	text.Wrapping = fyne.TextTruncate
	extText := canvas.NewText(i.ext, theme.BackgroundColor())
	extText.Alignment = fyne.TextAlignCenter
	extText.TextSize = theme.TextSize()
	return &fileItemRenderer{item: i,
		img: img, text: text, ext: extText, objects: []fyne.CanvasObject{img, text, extText}}
}

func (i *fileDialogItem) isDirectory() bool {
	return i.icon == theme.FolderIcon() || i.icon == theme.FolderOpenIcon()
}

func fileParts(path string) (name, ext string) {
	name = filepath.Base(path)
	ext = filepath.Ext(name[1:])
	name = name[:len(name)-len(ext)]

	if len(ext) > 1 {
		ext = ext[1:]
	}
	if len(ext) > 5 {
		ext = ext[:5]
	}
	return
}

func (f *fileDialog) newFileItem(icon fyne.Resource, path string) *fileDialogItem {
	name, ext := fileParts(path)
	if icon == theme.FolderOpenIcon() {
		name = "(Parent)"
		ext = ""
	}

	ret := &fileDialogItem{
		picker: f,
		icon:   icon,
		name:   name,
		ext:    ext,
		path:   path,
	}
	ret.ExtendBaseWidget(ret)
	return ret
}

type fileItemRenderer struct {
	item *fileDialogItem

	ext     *canvas.Text
	img     *canvas.Image
	text    *widget.Label
	objects []fyne.CanvasObject
}

func (s fileItemRenderer) Layout(size fyne.Size) {
	iconAlign := (size.Width - fileIconSize) / 2
	s.img.Resize(fyne.NewSize(fileIconSize, fileIconSize))
	s.img.Move(fyne.NewPos(iconAlign, 0))
	s.ext.Resize(fyne.NewSize(fileIconSize, fileTextSize))
	ratioDown := 0.45
	s.ext.Move(fyne.NewPos(iconAlign, int(float64(fileIconSize)*ratioDown)))

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
