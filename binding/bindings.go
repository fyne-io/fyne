// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"fyne.io/fyne"
	"net/url"
)

// BoolBinding implements a data binding for a bool.
type BoolBinding struct {
	ItemBinding
	value bool
}

// NewBoolBinding creates a new binding with the given value.
func NewBoolBinding(value bool) *BoolBinding {
	return &BoolBinding{value: value}
}

// Get returns the bound value.
func (b *BoolBinding) Get() bool {
	return b.value
}

// Set updates the bound value.
func (b *BoolBinding) Set(value bool) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddBoolListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *BoolBinding) AddBoolListener(listener func(bool)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// Float64Binding implements a data binding for a float64.
type Float64Binding struct {
	ItemBinding
	value float64
}

// NewFloat64Binding creates a new binding with the given value.
func NewFloat64Binding(value float64) *Float64Binding {
	return &Float64Binding{value: value}
}

// Get returns the bound value.
func (b *Float64Binding) Get() float64 {
	return b.value
}

// Set updates the bound value.
func (b *Float64Binding) Set(value float64) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddFloat64Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *Float64Binding) AddFloat64Listener(listener func(float64)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// IntBinding implements a data binding for a int.
type IntBinding struct {
	ItemBinding
	value int
}

// NewIntBinding creates a new binding with the given value.
func NewIntBinding(value int) *IntBinding {
	return &IntBinding{value: value}
}

// Get returns the bound value.
func (b *IntBinding) Get() int {
	return b.value
}

// Set updates the bound value.
func (b *IntBinding) Set(value int) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddIntListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *IntBinding) AddIntListener(listener func(int)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// Int64Binding implements a data binding for a int64.
type Int64Binding struct {
	ItemBinding
	value int64
}

// NewInt64Binding creates a new binding with the given value.
func NewInt64Binding(value int64) *Int64Binding {
	return &Int64Binding{value: value}
}

// Get returns the bound value.
func (b *Int64Binding) Get() int64 {
	return b.value
}

// Set updates the bound value.
func (b *Int64Binding) Set(value int64) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddInt64Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *Int64Binding) AddInt64Listener(listener func(int64)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// ResourceBinding implements a data binding for a fyne.Resource.
type ResourceBinding struct {
	ItemBinding
	value fyne.Resource
}

// NewResourceBinding creates a new binding with the given value.
func NewResourceBinding(value fyne.Resource) *ResourceBinding {
	return &ResourceBinding{value: value}
}

// Get returns the bound value.
func (b *ResourceBinding) Get() fyne.Resource {
	return b.value
}

// Set updates the bound value.
func (b *ResourceBinding) Set(value fyne.Resource) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddResourceListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *ResourceBinding) AddResourceListener(listener func(fyne.Resource)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// RuneBinding implements a data binding for a rune.
type RuneBinding struct {
	ItemBinding
	value rune
}

// NewRuneBinding creates a new binding with the given value.
func NewRuneBinding(value rune) *RuneBinding {
	return &RuneBinding{value: value}
}

// Get returns the bound value.
func (b *RuneBinding) Get() rune {
	return b.value
}

// Set updates the bound value.
func (b *RuneBinding) Set(value rune) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddRuneListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *RuneBinding) AddRuneListener(listener func(rune)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// StringBinding implements a data binding for a string.
type StringBinding struct {
	ItemBinding
	value string
}

// NewStringBinding creates a new binding with the given value.
func NewStringBinding(value string) *StringBinding {
	return &StringBinding{value: value}
}

// Get returns the bound value.
func (b *StringBinding) Get() string {
	return b.value
}

// Set updates the bound value.
func (b *StringBinding) Set(value string) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddStringListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *StringBinding) AddStringListener(listener func(string)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}

// URLBinding implements a data binding for a *url.URL.
type URLBinding struct {
	ItemBinding
	value *url.URL
}

// NewURLBinding creates a new binding with the given value.
func NewURLBinding(value *url.URL) *URLBinding {
	return &URLBinding{value: value}
}

// Get returns the bound value.
func (b *URLBinding) Get() *url.URL {
	return b.value
}

// Set updates the bound value.
func (b *URLBinding) Set(value *url.URL) {
	if b.value != value {
		b.value = value
		b.notify()
	}
}

// AddURLListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *URLBinding) AddURLListener(listener func(*url.URL)) *NotifyFunction {
	return b.AddListenerFunction(func(Binding) {
		listener(b.value)
	})
}
