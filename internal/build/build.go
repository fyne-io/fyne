// Package build contains information about they type of build currently running.
package build

import "fyne.io/fyne/v2"

func MigratedToFyneDo() bool {
	if DisableThreadChecks {
		return true
	}

	v, ok := fyne.CurrentApp().Metadata().Migrations["fyneDo"]
	if ok {
		return v
	}
	return false
}
