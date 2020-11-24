// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// Bool supports binding a bool value in a Fyne application
type Bool interface {
	DataItem
	Get() bool
	Set(bool)
}

// BoolPointer supports binding a bool value in a Fyne application
type BoolPointer interface {
	Bool
	DataPointer
}

// NewBool returns a bindable bool value that is managed internally.
func NewBool() Bool {
	blank := false
	return &boundBool{val: &blank}
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
func BindBool(v *bool) BoolPointer {
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
	if b.val == nil {
		return false
	}
	return *b.val
}

func (b *boundBool) Set(val bool) {
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		if *b.val == val {
			return
		}

		*b.val = val
	}

	b.trigger()
}

func (b *boundBool) Reload() {
	if b.val == nil {
		return
	}

	b.trigger()
}

// Float supports binding a float64 value in a Fyne application
type Float interface {
	DataItem
	Get() float64
	Set(float64)
}

// FloatPointer supports binding a float64 value in a Fyne application
type FloatPointer interface {
	Float
	DataPointer
}

// NewFloat returns a bindable float64 value that is managed internally.
func NewFloat() Float {
	blank := 0.0
	return &boundFloat{val: &blank}
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
func BindFloat(v *float64) FloatPointer {
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
	if b.val == nil {
		return 0.0
	}
	return *b.val
}

func (b *boundFloat) Set(val float64) {
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		if *b.val == val {
			return
		}

		*b.val = val
	}

	b.trigger()
}

func (b *boundFloat) Reload() {
	if b.val == nil {
		return
	}

	b.trigger()
}

// Int supports binding a int value in a Fyne application
type Int interface {
	DataItem
	Get() int
	Set(int)
}

// IntPointer supports binding a int value in a Fyne application
type IntPointer interface {
	Int
	DataPointer
}

// NewInt returns a bindable int value that is managed internally.
func NewInt() Int {
	blank := 0
	return &boundInt{val: &blank}
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
func BindInt(v *int) IntPointer {
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
	if b.val == nil {
		return 0
	}
	return *b.val
}

func (b *boundInt) Set(val int) {
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		if *b.val == val {
			return
		}

		*b.val = val
	}

	b.trigger()
}

func (b *boundInt) Reload() {
	if b.val == nil {
		return
	}

	b.trigger()
}

// Rune supports binding a rune value in a Fyne application
type Rune interface {
	DataItem
	Get() rune
	Set(rune)
}

// RunePointer supports binding a rune value in a Fyne application
type RunePointer interface {
	Rune
	DataPointer
}

// NewRune returns a bindable rune value that is managed internally.
func NewRune() Rune {
	blank := rune(0)
	return &boundRune{val: &blank}
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
func BindRune(v *rune) RunePointer {
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
	if b.val == nil {
		return rune(0)
	}
	return *b.val
}

func (b *boundRune) Set(val rune) {
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		if *b.val == val {
			return
		}

		*b.val = val
	}

	b.trigger()
}

func (b *boundRune) Reload() {
	if b.val == nil {
		return
	}

	b.trigger()
}

// String supports binding a string value in a Fyne application
type String interface {
	DataItem
	Get() string
	Set(string)
}

// StringPointer supports binding a string value in a Fyne application
type StringPointer interface {
	String
	DataPointer
}

// NewString returns a bindable string value that is managed internally.
func NewString() String {
	blank := ""
	return &boundString{val: &blank}
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
func BindString(v *string) StringPointer {
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
	if b.val == nil {
		return ""
	}
	return *b.val
}

func (b *boundString) Set(val string) {
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		if *b.val == val {
			return
		}

		*b.val = val
	}

	b.trigger()
}

func (b *boundString) Reload() {
	if b.val == nil {
		return
	}

	b.trigger()
}
