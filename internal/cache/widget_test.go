package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCachedRenderer(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	r, ok := CachedRenderer(nil)
	assert.False(t, ok)
	assert.Nil(t, r)

	inner := &dummyWidget{}
	want := Renderer(inner)

	outer := &dummyBaseWidget{impl: inner}
	r, ok = CachedRenderer(outer)
	assert.True(t, ok)
	assert.Equal(t, want, r)

	unwired := &dummyBaseWidget{}
	r, ok = CachedRenderer(unwired)
	assert.False(t, ok)
	assert.Nil(t, r)
}

func TestDestroyRendererAndIsRendered(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	destroyed := 0
	w := &dummyWidget{onDestroy: func() { destroyed++ }}
	assert.False(t, IsRendered(w))

	DestroyRenderer(w)
	assert.Zero(t, destroyed)

	Renderer(w)
	assert.True(t, IsRendered(w))

	DestroyRenderer(w)
	assert.Equal(t, 1, destroyed)
	assert.False(t, IsRendered(w))

	renderers.Store(w, nil)
	DestroyRenderer(w)
	assert.False(t, IsRendered(w))
}

func TestRendererNilWidget(t *testing.T) {
	assert.Nil(t, Renderer(nil))
}
