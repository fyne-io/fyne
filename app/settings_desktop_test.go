// +build !android,!ios,!mobile

package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/theme"
)

func TestDefaultTheme(t *testing.T) {
	assert.Equal(t, theme.DarkTheme(), defaultTheme())
}

func TestEnsureDir(t *testing.T) {
	tmpDir := filepath.Join(rootConfigDir(), "fynetest")

	ensureDirExists(tmpDir)
	if st, err := os.Stat(tmpDir); err != nil || !st.IsDir() {
		t.Error("Could not ensure directory exists")
	}

	os.Remove(tmpDir)
}
