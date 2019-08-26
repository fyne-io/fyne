package fyne

// Preferences describes the ways that an app can save and load user preferences
type Preferences interface {
	// Bool looks up a boolean value for the key
	Bool(key string) bool
	// SetBool saves a boolean value for the given key
	SetBool(key string, value bool)

	// Float looks up a float64 value for the key
	Float(key string) float64
	// SetFloat saves a float64 value for the given key
	SetFloat(key string, value float64)

	// Int looks up an integer value for the key
	Int(key string) int
	// SetInt saves an integer value for the given key
	SetInt(key string, value int)

	// String looks up a string value for the key
	String(key string) string
	// SetString saves a string value for the given key
	SetString(key string, value string)
}
