package internal

import (
	"sync"

	"fyne.io/fyne/v2"
)

// InMemoryPreferences provides an implementation of the fyne.Preferences API that is stored in memory.
type InMemoryPreferences struct {
	values          map[string]interface{}
	lock            sync.RWMutex
	changeListeners []func()
	wg              *sync.WaitGroup
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*InMemoryPreferences)(nil)

// AddChangeListener allows code to be notified when some preferences change. This will fire on any update.
// The passed 'listener' should not try to write values.
func (p *InMemoryPreferences) AddChangeListener(listener func()) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.changeListeners = append(p.changeListeners, listener)
}

// ReadValues provides read access to the underlying value map - for internal use only...
// You should not retain a reference to the map nor write to the values in the callback function
func (p *InMemoryPreferences) ReadValues(fn func(map[string]interface{})) {
	p.lock.RLock()
	fn(p.values)
	p.lock.RUnlock()
}

// WriteValues provides write access to the underlying value map - for internal use only...
// You should not retain a reference to the map passed to the callback function
func (p *InMemoryPreferences) WriteValues(fn func(map[string]interface{})) {
	p.lock.Lock()
	fn(p.values)
	p.lock.Unlock()

	p.fireChange()
}

func (p *InMemoryPreferences) set(key string, value interface{}) {
	p.lock.Lock()

	if stored, ok := p.values[key]; ok && stored == value {
		p.lock.Unlock()
		return
	}

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

func (p *InMemoryPreferences) fireChange() {
	p.lock.RLock()
	listeners := p.changeListeners
	p.lock.RUnlock()

	for _, l := range listeners {
		p.wg.Add(1)
		go func(listener func()) {
			defer p.wg.Done()
			listener()
		}(l)
	}

	p.wg.Wait()
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
	return &InMemoryPreferences{
		values: make(map[string]interface{}),
		wg:     &sync.WaitGroup{},
	}
}
