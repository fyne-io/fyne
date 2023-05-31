package fyne

// Preferences describes the ways that an app can save and load user preferences
type Preferences interface {
	// Bool looks up a bool value for the key
	Bool(key string) bool
	// BoolWithFallback looks up a bool value and returns the given fallback if not found
	BoolWithFallback(key string, fallback bool) bool
	// SetBool saves a bool value for the given key
	SetBool(key string, value bool)

	// BoolList looks up a list of bool values for the key
	//
	// Since: 2.4
	BoolList(key string) []bool
	// BoolListWithFallback looks up a list of bool values and returns the given fallback if not found
	//
	// Since: 2.4
	BoolListWithFallback(key string, fallback []bool) []bool
	// SetBoolList saves a list of bool values for the given key
	//
	// Since: 2.4
	SetBoolList(key string, value []bool)

	// Float looks up a float64 value for the key
	Float(key string) float64
	// FloatWithFallback looks up a float64 value and returns the given fallback if not found
	FloatWithFallback(key string, fallback float64) float64
	// SetFloat saves a float64 value for the given key
	SetFloat(key string, value float64)

	// FloatList looks up a list of float64 values for the key
	//
	// Since: 2.4
	FloatList(key string) []float64
	// FloatListWithFallback looks up a list of float64 values and returns the given fallback if not found
	//
	// Since: 2.4
	FloatListWithFallback(key string, fallback []float64) []float64
	// SetFloatList saves a list of float64 values for the given key
	//
	// Since: 2.4
	SetFloatList(key string, value []float64)

	// Int looks up an integer value for the key
	Int(key string) int
	// IntWithFallback looks up an integer value and returns the given fallback if not found
	IntWithFallback(key string, fallback int) int
	// SetInt saves an integer value for the given key
	SetInt(key string, value int)

	// IntList looks up a list of int values for the key
	//
	// Since: 2.4
	IntList(key string) []int
	// IntListWithFallback looks up a list of int values and returns the given fallback if not found
	//
	// Since: 2.4
	IntListWithFallback(key string, fallback []int) []int
	// SetIntList saves a list of string values for the given key
	//
	// Since: 2.4
	SetIntList(key string, value []int)

	// String looks up a string value for the key
	String(key string) string
	// StringWithFallback looks up a string value and returns the given fallback if not found
	StringWithFallback(key, fallback string) string
	// SetString saves a string value for the given key
	SetString(key string, value string)

	// StringList looks up a list of string values for the key
	//
	// Since: 2.4
	StringList(key string) []string
	// StringListWithFallback looks up a list of string values and returns the given fallback if not found
	//
	// Since: 2.4
	StringListWithFallback(key string, fallback []string) []string
	// SetStringList saves a list of string values for the given key
	//
	// Since: 2.4
	SetStringList(key string, value []string)

	// RemoveValue removes a value for the given key (not currently supported on iOS)
	RemoveValue(key string)

	// AddChangeListener allows code to be notified when some preferences change. This will fire on any update.
	AddChangeListener(func())

	// ChangeListeners returns a list of the known change listeners for this preference set.
	//
	// Since: 2.3
	ChangeListeners() []func()
}
