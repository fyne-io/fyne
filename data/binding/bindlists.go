// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

// BoolList supports binding a list of bool values.
//
// Since: 2.0.0
type BoolList interface {
	DataList

	Append(bool) error
	Get() ([]bool, error)
	GetValue(int) (bool, error)
	Prepend(bool) error
	Set([]bool) error
	SetValue(int, bool) error
}

// ExternalBoolList supports binding a list of bool values from an external variable.
//
// Since: 2.0.0
type ExternalBoolList interface {
	BoolList

	Reload() error
}

// NewBoolList returns a bindable list of bool values.
//
// Since: 2.0.0
func NewBoolList() BoolList {
	return &boundBoolList{val: &[]bool{}}
}

// BindBoolList returns a bound list of bool values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0.0
func BindBoolList(v *[]bool) ExternalBoolList {
	if v == nil {
		return NewBoolList().(ExternalBoolList)
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
	*l.val = append(*l.val, val)

	l.appendItem(BindBool(&val))
	return nil
}

func (l *boundBoolList) Get() ([]bool, error) {
	return *l.val, nil
}

func (l *boundBoolList) GetValue(i int) (bool, error) {
	if i < 0 || i >= l.Length() {
		return false, errOutOfBounds
	}
	return (*l.val)[i], nil
}

func (l *boundBoolList) Prepend(val bool) error {
	*l.val = append([]bool{val}, *l.val...)

	l.prependItem(BindBool(&val))
	return nil
}

func (l *boundBoolList) Set(v []bool) (retErr error) {
	oldLen := len(l.items)
	*l.val = v
	newLen := len(v)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(BindBool(&((*l.val)[i])))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		old, err := l.items[i].(Bool).Get()
		val := (*(l.val))[i]
		if err != nil || (*(l.val))[i] != old {
			err = item.(*boundBool).Set(val)
			if err != nil {
				retErr = err
			}
		}
	}
	return
}

func (l *boundBoolList) SetValue(i int, v bool) error {
	if i < 0 || i >= l.Length() {
		return errOutOfBounds
	}
	(*l.val)[i] = v

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Bool).Set(v)
}

// FloatList supports binding a list of float64 values.
//
// Since: 2.0.0
type FloatList interface {
	DataList

	Append(float64) error
	Get() ([]float64, error)
	GetValue(int) (float64, error)
	Prepend(float64) error
	Set([]float64) error
	SetValue(int, float64) error
}

// ExternalFloatList supports binding a list of float64 values from an external variable.
//
// Since: 2.0.0
type ExternalFloatList interface {
	FloatList

	Reload() error
}

// NewFloatList returns a bindable list of float64 values.
//
// Since: 2.0.0
func NewFloatList() FloatList {
	return &boundFloatList{val: &[]float64{}}
}

// BindFloatList returns a bound list of float64 values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0.0
func BindFloatList(v *[]float64) ExternalFloatList {
	if v == nil {
		return NewFloatList().(ExternalFloatList)
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
	*l.val = append(*l.val, val)

	l.appendItem(BindFloat(&val))
	return nil
}

func (l *boundFloatList) Get() ([]float64, error) {
	return *l.val, nil
}

func (l *boundFloatList) GetValue(i int) (float64, error) {
	if i < 0 || i >= l.Length() {
		return 0.0, errOutOfBounds
	}
	return (*l.val)[i], nil
}

func (l *boundFloatList) Prepend(val float64) error {
	*l.val = append([]float64{val}, *l.val...)

	l.prependItem(BindFloat(&val))
	return nil
}

func (l *boundFloatList) Set(v []float64) (retErr error) {
	oldLen := len(l.items)
	*l.val = v
	newLen := len(v)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(BindFloat(&((*l.val)[i])))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		old, err := l.items[i].(Float).Get()
		val := (*(l.val))[i]
		if err != nil || (*(l.val))[i] != old {
			err = item.(*boundFloat).Set(val)
			if err != nil {
				retErr = err
			}
		}
	}
	return
}

func (l *boundFloatList) SetValue(i int, v float64) error {
	if i < 0 || i >= l.Length() {
		return errOutOfBounds
	}
	(*l.val)[i] = v

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Float).Set(v)
}

// IntList supports binding a list of int values.
//
// Since: 2.0.0
type IntList interface {
	DataList

	Append(int) error
	Get() ([]int, error)
	GetValue(int) (int, error)
	Prepend(int) error
	Set([]int) error
	SetValue(int, int) error
}

// ExternalIntList supports binding a list of int values from an external variable.
//
// Since: 2.0.0
type ExternalIntList interface {
	IntList

	Reload() error
}

// NewIntList returns a bindable list of int values.
//
// Since: 2.0.0
func NewIntList() IntList {
	return &boundIntList{val: &[]int{}}
}

