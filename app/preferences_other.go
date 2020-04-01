// +build !ios,!android

package app

import "path/filepath"

// storagePath returns the location of the settings storage
func (p *preferences) storagePath() string {
	return filepath.Join(rootConfigDir(), p.uniqueID(), "preferences.json")
}
