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

type fileIcon struct {
	widget.BaseWidget

	extension, mimeType, mimeSubType string
	resource fyne.Resource
}

const ratioDown float64 = 0.45

func (i *fileIcon) CreateRenderer() fyne.WidgetRenderer {
	img := canvas.NewImageFromResource(i.resource)
	extText := canvas.NewText(i.extension, theme.BackgroundColor())
	extText.Alignment = fyne.TextAlignCenter
	extText.TextSize = theme.TextSize()
	return &fileIconRenderer{item: i,
		img: img, ext: extText, objects: []fyne.CanvasObject{img, extText}}
}

func mimeTypeGet(path string) (ext, mimeType, mimeSubType string) {
	ext = filepath.Ext(path[1:])
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
	mimeSubTypeSplit := strings.Split(mimeTypeSplit[1], ";")
	mimeSubType = mimeSubTypeSplit[0]

	if len(ext) > 5 {
		ext = ext[:5]
	}
	return
}

// NewFileIcon takes a filepath and creates an icon with an overlayed label using the detected mimetype and extension
func NewFileIcon(path string) fyne.CanvasObject {
	ext, mimeType, mimeSubType := mimeTypeGet(path)

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
		extension:   ext,
		resource:    res,
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
	alignHeight := float64(fileIconSize)*ratioDown
	s.ext.Move(fyne.NewPos(iconAlign, int(alignHeight)))
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
