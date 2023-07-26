// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"bytes"

	"fyne.io/fyne/v2"
)

// BoolTree supports binding a tree of bool values.
//
// Since: 2.4
type BoolTree interface {
	DataTree

	Append(parent, id string, value bool) error
	Get() (map[string][]string, map[string]bool, error)
	GetValue(id string) (bool, error)
	Prepend(parent, id string, value bool) error
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
	t := &boundBoolTree{val: &map[string]bool{}}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

// BindBoolTree returns a bound tree of bool values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindBoolTree(ids *map[string][]string, v *map[string]bool) ExternalBoolTree {
	if v == nil {
		return NewBoolTree().(ExternalBoolTree)
	}

	t := &boundBoolTree{val: v, updateExternal: true}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindBoolTreeItem(v, leaf, t.updateExternal), leaf, parent)
		}
	}

	return t
}

type boundBoolTree struct {
	treeBase

	updateExternal bool
	val            *map[string]bool
}

func (t *boundBoolTree) Append(parent, id string, val bool) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundBoolTree) Get() (map[string][]string, map[string]bool, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundBoolTree) GetValue(id string) (bool, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return false, errOutOfBounds
}

func (t *boundBoolTree) Prepend(parent, id string, val bool) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundBoolTree) Reload() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doReload()
}

func (t *boundBoolTree) Set(ids map[string][]string, v map[string]bool) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ids = ids
	*t.val = v

	return t.doReload()
}

func (t *boundBoolTree) doReload() (retErr error) {
	updated := []string{}
	fire := false
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
		t.appendItem(bindBoolTreeItem(t.val, id, t.updateExternal), id, parentIDFor(id, t.ids))
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
	if fire {
		t.trigger()
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			item.(*boundExternalBoolTreeItem).lock.Lock()
			err = item.(*boundExternalBoolTreeItem).setIfChanged((*t.val)[id])
			item.(*boundExternalBoolTreeItem).lock.Unlock()
		} else {
			item.(*boundBoolTreeItem).lock.Lock()
			err = item.(*boundBoolTreeItem).doSet((*t.val)[id])
			item.(*boundBoolTreeItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *boundBoolTree) SetValue(id string, v bool) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(Bool).Set(v)
}

func bindBoolTreeItem(v *map[string]bool, id string, external bool) Bool {
	if external {
		ret := &boundExternalBoolTreeItem{old: (*v)[id]}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundBoolTreeItem{id: id, val: v}
}

type boundBoolTreeItem struct {
	base

	val *map[string]bool
	id  string
}

func (t *boundBoolTreeItem) Get() (bool, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return false, errOutOfBounds
}

func (t *boundBoolTreeItem) Set(val bool) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doSet(val)
}

func (t *boundBoolTreeItem) doSet(val bool) error {
	(*t.val)[t.id] = val

	t.trigger()
	return nil
}

type boundExternalBoolTreeItem struct {
	boundBoolTreeItem

	old bool
}

func (t *boundExternalBoolTreeItem) setIfChanged(val bool) error {
	if val == t.old {
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val

	t.trigger()
	return nil
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
	t := &boundBytesTree{val: &map[string][]byte{}}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

// BindBytesTree returns a bound tree of []byte values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindBytesTree(ids *map[string][]string, v *map[string][]byte) ExternalBytesTree {
	if v == nil {
		return NewBytesTree().(ExternalBytesTree)
	}

	t := &boundBytesTree{val: v, updateExternal: true}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindBytesTreeItem(v, leaf, t.updateExternal), leaf, parent)
		}
	}

	return t
}

type boundBytesTree struct {
	treeBase

	updateExternal bool
	val            *map[string][]byte
}

func (t *boundBytesTree) Append(parent, id string, val []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundBytesTree) Get() (map[string][]string, map[string][]byte, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundBytesTree) GetValue(id string) ([]byte, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return nil, errOutOfBounds
}

func (t *boundBytesTree) Prepend(parent, id string, val []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundBytesTree) Reload() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doReload()
}

