// Package scheduler provides a persistent in-process notification scheduler used
// as the fallback for platforms without a native scheduling API.
//
// Schedules are stored as JSON in the application Cache so that pending entries
// survive an app restart. When the scheduler is started any entries whose
// delivery time has already passed are fired immediately; future entries are
// re-armed with a [time.Timer].
package scheduler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// scheduleFile is the cache entry name used to persist pending schedules.
const scheduleFile = "fyne-scheduled-notifications.json"

// FireFunc delivers a notification. It is invoked by the scheduler when a timer
// fires or when a stored schedule's delivery time has passed at startup.
type FireFunc func(n *fyne.Notification)

// Entry is the persisted record for a single pending notification.
type Entry struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	DeliveryTime time.Time `json:"deliveryTime"`
}

// Scheduler manages a set of pending in-process scheduled notifications, persisting
// them to a [fyne.Cache] so that schedules survive an app restart.
type Scheduler struct {
	cache fyne.Cache
	fire  FireFunc

	mu      sync.Mutex
	entries map[string]*Entry
	timers  map[string]*time.Timer
	started bool
}

// New returns a new Scheduler that persists to the given cache and delivers via fire.
// Call [Scheduler.Start] once the application is ready to deliver notifications.
func New(cache fyne.Cache, fire FireFunc) *Scheduler {
	return &Scheduler{
		cache:   cache,
		fire:    fire,
		entries: map[string]*Entry{},
		timers:  map[string]*time.Timer{},
	}
}

// Start loads any persisted entries and arms timers for those still in the future.
// Past-due entries are delivered immediately and removed from the persisted list
// so they do not refire on the next launch. Calling Start more than once is a no-op.
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return
	}
	s.started = true
	loaded := s.loadLocked()
	now := time.Now()
	var due []*Entry
	for _, e := range loaded {
		if !e.DeliveryTime.After(now) {
			due = append(due, e)
			continue
		}
		s.entries[e.ID] = e
		s.armLocked(e)
	}
	if len(due) > 0 {
		_ = s.saveLocked()
	}
	s.mu.Unlock()

	for _, e := range due {
		s.deliver(e)
	}
}

// Schedule queues a notification for delivery at the given time, returning the
// generated ID. An error is returned if when is in the past or persistence fails.
func (s *Scheduler) Schedule(n *fyne.Notification, when time.Time) (string, error) {
	if n == nil {
		return "", errors.New("nil notification")
	}
	if !when.After(time.Now()) {
		return "", errors.New("scheduled delivery time must be in the future")
	}

	id, err := NewID()
	if err != nil {
		return "", err
	}

	entry := &Entry{
		ID:           id,
		Title:        n.Title,
		Content:      n.Content,
		DeliveryTime: when,
	}

	s.mu.Lock()
	s.entries[id] = entry
	if s.started {
		s.armLocked(entry)
	}
	if err := s.saveLocked(); err != nil {
		delete(s.entries, id)
		if t, ok := s.timers[id]; ok {
			t.Stop()
			delete(s.timers, id)
		}
		s.mu.Unlock()
		return "", err
	}
	s.mu.Unlock()

	return id, nil
}

// Cancel removes a scheduled notification by ID. It is safe to call with an
// unknown ID or after the timer has already fired. An error is returned if the
// updated entry list could not be persisted.
func (s *Scheduler) Cancel(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.timers[id]; ok {
		t.Stop()
		delete(s.timers, id)
	}
	if _, ok := s.entries[id]; ok {
		delete(s.entries, id)
		return s.saveLocked()
	}
	return nil
}

// armLocked installs a timer for entry. Must be called with s.mu held.
func (s *Scheduler) armLocked(entry *Entry) {
	delay := time.Until(entry.DeliveryTime)
	if delay < 0 {
		delay = 0
	}
	id := entry.ID
	s.timers[id] = time.AfterFunc(delay, func() {
		s.mu.Lock()
		current, ok := s.entries[id]
		if ok {
			delete(s.entries, id)
		}
		delete(s.timers, id)
		_ = s.saveLocked()
		s.mu.Unlock()

		if ok {
			s.deliver(current)
		}
	})
}

// deliver invokes the fire function with a Notification built from entry.
func (s *Scheduler) deliver(e *Entry) {
	if s.fire == nil {
		return
	}
	s.fire(&fyne.Notification{Title: e.Title, Content: e.Content})
}

// saveLocked writes the current entries map to the cache. Must be called with s.mu held.
func (s *Scheduler) saveLocked() error {
	if s.cache == nil {
		return nil
	}
	list := make([]*Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	if len(list) == 0 {
		_ = s.cache.Remove(scheduleFile)
		return nil
	}
	w, err := s.cache.Write(scheduleFile)
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(list)
}

// loadLocked reads persisted entries from the cache. Must be called with s.mu held.
func (s *Scheduler) loadLocked() []*Entry {
	if s.cache == nil || !s.cache.Exists(scheduleFile) {
		return nil
	}
	r, err := s.cache.Read(scheduleFile)
	if err != nil {
		return nil
	}
	defer r.Close()

	var list []*Entry
	if err = json.NewDecoder(r).Decode(&list); err != nil {
		fyne.LogError("Failed to read scheduled notifications", err)
		return nil
	}
	return list
}

// NewID generates a random hex identifier suitable for a scheduled notification.
//
// Since: 2.8
func NewID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "fyne-sched-" + hex.EncodeToString(b[:]), nil
}
