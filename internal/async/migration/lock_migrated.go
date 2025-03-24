//go:build migrated_fynedo

package migration

// Mutex behaves like a lock when app has not migrated to Fyne Do, otherwise does nothing.
type Mutex struct{}

func (m *Mutex) Lock() {}

func (m *Mutex) Unlock() {}
