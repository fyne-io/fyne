package dialog

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type textWidget interface {
	fyne.Widget
	SetText(string)
}

type fileDialog struct {
	file       *File
	fileName   textWidget
	open       *widget.Button
	breadcrumb *widget.Box
	files      *fyne.Container
	fileScroll *widget.ScrollContainer

	win      *widget.PopUp
	selected *fileDialogItem
	dir      string
}

// File is a dialog containing a file picker for use in opening or saving files
type File struct {
	save     bool
	callback interface{}
	filter   FileFilter
	parent   fyne.Window
	dialog   *fileDialog
}

// FileFilter is an interface that can be implemented to provide a filter to a file dialog
type FileFilter interface {
	Matches(fyne.URI) bool
}

type extensionFileFilter struct {
	extensions []string
}

type mimeTypeFileFilter struct {
	mimeTypes []string
}

// Matches returns true if a file URI has one of the filtered extensions
func (e *extensionFileFilter) Matches(uri fyne.URI) bool {
	extension := uri.Extension()
	for _, ext := range e.extensions {
		if extension == ext {
			return true
		}
	}
	return false
}

// Matches returns true if a file URI has one of the filtered mimetypes
func (mt *mimeTypeFileFilter) Matches(uri fyne.URI) bool {
	_, mimeType, mimeSubType := mimeTypeGet(uri)
	for _, mimeTypeFull := range mt.mimeTypes {
		mimeTypeSplit := strings.Split(mimeTypeFull, "/")
		if len(mimeTypeSplit) <= 1 {
			continue
		}
		mType := mimeTypeSplit[0]
		mSubType := strings.Split(mimeTypeSplit[1], ";")[0]
		if mType == mimeType {
			if mSubType == mimeSubType || mSubType == "*" {
				return true
			}
		}
	}
	return false
}

func (f *fileDialog) makeUI() fyne.CanvasObject {
	if f.file.save {
		saveName := widget.NewEntry()
		saveName.OnChanged = func(s string) {
			if s == "" {
				f.open.Disable()
			} else {
				f.open.Enable()
			}
		}
		f.fileName = saveName
	} else {
		f.fileName = widget.NewLabel("")
	}

	label := "Open"
	if f.file.save {
		label = "Save"
	}
	f.open = widget.NewButton(label, func() {
		if f.file.callback == nil {
			f.win.Hide()
			return
		}

		if f.file.save {
			callback := f.file.callback.(func(fyne.FileWriteCloser, error))
			name := f.fileName.(*widget.Entry).Text
			path := filepath.Join(f.dir, name)

			info, err := os.Stat(path)
			if os.IsNotExist(err) {
				f.win.Hide()
				callback(storage.SaveFileToURI(storage.NewURI("file://" + path)))
				return
			} else if info.IsDir() {
				ShowInformation("Cannot overwrite",
					"Files cannot replace a directory,\ncheck the file name and try again", f.file.parent)
				return
			}

			ShowConfirm("Overwrite?", "Are you sure you want to overwrite the file\n"+name+"?",
				func(ok bool) {
					if !ok {
						callback(nil, nil)
						return
					}

					callback(storage.SaveFileToURI(storage.NewURI("file://" + path)))
					f.win.Hide()
				}, f.file.parent)
		} else if f.selected != nil {
			callback := f.file.callback.(func(fyne.FileReadCloser, error))
			f.win.Hide()
			callback(storage.OpenFileFromURI(storage.NewURI("file://" + f.selected.path)))
		}
	})
	f.open.Style = widget.PrimaryButton
	f.open.Disable()
	buttons := widget.NewHBox(
		widget.NewButton("Cancel", func() {
			f.win.Hide()
			if f.file.callback != nil {
				if f.file.save {
					f.file.callback.(func(fyne.FileWriteCloser, error))(nil, nil)
				} else {
					f.file.callback.(func(fyne.FileReadCloser, error))(nil, nil)
				}
			}
		}),
		f.open)
	footer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, buttons),
		buttons, widget.NewHScrollContainer(f.fileName))

	f.files = fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(fileIconCellWidth,
		fileIconSize+theme.Padding()+fileTextSize)),
	)
	f.fileScroll = widget.NewScrollContainer(f.files)
	verticalExtra := int(float64(fileIconSize) * 0.25)
	f.fileScroll.SetMinSize(fyne.NewSize(fileIconCellWidth*2+theme.Padding(),
		(fileIconSize+fileTextSize)+theme.Padding()*2+verticalExtra))

	f.breadcrumb = widget.NewHBox()
	scrollBread := widget.NewScrollContainer(f.breadcrumb)
	body := fyne.NewContainerWithLayout(layout.NewBorderLayout(scrollBread, nil, nil, nil),
		scrollBread, f.fileScroll)
	header := widget.NewLabelWithStyle(label+" File", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	favourites := widget.NewGroup("Favourites", f.loadFavourites()...)
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(header, footer, favourites, nil),
		favourites, header, footer, body)
}

func (f *fileDialog) loadFavourites() []fyne.CanvasObject {
	home, _ := os.UserHomeDir()
	places := []fyne.CanvasObject{
		widget.NewButton("Home", func() {
			f.setDirectory(home)
		}),
		widget.NewButton("Documents", func() {
			f.setDirectory(filepath.Join(home, "Documents"))
		}),
		widget.NewButton("Downloads", func() {
			f.setDirectory(filepath.Join(home, "Downloads"))
		}),
	}

	places = append(places, f.loadPlaces()...)
	return places
}

