// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"fmt"
	"fyne.io/fyne"
	"net/url"
)

type BoolBinding struct {
	BaseBinding
	Value bool
}

func (b *BoolBinding) GetBool() bool {
	return b.Value
}

func (b *BoolBinding) Set(value interface{}) {
	v, ok := value.(bool)
	if ok {
		b.SetBool(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'bool', got '%v'", value), nil)
	}
}

func (b *BoolBinding) SetBool(value bool) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *BoolBinding) AddBoolListener(listener func(bool)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(bool)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'bool', got '%v'", value), nil)
		}
	})
}

type ByteBinding struct {
	BaseBinding
	Value byte
}

func (b *ByteBinding) GetByte() byte {
	return b.Value
}

func (b *ByteBinding) Set(value interface{}) {
	v, ok := value.(byte)
	if ok {
		b.SetByte(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'byte', got '%v'", value), nil)
	}
}

func (b *ByteBinding) SetByte(value byte) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *ByteBinding) AddByteListener(listener func(byte)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(byte)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'byte', got '%v'", value), nil)
		}
	})
}

type Float32Binding struct {
	BaseBinding
	Value float32
}

func (b *Float32Binding) GetFloat32() float32 {
	return b.Value
}

func (b *Float32Binding) Set(value interface{}) {
	v, ok := value.(float32)
	if ok {
		b.SetFloat32(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'float32', got '%v'", value), nil)
	}
}

func (b *Float32Binding) SetFloat32(value float32) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Float32Binding) AddFloat32Listener(listener func(float32)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(float32)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'float32', got '%v'", value), nil)
		}
	})
}

type Float64Binding struct {
	BaseBinding
	Value float64
}

func (b *Float64Binding) GetFloat64() float64 {
	return b.Value
}

func (b *Float64Binding) Set(value interface{}) {
	v, ok := value.(float64)
	if ok {
		b.SetFloat64(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'float64', got '%v'", value), nil)
	}
}

func (b *Float64Binding) SetFloat64(value float64) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Float64Binding) AddFloat64Listener(listener func(float64)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(float64)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'float64', got '%v'", value), nil)
		}
	})
}

type IntBinding struct {
	BaseBinding
	Value int
}

func (b *IntBinding) GetInt() int {
	return b.Value
}

func (b *IntBinding) Set(value interface{}) {
	v, ok := value.(int)
	if ok {
		b.SetInt(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int', got '%v'", value), nil)
	}
}

func (b *IntBinding) SetInt(value int) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *IntBinding) AddIntListener(listener func(int)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(int)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int', got '%v'", value), nil)
		}
	})
}

type Int8Binding struct {
	BaseBinding
	Value int8
}

func (b *Int8Binding) GetInt8() int8 {
	return b.Value
}

func (b *Int8Binding) Set(value interface{}) {
	v, ok := value.(int8)
	if ok {
		b.SetInt8(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int8', got '%v'", value), nil)
	}
}

func (b *Int8Binding) SetInt8(value int8) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Int8Binding) AddInt8Listener(listener func(int8)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(int8)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int8', got '%v'", value), nil)
		}
	})
}

type Int16Binding struct {
	BaseBinding
	Value int16
}

func (b *Int16Binding) GetInt16() int16 {
	return b.Value
}

func (b *Int16Binding) Set(value interface{}) {
	v, ok := value.(int16)
	if ok {
		b.SetInt16(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int16', got '%v'", value), nil)
	}
}

func (b *Int16Binding) SetInt16(value int16) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Int16Binding) AddInt16Listener(listener func(int16)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(int16)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int16', got '%v'", value), nil)
		}
	})
}

type Int32Binding struct {
	BaseBinding
	Value int32
}

func (b *Int32Binding) GetInt32() int32 {
	return b.Value
}

func (b *Int32Binding) Set(value interface{}) {
	v, ok := value.(int32)
	if ok {
		b.SetInt32(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int32', got '%v'", value), nil)
	}
}

func (b *Int32Binding) SetInt32(value int32) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Int32Binding) AddInt32Listener(listener func(int32)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(int32)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int32', got '%v'", value), nil)
		}
	})
}

type Int64Binding struct {
	BaseBinding
	Value int64
}

func (b *Int64Binding) GetInt64() int64 {
	return b.Value
}

func (b *Int64Binding) Set(value interface{}) {
	v, ok := value.(int64)
	if ok {
		b.SetInt64(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int64', got '%v'", value), nil)
	}
}

func (b *Int64Binding) SetInt64(value int64) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Int64Binding) AddInt64Listener(listener func(int64)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(int64)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'int64', got '%v'", value), nil)
		}
	})
}

type UintBinding struct {
	BaseBinding
	Value uint
}

func (b *UintBinding) GetUint() uint {
	return b.Value
}

