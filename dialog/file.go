package dialog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	breadcrumb *fyne.Container
	files      *fyne.Container
	fileScroll *container.Scroll
	showHidden bool

	win      *widget.PopUp
	selected *fileDialogItem
	dir      fyne.ListableURI
	// this will be the initial filename in a FileDialog in save mode
	initialFileName string
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
	desiredSize      fyne.Size
	// this will be applied to dialog.dir when it's loaded
	startingLocation fyne.ListableURI
	// this will be the initial filename in a FileDialog in save mode
	initialFileName string
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
		saveName.SetPlaceHolder("enter filename")
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
			listable, err := storage.CanList(location)

			if !exists {
				f.win.Hide()
				if f.file.onClosedCallback != nil {
					f.file.onClosedCallback(true)
				}
				callback(storage.Writer(location))
				return
			} else if err == nil && listable {
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

					callback(storage.Writer(location))
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
			callback(storage.Reader(f.selected.location))
		} else if f.file.isDirectory() {
			callback := f.file.callback.(func(fyne.ListableURI, error))
			f.win.Hide()
			if f.file.onClosedCallback != nil {
				f.file.onClosedCallback(true)
			}
			callback(f.dir, nil)
		}
	})
	f.open.Importance = widget.HighImportance
	f.open.Disable()
	if f.file.save {
		f.fileName.SetText(f.initialFileName)
	}
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
	buttons := container.NewHBox(f.dismiss, f.open)

	footer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, buttons),
		buttons, container.NewHScroll(f.fileName))

	f.files = fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(fileIconCellWidth,
		fileIconSize+theme.Padding()+fileTextSize)),
	)
	f.fileScroll = container.NewScroll(f.files)
	verticalExtra := float32(float64(fileIconSize) * 0.25)
	f.fileScroll.SetMinSize(fyne.NewSize(fileIconCellWidth*2+theme.Padding(),
		(fileIconSize+fileTextSize)+theme.Padding()*2+verticalExtra))

	f.breadcrumb = container.NewHBox()
	scrollBread := container.NewHScroll(f.breadcrumb)
	body := fyne.NewContainerWithLayout(layout.NewBorderLayout(scrollBread, nil, nil, nil),
		scrollBread, f.fileScroll)
	title := label + " File"
	if f.file.isDirectory() {
		title = label + " Folder"
	}
	header := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	favorites := f.loadFavorites()

	favoritesGroup := container.NewVScroll(widget.NewCard("Favorites", "",
		container.NewVBox(favorites...)))
	var optionsButton *widget.Button
	optionsButton = widget.NewButtonWithIcon("Options", theme.SettingsIcon(), func() {
		f.optionsMenu(fyne.CurrentApp().Driver().AbsolutePositionForObject(optionsButton), optionsButton.Size())
	})

	left := container.NewBorder(nil, optionsButton, nil, nil, favoritesGroup)

	return container.NewBorder(header, footer, left, nil, body)
}

func (f *fileDialog) optionsMenu(position fyne.Position, buttonSize fyne.Size) {
	hiddenFiles := widget.NewCheck("Show Hidden Files", func(changed bool) {
		f.showHidden = changed
		f.refreshDir(f.dir)
	})
	hiddenFiles.SetChecked(f.showHidden)
	content := container.NewVBox(hiddenFiles)

	p := position.Add(buttonSize)
	pos := fyne.NewPos(p.X, p.Y-content.MinSize().Height-theme.Padding()*2)
	widget.ShowPopUpAtPosition(content, f.win.Canvas, pos)
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
	if err != nil && err != repository.ErrURIRoot {
		fyne.LogError("Unable to get parent of "+dir.String(), err)
		return
	}
	if parent != nil && parent.String() != dir.String() {
		fi := &fileDialogItem{picker: f, name: "(Parent)", location: parent, dir: true}
		fi.ExtendBaseWidget(fi)
		icons = append(icons, fi)
	}

	for _, file := range files {
		if !f.showHidden && isHidden(file) {
			continue
		}

		listable, err := storage.CanList(file)
		if f.file.isDirectory() && err != nil {
			continue
		} else if err == nil && listable { // URI points to a directory
			icons = append(icons, f.newFileItem(file, true)) // Pass the listable URI to avoid doing the same check in FileIcon
		} else if f.file.filter == nil || f.file.filter.Matches(file) {
			icons = append(icons, f.newFileItem(file, false))
		}
	}

	f.files.Objects = icons
	f.files.Refresh()
	f.fileScroll.Offset = fyne.NewPos(0, 0)
	f.fileScroll.Refresh()
}

func (f *fileDialog) setLocation(dir fyne.URI) error {
	if dir == nil {
		return fmt.Errorf("failed to open nil directory")
	}
	list, err := storage.ListerForURI(dir)
	if err != nil {
		return err
	}

	f.setSelected(nil)
	f.dir = list

	f.breadcrumb.Objects = nil

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

		newDir := storage.NewFileURI(buildDir)
		isDir, err := storage.CanList(newDir)
		if err != nil {
			return err
		}

		if !isDir {
			return errors.New("location was not a listable URI")
		}
		f.breadcrumb.Add(
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
	f.refreshDir(list)

	return nil
}

func (f *fileDialog) setSelected(file *fileDialogItem) {
	if f.selected != nil {
		f.selected.isCurrent = false
		f.selected.Refresh()
	}
	if file != nil && file.isDirectory() {
		listable, err := storage.CanList(file.location)
		if err != nil || !listable {
			fyne.LogError("Failed to create lister for URI"+file.location.String(), err)
		}
		f.setLocation(file.location)
		return
	}
	f.selected = file

	if file == nil || file.location.String()[len(file.location.Scheme())+3:] == "" {
		// keep user input while navigating
		// in a FileSave dialog
		if !f.file.save {
			f.fileName.SetText("")
			f.open.Disable()
		}
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
	d := &fileDialog{file: file, initialFileName: file.initialFileName}
	ui := d.makeUI()

	d.setLocation(file.effectiveStartingDir())

	size := ui.MinSize().Add(fyne.NewSize(fileIconCellWidth*2+theme.Padding()*4,
		(fileIconSize+fileTextSize)+theme.Padding()*4))

	d.win = widget.NewModalPopUp(ui, file.parent.Canvas())
	d.win.Resize(size)

	d.win.Show()
	return d
}

// MinSize returns the size that this dialog should not shrink below
//
// Since: 2.1
func (f *FileDialog) MinSize() fyne.Size {
	return f.dialog.win.MinSize()
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
	if !f.desiredSize.IsZero() {
		f.Resize(f.desiredSize)
	}
}

// Refresh causes this dialog to be updated
func (f *FileDialog) Refresh() {
	f.dialog.win.Refresh()
}

// Resize dialog to the requested size, if there is sufficient space.
// If the parent window is not large enough then the size will be reduced to fit.
func (f *FileDialog) Resize(size fyne.Size) {
	f.desiredSize = size
	if f.dialog == nil {
		return
	}
	f.dialog.win.Resize(size)
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
	f.dialog.win.Refresh()
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

// SetFileName sets the filename in a FileDialog in save mode.
// This is normally called before the dialog is shown.
func (f *FileDialog) SetFileName(fileName string) {
	if f.save {
		f.initialFileName = fileName
		//Update entry if fileDialog has already been created
		if f.dialog != nil {
			f.dialog.fileName.SetText(fileName)
		}
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
