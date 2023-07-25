package internal

import (
	"reflect"
	"sync"

	"fyne.io/fyne/v2"
)

// InMemoryPreferences provides an implementation of the fyne.Preferences API that is stored in memory.
type InMemoryPreferences struct {
	values          map[string]interface{}
	lock            sync.RWMutex
	changeListeners []func()
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

// Bool looks up a boolean value for the key
func (p *InMemoryPreferences) Bool(key string) bool {
	return p.BoolWithFallback(key, false)
}

func (p *InMemoryPreferences) BoolList(key string) []bool {
	return p.BoolListWithFallback(key, []bool{})
}

func (p *InMemoryPreferences) BoolListWithFallback(key string, fallback []bool) []bool {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}
	valb, ok := value.([]bool)
	if !ok {
		return fallback
	}
	return valb
}

// BoolWithFallback looks up a boolean value and returns the given fallback if not found
func (p *InMemoryPreferences) BoolWithFallback(key string, fallback bool) bool {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}
	valb, ok := value.(bool)
	if !ok {
		return fallback
	}
	return valb
}

// ChangeListeners returns the list of listeners registered for this set of preferences.
func (p *InMemoryPreferences) ChangeListeners() []func() {
	return p.changeListeners
}

// Float looks up a float64 value for the key
func (p *InMemoryPreferences) Float(key string) float64 {
	return p.FloatWithFallback(key, 0.0)
}

func (p *InMemoryPreferences) FloatList(key string) []float64 {
	return p.FloatListWithFallback(key, []float64{})
}

func (p *InMemoryPreferences) FloatListWithFallback(key string, fallback []float64) []float64 {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}
	valf, ok := value.([]float64)
	if ok {
		return valf
	}
	vali, ok := value.([]int)
	if ok {
		flts := make([]float64, len(vali))
		for i, f := range vali {
			flts[i] = float64(f)
		}
		return flts
	}
	return fallback
}

// FloatWithFallback looks up a float64 value and returns the given fallback if not found
func (p *InMemoryPreferences) FloatWithFallback(key string, fallback float64) float64 {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}
	valf, ok := value.(float64)
	if ok {
		return valf
	}
	vali, ok := value.(int)
	if ok {
		return float64(vali)
	}
	return fallback
}

// Int looks up an integer value for the key
func (p *InMemoryPreferences) Int(key string) int {
	return p.IntWithFallback(key, 0)
}

func (p *InMemoryPreferences) IntList(key string) []int {
	return p.IntListWithFallback(key, []int{})
}

func (p *InMemoryPreferences) IntListWithFallback(key string, fallback []int) []int {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}
	vali, ok := value.([]int)
	if ok {
		return vali
	}
	// integers can be de-serialised as floats, so support both
	valf, ok := value.([]float64)
	if ok {
		ints := make([]int, len(valf))
		for i, f := range valf {
			ints[i] = int(f)
		}
		return ints
	}
	return fallback
}

// IntWithFallback looks up an integer value and returns the given fallback if not found
func (p *InMemoryPreferences) IntWithFallback(key string, fallback int) int {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}
	vali, ok := value.(int)
	if ok {
		return vali
	}
	// integers can be de-serialised as floats, so support both
	valf, ok := value.(float64)
	if !ok {
		return fallback
	}
	return int(valf)
}

// ReadValues provides read access to the underlying value map - for internal use only...
// You should not retain a reference to the map nor write to the values in the callback function
func (p *InMemoryPreferences) ReadValues(fn func(map[string]interface{})) {
	p.lock.RLock()
	fn(p.values)
	p.lock.RUnlock()
}

// RemoveValue deletes a value on the given key
func (p *InMemoryPreferences) RemoveValue(key string) {
	p.remove(key)
}

// SetBool saves a boolean value for the given key
func (p *InMemoryPreferences) SetBool(key string, value bool) {
	p.set(key, value)
}

func (p *InMemoryPreferences) SetBoolList(key string, value []bool) {
	p.set(key, value)
}

// SetFloat saves a float64 value for the given key
func (p *InMemoryPreferences) SetFloat(key string, value float64) {
	p.set(key, value)
}

func (p *InMemoryPreferences) SetFloatList(key string, value []float64) {
	p.set(key, value)
}

// SetInt saves an integer value for the given key
func (p *InMemoryPreferences) SetInt(key string, value int) {
	p.set(key, value)
}

func (p *InMemoryPreferences) SetIntList(key string, value []int) {
	p.set(key, value)
}

// SetString saves a string value for the given key
func (p *InMemoryPreferences) SetString(key string, value string) {
	p.set(key, value)
}

func (p *InMemoryPreferences) SetStringList(key string, value []string) {
	p.set(key, value)
}

// String looks up a string value for the key
func (p *InMemoryPreferences) String(key string) string {
	return p.StringWithFallback(key, "")
}

func (p *InMemoryPreferences) StringList(key string) []string {
	return p.StringListWithFallback(key, []string{})
}

func (p *InMemoryPreferences) StringListWithFallback(key string, fallback []string) []string {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}
	vals, ok := value.([]string)
	if !ok {
		return fallback
	}
	return vals
}

// StringWithFallback looks up a string value and returns the given fallback if not found
func (p *InMemoryPreferences) StringWithFallback(key, fallback string) string {
	value, ok := p.get(key)
	if !ok {
		return fallback
	}
	vals, ok := value.(string)
	if !ok {
		return fallback
	}
	return vals
}

// WriteValues provides write access to the underlying value map - for internal use only...
// You should not retain a reference to the map passed to the callback function
func (p *InMemoryPreferences) WriteValues(fn func(map[string]interface{})) {
	p.lock.Lock()
	fn(p.values)
	p.lock.Unlock()

	p.fireChange()
}

// NewInMemoryPreferences creates a new preferences implementation stored in memory
func NewInMemoryPreferences() *InMemoryPreferences {
	return &InMemoryPreferences{values: make(map[string]interface{})}
}

func (p *InMemoryPreferences) fireChange() {
	p.lock.RLock()
	listeners := p.changeListeners
	p.lock.RUnlock()

	var wg sync.WaitGroup

	for _, l := range listeners {
		wg.Add(1)
		go func(listener func()) {
			defer wg.Done()
			listener()
		}(l)
	}

	wg.Wait()
}

func (p *InMemoryPreferences) get(key string) (interface{}, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	v, err := p.values[key]
	return v, err
}

func (p *InMemoryPreferences) remove(key string) {
	p.lock.Lock()
	delete(p.values, key)
	p.lock.Unlock()

	p.fireChange()
}

func (p *InMemoryPreferences) set(key string, value interface{}) {
	p.lock.Lock()

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		s := reflect.ValueOf(value)
		old := reflect.ValueOf(p.values[key])
		if p.values[key] != nil && s.Len() == old.Len() {
			changed := false
			for i := 0; i < s.Len(); i++ {
				if s.Index(i).Interface() != old.Index(i).Interface() {
					changed = true
					break
				}
			}
			if !changed {
				p.lock.Unlock()
				return
			}
		}
	} else {
		if stored, ok := p.values[key]; ok && stored == value {
			p.lock.Unlock()
			return
		}
	}

	p.values[key] = value
	p.lock.Unlock()

	p.fireChange()
}
