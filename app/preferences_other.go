//go:build !ios && !android && !mobile && !wasm

package app

import (
	"path/filepath"

	"fyne.io/fyne/v2/internal/app"
)

// storagePath returns the location of the settings storage
func (p *preferences) storagePath() string {
	return filepath.Join(p.app.storageRoot(), "preferences.json")
}

// storageRoot returns the location of the app storage
func (a *fyneApp) storageRoot() string {
	return filepath.Join(app.RootConfigDir(), a.UniqueID())
}

func (p *preferences) watch() {
	watchFile(p.storagePath(), func() {
		p.prefLock.RLock()
		shouldIgnoreChange := p.savedRecently
		p.prefLock.RUnlock()
		if shouldIgnoreChange {
			return
		}

		p.load()
	})
}