func (t *boundBytesTree) Set(ids map[string][]string, v map[string][]byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ids = ids
	*t.val = v

	return t.doReload()
}

func (t *boundBytesTree) doReload() (retErr error) {
	updated := []string{}
	fire := false
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
		t.appendItem(bindBytesTreeItem(t.val, id, t.updateExternal), id, parentIDFor(id, t.ids))
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
	if fire {
		t.trigger()
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			item.(*boundExternalBytesTreeItem).lock.Lock()
			err = item.(*boundExternalBytesTreeItem).setIfChanged((*t.val)[id])
			item.(*boundExternalBytesTreeItem).lock.Unlock()
		} else {
			item.(*boundBytesTreeItem).lock.Lock()
			err = item.(*boundBytesTreeItem).doSet((*t.val)[id])
			item.(*boundBytesTreeItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *boundBytesTree) SetValue(id string, v []byte) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(Bytes).Set(v)
}

func bindBytesTreeItem(v *map[string][]byte, id string, external bool) Bytes {
	if external {
		ret := &boundExternalBytesTreeItem{old: (*v)[id]}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundBytesTreeItem{id: id, val: v}
}

type boundBytesTreeItem struct {
	base

	val *map[string][]byte
	id  string
}

func (t *boundBytesTreeItem) Get() ([]byte, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return nil, errOutOfBounds
}

func (t *boundBytesTreeItem) Set(val []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doSet(val)
}

func (t *boundBytesTreeItem) doSet(val []byte) error {
	(*t.val)[t.id] = val

	t.trigger()
	return nil
}

type boundExternalBytesTreeItem struct {
	boundBytesTreeItem

	old []byte
}

func (t *boundExternalBytesTreeItem) setIfChanged(val []byte) error {
	if bytes.Equal(val, t.old) {
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val

	t.trigger()
	return nil
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
	t := &boundFloatTree{val: &map[string]float64{}}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

// BindFloatTree returns a bound tree of float64 values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindFloatTree(ids *map[string][]string, v *map[string]float64) ExternalFloatTree {
	if v == nil {
		return NewFloatTree().(ExternalFloatTree)
	}

	t := &boundFloatTree{val: v, updateExternal: true}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindFloatTreeItem(v, leaf, t.updateExternal), leaf, parent)
		}
	}

	return t
}

type boundFloatTree struct {
	treeBase

	updateExternal bool
	val            *map[string]float64
}

func (t *boundFloatTree) Append(parent, id string, val float64) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundFloatTree) Get() (map[string][]string, map[string]float64, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundFloatTree) GetValue(id string) (float64, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return 0.0, errOutOfBounds
}

func (t *boundFloatTree) Prepend(parent, id string, val float64) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundFloatTree) Reload() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doReload()
}

func (t *boundFloatTree) Set(ids map[string][]string, v map[string]float64) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ids = ids
	*t.val = v

	return t.doReload()
}

func (t *boundFloatTree) doReload() (retErr error) {
	updated := []string{}
	fire := false
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
		t.appendItem(bindFloatTreeItem(t.val, id, t.updateExternal), id, parentIDFor(id, t.ids))
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
	if fire {
		t.trigger()
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			item.(*boundExternalFloatTreeItem).lock.Lock()
			err = item.(*boundExternalFloatTreeItem).setIfChanged((*t.val)[id])
			item.(*boundExternalFloatTreeItem).lock.Unlock()
		} else {
			item.(*boundFloatTreeItem).lock.Lock()
			err = item.(*boundFloatTreeItem).doSet((*t.val)[id])
			item.(*boundFloatTreeItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *boundFloatTree) SetValue(id string, v float64) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(Float).Set(v)
}

func bindFloatTreeItem(v *map[string]float64, id string, external bool) Float {
	if external {
		ret := &boundExternalFloatTreeItem{old: (*v)[id]}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundFloatTreeItem{id: id, val: v}
}

type boundFloatTreeItem struct {
	base

	val *map[string]float64
	id  string
}

func (t *boundFloatTreeItem) Get() (float64, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return 0.0, errOutOfBounds
}

func (t *boundFloatTreeItem) Set(val float64) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doSet(val)
}

func (t *boundFloatTreeItem) doSet(val float64) error {
	(*t.val)[t.id] = val

	t.trigger()
	return nil
}

type boundExternalFloatTreeItem struct {
	boundFloatTreeItem

	old float64
}

func (t *boundExternalFloatTreeItem) setIfChanged(val float64) error {
	if val == t.old {
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val

	t.trigger()
	return nil
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
	t := &boundIntTree{val: &map[string]int{}}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

// BindIntTree returns a bound tree of int values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindIntTree(ids *map[string][]string, v *map[string]int) ExternalIntTree {
	if v == nil {
		return NewIntTree().(ExternalIntTree)
	}

	t := &boundIntTree{val: v, updateExternal: true}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindIntTreeItem(v, leaf, t.updateExternal), leaf, parent)
		}
	}

	return t
}

