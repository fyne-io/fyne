package dialog

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
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
	file       *FileDialog
	fileName   textWidget
	dismiss    *widget.Button
	open       *widget.Button
	breadcrumb *widget.Box
	files      *fyne.Container
	fileScroll *widget.ScrollContainer

	win      *widget.PopUp
	selected *fileDialogItem
	dir      string
}

// FileDialog is a dialog containing a file picker for use in opening or saving files.
type FileDialog struct {
	save             bool
	callback         interface{}
	onClosedCallback func(bool)
	filter           storage.FileFilter
	parent           fyne.Window
	dialog           *fileDialog
	dismissText      string
}

// Declare conformity to Dialog interface
var _ Dialog = (*FileDialog)(nil)

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
			if f.file.onClosedCallback != nil {
				f.file.onClosedCallback(false)
			}
			return
		}

		if f.file.save {
			callback := f.file.callback.(func(fyne.URIWriteCloser, error))
			name := f.fileName.(*widget.Entry).Text
			path := filepath.Join(f.dir, name)

			info, err := os.Stat(path)
			if os.IsNotExist(err) {
				f.win.Hide()
				if f.file.onClosedCallback != nil {
					f.file.onClosedCallback(true)
				}
				callback(storage.SaveFileToURI(storage.NewURI("file://" + path)))
				return
			} else if info.IsDir() {
				ShowInformation("Cannot overwrite",
					"Files cannot replace a directory,\ncheck the file name and try again", f.file.parent)
				return
			}

			ShowConfirm("Overwrite?", "Are you sure you want to overwrite the file\n"+name+"?",
				func(ok bool) {
					f.win.Hide()
					if !ok {
						callback(nil, nil)
						return
					}

					callback(storage.SaveFileToURI(storage.NewURI("file://" + path)))
					if f.file.onClosedCallback != nil {
						f.file.onClosedCallback(true)
					}
				}, f.file.parent)
		} else if f.selected != nil {
			callback := f.file.callback.(func(fyne.URIReadCloser, error))
			f.win.Hide()
			if f.file.onClosedCallback != nil {
				f.file.onClosedCallback(true)
			}
			callback(storage.OpenFileFromURI(storage.NewURI("file://" + f.selected.path)))
		}
	})
	f.open.Style = widget.PrimaryButton
	f.open.Disable()
	dismissLabel := "Cancel"
	if f.file.dismissText != "" {
		dismissLabel = f.file.dismissText
	}
	f.dismiss = widget.NewButton(dismissLabel, func() {
		f.win.Hide()
		if f.file.onClosedCallback != nil {
			f.file.onClosedCallback(false)
		}
		if f.file.callback != nil {
			if f.file.save {
				f.file.callback.(func(fyne.URIWriteCloser, error))(nil, nil)
			} else {
				f.file.callback.(func(fyne.URIReadCloser, error))(nil, nil)
			}
		}
	})
	buttons := widget.NewHBox(f.dismiss, f.open)
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
	scrollBread := widget.NewHScrollContainer(f.breadcrumb)
	body := fyne.NewContainerWithLayout(layout.NewBorderLayout(scrollBread, nil, nil, nil),
		scrollBread, f.fileScroll)
	header := widget.NewLabelWithStyle(label+" File", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	favorites := widget.NewGroup("Favorites", f.loadFavorites()...)
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(header, footer, favorites, nil),
		favorites, header, footer, body)
}

