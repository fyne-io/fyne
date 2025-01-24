package app

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
)

type preferences struct {
	*internal.InMemoryPreferences

	prefLock            sync.RWMutex
	savedRecently       bool
	changedDuringSaving bool

	app                 *fyneApp
	needsSaveBeforeExit bool
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*preferences)(nil)

// sentinel error to signal an empty preferences storage backend was loaded
var errEmptyPreferencesStore = errors.New("empty preferences store")

// returned from storageWriter() - may be a file, browser local storage, etc
type writeSyncCloser interface {
	io.WriteCloser
	Sync() error
}

// forceImmediateSave writes preferences to storage immediately, ignoring the debouncing
// logic in the change listener. Does nothing if preferences are not backed with a persistent store.
func (p *preferences) forceImmediateSave() {
	if !p.needsSaveBeforeExit {
		return
	}
	err := p.save()
	if err != nil {
		fyne.LogError("Failed on force saving preferences", err)
	}
}

func (p *preferences) resetSavedRecently() {
	go func() {
		time.Sleep(time.Millisecond * 100) // writes are not always atomic. 10ms worked, 100 is safer.

		// For test reasons we need to use current app not what we were initialised with as they can differ
		fyne.DoAndWait(func() {
			p.prefLock.Lock()
			p.savedRecently = false
			changedDuringSaving := p.changedDuringSaving
			p.changedDuringSaving = false
			p.prefLock.Unlock()

			if changedDuringSaving {
				p.save()
			}
		})
	}()
}

func (p *preferences) save() error {
	storage, err := p.storageWriter()
	if err != nil {
		return err
	}
	return p.saveToStorage(storage)
}

func (p *preferences) saveToStorage(writer writeSyncCloser) error {
	p.prefLock.Lock()
	p.savedRecently = true
	p.prefLock.Unlock()
	defer p.resetSavedRecently()

	defer writer.Close()
	encode := json.NewEncoder(writer)

	var err error
	p.InMemoryPreferences.ReadValues(func(values map[string]any) {
		err = encode.Encode(&values)
	})

	err2 := writer.Sync()
	if err == nil {
		err = err2
	}
	return err
}

func (p *preferences) load() {
	storage, err := p.storageReader()
	if err == nil {
		err = p.loadFromStorage(storage)
	}
	if err != nil && err != errEmptyPreferencesStore {
		fyne.LogError("Preferences load error:", err)
	}
}

func (p *preferences) loadFromStorage(storage io.ReadCloser) (err error) {
	defer func() {
		if r := storage.Close(); r != nil && err == nil {
			err = r
		}
	}()
	decode := json.NewDecoder(storage)

	p.InMemoryPreferences.WriteValues(func(values map[string]any) {
		err = decode.Decode(&values)
		if err != nil {
			return
		}
		convertLists(values)
	})

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

	p.needsSaveBeforeExit = true
	p.AddChangeListener(func() {
		if p != app.prefs {
			return
		}
		p.prefLock.Lock()
		shouldIgnoreChange := p.savedRecently
		if p.savedRecently {
			p.changedDuringSaving = true
		}
		p.prefLock.Unlock()

		if shouldIgnoreChange { // callback after loading from storage, or too many updates in a row
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

func convertLists(values map[string]any) {
	for k, v := range values {
		if items, ok := v.([]any); ok {
			if len(items) == 0 {
				continue
			}

			switch items[0].(type) {
			case bool:
				bools := make([]bool, len(items))
				for i, item := range items {
					bools[i] = item.(bool)
				}
				values[k] = bools
			case float64:
				floats := make([]float64, len(items))
				for i, item := range items {
					floats[i] = item.(float64)
				}
				values[k] = floats
			//case int: // json has no int!
			case string:
				strings := make([]string, len(items))
				for i, item := range items {
					strings[i] = item.(string)
				}
				values[k] = strings
			}
		}
	}
}
