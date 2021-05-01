package binding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
}

func newTime() *timeBinding {
	return &timeBinding{Int: NewInt()}
}

func (t *timeBinding) Get() (time.Time, error) {
	i, err := t.Int.Get()
	return time.Unix(int64(i), 0), err
}

func (t *timeBinding) Set(time time.Time) error {
	i := time.Unix()
	return t.Int.Set(int(i))
}
