package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/cmd/fyne/internal/metadata"
)

func Test_calculateExeName(t *testing.T) {
	modulesApp := calculateExeName("testdata/modules_app", "windows")
	assert.Equal(t, "module.exe", modulesApp)

	modulesShortName := calculateExeName("testdata/short_module", "linux")
	assert.Equal(t, "app", modulesShortName)

	nonModulesApp := calculateExeName("testdata", "linux")
	assert.Equal(t, "testdata", nonModulesApp)
}

func Test_isValidVersion(t *testing.T) {
	assert.True(t, isValidVersion("1"))
	assert.True(t, isValidVersion("1.2"))
	assert.True(t, isValidVersion("1.2.3"))

	assert.False(t, isValidVersion("1.2.3.4"))
	assert.False(t, isValidVersion(""))
	assert.False(t, isValidVersion("1.2-alpha3"))
	assert.False(t, isValidVersion("pre1"))
	assert.False(t, isValidVersion("1..2"))
}

func Test_MergeMetata(t *testing.T) {
	p := &Packager{appVersion: "v0.1"}
	data := &metadata.FyneApp{
		Details: metadata.AppDetails{
			Icon:    "test.png",
			Build:   3,
			Version: "v0.0.1",
		},
	}

	mergeMetadata(p, data)
	assert.Equal(t, "v0.1", p.appVersion)
	assert.Equal(t, 3, p.appBuild)
	assert.Equal(t, "test.png", p.icon)
}

func Test_validateAppID(t *testing.T) {
	id, err := validateAppID("myApp", "windows", "myApp", false)
	assert.Nil(t, err)
	assert.Equal(t, "myApp", id)

	id, err = validateAppID("", "darwin", "myApp", true)
	assert.Nil(t, err)
	assert.Equal(t, "com.example.myApp", id) // this was in for compatibility

	id, err = validateAppID("com.myApp", "darwin", "myApp", true)
	assert.Nil(t, err)
	assert.Equal(t, "com.myApp", id)

	_, err = validateAppID("", "ios", "myApp", false)
	assert.NotNil(t, err)

	_, err = validateAppID("myApp", "ios", "myApp", false)
	assert.NotNil(t, err)

	id, err = validateAppID("com.myApp", "android", "myApp", true)
	assert.Nil(t, err)
	assert.Equal(t, "com.myApp", id)

	_, err = validateAppID("myApp", "android", "myApp", true)
	assert.NotNil(t, err)
}
