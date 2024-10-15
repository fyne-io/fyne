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
	a := &fyneApp{uniqueID: "icontest"}
	assert.Equal(t, "", a.cachedIconPath())

	a.SetIcon(fyne.NewStaticResource("dummy", iconData))
	path := a.cachedIconPath()
	assert.NotEqual(t, "", path)

	info, err := os.Stat(path)
	assert.Nil(t, err)
	assert.Equal(t, "icon.png", info.Name())

	err = os.Remove(path)
	assert.Nil(t, err)
}
