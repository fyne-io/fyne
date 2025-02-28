package binding

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

type preferenceItem interface {
	checkForChange()
}

type preferenceBindings struct {
	async.Map[string, preferenceItem]
}

func (b *preferenceBindings) list() []preferenceItem {
	ret := []preferenceItem{}
	b.Range(func(_ string, item preferenceItem) bool {
		ret = append(ret, item)
		return true
	})
	return ret
}

type preferencesMap struct {
	prefs async.Map[fyne.Preferences, *preferenceBindings]

	appPrefs fyne.Preferences // the main application prefs, to check if it changed...
	appLock  sync.Mutex
}

func newPreferencesMap() *preferencesMap {
	return &preferencesMap{}
}

func (m *preferencesMap) ensurePreferencesAttached(p fyne.Preferences) *preferenceBindings {
	binds, loaded := m.prefs.LoadOrStore(p, &preferenceBindings{})
	if loaded {
		return binds
	}

	p.AddChangeListener(func() { m.preferencesChanged(p) })
	return binds
}

func (m *preferencesMap) getBindings(p fyne.Preferences) *preferenceBindings {
	if p == fyne.CurrentApp().Preferences() {
		m.appLock.Lock()
		prefs := m.appPrefs
		if m.appPrefs == nil {
			m.appPrefs = p
		}
		m.appLock.Unlock()
		if prefs != p {
			m.migratePreferences(prefs, p)
		}
	}
	binds, _ := m.prefs.Load(p)
	return binds
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

func (m *preferencesMap) migratePreferences(src, dst fyne.Preferences) {
	old, loaded := m.prefs.Load(src)
	if !loaded {
		return
	}

	m.prefs.Store(dst, old)
	m.prefs.Delete(src)
	m.appLock.Lock()
	m.appPrefs = dst
	m.appLock.Unlock()

	binds := m.getBindings(dst)
	if binds == nil {
		return
	}
	for _, b := range binds.list() {
		if backed, ok := b.(interface{ replaceProvider(fyne.Preferences) }); ok {
			backed.replaceProvider(dst)
		}
	}

	m.preferencesChanged(dst)
}
