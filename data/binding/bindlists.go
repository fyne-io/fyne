// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"bytes"

	"fyne.io/fyne/v2"
)

// BoolList supports binding a list of bool values.
//
// Since: 2.0
type BoolList interface {
	DataList

	Append(value bool) error
	Get() ([]bool, error)
	GetValue(index int) (bool, error)
	Prepend(value bool) error
	Set(list []bool) error
	SetValue(index int, value bool) error
}

// ExternalBoolList supports binding a list of bool values from an external variable.
//
// Since: 2.0
type ExternalBoolList interface {
	BoolList

	Reload() error
}

// NewBoolList returns a bindable list of bool values.
//
// Since: 2.0
func NewBoolList() BoolList {
	return &boundBoolList{val: &[]bool{}}
}

// BindBoolList returns a bound list of bool values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindBoolList(v *[]bool) ExternalBoolList {
	if v == nil {
		return NewBoolList().(ExternalBoolList)
	}

	b := &boundBoolList{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bindBoolListItem(v, i, b.updateExternal))
	}

	return b
}

type boundBoolList struct {
	listBase

	updateExternal bool
	val            *[]bool
}

func (l *boundBoolList) Append(val bool) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *boundBoolList) Get() ([]bool, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundBoolList) GetValue(i int) (bool, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return false, errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundBoolList) Prepend(val bool) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([]bool{val}, *l.val...)

	return l.doReload()
}

func (l *boundBoolList) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

func (l *boundBoolList) Set(v []bool) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *boundBoolList) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindBoolListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternalBoolListItem).lock.Lock()
			err = item.(*boundExternalBoolListItem).setIfChanged((*l.val)[i])
			item.(*boundExternalBoolListItem).lock.Unlock()
		} else {
			item.(*boundBoolListItem).lock.Lock()
			err = item.(*boundBoolListItem).doSet((*l.val)[i])
			item.(*boundBoolListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundBoolList) SetValue(i int, v bool) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Bool).Set(v)
}

func bindBoolListItem(v *[]bool, i int, external bool) Bool {
	if external {
		ret := &boundExternalBoolListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &boundBoolListItem{val: v, index: i}
}

type boundBoolListItem struct {
	base

	val   *[]bool
	index int
}

func (b *boundBoolListItem) Get() (bool, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return false, errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundBoolListItem) Set(val bool) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *boundBoolListItem) doSet(val bool) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternalBoolListItem struct {
	boundBoolListItem

	old bool
}

func (b *boundExternalBoolListItem) setIfChanged(val bool) error {
	if val == b.old {
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}

// BytesList supports binding a list of []byte values.
//
// Since: 2.2
type BytesList interface {
	DataList

	Append(value []byte) error
	Get() ([][]byte, error)
	GetValue(index int) ([]byte, error)
	Prepend(value []byte) error
	Set(list [][]byte) error
	SetValue(index int, value []byte) error
}

// ExternalBytesList supports binding a list of []byte values from an external variable.
//
// Since: 2.2
type ExternalBytesList interface {
	BytesList

	Reload() error
}

// NewBytesList returns a bindable list of []byte values.
//
// Since: 2.2
func NewBytesList() BytesList {
	return &boundBytesList{val: &[][]byte{}}
}

// BindBytesList returns a bound list of []byte values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.2
func BindBytesList(v *[][]byte) ExternalBytesList {
	if v == nil {
		return NewBytesList().(ExternalBytesList)
	}

	b := &boundBytesList{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bindBytesListItem(v, i, b.updateExternal))
	}

	return b
}

type boundBytesList struct {
	listBase

	updateExternal bool
	val            *[][]byte
}

func (l *boundBytesList) Append(val []byte) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *boundBytesList) Get() ([][]byte, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundBytesList) GetValue(i int) ([]byte, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return nil, errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundBytesList) Prepend(val []byte) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([][]byte{val}, *l.val...)

	return l.doReload()
}

func (l *boundBytesList) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

