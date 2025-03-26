package binding

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2/storage"
)

func init() {
	runtime.LockOSThread()
}

func TestBindFloat(t *testing.T) {
	val := 0.5
	f := BindFloat(&val)
	v, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 0.5, v)

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	assert.True(t, called)

	called = false
	err = f.Set(0.3)
	require.NoError(t, err)
	assert.Equal(t, 0.3, val)
	assert.True(t, called)

	called = false
	val = 1.2
	f.Reload()
	assert.True(t, called)
	v, err = f.Get()
	require.NoError(t, err)
	assert.Equal(t, 1.2, v)
}

func TestBindUntyped(t *testing.T) {
	var val any
	val = 0.5
	f := BindUntyped(&val)
	v, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 0.5, v)

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	assert.True(t, called)

	called = false
	err = f.Set(0.3)
	require.NoError(t, err)
	assert.Equal(t, 0.3, val)
	assert.True(t, called)

	called = false
	val = 1.2
	f.Reload()
	assert.True(t, called)
	v, err = f.Get()
	require.NoError(t, err)
	assert.Equal(t, 1.2, v)
}

func TestNewFloat(t *testing.T) {
	f := NewFloat()
	v, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, 0.0, v)

	err = f.Set(0.3)
	require.NoError(t, err)
	v, err = f.Get()
	require.NoError(t, err)
	assert.Equal(t, 0.3, v)
}

func TestBindFloat_TriggerOnlyWhenChange(t *testing.T) {
	v := 9.8
	f := BindFloat(&v)
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := f.Set(9.0)
		require.NoError(t, err)
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := f.Set(9.1)
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = f.Set(8.0)
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)

	triggered = 0
	v = 8.0
	err = f.Reload()
	require.NoError(t, err)
	assert.Equal(t, 0, triggered)

	triggered = 0
	v = 9.2
	err = f.Reload()
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)
}

func TestNewFloat_TriggerOnlyWhenChange(t *testing.T) {
	f := NewFloat()
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := f.Set(9.0)
		require.NoError(t, err)
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := f.Set(9.1)
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = f.Set(8.0)
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)
}

func TestBindURI(t *testing.T) {
	val := storage.NewFileURI("/tmp")
	f := BindURI(&val)
	v, err := f.Get()
	require.NoError(t, err)
	assert.Equal(t, "file:///tmp", v.String())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	assert.True(t, called)

	called = false
	err = f.Set(storage.NewFileURI("/tmp/test.txt"))
	require.NoError(t, err)
	assert.Equal(t, "file:///tmp/test.txt", val.String())
	assert.True(t, called)

	called = false
	val = storage.NewFileURI("/hello")
	_ = f.Reload()
	assert.True(t, called)
	v, err = f.Get()
	require.NoError(t, err)
	assert.Equal(t, "file:///hello", v.String())
}

func TestBindURI_TriggerOnlyWhenChange(t *testing.T) {
	v := storage.NewFileURI("first")
	b := BindURI(&v)
	triggered := 0
	b.AddListener(NewDataListener(func() {
		triggered++
	}))
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := b.Set(storage.NewFileURI("second"))
		require.NoError(t, err)
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := b.Set(storage.NewFileURI("third"))
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = b.Set(storage.NewFileURI("fourth"))
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)

	triggered = 0
	v = storage.NewFileURI("fourth")
	err = b.Reload()
	require.NoError(t, err)
	assert.Equal(t, 0, triggered)

	triggered = 0
	v = storage.NewFileURI("fifth")
	err = b.Reload()
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)
}

func TestNewURI(t *testing.T) {
	f := NewURI()
	v, err := f.Get()
	require.NoError(t, err)
	assert.Nil(t, v)

	err = f.Set(storage.NewFileURI("/var"))
	require.NoError(t, err)
	v, err = f.Get()
	require.NoError(t, err)
	assert.Equal(t, "file:///var", v.String())
}

func TestNewURI_TriggerOnlyWhenChange(t *testing.T) {
	b := NewURI()
	triggered := 0
	b.AddListener(NewDataListener(func() {
		triggered++
	}))
	assert.Equal(t, 1, triggered)

	triggered = 0
	for i := 0; i < 5; i++ {
		err := b.Set(storage.NewFileURI("first"))
		require.NoError(t, err)
		assert.Equal(t, 1, triggered)
	}

	triggered = 0
	err := b.Set(storage.NewFileURI("second"))
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)

	triggered = 0
	err = b.Set(storage.NewFileURI("third"))
	require.NoError(t, err)
	assert.Equal(t, 1, triggered)
}
