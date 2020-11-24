// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// Bool supports binding a bool value in a Fyne application
type Bool interface {
	DataItem
	Get() bool
	Set(bool)
}

// BoolFromPointer supports binding a bool value in a Fyne application
type BoolFromPointer interface {
	Bool

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewBool returns a bindable bool value that is managed internally.
func NewBool() Bool {
	blank := false
	return &boundBool{val: &blank}
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
func BindBool(v *bool) BoolFromPointer {
	if v == nil {
		return NewBool().(*boundBool) // never allow a nil value pointer
	}

	return &boundBool{val: v}
}

type boundBool struct {
	base

	val *bool
}

func (b *boundBool) Get() bool {
	return *b.val
}

func (b *boundBool) Set(val bool) {
	if *b.val == val {
		return
	}

	*b.val = val
	b.trigger()
}

func (b *boundBool) Reload() {
	b.trigger()
}

// Float supports binding a float64 value in a Fyne application
type Float interface {
	DataItem
	Get() float64
	Set(float64)
}

// FloatFromPointer supports binding a float64 value in a Fyne application
type FloatFromPointer interface {
	Float

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewFloat returns a bindable float64 value that is managed internally.
func NewFloat() Float {
	blank := 0.0
	return &boundFloat{val: &blank}
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
func BindFloat(v *float64) FloatFromPointer {
	if v == nil {
		return NewFloat().(*boundFloat) // never allow a nil value pointer
	}

	return &boundFloat{val: v}
}

type boundFloat struct {
	base

	val *float64
}

func (b *boundFloat) Get() float64 {
	return *b.val
}

func (b *boundFloat) Set(val float64) {
	if *b.val == val {
		return
	}

	*b.val = val
	b.trigger()
}

func (b *boundFloat) Reload() {
	b.trigger()
}

// Int supports binding a int value in a Fyne application
type Int interface {
	DataItem
	Get() int
	Set(int)
}

// IntFromPointer supports binding a int value in a Fyne application
type IntFromPointer interface {
	Int

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewInt returns a bindable int value that is managed internally.
func NewInt() Int {
	blank := 0
	return &boundInt{val: &blank}
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
func BindInt(v *int) IntFromPointer {
	if v == nil {
		return NewInt().(*boundInt) // never allow a nil value pointer
	}

	return &boundInt{val: v}
}

type boundInt struct {
	base

	val *int
}

func (b *boundInt) Get() int {
	return *b.val
}

func (b *boundInt) Set(val int) {
	if *b.val == val {
		return
	}

	*b.val = val
	b.trigger()
}

func (b *boundInt) Reload() {
	b.trigger()
}

// Rune supports binding a rune value in a Fyne application
type Rune interface {
	DataItem
	Get() rune
	Set(rune)
}

// RuneFromPointer supports binding a rune value in a Fyne application
type RuneFromPointer interface {
	Rune

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewRune returns a bindable rune value that is managed internally.
func NewRune() Rune {
	blank := rune(0)
	return &boundRune{val: &blank}
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
func BindRune(v *rune) RuneFromPointer {
	if v == nil {
		return NewRune().(*boundRune) // never allow a nil value pointer
	}

	return &boundRune{val: v}
}

type boundRune struct {
	base

	val *rune
}

func (b *boundRune) Get() rune {
	return *b.val
}

func (b *boundRune) Set(val rune) {
	if *b.val == val {
		return
	}

	*b.val = val
	b.trigger()
}

func (b *boundRune) Reload() {
	b.trigger()
}

// String supports binding a string value in a Fyne application
type String interface {
	DataItem
	Get() string
	Set(string)
}

// StringFromPointer supports binding a string value in a Fyne application
type StringFromPointer interface {
	String

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewString returns a bindable string value that is managed internally.
func NewString() String {
	blank := ""
	return &boundString{val: &blank}
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
func BindString(v *string) StringFromPointer {
	if v == nil {
		return NewString().(*boundString) // never allow a nil value pointer
	}

	return &boundString{val: v}
}

type boundString struct {
	base

	val *string
}

func (b *boundString) Get() string {
	return *b.val
}

func (b *boundString) Set(val string) {
	if *b.val == val {
		return
	}

	*b.val = val
	b.trigger()
}

func (b *boundString) Reload() {
	b.trigger()
}
