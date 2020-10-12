package dialog

import (
	"os"
	"path/filepath"
)

func getFavoriteLocations() (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Home":      homeDir,
		"Documents": filepath.Join(homeDir, "Documents"),
		"Downloads": filepath.Join(homeDir, "Downloads"),
	}, nil
}
