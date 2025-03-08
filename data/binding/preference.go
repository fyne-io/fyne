package binding

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
)

// Work around Go not supporting generic methods on non-generic types:
type preferenceLookupSetter[T any] func(fyne.Preferences) (func(string) T, func(string, T))

const keyTypeMismatchError = "A previous preference binding exists with different type for key: "

// BindPreferenceBool returns a bindable bool value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceBool(key string, p fyne.Preferences) Bool {
	return bindPreferenceItem(key, p,
		func(p fyne.Preferences) (func(string) bool, func(string, bool)) {
			return p.Bool, p.SetBool
		})
}

// BindPreferenceBoolList returns a bound list of bool values that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.6
func BindPreferenceBoolList(key string, p fyne.Preferences) BoolList {
	return bindPreferenceListComparable[bool](key, p,
		func(p fyne.Preferences) (func(string) []bool, func(string, []bool)) {
			return p.BoolList, p.SetBoolList
		},
	)
}

// BindPreferenceFloat returns a bindable float64 value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceFloat(key string, p fyne.Preferences) Float {
	return bindPreferenceItem(key, p,
		func(p fyne.Preferences) (func(string) float64, func(string, float64)) {
			return p.Float, p.SetFloat
		})
}

// BindPreferenceFloatList returns a bound list of float64 values that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.6
func BindPreferenceFloatList(key string, p fyne.Preferences) FloatList {
	return bindPreferenceListComparable[float64](key, p,
		func(p fyne.Preferences) (func(string) []float64, func(string, []float64)) {
			return p.FloatList, p.SetFloatList
		},
	)
}

// BindPreferenceInt returns a bindable int value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceInt(key string, p fyne.Preferences) Int {
	return bindPreferenceItem(key, p,
		func(p fyne.Preferences) (func(string) int, func(string, int)) {
			return p.Int, p.SetInt
		})
}

// BindPreferenceIntList returns a bound list of int values that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.6
func BindPreferenceIntList(key string, p fyne.Preferences) IntList {
	return bindPreferenceListComparable[int](key, p,
		func(p fyne.Preferences) (func(string) []int, func(string, []int)) {
			return p.IntList, p.SetIntList
		},
	)
}

// BindPreferenceString returns a bindable string value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceString(key string, p fyne.Preferences) String {
	return bindPreferenceItem(key, p,
		func(p fyne.Preferences) (func(string) string, func(string, string)) {
			return p.String, p.SetString
		})
}

// BindPreferenceStringList returns a bound list of string values that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.6
func BindPreferenceStringList(key string, p fyne.Preferences) StringList {
	return bindPreferenceListComparable[string](key, p,
		func(p fyne.Preferences) (func(string) []string, func(string, []string)) {
			return p.StringList, p.SetStringList
		},
	)
}

func bindPreferenceItem[T bool | float64 | int | string](key string, p fyne.Preferences, setLookup preferenceLookupSetter[T]) Item[T] {
	if found, ok := lookupExistingBinding[T](key, p); ok {
		return found
	}

	listen := &prefBoundBase[T]{key: key, setLookup: setLookup}
	listen.replaceProvider(p)
	binds := prefBinds.ensurePreferencesAttached(p)
	binds.Store(key, listen)
	return listen
}

func lookupExistingBinding[T any](key string, p fyne.Preferences) (Item[T], bool) {
	binds := prefBinds.getBindings(p)
	if binds == nil {
		return nil, false
	}

	if listen, ok := binds.Load(key); listen != nil && ok {
		if l, ok := listen.(Item[T]); ok {
			return l, ok
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	return nil, false
}

func lookupExistingListBinding[T bool | float64 | int | string](key string, p fyne.Preferences) (*prefBoundList[T], bool) {
	binds := prefBinds.getBindings(p)
	if binds == nil {
		return nil, false
	}

	if listen, ok := binds.Load(key); listen != nil && ok {
		if l, ok := listen.(*prefBoundList[T]); ok {
			return l, ok
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	return nil, false
}

type prefBoundBase[T bool | float64 | int | string] struct {
	base
	key string

	get       func(string) T
	set       func(string, T)
	setLookup preferenceLookupSetter[T]
	cache     atomic.Pointer[T]
}

func (b *prefBoundBase[T]) Get() (T, error) {
	cache := b.get(b.key)
	b.cache.Store(&cache)
	return cache, nil
}

func (b *prefBoundBase[T]) Set(v T) error {
	b.set(b.key, v)

	b.lock.RLock()
	defer b.lock.RUnlock()
	b.trigger()
	return nil
}

func (b *prefBoundBase[T]) checkForChange() {
	val := b.cache.Load()
	if val != nil && b.get(b.key) == *val {
		return
	}
	b.trigger()
}

func (b *prefBoundBase[T]) replaceProvider(p fyne.Preferences) {
	b.get, b.set = b.setLookup(p)
}

type prefBoundList[T bool | float64 | int | string] struct {
	boundList[T]
	key string

	get       func(string) []T
	set       func(string, []T)
	setLookup preferenceLookupSetter[[]T]
}

func (b *prefBoundList[T]) checkForChange() {
	val := *b.val
	updated := b.get(b.key)
	if val == nil || len(updated) != len(val) {
		b.Set(updated)
		return
	}
	if val == nil {
		return
	}

	// incoming changes to a preference list are not at the child level
	for i, v := range val {
		if i >= len(updated) {
			break
		}

		if !b.comparator(v, updated[i]) {
			_ = b.items[i].(Item[T]).Set(updated[i])
		}
	}
}

func (b *prefBoundList[T]) replaceProvider(p fyne.Preferences) {
	b.get, b.set = b.setLookup(p)
}

type internalPrefs = interface{ WriteValues(func(map[string]any)) }

func bindPreferenceListComparable[T bool | float64 | int | string](key string, p fyne.Preferences,
	setLookup preferenceLookupSetter[[]T]) *prefBoundList[T] {
	if found, ok := lookupExistingListBinding[T](key, p); ok {
		return found
	}

	listen := &prefBoundList[T]{key: key, setLookup: setLookup}
	listen.replaceProvider(p)

	items := listen.get(listen.key)
	listen.boundList = *bindList(nil, func(t1, t2 T) bool { return t1 == t2 })

	listen.boundList.AddListener(NewDataListener(func() {
		cached := *listen.val
		replaced := listen.get(listen.key)
		if len(cached) == len(replaced) {
			return
		}

		listen.set(listen.key, *listen.val)
		listen.trigger()
	}))

	listen.boundList.parentListener = func(index int) {
		listen.set(listen.key, *listen.val)

		// the child changes are not seen on the write end so force it
		if prefs, ok := p.(internalPrefs); ok {
			prefs.WriteValues(func(map[string]any) {})
		}
	}
	listen.boundList.Set(items)

	binds := prefBinds.ensurePreferencesAttached(p)
	binds.Store(key, listen)
	return listen
}