func (l *boundBytesList) Set(v [][]byte) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *boundBytesList) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindBytesListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternalBytesListItem).lock.Lock()
			err = item.(*boundExternalBytesListItem).setIfChanged((*l.val)[i])
			item.(*boundExternalBytesListItem).lock.Unlock()
		} else {
			item.(*boundBytesListItem).lock.Lock()
			err = item.(*boundBytesListItem).doSet((*l.val)[i])
			item.(*boundBytesListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundBytesList) SetValue(i int, v []byte) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Bytes).Set(v)
}

func bindBytesListItem(v *[][]byte, i int, external bool) Bytes {
	if external {
		ret := &boundExternalBytesListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &boundBytesListItem{val: v, index: i}
}

type boundBytesListItem struct {
	base

	val   *[][]byte
	index int
}

func (b *boundBytesListItem) Get() ([]byte, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return nil, errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundBytesListItem) Set(val []byte) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *boundBytesListItem) doSet(val []byte) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternalBytesListItem struct {
	boundBytesListItem

	old []byte
}

func (b *boundExternalBytesListItem) setIfChanged(val []byte) error {
	if bytes.Equal(val, b.old) {
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}

// FloatList supports binding a list of float64 values.
//
// Since: 2.0
type FloatList interface {
	DataList

	Append(value float64) error
	Get() ([]float64, error)
	GetValue(index int) (float64, error)
	Prepend(value float64) error
	Set(list []float64) error
	SetValue(index int, value float64) error
}

// ExternalFloatList supports binding a list of float64 values from an external variable.
//
// Since: 2.0
type ExternalFloatList interface {
	FloatList

	Reload() error
}

// NewFloatList returns a bindable list of float64 values.
//
// Since: 2.0
func NewFloatList() FloatList {
	return &boundFloatList{val: &[]float64{}}
}

// BindFloatList returns a bound list of float64 values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindFloatList(v *[]float64) ExternalFloatList {
	if v == nil {
		return NewFloatList().(ExternalFloatList)
	}

	b := &boundFloatList{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bindFloatListItem(v, i, b.updateExternal))
	}

	return b
}

type boundFloatList struct {
	listBase

	updateExternal bool
	val            *[]float64
}

func (l *boundFloatList) Append(val float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *boundFloatList) Get() ([]float64, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundFloatList) GetValue(i int) (float64, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return 0.0, errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundFloatList) Prepend(val float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([]float64{val}, *l.val...)

	return l.doReload()
}

func (l *boundFloatList) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

func (l *boundFloatList) Set(v []float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *boundFloatList) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindFloatListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternalFloatListItem).lock.Lock()
			err = item.(*boundExternalFloatListItem).setIfChanged((*l.val)[i])
			item.(*boundExternalFloatListItem).lock.Unlock()
		} else {
			item.(*boundFloatListItem).lock.Lock()
			err = item.(*boundFloatListItem).doSet((*l.val)[i])
			item.(*boundFloatListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundFloatList) SetValue(i int, v float64) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Float).Set(v)
}

func bindFloatListItem(v *[]float64, i int, external bool) Float {
	if external {
		ret := &boundExternalFloatListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &boundFloatListItem{val: v, index: i}
}

type boundFloatListItem struct {
	base

	val   *[]float64
	index int
}

func (b *boundFloatListItem) Get() (float64, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return 0.0, errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundFloatListItem) Set(val float64) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *boundFloatListItem) doSet(val float64) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternalFloatListItem struct {
	boundFloatListItem

	old float64
}

func (b *boundExternalFloatListItem) setIfChanged(val float64) error {
	if val == b.old {
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}

// IntList supports binding a list of int values.
//
// Since: 2.0
type IntList interface {
	DataList

	Append(value int) error
	Get() ([]int, error)
	GetValue(index int) (int, error)
	Prepend(value int) error
	Set(list []int) error
	SetValue(index int, value int) error
}

// ExternalIntList supports binding a list of int values from an external variable.
//
// Since: 2.0
type ExternalIntList interface {
	IntList

	Reload() error
}

// NewIntList returns a bindable list of int values.
//
// Since: 2.0
func NewIntList() IntList {
	return &boundIntList{val: &[]int{}}
}

// BindIntList returns a bound list of int values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindIntList(v *[]int) ExternalIntList {
	if v == nil {
		return NewIntList().(ExternalIntList)
	}

	b := &boundIntList{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bindIntListItem(v, i, b.updateExternal))
	}

	return b
}