type boundIntTree struct {
	treeBase

	updateExternal bool
	val            *map[string]int
}

func (t *boundIntTree) Append(parent, id string, val int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundIntTree) Get() (map[string][]string, map[string]int, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundIntTree) GetValue(id string) (int, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return 0, errOutOfBounds
}

func (t *boundIntTree) Prepend(parent, id string, val int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundIntTree) Reload() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doReload()
}

func (t *boundIntTree) Set(ids map[string][]string, v map[string]int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ids = ids
	*t.val = v

	return t.doReload()
}

func (t *boundIntTree) doReload() (retErr error) {
	updated := []string{}
	fire := false
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
		t.appendItem(bindIntTreeItem(t.val, id, t.updateExternal), id, parentIDFor(id, t.ids))
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
	if fire {
		t.trigger()
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			item.(*boundExternalIntTreeItem).lock.Lock()
			err = item.(*boundExternalIntTreeItem).setIfChanged((*t.val)[id])
			item.(*boundExternalIntTreeItem).lock.Unlock()
		} else {
			item.(*boundIntTreeItem).lock.Lock()
			err = item.(*boundIntTreeItem).doSet((*t.val)[id])
			item.(*boundIntTreeItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *boundIntTree) SetValue(id string, v int) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(Int).Set(v)
}

func bindIntTreeItem(v *map[string]int, id string, external bool) Int {
	if external {
		ret := &boundExternalIntTreeItem{old: (*v)[id]}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundIntTreeItem{id: id, val: v}
}

type boundIntTreeItem struct {
	base

	val *map[string]int
	id  string
}

func (t *boundIntTreeItem) Get() (int, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return 0, errOutOfBounds
}

func (t *boundIntTreeItem) Set(val int) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doSet(val)
}

func (t *boundIntTreeItem) doSet(val int) error {
	(*t.val)[t.id] = val

	t.trigger()
	return nil
}

type boundExternalIntTreeItem struct {
	boundIntTreeItem

	old int
}

func (t *boundExternalIntTreeItem) setIfChanged(val int) error {
	if val == t.old {
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val

	t.trigger()
	return nil
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
	t := &boundRuneTree{val: &map[string]rune{}}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

// BindRuneTree returns a bound tree of rune values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindRuneTree(ids *map[string][]string, v *map[string]rune) ExternalRuneTree {
	if v == nil {
		return NewRuneTree().(ExternalRuneTree)
	}

	t := &boundRuneTree{val: v, updateExternal: true}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindRuneTreeItem(v, leaf, t.updateExternal), leaf, parent)
		}
	}

	return t
}

type boundRuneTree struct {
	treeBase

	updateExternal bool
	val            *map[string]rune
}

func (t *boundRuneTree) Append(parent, id string, val rune) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundRuneTree) Get() (map[string][]string, map[string]rune, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundRuneTree) GetValue(id string) (rune, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return rune(0), errOutOfBounds
}

func (t *boundRuneTree) Prepend(parent, id string, val rune) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundRuneTree) Reload() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doReload()
}

func (t *boundRuneTree) Set(ids map[string][]string, v map[string]rune) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ids = ids
	*t.val = v

	return t.doReload()
}

