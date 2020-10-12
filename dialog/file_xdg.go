// +build linux openbsd freebsd netbsd
// +build !android

package dialog

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func getFavoriteLocation(homeDir, name string) (string, error) {
	cmdName := "xdg-user-dir"
	fallback := filepath.Join(homeDir, name)

	if _, err := exec.LookPath(cmdName); err != nil {
		return fallback, fmt.Errorf("%s not found in PATH. using fallback paths", cmdName)
	}

	cmd := exec.Command(cmdName, name)
	stdout := bytes.NewBufferString("")

	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return fallback, err
	}

	loc := stdout.String()
	// Remove \n at the end
	loc = loc[:len(loc)-1]

	if loc == homeDir {
		return fallback, fmt.Errorf("this computer does not have a %s folder", name)
	}

	return stdout.String(), nil
}

func getFavoriteLocations() (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	rv := map[string]string{
		"Home": homeDir,
	}
	documentsDir, err1 := getFavoriteLocation(homeDir, "DOCUMENTS")
	if err == nil {
		rv["Documents"] = documentsDir
	} else {
		rv["Documents"] = filepath.Join(homeDir, "Documents")
		err = err1
	}
	downloadsDir, err1 := getFavoriteLocation(homeDir, "DOWNLOADS")
	if err == nil {
		rv["Downloads"] = downloadsDir
	} else {
		rv["Downloads"] = filepath.Join(homeDir, "Downloads")
		err = err1
	}

	return rv, err
}
