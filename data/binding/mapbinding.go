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

// ExternalUntypedMap is a map data binding with all values untyped (interface{}), connected to an external data source.
//
// Since: 2.0.0
type ExternalUntypedMap interface {
	UntypedMap
	Reload() error
}

// UntypedMap is a map data binding with all values Untyped (interface{}).
//
// Since: 2.0.0
type UntypedMap interface {
	DataMap
	Delete(string)
	Get() (map[string]interface{}, error)
	GetValue(string) (interface{}, error)
	Set(map[string]interface{}) error
	SetValue(string, interface{}) error
}

// NewUntypedMap creates a new, empty map binding of string to interface{}.
//
// Since: 2.0.0
func NewUntypedMap() UntypedMap {
	return &mapBase{items: make(map[string]DataItem)}
}

// BindUntypedMap creates a new map binding of string to interface{} based on the data passed.
// If your code changes the content of the map this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0.0
func BindUntypedMap(d *map[string]interface{}) ExternalUntypedMap {
	if d == nil {
		return NewUntypedMap().(ExternalUntypedMap)
	}
	m := &mapBase{items: make(map[string]DataItem), val: d}

	for k := range *d {
		m.setItem(k, bindUntypedMapValue(d, k))
	}

	return m
}

// Struct is the base interface for a bound struct type.
//
// Since: 2.0.0
type Struct interface {
	DataMap
	GetValue(string) (interface{}, error)
	SetValue(string, interface{}) error
	Reload() error
}

// BindStruct creates a new map binding of string to interface{} using the struct passed as data.
// The key for each item is a string representation of each exported field with the value set as an interface{}.
// Only exported fields are included.
//
// Since: 2.0.0
func BindStruct(i interface{}) Struct {
	if i == nil {
		return NewUntypedMap().(Struct)
	}
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr ||
		(reflect.TypeOf(reflect.ValueOf(i).Elem()).Kind() != reflect.Struct) {
		fyne.LogError("Invalid type passed to BindStruct, must be pointer to struct", nil)
		return NewUntypedMap().(Struct)
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

// Untyped id used tpo represent binding an interface{} value.
//
// Since: 2.0.0
type Untyped interface {
	DataItem
	get() (interface{}, error)
	set(interface{}) error
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

func (b *mapBase) Get() (map[string]interface{}, error) {
	if b.val == nil {
		return map[string]interface{}{}, nil
	}

	return *b.val, nil
}

// Get returns the value stored at the specified key.
//
// Since: 2.0.0
func (b *mapBase) GetValue(key string) (interface{}, error) {
	if i, ok := b.items[key]; ok {
		return i.(Untyped).get()
	}

	return nil, errKeyNotFound
}

func (b *mapBase) Reload() error {
	return b.doReload()
}

func (b *mapBase) Set(v map[string]interface{}) error {
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &v
		b.trigger()
		return nil
	}

	*b.val = v
	return b.doReload()
}

func (b *mapBase) doReload() (retErr error) {
	changed := false
	// add new
	for key := range *b.val {
		found := false
		for newKey := range b.items {
			if newKey == key {
				found = true
			}
		}

		if !found {
			b.setItem(key, bindUntypedMapValue(b.val, key))
			changed = true
		}
	}

	// remove old
	for key := range b.items {
		found := false
		for newKey := range *b.val {
			if newKey == key {
				found = true
			}
		}
		if !found {
			b.Delete(key)
			changed = true
		}
	}
	if changed {
		b.trigger()
	}

	for _, item := range b.items {
		// TODO cache values and do comparison - for now we just always trigger child elements
		//		old, err := b.items[k].(Untyped).get()
		//		val := (*(b.val))[k]
		//		if err != nil || (*(b.val))[k] != old {
		//			err = item.(Untyped).set(val)
		//			if err != nil {
		//				retErr = err
		//			}
		//		}
		item.(*boundUntyped).trigger()
	}
	return
}

// Set stores the value d at the specified key.
// If the key is not present it will create a new binding internally.
//
// Since: 2.0.0
func (b *mapBase) SetValue(key string, d interface{}) error {
	if i, ok := b.items[key]; ok {
		i.(Untyped).set(d)
		return nil
	}

	(*b.val)[key] = d
	item := bindUntypedMapValue(b.val, key)
	b.setItem(key, item)
	return nil
}

func (b *mapBase) setItem(key string, d DataItem) {
	b.items[key] = d

	b.trigger()
}

func bindUntypedMapValue(m *map[string]interface{}, k string) Untyped {
	return &boundUntyped{val: m, key: k}
}

type boundUntyped struct {
	base

	val *map[string]interface{}
	key string
}

func (b *boundUntyped) get() (interface{}, error) {
	if v, ok := (*b.val)[b.key]; ok {
		return v, nil
	}

	return nil, errKeyNotFound
}

func (b *boundUntyped) set(val interface{}) error {
	(*b.val)[b.key] = val

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
