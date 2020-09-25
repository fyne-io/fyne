package widget

import (
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

const ratioDown = 0.45

// FileIcon is an adaption of widget.Icon for showing file information
type FileIcon struct {
	BaseWidget

	URI                fyne.URI
	IconSize, TextSize int

	extension, mimeType, mimeSubType string
	resource                         fyne.Resource
	current                          bool
}

// NewFileIcon takes a filepath and creates an icon with an overlayed label using the detected mimetype and extension
func NewFileIcon(uri fyne.URI) *FileIcon {
	i := &FileIcon{IconSize: 64, TextSize: 24}
	i.ExtendBaseWidget(i)
	i.UpdateURI(uri)
	return i
}

// UpdateURI takes a new URI and updates the icon from that
func (i *FileIcon) UpdateURI(uri fyne.URI) {
	if i.URI == uri {
		return
	}

	i.mimeType, i.mimeSubType = splitMimeType(uri)
	switch i.mimeType {
	case "application":
		i.resource = theme.FileApplicationIcon()
	case "audio":
		i.resource = theme.FileAudioIcon()
	case "image":
		i.resource = theme.FileImageIcon()
	case "text":
		i.resource = theme.FileTextIcon()
	case "video":
		i.resource = theme.FileVideoIcon()
	default:
		i.resource = theme.FileIcon()
	}

	i.extension = trimmedExtension(uri)
	i.URI = uri

	i.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (i *FileIcon) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return i.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (i *FileIcon) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)
	i.propertyLock.RLock()
	defer i.propertyLock.RUnlock()

	img := canvas.NewImageFromResource(i.resource)
	extText := canvas.NewText(i.extension, theme.BackgroundColor())
	extText.Alignment = fyne.TextAlignCenter
	extText.TextSize = theme.TextSize()
	s := &fileIconRenderer{item: i, img: img, ext: extText}
	s.updateObjects()

	return s
}

// SetCurrent makes the icon look like it has been tapped
func (i *FileIcon) SetCurrent(cur bool) {
	i.current = cur
	i.Refresh()
}

// SetSize sets the size for the icon and text
func (i *FileIcon) SetSize(icon, text int) {
	i.IconSize = icon
	i.TextSize = text
	i.Refresh()
}

type fileIconRenderer struct {
	widget.BaseRenderer

	item *FileIcon

	ext *canvas.Text
	img *canvas.Image
}

func (s *fileIconRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (s *fileIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(s.item.IconSize, s.item.IconSize+s.item.TextSize+theme.Padding())
}

func (s *fileIconRenderer) Layout(size fyne.Size) {
	iconAlign := (size.Width - s.item.IconSize) / 2
	s.ext.Resize(fyne.NewSize(s.item.IconSize, s.item.TextSize))
	alignHeight := float64(s.item.IconSize) * ratioDown
	s.ext.Move(fyne.NewPos(iconAlign, int(alignHeight)))

	s.Objects()[0].Resize(size)
}

func (s *fileIconRenderer) Refresh() {
	s.item.propertyLock.RLock()
	s.updateObjects()
	s.item.propertyLock.RUnlock()

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

	//s.Layout(s.img.Size())
	canvas.Refresh(s.item.super())
}

func (s *fileIconRenderer) updateObjects() {
	var objects []fyne.CanvasObject
	img := canvas.NewImageFromResource(s.item.resource)
	s.ext = canvas.NewText(s.item.extension, theme.BackgroundColor())
	s.ext.Alignment = fyne.TextAlignCenter
	s.ext.TextSize = theme.TextSize()

	objects = append(objects, img, s.ext)
	s.SetObjects(objects)
}

func trimmedExtension(uri fyne.URI) string {
	ext := uri.Extension()
	if len(ext) > 5 {
		ext = ext[:5]
	}
	return ext
}

func splitMimeType(uri fyne.URI) (string, string) {
	mimeTypeSplit := strings.Split(uri.MimeType(), "/")
	if len(mimeTypeSplit) <= 1 {
		return "", ""
	}
	return mimeTypeSplit[0], mimeTypeSplit[1]
}
