package repository

import (
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
)

// declare conformance to interfaces
var (
	_ io.ReadCloser       = (*nodeReaderWriter)(nil)
	_ io.WriteCloser      = (*nodeReaderWriter)(nil)
	_ fyne.URIReadCloser  = (*nodeReaderWriter)(nil)
	_ fyne.URIWriteCloser = (*nodeReaderWriter)(nil)
)

// declare conformance with repository types
var (
	_ repository.Repository             = (*InMemoryRepository)(nil)
	_ repository.WritableRepository     = (*InMemoryRepository)(nil)
	_ repository.AppendableRepository   = (*InMemoryRepository)(nil)
	_ repository.HierarchicalRepository = (*InMemoryRepository)(nil)
	_ repository.CopyableRepository     = (*InMemoryRepository)(nil)
	_ repository.MovableRepository      = (*InMemoryRepository)(nil)
	_ repository.ListableRepository     = (*InMemoryRepository)(nil)
)

// nodeReaderWriter allows reading or writing to elements in a InMemoryRepository
type nodeReaderWriter struct {
	path        string
	repo        *InMemoryRepository
	writing     bool
	readCursor  int
	writeCursor int
}

// InMemoryRepository implements an in-memory version of the
// repository.Repository type. It is useful for writing test cases, and may
// also be of use as a template for people wanting to implement their own
// "virtual repository". In future, we may consider moving this into the public
// API.
//
// Because of its design, this repository has several quirks:
//
// * The Parent() of a path that exists does not necessarily exist
//
//   - Listing takes O(number of extant paths in the repository), rather than
//     O(number of children of path being listed).
//
// This repository is not designed to be particularly fast or robust, but
// rather to be simple and easy to read. If you need performance, look
// elsewhere.
//
// Since: 2.0
type InMemoryRepository struct {
	// Data is exposed to allow tests to directly insert their own data
	// without having to go through the API
	Data map[string][]byte

	scheme string
}

// Read reads data from the repository into the provided buffer.
func (n *nodeReaderWriter) Read(p []byte) (int, error) {
	// first make sure the requested path actually exists
	data, ok := n.repo.Data[n.path]
	if !ok {
		return 0, fmt.Errorf("path '%s' not present in InMemoryRepository", n.path)
	}

	// copy it into p - we maintain counts since len(data) may be smaller
	// than len(p)
	count := 0
	j := 0 // index into p
	for ; (j < len(p)) && (n.readCursor < len(data)); n.readCursor++ {
		p[j] = data[n.readCursor]
		count++
		j++
	}

	// generate EOF if needed
	var err error = nil
	if n.readCursor >= len(data) {
		err = io.EOF
	}

	return count, err
}

// Close closes the reader and writer.
func (n *nodeReaderWriter) Close() error {
	n.readCursor = 0
	n.writeCursor = 0
	n.writing = false
	return nil
}

// Write writes data to the repository.
//
// This implementation automatically creates the path n.path if it does not
// exist. If it does exist, it is overwritten.
func (n *nodeReaderWriter) Write(p []byte) (int, error) {
	// overwrite the file if we haven't already started writing to it
	if !n.writing {
		n.repo.Data[n.path] = make([]byte, 0, len(p))
		n.writing = true
	}

	// copy the data into the node buffer
	for start := n.writeCursor; n.writeCursor < start+len(p); n.writeCursor++ {
		// extend the file if needed
		if len(n.repo.Data) < n.writeCursor+len(p) {
			n.repo.Data[n.path] = append(n.repo.Data[n.path], 0)
		}
		n.repo.Data[n.path][n.writeCursor] = p[n.writeCursor-start]
	}

	return len(p), nil
}

// URI returns the URI of the node.
func (n *nodeReaderWriter) URI() fyne.URI {
	// discarding the error because this should never fail
	u, _ := storage.ParseURI(n.repo.scheme + "://" + n.path)
	return u
}

// NewInMemoryRepository creates a new InMemoryRepository instance. It must be
// given the scheme it is registered for. The caller needs to call
// repository.Register() on the result of this function.
//
// Since: 2.0
func NewInMemoryRepository(scheme string) *InMemoryRepository {
	return &InMemoryRepository{
		Data:   make(map[string][]byte),
		scheme: scheme,
	}
}

// Exists checks if the given URI exists.
//
// Since: 2.0
func (m *InMemoryRepository) Exists(u fyne.URI) (bool, error) {
	path := u.Path()
	if path == "" {
		return false, fmt.Errorf("invalid path '%s'", path)
	}

	_, ok := m.Data[path]
	return ok, nil
}

// Reader reads the contents of the given URI.
//
// Since: 2.0
func (m *InMemoryRepository) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	path := u.Path()

	if path == "" {
		return nil, fmt.Errorf("invalid path '%s'", path)
	}

	_, ok := m.Data[path]
	if !ok {
		return nil, fmt.Errorf("no such path '%s' in InMemoryRepository", path)
	}

	return &nodeReaderWriter{path: path, repo: m}, nil
}

