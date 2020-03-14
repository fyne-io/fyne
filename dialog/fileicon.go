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

type fileDialogIcon struct {
	widget.BaseWidget
	picker    *fileDialog
	isCurrent bool

	path      string
	icon      fyne.Resource
	name, ext string
}

func (i *fileDialogIcon) Tapped(_ *fyne.PointEvent) {
	i.picker.setSelected(i)
	i.Refresh()
}

func (i *fileDialogIcon) TappedSecondary(_ *fyne.PointEvent) {
}

func (i *fileDialogIcon) CreateRenderer() fyne.WidgetRenderer {
	img := canvas.NewImageFromResource(i.icon)
	text := widget.NewLabelWithStyle(i.name, fyne.TextAlignCenter, fyne.TextStyle{})
	extText := canvas.NewText(i.ext, theme.BackgroundColor())
	extText.Alignment = fyne.TextAlignCenter
	extText.TextSize = theme.TextSize()
	return &fileIconRenderer{icon: i,
		img: img, text: text, ext: extText, objects: []fyne.CanvasObject{img, text, extText}}
}

func (i *fileDialogIcon) isDirectory() bool {
	return i.icon == theme.FolderIcon() || i.icon == theme.FolderOpenIcon()
}

func fileParts(path string) (name, ext string) {
	name = filepath.Base(path)
	ext = filepath.Ext(name[1:])
	name = name[:len(name)-len(ext)]

	if len(ext) > 1 {
		ext = ext[1:]
	}
	return
}

func (f *fileDialog) newFileIcon(icon fyne.Resource, path string) *fileDialogIcon {
	name, ext := fileParts(path)
	if icon == theme.FolderOpenIcon() {
		name = "(Parent)"
		ext = ""
	}

	ret := &fileDialogIcon{
		picker: f,
		icon:   icon,
		name:   name,
		ext:    ext,
		path:   path,
	}
	ret.ExtendBaseWidget(ret)
	return ret
}

type fileIconRenderer struct {
	icon *fileDialogIcon

	ext     *canvas.Text
	img     *canvas.Image
	text    *widget.Label
	objects []fyne.CanvasObject
}

func (s fileIconRenderer) Layout(size fyne.Size) {
	iconAlign := (size.Width - fileIconSize) / 2
	s.img.Resize(fyne.NewSize(fileIconSize, fileIconSize))
	s.img.Move(fyne.NewPos(iconAlign, 0))
	s.ext.Resize(fyne.NewSize(fileIconSize, fileTextSize))
	ratioDown := 0.45
	s.ext.Move(fyne.NewPos(iconAlign, int(float64(fileIconSize)*ratioDown)))

	s.text.Resize(fyne.NewSize(size.Width, fileTextSize))
	s.text.Move(fyne.NewPos(0, fileIconSize+theme.Padding()))
}

func (s fileIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(fileIconSize, fileIconSize+fileTextSize+theme.Padding())
}

func (s fileIconRenderer) Refresh() {
	canvas.Refresh(s.icon)
}

func (s fileIconRenderer) BackgroundColor() color.Color {
	if s.icon.isCurrent {
		return theme.PrimaryColor()
	}
	return theme.BackgroundColor()
}

func (s fileIconRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s fileIconRenderer) Destroy() {
}
