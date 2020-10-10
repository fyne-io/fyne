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

	var icon fyne.CanvasObject
	if i.dir {
		icon = canvas.NewImageFromResource(theme.FolderIcon())
	} else {
		icon = widget.NewFileIcon(i.location)
	}

	return &fileItemRenderer{item: i,
		icon: icon, text: text, objects: []fyne.CanvasObject{icon, text}}
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
	var name string
	if dir {
		name = location.Name()
	} else {
		name = fileName(location)
	}

	ret := &fileDialogItem{
		picker:   f,
		name:     name,
		location: location,
		dir:      dir,
	}
	ret.ExtendBaseWidget(ret)
	return ret
}

type fileItemRenderer struct {
	item *fileDialogItem

	icon    fyne.CanvasObject
	text    *widget.Label
	objects []fyne.CanvasObject
}

func (s fileItemRenderer) Layout(size fyne.Size) {
	iconAlign := (size.Width - fileIconSize) / 2
	s.icon.Resize(fyne.NewSize(fileIconSize, fileIconSize))
	s.icon.Move(fyne.NewPos(iconAlign, 0))

	s.text.Resize(fyne.NewSize(size.Width, fileTextSize))
	s.text.Move(fyne.NewPos(0, fileIconSize+theme.Padding()))
}

func (s fileItemRenderer) MinSize() fyne.Size {
	return fyne.NewSize(fileIconSize, fileIconSize+fileTextSize+theme.Padding())
}

func (s fileItemRenderer) Refresh() {
	if s.item.isCurrent {
		if i, ok := s.icon.(*canvas.Image); ok {
			if _, ok := i.Resource.(*theme.InvertedThemedResource); !ok {
				i.Resource = theme.NewInvertedThemedResource(i.Resource)
			}
		} else if i, ok := s.icon.(*widget.FileIcon); ok {
			i.SetSelected(true)
		}
	} else {
		if i, ok := s.icon.(*canvas.Image); ok {
			if res, ok := i.Resource.(*theme.InvertedThemedResource); ok {
				i.Resource = res.Original()
			}
		} else if i, ok := s.icon.(*widget.FileIcon); ok {
			i.SetSelected(false)
		}
	}
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
