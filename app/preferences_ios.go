//go:build ios
// +build ios

package app

import (
	"path/filepath"
)
import "C"

// storagePath returns the location of the settings storage
func (p *preferences) storagePath() string {
	ret := filepath.Join(p.app.storageRoot(), "preferences.json")
	return ret
}

// storageRoot returns the location of the app storage
func (a *fyneApp) storageRoot() string {
	return rootConfigDir() // we are in a sandbox, so no app ID added to this path
}

func (p *preferences) watch() {
	// no-op on mobile
}
