package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/repository/mime"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
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

	// Deprecated: Selection is now handled externally.
	Selected bool
	URI      fyne.URI

	resource  fyne.Resource
	extension string
}

// NewFileIcon takes a filepath and creates an icon with an overlaid label using the detected mimetype and extension
//
// Since: 1.4
func NewFileIcon(uri fyne.URI) *FileIcon {
	i := &FileIcon{URI: uri}
	i.ExtendBaseWidget(i)
	return i
}

// SetURI changes the URI and makes the icon reflect a different file
func (i *FileIcon) SetURI(uri fyne.URI) {
	i.URI = uri
	i.Refresh()
}

func (i *FileIcon) setURI(uri fyne.URI) {
	if uri == nil {
		i.resource = theme.FileIcon()
		return
	}

	i.URI = uri
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
	i.propertyLock.Lock()
	i.setURI(i.URI)
	i.propertyLock.Unlock()

	i.propertyLock.RLock()
	defer i.propertyLock.RUnlock()

	// TODO remove background when `SetSelected` is gone.
	background := canvas.NewRectangle(theme.SelectionColor())
	background.Hide()

	s := &fileIconRenderer{file: i, background: background}
	s.img = canvas.NewImageFromResource(s.file.resource)
	s.img.FillMode = canvas.ImageFillContain
	s.ext = canvas.NewText(s.file.extension, theme.BackgroundColor())
	s.ext.Alignment = fyne.TextAlignCenter

	s.SetObjects([]fyne.CanvasObject{s.background, s.img, s.ext})

	return s
}

// SetSelected makes the file look like it is selected.
//
// Deprecated: Selection is now handled externally.
func (i *FileIcon) SetSelected(selected bool) {
	i.Selected = selected
	i.Refresh()
}

func (i *FileIcon) lookupIcon(uri fyne.URI) fyne.Resource {
	if i.isDir(uri) {
		return theme.FolderIcon()
	}

	mainMimeType, _ := mime.Split(uri.MimeType())
	switch mainMimeType {
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

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		return false
	}

	i.URI = listable // Avoid having to call storage.ListerForURI(uri) the next time.
	return true
}

type fileIconRenderer struct {
	widget.BaseRenderer

	file *FileIcon

	background *canvas.Rectangle
	ext        *canvas.Text
	img        *canvas.Image
}

func (s *fileIconRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(theme.IconInlineSize())
}

func (s *fileIconRenderer) Layout(size fyne.Size) {
	isize := fyne.Min(size.Width, size.Height)

	xoff := float32(0)
	yoff := (size.Height - isize) / 2

	if size.Width > size.Height {
		xoff = (size.Width - isize) / 2
	}
	yoff += isize * ratioDown

	oldSize := s.ext.TextSize
	s.ext.TextSize = float32(int(isize * ratioTextSize))
	s.ext.Resize(fyne.NewSize(isize, s.ext.MinSize().Height))
	s.ext.Move(fyne.NewPos(xoff, yoff))
	if oldSize != s.ext.TextSize {
		s.ext.Refresh()
	}

	s.Objects()[0].Resize(size)
	s.Objects()[1].Resize(size)
}

func (s *fileIconRenderer) Refresh() {
	s.file.propertyLock.Lock()
	s.file.setURI(s.file.URI)
	s.file.propertyLock.Unlock()

	s.file.propertyLock.RLock()
	s.img.Resource = s.file.resource
	s.ext.Text = s.file.extension
	s.file.propertyLock.RUnlock()

	if s.file.Selected {
		s.background.Show()
		s.ext.Color = theme.SelectionColor()
		if _, ok := s.img.Resource.(*theme.InvertedThemedResource); !ok {
			s.img.Resource = theme.NewInvertedThemedResource(s.img.Resource)
		}
	} else {
		s.background.Hide()
		s.ext.Color = theme.BackgroundColor()
		if res, ok := s.img.Resource.(*theme.InvertedThemedResource); ok {
			s.img.Resource = res.Original()
		}
	}

	canvas.Refresh(s.file.super())
	canvas.Refresh(s.ext)
}

func trimmedExtension(uri fyne.URI) string {
	ext := uri.Extension()
	if len(ext) > 5 {
		ext = ext[:5]
	}
	return ext
}
