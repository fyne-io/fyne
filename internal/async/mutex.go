//go:build !migrated_fynedo

package async

import "sync"

// Mutex is an alias for sync.Mutex on non-migrated builds
// In migrated builds it is a no-op.
type Mutex = sync.Mutex
