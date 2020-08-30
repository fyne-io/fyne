// +build !android,!ios,!mobile,!nacl

package app

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func TestDefaultTheme(t *testing.T) {
	if runtime.GOOS != "darwin" { // system defines default for macOS
		assert.Equal(t, theme.DarkTheme(), defaultTheme())
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := testPath("fynetest")

	ensureDirExists(tmpDir)
	if st, err := os.Stat(tmpDir); err != nil || !st.IsDir() {
		t.Error("Could not ensure directory exists")
	}

	os.Remove(tmpDir)
}

func TestWatchSettings(t *testing.T) {
	settings := &settings{}
	listener := make(chan fyne.Settings, 1)
	settings.AddChangeListener(listener)

	settings.fileChanged() // simulate the settings file changing

	select {
	case <-listener:
	case <-time.After(100 * time.Millisecond):
		t.Error("Settings listener was not called")
	}
}

func TestWatchFile(t *testing.T) {
	path := testPath("fyne-temp-watch.txt")
	f, _ := os.Create(path)
	f.Close()
	defer os.Remove(path)

	called := make(chan interface{}, 1)
	watchFile(path, func() {
		called <- true
	})
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	file.WriteString(" ")
	file.Close()

	select {
	case <-called:
	case <-time.After(100 * time.Millisecond):
		t.Error("File watcher callback was not called")
	}
}

func TestFileWatcher_FileDeleted(t *testing.T) {
	path := testPath("fyne-temp-watch.txt")
	f, _ := os.Create(path)
	f.Close()
	defer os.Remove(path)

	called := make(chan interface{}, 1)
	watcher := watchFile(path, func() {
		called <- true
	})
	if watcher == nil {
		assert.Fail(t, "Could not start watcher")
		return
	}

	defer watcher.Close()
	os.Remove(path)
	f, _ = os.Create(path)

	select {
	case <-called:
	case <-time.After(100 * time.Millisecond):
		t.Error("File watcher callback was not called")
	}
	f.Close()
}

func testPath(child string) string {
	// TMPDIR would be more normal but fsnotify cannot watch that on macOS...
	return filepath.Join("testdata", child)
}
