//go:build !android && !ios && !mobile

package app

import (
	"os"
	"path/filepath"
)

// appDir returns the directory for this app's data given the "fyne" directory
// that used to hold all apps. Apps that already have data in there keep using
// it, new apps live directly under the parent (e.g. ~/.config/app_id).
// Apps without an ID have nothing to migrate, so they stay in the old place.
// So does everything under a root that is not an absolute "fyne" directory,
// which covers test builds and web where there is no platform location.
func (a *fyneApp) appDir(fyneDir string) string {
	old := filepath.Join(fyneDir, a.UniqueID())
	if a.missingID || !filepath.IsAbs(fyneDir) || filepath.Base(fyneDir) != "fyne" {
		return old
	}
	if _, err := os.Stat(old); err == nil {
		return old
	}
	return filepath.Join(filepath.Dir(fyneDir), a.UniqueID())
}
