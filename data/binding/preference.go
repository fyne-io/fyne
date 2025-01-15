package binding

import "fyne.io/fyne/v2"

const keyTypeMismatchError = "A previous preference binding exists with different type for key: "

type prefBoundBool struct {
	prefBoundBase[bool]
}

// BindPreferenceBool returns a bindable bool value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceBool(key string, p fyne.Preferences) Bool {
	if found, ok := lookupExistingBinding[bool](key, p); ok {
		return found
	}

	listen := &prefBoundBool{}
	listen.setKey(key)
	listen.replaceProvider(p)
	binds := prefBinds.ensurePreferencesAttached(p)
	binds.Store(key, listen)
	return listen
}

func (b *prefBoundBool) replaceProvider(p fyne.Preferences) {
	b.get = p.Bool
	b.set = p.SetBool
}

type prefBoundFloat struct {
	prefBoundBase[float64]
}

// BindPreferenceFloat returns a bindable float64 value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceFloat(key string, p fyne.Preferences) Float {
	if found, ok := lookupExistingBinding[float64](key, p); ok {
		return found
	}

	listen := &prefBoundFloat{}
	listen.setKey(key)
	listen.replaceProvider(p)
	binds := prefBinds.ensurePreferencesAttached(p)
	binds.Store(key, listen)
	return listen
}

func (b *prefBoundFloat) replaceProvider(p fyne.Preferences) {
	b.get = p.Float
	b.set = p.SetFloat
}

type prefBoundInt struct {
	prefBoundBase[int]
}

// BindPreferenceInt returns a bindable int value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceInt(key string, p fyne.Preferences) Int {
	if found, ok := lookupExistingBinding[int](key, p); ok {
		return found
	}

	listen := &prefBoundInt{}
	listen.setKey(key)
	listen.replaceProvider(p)
	binds := prefBinds.ensurePreferencesAttached(p)
	binds.Store(key, listen)
	return listen
}

func (b *prefBoundInt) replaceProvider(p fyne.Preferences) {
	b.get = p.Int
	b.set = p.SetInt
}

type prefBoundString struct {
	prefBoundBase[string]
}

// BindPreferenceString returns a bindable string value that is managed by the application preferences.
// Changes to this value will be saved to application storage and when the app starts the previous values will be read.
//
// Since: 2.0
func BindPreferenceString(key string, p fyne.Preferences) String {
	if found, ok := lookupExistingBinding[string](key, p); ok {
		return found
	}

	listen := &prefBoundString{}
	listen.setKey(key)
	listen.replaceProvider(p)
	binds := prefBinds.ensurePreferencesAttached(p)
	binds.Store(key, listen)
	return listen
}

func (b *prefBoundString) replaceProvider(p fyne.Preferences) {
	b.get = p.String
	b.set = p.SetString
}
