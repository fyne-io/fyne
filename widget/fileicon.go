package widget

import (
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
)

const (
	ratioDown     = 0.45
	ratioTextSize = 0.22
)

// FileIcon is an adaption of widget.Icon for showing files and folders
//
// Since: 1.4
type FileIcon struct {
	BaseWidget

	Selected bool
	URI      fyne.URI

	resource  fyne.Resource
	extension string
	cachedURI fyne.URI
}

// NewFileIcon takes a filepath and creates an icon with an overlayed label using the detected mimetype and extension
//
// Since: 1.4
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

	i.URI = uri
	i.cachedURI = nil
	i.resource = i.lookupIcon(i.URI)
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

func (i *FileIcon) lookupIcon(uri fyne.URI) fyne.Resource {
	if i.isDir(uri) {
		return theme.FolderIcon()
	}

	switch splitMimeType(uri) {
	case "application":
		return theme.FileApplicationIcon()
	case "audio":
		return theme.FileAudioIcon()
	case "image":
		return theme.FileImageIcon()
	case "text":
		return theme.FileTextIcon()
	case "video":
		return theme.FileVideoIcon()
	default:
		return theme.FileIcon()
	}
}

func (i *FileIcon) isDir(uri fyne.URI) bool {
	if _, ok := uri.(fyne.ListableURI); ok {
		return true
	}

	if luri, err := storage.ListerForURI(uri); err == nil {
		i.URI = luri // Optimization to avoid having to list it next time
		return true
	}

	return false
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
	size := theme.IconInlineSize()
	return fyne.NewSize(size, size)
}

func (s *fileIconRenderer) Layout(size fyne.Size) {
	isize := fyne.Min(size.Width, size.Height)

	xoff := 0
	yoff := (size.Height - isize) / 2

	if size.Width > size.Height {
		xoff = (size.Width - isize) / 2
	}
	yoff += int(float64(isize) * ratioDown)

	s.ext.TextSize = int(float64(isize) * ratioTextSize)
	s.ext.Resize(fyne.NewSize(isize, s.ext.MinSize().Height))
	s.ext.Move(fyne.NewPos(xoff, yoff))

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
	canvas.Refresh(s.ext)
}

func (s *fileIconRenderer) updateObjects() {
	s.img = canvas.NewImageFromResource(s.file.resource)
	s.ext = canvas.NewText(s.file.extension, theme.BackgroundColor())
	s.img.FillMode = canvas.ImageFillContain
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
