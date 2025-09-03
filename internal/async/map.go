//go:build !migrated_fynedo

package async

import "sync"

// Map is a generic wrapper around [sync.Map].
type Map[K, V any] struct {
	sync.Map
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key K) {
	m.Map.Delete(key)
}

// Len returns the length of the map. It is O(n) over the number of items.
func (m *Map[K, V]) Len() (count int) {
	m.Map.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

// Load returns the value stored in the map for a key, or nil if no value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	val, ok := m.Map.Load(key)
	if val == nil {
		return *new(V), ok
	}
	return val.(V), ok
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	val, loaded := m.Map.LoadAndDelete(key)
	if val == nil {
		return *new(V), loaded
	}
	return val.(V), loaded
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	act, loaded := m.Map.LoadOrStore(key, value)
	if act == nil {
		return *new(V), loaded
	}
	return act.(V), loaded
}

// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.Map.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}
