// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"net/url"

	"fyne.io/fyne"
)

// Bool implements a data binding for a bool.
type Bool struct {
	Base
	value bool
}

// NewBool creates a new binding with the given value.
func NewBool(value bool) *Bool {
	return &Bool{value: value}
}

// Get returns the bound value.
func (b *Bool) Get() bool {
	return b.value
}

// Set updates the bound value.
func (b *Bool) Set(value bool) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddBoolListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *Bool) AddBoolListener(listener func(bool)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// Float64 implements a data binding for a float64.
type Float64 struct {
	Base
	value float64
}

// NewFloat64 creates a new binding with the given value.
func NewFloat64(value float64) *Float64 {
	return &Float64{value: value}
}

// Get returns the bound value.
func (b *Float64) Get() float64 {
	return b.value
}

// Set updates the bound value.
func (b *Float64) Set(value float64) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddFloat64Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *Float64) AddFloat64Listener(listener func(float64)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// Int implements a data binding for a int.
type Int struct {
	Base
	value int
}

// NewInt creates a new binding with the given value.
func NewInt(value int) *Int {
	return &Int{value: value}
}

// Get returns the bound value.
func (b *Int) Get() int {
	return b.value
}

// Set updates the bound value.
func (b *Int) Set(value int) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddIntListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *Int) AddIntListener(listener func(int)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// Int64 implements a data binding for a int64.
type Int64 struct {
	Base
	value int64
}

// NewInt64 creates a new binding with the given value.
func NewInt64(value int64) *Int64 {
	return &Int64{value: value}
}

// Get returns the bound value.
func (b *Int64) Get() int64 {
	return b.value
}

// Set updates the bound value.
func (b *Int64) Set(value int64) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddInt64Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *Int64) AddInt64Listener(listener func(int64)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// Resource implements a data binding for a fyne.Resource.
type Resource struct {
	Base
	value fyne.Resource
}

// NewResource creates a new binding with the given value.
func NewResource(value fyne.Resource) *Resource {
	return &Resource{value: value}
}

// Get returns the bound value.
func (b *Resource) Get() fyne.Resource {
	return b.value
}

// Set updates the bound value.
func (b *Resource) Set(value fyne.Resource) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddResourceListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *Resource) AddResourceListener(listener func(fyne.Resource)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// Rune implements a data binding for a rune.
type Rune struct {
	Base
	value rune
}

// NewRune creates a new binding with the given value.
func NewRune(value rune) *Rune {
	return &Rune{value: value}
}

// Get returns the bound value.
func (b *Rune) Get() rune {
	return b.value
}

// Set updates the bound value.
func (b *Rune) Set(value rune) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddRuneListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *Rune) AddRuneListener(listener func(rune)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// String implements a data binding for a string.
type String struct {
	Base
	value string
}

// NewString creates a new binding with the given value.
func NewString(value string) *String {
	return &String{value: value}
}

// Get returns the bound value.
func (b *String) Get() string {
	return b.value
}

// Set updates the bound value.
func (b *String) Set(value string) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddStringListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *String) AddStringListener(listener func(string)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// URL implements a data binding for a *url.URL.
type URL struct {
	Base
	value *url.URL
}

// NewURL creates a new binding with the given value.
func NewURL(value *url.URL) *URL {
	return &URL{value: value}
}

// Get returns the bound value.
func (b *URL) Get() *url.URL {
	return b.value
}

// Set updates the bound value.
func (b *URL) Set(value *url.URL) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddURLListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *URL) AddURLListener(listener func(*url.URL)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}
