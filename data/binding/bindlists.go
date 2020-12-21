// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// BoolList supports binding a list of bool values.
//
// Since: 2.0.0
type BoolList interface {
	DataList

	Append(bool) error
	Get(int) (bool, error)
	Prepend(bool) error
	Set(int, bool) error
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

	for i := range *v {
		b.appendItem(BindBool(&((*v)[i])))
	}

	return b
}

type boundBoolList struct {
	listBase

	val *[]bool
}

func (l *boundBoolList) Append(val bool) error {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindBool(&val))
	return nil
}

func (l *boundBoolList) Get(i int) (bool, error) {
	if i < 0 || i >= l.Length() {
		return false, outOfBounds
	}
	if l.val != nil {
		return (*l.val)[i], nil
	}

	item, ok := l.GetItem(i)
	if !ok {
		return false, outOfBounds
	}
	return item.(Bool).Get()
}

func (l *boundBoolList) Prepend(val bool) error {
	if l.val != nil {
		*l.val = append([]bool{val}, *l.val...)
	}

	l.prependItem(BindBool(&val))
	return nil
}

func (l *boundBoolList) Set(i int, v bool) error {
	if i < 0 || i >= l.Length() {
		return outOfBounds
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	item, ok := l.GetItem(i)
	if !ok {
		return outOfBounds
	}
	return item.(Bool).Set(v)
}

// FloatList supports binding a list of float64 values.
//
// Since: 2.0.0
type FloatList interface {
	DataList

	Append(float64) error
	Get(int) (float64, error)
	Prepend(float64) error
	Set(int, float64) error
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

func (l *boundFloatList) Append(val float64) error {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindFloat(&val))
	return nil
}

func (l *boundFloatList) Get(i int) (float64, error) {
	if i < 0 || i >= l.Length() {
		return 0.0, outOfBounds
	}
	if l.val != nil {
		return (*l.val)[i], nil
	}

	item, ok := l.GetItem(i)
	if !ok {
		return 0.0, outOfBounds
	}
	return item.(Float).Get()
}

func (l *boundFloatList) Prepend(val float64) error {
	if l.val != nil {
		*l.val = append([]float64{val}, *l.val...)
	}

	l.prependItem(BindFloat(&val))
	return nil
}

func (l *boundFloatList) Set(i int, v float64) error {
	if i < 0 || i >= l.Length() {
		return outOfBounds
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	item, ok := l.GetItem(i)
	if !ok {
		return outOfBounds
	}
	return item.(Float).Set(v)
}

// IntList supports binding a list of int values.
//
// Since: 2.0.0
type IntList interface {
	DataList

	Append(int) error
	Get(int) (int, error)
	Prepend(int) error
	Set(int, int) error
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

	for i := range *v {
		b.appendItem(BindInt(&((*v)[i])))
	}

	return b
}

type boundIntList struct {
	listBase

	val *[]int
}

func (l *boundIntList) Append(val int) error {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindInt(&val))
	return nil
}

func (l *boundIntList) Get(i int) (int, error) {
	if i < 0 || i >= l.Length() {
		return 0, outOfBounds
	}
	if l.val != nil {
		return (*l.val)[i], nil
	}

	item, ok := l.GetItem(i)
	if !ok {
		return 0, outOfBounds
	}
	return item.(Int).Get()
}

func (l *boundIntList) Prepend(val int) error {
	if l.val != nil {
		*l.val = append([]int{val}, *l.val...)
	}

	l.prependItem(BindInt(&val))
	return nil
}

func (l *boundIntList) Set(i int, v int) error {
	if i < 0 || i >= l.Length() {
		return outOfBounds
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	item, ok := l.GetItem(i)
	if !ok {
		return outOfBounds
	}
	return item.(Int).Set(v)
}

// RuneList supports binding a list of rune values.
//
// Since: 2.0.0
type RuneList interface {
	DataList

	Append(rune) error
	Get(int) (rune, error)
	Prepend(rune) error
	Set(int, rune) error
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

	for i := range *v {
		b.appendItem(BindRune(&((*v)[i])))
	}

	return b
}

type boundRuneList struct {
	listBase

	val *[]rune
}

func (l *boundRuneList) Append(val rune) error {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindRune(&val))
	return nil
}

func (l *boundRuneList) Get(i int) (rune, error) {
	if i < 0 || i >= l.Length() {
		return rune(0), outOfBounds
	}
	if l.val != nil {
		return (*l.val)[i], nil
	}

	item, ok := l.GetItem(i)
	if !ok {
		return rune(0), outOfBounds
	}
	return item.(Rune).Get()
}

func (l *boundRuneList) Prepend(val rune) error {
	if l.val != nil {
		*l.val = append([]rune{val}, *l.val...)
	}

	l.prependItem(BindRune(&val))
	return nil
}

func (l *boundRuneList) Set(i int, v rune) error {
	if i < 0 || i >= l.Length() {
		return outOfBounds
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	item, ok := l.GetItem(i)
	if !ok {
		return outOfBounds
	}
	return item.(Rune).Set(v)
}

// StringList supports binding a list of string values.
//
// Since: 2.0.0
type StringList interface {
	DataList

	Append(string) error
	Get(int) (string, error)
	Prepend(string) error
	Set(int, string) error
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

	for i := range *v {
		b.appendItem(BindString(&((*v)[i])))
	}

	return b
}

type boundStringList struct {
	listBase

	val *[]string
}

func (l *boundStringList) Append(val string) error {
	if l.val != nil {
		*l.val = append(*l.val, val)
	}

	l.appendItem(BindString(&val))
	return nil
}

func (l *boundStringList) Get(i int) (string, error) {
	if i < 0 || i >= l.Length() {
		return "", outOfBounds
	}
	if l.val != nil {
		return (*l.val)[i], nil
	}

	item, ok := l.GetItem(i)
	if !ok {
		return "", outOfBounds
	}
	return item.(String).Get()
}

func (l *boundStringList) Prepend(val string) error {
	if l.val != nil {
		*l.val = append([]string{val}, *l.val...)
	}

	l.prependItem(BindString(&val))
	return nil
}

func (l *boundStringList) Set(i int, v string) error {
	if i < 0 || i >= l.Length() {
		return outOfBounds
	}
	if l.val != nil {
		(*l.val)[i] = v
	}

	item, ok := l.GetItem(i)
	if !ok {
		return outOfBounds
	}
	return item.(String).Set(v)
}
