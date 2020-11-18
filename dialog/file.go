package dialog

import (
	"fmt"
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
	desiredSize      *fyne.Size
	// this will be applied to dialog.dir when it's loaded
	startingLocation fyne.ListableURI
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
			location, _ := storage.Child(f.dir, name)

			exists, _ := storage.Exists(location)

			// check if a directory is selected
			_, err := storage.ListerForURI(location)

			if !exists {
				f.win.Hide()
				if f.file.onClosedCallback != nil {
					f.file.onClosedCallback(true)
				}
				callback(storage.SaveFileToURI(location))
				return
			} else if err == nil {
				// a directory has been selected
				ShowInformation("Cannot overwrite",
					"Files cannot replace a directory,\ncheck the file name and try again", f.file.parent)
				return
			}

			ShowConfirm("Overwrite?", "Are you sure you want to overwrite the file\n"+name+"?",
				func(ok bool) {
					if !ok {
						return
					}
					f.win.Hide()

					callback(storage.SaveFileToURI(location))
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
			callback(storage.OpenFileFromURI(f.selected.location))
		} else if f.file.isDirectory() {
			callback := f.file.callback.(func(fyne.ListableURI, error))
			f.win.Hide()
			if f.file.onClosedCallback != nil {
				f.file.onClosedCallback(true)
			}
			callback(f.dir, nil)
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
			} else if f.file.isDirectory() {
				f.file.callback.(func(fyne.ListableURI, error))(nil, nil)
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
	title := label + " File"
	if f.file.isDirectory() {
		title = label + " Folder"
	}
	header := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	favorites := f.loadFavorites()

	favoritesGroup := widget.NewGroup("Favorites", favorites...)
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(header, footer, favoritesGroup, nil),
		favoritesGroup, header, footer, body)

}

func (f *fileDialog) loadFavorites() []fyne.CanvasObject {
	favoriteLocations, err := getFavoriteLocations()
	if err != nil {
		fyne.LogError("Getting favorite locations", err)
	}
	favoriteIcons := getFavoriteIcons()
	favoriteOrder := getFavoriteOrder()

	var places []fyne.CanvasObject
	for _, locName := range favoriteOrder {
		loc, ok := favoriteLocations[locName]
		if !ok {
			continue
		}
		icon := favoriteIcons[locName]
		places = append(places, makeFavoriteButton(locName, icon, func() {
			f.setLocation(loc)
		}))
	}
	places = append(places, f.loadPlaces()...)
	return places
}

func (f *fileDialog) refreshDir(dir fyne.ListableURI) {
	f.files.Objects = nil

	files, err := dir.List()
	if err != nil {
		fyne.LogError("Unable to read ListableURI "+dir.String(), err)
		return
	}

	var icons []fyne.CanvasObject
	parent, err := storage.Parent(dir)
	if err != nil && err != storage.URIRootError {
		fyne.LogError("Unable to get parent of "+dir.String(), err)
		return
	}
	if parent != nil && parent.String() != dir.String() {
		fi := &fileDialogItem{picker: f, name: "(Parent)", location: parent, dir: true}
		fi.ExtendBaseWidget(fi)
		icons = append(icons, fi)
	}

	for _, file := range files {
		if isHidden(file) {
			continue
		}

		listable, err := storage.ListerForURI(file)
		if f.file.isDirectory() && err != nil {
			continue
		} else if err == nil { // URI points to a directory
			icons = append(icons, f.newFileItem(listable, true)) // Pass the listable URI to avoid doing the same check in FileIcon
		} else if f.file.filter == nil || f.file.filter.Matches(file) {
			icons = append(icons, f.newFileItem(file, false))
		}
	}

	f.files.Objects = icons
	f.files.Refresh()
	f.fileScroll.Offset = fyne.NewPos(0, 0)
	f.fileScroll.Refresh()
}

func (f *fileDialog) setLocation(dir fyne.ListableURI) error {
	if dir == nil {
		return fmt.Errorf("failed to open nil directory")
	}

	f.setSelected(nil)
	f.dir = dir

	f.breadcrumb.Children = nil

	localdir := dir.String()[len(dir.Scheme())+3:]

	buildDir := filepath.VolumeName(localdir)
	for i, d := range strings.Split(localdir, "/") {
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

		newDir, err := storage.ListerForURI(storage.NewFileURI(buildDir))
		if err != nil {
			return err
		}
		f.breadcrumb.Append(
			widget.NewButton(d, func() {
				err := f.setLocation(newDir)
				if err != nil {
					fyne.LogError("Failed to set directory", err)
				}
			}),
		)
	}

	if f.file.isDirectory() {
		f.fileName.SetText(dir.Name())
		f.open.Enable()
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
		lister, err := storage.ListerForURI(file.location)
		if err != nil {
			fyne.LogError("Failed to create lister for URI"+file.location.String(), err)
		}
		f.setLocation(lister)
		return
	}
	f.selected = file

	if file == nil || file.location.String()[len(file.location.Scheme())+3:] == "" {
		f.fileName.SetText("")
		f.open.Disable()
	} else {
		file.isCurrent = true
		f.fileName.SetText(file.location.Name())
		f.open.Enable()
	}
}

// effectiveStartingDir calculates the directory at which the file dialog should
// open, based on the values of startingDirectory, CWD, home, and any error
// conditions which occur.
//
// Order of precedence is:
//
// * file.startingDirectory if non-empty, os.Stat()-able, and uses the file://
//   URI scheme
// * os.UserHomeDir()
// * os.Getwd()
// * "/" (should be filesystem root on all supported platforms)
//
func (f *FileDialog) effectiveStartingDir() fyne.ListableURI {
	var startdir fyne.ListableURI = nil

	if f.startingLocation != nil {
		startdir = f.startingLocation
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

	// Try home dir
	dir, err := os.UserHomeDir()
	if err == nil {
		lister, err := storage.ListerForURI(storage.NewFileURI(dir))
		if err == nil {
			return lister
		}
		fyne.LogError("Could not create lister for user home dir", err)
	}
	fyne.LogError("Could not load user home dir", err)

	// Try to get ./
	wd, err := os.Getwd()
	if err == nil {
		lister, err := storage.ListerForURI(storage.NewFileURI(wd))
		if err == nil {
			return lister
		}
		fyne.LogError("Could not create lister for working dir", err)
	}

	lister, err := storage.ListerForURI(storage.NewFileURI("/"))
	if err != nil {
		fyne.LogError("could not create lister for /", err)
		return nil
	}
	return lister
}

func showFile(file *FileDialog) *fileDialog {
	d := &fileDialog{file: file}
	ui := d.makeUI()

	d.setLocation(file.effectiveStartingDir())

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
	if f.desiredSize != nil {
		f.Resize(*f.desiredSize)
		f.desiredSize = nil
	}
}

// Refresh causes this dialog to be updated
func (f *FileDialog) Refresh() {
	f.dialog.win.Refresh()
}

// Resize dialog, call this function after dialog show
func (f *FileDialog) Resize(size fyne.Size) {
	if f.dialog == nil {
		f.desiredSize = &size
		return
	}
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

// SetLocation tells this FileDirectory which location to display.
// This is normally called before the dialog is shown.
//
// Since: 1.4
func (f *FileDialog) SetLocation(u fyne.ListableURI) {
	f.startingLocation = u
	if f.dialog != nil {
		f.dialog.setLocation(u)
	}
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
	if f.isDirectory() {
		fyne.LogError("Cannot set a filter for a folder dialog", nil)
		return
	}
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

func getFavoriteIcons() map[string]fyne.Resource {
	return map[string]fyne.Resource{
		"Home":      theme.HomeIcon(),
		"Documents": theme.DocumentIcon(),
		"Downloads": theme.DownloadIcon(),
	}
}

func getFavoriteOrder() []string {
	return []string{
		"Home",
		"Documents",
		"Downloads",
	}
}
