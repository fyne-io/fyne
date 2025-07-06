//go:build wasm

package repository

import (
	"context"
	"syscall/js"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/hack-pad/go-indexeddb/idb"
)

// fileSchemePrefix is used for when we need a hard-coded version of "idbfile://"
// for string processing
const idbfileSchemePrefix string = "idbfile://"

var (
	_ repository.Repository             = (*IndexDBRepository)(nil)
	_ repository.WritableRepository     = (*IndexDBRepository)(nil)
	_ repository.AppendableRepository   = (*IndexDBRepository)(nil)
	_ repository.HierarchicalRepository = (*IndexDBRepository)(nil)
	_ repository.ListableRepository     = (*IndexDBRepository)(nil)
	_ repository.MovableRepository      = (*IndexDBRepository)(nil)
	_ repository.CopyableRepository     = (*IndexDBRepository)(nil)
)

type IndexDBRepository struct {
	db *idb.Database
}

func NewIndexDBRepository() (*IndexDBRepository, error) {
	ctx := context.Background()
	req, err := idb.Global().Open(ctx, "files", 1, func(db *idb.Database, oldVer, newVer uint) error {
		metastore, err := db.CreateObjectStore("meta", idb.ObjectStoreOptions{})
		if err != nil {
			return err
		}
		_, err = metastore.CreateIndex("path", js.ValueOf(""), idb.IndexOptions{Unique: true})
		if err != nil {
			return err
		}
		_, err = metastore.CreateIndex("parent", js.ValueOf("parent"), idb.IndexOptions{})
		if err != nil {
			return err
		}

		datastore, err := db.CreateObjectStore("data", idb.ObjectStoreOptions{})
		if err != nil {
			return err
		}
		_, err = datastore.CreateIndex("path", js.ValueOf(""), idb.IndexOptions{Unique: true})
		return err
	})
	db, err := req.Await(ctx)
	if err != nil {
		return nil, err
	}

	if err := mkdir(db, "/", ""); err != nil {
		return nil, err
	}

	return &IndexDBRepository{db: db}, nil
}

func (r *IndexDBRepository) Exists(u fyne.URI) (bool, error) {
	p := u.Path()
	ctx := context.Background()
	txn, err := r.db.Transaction(idb.TransactionReadOnly, "meta")
	if err != nil {
		return false, err
	}
	store, err := txn.ObjectStore("meta")
	if err != nil {
		return false, err
	}
	req, err := store.CountKey(js.ValueOf(p))
	if err != nil {
		return false, err
	}
	n, err := req.Await(ctx)
	if err != nil {
		return false, err
	}

	return n != 0, nil
}

func get(db *idb.Database, s, p string) (js.Value, error) {
	ctx := context.Background()

	txn, err := db.Transaction(idb.TransactionReadOnly, s)
	if err != nil {
		return js.Undefined(), err
	}
	store, err := txn.ObjectStore(s)
	if err != nil {
		return js.Undefined(), err
	}
	req, err := store.Get(js.ValueOf(p))
	if err != nil {
		return js.Undefined(), err
	}

	v, err := req.Await(ctx)
	if err != nil {
		return js.Undefined(), err
	}

	return v, nil
}

func (r *IndexDBRepository) CanList(u fyne.URI) (bool, error) {
	p := u.Path()

	v, err := get(r.db, "meta", p)
	if err != nil {
		return false, err
	}

	if v.IsUndefined() {
		return false, nil
	}

	isDir := v.Get("isDir")
	if isDir.IsUndefined() {
		return false, nil
	}

	return isDir.Bool(), nil
}

func mkdir(db *idb.Database, dir, parent string) error {
	ctx := context.Background()
	txn, err := db.Transaction(idb.TransactionReadWrite, "meta")
	if err != nil {
		return err
	}

	store, err := txn.ObjectStore("meta")
	if err != nil {
		return err
	}

	f := map[string]any{
		"isDir":  true,
		"parent": parent,
	}
	req, err := store.PutKey(js.ValueOf(dir), js.ValueOf(f))
	if err != nil {
		return err
	}

	_, err = req.Await(ctx)
	return err
}

