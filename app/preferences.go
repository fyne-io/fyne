package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
)

type preferences struct {
	*internal.InMemoryPreferences
	ignoreChange bool

	app *fyneApp
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*preferences)(nil)

func (p *preferences) resetIgnore() {
	go func() {
		time.Sleep(time.Millisecond * 100) // writes are not always atomic. 10ms worked, 100 is safer.
		p.ignoreChange = false
	}()
}

func (p *preferences) save() error {
	return p.saveToFile(p.storagePath())
}

func (p *preferences) saveToFile(path string) error {
	p.ignoreChange = true
	defer p.resetIgnore()
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

	p.InMemoryPreferences.ReadValues(func(values map[string]interface{}) {
		err = encode.Encode(&values)
	})
	return err
}

func (p *preferences) load() {
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
	defer file.Close()
	decode := json.NewDecoder(file)

	p.InMemoryPreferences.WriteValues(func(values map[string]interface{}) {
		err = decode.Decode(&values)
	})
	return err
}

func newPreferences(app *fyneApp) *preferences {
	p := &preferences{}
	p.app = app
	p.InMemoryPreferences = internal.NewInMemoryPreferences()

	// don't load or watch if not setup
	if app.uniqueID == "" {
		return p
	}

	p.AddChangeListener(func() {
		if p.ignoreChange { // callback after loading, no need to save
			return
		}

		err := p.save()
		if err != nil {
			fyne.LogError("Failed on saving preferences", err)
		}
	})
	p.watch()
	return p
}
