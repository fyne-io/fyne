// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"fmt"
	"fyne.io/fyne"
	"net/url"
)

func typeError(e string, a interface{}) {
	fyne.LogError(fmt.Sprintf("Incorrect type: expected '%s', got '%v'", e, a), nil)
}

type BoolBinding struct {
	itemBinding
	value bool
}

func NewBoolBinding(value bool) *BoolBinding {
	return &BoolBinding{value: value}
}

func (b *BoolBinding) Get() bool {
	return b.value
}

func (b *BoolBinding) Set(value bool) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *BoolBinding) AddListener(listener func(bool)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(bool)
		if ok {
			listener(v)
		} else {
			typeError("bool", value)
		}
	})
}

type ByteBinding struct {
	itemBinding
	value byte
}

func NewByteBinding(value byte) *ByteBinding {
	return &ByteBinding{value: value}
}

func (b *ByteBinding) Get() byte {
	return b.value
}

func (b *ByteBinding) Set(value byte) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *ByteBinding) AddListener(listener func(byte)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(byte)
		if ok {
			listener(v)
		} else {
			typeError("byte", value)
		}
	})
}

type Float64Binding struct {
	itemBinding
	value float64
}

func NewFloat64Binding(value float64) *Float64Binding {
	return &Float64Binding{value: value}
}

func (b *Float64Binding) Get() float64 {
	return b.value
}

func (b *Float64Binding) Set(value float64) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *Float64Binding) AddListener(listener func(float64)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(float64)
		if ok {
			listener(v)
		} else {
			typeError("float64", value)
		}
	})
}

type IntBinding struct {
	itemBinding
	value int
}

func NewIntBinding(value int) *IntBinding {
	return &IntBinding{value: value}
}

func (b *IntBinding) Get() int {
	return b.value
}

func (b *IntBinding) Set(value int) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *IntBinding) AddListener(listener func(int)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(int)
		if ok {
			listener(v)
		} else {
			typeError("int", value)
		}
	})
}

type Int64Binding struct {
	itemBinding
	value int64
}

func NewInt64Binding(value int64) *Int64Binding {
	return &Int64Binding{value: value}
}

func (b *Int64Binding) Get() int64 {
	return b.value
}

func (b *Int64Binding) Set(value int64) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *Int64Binding) AddListener(listener func(int64)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(int64)
		if ok {
			listener(v)
		} else {
			typeError("int64", value)
		}
	})
}

type UintBinding struct {
	itemBinding
	value uint
}

func NewUintBinding(value uint) *UintBinding {
	return &UintBinding{value: value}
}

func (b *UintBinding) Get() uint {
	return b.value
}

func (b *UintBinding) Set(value uint) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *UintBinding) AddListener(listener func(uint)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(uint)
		if ok {
			listener(v)
		} else {
			typeError("uint", value)
		}
	})
}

type Uint64Binding struct {
	itemBinding
	value uint64
}

func NewUint64Binding(value uint64) *Uint64Binding {
	return &Uint64Binding{value: value}
}

func (b *Uint64Binding) Get() uint64 {
	return b.value
}

func (b *Uint64Binding) Set(value uint64) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *Uint64Binding) AddListener(listener func(uint64)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(uint64)
		if ok {
			listener(v)
		} else {
			typeError("uint64", value)
		}
	})
}

type ResourceBinding struct {
	itemBinding
	value fyne.Resource
}

func NewResourceBinding(value fyne.Resource) *ResourceBinding {
	return &ResourceBinding{value: value}
}

func (b *ResourceBinding) Get() fyne.Resource {
	return b.value
}

func (b *ResourceBinding) Set(value fyne.Resource) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *ResourceBinding) AddListener(listener func(fyne.Resource)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(fyne.Resource)
		if ok {
			listener(v)
		} else {
			typeError("fyne.Resource", value)
		}
	})
}

type RuneBinding struct {
	itemBinding
	value rune
}

func NewRuneBinding(value rune) *RuneBinding {
	return &RuneBinding{value: value}
}

func (b *RuneBinding) Get() rune {
	return b.value
}

func (b *RuneBinding) Set(value rune) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *RuneBinding) AddListener(listener func(rune)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(rune)
		if ok {
			listener(v)
		} else {
			typeError("rune", value)
		}
	})
}

type StringBinding struct {
	itemBinding
	value string
}

func NewStringBinding(value string) *StringBinding {
	return &StringBinding{value: value}
}

func (b *StringBinding) Get() string {
	return b.value
}

func (b *StringBinding) Set(value string) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *StringBinding) AddListener(listener func(string)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(string)
		if ok {
			listener(v)
		} else {
			typeError("string", value)
		}
	})
}

type URLBinding struct {
	itemBinding
	value *url.URL
}

func NewURLBinding(value *url.URL) *URLBinding {
	return &URLBinding{value: value}
}

func (b *URLBinding) Get() *url.URL {
	return b.value
}

func (b *URLBinding) Set(value *url.URL) {
	if b.value != value {
		b.value = value
		b.notify(value)
	}
}

func (b *URLBinding) AddListener(listener func(*url.URL)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(*url.URL)
		if ok {
			listener(v)
		} else {
			typeError("*url.URL", value)
		}
	})
}
