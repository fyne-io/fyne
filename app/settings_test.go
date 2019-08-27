package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func TestSettingsLoad(t *testing.T) {
	settings := &settings{}

	err := settings.loadFromFile(filepath.Join("testdata", "light-theme.json"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "light", settings.schema.ThemeName)

	err = settings.loadFromFile(filepath.Join("testdata", "dark-theme.json"))
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

func TestWatchSettings(t *testing.T) {
	settings := &settings{}
	listener := make(chan fyne.Settings)
	settings.AddChangeListener(listener)

	settings.fileChanged() // simulate the settings file changing

	select {
	case _ = <-listener:
	case <-time.After(100 * time.Millisecond):
		t.Error("Settings listener was not called")
	}
}

func TestWatchFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "fyne-temp-watch.txt")
	os.Create(path)
	defer os.Remove(path)

	called := make(chan interface{})
	watchFile(path, func() {
		called <- true
	})
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	file.WriteString(" ")
	file.Close()

	select {
	case _ = <-called:
	case <-time.After(100 * time.Millisecond):
		t.Error("File watcher callback was not called")
	}
}

func TestFileWatcher_FileDeleted(t *testing.T) {
	path := filepath.Join(os.TempDir(), "fyne-temp-watch.txt")
	os.Create(path)
	defer os.Remove(path)

	called := make(chan interface{})
	watcher := watchFile(path, func() {
		called <- true
	})
	if watcher == nil {
		assert.Fail(t, "Could not start watcher")
		return
	}

	defer watcher.Close()
	os.Remove(path)
	os.Create(path)

	select {
	case _ = <-called:
	case <-time.After(100 * time.Millisecond):
		t.Error("File watcher callback was not called")
	}
}
