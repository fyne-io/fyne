package dialog

import (
	"fmt"
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
	file       *FileDialog
	fileName   textWidget
	dismiss    *widget.Button
	open       *widget.Button
	breadcrumb *widget.Box
	files      *fyne.Container
	fileScroll *widget.ScrollContainer

	win      *widget.PopUp
	selected *fileDialogItem
	dir      fyne.ListableURI
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

	// StartingLocation allows overriding the default location where the
	// file dialog should "view" when it is first opened.
	StartingLocation fyne.ListableURI
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
			path := storage.Join(f.dir, name)

			// this assumes the file:// type, which is enforced
			// when we set `dir`
			info, err := os.Stat(path.String()[len(path.Scheme())+3:])
			if os.IsNotExist(err) {
				f.win.Hide()
				if f.file.onClosedCallback != nil {
					f.file.onClosedCallback(true)
				}
				callback(storage.SaveFileToURI(path))
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

					callback(storage.SaveFileToURI(path))
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
			callback(storage.OpenFileFromURI(f.selected.path))
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

	favorites, err := f.loadFavorites()
	if err != nil {
		// only generate the Favorites group if we were able to load
		// them successfully
		favorites = []fyne.CanvasObject{}
		fyne.LogError("Unable to load favorites", err)
	}

	favoritesGroup := widget.NewGroup("Favorites", favorites...)
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(header, footer, favoritesGroup, nil),
		favoritesGroup, header, footer, body)

}

func (f *fileDialog) loadFavorites() ([]fyne.CanvasObject, error) {
	osHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	home, err := storage.ListerForURI(storage.NewURI("file://" + osHome))
	if err != nil {
		return nil, err
	}

	documents, err := storage.ListerForURI(storage.Join(home, "Documents"))
	if err != nil {
		return nil, err
	}

	downloads, err := storage.ListerForURI(storage.Join(home, "Downloads"))
	if err != nil {
		return nil, err
	}

	places := []fyne.CanvasObject{

		makeFavoriteButton("Home", theme.HomeIcon(), func() {
			f.setDirectory(home)
		}),
		makeFavoriteButton("Documents", theme.DocumentIcon(), func() {
			f.setDirectory(documents)
		}),
		makeFavoriteButton("Downloads", theme.DownloadIcon(), func() {
			f.setDirectory(downloads)
		}),
	}

	places = append(places, f.loadPlaces()...)
	return places, nil
}

func (f *fileDialog) refreshDir(dir fyne.ListableURI) {
	f.files.Objects = nil

	files, err := dir.List()
	if err != nil {
		fyne.LogError("Unable to read path "+dir.String(), err)
		return
	}

	var icons []fyne.CanvasObject
	parent, err := storage.Parent(dir)
	if err != nil {
		fyne.LogError("Unable to get parent of "+dir.String(), err)
		return
	}
	if parent.String() != dir.String() {
		fi := &fileDialogItem{picker: f, icon: canvas.NewImageFromResource(theme.FolderOpenIcon()),
			name: "(Parent)", path: parent, dir: true}
		fi.ExtendBaseWidget(fi)
		icons = append(icons, fi)
	}
	for _, file := range files {
		if isHidden(file.Name(), dir.Name()) {
			continue
		}

		_, err := storage.ListerForURI(file)
		if err == nil {
			// URI points to a directory
			icons = append(icons, f.newFileItem(file, true))

		} else if f.file.filter == nil || f.file.filter.Matches(file) {
			icons = append(icons, f.newFileItem(file, false))
		}
	}

	f.files.Objects = icons
	f.files.Refresh()
	f.fileScroll.Offset = fyne.NewPos(0, 0)
	f.fileScroll.Refresh()
}

