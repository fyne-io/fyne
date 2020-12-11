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
func BindUntypedMap(d map[string]interface{}) UntypedMap {
	m := &mapBase{val: make(map[string]DataItem)}

	for k, v := range d {
		m.Set(k, v)
	}

	return m
}

// BindStruct creates a new map biding of string to interface{} using the struct passed as data.
// The key in for each item is a string representation of each exported field with the value set as an interface{}.
// Only exported fields are included
//
// Since: 2.0.0
func BindStruct(i interface{}) UntypedMap {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Struct {
		fyne.LogError("Invalid type passed to BindStruct", nil)
	}

	data := make(map[string]interface{})
	v := reflect.ValueOf(i)
	for j := 0; j < v.NumField(); j++ {
		name := t.Field(j).Name
		f := v.Field(j)
		if !f.CanInterface() {
			continue
		}

		switch v.Field(j).Kind() {
		case reflect.Bool:
			data[name] = f.Bool()
		case reflect.Float32, reflect.Float64:
			data[name] = f.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			data[name] = f.Int()
		case reflect.String:
			data[name] = f.String()
		default:
			data[name] = f.Interface()
		}
	}

	return BindUntypedMap(data)
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

type untypedMap struct {
	mapBase
}

func (b *mapBase) Delete(key string) {
	b.val[key] = nil

	b.trigger()
}

func (b *mapBase) Get(key string) interface{} {
	if i, ok := b.val[key]; ok {
		return i.(untyped).Get()
	}

	return nil
}

func (b *mapBase) Set(key string, d interface{}) {
	if i, ok := b.val[key]; ok {
		i.(untyped).Set(d)
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
	Get() interface{}
	Set(interface{})
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

func (b *boundUntyped) Get() interface{} {
	if b.val == nil {
		return 0
	}
	return *b.val
}

func (b *boundUntyped) Set(val interface{}) {
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
