package dialog

import (
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type fileIcon struct {
	widget.BaseWidget

	extension, mimeType, mimeSubType string
	resource                         fyne.Resource
	current                          bool
}

type fileIconRenderer struct {
	item *fileIcon

	ext     *canvas.Text
	img     *canvas.Image
	objects []fyne.CanvasObject
}

const ratioDown float64 = 0.45

// NewFileIcon takes a filepath and creates an icon with an overlayed label using the detected mimetype and extension
func NewFileIcon(uri fyne.URI) fyne.CanvasObject {
	mimeType, mimeSubType := splitMimeType(uri)

	var res fyne.Resource
	switch mimeType {
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

	ret := &fileIcon{
		mimeType:    mimeType,
		mimeSubType: mimeSubType,
		extension:   trimmedExtension(uri),
		resource:    res,
	}

	ret.ExtendBaseWidget(ret)
	return ret
}

func (i *fileIcon) CreateRenderer() fyne.WidgetRenderer {
	img := canvas.NewImageFromResource(i.resource)
	extText := canvas.NewText(i.extension, theme.BackgroundColor())
	extText.Alignment = fyne.TextAlignCenter
	extText.TextSize = theme.TextSize()
	return &fileIconRenderer{item: i,
		img: img, ext: extText, objects: []fyne.CanvasObject{img, extText}}
}

func (i *fileIcon) setCurrent(cur bool) {
	i.current = cur
	i.Refresh()
}

func (s fileIconRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (s fileIconRenderer) Destroy() {
}

func (s fileIconRenderer) Layout(size fyne.Size) {
	iconAlign := (size.Width - fileIconSize) / 2
	s.img.Resize(fyne.NewSize(fileIconSize, fileIconSize))
	s.img.Move(fyne.NewPos(iconAlign, 0))
	s.ext.Resize(fyne.NewSize(fileIconSize, fileTextSize))
	alignHeight := float64(fileIconSize) * ratioDown
	s.ext.Move(fyne.NewPos(iconAlign, int(alignHeight)))
}

func (s fileIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(fileIconSize, fileIconSize+fileTextSize+theme.Padding())
}

func (s fileIconRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s fileIconRenderer) Refresh() {
	if s.item.current {
		s.ext.Color = theme.PrimaryColor()
		if _, ok := s.img.Resource.(*theme.InvertedThemedResource); !ok {
			s.img.Resource = theme.NewInvertedThemedResource(s.img.Resource)
		}
	} else {
		s.ext.Color = theme.BackgroundColor()
		if res, ok := s.img.Resource.(*theme.InvertedThemedResource); ok {
			s.img.Resource = res.Original()
		}
	}
	canvas.Refresh(s.item)
}

func trimmedExtension(uri fyne.URI) string {
	ext := uri.Extension()
	if len(ext) > 5 {
		ext = ext[:5]
	}
	return ext
}

func splitMimeType(uri fyne.URI) (string, string) {
	mimeTypeFull := uri.MimeType()
	mimeTypeSplit := strings.Split(mimeTypeFull, "/")
	if len(mimeTypeSplit) <= 1 {
		return "", ""
	}
	return mimeTypeSplit[0], mimeTypeSplit[1]
}