func (f *fileDialog) loadFavorites() []fyne.CanvasObject {
	home, _ := os.UserHomeDir()
	places := []fyne.CanvasObject{
		makeFavoriteButton("Home", theme.HomeIcon(), func() {
			f.setDirectory(home)
		}),
		makeFavoriteButton("Documents", theme.DocumentIcon(), func() {
			f.setDirectory(filepath.Join(home, "Documents"))
		}),
		makeFavoriteButton("Downloads", theme.DownloadIcon(), func() {
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
		fi := &fileDialogItem{picker: f, name: "(Parent)", path: filepath.Dir(dir), dir: true}
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
		} else if f.file.filter == nil || f.file.filter.Matches(storage.NewURI("file://"+itemPath)) {
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

// effectiveStartingDir calculates the directory at which the file dialog
// should open, based on the values of  CWD, home, and any error conditions
// which occur.
//
// Order of precedence is:
//
// * os.UserHomeDir()
// * "/" (should be filesystem root on all supported platforms)
func (f *FileDialog) effectiveStartingDir() string {

	// Try home dir
	dir, err := os.UserHomeDir()
	if err == nil {
		return dir
	}
	fyne.LogError("Could not load user home dir", err)

	return "/"
}

func showFile(file *FileDialog) *fileDialog {
	d := &fileDialog{file: file}
	ui := d.makeUI()

	d.setDirectory(file.effectiveStartingDir())

	size := ui.MinSize().Add(fyne.NewSize(fileIconCellWidth*2+theme.Padding()*4,
		(fileIconSize+fileTextSize)+theme.Padding()*4))

	d.win = widget.NewModalPopUp(ui, file.parent.Canvas())
	d.win.Resize(size)

	d.win.Show()
	return d
}

// Show shows the file dialog.
func (f *FileDialog) Show() {
	if f.save {
		if fileSaveOSOverride(f) {
			return
		}
	} else {
		if fileOpenOSOverride(f) {
			return
		}
	}
	if f.dialog != nil {
		f.dialog.win.Show()
		return
	}
	f.dialog = showFile(f)
}

// Resize dialog, call this function after dialog show
func (f *FileDialog) Resize(size fyne.Size) {
	maxSize := f.dialog.win.Size()
	minSize := f.dialog.win.MinSize()
	newWidth := size.Width
	if size.Width > maxSize.Width {
		newWidth = maxSize.Width
	} else if size.Width < minSize.Width {
		newWidth = minSize.Width
	}
	newHeight := size.Height
	if size.Height > maxSize.Height {
		newHeight = maxSize.Height
	} else if size.Height < minSize.Height {
		newHeight = minSize.Height
	}
	f.dialog.win.Resize(fyne.NewSize(newWidth, newHeight))
}

// Hide hides the file dialog.
func (f *FileDialog) Hide() {
	if f.dialog == nil {
		return
	}
	f.dialog.win.Hide()
	if f.onClosedCallback != nil {
		f.onClosedCallback(false)
	}
}

// SetDismissText allows custom text to be set in the confirmation button
func (f *FileDialog) SetDismissText(label string) {
	if f.dialog == nil {
		return
	}
	f.dialog.dismiss.SetText(label)
	widget.Refresh(f.dialog.win)
}

// SetOnClosed sets a callback function that is called when
// the dialog is closed.
func (f *FileDialog) SetOnClosed(closed func()) {
	if f.dialog == nil {
		return
	}
	// If there is already a callback set, remember it and call both.
	originalCallback := f.onClosedCallback

	f.onClosedCallback = func(response bool) {
		closed()
		if originalCallback != nil {
			originalCallback(response)
		}
	}
}

// SetFilter sets a filter for limiting files that can be chosen in the file dialog.
func (f *FileDialog) SetFilter(filter storage.FileFilter) {
	f.filter = filter
	if f.dialog != nil {
		f.dialog.refreshDir(f.dialog.dir)
	}
}

// NewFileOpen creates a file dialog allowing the user to choose a file to open.
// The dialog will appear over the window specified when Show() is called.
func NewFileOpen(callback func(fyne.URIReadCloser, error), parent fyne.Window) *FileDialog {
	dialog := &FileDialog{callback: callback, parent: parent}
	return dialog
}

// NewFileSave creates a file dialog allowing the user to choose a file to save to (new or overwrite).
// If the user chooses an existing file they will be asked if they are sure.
// The dialog will appear over the window specified when Show() is called.
func NewFileSave(callback func(fyne.URIWriteCloser, error), parent fyne.Window) *FileDialog {
	dialog := &FileDialog{callback: callback, parent: parent, save: true}
	return dialog
}

// ShowFileOpen creates and shows a file dialog allowing the user to choose a file to open.
// The dialog will appear over the window specified.
func ShowFileOpen(callback func(fyne.URIReadCloser, error), parent fyne.Window) {
	dialog := NewFileOpen(callback, parent)
	if fileOpenOSOverride(dialog) {
		return
	}
	dialog.Show()
}

// ShowFileSave creates and shows a file dialog allowing the user to choose a file to save to (new or overwrite).
// If the user chooses an existing file they will be asked if they are sure.
// The dialog will appear over the window specified.
func ShowFileSave(callback func(fyne.URIWriteCloser, error), parent fyne.Window) {
	dialog := NewFileSave(callback, parent)
	if fileSaveOSOverride(dialog) {
		return
	}
	dialog.Show()
}

func makeFavoriteButton(title string, icon fyne.Resource, f func()) *widget.Button {
	b := widget.NewButtonWithIcon(title, icon, f)

	b.Alignment = widget.ButtonAlignLeading
	return b
}
