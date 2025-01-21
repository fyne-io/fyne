// Package build contains information about they type of build currently running.
package build

import (
	"sync"

	"fyne.io/fyne/v2"
)

var (
	migrateCheck sync.Once

	migratedFyneDo bool
)

func MigratedToFyneDo() bool {
	if DisableThreadChecks {
		return true
	}

	migrateCheck.Do(func() {
		v, ok := fyne.CurrentApp().Metadata().Migrations["fyneDo"]
		if ok {
			migratedFyneDo = v
		}
	})

	return migratedFyneDo
}
