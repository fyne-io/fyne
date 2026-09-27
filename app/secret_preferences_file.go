//go:build !wasm

package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/internal/repository"
)

const secretPreferencesFile = "preferences.secret"

// newPlainSecretStore returns the file that holds the (already protected) secret preferences data.
func (a *fyneApp) newPlainSecretStore() secretStore {
	return &fileSecretStore{path: filepath.Join(a.storageRoot(), secretPreferencesFile)}
}

type fileSecretStore struct {
	path string
}

func (f *fileSecretStore) load() ([]byte, error) {
	data, err := os.ReadFile(f.path) // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errEmptyPreferencesStore
		}
		return nil, err
	}
	return data, nil
}

func (f *fileSecretStore) save(data []byte) error {
	err := os.MkdirAll(filepath.Dir(f.path), repository.PermUserReadWriteExec)
	if err != nil {
		return err
	}

	// write to a temporary file and rename so a crash mid-write cannot corrupt the store
	tmp := f.path + ".tmp"
	file, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, repository.PermUserRead|repository.PermUserWrite) // #nosec
	if err != nil {
		return err
	}
	if _, err = file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err = file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, f.path)
}
