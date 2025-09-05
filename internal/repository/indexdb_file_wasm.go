//go:build wasm

package repository

import (
	"context"
	"fmt"
	"io"
	"syscall/js"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"

	"github.com/hack-pad/go-indexeddb/idb"
)

var (
	blob       js.Value
	uint8Array js.Value
)

func init() {
	blob = js.Global().Get("Blob")
	uint8Array = js.Global().Get("Uint8Array")
}

var (
	_ fyne.URIReadCloser  = (*idbfile)(nil)
	_ fyne.URIWriteCloser = (*idbfile)(nil)
)

type idbfile struct {
	db          *idb.Database
	path        string
	parent      string
	isDir       bool
	truncate    bool
	isTruncated bool
	parts       []any
	add         bool
	isAdding    bool
}

func (f *idbfile) Close() error {
	return nil
}

func (f *idbfile) URI() fyne.URI {
	u, _ := storage.ParseURI(idbfileSchemePrefix + f.path)
	return u
}

func (f *idbfile) rwstore(name string) (*idb.ObjectStore, error) {
	txn, err := f.db.Transaction(idb.TransactionReadWrite, name)
	if err != nil {
		return nil, err
	}
	store, err := txn.ObjectStore(name)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (f *idbfile) Write(data []byte) (int, error) {
	p := f.path
	ctx := context.Background()

	m := map[string]any{
		"parent": f.parent,
		"size":   0,
		"ctime":  0,
		"mtime":  0,
	}

	if f.truncate && !f.isTruncated {
		store, err := f.rwstore("data")
		if err != nil {
			return 0, err
		}
		delreq, err := store.Delete(js.ValueOf(p))
		if err != nil {
			return 0, err
		}
		if err := delreq.Await(ctx); err != nil {
			return 0, err
		}
		f.isTruncated = true

		m["ctime"] = time.Now().UnixMilli()
		m["mtime"] = m["ctime"]
	}

	if f.add && !f.isAdding {
		b, err := get(f.db, "data", f.path)
		if err != nil {
			return 0, err
		}

		f.parts = []any{getBytes(b)}
		f.isAdding = true

		meta, err := get(f.db, "meta", f.path)
		if err != nil {
			return 0, err
		}

		m["ctime"] = meta.Get("ctime").Int()
		m["mtime"] = time.Now().UnixMilli()
	}

	a := uint8Array.New(len(data))
	n := js.CopyBytesToJS(a, data)
	f.parts = append(f.parts, a)
	b := blob.New(js.ValueOf(f.parts))

	m["size"] = b.Get("size").Int()

	metastore, err := f.rwstore("meta")
	if err != nil {
		return 0, err
	}

	metareq, err := metastore.PutKey(js.ValueOf(p), js.ValueOf(m))
	if err != nil {
		return 0, err
	}

	store, err := f.rwstore("data")
	if err != nil {
		return 0, err
	}
	req, err := store.PutKey(js.ValueOf(p), b)
	if err != nil {
		return 0, err
	}

	if _, err := metareq.Await(ctx); err != nil {
		return 0, err
	}

	_, err = req.Await(ctx)
	return n, err
}

func getBytes(b js.Value) js.Value {
	outch := make(chan js.Value)
	send := js.FuncOf(func(this js.Value, args []js.Value) any {
		outch <- args[0]
		return nil
	})
	defer send.Release()

	b.Call("arrayBuffer").Call("then", send)
	buf := <-outch
	return uint8Array.New(buf)
}

func (f *idbfile) Read(data []byte) (int, error) {
	b, err := get(f.db, "data", f.path)
	if err != nil {
		return 0, err
	}

	if b.IsUndefined() {
		return 0, fmt.Errorf("idbfile undefined")
	}

	if !b.InstanceOf(blob) {
		return 0, fmt.Errorf("returned object not of type blob")
	}

	return js.CopyBytesToGo(data, getBytes(b)), io.EOF
}