func (t *boundRuneTree) doReload() (retErr error) {
	updated := []string{}
	fire := false
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
		t.appendItem(bindRuneTreeItem(t.val, id, t.updateExternal), id, parentIDFor(id, t.ids))
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
	if fire {
		t.trigger()
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			item.(*boundExternalRuneTreeItem).lock.Lock()
			err = item.(*boundExternalRuneTreeItem).setIfChanged((*t.val)[id])
			item.(*boundExternalRuneTreeItem).lock.Unlock()
		} else {
			item.(*boundRuneTreeItem).lock.Lock()
			err = item.(*boundRuneTreeItem).doSet((*t.val)[id])
			item.(*boundRuneTreeItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *boundRuneTree) SetValue(id string, v rune) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(Rune).Set(v)
}

func bindRuneTreeItem(v *map[string]rune, id string, external bool) Rune {
	if external {
		ret := &boundExternalRuneTreeItem{old: (*v)[id]}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundRuneTreeItem{id: id, val: v}
}

type boundRuneTreeItem struct {
	base

	val *map[string]rune
	id  string
}

func (t *boundRuneTreeItem) Get() (rune, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return rune(0), errOutOfBounds
}

func (t *boundRuneTreeItem) Set(val rune) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doSet(val)
}

func (t *boundRuneTreeItem) doSet(val rune) error {
	(*t.val)[t.id] = val

	t.trigger()
	return nil
}

type boundExternalRuneTreeItem struct {
	boundRuneTreeItem

	old rune
}

func (t *boundExternalRuneTreeItem) setIfChanged(val rune) error {
	if val == t.old {
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val

	t.trigger()
	return nil
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
	t := &boundStringTree{val: &map[string]string{}}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

// BindStringTree returns a bound tree of string values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindStringTree(ids *map[string][]string, v *map[string]string) ExternalStringTree {
	if v == nil {
		return NewStringTree().(ExternalStringTree)
	}

	t := &boundStringTree{val: v, updateExternal: true}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindStringTreeItem(v, leaf, t.updateExternal), leaf, parent)
		}
	}

	return t
}

type boundStringTree struct {
	treeBase

	updateExternal bool
	val            *map[string]string
}

func (t *boundStringTree) Append(parent, id string, val string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundStringTree) Get() (map[string][]string, map[string]string, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundStringTree) GetValue(id string) (string, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return "", errOutOfBounds
}

func (t *boundStringTree) Prepend(parent, id string, val string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundStringTree) Reload() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doReload()
}

func (t *boundStringTree) Set(ids map[string][]string, v map[string]string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ids = ids
	*t.val = v

	return t.doReload()
}

