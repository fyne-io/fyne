//go:build !ci && android

package app

import (
	"os"
	"testing"

	"fyne.io/fyne/v2/internal/app"

	"github.com/stretchr/testify/assert"
)

func Test_RootConfigDir(t *testing.T) {
	oldEnv := os.Getenv("FILESPATH")
	os.Setenv("FILESPATH", "/tmp")

	assert.Equal(t, "/tmp", app.RootConfigDir())
	os.Setenv("FILESPATH", oldEnv)
}

func Test_StoragePath(t *testing.T) {
	oldEnv := os.Getenv("FILESPATH")
	os.Setenv("FILESPATH", "/tmp")
	defer os.Setenv("FILESPATH", oldEnv)

	p := &preferences{}
	expected := "/tmp/storage/preferences.json"
	assert.Equal(t, expected, p.storagePath())
}
