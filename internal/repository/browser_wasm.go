//go:build wasm

package repository

import (
	"bytes"
	"errors"
	"path"
	"sort"
	"strings"
	"sync"
	"syscall/js"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage/repository"
)

var (
	_ repository.Repository             = (*BrowserFileRepository)(nil)
	_ repository.ListableRepository     = (*BrowserFileRepository)(nil)
	_ repository.HierarchicalRepository = (*BrowserFileRepository)(nil)
)

// BrowserFileRepository is the repository for the "file" scheme in a web browser.
//
// A page cannot open files by where they are. What it has instead are the files and folders that
// the user chose with the browser's picker: the browser gives a handle to each, from which it can
// be read for as long as the page is open, but nothing of where it is. So each is given a URI
// by name alone - file:///<file> or file:///<folder>/<path within it> - that can be listed and
// read here. Choosing something of the same name again replaces what was chosen before.
type BrowserFileRepository struct {
	mu      sync.Mutex
	files   map[string]js.Value // by path
	folders map[string][]string // the names of what is in each folder, by path
}

var browserFiles = &BrowserFileRepository{
	files:   map[string]js.Value{},
	folders: map[string][]string{"/": nil},
}

// BrowserFiles returns the repository that holds what has been chosen with the browser's picker.
func BrowserFiles() *BrowserFileRepository {
	return browserFiles
}

// AddFile records a JavaScript File that the user chose, returning the URI it can be read by.
func (r *BrowserFileRepository) AddFile(file js.Value) fyne.URI {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := "/" + path.Base(path.Clean("/"+file.Get("name").String()))
	r.remove(p)
	r.add(p, file)
	return repository.NewFileURI(p)
}

// AddFolder records a JavaScript FileList, the files of a folder that the user chose, returning the
// URI of that folder. The list must not be empty.
func (r *BrowserFileRepository) AddFolder(files js.Value) fyne.URI {
	r.mu.Lock()
	defer r.mu.Unlock()

	paths := make([]string, files.Length())
	for i := range paths {
		file := files.Index(i)
		paths[i] = path.Clean("/" + file.Get("webkitRelativePath").String())
		if strings.Count(paths[i], "/") < 2 {
			// not within a folder, as when a browser cannot choose folders and offers files instead
			paths[i] = "/files/" + path.Base(path.Clean("/"+file.Get("name").String()))
		}
	}

	root := "/" + strings.Split(paths[0], "/")[1]
	r.remove(root)
	for i, p := range paths {
		r.add(p, files.Index(i))
	}
	return repository.NewFileURI(root)
}

// add records a file and the folders it is in. The lock must be held.
func (r *BrowserFileRepository) add(p string, file js.Value) {
	r.files[p] = file
	for child := p; child != "/"; child = path.Dir(child) {
		parent, name := path.Dir(child), path.Base(child)
		names := r.folders[parent]
		at := sort.SearchStrings(names, name)
		if at < len(names) && names[at] == name {
			return // the folders above are known too
		}
		names = append(names, "")
		copy(names[at+1:], names[at:])
		names[at] = name
		r.folders[parent] = names
	}
}

// remove forgets a file or folder at the top level, and everything in it. The lock must be held.
func (r *BrowserFileRepository) remove(root string) {
	for p := range r.files {
		if p == root || strings.HasPrefix(p, root+"/") {
			delete(r.files, p)
		}
	}
	for p := range r.folders {
		if p == root || strings.HasPrefix(p, root+"/") {
			delete(r.folders, p)
		}
	}

	names := r.folders["/"]
	name := path.Base(root)
	if at := sort.SearchStrings(names, name); at < len(names) && names[at] == name {
		r.folders["/"] = append(names[:at], names[at+1:]...)
	}
}

func (r *BrowserFileRepository) lookup(u fyne.URI) (file js.Value, isFile, isFolder bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := path.Clean("/" + u.Path())
	file, isFile = r.files[p]
	_, isFolder = r.folders[p]
	return file, isFile, isFolder
}

