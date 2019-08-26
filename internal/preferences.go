package internal

import (
	"fyne.io/fyne"
)

// InMemoryPreferences provides an implementation of the fyne.Preferences API that is stored in memory.
type InMemoryPreferences struct {
	Values   map[string]interface{}
	OnChange func()
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*InMemoryPreferences)(nil)

func (p *InMemoryPreferences) fireChange() {
	if p.OnChange == nil {
		return
	}

	p.OnChange()
}

// Bool looks up a boolean value for the key
func (p *InMemoryPreferences) Bool(key string) bool {
	value, ok := p.Values[key]
	if !ok {
		return false
	}

	return value.(bool)
}

// SetBool saves a boolean value for the given key
func (p *InMemoryPreferences) SetBool(key string, value bool) {
	p.Values[key] = value
	p.fireChange()
}

// Float looks up a float64 value for the key
func (p *InMemoryPreferences) Float(key string) float64 {
	value, ok := p.Values[key]
	if !ok {
		return 0.0
	}

	return value.(float64)
}

// SetFloat saves a float64 value for the given key
func (p *InMemoryPreferences) SetFloat(key string, value float64) {
	p.Values[key] = value
	p.fireChange()
}

// Int looks up an integer value for the key
func (p *InMemoryPreferences) Int(key string) int {
	value, ok := p.Values[key]
	if !ok {
		return 0
	}

	// integers can be de-serialised as floats, so support both
	if intVal, ok := value.(int); ok {
		return intVal
	}
	return int(value.(float64))
}

// SetInt saves an integer value for the given key
func (p *InMemoryPreferences) SetInt(key string, value int) {
	p.Values[key] = value
	p.fireChange()
}

// String looks up a string value for the key
func (p *InMemoryPreferences) String(key string) string {
	value, ok := p.Values[key]
	if !ok {
		return ""
	}

	return value.(string)
}

// SetString saves a string value for the given key
func (p *InMemoryPreferences) SetString(key string, value string) {
	p.Values[key] = value
	p.fireChange()
}

// NewInMemoryPreferences creates a new preferences implementation stored in memory
func NewInMemoryPreferences() *InMemoryPreferences {
	p := &InMemoryPreferences{}
	p.Values = make(map[string]interface{})
	return p
}
