//go:build flatpak && !windows && !android && !ios && !wasm && !js

package dialog

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/storage"

	"github.com/rymdport/portal"
	"github.com/rymdport/portal/filechooser"
)

func openFile(parentWindowHandle string, options *filechooser.OpenFileOptions) (fyne.URIReadCloser, error) {
	title := lang.L("Open") + " " + lang.L("File")
	uri, err := open(parentWindowHandle, title, options)
	if err != nil || uri == nil {
		return nil, err
	}

	return storage.Reader(uri)
}

func openFolder(parentWindowHandle string, options *filechooser.OpenFileOptions) (fyne.ListableURI, error) {
	title := lang.L("Open") + " " + lang.L("Folder")
	uri, err := open(parentWindowHandle, title, options)
	if err != nil || uri == nil {
		return nil, err
	}

	return storage.ListerForURI(uri)
}

func open(parentWindowHandle, title string, options *filechooser.OpenFileOptions) (fyne.URI, error) {
	uris, err := filechooser.OpenFile(parentWindowHandle, title, options)
	if err != nil {
		return nil, err
	}

	if len(uris) == 0 {
		return nil, nil
	}

	return storage.ParseURI(uris[0])
}

func saveFile(parentWindowHandle string, options *filechooser.SaveFileOptions) (fyne.URIWriteCloser, error) {
	title := lang.L("Save") + " " + lang.L("File")
	uris, err := filechooser.SaveFile(parentWindowHandle, title, options)
	if err != nil {
		return nil, err
	}

	if len(uris) == 0 {
		return nil, nil
	}

	uri, err := storage.ParseURI(uris[0])
	if err != nil {
		return nil, err
	}

	return storage.Writer(uri)
}

func fileOpenOSOverride(d *FileDialog) bool {
	options := &filechooser.OpenFileOptions{
		Directory:   d.isDirectory(),
		AcceptLabel: d.confirmText,
	}
	if d.startingLocation != nil {
		options.CurrentFolder = d.startingLocation.Path()
	}
	options.Filters, options.CurrentFilter = convertFilterForPortal(d.filter)

	windowHandle := windowHandleForPortal(d.parent)

	go func() {
		if options.Directory {
			folder, err := openFolder(windowHandle, options)
			fyne.Do(func() {
				folderCallback := d.callback.(func(fyne.ListableURI, error))
				folderCallback(folder, err)
			})
		} else {
			file, err := openFile(windowHandle, options)
			fyne.Do(func() {
				fileCallback := d.callback.(func(fyne.URIReadCloser, error))
				fileCallback(file, err)
			})
		}
	}()

	return true
}

func fileSaveOSOverride(d *FileDialog) bool {
	options := &filechooser.SaveFileOptions{
		AcceptLabel: d.confirmText,
		CurrentName: d.initialFileName,
	}
	if d.startingLocation != nil {
		options.CurrentFolder = d.startingLocation.Path()
	}
	options.Filters, options.CurrentFilter = convertFilterForPortal(d.filter)

	callback := d.callback.(func(fyne.URIWriteCloser, error))
	windowHandle := windowHandleForPortal(d.parent)

	go func() {
		file, err := saveFile(windowHandle, options)
		fyne.Do(func() {
			callback(file, err)
		})
	}()

	return true
}

func windowHandleForPortal(window fyne.Window) string {
	windowHandle := ""
	if !build.IsWayland {
		window.(driver.NativeWindow).RunNative(func(context any) {
			handle := context.(driver.X11WindowContext).WindowHandle
			windowHandle = portal.FormatX11WindowHandle(handle)
		})
	}

	// TODO: We need to get the Wayland handle from the xdg_foreign protocol and convert to string on the form "wayland:{id}".
	return windowHandle
}

func convertFilterForPortal(fyneFilter storage.FileFilter) (list []*filechooser.Filter, current *filechooser.Filter) {
	if fyneFilter == nil || fyneFilter == folderFilter {
		return nil, nil
	}

	if filter, ok := fyneFilter.(*storage.ExtensionFileFilter); ok {
		rules := make([]filechooser.Rule, 0, 2*len(filter.Extensions))
		for _, ext := range filter.Extensions {
			lowercase := filechooser.Rule{
				Type:    filechooser.GlobPattern,
				Pattern: "*" + strings.ToLower(ext),
			}
			uppercase := filechooser.Rule{
				Type:    filechooser.GlobPattern,
				Pattern: "*" + strings.ToUpper(ext),
			}
			rules = append(rules, lowercase, uppercase)
		}

		name := formatFilterName(filter.Extensions, 3)
		converted := &filechooser.Filter{Name: name, Rules: rules}
		return []*filechooser.Filter{converted}, converted
	}

	if filter, ok := fyneFilter.(*storage.MimeTypeFileFilter); ok {
		rules := make([]filechooser.Rule, len(filter.MimeTypes))
		for i, mime := range filter.MimeTypes {
			rules[i] = filechooser.Rule{
				Type:    filechooser.MIMEType,
				Pattern: mime,
			}
		}

		name := formatFilterName(filter.MimeTypes, 3)
		converted := &filechooser.Filter{Name: name, Rules: rules}
		return []*filechooser.Filter{converted}, converted
	}

	return nil, nil
}

func formatFilterName(patterns []string, count int) string {
	if len(patterns) < count {
		count = len(patterns)
	}

	name := strings.Join(patterns[:count], ", ")
	if len(patterns) > count {
		name += "…"
	}

	return name
}
