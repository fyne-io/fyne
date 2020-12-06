// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// BoolList supports binding a list of bool values.
type BoolList interface {
	DataList

	Append(bool)
	Get(int) bool
	Prepend(bool)
	Set(int, bool)
}

// NewBoolList returns a bindable list of bool values.
func NewBoolList() BoolList {
	return &boundBoolList{}
}

type boundBoolList struct {
	listBase
}

func (l *boundBoolList) Append(val bool) {
	l.appendItem(BindBool(&val))
}

func (l *boundBoolList) Get(i int) bool {
	if i > l.Length() {
		return false
	}

	return l.GetItem(i).(Bool).Get()
}

func (l *boundBoolList) Prepend(val bool) {
	l.prependItem(BindBool(&val))
}

func (l *boundBoolList) Set(i int, v bool) {
	if i > l.Length() {
		return
	}

	l.GetItem(i).(Bool).Set(v)
}

// FloatList supports binding a list of float64 values.
type FloatList interface {
	DataList

	Append(float64)
	Get(int) float64
	Prepend(float64)
	Set(int, float64)
}

// NewFloatList returns a bindable list of float64 values.
func NewFloatList() FloatList {
	return &boundFloatList{}
}

type boundFloatList struct {
	listBase
}

func (l *boundFloatList) Append(val float64) {
	l.appendItem(BindFloat(&val))
}

func (l *boundFloatList) Get(i int) float64 {
	if i > l.Length() {
		return 0.0
	}

	return l.GetItem(i).(Float).Get()
}

func (l *boundFloatList) Prepend(val float64) {
	l.prependItem(BindFloat(&val))
}

func (l *boundFloatList) Set(i int, v float64) {
	if i > l.Length() {
		return
	}

	l.GetItem(i).(Float).Set(v)
}

// IntList supports binding a list of int values.
type IntList interface {
	DataList

	Append(int)
	Get(int) int
	Prepend(int)
	Set(int, int)
}

// NewIntList returns a bindable list of int values.
func NewIntList() IntList {
	return &boundIntList{}
}

type boundIntList struct {
	listBase
}

func (l *boundIntList) Append(val int) {
	l.appendItem(BindInt(&val))
}

func (l *boundIntList) Get(i int) int {
	if i > l.Length() {
		return 0
	}

	return l.GetItem(i).(Int).Get()
}

func (l *boundIntList) Prepend(val int) {
	l.prependItem(BindInt(&val))
}

func (l *boundIntList) Set(i int, v int) {
	if i > l.Length() {
		return
	}

	l.GetItem(i).(Int).Set(v)
}

// RuneList supports binding a list of rune values.
type RuneList interface {
	DataList

	Append(rune)
	Get(int) rune
	Prepend(rune)
	Set(int, rune)
}

// NewRuneList returns a bindable list of rune values.
func NewRuneList() RuneList {
	return &boundRuneList{}
}

type boundRuneList struct {
	listBase
}

func (l *boundRuneList) Append(val rune) {
	l.appendItem(BindRune(&val))
}

func (l *boundRuneList) Get(i int) rune {
	if i > l.Length() {
		return rune(0)
	}

	return l.GetItem(i).(Rune).Get()
}

func (l *boundRuneList) Prepend(val rune) {
	l.prependItem(BindRune(&val))
}

func (l *boundRuneList) Set(i int, v rune) {
	if i > l.Length() {
		return
	}

	l.GetItem(i).(Rune).Set(v)
}

// StringList supports binding a list of string values.
type StringList interface {
	DataList

	Append(string)
	Get(int) string
	Prepend(string)
	Set(int, string)
}

// NewStringList returns a bindable list of string values.
func NewStringList() StringList {
	return &boundStringList{}
}

type boundStringList struct {
	listBase
}

func (l *boundStringList) Append(val string) {
	l.appendItem(BindString(&val))
}

func (l *boundStringList) Get(i int) string {
	if i > l.Length() {
		return ""
	}

	return l.GetItem(i).(String).Get()
}

func (l *boundStringList) Prepend(val string) {
	l.prependItem(BindString(&val))
}

func (l *boundStringList) Set(i int, v string) {
	if i > l.Length() {
		return
	}

	l.GetItem(i).(String).Set(v)
}
