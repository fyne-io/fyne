package widget

import (
	"fmt"
	"reflect"
	"strconv"

	"fyne.io/fyne/dataapi"
)

// DataSetter interface
type DataSetter interface {
	SetFromData(dataapi.DataItem)
}

// DataMapSetter interface
type DataMapSetter interface {
	SetFromData(dataapi.DataMap)
}

// DataSourceSetter interface, for things that can be set from a source
type DataSourceSetter interface {
	SetFromSource(source dataapi.DataSource)
}

// DataListener is a base struct that all DataAPI aware widgets can embed
type DataListener struct {
	data       dataapi.DataItem
	source     dataapi.DataSource
	onChanged  reflect.Value
	listenerID int
}

// StringToInter objects can get an int index from a string (Radio / Select, etc)
type StringToInter interface {
	AsInt(string) int
}

// This gets a bit sticky, because the OnChange handlers on each widget are not methods, but attribs
// so some reflection magic is needed to see if the object has the appropriate OnChange handler,
// and set it accordingly
func getFunction(obj interface{}, name string) (reflect.Value, bool) {
	f := reflect.ValueOf(obj).Elem().FieldByName(name)
	return f, f.IsValid() && f.CanSet()
}

// Bind will Bind this widget to the given DataItem
// It takes a DataItem, and an object that can be set from a DataItem
// It sets the UpdateBoundData() callback function on the setter object passed in
// to a new function that correctly handles the type of the DataItem
func (d *DataListener) Bind(data dataapi.DataItem, setter DataSetter) {
	d.data = data
	d.listenerID = data.AddListener(setter.SetFromData)
	setter.SetFromData(data)

	if f, ok := getFunction(setter, "UpdateBinding"); ok {
		ftype := f.Type().String()
		switch ftype {
		case "func(string)":
			if s, ok := data.(dataapi.Settable); ok {
				f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
					s.SetStringWithMask(in[0].String(), d.listenerID)
					return nil
				}))
			} else if s, ok := data.(dataapi.SettableInt); ok {
				if ss, ok := setter.(StringToInter); ok {
					f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
						s.SetIntWithMask(ss.AsInt(in[0].String()), d.listenerID)
						return nil
					}))
				} else {
					f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
						vv, _ := strconv.Atoi(in[0].String())
						s.SetIntWithMask(vv, d.listenerID)
						return nil
					}))
				}
			} else if s, ok := data.(dataapi.SettableFloat); ok {
				f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
					vv, _ := strconv.ParseFloat(in[0].String(), 64)
					s.SetFloatWithMask(vv, d.listenerID)
					return nil
				}))
			}
		case "func(bool)":
			if s, ok := data.(dataapi.SettableBool); ok {
				f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
					s.SetBoolWithMask(in[0].Bool(), d.listenerID)
					return []reflect.Value{}
				}))
			} else if s, ok := data.(dataapi.Settable); ok {
				f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
					bb := in[0].Bool()
					ss := "false"
					if bb {
						ss = "true"
					}
					s.SetStringWithMask(ss, d.listenerID)
					return nil
				}))
			}
		case "func(float64)":
			if s, ok := data.(dataapi.SettableFloat); ok {
				f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
					s.SetFloatWithMask(in[0].Float(), d.listenerID)
					return nil
				}))
			} else if s, ok := data.(dataapi.SettableInt); ok {
				f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
					s.SetIntWithMask(int(in[0].Float()), d.listenerID)
					return nil
				}))
			} else if s, ok := data.(dataapi.Settable); ok {
				f.Set(reflect.MakeFunc(f.Type(), func(in []reflect.Value) []reflect.Value {
					s.SetStringWithMask(fmt.Sprintf("%v", in[0].Float()), d.listenerID)
					return nil
				}))
			}
		}
	}
}

// Source will connect this listener to a DataSource, which is a 1 way binding
func (d *DataListener) Source(src dataapi.DataSource, setter DataSourceSetter) {
	d.source = src
	d.listenerID = src.AddListener(setter.SetFromSource)
	setter.SetFromSource(src)
}
