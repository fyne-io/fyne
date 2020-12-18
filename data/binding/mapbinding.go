package binding

import (
	"reflect"

	"fyne.io/fyne"
)

// DataMap is the base interface for all bindable data lists.
//
// Since: 2.0.0
type DataMap interface {
	DataItem
	GetItem(string) DataItem
	Keys() []string
}

// UntypedMap is a map data binding with all values untyped (interface{}).
//
// Since: 2.0.0
type UntypedMap interface {
	DataMap
	Delete(string)
	Get(string) interface{}
	Set(string, interface{})
}

// NewUntypedMap creates a new, empty map binding of string to interface{}.
//
// Since: 2.0.0
func NewUntypedMap() UntypedMap {
	return &mapBase{val: make(map[string]DataItem)}
}

// BindUntypedMap creates a new map binding of string to interface{} based on the data passed.
//
// Since: 2.0.0
func BindUntypedMap(d *map[string]interface{}) UntypedMap {
	if d == nil {
		return NewUntypedMap()
	}
	m := &mapBase{val: make(map[string]DataItem)}

	for k, v := range *d {
		m.Set(k, v)
	}

	return m
}

// BindStruct creates a new map biding of string to interface{} using the struct passed as data.
// The key in for each item is a string representation of each exported field with the value set as an interface{}.
// Only exported fields are included.
//
// Since: 2.0.0
func BindStruct(i interface{}) DataMap {
	if i == nil {
		return NewUntypedMap()
	}
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr ||
		(reflect.TypeOf(reflect.ValueOf(i).Elem()).Kind() != reflect.Struct) {
		fyne.LogError("Invalid type passed to BindStruct, must be pointer to struct", nil)
		return NewUntypedMap()
	}

	m := &mapBase{val: make(map[string]DataItem)}
	v := reflect.ValueOf(i).Elem()
	t = v.Type()
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		if !f.CanSet() {
			continue
		}

		m.setItem(t.Field(j).Name, bindReflect(f))
	}

	return m
}

type mapBase struct {
	base
	val map[string]DataItem
}

// GetItem returns the DataItem at the specified key.
// It will return nil if the key was not found.
//
// Since: 2.0.0
func (b *mapBase) GetItem(key string) DataItem {
	if v, ok := b.val[key]; ok {
		return v
	}

	return nil
}

// Keys returns a list of all the keys in this data map.
//
// Since: 2.0.0
func (b *mapBase) Keys() []string {
	ret := make([]string, len(b.val))
	// TODO lock
	i := 0
	for k := range b.val {
		ret[i] = k
		i++
	}

	return ret
}

func (b *mapBase) Delete(key string) {
	delete(b.val, key)

	b.trigger()
}

func (b *mapBase) Get(key string) interface{} {
	if i, ok := b.val[key]; ok {
		return i.(untyped).get()
	}

	return nil
}

func (b *mapBase) Set(key string, d interface{}) {
	if i, ok := b.val[key]; ok {
		i.(untyped).set(d)
		return
	}

	b.setItem(key, bindUntyped(&d))
}

func (b *mapBase) setItem(key string, d DataItem) {
	b.val[key] = d

	b.trigger()
}

type untyped interface {
	DataItem
	get() interface{}
	set(interface{})
}

func bindUntyped(v *interface{}) untyped {
	if v == nil {
		return &boundUntyped{val: nil}
	}

	return &boundUntyped{val: v}
}

type boundUntyped struct {
	base

	val *interface{}
}

func (b *boundUntyped) get() interface{} {
	if b.val == nil {
		return 0
	}
	return *b.val
}

func (b *boundUntyped) set(val interface{}) {
	if *b.val == val {
		return
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
}

type boundReflect struct {
	base

	val reflect.Value
}

func (b *boundReflect) get() interface{} {
	return b.val.Interface()
}

func (b *boundReflect) set(val interface{}) {
	b.val.Set(reflect.ValueOf(val))

	b.trigger()
}

type reflectBool struct {
	boundReflect
}

func (r *reflectBool) Get() bool {
	return r.val.Bool()
}

func (r *reflectBool) Set(b bool) {
	r.val.SetBool(b)
}

func bindReflectBool(f reflect.Value) DataItem {
	r := &reflectBool{}
	r.val = f
	return r
}

type reflectFloat struct {
	boundReflect
}

func (r *reflectFloat) Get() float64 {
	return r.val.Float()
}

func (r *reflectFloat) Set(f float64) {
	r.val.SetFloat(f)
}

func bindReflectFloat(f reflect.Value) DataItem {
	r := &reflectFloat{}
	r.val = f
	return r
}

type reflectInt struct {
	boundReflect
}

func (r *reflectInt) Get() int {
	return int(r.val.Int())
}

func (r *reflectInt) Set(i int) {
	r.val.SetInt(int64(i))
}

func bindReflectInt(f reflect.Value) DataItem {
	r := &reflectInt{}
	r.val = f
	return r
}

type reflectString struct {
	boundReflect
}

func (r *reflectString) Get() string {
	return r.val.String()
}

func (r *reflectString) Set(s string) {
	r.val.SetString(s)
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
