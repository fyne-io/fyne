//go:build !migrated_fynedo

package migrated

import "fyne.io/fyne/v2/internal/async"

// Map is a generic wrapper around [sync.Map].
type Map[K, V any] struct {
	async.Map[K, V]
}
