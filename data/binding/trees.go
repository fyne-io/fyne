package binding

import (
	"bytes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// DataTreeRootID const is the value used as ID for the root of any tree binding.
const DataTreeRootID = ""

// Tree supports binding a tree of values with type T.
//
// Since: 2.7
type Tree[T any] interface {
	DataTree

	Append(parent, id string, value T) error
	Get() (map[string][]string, map[string]T, error)
	GetValue(id string) (T, error)
	Prepend(parent, id string, value T) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]T) error
	SetValue(id string, value T) error
}

// ExternalTree supports binding a tree of values, of type T, from an external variable.
//
// Since: 2.7
type ExternalTree[T any] interface {
	Tree[T]

	Reload() error
}

// NewTree returns a bindable tree of values with type T.
//
// Since: 2.7
func NewTree[T any](comparator func(T, T) bool) Tree[T] {
	return newTree[T](comparator)
}

// BindTree returns a bound tree of values with type T, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.7
func BindTree[T any](ids *map[string][]string, v *map[string]T, comparator func(T, T) bool) ExternalTree[T] {
	return bindTree(ids, v, comparator)
}

// DataTree is the base interface for all bindable data trees.
//
// Since: 2.4
type DataTree interface {
	DataItem
	GetItem(id string) (DataItem, error)
	ChildIDs(string) []string
}

// BoolTree supports binding a tree of bool values.
//
// Since: 2.4
type BoolTree = Tree[bool]

// ExternalBoolTree supports binding a tree of bool values from an external variable.
//
// Since: 2.4
type ExternalBoolTree = ExternalTree[bool]

// NewBoolTree returns a bindable tree of bool values.
//
// Since: 2.4
func NewBoolTree() Tree[bool] {
	return newTreeComparable[bool]()
}

// BindBoolTree returns a bound tree of bool values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindBoolTree(ids *map[string][]string, v *map[string]bool) ExternalTree[bool] {
	return bindTreeComparable(ids, v)
}

// BytesTree supports binding a tree of []byte values.
//
// Since: 2.4
type BytesTree = Tree[[]byte]

// ExternalBytesTree supports binding a tree of []byte values from an external variable.
//
// Since: 2.4
type ExternalBytesTree = ExternalTree[[]byte]

// NewBytesTree returns a bindable tree of []byte values.
//
// Since: 2.4
func NewBytesTree() Tree[[]byte] {
	return newTree(bytes.Equal)
}

// BindBytesTree returns a bound tree of []byte values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindBytesTree(ids *map[string][]string, v *map[string][]byte) ExternalTree[[]byte] {
	return bindTree(ids, v, bytes.Equal)
}

// FloatTree supports binding a tree of float64 values.
//
// Since: 2.4
type FloatTree = Tree[float64]

// ExternalFloatTree supports binding a tree of float64 values from an external variable.
//
// Since: 2.4
type ExternalFloatTree = ExternalTree[float64]

// NewFloatTree returns a bindable tree of float64 values.
//
// Since: 2.4
func NewFloatTree() Tree[float64] {
	return newTreeComparable[float64]()
}

// BindFloatTree returns a bound tree of float64 values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindFloatTree(ids *map[string][]string, v *map[string]float64) ExternalTree[float64] {
	return bindTreeComparable(ids, v)
}

// IntTree supports binding a tree of int values.
//
// Since: 2.4
type IntTree = Tree[int]

// ExternalIntTree supports binding a tree of int values from an external variable.
//
// Since: 2.4
type ExternalIntTree = ExternalTree[int]

// NewIntTree returns a bindable tree of int values.
//
// Since: 2.4
func NewIntTree() Tree[int] {
	return newTreeComparable[int]()
}

// BindIntTree returns a bound tree of int values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindIntTree(ids *map[string][]string, v *map[string]int) ExternalTree[int] {
	return bindTreeComparable(ids, v)
}

// RuneTree supports binding a tree of rune values.
//
// Since: 2.4
type RuneTree = Tree[rune]

// ExternalRuneTree supports binding a tree of rune values from an external variable.
//
// Since: 2.4
type ExternalRuneTree = ExternalTree[rune]

