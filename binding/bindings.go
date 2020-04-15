// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"net/url"

	"fyne.io/fyne"
)

// Bool defines a data binding for a bool.
type Bool interface {
	Binding
	Get() bool
	Set(bool)
	AddBoolListener(func(bool)) *NotifyFunction
}

// baseBool implements a data binding for a bool.
type baseBool struct {
	Base
	value bool
}

// NewBool creates a new binding with the given value.
func NewBool(value bool) Bool {
	return &baseBool{value: value}
}

// Get returns the bound value.
func (b *baseBool) Get() bool {
	return b.value
}

// Set updates the bound value.
func (b *baseBool) Set(value bool) {
	if b.value != value {
		b.value = value
		b.Update()
	}
}

// AddBoolListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseBool) AddBoolListener(listener func(bool)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(b.value)
	})
	b.AddListener(n)
	return n
}

// Float64 defines a data binding for a float64.
type Float64 interface {
	Binding
	Get() float64
	Set(float64)
	AddFloat64Listener(func(float64)) *NotifyFunction
}

// baseFloat64 implements a data binding for a float64.
type baseFloat64 struct {
	Base
	value float64
}

// NewFloat64 creates a new binding with the given value.
func NewFloat64(value float64) Float64 {
	return &baseFloat64{value: value}
}

// Get returns the bound value.
func (b *baseFloat64) Get() float64 {
	return b.value
}

// Set updates the bound value.
func (b *baseFloat64) Set(value float64) {
	if b.value != value {
		b.value = value
		b.Update()
	}
}

// AddFloat64Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseFloat64) AddFloat64Listener(listener func(float64)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(b.value)
	})
	b.AddListener(n)
	return n
}

// Int defines a data binding for a int.
type Int interface {
	Binding
	Get() int
	Set(int)
	AddIntListener(func(int)) *NotifyFunction
}

// baseInt implements a data binding for a int.
type baseInt struct {
	Base
	value int
}

// NewInt creates a new binding with the given value.
func NewInt(value int) Int {
	return &baseInt{value: value}
}

// Get returns the bound value.
func (b *baseInt) Get() int {
	return b.value
}

// Set updates the bound value.
func (b *baseInt) Set(value int) {
	if b.value != value {
		b.value = value
		b.Update()
	}
}

// AddIntListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseInt) AddIntListener(listener func(int)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(b.value)
	})
	b.AddListener(n)
	return n
}

// Int64 defines a data binding for a int64.
type Int64 interface {
	Binding
	Get() int64
	Set(int64)
	AddInt64Listener(func(int64)) *NotifyFunction
}

// baseInt64 implements a data binding for a int64.
type baseInt64 struct {
	Base
	value int64
}

// NewInt64 creates a new binding with the given value.
func NewInt64(value int64) Int64 {
	return &baseInt64{value: value}
}

// Get returns the bound value.
func (b *baseInt64) Get() int64 {
	return b.value
}

// Set updates the bound value.
func (b *baseInt64) Set(value int64) {
	if b.value != value {
		b.value = value
		b.Update()
	}
}

// AddInt64Listener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseInt64) AddInt64Listener(listener func(int64)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(b.value)
	})
	b.AddListener(n)
	return n
}

// Resource defines a data binding for a fyne.Resource.
type Resource interface {
	Binding
	Get() fyne.Resource
	Set(fyne.Resource)
	AddResourceListener(func(fyne.Resource)) *NotifyFunction
}

// baseResource implements a data binding for a fyne.Resource.
type baseResource struct {
	Base
	value fyne.Resource
}

// NewResource creates a new binding with the given value.
func NewResource(value fyne.Resource) Resource {
	return &baseResource{value: value}
}

// Get returns the bound value.
func (b *baseResource) Get() fyne.Resource {
	return b.value
}

// Set updates the bound value.
func (b *baseResource) Set(value fyne.Resource) {
	if b.value != value {
		b.value = value
		b.Update()
	}
}

// AddResourceListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseResource) AddResourceListener(listener func(fyne.Resource)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(b.value)
	})
	b.AddListener(n)
	return n
}

// Rune defines a data binding for a rune.
type Rune interface {
	Binding
	Get() rune
	Set(rune)
	AddRuneListener(func(rune)) *NotifyFunction
}

// baseRune implements a data binding for a rune.
type baseRune struct {
	Base
	value rune
}

// NewRune creates a new binding with the given value.
func NewRune(value rune) Rune {
	return &baseRune{value: value}
}

// Get returns the bound value.
func (b *baseRune) Get() rune {
	return b.value
}

// Set updates the bound value.
func (b *baseRune) Set(value rune) {
	if b.value != value {
		b.value = value
		b.Update()
	}
}

// AddRuneListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseRune) AddRuneListener(listener func(rune)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(b.value)
	})
	b.AddListener(n)
	return n
}

// String defines a data binding for a string.
type String interface {
	Binding
	Get() string
	Set(string)
	AddStringListener(func(string)) *NotifyFunction
}

// baseString implements a data binding for a string.
type baseString struct {
	Base
	value string
}

// NewString creates a new binding with the given value.
func NewString(value string) String {
	return &baseString{value: value}
}

// Get returns the bound value.
func (b *baseString) Get() string {
	return b.value
}

// Set updates the bound value.
func (b *baseString) Set(value string) {
	if b.value != value {
		b.value = value
		b.Update()
	}
}

// AddStringListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseString) AddStringListener(listener func(string)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(b.value)
	})
	b.AddListener(n)
	return n
}

// URL defines a data binding for a *url.URL.
type URL interface {
	Binding
	Get() *url.URL
	Set(*url.URL)
	AddURLListener(func(*url.URL)) *NotifyFunction
}

// baseURL implements a data binding for a *url.URL.
type baseURL struct {
	Base
	value *url.URL
}

// NewURL creates a new binding with the given value.
func NewURL(value *url.URL) URL {
	return &baseURL{value: value}
}

// Get returns the bound value.
func (b *baseURL) Get() *url.URL {
	return b.value
}

// Set updates the bound value.
func (b *baseURL) Set(value *url.URL) {
	if b.value != value {
		b.value = value
		b.Update()
	}
}

// AddURLListener adds the given function as a listener to the binding.
// The function is wrapped in the returned NotifyFunction which can be passed to DeleteListener.
func (b *baseURL) AddURLListener(listener func(*url.URL)) *NotifyFunction {
	n := NewNotifyFunction(func(Binding) {
		listener(b.value)
	})
	b.AddListener(n)
	return n
}
