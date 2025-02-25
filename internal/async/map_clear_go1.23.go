//go:build go1.23 && !migrated_fynedo

package async

// Clear deletes all the entries, resulting in an empty Map.
func (m *Map[K, V]) Clear() {
	m.Map.Clear() // More efficient than O(n) range and delete in older Go.
}
