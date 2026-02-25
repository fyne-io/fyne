package repository

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
)

// declare conformance with repository types
var (
	_ repository.Repository             = (*FileRepository)(nil)
	_ repository.WritableRepository     = (*FileRepository)(nil)
	_ repository.DeleteAllRepository    = (*FileRepository)(nil)
	_ repository.AppendableRepository   = (*FileRepository)(nil)
	_ repository.HierarchicalRepository = (*FileRepository)(nil)
	_ repository.ListableRepository     = (*FileRepository)(nil)
	_ repository.MovableRepository      = (*FileRepository)(nil)
	_ repository.CopyableRepository     = (*FileRepository)(nil)
)

var (
	_ fyne.URIReadCloser  = (*file)(nil)
	_ fyne.URIWriteCloser = (*file)(nil)
)

type file struct {
	*os.File
	uri fyne.URI
}

func (f *file) URI() fyne.URI {
	return f.uri
}

// FileRepository implements a simple wrapper around Go's filesystem
// interface libraries. It should be registered by the driver on platforms
// where it is appropriate to do so.
//
// This repository is suitable to handle the file:// scheme.
//
// Since: 2.0
type FileRepository struct{}

// NewFileRepository creates a new FileRepository instance.
// The caller needs to call repository.Register() with the result of this function.
//
// Since: 2.0
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

// Exists checks if the given URI exists.
//
// Since: 2.0
func (r *FileRepository) Exists(u fyne.URI) (bool, error) {
	p := u.Path()
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Reader returns a reader for the given URI.
//
// Since: 2.0
func (r *FileRepository) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	return openFile(u, false, false)
}

// CanRead checks if the given URI can be read.
//
// Since: 2.0
func (r *FileRepository) CanRead(u fyne.URI) (bool, error) {
	f, err := os.OpenFile(u.Path(), os.O_RDONLY, 0o666)
	if err != nil {
		if os.IsPermission(err) || os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, f.Close()
}

// Destroy tears down the repository for the specified scheme.
func (r *FileRepository) Destroy(scheme string) {
	// do nothing
}

// Writer returns a truncating writer for the given URI.
//
// Since: 2.0
func (r *FileRepository) Writer(u fyne.URI) (fyne.URIWriteCloser, error) {
	return openFile(u, true, true)
}

// Appender returns a writer that appends to the given URI.
//
// Since: 2.6
func (r *FileRepository) Appender(u fyne.URI) (fyne.URIWriteCloser, error) {
	return openFile(u, true, false)
}

// CanWrite checks if the given URI can be written.
//
// Since: 2.0
func (r *FileRepository) CanWrite(u fyne.URI) (bool, error) {
	f, err := os.OpenFile(u.Path(), os.O_WRONLY, 0o666)
	if err != nil {
		if os.IsPermission(err) {
			return false, nil
		}

		if os.IsNotExist(err) {
			// We may need to do extra logic to check if the
			// directory is writable, but presumably the
			// IsPermission check covers this.
			return true, nil
		}

		return false, err
	}

	return true, f.Close()
}

// Delete deletes the given URI.
//
// Since: 2.0
func (r *FileRepository) Delete(u fyne.URI) error {
	return os.Remove(u.Path())
}

// DeleteAll deletes the given URI and all its children.
//
// Since: 2.7
func (r *FileRepository) DeleteAll(u fyne.URI) error {
	return os.RemoveAll(u.Path())
}

// Parent returns the parent URI of the given URI.
//
// Since: 2.0
func (r *FileRepository) Parent(u fyne.URI) (fyne.URI, error) {
	child := path.Clean(u.Path())
	if child == "." || // Clean ending up empty returns ".".
		strings.HasSuffix(child, "/") || // Only root has trailing slash.
		runtime.GOOS == "windows" && len(child) == 2 && child[1] == ':' {
		return nil, repository.ErrURIRoot
	}

	parent := path.Dir(child)
	if parent == "/" {
		return storage.NewFileURI("/"), nil
	}

	return storage.NewFileURI(parent + "/"), nil
}

// Child creates a child URI from the given URI and component.
//
// Since: 2.0
func (r *FileRepository) Child(u fyne.URI, component string) (fyne.URI, error) {
	return storage.NewFileURI(path.Join(u.Path(), component)), nil
}

// List returns a list of all child URIs of the given URI.
//
// Since: 2.0
func (r *FileRepository) List(u fyne.URI) ([]fyne.URI, error) {
	p := u.Path()
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}

	urilist := make([]fyne.URI, len(files))
	for i, f := range files {
		urilist[i] = storage.NewFileURI(path.Join(p, f.Name()))
	}

	return urilist, nil
}

// CreateListable creates a new directory at the given URI.
func (r *FileRepository) CreateListable(u fyne.URI) error {
	path := u.Path()
	return os.Mkdir(path, 0o755)
}

// CanList checks if the given URI can be listed.
//
// Since: 2.0
func (r *FileRepository) CanList(u fyne.URI) (bool, error) {
	p := u.Path()
	info, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if !info.IsDir() {
		return false, nil
	}

	if runtime.GOOS == "windows" && len(p) <= 3 {
		return true, nil // assume drives can be read, avoids hang if the drive is temporarily unresponsive
	}

	// We know it is a directory, but we don't know if we can read it, so
	// we'll just try to do so and see if we get a permission error.
	f, err := os.Open(p)
	if err == nil {
		_, err = f.Readdir(1)
		f.Close()
	}

	if err != nil && err != io.EOF {
		return false, err
	}

	if os.IsPermission(err) {
		return false, nil
	}

	// it is a directory, and checking the permissions did not error out
	return true, nil
}

// Copy copies the contents of the source URI to the destination URI.
//
// Since: 2.0
func (r *FileRepository) Copy(source, destination fyne.URI) error {
	err := fastCopy(destination.Path(), source.Path())
	if err == nil {
		return nil
	}

	return repository.GenericCopy(source, destination)
}

// Move moves the contents of the source URI to the destination URI.
//
// Since: 2.0
func (r *FileRepository) Move(source, destination fyne.URI) error {
	err := os.Rename(source.Path(), destination.Path())
	if err == nil {
		return nil
	}

	return repository.GenericMove(source, destination)
}

func copyFile(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func fastCopy(dst, src string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return copyFile(dst, src)
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, rel)
		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(dstPath, path)
	})
}

func openFile(uri fyne.URI, write bool, truncate bool) (*file, error) {
	path := uri.Path()
	var f *os.File
	var err error
	if write {
		if truncate {
			f, err = os.Create(path) // If it exists this will truncate which is what we wanted
		} else {
			f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
		}
	} else {
		f, err = os.Open(path)
	}
	return &file{File: f, uri: uri}, err
}
