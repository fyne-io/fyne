package binding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBindTime(t *testing.T) {
	val := time.Now()
	f := bindTime(&val)
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, val.Unix(), v.Unix())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	newTime := val.Add(time.Hour)
	called = false
	err = f.Set(newTime)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, newTime.Unix(), val.Unix())
	assert.True(t, called)

	newTime = newTime.Add(time.Minute)
	called = false
	val = newTime
	_ = f.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, newTime.Unix(), v.Unix())
}

func TestNewTime(t *testing.T) {
	f := newTime()
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, time.Unix(0, 0), v)

	now := time.Now()
	err = f.Set(now)
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, now.Unix(), v.Unix())
}

type timeBinding struct {
	Int
	src *time.Time
}

func bindTime(t *time.Time) *timeBinding {
	return &timeBinding{Int: NewInt(), src: t}
}

func newTime() *timeBinding {
	return &timeBinding{Int: NewInt()}
}

func (t *timeBinding) Get() (time.Time, error) {
	if t.src != nil {
		return *t.src, nil
	}

	i, err := t.Int.Get()
	return time.Unix(int64(i), 0), err
}

func (t *timeBinding) Reload() error {
	if t.src == nil {
		return nil
	}

	return t.Set(*t.src)
}

func (t *timeBinding) Set(time time.Time) error {
	if t.src != nil {
		*t.src = time
	}

	i := time.Unix()
	return t.Int.Set(int(i))
}
