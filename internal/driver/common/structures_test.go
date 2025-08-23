package common

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/layout"
	"github.com/stretchr/testify/assert"
)

func TestDeduplicatedObjectQueue(t *testing.T) {
	queue := deduplicatedObjectQueue{queue: async.NewCanvasObjectQueue()}
	assert.Zero(t, queue.Len())
	assert.Nil(t, queue.Out())

	obj1 := &fyne.Container{Layout: layout.NewCenterLayout()}
	queue.In(obj1)
	assert.Equal(t, uint64(1), queue.Len())
	queue.In(obj1)
	assert.Equal(t, uint64(1), queue.Len())

	obj2 := &fyne.Container{Layout: layout.NewFormLayout()}
	queue.In(obj2)
	assert.Equal(t, uint64(2), queue.Len())

	assert.Equal(t, obj1, queue.Out())
	assert.Equal(t, obj2, queue.Out())
	assert.Equal(t, nil, queue.Out())
	assert.Zero(t, queue.Len())

	queue.In(obj1)
	assert.Equal(t, uint64(1), queue.Len())
	queue.In(obj1)
	assert.Equal(t, uint64(1), queue.Len())
}
