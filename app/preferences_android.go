// +build android

package app

import "path/filepath"

// storagePath returns the location of the settings storage
func (p *preferences) storagePath() string {
	// we have no global storage, use app global instead - rootConfigDir looks up in app_mobile_and.go
	return filepath.Join(rootConfigDir(), "storage", "preferences.json")
}
