package util

import (
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

// Exists will return true if the passed path exists on the current system.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// CopyFile copies the content of a regular file, source, into target path.
func CopyFile(source, target string) error {
	return copyFileMode(source, target, 0644)
}

// CopyExeFile copies the content of an executable file, source, into target path.
func CopyExeFile(src, tgt string) error {
	return copyFileMode(src, tgt, 0755)
}

// EnsureSubDir will make sure a named directory exists within the parent - creating it if not.
func EnsureSubDir(parent, name string) string {
	path := filepath.Join(parent, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fyne.LogError("Failed to create directory", err)
		}
	}
	return path
}

// EnsureAbsPath returns the absolute path of the given file if it is not already absolute
func EnsureAbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		fyne.LogError("Failed to find absolute path", err)
		return path
	}
	return abs
}

// MakePathRelativeTo joins root with path if path is not absolute and exists in root
func MakePathRelativeTo(root, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	joined := filepath.Join(root, path)
	if !Exists(joined) {
		return path
	}
	return joined
}

func copyFileMode(src, tgt string, perm os.FileMode) (err error) {
	if _, err := os.Stat(src); err != nil {
		return err
	}

	source, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.OpenFile(tgt, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer func() {
		if r := target.Close(); r != nil && err == nil {
			err = r
		}
	}()

	_, err = io.Copy(target, source)
	return err
}
