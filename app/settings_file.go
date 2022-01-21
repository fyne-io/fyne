//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package app

import (
	"encoding/json"
	"io"
	"os"

	"fyne.io/fyne/v2"
)

func (s *settings) load() {
	err := s.loadFromFile(s.schema.StoragePath())
	if err != nil && err != io.EOF { // we can get an EOF in windows settings writes
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
