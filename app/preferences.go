package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
)

type preferences struct {
	*internal.InMemoryPreferences

	prefLock            sync.RWMutex
	loadingInProgress   bool
	savedRecently       bool
	changedDuringSaving bool

	app *fyneApp
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*preferences)(nil)

func (p *preferences) resetSavedRecently() {
	go func() {
		time.Sleep(time.Millisecond * 100) // writes are not always atomic. 10ms worked, 100 is safer.
		p.prefLock.Lock()
		p.savedRecently = false
		changedDuringSaving := p.changedDuringSaving
		p.changedDuringSaving = false
		p.prefLock.Unlock()

		if changedDuringSaving {
			p.save()
		}
	}()
}

func (p *preferences) save() error {
	return p.saveToFile(p.storagePath())
}

func (p *preferences) saveToFile(path string) error {
	p.prefLock.Lock()
	p.savedRecently = true
	p.prefLock.Unlock()
	defer p.resetSavedRecently()
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
	defer file.Close()
	encode := json.NewEncoder(file)

	p.InMemoryPreferences.ReadValues(func(values map[string]interface{}) {
		err = encode.Encode(&values)
	})

	err2 := file.Sync()
	if err == nil {
		err = err2
	}
	return err
}

func (p *preferences) load() {
	err := p.loadFromFile(p.storagePath())
	if err != nil {
		fyne.LogError("Preferences load error:", err)
	}
}

func (p *preferences) loadFromFile(path string) (err error) {
	file, err := os.Open(path) // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	defer func() {
		if r := file.Close(); r != nil && err == nil {
			err = r
		}
	}()
	decode := json.NewDecoder(file)

	p.prefLock.Lock()
	p.loadingInProgress = true
	p.prefLock.Unlock()

	p.InMemoryPreferences.WriteValues(func(values map[string]interface{}) {
		err = decode.Decode(&values)
	})

	p.prefLock.Lock()
	p.loadingInProgress = false
	p.prefLock.Unlock()

	return err
}

func newPreferences(app *fyneApp) *preferences {
	p := &preferences{}
	p.app = app
	p.InMemoryPreferences = internal.NewInMemoryPreferences()

	// don't load or watch if not setup
	if app.uniqueID == "" && app.Metadata().ID == "" {
		return p
	}

	p.AddChangeListener(func() {
		p.prefLock.Lock()
		shouldIgnoreChange := p.savedRecently || p.loadingInProgress
		if p.savedRecently && !p.loadingInProgress {
			p.changedDuringSaving = true
		}
		p.prefLock.Unlock()

		if shouldIgnoreChange { // callback after loading file, or too many updates in a row
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
