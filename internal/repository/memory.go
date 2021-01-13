package repository

import (
	"fyne.io/fyne"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/storage/repository"

	"fmt"
	"io"
)

// declare conformance to interfaces
var _ io.ReadCloser = &nodeReaderWriter{}
var _ io.WriteCloser = &nodeReaderWriter{}
var _ fyne.URIReadCloser = &nodeReaderWriter{}
var _ fyne.URIWriteCloser = &nodeReaderWriter{}

// nodeReaderWriter allows reading or writing to elements in a InMemoryRepository
type nodeReaderWriter struct {
	path        string
	repo        *InMemoryRepository
	writing     bool
	readCursor  int
	writeCursor int
}

// Read implements io.Reader.Read
func (n *nodeReaderWriter) Read(p []byte) (int, error) {

	// first make sure the requested path actually exists
	data, ok := n.repo.data[n.path]
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

// Close implements io.Closer.Close
func (n *nodeReaderWriter) Close() error {
	n.readCursor = 0
	n.writeCursor = 0
	n.writing = false
	return nil
}

// Write implements io.Writer.Write
//
// This implementation automatically creates the path n.path if it does not
// exist. If it does exist, it is overwritten.
func (n *nodeReaderWriter) Write(p []byte) (int, error) {

	// guarantee that the path exists
	_, ok := n.repo.data[n.path]
	if !ok {
		n.repo.data[n.path] = []byte{}
	}

	// overwrite the file if we haven't already started writing to it
	if !n.writing {
		n.repo.data[n.path] = make([]byte, 0)
		n.writing = true
	}

	// copy the data into the node buffer
	count := 0
	start := n.writeCursor
	for ; n.writeCursor < start+len(p); n.writeCursor++ {
		// extend the file if needed
		if len(n.repo.data) < n.writeCursor+len(p) {
			n.repo.data[n.path] = append(n.repo.data[n.path], 0)
		}
		n.repo.data[n.path][n.writeCursor] = p[n.writeCursor-start]
		count++
	}

	return count, nil
}

// Name implements fyne.URI*Closer.URI
func (n *nodeReaderWriter) URI() fyne.URI {

	// discarding the error because this should never fail
	u, _ := storage.ParseURI(n.repo.scheme + "://" + n.path)

	return u
}

// declare conformance with repository types
var _ repository.Repository = &InMemoryRepository{}
var _ repository.WriteableRepository = &InMemoryRepository{}
var _ repository.HierarchicalRepository = &InMemoryRepository{}

// InMemoryRepository implements an in-memory version of the
// repository.Repository type. It is useful for writing test cases, and may
// also be of use as a template for people wanting to implement their own
// "virtual repository". In future, we may consider moving this into the public
// API.
//
// Because of it's design, this repository has several quirks:
//
// * The Parent() of a path that exists does not necessarily exist
//
// * Listing takes O(number of extant paths in the repository), rather than
//   O(number of children of path being listed).
//
// This repository is not designed to be particularly fast or robust, but
// rather to be simple and easy to read. If you need performance, look
// elsewhere.
//
// Since 2.0.0
type InMemoryRepository struct {
	data   map[string][]byte
	scheme string
}

// NewInMemoryRepository creates a new InMemoryRepository instance. It must be
// given the scheme it is registered for. The caller needs to call
// repository.Register() on the result of this function.
//
// Since 2.0.0
func NewInMemoryRepository(scheme string) *InMemoryRepository {
	return &InMemoryRepository{
		data:   make(map[string][]byte),
		scheme: scheme,
	}
}

// Exists implements repository.Repository.Exists
//
// Since 2.0.0
func (m *InMemoryRepository) Exists(u fyne.URI) (bool, error) {
	if u.Path() == "" {
		return false, fmt.Errorf("invalid path '%s'", u.Path())
	}

	_, ok := m.data[u.Path()]
	return ok, nil
}

// Reader implements repository.Repository.Reader
//
// Since 2.0.0
func (m *InMemoryRepository) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	if u.Path() == "" {
		return nil, fmt.Errorf("invalid path '%s'", u.Path())
	}

	_, ok := m.data[u.Path()]
	if !ok {
		return nil, fmt.Errorf("no such path '%s' in InMemoryRepository", u.Path())
	}

	return &nodeReaderWriter{path: u.Path(), repo: m}, nil
}

// CanRead implements repository.Repository.CanRead
//
// Since 2.0.0
func (m *InMemoryRepository) CanRead(u fyne.URI) (bool, error) {
	if u.Path() == "" {
		return false, fmt.Errorf("invalid path '%s'", u.Path())
	}

	_, ok := m.data[u.Path()]
	if !ok {
		return false, fmt.Errorf("no such path '%s' in InMemoryRepository", u.Path())
	}

	return true, nil
}

// Destroy implements repository.Repository.Destroy
func (m *InMemoryRepository) Destroy(scheme string) {
	// do nothing
}

// Writer implements repository.WriteableRepository.Writer
//
// Since 2.0.0
func (m *InMemoryRepository) Writer(u fyne.URI) (fyne.URIWriteCloser, error) {
	if u.Path() == "" {
		return nil, fmt.Errorf("invalid path '%s'", u.Path())
	}

	return &nodeReaderWriter{path: u.Path(), repo: m}, nil
}

// CanWrite implements repository.WriteableRepository.CanWrite
//
// Since 2.0.0
func (m *InMemoryRepository) CanWrite(u fyne.URI) (bool, error) {
	if u.Path() == "" {
		return false, fmt.Errorf("invalid path '%s'", u.Path())
	}

	return true, nil
}

// Delete implements repository.WriteableRepository.Delete
//
// Since 2.0.0
func (m *InMemoryRepository) Delete(u fyne.URI) error {
	_, ok := m.data[u.Path()]
	if ok {
		delete(m.data, u.Path())
	}

	return nil
}

// Parent implements repository.HierarchicalRepository.Parent
//
// Since 2.0.0
func (m *InMemoryRepository) Parent(u fyne.URI) (fyne.URI, error) {
	return repository.GenericParent(u, storage.ParseURI)
}

// Child implements repository.HierarchicalRepository.Child
//
// Since 2.0.0
func (m *InMemoryRepository) Child(u fyne.URI, component string) (fyne.URI, error) {
	return repository.GenericChild(u, component, storage.ParseURI)
}
