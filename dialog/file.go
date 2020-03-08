package dialog

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type fileDialog struct {
	fileName   *widget.Label
	open       *widget.Button
	breadcrumb *widget.Box
	files      *fyne.Container

	win      *widget.PopUp
	current  *fileIcon
	callback func(string)
}

func (f *fileDialog) makeUI() fyne.CanvasObject {
	fav := widget.NewGroup("Favourites", f.loadFavourites()...)

	f.fileName = widget.NewLabel("")
	f.open = widget.NewButton("Open", func() {
		f.win.Hide()
		if f.callback != nil && f.current != nil {
			f.callback(f.current.path)
		}
	})
	f.open.Style = widget.PrimaryButton
	f.open.Disable()
	footer := widget.NewHBox(f.fileName, layout.NewSpacer(),
		widget.NewButton("Cancel", func() {
			f.win.Hide()
			if f.callback != nil {
				f.callback("")
			}
		}),
		f.open)

	f.breadcrumb = widget.NewHBox()

	f.files = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(fileIconSize+16,
		fileIconSize+theme.Padding()+fileTextSize)),
	)

	scrollBread := widget.NewScrollContainer(f.breadcrumb)
	body := fyne.NewContainerWithLayout(layout.NewBorderLayout(scrollBread, nil, nil, nil),
		scrollBread, widget.NewScrollContainer(f.files))
	header := widget.NewLabelWithStyle("Open File", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(header, footer, fav, nil),
		fav, header, footer, body)
}

func (f *fileDialog) loadFavourites() []fyne.CanvasObject {
	home, _ := os.UserHomeDir()
	return []fyne.CanvasObject{
		widget.NewButton("Home", func() {
			f.setDirectory(home)
		}),
		widget.NewButton("Documents", func() {
			f.setDirectory(path.Join(home, "Documents"))
		}),
		widget.NewButton("Downloads", func() {
			f.setDirectory(path.Join(home, "Downloads"))
		}),
	}
}

func (f *fileDialog) setFileDir(dir *fileIcon) {
	f.setDirectory(dir.path)
}

func (f *fileDialog) setDirectory(dir string) {
	f.setFile(nil)

	f.breadcrumb.Children = nil
	buildDir := filepath.VolumeName(dir)
	for i, d := range strings.Split(dir, string(filepath.Separator)) {
		if d == "" {
			if i > 0 { // what we get if we split "/"
				break
			}
			buildDir = "/"
			d = "/"
		} else {
			buildDir = filepath.Join(buildDir, d)
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

func (f *fileDialog) refreshDir(dir string) {
	f.files.Objects = nil

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fyne.LogError("Unable to read path "+dir, err)
		return
	}

	var icons []fyne.CanvasObject
	parent := path.Dir(dir)
	if parent != dir {
		icons = append(icons, f.newFileIcon(theme.FolderOpenIcon(), path.Dir(dir), f.setFileDir))
	}
	for _, file := range files {
		if len(file.Name()) == 0 || file.Name()[0] == '.' {
			continue
		}

		itemPath := path.Join(dir, file.Name())
		if file.IsDir() {
			icons = append(icons, f.newFileIcon(theme.FolderIcon(), itemPath, f.setFileDir))
		} else {
			icons = append(icons, f.newFileIcon(theme.FileIcon(), itemPath, f.setFile))
		}
	}

	f.files.Objects = icons
	f.files.Refresh()
}

func (f *fileDialog) setFile(file *fileIcon) {
	if f.current != nil {
		f.current.current = false
		f.current.Refresh()
	}
	f.current = file

	if file == nil || file.path == "" {
		f.fileName.SetText("")
		f.open.Disable()
	} else {
		f.fileName.SetText(path.Base(file.path))
		f.open.Enable()
	}
}

// ShowFileOpen shows a file dialog allowing the user to choose a file to open.
// The dialog will appear over the window specified.
func ShowFileOpen(callback func(string), parent fyne.Window) {
	d := &fileDialog{callback: callback}
	ui := d.makeUI()
	dir, err := os.UserHomeDir()
	if err != nil {
		fyne.LogError("Could not load user home dir", err)
		dir, _ = os.Getwd() //fallback
	}
	d.setDirectory(dir)

	spacer := canvas.NewRectangle(theme.BackgroundColor())
	spacer.SetMinSize(fyne.NewSize(436, 300))
	content := fyne.NewContainerWithLayout(layout.NewMaxLayout(), spacer, ui)
	d.win = widget.NewModalPopUp(content, parent.Canvas())

	d.win.Show()
}
