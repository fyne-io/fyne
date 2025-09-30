package repository

import (
	"io"
	"path"
	"strings"

	"fyne.io/fyne/v2"
)

// GenericParent can be used as a common-case implementation of
// HierarchicalRepository.Parent(). It will create a parent URI based on
// IETF RFC3986.
//
// In short, the URI is separated into its component parts, the path component
// is split along instances of '/', and the trailing element is removed. The
// result is concatenated and parsed as a new URI.
//
// If the URI path is empty or '/', then a nil URI is returned, along with
// ErrURIRoot.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0
func GenericParent(u fyne.URI) (fyne.URI, error) {
	p := strings.TrimSuffix(u.Path(), "/")
	if p == "" {
		return nil, ErrURIRoot
	}

	newURI := uri{
		scheme:    u.Scheme(),
		authority: u.Authority(),
		path:      path.Dir(p),
		query:     u.Query(),
		fragment:  u.Fragment(),
	}

	// NOTE: we specifically want to use ParseURI, rather than &uri{},
	// since the repository for the URI we just created might be a
	// CustomURIRepository that implements its own ParseURI.
	// However, we can reuse &uri.String() to not duplicate string creation.
	return ParseURI(newURI.String())
}

// GenericChild can be used as a common-case implementation of
// HierarchicalRepository.Child(). It will create a child URI by separating the
// URI into its component parts as described in IETF RFC 3986, then appending
// "/" + component to the path, then concatenating the result and parsing it as
// a new URI.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0
func GenericChild(u fyne.URI, component string) (fyne.URI, error) {
	newURI := uri{
		scheme:    u.Scheme(),
		authority: u.Authority(),
		path:      path.Join(u.Path(), component),
		query:     u.Query(),
		fragment:  u.Fragment(),
	}

	// NOTE: we specifically want to use ParseURI, rather than &uri{},
	// since the repository for the URI we just created might be a
	// CustomURIRepository that implements its own ParseURI.
	// However, we can reuse &uri.String() to not duplicate string creation.
	return ParseURI(newURI.String())
}

// GenericCopy can be used a common-case implementation of
// CopyableRepository.Copy(). It will perform the copy by obtaining a reader
// for the source URI, a writer for the destination URI, then writing the
// contents of the source to the destination.
//
// For obvious reasons, the destination URI must have a registered
// WritableRepository.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0
func GenericCopy(source fyne.URI, destination fyne.URI) error {
	// Look up repositories for the source and destination.
	srcrepo, err := ForURI(source)
	if err != nil {
		return err
	}

	dstrepo, err := ForURI(destination)
	if err != nil {
		return err
	}

	// The destination must be writable.
	destwrepo, ok := dstrepo.(WritableRepository)
	if !ok {
		return ErrOperationNotSupported
	}

	if listable, ok := srcrepo.(ListableRepository); ok {
		isParent, err := listable.CanList(source)
		if err == nil && isParent {
			if srcrepo != destwrepo { // cannot copy folders between repositories
				return ErrOperationNotSupported
			}

			return genericCopyMoveListable(source, destination, srcrepo, false)
		}
	}

	// Create a reader and a writer.
	srcReader, err := srcrepo.Reader(source)
	if err != nil {
		return err
	}
	defer srcReader.Close()

	dstWriter, err := destwrepo.Writer(destination)
	if err != nil {
		return err
	}
	defer dstWriter.Close()

	// Perform the copy.
	_, err = io.Copy(dstWriter, srcReader)
	return err
}

// GenericDeleteAll can be used a common-case implementation of
// DeletableRepository.DeleteAll(). It will perform the deletion by obtaining
// a list of all items in the URI, then deleting each one.
//
// For obvious reasons, the URI must be writable.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.7
func GenericDeleteAll(u fyne.URI) error {
	repo, err := ForURI(u)
	if err != nil {
		return err
	}

	wrepo, ok := repo.(WritableRepository)
	if !ok {
		return ErrOperationNotSupported
	}

	lrepo, ok := repo.(ListableRepository)
	if !ok {
		return wrepo.Delete(u)
	}

	return genericDeleteAll(u, wrepo, lrepo)
}

