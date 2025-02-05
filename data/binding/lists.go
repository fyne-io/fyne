package binding

import (
	"bytes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// DataList is the base interface for all bindable data lists.
//
// Since: 2.0
type DataList interface {
	DataItem
	GetItem(index int) (DataItem, error)
	Length() int
}

// BoolList supports binding a list of bool values.
//
// Since: 2.0
type BoolList interface {
	DataList

	Append(value bool) error
	Get() ([]bool, error)
	GetValue(index int) (bool, error)
	Prepend(value bool) error
	Remove(value bool) error
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
	return newListComparable[bool]()
}

// BindBoolList returns a bound list of bool values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindBoolList(v *[]bool) ExternalBoolList {
	return bindListComparable(v)
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
	Remove(value []byte) error
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
	return newList(bytes.Equal)
}

// BindBytesList returns a bound list of []byte values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.2
func BindBytesList(v *[][]byte) ExternalBytesList {
	return bindList(v, bytes.Equal)
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
	Remove(value float64) error
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
	return newListComparable[float64]()
}

// BindFloatList returns a bound list of float64 values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindFloatList(v *[]float64) ExternalFloatList {
	return bindListComparable(v)
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
	Remove(value int) error
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
	return newListComparable[int]()
}

// BindIntList returns a bound list of int values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindIntList(v *[]int) ExternalIntList {
	return bindListComparable(v)
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
	Remove(value rune) error
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
	return newListComparable[rune]()
}

// BindRuneList returns a bound list of rune values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindRuneList(v *[]rune) ExternalRuneList {
	if v == nil {
		return NewRuneList().(ExternalRuneList)
	}

	b := newListComparable[rune]()
	b.val = v
	b.updateExternal = true

	for i := range *v {
		b.appendItem(bindListItemComparable(v, i, b.updateExternal))
	}

	return b
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
	Remove(value string) error
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
	return newListComparable[string]()
}

// BindStringList returns a bound list of string values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindStringList(v *[]string) ExternalStringList {
	return bindListComparable(v)
}

// UntypedList supports binding a list of any values.
//
// Since: 2.1
type UntypedList interface {
	DataList

	Append(value any) error
	Get() ([]any, error)
	GetValue(index int) (any, error)
	Prepend(value any) error
	Remove(value any) error
	Set(list []any) error
	SetValue(index int, value any) error
}

// ExternalUntypedList supports binding a list of any values from an external variable.
//
// Since: 2.1
type ExternalUntypedList interface {
	UntypedList

	Reload() error
}

// NewUntypedList returns a bindable list of any values.
//
// Since: 2.1
func NewUntypedList() UntypedList {
	return newList(func(t1, t2 any) bool { return t1 == t2 })
}

// BindUntypedList returns a bound list of any values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindUntypedList(v *[]any) ExternalUntypedList {
	if v == nil {
		return NewUntypedList().(ExternalUntypedList)
	}

	comparator := func(t1, t2 any) bool { return t1 == t2 }
	b := newExternalList(v, comparator)
	for i := range *v {
		b.appendItem(bindListItem(v, i, b.updateExternal, comparator))
	}

	return b
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
	Remove(value fyne.URI) error
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
	return newList(storage.EqualURI)
}

// BindURIList returns a bound list of fyne.URI values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindURIList(v *[]fyne.URI) ExternalURIList {
	return bindList(v, storage.EqualURI)
}

type listBase struct {
	base
	items []DataItem
}

// GetItem returns the DataItem at the specified index.
func (b *listBase) GetItem(i int) (DataItem, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if i < 0 || i >= len(b.items) {
		return nil, errOutOfBounds
	}

	return b.items[i], nil
}

// Length returns the number of items in this data list.
func (b *listBase) Length() int {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return len(b.items)
}

func (b *listBase) appendItem(i DataItem) {
	b.items = append(b.items, i)
}

func (b *listBase) deleteItem(i int) {
	b.items = append(b.items[:i], b.items[i+1:]...)
}
