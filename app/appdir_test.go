//go:build !android && !ios && !mobile

package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFyneApp_appDir(t *testing.T) {
	parent := t.TempDir()
	fyneDir := filepath.Join(parent, "fyne")
	testDir := filepath.Join(parent, "fyne-test")
	require.NoError(t, os.MkdirAll(filepath.Join(fyneDir, "io.fyne.old"), 0o700))

	noID, noIDOld := &fyneApp{}, &fyneApp{}
	require.NoError(t, os.MkdirAll(filepath.Join(fyneDir, noIDOld.UniqueID()), 0o700))

	for name, tt := range map[string]struct {
		app  *fyneApp
		root string
		want string
	}{
		"existing app keeps fyne subdirectory":      {&fyneApp{uniqueID: "io.fyne.old"}, fyneDir, filepath.Join(fyneDir, "io.fyne.old")},
		"new app moves out of fyne subdirectory":    {&fyneApp{uniqueID: "io.fyne.new"}, fyneDir, filepath.Join(parent, "io.fyne.new")},
		"root not named fyne is left alone":         {&fyneApp{uniqueID: "io.fyne.new"}, testDir, filepath.Join(testDir, "io.fyne.new")},
		"relative root is left alone":               {&fyneApp{uniqueID: "io.fyne.new"}, "fyne", filepath.Join("fyne", "io.fyne.new")},
		"app without ID stays in fyne subdirectory": {noID, fyneDir, filepath.Join(fyneDir, noID.UniqueID())},
		"app without ID keeps existing data":        {noIDOld, fyneDir, filepath.Join(fyneDir, noIDOld.UniqueID())},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.app.appDir(tt.root))
		})
	}
}