func (b *UintBinding) Set(value interface{}) {
	v, ok := value.(uint)
	if ok {
		b.SetUint(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint', got '%v'", value), nil)
	}
}

func (b *UintBinding) SetUint(value uint) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *UintBinding) AddUintListener(listener func(uint)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(uint)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint', got '%v'", value), nil)
		}
	})
}

type Uint8Binding struct {
	BaseBinding
	Value uint8
}

func (b *Uint8Binding) GetUint8() uint8 {
	return b.Value
}

func (b *Uint8Binding) Set(value interface{}) {
	v, ok := value.(uint8)
	if ok {
		b.SetUint8(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint8', got '%v'", value), nil)
	}
}

func (b *Uint8Binding) SetUint8(value uint8) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Uint8Binding) AddUint8Listener(listener func(uint8)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(uint8)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint8', got '%v'", value), nil)
		}
	})
}

type Uint16Binding struct {
	BaseBinding
	Value uint16
}

func (b *Uint16Binding) GetUint16() uint16 {
	return b.Value
}

func (b *Uint16Binding) Set(value interface{}) {
	v, ok := value.(uint16)
	if ok {
		b.SetUint16(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint16', got '%v'", value), nil)
	}
}

func (b *Uint16Binding) SetUint16(value uint16) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Uint16Binding) AddUint16Listener(listener func(uint16)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(uint16)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint16', got '%v'", value), nil)
		}
	})
}

type Uint32Binding struct {
	BaseBinding
	Value uint32
}

func (b *Uint32Binding) GetUint32() uint32 {
	return b.Value
}

func (b *Uint32Binding) Set(value interface{}) {
	v, ok := value.(uint32)
	if ok {
		b.SetUint32(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint32', got '%v'", value), nil)
	}
}

func (b *Uint32Binding) SetUint32(value uint32) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Uint32Binding) AddUint32Listener(listener func(uint32)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(uint32)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint32', got '%v'", value), nil)
		}
	})
}

type Uint64Binding struct {
	BaseBinding
	Value uint64
}

func (b *Uint64Binding) GetUint64() uint64 {
	return b.Value
}

func (b *Uint64Binding) Set(value interface{}) {
	v, ok := value.(uint64)
	if ok {
		b.SetUint64(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint64', got '%v'", value), nil)
	}
}

func (b *Uint64Binding) SetUint64(value uint64) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *Uint64Binding) AddUint64Listener(listener func(uint64)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(uint64)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'uint64', got '%v'", value), nil)
		}
	})
}

type ResourceBinding struct {
	BaseBinding
	Value fyne.Resource
}

func (b *ResourceBinding) GetResource() fyne.Resource {
	return b.Value
}

func (b *ResourceBinding) Set(value interface{}) {
	v, ok := value.(fyne.Resource)
	if ok {
		b.SetResource(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'fyne.Resource', got '%v'", value), nil)
	}
}

func (b *ResourceBinding) SetResource(value fyne.Resource) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *ResourceBinding) AddResourceListener(listener func(fyne.Resource)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(fyne.Resource)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'fyne.Resource', got '%v'", value), nil)
		}
	})
}

type RuneBinding struct {
	BaseBinding
	Value rune
}

func (b *RuneBinding) GetRune() rune {
	return b.Value
}

func (b *RuneBinding) Set(value interface{}) {
	v, ok := value.(rune)
	if ok {
		b.SetRune(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'rune', got '%v'", value), nil)
	}
}

func (b *RuneBinding) SetRune(value rune) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *RuneBinding) AddRuneListener(listener func(rune)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(rune)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'rune', got '%v'", value), nil)
		}
	})
}

type StringBinding struct {
	BaseBinding
	Value string
}

func (b *StringBinding) GetString() string {
	return b.Value
}

func (b *StringBinding) Set(value interface{}) {
	v, ok := value.(string)
	if ok {
		b.SetString(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected 'string', got '%v'", value), nil)
	}
}

func (b *StringBinding) SetString(value string) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *StringBinding) AddStringListener(listener func(string)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(string)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected 'string', got '%v'", value), nil)
		}
	})
}

type URLBinding struct {
	BaseBinding
	Value *url.URL
}

func (b *URLBinding) GetURL() *url.URL {
	return b.Value
}

func (b *URLBinding) Set(value interface{}) {
	v, ok := value.(*url.URL)
	if ok {
		b.SetURL(v)
	} else {
		fyne.LogError(fmt.Sprintf("Incorrect type: expected '*url.URL', got '%v'", value), nil)
	}
}

func (b *URLBinding) SetURL(value *url.URL) {
	if b.Value != value {
		b.Value = value
		b.notify(value)
	}
}

func (b *URLBinding) AddURLListener(listener func(*url.URL)) {
	b.addListener(func(value interface{}) {
		v, ok := value.(*url.URL)
		if ok {
			listener(v)
		} else {
			fyne.LogError(fmt.Sprintf("Incorrect type: expected '*url.URL', got '%v'", value), nil)
		}
	})
}