type boundIntList struct {
	listBase

	updateExternal bool
	val            *[]int
}

func (l *boundIntList) Append(val int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *boundIntList) Get() ([]int, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundIntList) GetValue(i int) (int, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return 0, errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundIntList) Prepend(val int) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([]int{val}, *l.val...)

	return l.doReload()
}

func (l *boundIntList) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

func (l *boundIntList) Set(v []int) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *boundIntList) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindIntListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternalIntListItem).lock.Lock()
			err = item.(*boundExternalIntListItem).setIfChanged((*l.val)[i])
			item.(*boundExternalIntListItem).lock.Unlock()
		} else {
			item.(*boundIntListItem).lock.Lock()
			err = item.(*boundIntListItem).doSet((*l.val)[i])
			item.(*boundIntListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundIntList) SetValue(i int, v int) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Int).Set(v)
}

func bindIntListItem(v *[]int, i int, external bool) Int {
	if external {
		ret := &boundExternalIntListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &boundIntListItem{val: v, index: i}
}

type boundIntListItem struct {
	base

	val   *[]int
	index int
}

func (b *boundIntListItem) Get() (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return 0, errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundIntListItem) Set(val int) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *boundIntListItem) doSet(val int) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternalIntListItem struct {
	boundIntListItem

	old int
}

func (b *boundExternalIntListItem) setIfChanged(val int) error {
	if val == b.old {
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}

// RuneList supports binding a list of rune values.
//
// Since: 2.0
type RuneList interface {
	DataList

	Append(value rune) error
	Get() ([]rune, error)
	GetValue(index int) (rune, error)
	Prepend(value rune) error
	Set(list []rune) error
	SetValue(index int, value rune) error
}

// ExternalRuneList supports binding a list of rune values from an external variable.
//
// Since: 2.0
type ExternalRuneList interface {
	RuneList

	Reload() error
}

// NewRuneList returns a bindable list of rune values.
//
// Since: 2.0
func NewRuneList() RuneList {
	return &boundRuneList{val: &[]rune{}}
}

// BindRuneList returns a bound list of rune values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindRuneList(v *[]rune) ExternalRuneList {
	if v == nil {
		return NewRuneList().(ExternalRuneList)
	}

	b := &boundRuneList{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bindRuneListItem(v, i, b.updateExternal))
	}

	return b
}

type boundRuneList struct {
	listBase

	updateExternal bool
	val            *[]rune
}

func (l *boundRuneList) Append(val rune) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *boundRuneList) Get() ([]rune, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundRuneList) GetValue(i int) (rune, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return rune(0), errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundRuneList) Prepend(val rune) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([]rune{val}, *l.val...)

	return l.doReload()
}

func (l *boundRuneList) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

func (l *boundRuneList) Set(v []rune) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *boundRuneList) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindRuneListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternalRuneListItem).lock.Lock()
			err = item.(*boundExternalRuneListItem).setIfChanged((*l.val)[i])
			item.(*boundExternalRuneListItem).lock.Unlock()
		} else {
			item.(*boundRuneListItem).lock.Lock()
			err = item.(*boundRuneListItem).doSet((*l.val)[i])
			item.(*boundRuneListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundRuneList) SetValue(i int, v rune) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Rune).Set(v)
}

func bindRuneListItem(v *[]rune, i int, external bool) Rune {
	if external {
		ret := &boundExternalRuneListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &boundRuneListItem{val: v, index: i}
}

type boundRuneListItem struct {
	base

	val   *[]rune
	index int
}

func (b *boundRuneListItem) Get() (rune, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return rune(0), errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundRuneListItem) Set(val rune) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *boundRuneListItem) doSet(val rune) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternalRuneListItem struct {
	boundRuneListItem

	old rune
}