// NewRuneTree returns a bindable tree of rune values.
//
// Since: 2.4
func NewRuneTree() Tree[rune] {
	return newTreeComparable[rune]()
}

// BindRuneTree returns a bound tree of rune values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindRuneTree(ids *map[string][]string, v *map[string]rune) ExternalTree[rune] {
	return bindTreeComparable(ids, v)
}

// StringTree supports binding a tree of string values.
//
// Since: 2.4
type StringTree = Tree[string]

// ExternalStringTree supports binding a tree of string values from an external variable.
//
// Since: 2.4
type ExternalStringTree = ExternalTree[string]

// NewStringTree returns a bindable tree of string values.
//
// Since: 2.4
func NewStringTree() Tree[string] {
	return newTreeComparable[string]()
}

// BindStringTree returns a bound tree of string values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindStringTree(ids *map[string][]string, v *map[string]string) ExternalTree[string] {
	return bindTreeComparable(ids, v)
}

// UntypedTree supports binding a tree of any values.
//
// Since: 2.5
type UntypedTree = Tree[any]

// ExternalUntypedTree supports binding a tree of any values from an external variable.
//
// Since: 2.5
type ExternalUntypedTree = ExternalTree[any]

// NewUntypedTree returns a bindable tree of any values.
//
// Since: 2.5
func NewUntypedTree() Tree[any] {
	return newTree(func(a1, a2 any) bool { return a1 == a2 })
}

// BindUntypedTree returns a bound tree of any values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindUntypedTree(ids *map[string][]string, v *map[string]any) ExternalTree[any] {
	return bindTree(ids, v, func(a1, a2 any) bool { return a1 == a2 })
}

// URITree supports binding a tree of fyne.URI values.
//
// Since: 2.4
type URITree = Tree[fyne.URI]

// ExternalURITree supports binding a tree of fyne.URI values from an external variable.
//
// Since: 2.4
type ExternalURITree = ExternalTree[fyne.URI]

// NewURITree returns a bindable tree of fyne.URI values.
//
// Since: 2.4
func NewURITree() Tree[fyne.URI] {
	return newTree(storage.EqualURI)
}

// BindURITree returns a bound tree of fyne.URI values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindURITree(ids *map[string][]string, v *map[string]fyne.URI) ExternalTree[fyne.URI] {
	return bindTree(ids, v, storage.EqualURI)
}

type treeBase struct {
	base

	ids   map[string][]string
	items map[string]DataItem
}

// GetItem returns the DataItem at the specified id.
func (t *treeBase) GetItem(id string) (DataItem, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := t.items[id]; ok {
		return item, nil
	}

	return nil, errOutOfBounds
}

// ChildIDs returns the ordered IDs of items in this data tree that are children of the specified ID.
func (t *treeBase) ChildIDs(id string) []string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if ids, ok := t.ids[id]; ok {
		return ids
	}

	return []string{}
}

func (t *treeBase) appendItem(i DataItem, id, parent string) {
	t.items[id] = i

	ids := t.ids[parent]
	for _, in := range ids {
		if in == id {
			return
		}
	}
	t.ids[parent] = append(ids, id)
}

func (t *treeBase) deleteItem(id, parent string) {
	delete(t.items, id)

	ids, ok := t.ids[parent]
	if !ok {
		return
	}

	off := -1
	for i, id2 := range ids {
		if id2 == id {
			off = i
			break
		}
	}
	if off == -1 {
		return
	}
	t.ids[parent] = append(ids[:off], ids[off+1:]...)
}

func parentIDFor(id string, ids map[string][]string) string {
	for parent, list := range ids {
		for _, child := range list {
			if child == id {
				return parent
			}
		}
	}

	return ""
}

func newTree[T any](comparator func(T, T) bool) *boundTree[T] {
	t := &boundTree[T]{val: &map[string]T{}, comparator: comparator}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

func newTreeComparable[T comparable]() *boundTree[T] {
	return newTree(func(t1, t2 T) bool { return t1 == t2 })
}

func bindTree[T any](ids *map[string][]string, v *map[string]T, comparator func(T, T) bool) *boundTree[T] {
	if v == nil {
		return newTree(comparator)
	}

	t := &boundTree[T]{val: v, updateExternal: true, comparator: comparator}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindTreeItem(v, leaf, t.updateExternal, t.comparator), leaf, parent)
		}
	}

	return t
}

