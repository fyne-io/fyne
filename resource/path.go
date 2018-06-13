package resource

import "os"
import "path"
import "io/ioutil"
import "os/user"

var rootDir string

func lookupRoot() string {
	// TODO replace this with io.UserCacheDir once Go 11 is out
	user, err := user.Current()
	if err == nil {
		return path.Join(user.HomeDir, ".fyne")
	}

	root, _ := ioutil.TempDir("", "fyne")
	return root
}

func rootPath() string {
	if rootDir == "" {
		rootDir = lookupRoot()
		if !pathExists(rootDir) {
			os.Mkdir(rootDir, 0700)
		}
	}

	return rootDir
}

func filePath(name string) string {
	return path.Join(rootPath(), name)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
