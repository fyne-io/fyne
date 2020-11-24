// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// Bool supports binding a new bool value in a Fyne application.
type Bool interface {
	DataItem
	Get() bool
	Set(bool)
}

// BoolPointer supports binding to an existing bool value in a Fyne application.
type BoolPointer interface {
	Bool

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewBool returns a bindable bool value that is managed internally.
func NewBool() Bool {
	blank := false
	return &boundBool{val: &blank}
}

// BindBool returns a binding to an existing bool value passed into this function.
// If you modify the bool value directly you should call Reload() to inform the binding.
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

// Float supports binding a new float64 value in a Fyne application.
type Float interface {
	DataItem
	Get() float64
	Set(float64)
}

// FloatPointer supports binding to an existing float64 value in a Fyne application.
type FloatPointer interface {
	Float

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewFloat returns a bindable float64 value that is managed internally.
func NewFloat() Float {
	blank := 0.0
	return &boundFloat{val: &blank}
}

// BindFloat returns a binding to an existing float64 value passed into this function.
// If you modify the float64 value directly you should call Reload() to inform the binding.
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

// Int supports binding a new int value in a Fyne application.
type Int interface {
	DataItem
	Get() int
	Set(int)
}

// IntPointer supports binding to an existing int value in a Fyne application.
type IntPointer interface {
	Int

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewInt returns a bindable int value that is managed internally.
func NewInt() Int {
	blank := 0
	return &boundInt{val: &blank}
}

// BindInt returns a binding to an existing int value passed into this function.
// If you modify the int value directly you should call Reload() to inform the binding.
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

// Rune supports binding a new rune value in a Fyne application.
type Rune interface {
	DataItem
	Get() rune
	Set(rune)
}

// RunePointer supports binding to an existing rune value in a Fyne application.
type RunePointer interface {
	Rune

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewRune returns a bindable rune value that is managed internally.
func NewRune() Rune {
	blank := rune(0)
	return &boundRune{val: &blank}
}

// BindRune returns a binding to an existing rune value passed into this function.
// If you modify the rune value directly you should call Reload() to inform the binding.
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

// String supports binding a new string value in a Fyne application.
type String interface {
	DataItem
	Get() string
	Set(string)
}

// StringPointer supports binding to an existing string value in a Fyne application.
type StringPointer interface {
	String

	// Reload should be called if the value this binding points to has been directly edited and should be reloaded.
	Reload()
}

// NewString returns a bindable string value that is managed internally.
func NewString() String {
	blank := ""
	return &boundString{val: &blank}
}

// BindString returns a binding to an existing string value passed into this function.
// If you modify the string value directly you should call Reload() to inform the binding.
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
