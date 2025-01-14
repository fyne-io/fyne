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
	if found, ok := lookupExistingBinding(key, p); ok {
		if l, ok := found.(Bool); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBoundBool{}
	setupPrefItem(listen, key, p)
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
	if found, ok := lookupExistingBinding(key, p); ok {
		if l, ok := found.(Float); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBoundFloat{}
	setupPrefItem(listen, key, p)
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
	if found, ok := lookupExistingBinding(key, p); ok {
		if l, ok := found.(Int); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBoundInt{}
	setupPrefItem(listen, key, p)
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
	if found, ok := lookupExistingBinding(key, p); ok {
		if l, ok := found.(String); ok {
			return l
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	listen := &prefBoundString{}
	setupPrefItem(listen, key, p)
	return listen
}

func (b *prefBoundString) replaceProvider(p fyne.Preferences) {
	b.get = p.String
	b.set = p.SetString
}
