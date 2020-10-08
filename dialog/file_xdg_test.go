// +build linux openbsd freebsd netbsd
// +build !android

package dialog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFavoriteLocations(t *testing.T) {
	favoriteLocations, _ := getFavoriteLocations()

	// This is a map containing the base name
	// of every key that should be in favoriteLocations
	expected := map[string]string{
		"Home":      "",
		"Documents": "Documents",
		"Downloads": "Downloads",
	}

	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)

	for e := range expected {
		fav, ok := favoriteLocations[e]
		if assert.True(t, ok) {
			if expected[e] == "" {
				assert.Equal(t, homeDir, fav)
			} else {
				assert.Equal(t, homeDir, filepath.Dir(fav))
				if _, err := os.Stat(fav); os.IsNotExist(err) {
					assert.Equal(t, expected[e], filepath.Base(fav))
				}
			}
		}
	}
}