func (f *fileDialog) refreshDir(dir string) {
	f.files.Objects = nil

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fyne.LogError("Unable to read path "+dir, err)
		return
	}

	var icons []fyne.CanvasObject
	parent := filepath.Dir(dir)
	if parent != dir {
		fi := &fileDialogItem{picker: f, icon: canvas.NewImageFromResource(theme.FolderOpenIcon()),
			name: "(Parent)", path: filepath.Dir(dir), dir: true}
		fi.ExtendBaseWidget(fi)
		icons = append(icons, fi)
	}
	for _, file := range files {
		if isHidden(file.Name(), dir) {
			continue
		}
		itemPath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			icons = append(icons, f.newFileItem(itemPath, true))
		} else if f.file.filter == nil || f.file.filter.Matches(storage.NewURI(itemPath)) {
			icons = append(icons, f.newFileItem(itemPath, false))
		}
	}

	f.files.Objects = icons
	f.files.Refresh()
	f.fileScroll.Offset = fyne.NewPos(0, 0)
	f.fileScroll.Refresh()
}

func (f *fileDialog) setDirectory(dir string) {
	f.setSelected(nil)
	f.dir = dir

	f.breadcrumb.Children = nil
	buildDir := filepath.VolumeName(dir)
	for i, d := range strings.Split(dir, string(filepath.Separator)) {
		if d == "" {
			if i > 0 { // what we get if we split "/"
				break
			}
			buildDir = "/"
			d = "/"
		} else if i > 0 {
			buildDir = filepath.Join(buildDir, d)
		} else {
			d = buildDir
			buildDir = d + string(os.PathSeparator)
		}

		newDir := buildDir
		f.breadcrumb.Append(
			widget.NewButton(d, func() {
				f.setDirectory(newDir)
			}),
		)
	}

	f.refreshDir(dir)
}

func (f *fileDialog) setSelected(file *fileDialogItem) {
	if f.selected != nil {
		f.selected.isCurrent = false
		f.selected.Refresh()
	}
	if file != nil && file.isDirectory() {
		f.setDirectory(file.path)
		return
	}
	f.selected = file

	if file == nil || file.path == "" {
		f.fileName.SetText("")
		f.open.Disable()
	} else {
		file.isCurrent = true
		f.fileName.SetText(filepath.Base(file.path))
		f.open.Enable()
	}
}

func showFile(file *File) *fileDialog {
	d := &fileDialog{file: file}
	ui := d.makeUI()
	dir, err := os.UserHomeDir()
	if err != nil {
		fyne.LogError("Could not load user home dir", err)
		dir, _ = os.Getwd() //fallback
	}
	d.setDirectory(dir)

	size := ui.MinSize().Add(fyne.NewSize(fileIconCellWidth*2+theme.Padding()*4,
		(fileIconSize+fileTextSize)+theme.Padding()*4))

	d.win = widget.NewModalPopUp(ui, file.parent.Canvas())
	d.win.Resize(size)

	d.win.Show()
	return d
}

// SetFilter sets a filter for limiting files that can be chosen in the file dialog
func (f *File) SetFilter(filter FileFilter) {
	f.filter = filter
	if f.dialog != nil {
		f.dialog.refreshDir(f.dialog.dir)
	}
}

// Show shows the file dialog
func (f *File) Show() {
	if !f.save {
		if fileOpenOSOverride(f.callback.(func(fyne.FileReadCloser, error)), f.parent) {
			return
		}
		f.dialog = showFile(f)
		return
	}
	if fileSaveOSOverride(f.callback.(func(fyne.FileWriteCloser, error)), f.parent) {
		return
	}
	f.dialog = showFile(f)
}

// NewExtensionFileFilter takes a string slice of extensions with a leading . and creates a filter for the file dialog.
// Example: .jpg, .mp3, .txt, .sh
func NewExtensionFileFilter(extensions []string) FileFilter {
	return &extensionFileFilter{extensions: extensions}
}

// NewMimeTypeFileFilter takes a string slice of mimetypes, including globs, and creates a filter for the file dialog.
// Example: image/*, audio/mp3, text/plain, application/*
func NewMimeTypeFileFilter(mimeTypes []string) FileFilter {
	return &mimeTypeFileFilter{mimeTypes: mimeTypes}
}

// NewFileOpen creates a file dialog allowing the user to choose a file to open.
// The dialog will appear over the window specified when Show() is called.
func NewFileOpen(callback func(fyne.FileReadCloser, error), parent fyne.Window) *File {
	dialog := &File{callback: callback, parent: parent}
	return dialog
}

// NewFileSave creates a file dialog allowing the user to choose a file to save to (new or overwrite).
// If the user chooses an existing file they will be asked if they are sure.
// The dialog will appear over the window specified when Show() is called.
func NewFileSave(callback func(fyne.FileWriteCloser, error), parent fyne.Window) *File {
	dialog := &File{callback: callback, parent: parent, save: true}
	return dialog
}

// ShowFileOpen creates and shows a file dialog allowing the user to choose a file to open.
// The dialog will appear over the window specified.
func ShowFileOpen(callback func(fyne.FileReadCloser, error), parent fyne.Window) *File {
	dialog := NewFileOpen(callback, parent)
	if fileOpenOSOverride(callback, parent) {
		return nil
	}
	dialog.Show()
	return dialog
}

// ShowFileSave creates and shows a file dialog allowing the user to choose a file to save to (new or overwrite).
// If the user chooses an existing file they will be asked if they are sure.
// The dialog will appear over the window specified.
func ShowFileSave(callback func(fyne.FileWriteCloser, error), parent fyne.Window) *File {
	dialog := NewFileSave(callback, parent)
	if fileSaveOSOverride(callback, parent) {
		return nil
	}
	dialog.Show()
	return dialog
}