func (t *boundStringTree) doReload() (retErr error) {
	updated := []string{}
	fire := false
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
		t.appendItem(bindStringTreeItem(t.val, id, t.updateExternal), id, parentIDFor(id, t.ids))
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
	if fire {
		t.trigger()
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			item.(*boundExternalStringTreeItem).lock.Lock()
			err = item.(*boundExternalStringTreeItem).setIfChanged((*t.val)[id])
			item.(*boundExternalStringTreeItem).lock.Unlock()
		} else {
			item.(*boundStringTreeItem).lock.Lock()
			err = item.(*boundStringTreeItem).doSet((*t.val)[id])
			item.(*boundStringTreeItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *boundStringTree) SetValue(id string, v string) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(String).Set(v)
}

func bindStringTreeItem(v *map[string]string, id string, external bool) String {
	if external {
		ret := &boundExternalStringTreeItem{old: (*v)[id]}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundStringTreeItem{id: id, val: v}
}

type boundStringTreeItem struct {
	base

	val *map[string]string
	id  string
}

func (t *boundStringTreeItem) Get() (string, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return "", errOutOfBounds
}

func (t *boundStringTreeItem) Set(val string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doSet(val)
}

func (t *boundStringTreeItem) doSet(val string) error {
	(*t.val)[t.id] = val

	t.trigger()
	return nil
}

type boundExternalStringTreeItem struct {
	boundStringTreeItem

	old string
}

func (t *boundExternalStringTreeItem) setIfChanged(val string) error {
	if val == t.old {
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val

	t.trigger()
	return nil
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
	t := &boundURITree{val: &map[string]fyne.URI{}}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)
	return t
}

// BindURITree returns a bound tree of fyne.URI values, based on the contents of the passed values.
// The ids map specifies how each item relates to its parent (with id ""), with the values being in the v map.
// If your code changes the content of the maps this refers to you should call Reload() to inform the bindings.
//
// Since: 2.4
func BindURITree(ids *map[string][]string, v *map[string]fyne.URI) ExternalURITree {
	if v == nil {
		return NewURITree().(ExternalURITree)
	}

	t := &boundURITree{val: v, updateExternal: true}
	t.ids = make(map[string][]string)
	t.items = make(map[string]DataItem)

	for parent, children := range *ids {
		for _, leaf := range children {
			t.appendItem(bindURITreeItem(v, leaf, t.updateExternal), leaf, parent)
		}
	}

	return t
}

type boundURITree struct {
	treeBase

	updateExternal bool
	val            *map[string]fyne.URI
}

func (t *boundURITree) Append(parent, id string, val fyne.URI) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append(ids, id)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundURITree) Get() (map[string][]string, map[string]fyne.URI, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.ids, *t.val, nil
}

func (t *boundURITree) GetValue(id string) (fyne.URI, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if item, ok := (*t.val)[id]; ok {
		return item, nil
	}

	return fyne.URI(nil), errOutOfBounds
}

func (t *boundURITree) Prepend(parent, id string, val fyne.URI) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ids, ok := t.ids[parent]
	if !ok {
		ids = make([]string, 0)
	}

	t.ids[parent] = append([]string{id}, ids...)
	v := *t.val
	v[id] = val

	return t.doReload()
}

func (t *boundURITree) Reload() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doReload()
}

func (t *boundURITree) Set(ids map[string][]string, v map[string]fyne.URI) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ids = ids
	*t.val = v

	return t.doReload()
}

func (t *boundURITree) doReload() (retErr error) {
	updated := []string{}
	fire := false
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
		t.appendItem(bindURITreeItem(t.val, id, t.updateExternal), id, parentIDFor(id, t.ids))
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
	if fire {
		t.trigger()
	}

	for id, item := range t.items {
		var err error
		if t.updateExternal {
			item.(*boundExternalURITreeItem).lock.Lock()
			err = item.(*boundExternalURITreeItem).setIfChanged((*t.val)[id])
			item.(*boundExternalURITreeItem).lock.Unlock()
		} else {
			item.(*boundURITreeItem).lock.Lock()
			err = item.(*boundURITreeItem).doSet((*t.val)[id])
			item.(*boundURITreeItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (t *boundURITree) SetValue(id string, v fyne.URI) error {
	t.lock.Lock()
	(*t.val)[id] = v
	t.lock.Unlock()

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}
	return item.(URI).Set(v)
}

func bindURITreeItem(v *map[string]fyne.URI, id string, external bool) URI {
	if external {
		ret := &boundExternalURITreeItem{old: (*v)[id]}
		ret.val = v
		ret.id = id
		return ret
	}

	return &boundURITreeItem{id: id, val: v}
}

type boundURITreeItem struct {
	base

	val *map[string]fyne.URI
	id  string
}

func (t *boundURITreeItem) Get() (fyne.URI, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	v := *t.val
	if item, ok := v[t.id]; ok {
		return item, nil
	}

	return fyne.URI(nil), errOutOfBounds
}

func (t *boundURITreeItem) Set(val fyne.URI) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.doSet(val)
}

func (t *boundURITreeItem) doSet(val fyne.URI) error {
	(*t.val)[t.id] = val

	t.trigger()
	return nil
}

type boundExternalURITreeItem struct {
	boundURITreeItem

	old fyne.URI
}

func (t *boundExternalURITreeItem) setIfChanged(val fyne.URI) error {
	if compareURI(val, t.old) {
		return nil
	}
	(*t.val)[t.id] = val
	t.old = val

	t.trigger()
	return nil
}
