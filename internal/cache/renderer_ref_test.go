package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRendererFor(t *testing.T) {
	wid := &dummyWidget{}
	var ref RendererRef
	r := RendererFor(&ref, wid)
	assert.NotNil(t, r)
	assert.Same(t, r, RendererFor(&ref, wid), "the renderer kept is used again")
	cached, ok := CachedRenderer(wid)
	assert.True(t, ok)
	assert.Same(t, r, cached)

	DestroyRenderer(wid)
	again := RendererFor(&ref, wid)
	assert.NotSame(t, r, again, "destroyed, the widget gets a new renderer")
	cached, _ = CachedRenderer(wid)
	assert.Same(t, again, cached)

	ref.info.expires.Store(0)
	destroyExpiredRenderers(time.Now())
	assert.NotSame(t, again, RendererFor(&ref, wid), "expired and cleaned out, so is a new one")
}
