package app

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/theme"
)

func TestSettingsLoad(t *testing.T) {
	testPath := filepath.Join(os.TempDir(), "fyne-settings-test.json")

	file, _ := os.Create(testPath)
	defer os.Remove(testPath)
	buf := bufio.NewWriter(file)
	buf.WriteString("{\"theme\":\"light\"}")
	buf.Flush()
	file.Close()

	settings := &settings{}
	err := settings.loadFromFile(testPath)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "light", settings.schema.ThemeName)

	os.Remove(testPath)
	file, _ = os.Create(testPath)
	buf = bufio.NewWriter(file)
	buf.WriteString("{\"theme\":\"dark\"}")
	buf.Flush()
	file.Close()

	err = settings.loadFromFile(testPath)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "dark", settings.schema.ThemeName)
}

func TestOverrideTheme(t *testing.T) {
	set := &settings{}
	set.setupTheme()
	assert.Equal(t, theme.DarkTheme(), set.Theme())

	set.schema.ThemeName = "light"
	set.setupTheme()
	assert.Equal(t, theme.LightTheme(), set.Theme())

	set = &settings{}
	set.setupTheme()
	assert.Equal(t, theme.DarkTheme(), set.Theme())

	err := os.Setenv("FYNE_THEME", "light")
	if err != nil {
		t.Error(err)
	}
	set.setupTheme()
	assert.Equal(t, theme.LightTheme(), set.Theme())
	err = os.Setenv("FYNE_THEME", "")
	if err != nil {
		t.Error(err)
	}
}
