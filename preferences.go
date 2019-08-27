package fyne

// Preferences describes the ways that an app can save and load user preferences
type Preferences interface {
	// Bool looks up a boolean value for the key
	Bool(key string) bool
	// BoolWithFallback looks up a boolean value and returns the given fallback if not found
	BoolWithFallback(key string, fallback bool) bool
	// SetBool saves a boolean value for the given key
	SetBool(key string, value bool)

	// Float looks up a float64 value for the key
	Float(key string) float64
	// FloatWithFallback looks up a float64 value and returns the given fallback if not found
	FloatWithFallback(key string, fallback float64) float64
	// SetFloat saves a float64 value for the given key
	SetFloat(key string, value float64)

	// Int looks up an integer value for the key
	Int(key string) int
	// IntWithFallback looks up an integer value and returns the given fallback if not found
	IntWithFallback(key string, fallback int) int
	// SetInt saves an integer value for the given key
	SetInt(key string, value int)

	// String looks up a string value for the key
	String(key string) string
	// StringWithFallback looks up a string value and returns the given fallback if not found
	StringWithFallback(key, fallback string) string
	// SetString saves a string value for the given key
	SetString(key string, value string)
}