// Exists implements repository.Repository.Exists
func (r *BrowserFileRepository) Exists(u fyne.URI) (bool, error) {
	_, isFile, isFolder := r.lookup(u)
	return isFile || isFolder, nil
}

// Reader implements repository.Repository.Reader
//
// The browser has to be waited for to read a file, so this must not be called from an event handler.
func (r *BrowserFileRepository) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	file, isFile, _ := r.lookup(u)
	if !isFile {
		return nil, repository.ErrOperationNotSupported
	}
	return ReadBrowserFile(file, u)
}

// CanRead implements repository.Repository.CanRead
func (r *BrowserFileRepository) CanRead(u fyne.URI) (bool, error) {
	_, isFile, _ := r.lookup(u)
	return isFile, nil
}

// Destroy implements repository.Repository.Destroy
func (r *BrowserFileRepository) Destroy(string) {
}

// CanList implements repository.ListableRepository.CanList
func (r *BrowserFileRepository) CanList(u fyne.URI) (bool, error) {
	_, _, isFolder := r.lookup(u)
	return isFolder, nil
}

// List implements repository.ListableRepository.List
func (r *BrowserFileRepository) List(u fyne.URI) ([]fyne.URI, error) {
	p := path.Clean("/" + u.Path())
	r.mu.Lock()
	names, isFolder := r.folders[p]
	names = append([]string(nil), names...)
	r.mu.Unlock()
	if !isFolder {
		return nil, repository.ErrOperationNotSupported
	}

	list := make([]fyne.URI, len(names))
	for i, name := range names {
		list[i] = repository.NewFileURI(path.Join(p, name))
	}
	return list, nil
}

// CreateListable implements repository.ListableRepository.CreateListable
func (r *BrowserFileRepository) CreateListable(fyne.URI) error {
	return errors.New("files chosen in the browser cannot be changed")
}

// Parent implements repository.HierarchicalRepository.Parent
func (r *BrowserFileRepository) Parent(u fyne.URI) (fyne.URI, error) {
	return repository.GenericParent(u)
}

// Child implements repository.HierarchicalRepository.Child
func (r *BrowserFileRepository) Child(u fyne.URI, component string) (fyne.URI, error) {
	return repository.GenericChild(u, component)
}

// browserFile is the content of a file the user chose with the browser's file picker, held in memory.
type browserFile struct {
	*bytes.Reader
	uri fyne.URI
}

func (f *browserFile) Close() error {
	return nil
}

func (f *browserFile) URI() fyne.URI {
	return f.uri
}

// browserURI names a file chosen in the browser, reporting the type the browser gave it.
type browserURI struct {
	fyne.URI
	mime string
}

func (u *browserURI) MimeType() string {
	return u.mime
}

// ReadBrowserFile loads the content of a JavaScript File, to be known by the given URI.
// It waits for the browser and so must not be called from an event handler.
func ReadBrowserFile(file js.Value, uri fyne.URI) (fyne.URIReadCloser, error) {
	loaded := make(chan js.Value, 1)
	failed := make(chan error, 1)
	onLoad := js.FuncOf(func(_ js.Value, args []js.Value) any {
		loaded <- args[0]
		return nil
	})
	defer onLoad.Release()
	onFail := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) > 0 && args[0].Truthy() {
			failed <- js.Error{Value: args[0]}
		} else {
			failed <- errors.New("failed to read the chosen file")
		}
		return nil
	})
	defer onFail.Release()
	file.Call("arrayBuffer").Call("then", onLoad, onFail)

	var buffer js.Value
	select {
	case buffer = <-loaded:
	case err := <-failed:
		return nil, err
	}

	array := uint8Array.New(buffer)
	data := make([]byte, array.Length())
	js.CopyBytesToGo(data, array)

	if mime := file.Get("type").String(); mime != "" {
		uri = &browserURI{URI: uri, mime: mime}
	}
	return &browserFile{Reader: bytes.NewReader(data), uri: uri}, nil
}