// GenericMove can be used a common-case implementation of
// MovableRepository.Move(). It will perform the move by obtaining a reader
// for the source URI, a writer for the destination URI, then writing the
// contents of the source to the destination. Following this, the source
// will be deleted using WritableRepository.Delete.
//
// For obvious reasons, the source and destination URIs must both be writable.
//
// NOTE: this function should not be called except by an implementation of
// the Repository interface - using this for unknown URIs may break.
//
// Since: 2.0
func GenericMove(source fyne.URI, destination fyne.URI) error {
	// This looks a lot like GenericCopy(), but I duplicated the code
	// to avoid having to look up the repositories more than once.

	// Look up repositories for the source and destination.
	srcrepo, err := ForURI(source)
	if err != nil {
		return err
	}

	dstrepo, err := ForURI(destination)
	if err != nil {
		return err
	}

	// The source and destination must both be writable, since the source
	// is being deleted, which requires WritableRepository.
	destwrepo, ok := dstrepo.(WritableRepository)
	if !ok {
		return ErrOperationNotSupported
	}

	srcwrepo, ok := srcrepo.(WritableRepository)
	if !ok {
		return ErrOperationNotSupported
	}

	if listable, ok := srcrepo.(ListableRepository); ok {
		isParent, err := listable.CanList(source)
		if err == nil && isParent {
			if srcrepo != destwrepo { // cannot move between repositories
				return ErrOperationNotSupported
			}

			return genericCopyMoveListable(source, destination, srcrepo, true)
		}
	}

	// Create the reader and writer to perform the copy operation.
	srcReader, err := srcrepo.Reader(source)
	if err != nil {
		return err
	}

	dstWriter, err := destwrepo.Writer(destination)
	if err != nil {
		return err
	}
	defer dstWriter.Close()

	// Perform the copy.
	_, err = io.Copy(dstWriter, srcReader)
	if err != nil {
		return err
	}

	// Finally, delete the source only if the move finished without error.
	srcReader.Close()
	return srcwrepo.Delete(source)
}

func genericCopyMoveListable(source, destination fyne.URI, repo Repository, deleteSource bool) error {
	lister, ok1 := repo.(ListableRepository)
	mover, ok2 := repo.(MovableRepository)
	copier, ok3 := repo.(CopyableRepository)

	if !ok1 || (deleteSource && !ok2) || (!deleteSource && !ok3) {
		return ErrOperationNotSupported // cannot move a lister in a non-listable/movable repo
	}

	err := lister.CreateListable(destination)
	if err != nil {
		return err
	}

	list, err := lister.List(source)
	if err != nil {
		return err
	}
	for _, child := range list {
		newChild, _ := repo.(HierarchicalRepository).Child(destination, child.Name())
		if deleteSource {
			err = mover.Move(child, newChild)
		} else {
			err = copier.Copy(child, newChild)
		}
		if err != nil {
			return err
		}
	}

	if !deleteSource {
		return nil
	}
	// we know the repo is writable as well from earlier checks
	writer, _ := repo.(WritableRepository)
	return writer.Delete(source)
}

func genericDeleteAll(u fyne.URI, wrepo WritableRepository, lrepo ListableRepository) error {
	listable, err := lrepo.CanList(u)
	if err != nil {
		return err
	} else if !listable {
		return wrepo.Delete(u)
	}

	children, err := lrepo.List(u)
	if err != nil {
		return err
	} else if len(children) == 0 {
		return wrepo.Delete(u)
	}

	var folders []fyne.URI
	var files []fyne.URI
	for i := 0; i < len(children); i++ {
		listable, err = lrepo.CanList(children[i])
		if err != nil {
			return err
		}

		if listable {
			grandChildren, err := lrepo.List(children[i])
			if err != nil {
				return err
			}
			folders = append(folders, children[i])
			children = append(children, grandChildren...)
		} else {
			files = append(files, children[i])
		}
	}

	for i := len(files) - 1; i >= 0; i-- {
		err = wrepo.Delete(files[i])
		if err != nil {
			return err
		}
	}

	for i := len(folders) - 1; i >= 0; i-- {
		err = wrepo.Delete(folders[i])
		if err != nil {
			return err
		}
	}

	return wrepo.Delete(u)
}
