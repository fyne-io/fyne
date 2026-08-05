//go:build !go1.23 && !migrated_fynedo

package async

// Clear deletes all the entries, resulting in an empty Map.
// This is O(n) over the number of entries when not using Go 1.23 or newer.
func (m *Map[K, V]) Clear() {
	m.Map.Range(func(key, _ any) bool {
		m.Map.Delete(key)
		return true
	})
}
