// +build linux openbsd freebsd netbsd
// +build !android

package dialog

import (
	"os"
	"path/filepath"
	"strings"
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

	for name, subdir := range expected {
		fav, ok := favoriteLocations[name]
		if !ok {
			t.Errorf("missing favourite location: %s", name)
			continue
		}

		if subdir == "" {
			assert.Equal(t, homeDir, fav)
		} else {
			//			assert.Equal(t, homeDir, filepath.Dir(fav)) // I don't think this assertion must be true, the user could specify a different parent directory
			_, err := os.Stat(fav)
			if err == nil {
				continue // folder found
			}
			if !os.IsNotExist(err) {
				t.Errorf("failed to read directory %s", fav)
				continue
			}
			fallbackName := filepath.Base(fav)
			if !strings.EqualFold(subdir, fallbackName) {
				t.Errorf("%s should equal fold with %s", subdir, fallbackName)
			}
		}
	}
}
