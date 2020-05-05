package dialog

import (
	"bufio"
	"image/color"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

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

	path                  string
	icon                  fyne.Resource
	name, ext             string
	mimeType, mimeSubType string
}

func (i *fileDialogItem) Tapped(_ *fyne.PointEvent) {
	i.picker.setSelected(i)
	i.Refresh()
}

func (i *fileDialogItem) TappedSecondary(_ *fyne.PointEvent) {
}

func (i *fileDialogItem) CreateRenderer() fyne.WidgetRenderer {
	var img fyne.CanvasObject
	if i.icon == nil {
		mimeSub := i.mimeSubType
		if len(mimeSub) > 5 {
			mimeSub = mimeSub[:5]
		}
		img = NewFileIcon(i.mimeType, mimeSub)
	} else {
		img = canvas.NewImageFromResource(i.icon)
	}
	text := widget.NewLabelWithStyle(i.name, fyne.TextAlignCenter, fyne.TextStyle{})
	text.Wrapping = fyne.TextTruncate

	return &fileItemRenderer{item: i,
		img: img, text: text, objects: []fyne.CanvasObject{img, text}}
}

func (i *fileDialogItem) isDirectory() bool {
	return i.icon == theme.FolderIcon() || i.icon == theme.FolderOpenIcon()
}

func fileParts(path string) (name, ext, mimeType, mimeSubType string) {
	name = filepath.Base(path)
	ext = filepath.Ext(name[1:])
	name = name[:len(name)-len(ext)]

	mimeTypeFull := mime.TypeByExtension(ext)
	if mimeTypeFull == "" {
		mimeTypeFull = "text/plain"
		file, err := os.Open(path)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			if scanner.Scan() && !utf8.Valid(scanner.Bytes()) {
				mimeTypeFull = "application/octet-stream"
			}
		}
	}

	mimeTypeSplit := strings.Split(mimeTypeFull, "/")
	mimeType = mimeTypeSplit[0]
	mimeSubType = mimeTypeSplit[1]

	return
}

func (f *fileDialog) newFileItem(icon fyne.Resource, path string) *fileDialogItem {
	name, ext, mimeType, mimeSubType := fileParts(path)
	if icon == theme.FolderOpenIcon() {
		name = "(Parent)"
		ext = ""
	}

	ret := &fileDialogItem{
		picker:      f,
		icon:        icon,
		name:        name,
		ext:         ext,
		path:        path,
		mimeType:    mimeType,
		mimeSubType: mimeSubType,
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
