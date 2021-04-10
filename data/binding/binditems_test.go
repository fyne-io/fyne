package binding

import (
	"testing"

	"fyne.io/fyne/v2/storage"
	"github.com/stretchr/testify/assert"
)

func TestBindFloat(t *testing.T) {
	val := 0.5
	f := BindFloat(&val)
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.5, v)

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	called = false
	err = f.Set(0.3)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 0.3, val)
	assert.True(t, called)

	called = false
	val = 1.2
	f.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 1.2, v)
}

func TestNewFloat(t *testing.T) {
	f := NewFloat()
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, v)

	err = f.Set(0.3)
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.3, v)
}

func TestBindFloat_TriggerOnlyWhenChange(t *testing.T) {
	v := 9.8
	f := BindFloat(&v)
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := f.Set(9.0)
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := f.Set(9.1)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = f.Set(8.0)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	v = 8.0
	err = f.Reload()
	assert.NoError(t, err)
	waitForItems()
	assert.Equal(t, 0, triggered)

	triggered = 0
	v = 9.2
	err = f.Reload()
	assert.NoError(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)
}

func TestNewFloat_TriggerOnlyWhenChange(t *testing.T) {
	f := NewFloat()
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := f.Set(9.0)
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := f.Set(9.1)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = f.Set(8.0)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)
}

func TestBindURI_TriggerOnlyWhenChange(t *testing.T) {
	v := storage.NewURI("first")
	b := BindURI(&v)
	triggered := 0
	b.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := b.Set(storage.NewURI("second"))
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := b.Set(storage.NewURI("third"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = b.Set(storage.NewURI("fourth"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	v = storage.NewURI("fourth")
	err = b.Reload()
	assert.NoError(t, err)
	waitForItems()
	assert.Equal(t, 0, triggered)

	triggered = 0
	v = storage.NewURI("fifth")
	err = b.Reload()
	assert.NoError(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)
}

func TestNewURI_TriggerOnlyWhenChange(t *testing.T) {
	b := NewURI()
	triggered := 0
	b.AddListener(NewDataListener(func() {
		triggered++
	}))
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := b.Set(storage.NewURI("first"))
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := b.Set(storage.NewURI("second"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = b.Set(storage.NewURI("third"))
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 1, triggered)
}
