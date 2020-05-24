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
	fileName   textWidget
	open       *widget.Button
	breadcrumb *widget.Box
	files      *fyne.Container
	fileScroll *widget.ScrollContainer
	parent     fyne.Window

	win      *widget.PopUp
	selected *fileDialogItem
	callback interface{}
	filter   FileFilter
	dir      string
	save     bool
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
		mSubType := mimeTypeSplit[1]
		if mType == mimeType {
			if mSubType == mimeSubType || mSubType == "*" {
				return true
			}
		}
	}
	return false
}

func (f *fileDialog) makeUI() fyne.CanvasObject {
	if f.save {
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
	if f.save {
		label = "Save"
	}
	f.open = widget.NewButton(label, func() {
		if f.callback == nil {
			f.win.Hide()
			return
		}

		if f.save {
			callback := f.callback.(func(fyne.FileWriteCloser, error))
			name := f.fileName.(*widget.Entry).Text
			path := filepath.Join(f.dir, name)

			info, err := os.Stat(path)
			if os.IsNotExist(err) {
				f.win.Hide()
				callback(storage.SaveFileToURI(storage.NewURI("file://" + path)))
				return
			} else if info.IsDir() {
				ShowInformation("Cannot overwrite",
					"Files cannot replace a directory,\ncheck the file name and try again", f.parent)
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
				}, f.parent)
		} else if f.selected != nil {
			callback := f.callback.(func(fyne.FileReadCloser, error))
			f.win.Hide()
			callback(storage.OpenFileFromURI(storage.NewURI("file://" + f.selected.path)))
		}
	})
	f.open.Style = widget.PrimaryButton
	f.open.Disable()
	buttons := widget.NewHBox(
		widget.NewButton("Cancel", func() {
			f.win.Hide()
			if f.callback != nil {
				if f.save {
					f.callback.(func(fyne.FileWriteCloser, error))(nil, nil)
				} else {
					f.callback.(func(fyne.FileReadCloser, error))(nil, nil)
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
		} else {
			if f.filter != nil && !f.filter.Matches(storage.NewURI(itemPath)) {
				continue
			}
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

func showFileDialog(save bool, callback interface{}, parent fyne.Window, filter FileFilter) {
	d := &fileDialog{callback: callback, save: save, parent: parent, filter: filter}
	ui := d.makeUI()
	dir, err := os.UserHomeDir()
	if err != nil {
		fyne.LogError("Could not load user home dir", err)
		dir, _ = os.Getwd() //fallback
	}
	d.setDirectory(dir)

	size := ui.MinSize().Add(fyne.NewSize(fileIconCellWidth*2+theme.Padding()*4,
		(fileIconSize+fileTextSize)+theme.Padding()*4))

	d.win = widget.NewModalPopUp(ui, parent.Canvas())
	d.win.Resize(size)

	d.win.Show()
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

// ShowFileOpen shows a file dialog allowing the user to choose a file to open.
// The dialog will appear over the window specified.
func ShowFileOpen(callback func(fyne.FileReadCloser, error), parent fyne.Window) {
	if fileOpenOSOverride(callback, parent) {
		return
	}
	showFileDialog(false, callback, parent, nil)
}

// ShowFileOpenWithFilter shows a file dialog with a filter limiting displayed files to choose a file to open.
// The dialog will appear over the window specified.
func ShowFileOpenWithFilter(callback func(fyne.FileReadCloser, error), parent fyne.Window, filter FileFilter) {
	if fileOpenOSOverride(callback, parent) {
		return
	}
	showFileDialog(false, callback, parent, filter)
}

// ShowFileSave shows a file dialog allowing the user to choose a file to save to (new or overwrite).
// If the user chooses an existing file they will be asked if they are sure.
// The dialog will appear over the window specified.
func ShowFileSave(callback func(fyne.FileWriteCloser, error), parent fyne.Window) {
	if fileSaveOSOverride(callback, parent) {
		return
	}
	showFileDialog(true, callback, parent, nil)
}

// ShowFileSaveWithFilter shows a file dialog with a filter limiting displayed files to choose a file to save to (new or overwrite).
// If the user chooses an existing file they will be asked if they are sure.
// The dialog will appear over the window specified.
func ShowFileSaveWithFilter(callback func(fyne.FileWriteCloser, error), parent fyne.Window, filter FileFilter) {
	if fileSaveOSOverride(callback, parent) {
		return
	}
	showFileDialog(true, callback, parent, filter)
}
