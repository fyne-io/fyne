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
	p.watcherPath = p.storagePath()
	p.watcher = watchFile(p.watcherPath, func() {
		p.prefLock.RLock()
		shouldIgnoreChange := p.savedRecently
		p.prefLock.RUnlock()
		if shouldIgnoreChange {
			return
		}

		p.load()
	})
}

// ensureWatching adds the watcher target if it wasn't added during watch() due to
// the directory not existing yet. Called after directory creation during first save.
func (p *preferences) ensureWatching() {
	if p.watcher == nil || p.watcherPath == "" {
		return
	}
	addWatcherTarget(p.watcher, p.watcherPath)
}
