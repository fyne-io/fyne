//go:build wasm

package repository

import (
	"context"
	"fmt"
	"io"
	"syscall/js"

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

var _ fyne.URIReadCloser = (*idbfile)(nil)
var _ fyne.URIWriteCloser = (*idbfile)(nil)

type idbfile struct {
	db          *idb.Database
	path        string
	parent      string
	isDir       bool
	truncate    bool
	isTruncated bool
	parts       []interface{}
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

	if f.truncate && !f.isTruncated || f.add && !f.isAdding {
		store, err := f.rwstore("meta")
		if err != nil {
			return 0, err
		}

		f := map[string]interface{}{
			"parent": f.parent,
		}
		req, err := store.PutKey(js.ValueOf(p), js.ValueOf(f))
		if err != nil {
			return 0, err
		}

		if _, err := req.Await(ctx); err != nil {
			return 0, err
		}
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
	}

	if f.add && !f.isAdding {
		store, err := f.rwstore("data")
		if err != nil {
			return 0, err
		}
		req, err := store.Get(js.ValueOf(f.path))
		if err != nil {
			return 0, err
		}
		b, err := req.Await(ctx)
		if err != nil {
			return 0, err
		}

		bytes := b.Call("bytes")
		outch := make(chan js.Value)
		bytes.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			outch <- args[0]
			return nil
		}))
		a := <-outch
		f.parts = []interface{}{a}
		f.isAdding = true
	}

	a := uint8Array.New(len(data))
	n := js.CopyBytesToJS(a, data)
	f.parts = append(f.parts, a)
	b := blob.New(js.ValueOf(f.parts))

	store, err := f.rwstore("data")
	if err != nil {
		return 0, err
	}
	req, err := store.PutKey(js.ValueOf(p), b)
	if err != nil {
		return 0, err
	}

	_, err = req.Await(ctx)
	return n, err
}

func (f *idbfile) Read(data []byte) (int, error) {
	ctx := context.Background()
	txn, err := f.db.Transaction(idb.TransactionReadOnly, "data")
	if err != nil {
		return 0, err
	}

	store, err := txn.ObjectStore("data")
	if err != nil {
		return 0, err
	}
	req, err := store.Get(js.ValueOf(f.path))
	if err != nil {
		return 0, err
	}

	b, err := req.Await(ctx)
	if err != nil {
		return 0, err
	}

	if b.IsUndefined() {
		return 0, fmt.Errorf("idbfile undefined")
	}

	bytes := b.Call("bytes")
	outch := make(chan js.Value)
	bytes.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		outch <- args[0]
		return nil
	}))
	out := <-outch

	return js.CopyBytesToGo(data, out), io.EOF
}
