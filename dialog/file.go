package dialog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ViewLayout can be passed to SetView() to set the view of
// a FileDialog
//
// Since: 2.5
type ViewLayout int

const (
	defaultView ViewLayout = iota
	ListView
	GridView
)

const viewLayoutKey = "fyne:fileDialogViewLayout"

type textWidget interface {
	fyne.Widget
	SetText(string)
}

type favoriteItem struct {
	locName string
	locIcon fyne.Resource
	loc     fyne.URI
}

type fileDialogPanel interface {
	fyne.Widget

	Unselect(int)
}

type fileDialog struct {
	file             *FileDialog
	fileName         textWidget
	dismiss          *widget.Button
	open             *widget.Button
	breadcrumb       *fyne.Container
	breadcrumbScroll *container.Scroll
	files            fileDialogPanel
	filesScroll      *container.Scroll
	favorites        []favoriteItem
	favoritesList    *widget.List
	showHidden       bool

	view ViewLayout

	data     []fyne.URI
	dataLock sync.RWMutex

	win        *widget.PopUp
	selected   fyne.URI
	selectedID int
	dir        fyne.ListableURI
	// this will be the initial filename in a FileDialog in save mode
	initialFileName string

	toggleViewButton *widget.Button
}

// FileDialog is a dialog containing a file picker for use in opening or saving files.
type FileDialog struct {
	callback         any
	onClosedCallback func(bool)
	parent           fyne.Window
	dialog           *fileDialog

	confirmText, dismissText string
	desiredSize              fyne.Size
	filter                   storage.FileFilter
	save                     bool
	// this will be applied to dialog.dir when it's loaded
	startingLocation fyne.ListableURI
	// this will be the initial filename in a FileDialog in save mode
	initialFileName string
	// this will be the initial view in a FileDialog
	initialView ViewLayout
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
		saveName.SetPlaceHolder(lang.L("Enter filename"))
		saveName.OnSubmitted = func(s string) {
			f.open.OnTapped()
		}
		f.fileName = saveName
	} else {
		f.fileName = widget.NewLabel("")
	}

	label := lang.L("Open")
	if f.file.save {
		label = lang.L("Save")
	}
	if f.file.confirmText != "" {
		label = f.file.confirmText
	}
	f.open = f.makeOpenButton(label)

	if f.file.save {
		f.fileName.SetText(f.initialFileName)
	}

	dismissLabel := lang.L("Cancel")
	if f.file.dismissText != "" {
		dismissLabel = f.file.dismissText
	}
	f.dismiss = f.makeDismissButton(dismissLabel)

	buttons := container.NewGridWithRows(1, f.dismiss, f.open)

	f.filesScroll = container.NewScroll(nil) // filesScroll's content will be set by setView function.
	verticalExtra := float32(float64(fileIconSize) * 0.25)
	itemMin := f.newFileItem(storage.NewFileURI("filename.txt"), false, false).MinSize()
	f.filesScroll.SetMinSize(itemMin.AddWidthHeight(itemMin.Width+theme.Padding()*3, verticalExtra))

	f.breadcrumb = container.NewHBox()
	f.breadcrumbScroll = container.NewHScroll(container.NewPadded(f.breadcrumb))
	title := label + " " + lang.L("File")
	if f.file.isDirectory() {
		title = label + " " + lang.L("Folder")
	}

	view := ViewLayout(fyne.CurrentApp().Preferences().Int(viewLayoutKey))

	// handle invalid values
	if view != GridView && view != ListView {
		view = defaultView
	}

	if view == defaultView {
		// set GridView as default
		view = GridView

		if f.file.initialView != defaultView {
			view = f.file.initialView
		}
	}

	// icon of button is set in subsequent setView() call
	f.toggleViewButton = widget.NewButtonWithIcon("", nil, func() {
		if f.view == GridView {
			f.setView(ListView)
		} else {
			f.setView(GridView)
		}
	})
	f.setView(view)

	f.loadFavorites()

	f.favoritesList = widget.NewList(
		func() int {
			return len(f.favorites)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(container.New(&iconPaddingLayout{}, widget.NewIcon(theme.DocumentIcon())), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Icon).SetResource(f.favorites[id].locIcon)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(f.favorites[id].locName)
		},
	)
	f.favoritesList.OnSelected = func(id widget.ListItemID) {
		f.setLocation(f.favorites[id].loc)
	}

	var optionsButton *widget.Button
	optionsButton = widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		f.optionsMenu(fyne.CurrentApp().Driver().AbsolutePositionForObject(optionsButton), optionsButton.Size())
	})

	newFolderButton := widget.NewButtonWithIcon("", theme.FolderNewIcon(), func() {
		newFolderEntry := widget.NewEntry()
		ShowForm(lang.L("New Folder"), lang.L("Create Folder"), lang.L("Cancel"), []*widget.FormItem{
			{
				Text:   lang.X("file.name", "Name"),
				Widget: newFolderEntry,
			},
		}, func(s bool) {
			if !s || newFolderEntry.Text == "" {
				return
			}

			newFolderPath := filepath.Join(f.dir.Path(), newFolderEntry.Text)
			createFolderErr := os.MkdirAll(newFolderPath, 0750)
			if createFolderErr != nil {
				fyne.LogError(
					fmt.Sprintf("Failed to create folder with path %s", newFolderPath),
					createFolderErr,
				)
				ShowError(errors.New("folder cannot be created"), f.file.parent)
			}
			f.refreshDir(f.dir)
		}, f.file.parent)
	})

	optionsbuttons := container.NewHBox(
		newFolderButton,
		f.toggleViewButton,
		optionsButton,
	)

	header := container.NewBorder(nil, nil, nil, optionsbuttons,
		optionsbuttons, widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	footer := container.NewBorder(nil, nil, nil, buttons,
		buttons, container.NewHScroll(f.fileName),
	)

	body := container.NewHSplit(
		f.favoritesList,
		container.NewBorder(f.breadcrumbScroll, nil, nil, nil,
			f.breadcrumbScroll, f.filesScroll,
		),
	)
	body.SetOffset(0) // Set the minimum offset so that the favoritesList takes only it's minimal width

	return container.NewBorder(header, footer, nil, nil, body)
}

