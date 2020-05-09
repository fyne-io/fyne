package internal

import (
	"sync"

	"fyne.io/fyne"
)

// InMemoryPreferences provides an implementation of the fyne.Preferences API that is stored in memory.
type InMemoryPreferences struct {
	values   map[string]interface{}
	lock     sync.RWMutex
	OnChange func()
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*InMemoryPreferences)(nil)

func (p *InMemoryPreferences) set(key string, value interface{}) {
	p.lock.Lock()

	p.values[key] = value
	p.lock.Unlock()

	p.fireChange()
}

func (p *InMemoryPreferences) get(key string) (interface{}, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	v, err := p.values[key]
	return v, err
}

func (p *InMemoryPreferences) remove(key string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.values, key)
}

// Values provides access to the underlying value map - for internal use only...
func (p *InMemoryPreferences) Values() map[string]interface{} {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.values
}

func (p *InMemoryPreferences) fireChange() {
	if p.OnChange == nil {
		return
	}

	p.OnChange()
}

// Bool looks up a boolean value for the key
func (p *InMemoryPreferences) Bool(key string) bool {
	return p.BoolWithFallback(key, false)
}

// BoolWithFallback looks up a boolean value and returns the given fallback if not found
func (p *InMemoryPreferences) BoolWithFallback(key string, fallback bool) bool {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}

	val, ok := value.(bool)
	if !ok {
		return false
	}
	return val
}

// SetBool saves a boolean value for the given key
func (p *InMemoryPreferences) SetBool(key string, value bool) {
	p.set(key, value)
}

// Float looks up a float64 value for the key
func (p *InMemoryPreferences) Float(key string) float64 {
	return p.FloatWithFallback(key, 0.0)
}

// FloatWithFallback looks up a float64 value and returns the given fallback if not found
func (p *InMemoryPreferences) FloatWithFallback(key string, fallback float64) float64 {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}

	val, ok := value.(float64)
	if !ok {
		return 0.0
	}
	return val
}

// SetFloat saves a float64 value for the given key
func (p *InMemoryPreferences) SetFloat(key string, value float64) {
	p.set(key, value)
}

// Int looks up an integer value for the key
func (p *InMemoryPreferences) Int(key string) int {
	return p.IntWithFallback(key, 0)
}

// IntWithFallback looks up an integer value and returns the given fallback if not found
func (p *InMemoryPreferences) IntWithFallback(key string, fallback int) int {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}

	// integers can be de-serialised as floats, so support both
	if intVal, ok := value.(int); ok {
		return intVal
	}
	val, ok := value.(float64)
	if !ok {
		return 0
	}
	return int(val)
}

// SetInt saves an integer value for the given key
func (p *InMemoryPreferences) SetInt(key string, value int) {
	p.set(key, value)
}

// String looks up a string value for the key
func (p *InMemoryPreferences) String(key string) string {
	return p.StringWithFallback(key, "")
}

// StringWithFallback looks up a string value and returns the given fallback if not found
func (p *InMemoryPreferences) StringWithFallback(key, fallback string) string {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}

	return value.(string)
}

// SetString saves a string value for the given key
func (p *InMemoryPreferences) SetString(key string, value string) {
	p.set(key, value)
}

// RemoveValue deletes a value on the given key
func (p *InMemoryPreferences) RemoveValue(key string) {
	p.remove(key)
}

// NewInMemoryPreferences creates a new preferences implementation stored in memory
func NewInMemoryPreferences() *InMemoryPreferences {
	p := &InMemoryPreferences{}
	p.values = make(map[string]interface{})
	return p
}
