// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

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
	var blank bool = false
	return &boundBool{val: &blank}
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindBool(v *bool) ExternalBool {
	if v == nil {
		var blank bool = false
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalBool{}
	b.val = v
	b.old = *v
	return b
}

type boundBool struct {
	base

	val *bool
}

func (b *boundBool) Get() (bool, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return false, nil
	}
	return *b.val, nil
}

func (b *boundBool) Set(val bool) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return nil
	}
	*b.val = val

	b.trigger()
	return nil
}

type boundExternalBool struct {
	boundBool

	old bool
}

func (b *boundExternalBool) Set(val bool) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return nil
	}
	*b.val = val
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternalBool) Reload() error {
	return b.Set(*b.val)
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
	var blank []byte = nil
	return &boundBytes{val: &blank}
}

// BindBytes returns a new bindable value that controls the contents of the provided []byte variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.2
func BindBytes(v *[]byte) ExternalBytes {
	if v == nil {
		var blank []byte = nil
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalBytes{}
	b.val = v
	b.old = *v
	return b
}

type boundBytes struct {
	base

	val *[]byte
}

func (b *boundBytes) Get() ([]byte, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return nil, nil
	}
	return *b.val, nil
}

func (b *boundBytes) Set(val []byte) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if bytes.Equal(*b.val, val) {
		return nil
	}
	*b.val = val

	b.trigger()
	return nil
}

type boundExternalBytes struct {
	boundBytes

	old []byte
}

func (b *boundExternalBytes) Set(val []byte) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if bytes.Equal(b.old, val) {
		return nil
	}
	*b.val = val
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternalBytes) Reload() error {
	return b.Set(*b.val)
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
	var blank float64 = 0.0
	return &boundFloat{val: &blank}
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindFloat(v *float64) ExternalFloat {
	if v == nil {
		var blank float64 = 0.0
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalFloat{}
	b.val = v
	b.old = *v
	return b
}

type boundFloat struct {
	base

	val *float64
}

func (b *boundFloat) Get() (float64, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return 0.0, nil
	}
	return *b.val, nil
}

func (b *boundFloat) Set(val float64) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return nil
	}
	*b.val = val

	b.trigger()
	return nil
}

type boundExternalFloat struct {
	boundFloat

	old float64
}

func (b *boundExternalFloat) Set(val float64) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return nil
	}
	*b.val = val
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternalFloat) Reload() error {
	return b.Set(*b.val)
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
	var blank int = 0
	return &boundInt{val: &blank}
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindInt(v *int) ExternalInt {
	if v == nil {
		var blank int = 0
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalInt{}
	b.val = v
	b.old = *v
	return b
}

type boundInt struct {
	base

	val *int
}

func (b *boundInt) Get() (int, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return 0, nil
	}
	return *b.val, nil
}

func (b *boundInt) Set(val int) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return nil
	}
	*b.val = val

	b.trigger()
	return nil
}

type boundExternalInt struct {
	boundInt

	old int
}

func (b *boundExternalInt) Set(val int) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return nil
	}
	*b.val = val
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternalInt) Reload() error {
	return b.Set(*b.val)
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
	var blank rune = rune(0)
	return &boundRune{val: &blank}
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindRune(v *rune) ExternalRune {
	if v == nil {
		var blank rune = rune(0)
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalRune{}
	b.val = v
	b.old = *v
	return b
}

type boundRune struct {
	base

	val *rune
}

func (b *boundRune) Get() (rune, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return rune(0), nil
	}
	return *b.val, nil
}

func (b *boundRune) Set(val rune) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return nil
	}
	*b.val = val

	b.trigger()
	return nil
}

type boundExternalRune struct {
	boundRune

	old rune
}

func (b *boundExternalRune) Set(val rune) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return nil
	}
	*b.val = val
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternalRune) Reload() error {
	return b.Set(*b.val)
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
	var blank string = ""
	return &boundString{val: &blank}
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindString(v *string) ExternalString {
	if v == nil {
		var blank string = ""
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalString{}
	b.val = v
	b.old = *v
	return b
}

type boundString struct {
	base

	val *string
}

func (b *boundString) Get() (string, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return "", nil
	}
	return *b.val, nil
}

func (b *boundString) Set(val string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return nil
	}
	*b.val = val

	b.trigger()
	return nil
}

type boundExternalString struct {
	boundString

	old string
}

func (b *boundExternalString) Set(val string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return nil
	}
	*b.val = val
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternalString) Reload() error {
	return b.Set(*b.val)
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
	var blank fyne.URI = fyne.URI(nil)
	return &boundURI{val: &blank}
}

// BindURI returns a new bindable value that controls the contents of the provided fyne.URI variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindURI(v *fyne.URI) ExternalURI {
	if v == nil {
		var blank fyne.URI = fyne.URI(nil)
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalURI{}
	b.val = v
	b.old = *v
	return b
}

type boundURI struct {
	base

	val *fyne.URI
}

func (b *boundURI) Get() (fyne.URI, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return fyne.URI(nil), nil
	}
	return *b.val, nil
}

func (b *boundURI) Set(val fyne.URI) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if compareURI(*b.val, val) {
		return nil
	}
	*b.val = val

	b.trigger()
	return nil
}

type boundExternalURI struct {
	boundURI

	old fyne.URI
}

func (b *boundExternalURI) Set(val fyne.URI) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if compareURI(b.old, val) {
		return nil
	}
	*b.val = val
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternalURI) Reload() error {
	return b.Set(*b.val)
}
