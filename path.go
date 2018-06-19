package fyne

import "os"
import "path"
import "io/ioutil"
import "os/user"

var cacheDirPath string

func lookupCacheDir() string {
	// TODO replace this with io.UserCacheDir once Go 11 is out
	user, err := user.Current()
	if err == nil {
		return path.Join(user.HomeDir, ".fyne")
	}

	root, _ := ioutil.TempDir("", "fyne")
	return root
}

func cacheDir() string {
	if cacheDirPath == "" {
		cacheDirPath = lookupCacheDir()
		if !pathExists(cacheDirPath) {
			os.Mkdir(cacheDirPath, 0700)
		}
	}

	return cacheDirPath
}

func cachePath(name string) string {
	return path.Join(cacheDir(), name)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
