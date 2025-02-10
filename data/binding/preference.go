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

func bindPreferenceItem[T bool | float64 | int | string](key string, p fyne.Preferences, setLookup preferenceLookupSetter[T]) bindableItem[T] {
	if found, ok := lookupExistingBinding[T](key, p); ok {
		return found
	}

	listen := &prefBoundBase[T]{key: key, setLookup: setLookup}
	listen.replaceProvider(p)
	binds := prefBinds.ensurePreferencesAttached(p)
	binds.Store(key, listen)
	return listen
}

func lookupExistingBinding[T any](key string, p fyne.Preferences) (bindableItem[T], bool) {
	binds := prefBinds.getBindings(p)
	if binds == nil {
		return nil, false
	}

	if listen, ok := binds.Load(key); listen != nil && ok {
		if l, ok := listen.(bindableItem[T]); ok {
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
