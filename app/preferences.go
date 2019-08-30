package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/internal"
)

type preferences struct {
	*internal.InMemoryPreferences

	appID string
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*preferences)(nil)

func (p *preferences) uniqueID() string {
	if p.appID != "" {
		return p.appID
	}

	fyne.LogError("Preferences API requires a unique ID, use app.NewWithID()", nil)
	p.appID = fmt.Sprintf("missing-id-%d", time.Now().Unix()) // This is a fake unique - it just has to not be reused...
	return p.appID
}

// storagePath returns the location of the settings storage
func (p *preferences) storagePath() string {
	return filepath.Join(rootConfigDir(), p.uniqueID(), "preferences.json")
}

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
		file, err = os.Open(path)
		if err != nil {
			return err
		}
	}
	encode := json.NewEncoder(file)

	values := p.InMemoryPreferences.Values()
	return encode.Encode(&values)
}

func (p *preferences) load() {
	err := p.loadFromFile(p.storagePath())
	if err != nil {
		fyne.LogError("Preferences load error:", err)
	}
}

func (p *preferences) loadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(path), 0700)
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	values := p.InMemoryPreferences.Values()
	return decode.Decode(&values)
}

func newPreferences() *preferences {
	p := &preferences{}
	p.InMemoryPreferences = internal.NewInMemoryPreferences()

	p.OnChange = func() {
		p.save()
	}
	return p
}

func loadPreferences(id string) *preferences {
	p := newPreferences()
	p.appID = id

	p.load()
	return p
}
