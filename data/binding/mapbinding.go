package binding

import (
	"errors"
	"reflect"

	"fyne.io/fyne"
)

// DataMap is the base interface for all bindable data maps.
//
// Since: 2.0.0
type DataMap interface {
	DataItem
	GetItem(string) (DataItem, error)
	Keys() []string
}

// UntypedMap is a map data binding with all values untyped (interface{}).
//
// Since: 2.0.0
type UntypedMap interface {
	DataMap
	Delete(string)
	Get(string) (interface{}, error)
	Set(string, interface{}) error
}

// NewUntypedMap creates a new, empty map binding of string to interface{}.
//
// Since: 2.0.0
func NewUntypedMap() UntypedMap {
	return &mapBase{items: make(map[string]DataItem)}
}

// BindUntypedMap creates a new map binding of string to interface{} based on the data passed.
//
// Since: 2.0.0
func BindUntypedMap(d *map[string]interface{}) UntypedMap {
	if d == nil {
		return NewUntypedMap()
	}
	m := &mapBase{items: make(map[string]DataItem), val: d}

	for k := range *d {
		m.setItem(k, bindUntyped((*d)[k]))
	}

	return m
}

// Struct is the base interface for a bound struct type.
//
// Since: 2.0.0
type Struct interface {
	DataMap
	Get(string) (interface{}, error)
	Set(string, interface{}) error
}

// BindStruct creates a new map binding of string to interface{} using the struct passed as data.
// The key for each item is a string representation of each exported field with the value set as an interface{}.
// Only exported fields are included.
//
// Since: 2.0.0
func BindStruct(i interface{}) Struct {
	if i == nil {
		return NewUntypedMap()
	}
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr ||
		(reflect.TypeOf(reflect.ValueOf(i).Elem()).Kind() != reflect.Struct) {
		fyne.LogError("Invalid type passed to BindStruct, must be pointer to struct", nil)
		return NewUntypedMap()
	}

	m := &mapBase{items: make(map[string]DataItem)}
	v := reflect.ValueOf(i).Elem()
	t = v.Type()
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		if !f.CanSet() {
			continue
		}

		m.items[t.Field(j).Name] = bindReflect(f)
	}

	return m
}

type mapBase struct {
	base
	items map[string]DataItem
	val   *map[string]interface{}
}

// GetItem returns the DataItem at the specified key.
// It will return nil if the key was not found.
//
// Since: 2.0.0
func (b *mapBase) GetItem(key string) (DataItem, error) {
	if v, ok := b.items[key]; ok {
		return v, nil
	}

	return nil, errKeyNotFound
}

// Keys returns a list of all the keys in this data map.
//
// Since: 2.0.0
func (b *mapBase) Keys() []string {
	ret := make([]string, len(b.items))
	// TODO lock
	i := 0
	for k := range b.items {
		ret[i] = k
		i++
	}

	return ret
}

// Delete removes the specified key and tha value associated with it.
//
// Since: 2.0.0
func (b *mapBase) Delete(key string) {
	delete(b.items, key)

	b.trigger()
}

// Get returns the value stored at the specified key.
//
// Since: 2.0.0
func (b *mapBase) Get(key string) (interface{}, error) {
	if i, ok := b.items[key]; ok {
		return i.(untyped).get()
	}

	return nil, errKeyNotFound
}

// Set stores the value d at the specified key.
// If the key is not present it will create a new binding internally.
//
// Since: 2.0.0
func (b *mapBase) Set(key string, d interface{}) error {
	if i, ok := b.items[key]; ok {
		i.(untyped).set(d)
		return nil
	}

	item := bindUntyped(d)
	b.setItem(key, item)
	return nil
}

func (b *mapBase) setItem(key string, d DataItem) {
	b.items[key] = d

	b.trigger()
}

type untyped interface {
	DataItem
	get() (interface{}, error)
	set(interface{}) error
}

func bindUntyped(m interface{}) untyped {
	return &boundUntyped{val: m}
}

type boundUntyped struct {
	base

	val interface{}
}

func (b *boundUntyped) get() (interface{}, error) {
	return b.val, nil
}

func (b *boundUntyped) set(val interface{}) error {
	b.val = val

	b.trigger()
	return nil
}

type boundReflect struct {
	base

	val reflect.Value
}

func (b *boundReflect) get() (interface{}, error) {
	return b.val.Interface(), nil
}

func (b *boundReflect) set(val interface{}) error {
	// TODO catch the panic and return as error
	b.val.Set(reflect.ValueOf(val))

	b.trigger()
	return nil
}

type reflectBool struct {
	boundReflect
}

func (r *reflectBool) Get() (val bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid bool value in data binding")
		}
	}()

	val = r.val.Bool()
	return
}

func (r *reflectBool) Set(b bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to set bool in data binding")
		}
	}()

	r.val.SetBool(b)
	return
}

func bindReflectBool(f reflect.Value) DataItem {
	r := &reflectBool{}
	r.val = f
	return r
}

type reflectFloat struct {
	boundReflect
}

func (r *reflectFloat) Get() (val float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid float64 value in data binding")
		}
	}()

	val = r.val.Float()
	return
}

func (r *reflectFloat) Set(f float64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to set float64 in data binding")
		}
	}()

	r.val.SetFloat(f)
	return
}

func bindReflectFloat(f reflect.Value) DataItem {
	r := &reflectFloat{}
	r.val = f
	return r
}

type reflectInt struct {
	boundReflect
}

func (r *reflectInt) Get() (val int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid int value in data binding")
		}
	}()

	val = int(r.val.Int())
	return
}

func (r *reflectInt) Set(i int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to set int in data binding")
		}
	}()

	r.val.SetInt(int64(i))
	return
}

func bindReflectInt(f reflect.Value) DataItem {
	r := &reflectInt{}
	r.val = f
	return r
}

type reflectString struct {
	boundReflect
}

func (r *reflectString) Get() (val string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("invalid string value in data binding")
		}
	}()

	val = r.val.String()
	return
}

func (r *reflectString) Set(s string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to set string in data binding")
		}
	}()

	r.val.SetString(s)
	return
}

func bindReflectString(f reflect.Value) DataItem {
	r := &reflectString{}
	r.val = f
	return r
}

func bindReflect(field reflect.Value) DataItem {
	switch field.Kind() {
	case reflect.Bool:
		return bindReflectBool(field)
	case reflect.Float32, reflect.Float64:
		return bindReflectFloat(field)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return bindReflectInt(field)
	case reflect.String:
		return bindReflectString(field)
	}
	return &boundReflect{val: field}
}