func (b *boundExternalRuneListItem) setIfChanged(val rune) error {
	if val == b.old {
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}

// StringList supports binding a list of string values.
//
// Since: 2.0
type StringList interface {
	DataList

	Append(value string) error
	Get() ([]string, error)
	GetValue(index int) (string, error)
	Prepend(value string) error
	Set(list []string) error
	SetValue(index int, value string) error
}

// ExternalStringList supports binding a list of string values from an external variable.
//
// Since: 2.0
type ExternalStringList interface {
	StringList

	Reload() error
}

// NewStringList returns a bindable list of string values.
//
// Since: 2.0
func NewStringList() StringList {
	return &boundStringList{val: &[]string{}}
}

// BindStringList returns a bound list of string values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindStringList(v *[]string) ExternalStringList {
	if v == nil {
		return NewStringList().(ExternalStringList)
	}

	b := &boundStringList{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bindStringListItem(v, i, b.updateExternal))
	}

	return b
}

type boundStringList struct {
	listBase

	updateExternal bool
	val            *[]string
}

func (l *boundStringList) Append(val string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *boundStringList) Get() ([]string, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundStringList) GetValue(i int) (string, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return "", errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundStringList) Prepend(val string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([]string{val}, *l.val...)

	return l.doReload()
}

func (l *boundStringList) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

func (l *boundStringList) Set(v []string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *boundStringList) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindStringListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternalStringListItem).lock.Lock()
			err = item.(*boundExternalStringListItem).setIfChanged((*l.val)[i])
			item.(*boundExternalStringListItem).lock.Unlock()
		} else {
			item.(*boundStringListItem).lock.Lock()
			err = item.(*boundStringListItem).doSet((*l.val)[i])
			item.(*boundStringListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundStringList) SetValue(i int, v string) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(String).Set(v)
}

func bindStringListItem(v *[]string, i int, external bool) String {
	if external {
		ret := &boundExternalStringListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &boundStringListItem{val: v, index: i}
}

type boundStringListItem struct {
	base

	val   *[]string
	index int
}

func (b *boundStringListItem) Get() (string, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return "", errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundStringListItem) Set(val string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *boundStringListItem) doSet(val string) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternalStringListItem struct {
	boundStringListItem

	old string
}

func (b *boundExternalStringListItem) setIfChanged(val string) error {
	if val == b.old {
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}

// UntypedList supports binding a list of interface{} values.
//
// Since: 2.1
type UntypedList interface {
	DataList

	Append(value interface{}) error
	Get() ([]interface{}, error)
	GetValue(index int) (interface{}, error)
	Prepend(value interface{}) error
	Set(list []interface{}) error
	SetValue(index int, value interface{}) error
}

// ExternalUntypedList supports binding a list of interface{} values from an external variable.
//
// Since: 2.1
type ExternalUntypedList interface {
	UntypedList

	Reload() error
}

// NewUntypedList returns a bindable list of interface{} values.
//
// Since: 2.1
func NewUntypedList() UntypedList {
	return &boundUntypedList{val: &[]interface{}{}}
}

// BindUntypedList returns a bound list of interface{} values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindUntypedList(v *[]interface{}) ExternalUntypedList {
	if v == nil {
		return NewUntypedList().(ExternalUntypedList)
	}

	b := &boundUntypedList{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bindUntypedListItem(v, i, b.updateExternal))
	}

	return b
}

type boundUntypedList struct {
	listBase

	updateExternal bool
	val            *[]interface{}
}

func (l *boundUntypedList) Append(val interface{}) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *boundUntypedList) Get() ([]interface{}, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundUntypedList) GetValue(i int) (interface{}, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return nil, errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundUntypedList) Prepend(val interface{}) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([]interface{}{val}, *l.val...)

	return l.doReload()
}

func (l *boundUntypedList) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

func (l *boundUntypedList) Set(v []interface{}) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *boundUntypedList) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindUntypedListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternalUntypedListItem).lock.Lock()
			err = item.(*boundExternalUntypedListItem).setIfChanged((*l.val)[i])
			item.(*boundExternalUntypedListItem).lock.Unlock()
		} else {
			item.(*boundUntypedListItem).lock.Lock()
			err = item.(*boundUntypedListItem).doSet((*l.val)[i])
			item.(*boundUntypedListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundUntypedList) SetValue(i int, v interface{}) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Untyped).Set(v)
}

