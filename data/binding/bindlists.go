// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// BoolList supports binding a list of bool values.
//
// Since: 2.0.0
type BoolList interface {
	DataList

	Append(bool)
	Get(int) bool
	Prepend(bool)
	Set(int, bool)
}

// NewBoolList returns a bindable list of bool values.
//
// Since: 2.0.0
func NewBoolList() BoolList {
	return &boundBoolList{}
}

// BindBoolList returns a bound list of bool values, based on the contents of the passed slice.
//
// Since: 2.0.0
func BindBoolList(v *[]bool) BoolList {
	if v == nil {
		return NewBoolList()
	}

	b := &boundBoolList{val: v}

	for _, i := range *v {
		b.appendItem(BindBool(&i))
	}

	return b
}

type boundBoolList struct {
	listBase

	val *[]bool
}

func (l *boundBoolList) Append(val bool) {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindBool(&val))
}

func (l *boundBoolList) Get(i int) bool {
	if i < 0 || i > l.Length() {
		return false
	}
	if l.val != nil {
		return (*l.val)[i]
	}

	return l.GetItem(i).(Bool).Get()
}

func (l *boundBoolList) Prepend(val bool) {
	if l.val != nil {
		*l.val = append([]bool{val}, *l.val...)
	}

	l.prependItem(BindBool(&val))
}

func (l *boundBoolList) Set(i int, v bool) {
	if i > l.Length() {
		return
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	l.GetItem(i).(Bool).Set(v)
}

// FloatList supports binding a list of float64 values.
//
// Since: 2.0.0
type FloatList interface {
	DataList

	Append(float64)
	Get(int) float64
	Prepend(float64)
	Set(int, float64)
}

// NewFloatList returns a bindable list of float64 values.
//
// Since: 2.0.0
func NewFloatList() FloatList {
	return &boundFloatList{}
}

// BindFloatList returns a bound list of float64 values, based on the contents of the passed slice.
//
// Since: 2.0.0
func BindFloatList(v *[]float64) FloatList {
	if v == nil {
		return NewFloatList()
	}

	b := &boundFloatList{val: v}

	for i := range *v {
		b.appendItem(BindFloat(&((*v)[i])))
	}

	return b
}

type boundFloatList struct {
	listBase

	val *[]float64
}

func (l *boundFloatList) Append(val float64) {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindFloat(&val))
}

func (l *boundFloatList) Get(i int) float64 {
	if i < 0 || i > l.Length() {
		return 0.0
	}
	if l.val != nil {
		return (*l.val)[i]
	}

	return l.GetItem(i).(Float).Get()
}

func (l *boundFloatList) Prepend(val float64) {
	if l.val != nil {
		*l.val = append([]float64{val}, *l.val...)
	}

	l.prependItem(BindFloat(&val))
}

func (l *boundFloatList) Set(i int, v float64) {
	if i > l.Length() {
		return
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	l.GetItem(i).(Float).Set(v)
}

// IntList supports binding a list of int values.
//
// Since: 2.0.0
type IntList interface {
	DataList

	Append(int)
	Get(int) int
	Prepend(int)
	Set(int, int)
}

// NewIntList returns a bindable list of int values.
//
// Since: 2.0.0
func NewIntList() IntList {
	return &boundIntList{}
}

// BindIntList returns a bound list of int values, based on the contents of the passed slice.
//
// Since: 2.0.0
func BindIntList(v *[]int) IntList {
	if v == nil {
		return NewIntList()
	}

	b := &boundIntList{val: v}

	for _, i := range *v {
		b.appendItem(BindInt(&i))
	}

	return b
}

type boundIntList struct {
	listBase

	val *[]int
}

func (l *boundIntList) Append(val int) {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindInt(&val))
}

func (l *boundIntList) Get(i int) int {
	if i < 0 || i > l.Length() {
		return 0
	}
	if l.val != nil {
		return (*l.val)[i]
	}

	return l.GetItem(i).(Int).Get()
}

func (l *boundIntList) Prepend(val int) {
	if l.val != nil {
		*l.val = append([]int{val}, *l.val...)
	}

	l.prependItem(BindInt(&val))
}

func (l *boundIntList) Set(i int, v int) {
	if i > l.Length() {
		return
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	l.GetItem(i).(Int).Set(v)
}

// RuneList supports binding a list of rune values.
//
// Since: 2.0.0
type RuneList interface {
	DataList

	Append(rune)
	Get(int) rune
	Prepend(rune)
	Set(int, rune)
}

// NewRuneList returns a bindable list of rune values.
//
// Since: 2.0.0
func NewRuneList() RuneList {
	return &boundRuneList{}
}

// BindRuneList returns a bound list of rune values, based on the contents of the passed slice.
//
// Since: 2.0.0
func BindRuneList(v *[]rune) RuneList {
	if v == nil {
		return NewRuneList()
	}

	b := &boundRuneList{val: v}

	for _, i := range *v {
		b.appendItem(BindRune(&i))
	}

	return b
}

type boundRuneList struct {
	listBase

	val *[]rune
}

func (l *boundRuneList) Append(val rune) {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindRune(&val))
}

func (l *boundRuneList) Get(i int) rune {
	if i < 0 || i > l.Length() {
		return rune(0)
	}
	if l.val != nil {
		return (*l.val)[i]
	}

	return l.GetItem(i).(Rune).Get()
}

func (l *boundRuneList) Prepend(val rune) {
	if l.val != nil {
		*l.val = append([]rune{val}, *l.val...)
	}

	l.prependItem(BindRune(&val))
}

func (l *boundRuneList) Set(i int, v rune) {
	if i > l.Length() {
		return
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	l.GetItem(i).(Rune).Set(v)
}

// StringList supports binding a list of string values.
//
// Since: 2.0.0
type StringList interface {
	DataList

	Append(string)
	Get(int) string
	Prepend(string)
	Set(int, string)
}

// NewStringList returns a bindable list of string values.
//
// Since: 2.0.0
func NewStringList() StringList {
	return &boundStringList{}
}

// BindStringList returns a bound list of string values, based on the contents of the passed slice.
//
// Since: 2.0.0
func BindStringList(v *[]string) StringList {
	if v == nil {
		return NewStringList()
	}

	b := &boundStringList{val: v}

	for _, i := range *v {
		b.appendItem(BindString(&i))
	}

	return b
}

type boundStringList struct {
	listBase

	val *[]string
}

func (l *boundStringList) Append(val string) {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindString(&val))
}

func (l *boundStringList) Get(i int) string {
	if i < 0 || i > l.Length() {
		return ""
	}
	if l.val != nil {
		return (*l.val)[i]
	}

	return l.GetItem(i).(String).Get()
}

func (l *boundStringList) Prepend(val string) {
	if l.val != nil {
		*l.val = append([]string{val}, *l.val...)
	}

	l.prependItem(BindString(&val))
}

func (l *boundStringList) Set(i int, v string) {
	if i > l.Length() {
		return
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	l.GetItem(i).(String).Set(v)
}
