//go:build !migrated_fynedo

package migrated

import "fyne.io/fyne/v2/internal/async"

// Pool is the generic version of sync.Pool.
type Pool[T any] struct {
	async.Pool[T]
}