func bindUntypedListItem(v *[]interface{}, i int, external bool) Untyped {
	if external {
		ret := &boundExternalUntypedListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &boundUntypedListItem{val: v, index: i}
}

type boundUntypedListItem struct {
	base

	val   *[]interface{}
	index int
}

func (b *boundUntypedListItem) Get() (interface{}, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return nil, errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundUntypedListItem) Set(val interface{}) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *boundUntypedListItem) doSet(val interface{}) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternalUntypedListItem struct {
	boundUntypedListItem

	old interface{}
}

func (b *boundExternalUntypedListItem) setIfChanged(val interface{}) error {
	if val == b.old {
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}

// URIList supports binding a list of fyne.URI values.
//
// Since: 2.1
type URIList interface {
	DataList

	Append(value fyne.URI) error
	Get() ([]fyne.URI, error)
	GetValue(index int) (fyne.URI, error)
	Prepend(value fyne.URI) error
	Set(list []fyne.URI) error
	SetValue(index int, value fyne.URI) error
}

// ExternalURIList supports binding a list of fyne.URI values from an external variable.
//
// Since: 2.1
type ExternalURIList interface {
	URIList

	Reload() error
}

// NewURIList returns a bindable list of fyne.URI values.
//
// Since: 2.1
func NewURIList() URIList {
	return &boundURIList{val: &[]fyne.URI{}}
}

// BindURIList returns a bound list of fyne.URI values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindURIList(v *[]fyne.URI) ExternalURIList {
	if v == nil {
		return NewURIList().(ExternalURIList)
	}

	b := &boundURIList{val: v, updateExternal: true}

	for i := range *v {
		b.appendItem(bindURIListItem(v, i, b.updateExternal))
	}

	return b
}

type boundURIList struct {
	listBase

	updateExternal bool
	val            *[]fyne.URI
}

func (l *boundURIList) Append(val fyne.URI) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	*l.val = append(*l.val, val)

	return l.doReload()
}

func (l *boundURIList) Get() ([]fyne.URI, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundURIList) GetValue(i int) (fyne.URI, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return fyne.URI(nil), errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundURIList) Prepend(val fyne.URI) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = append([]fyne.URI{val}, *l.val...)

	return l.doReload()
}

func (l *boundURIList) Reload() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.doReload()
}

func (l *boundURIList) Set(v []fyne.URI) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	*l.val = v

	return l.doReload()
}

func (l *boundURIList) doReload() (retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindURIListItem(l.val, i, l.updateExternal))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			item.(*boundExternalURIListItem).lock.Lock()
			err = item.(*boundExternalURIListItem).setIfChanged((*l.val)[i])
			item.(*boundExternalURIListItem).lock.Unlock()
		} else {
			item.(*boundURIListItem).lock.Lock()
			err = item.(*boundURIListItem).doSet((*l.val)[i])
			item.(*boundURIListItem).lock.Unlock()
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundURIList) SetValue(i int, v fyne.URI) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(URI).Set(v)
}

func bindURIListItem(v *[]fyne.URI, i int, external bool) URI {
	if external {
		ret := &boundExternalURIListItem{old: (*v)[i]}
		ret.val = v
		ret.index = i
		return ret
	}

	return &boundURIListItem{val: v, index: i}
}

type boundURIListItem struct {
	base

	val   *[]fyne.URI
	index int
}

func (b *boundURIListItem) Get() (fyne.URI, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return fyne.URI(nil), errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundURIListItem) Set(val fyne.URI) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doSet(val)
}

func (b *boundURIListItem) doSet(val fyne.URI) error {
	(*b.val)[b.index] = val

	b.trigger()
	return nil
}

type boundExternalURIListItem struct {
	boundURIListItem

	old fyne.URI
}

func (b *boundExternalURIListItem) setIfChanged(val fyne.URI) error {
	if compareURI(val, b.old) {
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.trigger()
	return nil
}
