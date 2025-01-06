//go:build (linux || openbsd || freebsd || netbsd) && !android

package dialog

import (
	"os"
	"testing"

	"fyne.io/fyne/v2/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	homeURI := storage.NewFileURI(homeDir)

	for name, subdir := range expected {
		fav, ok := favoriteLocations[name]
		if !ok {
			continue
		}

		if subdir == "" {
			assert.Equal(t, homeURI.String(), fav.String())
		}
	}
}
