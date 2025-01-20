package binding

import (
	"bytes"

	"fyne.io/fyne/v2"
)

// Bool supports binding a bool value.
//
// Since: 2.0
type Bool interface {
	DataItem
	Get() (bool, error)
	Set(bool) error
}

// ExternalBool supports binding a bool value to an external value.
//
// Since: 2.0
type ExternalBool interface {
	Bool
	Reload() error
}

// NewBool returns a bindable bool value that is managed internally.
//
// Since: 2.0
func NewBool() Bool {
	return newBaseItemComparable[bool]()
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindBool(v *bool) ExternalBool {
	return baseBindExternalComparable(v)
}

// Bytes supports binding a []byte value.
//
// Since: 2.2
type Bytes interface {
	DataItem
	Get() ([]byte, error)
	Set([]byte) error
}

// ExternalBytes supports binding a []byte value to an external value.
//
// Since: 2.2
type ExternalBytes interface {
	Bytes
	Reload() error
}

// NewBytes returns a bindable []byte value that is managed internally.
//
// Since: 2.2
func NewBytes() Bytes {
	return newBaseItem(bytes.Equal)
}

// BindBytes returns a new bindable value that controls the contents of the provided []byte variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.2
func BindBytes(v *[]byte) ExternalBytes {
	return baseBindExternal(v, bytes.Equal)
}

// Float supports binding a float64 value.
//
// Since: 2.0
type Float interface {
	DataItem
	Get() (float64, error)
	Set(float64) error
}

// ExternalFloat supports binding a float64 value to an external value.
//
// Since: 2.0
type ExternalFloat interface {
	Float
	Reload() error
}

// NewFloat returns a bindable float64 value that is managed internally.
//
// Since: 2.0
func NewFloat() Float {
	return newBaseItemComparable[float64]()
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindFloat(v *float64) ExternalFloat {
	return baseBindExternalComparable(v)
}

// Int supports binding a int value.
//
// Since: 2.0
type Int interface {
	DataItem
	Get() (int, error)
	Set(int) error
}

// ExternalInt supports binding a int value to an external value.
//
// Since: 2.0
type ExternalInt interface {
	Int
	Reload() error
}

// NewInt returns a bindable int value that is managed internally.
//
// Since: 2.0
func NewInt() Int {
	return newBaseItemComparable[int]()
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindInt(v *int) ExternalInt {
	return baseBindExternalComparable(v)
}

// Rune supports binding a rune value.
//
// Since: 2.0
type Rune interface {
	DataItem
	Get() (rune, error)
	Set(rune) error
}

// ExternalRune supports binding a rune value to an external value.
//
// Since: 2.0
type ExternalRune interface {
	Rune
	Reload() error
}

// NewRune returns a bindable rune value that is managed internally.
//
// Since: 2.0
func NewRune() Rune {
	return newBaseItemComparable[rune]()
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindRune(v *rune) ExternalRune {
	return baseBindExternalComparable(v)
}

// String supports binding a string value.
//
// Since: 2.0
type String interface {
	DataItem
	Get() (string, error)
	Set(string) error
}

// ExternalString supports binding a string value to an external value.
//
// Since: 2.0
type ExternalString interface {
	String
	Reload() error
}

// NewString returns a bindable string value that is managed internally.
//
// Since: 2.0
func NewString() String {
	return newBaseItemComparable[string]()
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindString(v *string) ExternalString {
	return baseBindExternalComparable(v)
}

// URI supports binding a fyne.URI value.
//
// Since: 2.1
type URI interface {
	DataItem
	Get() (fyne.URI, error)
	Set(fyne.URI) error
}

// ExternalURI supports binding a fyne.URI value to an external value.
//
// Since: 2.1
type ExternalURI interface {
	URI
	Reload() error
}

// NewURI returns a bindable fyne.URI value that is managed internally.
//
// Since: 2.1
func NewURI() URI {
	return newBaseItem(compareURI)
}

// BindURI returns a new bindable value that controls the contents of the provided fyne.URI variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindURI(v *fyne.URI) ExternalURI {
	return baseBindExternal(v, compareURI)
}