// BindIntList returns a bound list of int values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0.0
func BindIntList(v *[]int) ExternalIntList {
	if v == nil {
		return NewIntList().(ExternalIntList)
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
	*l.val = append(*l.val, val)

	l.appendItem(BindInt(&val))
	return nil
}

func (l *boundIntList) Get() ([]int, error) {
	return *l.val, nil
}

func (l *boundIntList) GetValue(i int) (int, error) {
	if i < 0 || i >= l.Length() {
		return 0, errOutOfBounds
	}
	return (*l.val)[i], nil
}

func (l *boundIntList) Prepend(val int) error {
	*l.val = append([]int{val}, *l.val...)

	l.prependItem(BindInt(&val))
	return nil
}

func (l *boundIntList) Set(v []int) (retErr error) {
	oldLen := len(l.items)
	*l.val = v
	newLen := len(v)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(BindInt(&((*l.val)[i])))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		old, err := l.items[i].(Int).Get()
		val := (*(l.val))[i]
		if err != nil || (*(l.val))[i] != old {
			err = item.(*boundInt).Set(val)
			if err != nil {
				retErr = err
			}
		}
	}
	return
}

func (l *boundIntList) SetValue(i int, v int) error {
	if i < 0 || i >= l.Length() {
		return errOutOfBounds
	}
	(*l.val)[i] = v

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Int).Set(v)
}

// RuneList supports binding a list of rune values.
//
// Since: 2.0.0
type RuneList interface {
	DataList

	Append(rune) error
	Get() ([]rune, error)
	GetValue(int) (rune, error)
	Prepend(rune) error
	Set([]rune) error
	SetValue(int, rune) error
}

// ExternalRuneList supports binding a list of rune values from an external variable.
//
// Since: 2.0.0
type ExternalRuneList interface {
	RuneList

	Reload() error
}

// NewRuneList returns a bindable list of rune values.
//
// Since: 2.0.0
func NewRuneList() RuneList {
	return &boundRuneList{val: &[]rune{}}
}

// BindRuneList returns a bound list of rune values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0.0
func BindRuneList(v *[]rune) ExternalRuneList {
	if v == nil {
		return NewRuneList().(ExternalRuneList)
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
	*l.val = append(*l.val, val)

	l.appendItem(BindRune(&val))
	return nil
}

func (l *boundRuneList) Get() ([]rune, error) {
	return *l.val, nil
}

func (l *boundRuneList) GetValue(i int) (rune, error) {
	if i < 0 || i >= l.Length() {
		return rune(0), errOutOfBounds
	}
	return (*l.val)[i], nil
}

func (l *boundRuneList) Prepend(val rune) error {
	*l.val = append([]rune{val}, *l.val...)

	l.prependItem(BindRune(&val))
	return nil
}

func (l *boundRuneList) Set(v []rune) (retErr error) {
	oldLen := len(l.items)
	*l.val = v
	newLen := len(v)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(BindRune(&((*l.val)[i])))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		old, err := l.items[i].(Rune).Get()
		val := (*(l.val))[i]
		if err != nil || (*(l.val))[i] != old {
			err = item.(*boundRune).Set(val)
			if err != nil {
				retErr = err
			}
		}
	}
	return
}

func (l *boundRuneList) SetValue(i int, v rune) error {
	if i < 0 || i >= l.Length() {
		return errOutOfBounds
	}
	(*l.val)[i] = v

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(Rune).Set(v)
}

// StringList supports binding a list of string values.
//
// Since: 2.0.0
type StringList interface {
	DataList

	Append(string) error
	Get() ([]string, error)
	GetValue(int) (string, error)
	Prepend(string) error
	Set([]string) error
	SetValue(int, string) error
}

// ExternalStringList supports binding a list of string values from an external variable.
//
// Since: 2.0.0
type ExternalStringList interface {
	StringList

	Reload() error
}

// NewStringList returns a bindable list of string values.
//
// Since: 2.0.0
func NewStringList() StringList {
	return &boundStringList{val: &[]string{}}
}

// BindStringList returns a bound list of string values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0.0
func BindStringList(v *[]string) ExternalStringList {
	if v == nil {
		return NewStringList().(ExternalStringList)
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
	*l.val = append(*l.val, val)

	l.appendItem(BindString(&val))
	return nil
}

func (l *boundStringList) Get() ([]string, error) {
	return *l.val, nil
}

func (l *boundStringList) GetValue(i int) (string, error) {
	if i < 0 || i >= l.Length() {
		return "", errOutOfBounds
	}
	return (*l.val)[i], nil
}

func (l *boundStringList) Prepend(val string) error {
	*l.val = append([]string{val}, *l.val...)

	l.prependItem(BindString(&val))
	return nil
}

func (l *boundStringList) Set(v []string) (retErr error) {
	oldLen := len(l.items)
	*l.val = v
	newLen := len(v)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		l.trigger()
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(BindString(&((*l.val)[i])))
		}
		l.trigger()
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		old, err := l.items[i].(String).Get()
		val := (*(l.val))[i]
		if err != nil || (*(l.val))[i] != old {
			err = item.(*boundString).Set(val)
			if err != nil {
				retErr = err
			}
		}
	}
	return
}

func (l *boundStringList) SetValue(i int, v string) error {
	if i < 0 || i >= l.Length() {
		return errOutOfBounds
	}
	(*l.val)[i] = v

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(String).Set(v)
}
