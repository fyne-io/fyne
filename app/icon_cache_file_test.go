package app

import (
	_ "embed"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/fyne.png
var iconData []byte

func TestCachedIcon_PATH(t *testing.T) {
	SetMetadata(fyne.AppMetadata{})
	a := &fyneApp{uniqueID: "icontest"}
	assert.Equal(t, "", a.cachedIconPath())

	a.SetIcon(fyne.NewStaticResource("dummy", iconData))
	path := a.cachedIconPath()
	if path == "" {
		t.Error("cache path not constructed")
		return
	} else {
		defer os.Remove(path)
	}

	info, err := os.Stat(path)
	assert.NoError(t, err)
	assert.Equal(t, "icon.png", info.Name())

	assert.NoError(t, err)
}
