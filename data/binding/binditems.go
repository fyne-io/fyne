// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// Bool supports binding a bool value.
//
// Since: 2.0.0
type Bool interface {
	DataItem
	Get() (bool, error)
	Set(bool) error
}

// NewBool returns a bindable bool value that is managed internally.
//
// Since: 2.0.0
func NewBool() Bool {
	blank := false
	return &boundBool{val: &blank}
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
//
// Since: 2.0.0
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

func (b *boundBool) Get() (bool, error) {
	if b.val == nil {
		return false, nil
	}
	return *b.val, nil
}

func (b *boundBool) Set(val bool) error {
	if *b.val == val {
		return nil
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
	return nil
}

// Float supports binding a float64 value.
//
// Since: 2.0.0
type Float interface {
	DataItem
	Get() (float64, error)
	Set(float64) error
}

// NewFloat returns a bindable float64 value that is managed internally.
//
// Since: 2.0.0
func NewFloat() Float {
	blank := 0.0
	return &boundFloat{val: &blank}
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
//
// Since: 2.0.0
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

func (b *boundFloat) Get() (float64, error) {
	if b.val == nil {
		return 0.0, nil
	}
	return *b.val, nil
}

func (b *boundFloat) Set(val float64) error {
	if *b.val == val {
		return nil
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
	return nil
}

// Int supports binding a int value.
//
// Since: 2.0.0
type Int interface {
	DataItem
	Get() (int, error)
	Set(int) error
}

// NewInt returns a bindable int value that is managed internally.
//
// Since: 2.0.0
func NewInt() Int {
	blank := 0
	return &boundInt{val: &blank}
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
//
// Since: 2.0.0
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

func (b *boundInt) Get() (int, error) {
	if b.val == nil {
		return 0, nil
	}
	return *b.val, nil
}

func (b *boundInt) Set(val int) error {
	if *b.val == val {
		return nil
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
	return nil
}

// Rune supports binding a rune value.
//
// Since: 2.0.0
type Rune interface {
	DataItem
	Get() (rune, error)
	Set(rune) error
}

// NewRune returns a bindable rune value that is managed internally.
//
// Since: 2.0.0
func NewRune() Rune {
	blank := rune(0)
	return &boundRune{val: &blank}
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
//
// Since: 2.0.0
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

func (b *boundRune) Get() (rune, error) {
	if b.val == nil {
		return rune(0), nil
	}
	return *b.val, nil
}

func (b *boundRune) Set(val rune) error {
	if *b.val == val {
		return nil
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
	return nil
}

// String supports binding a string value.
//
// Since: 2.0.0
type String interface {
	DataItem
	Get() (string, error)
	Set(string) error
}

// NewString returns a bindable string value that is managed internally.
//
// Since: 2.0.0
func NewString() String {
	blank := ""
	return &boundString{val: &blank}
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
//
// Since: 2.0.0
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

func (b *boundString) Get() (string, error) {
	if b.val == nil {
		return "", nil
	}
	return *b.val, nil
}

func (b *boundString) Set(val string) error {
	if *b.val == val {
		return nil
	}
	if b.val == nil { // was not initialized with a blank value, recover
		b.val = &val
	} else {
		*b.val = val
	}

	b.trigger()
	return nil
}
