package shaping

import (
	"github.com/go-text/typesetting/harfbuzz"
	"github.com/go-text/typesetting/opentype/api/font"
)

// fontEntry holds a single key-value pair for an LRU cache.
type fontEntry struct {
	next, prev *fontEntry
	key        *font.Font
	v          *harfbuzz.Font
}

// fontLRU is a least-recently-used cache for harfbuzz fonts built from
// font.Fonts. It uses a doubly-linked list to track how recently elements have
// been used and a map to store element data for quick access.
type fontLRU struct {
	// This implementation is derived from the one here under the terms of the UNLICENSE:
	//
	// https://git.sr.ht/~eliasnaur/gio/tree/e768fe347a732056031100f2c66987d6db258ea4/item/text/lru.go
	m          map[*font.Font]*fontEntry
	head, tail *fontEntry
	maxSize    int
}

// Get fetches the value associated with the given key, if any.
func (l *fontLRU) Get(k *font.Font) (*harfbuzz.Font, bool) {
	if lt, ok := l.m[k]; ok {
		l.remove(lt)
		l.insert(lt)
		return lt.v, true
	}
	return nil, false
}

// Put inserts the given value with the given key, evicting old
// cache entries if necessary.
func (l *fontLRU) Put(k *font.Font, v *harfbuzz.Font) {
	if l.m == nil {
		l.m = make(map[*font.Font]*fontEntry)
		l.head = new(fontEntry)
		l.tail = new(fontEntry)
		l.head.prev = l.tail
		l.tail.next = l.head
	}
	val := &fontEntry{key: k, v: v}
	l.m[k] = val
	l.insert(val)
	if len(l.m) > l.maxSize {
		oldest := l.tail.next
		l.remove(oldest)
		delete(l.m, oldest.key)
	}
}

// remove cuts e out of the lru linked list.
func (l *fontLRU) remove(e *fontEntry) {
	e.next.prev = e.prev
	e.prev.next = e.next
}

// insert adds e to the lru linked list.
func (l *fontLRU) insert(e *fontEntry) {
	e.next = l.head
	e.prev = l.head.prev
	e.prev.next = e
	e.next.prev = e
}