func (f *fileDialog) setDirectory(dir fyne.ListableURI) error {
	f.setSelected(nil)
	f.dir = dir

	f.breadcrumb.Children = nil

	if dir.Scheme() != "file" {
		return fmt.Errorf("Scheme for directory was not file://")
	}

	localdir := dir.String()[len(dir.Scheme())+3:]
	buildDir := filepath.VolumeName(localdir)
	for i, d := range strings.Split(localdir, string(filepath.Separator)) {
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

		newDir, err := storage.ListerForURI(storage.NewURI("file://" + buildDir))
		if err != nil {
			return err
		}
		f.breadcrumb.Append(
			widget.NewButton(d, func() {
				err := f.setDirectory(newDir)
				if err != nil {
					fyne.LogError("Failed to set directory", err)
				}
			}),
		)
	}

	f.refreshDir(dir)

	return nil
}

func (f *fileDialog) setSelected(file *fileDialogItem) {
	if f.selected != nil {
		f.selected.isCurrent = false
		f.selected.Refresh()
	}
	if file != nil && file.isDirectory() {
		lister, err := storage.ListerForURI(file.path)
		if err != nil {
			fyne.LogError("Failed to create lister for URI"+file.path.String(), err)
		}
		f.setDirectory(lister)
		return
	}
	f.selected = file

	if file == nil || file.path.String()[len(file.path.Scheme())+3:] == "" {
		f.fileName.SetText("")
		f.open.Disable()
	} else {
		file.isCurrent = true
		f.fileName.SetText(file.path.Name())
		f.open.Enable()
	}
}

// effectiveStartingDir calculates the directory at which the file dialog should
// open, based on the values of StartingDirectory, CWD, home, and any error
// conditions which occur.
//
// Order of precedence is:
//
// * file.StartingDirectory if non-empty, os.Stat()-able, and uses the file://
//   URI scheme
// * os.Getwd()
// * os.UserHomeDir()
// * "/" (should be filesystem root on all supported platforms)
//
func (f *FileDialog) effectiveStartingDir() fyne.ListableURI {
	var startdir fyne.ListableURI = nil

	if f.StartingLocation != nil {
		startdir = f.StartingLocation
	}

	if startdir != nil {
		if startdir.Scheme() == "file" {
			path := startdir.String()[len(startdir.Scheme())+3:]

			// the starting directory is set explicitly
			if _, err := os.Stat(path); err != nil {
				fyne.LogError("Error with StartingLocation", err)
			} else {
				return startdir
			}
		}

	}

	// Try to get ./
	wd, err := os.Getwd()
	if err == nil {
		lister, err := storage.ListerForURI(storage.NewURI("file://" + wd))
		if err == nil {
			return lister
		}
		fyne.LogError("Could not create lister for working dir", err)
	}

	// Try home dir
	dir, err := os.UserHomeDir()
	if err == nil {
		lister, err := storage.ListerForURI(storage.NewURI("file://" + dir))
		if err == nil {
			return lister
		}
		fyne.LogError("Could not create lister for user home dir", err)
	}
	fyne.LogError("Could not load user home dir", err)

	lister, err := storage.ListerForURI(storage.NewURI("file:///"))
	if err != nil {
		fyne.LogError("could not create lister for /", err)
		return nil
	}
	return lister
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

// ShowFileOpenAt works similarly to ShowFileOpen(), but with a custom starting
// location.
func ShowFileOpenAt(callback func(fyne.URIReadCloser, error), parent fyne.Window, startlocation fyne.ListableURI) {
	dialog := NewFileOpen(callback, parent)
	if fileOpenOSOverride(dialog) {
		return
	}
	dialog.StartingLocation = startlocation
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

// ShowFileSaveAt works simialrly to ShowFileSave(), but with a custom starting
// location.
func ShowFileSaveAt(callback func(fyne.URIWriteCloser, error), parent fyne.Window, startlocation fyne.ListableURI) {
	dialog := NewFileSave(callback, parent)
	if fileSaveOSOverride(dialog) {
		return
	}
	dialog.StartingLocation = startlocation
	dialog.Show()
}

func makeFavoriteButton(title string, icon fyne.Resource, f func()) *widget.Button {
	b := widget.NewButtonWithIcon(title, icon, f)

	b.Alignment = widget.ButtonAlignLeading
	return b
}
