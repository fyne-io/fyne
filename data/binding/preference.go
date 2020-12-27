// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import "fyne.io/fyne"

const keyTypeMismatchError = "A previous preference binding exists with different type for key: "

type prefBoundBool struct {
	base
	key string
	p   fyne.Preferences
}

// BindPreferenceBool returns a bindable bool value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0.0
func BindPreferenceBool(key string, p fyne.Preferences) Bool {
	if listen, ok := prefBinds[key]; ok {
		if l, ok := listen.(Bool); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBoundBool{key: key, p: p}
	ensurePreferencesAttached(p)
	prefBinds[key] = listen
	return listen
}

func (b *prefBoundBool) Get() (bool, error) {
	return b.p.Bool(b.key), nil
}

func (b *prefBoundBool) Set(v bool) error {
	b.p.SetBool(b.key, v)

	b.trigger()
	return nil
}

type prefBoundFloat struct {
	base
	key string
	p   fyne.Preferences
}

// BindPreferenceFloat returns a bindable float64 value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0.0
func BindPreferenceFloat(key string, p fyne.Preferences) Float {
	if listen, ok := prefBinds[key]; ok {
		if l, ok := listen.(Float); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBoundFloat{key: key, p: p}
	ensurePreferencesAttached(p)
	prefBinds[key] = listen
	return listen
}

func (b *prefBoundFloat) Get() (float64, error) {
	return b.p.Float(b.key), nil
}

func (b *prefBoundFloat) Set(v float64) error {
	b.p.SetFloat(b.key, v)

	b.trigger()
	return nil
}

type prefBoundInt struct {
	base
	key string
	p   fyne.Preferences
}

// BindPreferenceInt returns a bindable int value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0.0
func BindPreferenceInt(key string, p fyne.Preferences) Int {
	if listen, ok := prefBinds[key]; ok {
		if l, ok := listen.(Int); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBoundInt{key: key, p: p}
	ensurePreferencesAttached(p)
	prefBinds[key] = listen
	return listen
}

func (b *prefBoundInt) Get() (int, error) {
	return b.p.Int(b.key), nil
}

func (b *prefBoundInt) Set(v int) error {
	b.p.SetInt(b.key, v)

	b.trigger()
	return nil
}

type prefBoundString struct {
	base
	key string
	p   fyne.Preferences
}

// BindPreferenceString returns a bindable string value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0.0
func BindPreferenceString(key string, p fyne.Preferences) String {
	if listen, ok := prefBinds[key]; ok {
		if l, ok := listen.(String); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBoundString{key: key, p: p}
	ensurePreferencesAttached(p)
	prefBinds[key] = listen
	return listen
}

func (b *prefBoundString) Get() (string, error) {
	return b.p.String(b.key), nil
}

func (b *prefBoundString) Set(v string) error {
	b.p.SetString(b.key, v)

	b.trigger()
	return nil
}
