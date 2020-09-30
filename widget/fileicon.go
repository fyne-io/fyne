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

	IconSize, TextSize int
	Selected           bool
	URI                fyne.URI

	resource  fyne.Resource
	extension string
	cachedURI fyne.URI
}

// NewFileIcon takes a filepath and creates an icon with an overlayed label using the detected mimetype and extension
func NewFileIcon(uri fyne.URI) *FileIcon {
	i := &FileIcon{URI: uri}
	i.ExtendBaseWidget(i)
	return i
}

// SetURI changes the URI and makes the icon reflect a different file
func (i *FileIcon) SetURI(uri fyne.URI) {
	if uri != i.URI {
		i.URI = uri
		i.Refresh()
	}
}

func (i *FileIcon) setURI(uri fyne.URI) {
	if uri == nil {
		i.resource = theme.FileIcon()
		return
	}

	switch splitMimeType(uri) {
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

	i.URI = uri
	i.cachedURI = nil
	i.extension = trimmedExtension(uri)
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

	if i.IconSize == 0 {
		i.IconSize = 64
	}
	if i.TextSize == 0 {
		i.TextSize = 24
	}

	i.setURI(i.URI)

	s := &fileIconRenderer{file: i}
	s.updateObjects()
	i.cachedURI = i.URI

	return s
}

// SetSelected makes the file look like it is selected
func (i *FileIcon) SetSelected(selected bool) {
	i.Selected = selected
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

	file *FileIcon

	ext *canvas.Text
	img *canvas.Image
}

func (s *fileIconRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (s *fileIconRenderer) MinSize() fyne.Size {
	return fyne.NewSize(s.file.IconSize, s.file.IconSize+s.file.TextSize+theme.Padding())
}

func (s *fileIconRenderer) Layout(size fyne.Size) {
	iconAlign := (size.Width - s.file.IconSize) / 2
	s.ext.Resize(fyne.NewSize(s.file.IconSize, s.file.TextSize))
	alignHeight := float64(s.file.IconSize) * ratioDown
	s.ext.Move(fyne.NewPos(iconAlign, int(alignHeight)))

	s.Objects()[0].Resize(size)
}

func (s *fileIconRenderer) Refresh() {
	if s.file.URI != s.file.cachedURI {
		s.file.propertyLock.RLock()
		s.file.setURI(s.file.URI)
		s.updateObjects()
		s.file.cachedURI = s.file.URI
		s.file.propertyLock.RUnlock()
	}

	if s.file.Selected {
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

	canvas.Refresh(s.file.super())
}

func (s *fileIconRenderer) updateObjects() {
	s.img = canvas.NewImageFromResource(s.file.resource)
	s.ext = canvas.NewText(s.file.extension, theme.BackgroundColor())
	s.ext.Alignment = fyne.TextAlignCenter
	s.ext.TextSize = theme.TextSize()

	objects := make([]fyne.CanvasObject, 2) // Length is known, we can pre-allocate
	objects[0], objects[1] = s.img, s.ext

	s.SetObjects(objects)
}

func trimmedExtension(uri fyne.URI) string {
	ext := uri.Extension()
	if len(ext) > 5 {
		ext = ext[:5]
	}
	return ext
}

func splitMimeType(uri fyne.URI) string {
	mimeTypeSplit := strings.Split(uri.MimeType(), "/")
	if len(mimeTypeSplit) <= 1 {
		return ""
	}
	return mimeTypeSplit[0]
}