func bindTreeComparable[T comparable](ids *map[string][]string, v *map[string]T) *boundTree[T] {
	return bindTree(ids, v, func(t1, t2 T) bool { return t1 == t2 })
}

type boundTree[T any] struct {
	treeBase

	comparator     func(T, T) bool
	val            *map[string]T
	updateExternal bool
}

func (t *boundTree[T]) Append(parent, id string, val T) error {
	t.lock.Lock()

	t.ids[parent] = append(t.ids[parent], id)
	v := *t.val
	v[id] = val

	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) Get() (map[string][]string, map[string]T, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundTree[T]) GetValue(id string) (T, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return *new(T), errOutOfBounds
}

func (t *boundTree[T]) Prepend(parent, id string, val T) error {
	t.lock.Lock()

	t.ids[parent] = append([]string{id}, t.ids[parent]...)
	v := *t.val
	v[id] = val

	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) Remove(id string) error {
	t.lock.Lock()
	t.removeChildren(id)
	delete(t.ids, id)
	v := *t.val
	delete(v, id)

	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) removeChildren(id string) {
	for _, cid := range t.ids[id] {
		t.removeChildren(cid)

		delete(t.ids, cid)
		v := *t.val
		delete(v, cid)
	}
}

func (t *boundTree[T]) Reload() error {
	t.lock.Lock()
	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) Set(ids map[string][]string, v map[string]T) error {
	t.lock.Lock()
	t.ids = ids
	*t.val = v

	trigger, err := t.doReload()
	t.lock.Unlock()

	if trigger {
		t.trigger()
	}

	return err
}

func (t *boundTree[T]) doReload() (fire bool, retErr error) {
	updated := []string{}
	for id := range *t.val {
		found := false
		for child := range t.items {
			if child == id { // update existing
				updated = append(updated, id)
				found = true
				break
			}
		}
		if found {
			continue
		}

		// append new
		t.appendItem(bindTreeItem(t.val, id, t.updateExternal, t.comparator), id, parentIDFor(id, t.ids))
		updated = append(updated, id)
		fire = true
	}

	for id := range t.items {
		remove := true
		for _, done := range updated {
			if done == id {
				remove = false
				break
			}
		}

		if remove { // remove item no longer present
			fire = true
			t.deleteItem(id, parentIDFor(id, t.ids))
		}
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			err = item.(*boundExternalTreeItem[T]).setIfChanged((*t.val)[id])
		} else {
			err = item.(*boundTreeItem[T]).doSet((*t.val)[id])
		}
		if err != nil {
			retErr = err
		}
	}
	return fire, retErr
}

func (t *boundTree[T]) SetValue(id string, v T) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(Item[T]).Set(v)
}

func bindTreeItem[T any](v *map[string]T, id string, external bool, comparator func(T, T) bool) Item[T] {
	if external {
		ret := &boundExternalTreeItem[T]{old: (*v)[id], comparator: comparator}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundTreeItem[T]{id: id, val: v}
}

type boundTreeItem[T any] struct {
	base

	val *map[string]T
	id  string
}

func (t *boundTreeItem[T]) Get() (T, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return *new(T), errOutOfBounds
}

func (t *boundTreeItem[T]) Set(val T) error {
	return t.doSet(val)
}

func (t *boundTreeItem[T]) doSet(val T) error {
	t.lock.Lock()
	(*t.val)[t.id] = val
	t.lock.Unlock()

	t.trigger()
	return nil
}

type boundExternalTreeItem[T any] struct {
	boundTreeItem[T]

	comparator func(T, T) bool
	old        T
}

func (t *boundExternalTreeItem[T]) setIfChanged(val T) error {
	t.lock.Lock()
	if t.comparator(val, t.old) {
		t.lock.Unlock()
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val
	t.lock.Unlock()

	t.trigger()
	return nil
}
