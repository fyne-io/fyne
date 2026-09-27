//go:build wasm

package app

import (
	"encoding/base64"
	"syscall/js"
)

const secretPreferencesLocalStorageKey = "fyne-secret-preferences"

// newPlainSecretStore returns the browser local storage entry that holds the encrypted secret preferences.
func (a *fyneApp) newPlainSecretStore() secretStore {
	return &localStorageSecretStore{key: secretPreferencesLocalStorageKey}
}

type localStorageSecretStore struct {
	key string
}

func (s *localStorageSecretStore) load() ([]byte, error) {
	data := js.Global().Get("localStorage").Call("getItem", js.ValueOf(s.key))
	if data.IsNull() || data.IsUndefined() {
		return nil, errEmptyPreferencesStore
	}
	return base64.StdEncoding.DecodeString(data.String())
}

func (s *localStorageSecretStore) save(data []byte) error {
	encoded := base64.StdEncoding.EncodeToString(data)
	js.Global().Get("localStorage").Call("setItem", js.ValueOf(s.key), js.ValueOf(encoded))
	return nil
}
