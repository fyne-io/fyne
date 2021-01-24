// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

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
	blank := false
	return &boundBool{val: &blank}
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindBool(v *bool) ExternalBool {
	if v == nil {
		return NewBool().(ExternalBool) // never allow a nil value pointer
	}

	return &boundBool{val: v}
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

func (b *boundBool) Reload() error {
	return b.Set(*b.val)
}

func (b *boundBool) Set(val bool) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	*b.val = val

	b.trigger()
	return nil
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
	blank := 0.0
	return &boundFloat{val: &blank}
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindFloat(v *float64) ExternalFloat {
	if v == nil {
		return NewFloat().(ExternalFloat) // never allow a nil value pointer
	}

	return &boundFloat{val: v}
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

func (b *boundFloat) Reload() error {
	return b.Set(*b.val)
}

func (b *boundFloat) Set(val float64) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	*b.val = val

	b.trigger()
	return nil
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
	blank := 0
	return &boundInt{val: &blank}
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindInt(v *int) ExternalInt {
	if v == nil {
		return NewInt().(ExternalInt) // never allow a nil value pointer
	}

	return &boundInt{val: v}
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

func (b *boundInt) Reload() error {
	return b.Set(*b.val)
}

func (b *boundInt) Set(val int) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	*b.val = val

	b.trigger()
	return nil
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
	blank := rune(0)
	return &boundRune{val: &blank}
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindRune(v *rune) ExternalRune {
	if v == nil {
		return NewRune().(ExternalRune) // never allow a nil value pointer
	}

	return &boundRune{val: v}
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

func (b *boundRune) Reload() error {
	return b.Set(*b.val)
}

func (b *boundRune) Set(val rune) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	*b.val = val

	b.trigger()
	return nil
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
	blank := ""
	return &boundString{val: &blank}
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindString(v *string) ExternalString {
	if v == nil {
		return NewString().(ExternalString) // never allow a nil value pointer
	}

	return &boundString{val: v}
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

func (b *boundString) Reload() error {
	return b.Set(*b.val)
}

func (b *boundString) Set(val string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	*b.val = val

	b.trigger()
	return nil
}
