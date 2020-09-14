// +build !ios

package app

import (
	"encoding/json"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/internal"
)

type preferences struct {
	*internal.InMemoryPreferences

	app *fyneApp
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*preferences)(nil)

func (p *preferences) save() error {
	return p.saveToFile(p.storagePath())
}

func (p *preferences) saveToFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		file, err = os.Open(path) // #nosec
		if err != nil {
			return err
		}
	}
	encode := json.NewEncoder(file)

	values := p.InMemoryPreferences.Values()
	return encode.Encode(&values)
}

func (p *preferences) load(_ string) {
	err := p.loadFromFile(p.storagePath())
	if err != nil {
		fyne.LogError("Preferences load error:", err)
	}
}

func (p *preferences) loadFromFile(path string) error {
	file, err := os.Open(path) // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(path), 0700)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	values := p.InMemoryPreferences.Values()
	return decode.Decode(&values)
}

func newPreferences(app *fyneApp) *preferences {
	p := &preferences{}
	p.app = app
	p.InMemoryPreferences = internal.NewInMemoryPreferences()

	p.OnChange = func() {
		err := p.save()
		if err != nil {
			fyne.LogError("Failed on saving preferences", err)
		}
	}
	return p
}
