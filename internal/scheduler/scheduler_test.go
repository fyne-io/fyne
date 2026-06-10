package scheduler

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
)

func encodeJSON(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func TestScheduler_Schedule_FiresAtDeliveryTime(t *testing.T) {
	cache := newMemCache()
	var fired atomic.Int32
	done := make(chan struct{}, 1)
	s := New(cache, func(n *fyne.Notification) {
		fired.Add(1)
		done <- struct{}{}
	})
	s.Start()

	id, err := s.Schedule(fyne.NewNotification("hello", "world"), time.Now().Add(50*time.Millisecond))
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for scheduled fire")
	}
	assert.Equal(t, int32(1), fired.Load())
}

func TestScheduler_Schedule_RejectsPastTime(t *testing.T) {
	s := New(newMemCache(), func(*fyne.Notification) {})
	_, err := s.Schedule(fyne.NewNotification("a", "b"), time.Now().Add(-time.Second))
	assert.Error(t, err)
}

func TestScheduler_Cancel_StopsDelivery(t *testing.T) {
	cache := newMemCache()
	var fired atomic.Int32
	s := New(cache, func(*fyne.Notification) { fired.Add(1) })
	s.Start()

	id, err := s.Schedule(fyne.NewNotification("a", "b"), time.Now().Add(100*time.Millisecond))
	require.NoError(t, err)

	s.Cancel(id)
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, int32(0), fired.Load())
	assert.False(t, cache.exists(scheduleFile), "cache file should be removed when no entries remain")
}

func TestScheduler_Persistence_AcrossInstances(t *testing.T) {
	cache := newMemCache()

	first := New(cache, func(*fyne.Notification) {})
	first.Start()
	deliverAt := time.Now().Add(2 * time.Second)
	id, err := first.Schedule(fyne.NewNotification("p", "q"), deliverAt)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	assert.True(t, cache.exists(scheduleFile))

	// Second instance: should pick up the persisted entry and re-arm it.
	var fired atomic.Int32
	done := make(chan *fyne.Notification, 1)
	second := New(cache, func(n *fyne.Notification) {
		fired.Add(1)
		done <- n
	})
	second.Start()

	select {
	case got := <-done:
		assert.Equal(t, "p", got.Title)
		assert.Equal(t, "q", got.Content)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for restored schedule to fire")
	}
}

func TestScheduler_PastEntries_FireImmediatelyOnStart(t *testing.T) {
	cache := newMemCache()

	// Persist an entry directly with a delivery time in the past, simulating
	// an app that scheduled a notification and was killed before delivery.
	pastEntry := &Entry{
		ID:           "manual-id",
		Title:        "late",
		Content:      "load",
		DeliveryTime: time.Now().Add(-time.Hour),
	}
	w, err := cache.Write(scheduleFile)
	require.NoError(t, err)
	require.NoError(t, encodeJSON(w, []*Entry{pastEntry}))
	require.NoError(t, w.Close())

	done := make(chan *fyne.Notification, 1)
	s := New(cache, func(n *fyne.Notification) { done <- n })
	s.Start()

	select {
	case got := <-done:
		assert.Equal(t, "late", got.Title)
		assert.Equal(t, "load", got.Content)
	case <-time.After(2 * time.Second):
		t.Fatal("past-due entry was not delivered on Start")
	}

	// The delivered entry must not remain in the cache, otherwise a subsequent
	// launch would re-deliver it.
	assert.False(t, cache.exists(scheduleFile),
		"past-due entry must be cleared from cache after delivery on Start")
}

// memCache is an in-memory implementation of fyne.Cache for tests.
type memCache struct {
	mu    sync.Mutex
	files map[string][]byte
}

func newMemCache() *memCache {
	return &memCache{files: map[string][]byte{}}
}

func (m *memCache) exists(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.files[name]
	return ok
}

func (m *memCache) RootURI() fyne.URI       { return nil }
func (m *memCache) Exists(name string) bool { return m.exists(name) }

func (m *memCache) Read(name string) (io.ReadCloser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, ok := m.files[name]
	if !ok {
		return nil, io.EOF
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *memCache) Write(name string) (io.WriteCloser, error) {
	return &memWriter{cache: m, name: name}, nil
}

func (m *memCache) Remove(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.files, name)
	return nil
}

type memWriter struct {
	cache *memCache
	name  string
	buf   []byte
}

func (w *memWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func (w *memWriter) Close() error {
	w.cache.mu.Lock()
	defer w.cache.mu.Unlock()
	w.cache.files[w.name] = w.buf
	return nil
}
