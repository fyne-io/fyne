// +build !js,!wasm,!web

package app

import (
	"encoding/json"
	"os"

	"fyne.io/fyne"
)

func (s *settings) load() {
	err := s.loadFromFile(s.schema.StoragePath())
	if err != nil {
		fyne.LogError("Settings load error:", err)
	}

	s.setupTheme()
}

func (s *settings) loadFromFile(path string) error {
	file, err := os.Open(path) // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	return decode.Decode(&s.schema)
}
