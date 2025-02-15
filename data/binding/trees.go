package binding

import (
	"bytes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// DataTreeRootID const is the value used as ID for the root of any tree binding.
const DataTreeRootID = ""

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
type BoolTree interface {
	DataTree

	Append(parent, id string, value bool) error
	Get() (map[string][]string, map[string]bool, error)
	GetValue(id string) (bool, error)
	Prepend(parent, id string, value bool) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]bool) error
	SetValue(id string, value bool) error
}

// ExternalBoolTree supports binding a tree of bool values from an external variable.
//
// Since: 2.4
type ExternalBoolTree interface {
	BoolTree

	Reload() error
}

// NewBoolTree returns a bindable tree of bool values.
//
// Since: 2.4
func NewBoolTree() BoolTree {
	return newTreeComparable[bool]()
}

// BindBoolTree returns a bound tree of bool values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindBoolTree(ids *map[string][]string, v *map[string]bool) ExternalBoolTree {
	return bindTreeComparable(ids, v)
}

// BytesTree supports binding a tree of []byte values.
//
// Since: 2.4
type BytesTree interface {
	DataTree

	Append(parent, id string, value []byte) error
	Get() (map[string][]string, map[string][]byte, error)
	GetValue(id string) ([]byte, error)
	Prepend(parent, id string, value []byte) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string][]byte) error
	SetValue(id string, value []byte) error
}

// ExternalBytesTree supports binding a tree of []byte values from an external variable.
//
// Since: 2.4
type ExternalBytesTree interface {
	BytesTree

	Reload() error
}

// NewBytesTree returns a bindable tree of []byte values.
//
// Since: 2.4
func NewBytesTree() BytesTree {
	return newTree(bytes.Equal)
}

// BindBytesTree returns a bound tree of []byte values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindBytesTree(ids *map[string][]string, v *map[string][]byte) ExternalBytesTree {
	return bindTree(ids, v, bytes.Equal)
}

// FloatTree supports binding a tree of float64 values.
//
// Since: 2.4
type FloatTree interface {
	DataTree

	Append(parent, id string, value float64) error
	Get() (map[string][]string, map[string]float64, error)
	GetValue(id string) (float64, error)
	Prepend(parent, id string, value float64) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]float64) error
	SetValue(id string, value float64) error
}

// ExternalFloatTree supports binding a tree of float64 values from an external variable.
//
// Since: 2.4
type ExternalFloatTree interface {
	FloatTree

	Reload() error
}

// NewFloatTree returns a bindable tree of float64 values.
//
// Since: 2.4
func NewFloatTree() FloatTree {
	return newTreeComparable[float64]()
}

// BindFloatTree returns a bound tree of float64 values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindFloatTree(ids *map[string][]string, v *map[string]float64) ExternalFloatTree {
	return bindTreeComparable(ids, v)
}

// IntTree supports binding a tree of int values.
//
// Since: 2.4
type IntTree interface {
	DataTree

	Append(parent, id string, value int) error
	Get() (map[string][]string, map[string]int, error)
	GetValue(id string) (int, error)
	Prepend(parent, id string, value int) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]int) error
	SetValue(id string, value int) error
}

// ExternalIntTree supports binding a tree of int values from an external variable.
//
// Since: 2.4
type ExternalIntTree interface {
	IntTree

	Reload() error
}

// NewIntTree returns a bindable tree of int values.
//
// Since: 2.4
func NewIntTree() IntTree {
	return newTreeComparable[int]()
}

// BindIntTree returns a bound tree of int values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindIntTree(ids *map[string][]string, v *map[string]int) ExternalIntTree {
	return bindTreeComparable(ids, v)
}

// RuneTree supports binding a tree of rune values.
//
// Since: 2.4
type RuneTree interface {
	DataTree

	Append(parent, id string, value rune) error
	Get() (map[string][]string, map[string]rune, error)
	GetValue(id string) (rune, error)
	Prepend(parent, id string, value rune) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]rune) error
	SetValue(id string, value rune) error
}

