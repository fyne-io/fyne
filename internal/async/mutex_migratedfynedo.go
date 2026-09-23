//go:build migrated_fynedo

package async

// Mutex is a no-op sync.Locker for when we are running single-threaded
type Mutex struct{}

func (Mutex) Lock() {}

func (Mutex) Unlock() {}