func (f *fileDialog) makeOpenButton(label string) *widget.Button {
	btn := widget.NewButton(label, func() {
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
			callback(storage.Reader(f.selected))
		} else if f.file.isDirectory() {
			callback := f.file.callback.(func(fyne.ListableURI, error))
			f.win.Hide()
			if f.file.onClosedCallback != nil {
				f.file.onClosedCallback(true)
			}
			callback(f.dir, nil)
		}
	})

	btn.Importance = widget.HighImportance
	btn.Disable()

	return btn
}

func (f *fileDialog) makeDismissButton(label string) *widget.Button {
	btn := widget.NewButton(label, func() {
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

	return btn
}

func (f *fileDialog) optionsMenu(position fyne.Position, buttonSize fyne.Size) {
	hiddenFiles := widget.NewCheck(lang.L("Show Hidden Files"), func(changed bool) {
		f.showHidden = changed
		f.refreshDir(f.dir)
	})
	hiddenFiles.Checked = f.showHidden
	hiddenFiles.Refresh()
	content := container.NewVBox(hiddenFiles)

	p := position.Add(buttonSize)
	pos := fyne.NewPos(p.X-content.MinSize().Width-theme.Padding()*2, p.Y+theme.Padding()*2)
	widget.ShowPopUpAtPosition(content, f.win.Canvas, pos)
}

func (f *fileDialog) loadFavorites() {
	favoriteLocations, err := getFavoriteLocations()
	if err != nil {
		fyne.LogError("Getting favorite locations", err)
	}
	favoriteIcons := getFavoriteIcons()
	favoriteOrder := getFavoriteOrder()

	f.favorites = []favoriteItem{
		{locName: "Home", locIcon: theme.HomeIcon(), loc: favoriteLocations["Home"]}}
	app := fyne.CurrentApp()
	if hasAppFiles(app) {
		f.favorites = append(f.favorites,
			favoriteItem{locName: "App Files", locIcon: theme.FileIcon(), loc: storageURI(app)})
	}
	f.favorites = append(f.favorites, f.getPlaces()...)

	for _, locName := range favoriteOrder {
		loc, ok := favoriteLocations[locName]
		if !ok {
			continue
		}
		locIcon := favoriteIcons[locName]
		f.favorites = append(f.favorites,
			favoriteItem{locName: locName, locIcon: locIcon, loc: loc})
	}
}

func (f *fileDialog) refreshDir(dir fyne.ListableURI) {
	f.dataLock.Lock()
	f.data = nil
	f.dataLock.Unlock()

	files, err := dir.List()
	if err != nil {
		fyne.LogError("Unable to read ListableURI "+dir.String(), err)
		return
	}

	var icons []fyne.URI
	parent, err := storage.Parent(dir)
	if err != nil && err != repository.ErrURIRoot {
		fyne.LogError("Unable to get parent of "+dir.String(), err)
		return
	}
	if parent != nil && parent.String() != dir.String() {
		icons = append(icons, parent)
	}

	for _, file := range files {
		if !f.showHidden && isHidden(file) {
			continue
		}

		listable, err := storage.ListerForURI(file)
		if f.file.isDirectory() && err != nil {
			continue
		} else if err == nil { // URI points to a directory
			icons = append(icons, listable)
		} else if f.file.filter == nil || f.file.filter.Matches(file) {
			icons = append(icons, file)
		}
	}

	f.dataLock.Lock()
	f.data = icons
	f.dataLock.Unlock()

	f.files.Refresh()
	f.filesScroll.Offset = fyne.NewPos(0, 0)
	f.filesScroll.Refresh()
}

func (f *fileDialog) setLocation(dir fyne.URI) error {
	if f.selectedID > -1 {
		f.files.Unselect(f.selectedID)
	}
	if dir == nil {
		return errors.New("failed to open nil directory")
	}
	list, err := storage.ListerForURI(dir)
	if err != nil {
		return err
	}

	isFav := false
	for i, fav := range f.favorites {
		if fav.loc == nil {
			continue
		}
		if fav.loc.Path() == dir.Path() {
			f.favoritesList.Select(i)
			isFav = true
			break
		}
	}
	if !isFav {
		f.favoritesList.UnselectAll()
	}

	f.setSelected(nil, -1)
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

	f.breadcrumbScroll.Refresh()
	f.breadcrumbScroll.Offset.X = f.breadcrumbScroll.Content.Size().Width - f.breadcrumbScroll.Size().Width
	f.breadcrumbScroll.Refresh()

	if f.file.isDirectory() {
		f.fileName.SetText(dir.Name())
		f.open.Enable()
	}
	f.refreshDir(list)

	return nil
}

func (f *fileDialog) setSelected(file fyne.URI, id int) {
	if file != nil {
		if listable, err := storage.CanList(file); err == nil && listable {
			f.setLocation(file)
			return
		}
	}
	f.selected = file
	f.selectedID = id

	if file == nil || file.String()[len(file.Scheme())+3:] == "" {
		// keep user input while navigating
		// in a FileSave dialog
		if !f.file.save {
			f.fileName.SetText("")
			f.open.Disable()
		}
	} else {
		f.fileName.SetText(file.Name())
		f.open.Enable()
	}
}

func (f *fileDialog) setView(view ViewLayout) {
	f.view = view
	fyne.CurrentApp().Preferences().SetInt(viewLayoutKey, int(view))
	var selectF func(id int)
	choose := func(id int) {
		if file, ok := f.getDataItem(id); ok {
			f.selectedID = id
			f.setSelected(file, id)
		}
	}
	count := func() int {
		f.dataLock.RLock()
		defer f.dataLock.RUnlock()

		return len(f.data)
	}
	template := func() fyne.CanvasObject {
		return f.newFileItem(storage.NewFileURI("./tempfile"), true, false)
	}
	update := func(id widget.GridWrapItemID, o fyne.CanvasObject) {
		if dir, ok := f.getDataItem(id); ok {
			parent := id == 0 && len(dir.Path()) < len(f.dir.Path())
			_, isDir := dir.(fyne.ListableURI)
			o.(*fileDialogItem).setLocation(dir, isDir || parent, parent)
			o.(*fileDialogItem).choose = selectF
			o.(*fileDialogItem).id = id
			o.(*fileDialogItem).open = f.open.OnTapped
		}
	}
	// Actually, during the real interaction, the OnSelected won't be called.
	// It will be called only when we directly calls container.select(i)
	if f.view == GridView {
		grid := widget.NewGridWrap(count, template, update)
		grid.OnSelected = choose
		f.files = grid
		f.toggleViewButton.SetIcon(theme.ListIcon())
		selectF = grid.Select
	} else {
		list := widget.NewList(count, template, update)
		list.OnSelected = choose
		f.files = list
		f.toggleViewButton.SetIcon(theme.GridIcon())
		selectF = list.Select
	}

	if f.dir != nil {
		f.refreshDir(f.dir)
	}
	f.filesScroll.Content = container.NewPadded(f.files)
	f.filesScroll.Refresh()
}

func (f *fileDialog) getDataItem(id int) (fyne.URI, bool) {
	f.dataLock.RLock()
	defer f.dataLock.RUnlock()

	if id >= len(f.data) {
		return nil, false
	}

	return f.data[id], true
}

// effectiveStartingDir calculates the directory at which the file dialog should
// open, based on the values of startingDirectory, CWD, home, and any error
// conditions which occur.
//
// Order of precedence is:
//
//   - file.startingDirectory if non-empty, os.Stat()-able, and uses the file://
//     URI scheme
//   - os.UserHomeDir()
//   - os.Getwd()
//   - "/" (should be filesystem root on all supported platforms)
func (f *FileDialog) effectiveStartingDir() fyne.ListableURI {
	if f.startingLocation != nil {
		if f.startingLocation.Scheme() == "file" {
			path := f.startingLocation.Path()

			// the starting directory is set explicitly
			if _, err := os.Stat(path); err != nil {
				fyne.LogError("Error with StartingLocation", err)
			} else {
				return f.startingLocation
			}
		}
		return f.startingLocation
	}

	// Try app storage
	app := fyne.CurrentApp()
	if hasAppFiles(app) {
		list, _ := storage.ListerForURI(storageURI(app))
		return list
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
	d := &fileDialog{file: file, initialFileName: file.initialFileName, view: GridView}
	ui := d.makeUI()
	pad := theme.Padding()
	itemMin := d.newFileItem(storage.NewFileURI("filename.txt"), false, false).MinSize()
	size := ui.MinSize().Add(itemMin.AddWidthHeight(itemMin.Width+pad*4, pad*2))

	d.win = widget.NewModalPopUp(ui, file.parent.Canvas())
	d.win.Resize(size)

	d.setLocation(file.effectiveStartingDir())
	d.win.Show()
	if file.save {
		d.win.Canvas.Focus(d.fileName.(*widget.Entry))
	}
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

// SetConfirmText allows custom text to be set in the confirmation button
//
// Since: 2.2
func (f *FileDialog) SetConfirmText(label string) {
	f.confirmText = label
	if f.dialog == nil {
		return
	}
	f.dialog.open.SetText(label)
	f.dialog.win.Refresh()
}

// SetDismissText allows custom text to be set in the dismiss button
func (f *FileDialog) SetDismissText(label string) {
	f.dismissText = label
	if f.dialog == nil {
		return
	}
	f.dialog.dismiss.SetText(label)
	f.dialog.win.Refresh()
}

// SetLocation tells this FileDialog which location to display.
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
	// If there is already a callback set, remember it and call both.
	originalCallback := f.onClosedCallback

	f.onClosedCallback = func(response bool) {
		if f.dialog == nil {
			return
		}
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

// SetView changes the default display view of the FileDialog
// This is normally called before the dialog is shown.
//
// Since: 2.5
func (f *FileDialog) SetView(v ViewLayout) {
	f.initialView = v
	if f.dialog != nil {
		f.dialog.setView(v)
	}
}

// NewFileOpen creates a file dialog allowing the user to choose a file to open.
//
// The callback function will run when the dialog closes and provide a reader for the chosen file.
// The reader will be nil when the user cancels or when nothing is selected.
// When the reader isn't nil it must be closed by the callback.
//
// The dialog will appear over the window specified when Show() is called.
func NewFileOpen(callback func(reader fyne.URIReadCloser, err error), parent fyne.Window) *FileDialog {
	dialog := &FileDialog{callback: callback, parent: parent}
	return dialog
}

// NewFileSave creates a file dialog allowing the user to choose a file to save
// to (new or overwrite). If the user chooses an existing file they will be
// asked if they are sure.
//
// The callback function will run when the dialog closes and provide a writer for the chosen file.
// The writer will be nil when the user cancels or when nothing is selected.
// When the writer isn't nil it must be closed by the callback.
//
// The dialog will appear over the window specified when Show() is called.
func NewFileSave(callback func(writer fyne.URIWriteCloser, err error), parent fyne.Window) *FileDialog {
	dialog := &FileDialog{callback: callback, parent: parent, save: true}
	return dialog
}

// ShowFileOpen creates and shows a file dialog allowing the user to choose a
// file to open.
//
// The callback function will run when the dialog closes and provide a reader for the chosen file.
// The reader will be nil when the user cancels or when nothing is selected.
// When the reader isn't nil it must be closed by the callback.
//
// The dialog will appear over the window specified.
func ShowFileOpen(callback func(reader fyne.URIReadCloser, err error), parent fyne.Window) {
	dialog := NewFileOpen(callback, parent)
	if fileOpenOSOverride(dialog) {
		return
	}
	dialog.Show()
}

// ShowFileSave creates and shows a file dialog allowing the user to choose a
// file to save to (new or overwrite). If the user chooses an existing file they
// will be asked if they are sure.
//
// The callback function will run when the dialog closes and provide a writer for the chosen file.
// The writer will be nil when the user cancels or when nothing is selected.
// When the writer isn't nil it must be closed by the callback.
//
// The dialog will appear over the window specified.
func ShowFileSave(callback func(writer fyne.URIWriteCloser, err error), parent fyne.Window) {
	dialog := NewFileSave(callback, parent)
	if fileSaveOSOverride(dialog) {
		return
	}
	dialog.Show()
}

func getFavoriteIcons() map[string]fyne.Resource {
	if runtime.GOOS == "darwin" {
		return map[string]fyne.Resource{
			"Documents": theme.DocumentIcon(),
			"Desktop":   theme.DesktopIcon(),
			"Downloads": theme.DownloadIcon(),
			"Music":     theme.MediaMusicIcon(),
			"Pictures":  theme.MediaPhotoIcon(),
			"Movies":    theme.MediaVideoIcon(),
		}
	}

	return map[string]fyne.Resource{
		"Documents": theme.DocumentIcon(),
		"Desktop":   theme.DesktopIcon(),
		"Downloads": theme.DownloadIcon(),
		"Music":     theme.MediaMusicIcon(),
		"Pictures":  theme.MediaPhotoIcon(),
		"Videos":    theme.MediaVideoIcon(),
	}
}

func getFavoriteOrder() []string {
	order := []string{
		"Desktop",
		"Documents",
		"Downloads",
		"Music",
		"Pictures",
		"Videos",
	}

	if runtime.GOOS == "darwin" {
		order[5] = "Movies"
	}

	return order
}

func hasAppFiles(a fyne.App) bool {
	return len(a.Storage().List()) > 0
}

func storageURI(a fyne.App) fyne.URI {
	dir, _ := storage.Child(a.Storage().RootURI(), "Documents")
	return dir
}

// iconPaddingLayout adds padding to the left of a widget.Icon().
// NOTE: It assumes that the slice only contains one item.
type iconPaddingLayout struct {
}

func (i *iconPaddingLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padding := theme.Padding() * 2
	objects[0].Move(fyne.NewPos(padding, 0))
	objects[0].Resize(size.SubtractWidthHeight(padding, 0))
}

func (i *iconPaddingLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return objects[0].MinSize().AddWidthHeight(theme.Padding()*2, 0)
}