// ExternalRuneTree supports binding a tree of rune values from an external variable.
//
// Since: 2.4
type ExternalRuneTree interface {
	RuneTree

	Reload() error
}

// NewRuneTree returns a bindable tree of rune values.
//
// Since: 2.4
func NewRuneTree() RuneTree {
	return newTreeComparable[rune]()
}

// BindRuneTree returns a bound tree of rune values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindRuneTree(ids *map[string][]string, v *map[string]rune) ExternalRuneTree {
	return bindTreeComparable(ids, v)
}

// StringTree supports binding a tree of string values.
//
// Since: 2.4
type StringTree interface {
	DataTree

	Append(parent, id string, value string) error
	Get() (map[string][]string, map[string]string, error)
	GetValue(id string) (string, error)
	Prepend(parent, id string, value string) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]string) error
	SetValue(id string, value string) error
}

// ExternalStringTree supports binding a tree of string values from an external variable.
//
// Since: 2.4
type ExternalStringTree interface {
	StringTree

	Reload() error
}

// NewStringTree returns a bindable tree of string values.
//
// Since: 2.4
func NewStringTree() StringTree {
	return newTreeComparable[string]()
}

// BindStringTree returns a bound tree of string values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindStringTree(ids *map[string][]string, v *map[string]string) ExternalStringTree {
	return bindTreeComparable(ids, v)
}

// UntypedTree supports binding a tree of any values.
//
// Since: 2.5
type UntypedTree interface {
	DataTree

	Append(parent, id string, value any) error
	Get() (map[string][]string, map[string]any, error)
	GetValue(id string) (any, error)
	Prepend(parent, id string, value any) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]any) error
	SetValue(id string, value any) error
}

// ExternalUntypedTree supports binding a tree of any values from an external variable.
//
// Since: 2.5
type ExternalUntypedTree interface {
	UntypedTree

	Reload() error
}

// NewUntypedTree returns a bindable tree of any values.
//
// Since: 2.5
func NewUntypedTree() UntypedTree {
	return newTree(func(a1, a2 any) bool { return a1 == a2 })
}

// BindUntypedTree returns a bound tree of any values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindUntypedTree(ids *map[string][]string, v *map[string]any) ExternalUntypedTree {
	return bindTree(ids, v, func(a1, a2 any) bool { return a1 == a2 })
}

// URITree supports binding a tree of fyne.URI values.
//
// Since: 2.4
type URITree interface {
	DataTree

	Append(parent, id string, value fyne.URI) error
	Get() (map[string][]string, map[string]fyne.URI, error)
	GetValue(id string) (fyne.URI, error)
	Prepend(parent, id string, value fyne.URI) error
	Remove(id string) error
	Set(ids map[string][]string, values map[string]fyne.URI) error
	SetValue(id string, value fyne.URI) error
}

// ExternalURITree supports binding a tree of fyne.URI values from an external variable.
//
// Since: 2.4
type ExternalURITree interface {
	URITree

	Reload() error
}

// NewURITree returns a bindable tree of fyne.URI values.
//
// Since: 2.4
func NewURITree() URITree {
	return newTree(storage.EqualURI)
}

// BindURITree returns a bound tree of fyne.URI values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindURITree(ids *map[string][]string, v *map[string]fyne.URI) ExternalURITree {
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
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

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

func newTreeComparable[T bool | float64 | int | rune | string]() *boundTree[T] {
	return newTree(func(t1, t2 T) bool { return t1 == t2 })
}

func bindTree[T any](ids *map[string][]string, v *map[string]T, comparator func(T, T) bool) *boundTree[T] {
	if v == nil {
		return newTree[T](comparator)
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

func bindTreeComparable[T bool | float64 | int | rune | string](ids *map[string][]string, v *map[string]T) *boundTree[T] {
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
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
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
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
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
	return
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
