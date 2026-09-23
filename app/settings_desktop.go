//go:build !android && !ios && !mobile && !wasm && !test_web_driver && !tamago && !noos && !tinygo

package app

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/repository"
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

	err := os.MkdirAll(dir, repository.PermUserReadWriteExec)
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
				err = watcher.Remove(path)
				if err != nil {
					fyne.LogError("failed to stop watching removed settings file", err)
					// fsnotify used to return false positives (https://github.com/fsnotify/fsnotify/issues/268).
					// So, don’t return but just continue here.
				}

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
	s.watcher = watchFile(s.schema.StoragePath(), s.fileChanged)

	if s.explicitThemeVariantName() == "" {
		a := fyne.CurrentApp()
		if a != nil && s != nil && a.Settings() == s { // ignore if testing
			watchTheme(s)
		}
	}
}

func (s *settings) stopWatching() {
	if s.watcher == nil {
		return
	}

	s.watcher.(*fsnotify.Watcher).Close() // fsnotify returns false positives, see https://github.com/fsnotify/fsnotify/issues/268
}
