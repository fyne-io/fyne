package binding

import (
	"errors"
	"reflect"

	"fyne.io/fyne/v2"
)

// DataMap is the base interface for all bindable data maps.
//
// Since: 2.0
type DataMap interface {
	DataItem
	GetItem(string) (DataItem, error)
	Keys() []string
}

// ExternalUntypedMap is a map data binding with all values untyped (interface{}), connected to an external data source.
//
// Since: 2.0
type ExternalUntypedMap interface {
	UntypedMap
	Reload() error
}

// UntypedMap is a map data binding with all values Untyped (interface{}).
//
// Since: 2.0
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
// Since: 2.0
func NewUntypedMap() UntypedMap {
	return &mapBase{items: make(map[string]reflectUntyped), val: &map[string]interface{}{}}
}

// BindUntypedMap creates a new map binding of string to interface{} based on the data passed.
// If your code changes the content of the map this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindUntypedMap(d *map[string]interface{}) ExternalUntypedMap {
	if d == nil {
		return NewUntypedMap().(ExternalUntypedMap)
	}
	m := &mapBase{items: make(map[string]reflectUntyped), val: d, updateExternal: true}

	for k := range *d {
		m.setItem(k, bindUntypedMapValue(d, k, m.updateExternal))
	}

	return m
}

// Struct is the base interface for a bound struct type.
//
// Since: 2.0
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
// Since: 2.0
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

	s := &boundStruct{orig: i}
	s.items = make(map[string]reflectUntyped)
	s.val = &map[string]interface{}{}
	s.updateExternal = true

	v := reflect.ValueOf(i).Elem()
	t = v.Type()
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		if !f.CanSet() {
			continue
		}

		key := t.Field(j).Name
		s.items[key] = bindReflect(f)
		(*s.val)[key] = f.Interface()
	}

	return s
}

type reflectUntyped interface {
	DataItem
	get() (interface{}, error)
	set(interface{}) error
}

type mapBase struct {
	base

	updateExternal bool
	items          map[string]reflectUntyped
	val            *map[string]interface{}
}

func (b *mapBase) GetItem(key string) (DataItem, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if v, ok := b.items[key]; ok {
		return v, nil
	}

	return nil, errKeyNotFound
}

func (b *mapBase) Keys() []string {
	b.lock.Lock()
	defer b.lock.Unlock()

	ret := make([]string, len(b.items))
	i := 0
	for k := range b.items {
		ret[i] = k
		i++
	}

	return ret
}

func (b *mapBase) Delete(key string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	delete(b.items, key)

	b.trigger()
}

func (b *mapBase) Get() (map[string]interface{}, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return map[string]interface{}{}, nil
	}

	return *b.val, nil
}

func (b *mapBase) GetValue(key string) (interface{}, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if i, ok := b.items[key]; ok {
		return i.get()
	}

	return nil, errKeyNotFound
}

func (b *mapBase) Reload() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doReload()
}

func (b *mapBase) Set(v map[string]interface{}) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &v
		b.trigger()
		return nil
	}

	*b.val = v
	return b.doReload()
}

func (b *mapBase) SetValue(key string, d interface{}) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if i, ok := b.items[key]; ok {
		return i.set(d)
	}

	(*b.val)[key] = d
	item := bindUntypedMapValue(b.val, key, b.updateExternal)
	b.setItem(key, item)
	return nil
}

func (b *mapBase) doReload() (retErr error) {
	changed := false
	// add new
	for key := range *b.val {
		_, found := b.items[key]
		if !found {
			b.setItem(key, bindUntypedMapValue(b.val, key, b.updateExternal))
			changed = true
		}
	}

	// remove old
	for key := range b.items {
		_, found := (*b.val)[key]
		if !found {
			delete(b.items, key)
			changed = true
		}
	}
	if changed {
		b.trigger()
	}

	for k, item := range b.items {
		var err error

		if b.updateExternal {
			err = item.(*boundExternalMapValue).setIfChanged((*b.val)[k])
		} else {
			err = item.(*boundMapValue).set((*b.val)[k])
		}

		if err != nil {
			retErr = err
		}
	}
	return
}

func (b *mapBase) setItem(key string, d reflectUntyped) {
	b.items[key] = d

	b.trigger()
}

type boundStruct struct {
	mapBase

	orig interface{}
}

func (b *boundStruct) Reload() (retErr error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	v := reflect.ValueOf(b.orig).Elem()
	t := v.Type()
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		if !f.CanSet() {
			continue
		}
		kind := f.Kind()
		if kind == reflect.Slice || kind == reflect.Struct {
			fyne.LogError("Data binding does not yet support slice or struct elements in a struct", nil)
			continue
		}

		key := t.Field(j).Name
		old := (*b.val)[key]
		if f.Interface() == old {
			continue
		}

		var err error
		switch kind {
		case reflect.Bool:
			err = b.items[key].(*reflectBool).Set(f.Bool())
		case reflect.Float32, reflect.Float64:
			err = b.items[key].(*reflectFloat).Set(f.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = b.items[key].(*reflectInt).Set(int(f.Int()))
		case reflect.String:
			err = b.items[key].(*reflectString).Set(f.String())
		}
		if err != nil {
			retErr = err
		}
		(*b.val)[key] = f.Interface()
	}
	return
}

func bindUntypedMapValue(m *map[string]interface{}, k string, external bool) reflectUntyped {
	if external {
		ret := &boundExternalMapValue{old: (*m)[k]}
		ret.val = m
		ret.key = k
		return ret
	}

	return &boundMapValue{val: m, key: k}
}

type boundMapValue struct {
	base

	val *map[string]interface{}
	key string
}

func (b *boundMapValue) get() (interface{}, error) {
	if v, ok := (*b.val)[b.key]; ok {
		return v, nil
	}

	return nil, errKeyNotFound
}

func (b *boundMapValue) set(val interface{}) error {
	(*b.val)[b.key] = val

	b.trigger()
	return nil
}

type boundExternalMapValue struct {
	boundMapValue

	old interface{}
}

func (b *boundExternalMapValue) setIfChanged(val interface{}) error {
	if val == b.old {
		return nil
	}
	b.old = val

	return b.set(val)
}

type boundReflect struct {
	base

	val reflect.Value
}

func (b *boundReflect) get() (interface{}, error) {
	return b.val.Interface(), nil
}

func (b *boundReflect) set(val interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to set bool in data binding")
		}
	}()
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
	r.trigger()
	return
}

func bindReflectBool(f reflect.Value) reflectUntyped {
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
	r.trigger()
	return
}

func bindReflectFloat(f reflect.Value) reflectUntyped {
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
	r.trigger()
	return
}

func bindReflectInt(f reflect.Value) reflectUntyped {
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
	r.trigger()
	return
}

func bindReflectString(f reflect.Value) reflectUntyped {
	r := &reflectString{}
	r.val = f
	return r
}

func bindReflect(field reflect.Value) reflectUntyped {
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
