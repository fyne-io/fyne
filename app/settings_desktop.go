//go:build !android && !ios && !mobile && !wasm && !test_web_driver && !tamago && !noos && !tinygo

package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"github.com/fsnotify/fsnotify"
)

func watchFileAddTarget(watcher *fsnotify.Watcher, path string) {
	dir := filepath.Dir(path)
	ensureDirExists(dir)

	err := watcher.Add(dir)
	if err != nil {
		fyne.LogError("Settings watch error:", err)
	}
}

func ensureDirExists(dir string) {
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		return
	}

	err := os.MkdirAll(dir, 0o700)
	if err != nil {
		fyne.LogError("Unable to create settings storage:", err)
	}
}

func watchFile(path string, callback func()) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fyne.LogError("Failed to watch settings file:", err)
		return nil
	}

	go func() {
		for event := range watcher.Events {
			if event.Op.Has(fsnotify.Remove) { // if it was deleted then watch again
				watcher.Remove(path) // fsnotify returns false positives, see https://github.com/fsnotify/fsnotify/issues/268

				watchFileAddTarget(watcher, path)
			} else {
				fyne.Do(callback)
			}
		}

		err = watcher.Close()
		if err != nil {
			fyne.LogError("Settings un-watch error:", err)
		}
	}()

	watchFileAddTarget(watcher, path)
	return watcher
}

func (s *settings) watchSettings() {
	if s.themeSpecified {
		return // we only watch for theme changes at this time so don't bother
	}
	s.watcher = watchFile(s.schema.StoragePath(), s.fileChanged)

	a := fyne.CurrentApp()
	if a != nil && s != nil && a.Settings() == s { // ignore if testing
		watchTheme(s)
	}
}

func (s *settings) stopWatching() {
	if s.watcher == nil {
		return
	}

	s.watcher.(*fsnotify.Watcher).Close() // fsnotify returns false positives, see https://github.com/fsnotify/fsnotify/issues/268
}
