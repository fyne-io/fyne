// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import "fyne.io/fyne/v2"

const keyTypeMismatchError = "A previous preference binding exists with different type for key: "

type prefBoundBool struct {
	base
	key   string
	p     fyne.Preferences
	cache bool
}

// BindPreferenceBool returns a bindable bool value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceBool(key string, p fyne.Preferences) Bool {
	if prefBinds[p] != nil {
		if listen, ok := prefBinds[p][key]; ok {
			if l, ok := listen.(Bool); ok {
				return l
			}
			fyne.LogError(keyTypeMismatchError+key, nil)
		}
	}

	listen := &prefBoundBool{key: key, p: p}
	ensurePreferencesAttached(p)
	prefBinds[p][key] = listen
	return listen
}

func (b *prefBoundBool) Get() (bool, error) {
	b.cache = b.p.Bool(b.key)
	return b.cache, nil
}

func (b *prefBoundBool) Set(v bool) error {
	b.p.SetBool(b.key, v)

	b.lock.RLock()
	defer b.lock.RUnlock()
	b.trigger()
	return nil
}

func (b *prefBoundBool) checkForChange() {
	if b.p.Bool(b.key) == b.cache {
		return
	}

	b.trigger()
}

type prefBoundFloat struct {
	base
	key   string
	p     fyne.Preferences
	cache float64
}

// BindPreferenceFloat returns a bindable float64 value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceFloat(key string, p fyne.Preferences) Float {
	if prefBinds[p] != nil {
		if listen, ok := prefBinds[p][key]; ok {
			if l, ok := listen.(Float); ok {
				return l
			}
			fyne.LogError(keyTypeMismatchError+key, nil)
		}
	}

	listen := &prefBoundFloat{key: key, p: p}
	ensurePreferencesAttached(p)
	prefBinds[p][key] = listen
	return listen
}

func (b *prefBoundFloat) Get() (float64, error) {
	b.cache = b.p.Float(b.key)
	return b.cache, nil
}

func (b *prefBoundFloat) Set(v float64) error {
	b.p.SetFloat(b.key, v)

	b.lock.RLock()
	defer b.lock.RUnlock()
	b.trigger()
	return nil
}

func (b *prefBoundFloat) checkForChange() {
	if b.p.Float(b.key) == b.cache {
		return
	}

	b.trigger()
}

type prefBoundInt struct {
	base
	key   string
	p     fyne.Preferences
	cache int
}

// BindPreferenceInt returns a bindable int value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceInt(key string, p fyne.Preferences) Int {
	if prefBinds[p] != nil {
		if listen, ok := prefBinds[p][key]; ok {
			if l, ok := listen.(Int); ok {
				return l
			}
			fyne.LogError(keyTypeMismatchError+key, nil)
		}
	}

	listen := &prefBoundInt{key: key, p: p}
	ensurePreferencesAttached(p)
	prefBinds[p][key] = listen
	return listen
}

func (b *prefBoundInt) Get() (int, error) {
	b.cache = b.p.Int(b.key)
	return b.cache, nil
}

func (b *prefBoundInt) Set(v int) error {
	b.p.SetInt(b.key, v)

	b.lock.RLock()
	defer b.lock.RUnlock()
	b.trigger()
	return nil
}

func (b *prefBoundInt) checkForChange() {
	if b.p.Int(b.key) == b.cache {
		return
	}

	b.trigger()
}

type prefBoundString struct {
	base
	key   string
	p     fyne.Preferences
	cache string
}

// BindPreferenceString returns a bindable string value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceString(key string, p fyne.Preferences) String {
	if prefBinds[p] != nil {
		if listen, ok := prefBinds[p][key]; ok {
			if l, ok := listen.(String); ok {
				return l
			}
			fyne.LogError(keyTypeMismatchError+key, nil)
		}
	}

	listen := &prefBoundString{key: key, p: p}
	ensurePreferencesAttached(p)
	prefBinds[p][key] = listen
	return listen
}

func (b *prefBoundString) Get() (string, error) {
	b.cache = b.p.String(b.key)
	return b.cache, nil
}

func (b *prefBoundString) Set(v string) error {
	b.p.SetString(b.key, v)

	b.lock.RLock()
	defer b.lock.RUnlock()
	b.trigger()
	return nil
}

func (b *prefBoundString) checkForChange() {
	if b.p.String(b.key) == b.cache {
		return
	}

	b.trigger()
}
