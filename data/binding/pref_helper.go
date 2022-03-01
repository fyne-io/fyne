package binding

import (
	"sync"

	"fyne.io/fyne/v2"
)

type preferenceItem interface {
	checkForChange()
}

type preferenceBindings struct {
	items sync.Map // map[string]preferenceItem
}

func (b *preferenceBindings) getItem(key string) preferenceItem {
	val, loaded := b.items.Load(key)
	if !loaded {
		return nil
	}
	return val.(preferenceItem)
}

func (b *preferenceBindings) list() []preferenceItem {
	ret := []preferenceItem{}
	b.items.Range(func(_, val interface{}) bool {
		ret = append(ret, val.(preferenceItem))
		return true
	})
	return ret
}

func (b *preferenceBindings) setItem(key string, item preferenceItem) {
	b.items.Store(key, item)
}

type preferencesMap struct {
	prefs sync.Map // map[fyne.Preferences]*preferenceBindings
}

func newPreferencesMap() *preferencesMap {
	return &preferencesMap{}
}

func (m *preferencesMap) ensurePreferencesAttached(p fyne.Preferences) *preferenceBindings {
	binds, loaded := m.prefs.LoadOrStore(p, &preferenceBindings{})
	if loaded {
		return binds.(*preferenceBindings)
	}

	p.AddChangeListener(func() { m.preferencesChanged(p) })
	return binds.(*preferenceBindings)
}

func (m *preferencesMap) getBindings(p fyne.Preferences) *preferenceBindings {
	binds, loaded := m.prefs.Load(p)
	if !loaded {
		return nil
	}
	return binds.(*preferenceBindings)
}

func (m *preferencesMap) preferencesChanged(p fyne.Preferences) {
	binds := m.getBindings(p)
	if binds == nil {
		return
	}
	for _, item := range binds.list() {
		item.checkForChange()
	}
}