func (r *IndexDBRepository) CreateListable(u fyne.URI) error {
	pu, err := storage.Parent(u)
	if err != nil {
		return err
	}
	return mkdir(r.db, u.Path(), pu.Path())
}

func (r *IndexDBRepository) CanRead(u fyne.URI) (bool, error) {
	p := u.Path()

	v, err := get(r.db, "meta", p)
	if err != nil {
		return false, err
	}

	if v.IsUndefined() {
		return false, nil
	}

	return true, nil
}

func (r *IndexDBRepository) Destroy(scheme string) {
	// do nothing
}

func (r *IndexDBRepository) List(u fyne.URI) ([]fyne.URI, error) {
	p := u.Path()
	ctx := context.Background()
	txn, err := r.db.Transaction(idb.TransactionReadOnly, "meta")
	if err != nil {
		return nil, err
	}

	store, err := txn.ObjectStore("meta")
	if err != nil {
		return nil, err
	}

	idx, err := store.Index("parent")
	if err != nil {
		return nil, err
	}

	creq, err := idx.OpenCursorKey(js.ValueOf(p), idb.CursorNext)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	if err := creq.Iter(ctx, func(cwv *idb.CursorWithValue) error {
		k, err := cwv.PrimaryKey()
		if err != nil {
			return err
		}
		paths = append(paths, idbfileSchemePrefix+k.String())
		return nil
	}); err != nil {
		return nil, err
	}

	us := make([]fyne.URI, len(paths))
	for n, path := range paths {
		us[n], err = storage.ParseURI(path)
		if err != nil {
			return nil, err
		}
	}
	return us, nil
}

func (r *IndexDBRepository) CanWrite(u fyne.URI) (bool, error) {
	p := u.Path()
	v, err := get(r.db, "meta", p)
	if err != nil {
		return false, err
	}

	if v.IsUndefined() {
		return true, nil
	}

	isDir := v.Get("isDir")
	if isDir.IsUndefined() {
		return true, nil
	}

	return !isDir.Bool(), nil
}

func (r *IndexDBRepository) Delete(u fyne.URI) error {
	p := u.Path()
	ctx := context.Background()
	txn, err := r.db.Transaction(idb.TransactionReadWrite, "meta", "data")
	if err != nil {
		return err
	}

	metastore, err := txn.ObjectStore("meta")
	if err != nil {
		return err
	}

	metareq, err := metastore.Delete(js.ValueOf(p))
	if err != nil {
		return err
	}
	if err := metareq.Await(ctx); err != nil {
		return err
	}

	datastore, err := txn.ObjectStore("data")
	if err != nil {
		return err
	}
	datareq, err := datastore.Delete(js.ValueOf(p))
	if err != nil {
		return err
	}
	return datareq.Await(ctx)
}

func (r *IndexDBRepository) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	pu, err := storage.Parent(u)
	if err != nil {
		return nil, err
	}

	return &idbfile{
		db:     r.db,
		path:   u.Path(),
		parent: pu.Path(),
	}, nil
}

func (r *IndexDBRepository) Writer(u fyne.URI) (fyne.URIWriteCloser, error) {
	pu, err := storage.Parent(u)
	if err != nil {
		return nil, err
	}

	return &idbfile{
		db:       r.db,
		path:     u.Path(),
		parent:   pu.Path(),
		truncate: true,
	}, nil
}

func (r *IndexDBRepository) Appender(u fyne.URI) (fyne.URIWriteCloser, error) {
	pu, err := storage.Parent(u)
	if err != nil {
		return nil, err
	}

	return &idbfile{
		db:     r.db,
		path:   u.Path(),
		parent: pu.Path(),
		add:    true,
	}, nil
}

func (r *IndexDBRepository) Copy(src, dst fyne.URI) error {
	return repository.GenericCopy(src, dst)
}

func (r *IndexDBRepository) Move(source, destination fyne.URI) error {
	return repository.GenericMove(source, destination)
}

func (r *IndexDBRepository) Child(u fyne.URI, component string) (fyne.URI, error) {
	return repository.GenericChild(u, component)
}

func (r *IndexDBRepository) Parent(u fyne.URI) (fyne.URI, error) {
	return repository.GenericParent(u)
}
