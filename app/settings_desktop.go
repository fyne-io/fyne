// +build !android,!ios,!mobile,!nacl

package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne"
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

	err := os.Mkdir(dir, 0700)
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
			if event.Op&fsnotify.Remove != 0 { // if it was deleted then watch again
				watcher.Remove(path)
				watchFileAddTarget(watcher, path)
			} else {
				callback()
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
}

func (s *settings) stopWatching() {
	if s.watcher == nil {
		return
	}

	s.watcher.(*fsnotify.Watcher).Close()
}
