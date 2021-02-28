package binding

import (
	"sync"

	"fyne.io/fyne/v2"
)

type preferenceItem interface {
	checkForChange()
}

type preferenceBindings struct {
	lock  sync.RWMutex
	items map[string]preferenceItem
}

func (b *preferenceBindings) setItem(key string, item preferenceItem) {
	b.lock.Lock()
	b.items[key] = item
	b.lock.Unlock()
}

func (b *preferenceBindings) getItem(key string) preferenceItem {
	b.lock.RLock()
	item := b.items[key]
	b.lock.RUnlock()
	return item
}

func (b *preferenceBindings) list() []preferenceItem {
	b.lock.RLock()
	items := make([]preferenceItem, 0, len(b.items))
	for _, i := range b.items {
		items = append(items, i)
	}
	b.lock.RUnlock()
	return items
}

type preferencesMap struct {
	lock  sync.RWMutex
	prefs map[fyne.Preferences]*preferenceBindings
}

func newPreferencesMap() *preferencesMap {
	return &preferencesMap{
		prefs: make(map[fyne.Preferences]*preferenceBindings),
	}
}

func (m *preferencesMap) ensurePreferencesAttached(p fyne.Preferences) *preferenceBindings {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.prefs[p] != nil {
		return m.prefs[p]
	}
	m.prefs[p] = &preferenceBindings{
		items: make(map[string]preferenceItem),
	}
	p.AddChangeListener(func() {
		m.preferencesChanged(p)
	})
	return m.prefs[p]
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

func (m *preferencesMap) getBindings(p fyne.Preferences) *preferenceBindings {
	m.lock.RLock()
	binds := m.prefs[p]
	m.lock.RUnlock()
	return binds
}
