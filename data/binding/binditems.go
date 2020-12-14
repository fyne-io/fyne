// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// Bool supports binding a bool value in a Fyne application
type Bool interface {
	DataItem
	Get() bool
	Set(bool)
}

// NewBool returns a bindable bool value that is managed internally.
func NewBool() Bool {
	blank := false
	return &boundBool{val: &blank}
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
func BindBool(v *bool) Bool {
	if v == nil {
		return NewBool() // never allow a nil value pointer
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
	if *b.val == val {
		return
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
}

// Float supports binding a float64 value in a Fyne application
type Float interface {
	DataItem
	Get() float64
	Set(float64)
}

// NewFloat returns a bindable float64 value that is managed internally.
func NewFloat() Float {
	blank := 0.0
	return &boundFloat{val: &blank}
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
func BindFloat(v *float64) Float {
	if v == nil {
		return NewFloat() // never allow a nil value pointer
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
	if *b.val == val {
		return
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
}

// Int supports binding a int value in a Fyne application
type Int interface {
	DataItem
	Get() int
	Set(int)
}

// NewInt returns a bindable int value that is managed internally.
func NewInt() Int {
	blank := 0
	return &boundInt{val: &blank}
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
func BindInt(v *int) Int {
	if v == nil {
		return NewInt() // never allow a nil value pointer
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
	if *b.val == val {
		return
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
}

// Rune supports binding a rune value in a Fyne application
type Rune interface {
	DataItem
	Get() rune
	Set(rune)
}

// NewRune returns a bindable rune value that is managed internally.
func NewRune() Rune {
	blank := rune(0)
	return &boundRune{val: &blank}
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
func BindRune(v *rune) Rune {
	if v == nil {
		return NewRune() // never allow a nil value pointer
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
	if *b.val == val {
		return
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
}

// String supports binding a string value in a Fyne application
type String interface {
	DataItem
	Get() string
	Set(string)
}

// NewString returns a bindable string value that is managed internally.
func NewString() String {
	blank := ""
	return &boundString{val: &blank}
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
func BindString(v *string) String {
	if v == nil {
		return NewString() // never allow a nil value pointer
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
	if *b.val == val {
		return
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
}
