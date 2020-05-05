package dialog

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type fileIcon struct {
	widget.BaseWidget

	mimeType, mimeSubType string
}

func (i *fileIcon) CreateRenderer() fyne.WidgetRenderer {
	var res fyne.Resource
	switch i.mimeType {
	case "application":
		res = theme.FileApplicationIcon()
	case "audio":
		res = theme.FileAudioIcon()
	case "image":
		res = theme.FileImageIcon()
	case "text":
		res = theme.FileTextIcon()
	case "video":
		res = theme.FileVideoIcon()
	default:
		res = theme.FileIcon()
	}
	img := canvas.NewImageFromResource(res)
	extText := canvas.NewText(i.mimeSubType, theme.BackgroundColor())
	extText.Alignment = fyne.TextAlignCenter
	extText.TextSize = theme.TextSize()
	return &fileIconRenderer{item: i,
		img: img, ext: extText, objects: []fyne.CanvasObject{img, extText}}
}

// NewFileIcon takes mimetype information and creates an icon that can be used to represent files
func NewFileIcon(mimeType string, mimeSubType string) fyne.CanvasObject {

	ret := &fileIcon{
		mimeType:    mimeType,
		mimeSubType: mimeSubType,
	}

	ret.ExtendBaseWidget(ret)
	return ret
}

type fileIconRenderer struct {
	item *fileIcon

	ext     *canvas.Text
	img     *canvas.Image
	objects []fyne.CanvasObject
}

func (s fileIconRenderer) Layout(size fyne.Size) {
	iconAlign := (size.Width - fileIconSize) / 2
	s.img.Resize(fyne.NewSize(fileIconSize, fileIconSize))
	s.img.Move(fyne.NewPos(iconAlign, 0))
	s.ext.Resize(fyne.NewSize(fileIconSize, fileTextSize))
	ratioDown := 0.45
	s.ext.Move(fyne.NewPos(iconAlign, int(float64(fileIconSize)*ratioDown)))
}

func (s fileIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(fileIconSize, fileIconSize+fileTextSize+theme.Padding())
}

func (s fileIconRenderer) Refresh() {
	canvas.Refresh(s.item)
}

func (s fileIconRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (s fileIconRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s fileIconRenderer) Destroy() {
}