// CanRead checks if the given URI can be read.
//
// Since: 2.0
func (m *InMemoryRepository) CanRead(u fyne.URI) (bool, error) {
	path := u.Path()
	if path == "" {
		return false, fmt.Errorf("invalid path '%s'", path)
	}

	_, ok := m.Data[path]
	return ok, nil
}

// Destroy tears down the InMemoryRepository.
func (m *InMemoryRepository) Destroy(scheme string) {
	// do nothing
}

// Writer writes to the given URI.
//
// Since: 2.0
func (m *InMemoryRepository) Writer(u fyne.URI) (fyne.URIWriteCloser, error) {
	path := u.Path()
	if path == "" {
		return nil, fmt.Errorf("invalid path '%s'", path)
	}

	return &nodeReaderWriter{path: path, repo: m}, nil
}

// Appender returns a writer that appends to the given URI.
//
// Since: 2.6
func (m *InMemoryRepository) Appender(u fyne.URI) (fyne.URIWriteCloser, error) {
	path := u.Path()
	if path == "" {
		return nil, fmt.Errorf("invalid path '%s'", path)
	}

	return &nodeReaderWriter{path: path, repo: m, writing: true, writeCursor: len(m.Data[path])}, nil
}

// CanWrite checks if the given URI can be written to.
//
// Since: 2.0
func (m *InMemoryRepository) CanWrite(u fyne.URI) (bool, error) {
	if p := u.Path(); p == "" {
		return false, fmt.Errorf("invalid path '%s'", p)
	}

	return true, nil
}

// Delete deletes the given URI.
//
// Since: 2.0
func (m *InMemoryRepository) Delete(u fyne.URI) error {
	path := u.Path()
	_, ok := m.Data[path]
	if ok {
		delete(m.Data, path)
	}

	return nil
}

// Parent returns the parent URI of the given URI.
//
// Since: 2.0
func (m *InMemoryRepository) Parent(u fyne.URI) (fyne.URI, error) {
	return repository.GenericParent(u)
}

// Child returns the child URI created from the given URI and component.
//
// Since: 2.0
func (m *InMemoryRepository) Child(u fyne.URI, component string) (fyne.URI, error) {
	return repository.GenericChild(u, component)
}

// Copy copies the source URI to the destination URI.
//
// Since: 2.0
func (m *InMemoryRepository) Copy(source, destination fyne.URI) error {
	return repository.GenericCopy(source, destination)
}

// Move moves the contents of the source URI to the destination.
//
// Since: 2.0
func (m *InMemoryRepository) Move(source, destination fyne.URI) error {
	return repository.GenericMove(source, destination)
}

// CanList checks if the given URI can be listed.
//
// Since: 2.0
func (m *InMemoryRepository) CanList(u fyne.URI) (bool, error) {
	path := u.Path()
	exist, err := m.Exists(u)
	if err != nil || !exist {
		return false, err
	}

	if path == "" || path[len(path)-1] == '/' {
		return true, nil
	}

	children, err := m.List(u)
	return len(children) > 0, err
}

// List returns a list of URIs that are children of the given URI.
//
// Since: 2.0
func (m *InMemoryRepository) List(u fyne.URI) ([]fyne.URI, error) {
	// Get the prefix, and make sure it ends with a path separator so that
	// HasPrefix() will only find things that are children of it - this
	// solves the edge case where you have say '/foo/bar' and
	// '/foo/barbaz'.
	prefix := u.Path()

	if len(prefix) > 0 && prefix[len(prefix)-1] != '/' {
		prefix = prefix + "/"
	}

	prefixSplit := strings.Split(prefix, "/")
	prefixSplitLen := len(prefixSplit)

	// Now we can simply loop over all the paths and find the ones with an
	// appropriate prefix, then eliminate those with too many path
	// components.
	listing := []fyne.URI{}
	for p := range m.Data {
		// We are going to compare ncomp with the number of elements in
		// prefixSplit, which is guaranteed to have a trailing slash,
		// so we want to also make pSplit be counted in ncomp like it
		// does not have one.
		pSplit := strings.Split(p, "/")
		ncomp := len(pSplit)
		if len(p) > 0 && p[len(p)-1] == '/' {
			ncomp--
		}

		if strings.HasPrefix(p, prefix) && ncomp == prefixSplitLen {
			uri, err := storage.ParseURI(m.scheme + "://" + p)
			if err != nil {
				return nil, err
			}

			listing = append(listing, uri)
		}
	}

	return listing, nil
}

// CreateListable makes the given URI a listable URI.
//
// Since: 2.0
func (m *InMemoryRepository) CreateListable(u fyne.URI) error {
	ex, err := m.Exists(u)
	if err != nil {
		return err
	}
	path := u.Path()
	if ex {
		return fmt.Errorf("cannot create '%s' as a listable path because it already exists", path)
	}
	m.Data[path] = []byte{}
	return nil
}
