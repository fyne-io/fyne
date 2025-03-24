//go:build !migrated_fynedo

package migration

import "sync"

// Mutex behaves like a lock when app has not migrated to Fyne Do, otherwise does nothing.
type Mutex = sync.Mutex
