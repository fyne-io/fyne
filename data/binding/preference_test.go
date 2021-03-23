package binding

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestBindPreferenceBool(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-bool"

	p.SetBool(key, true)
	bind := BindPreferenceBool(key, p)
	v, err := bind.Get()
	assert.Nil(t, err)
	assert.True(t, v)

	err = bind.Set(false)
	assert.Nil(t, err)
	v, err = bind.Get()
	assert.Nil(t, err)
	assert.False(t, v)
	assert.False(t, p.Bool(key))
}

func TestBindPreferenceFloat(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-float"

	p.SetFloat(key, 1.3)
	bind := BindPreferenceFloat(key, p)
	v, err := bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, 1.3, v)

	err = bind.Set(2.5)
	assert.Nil(t, err)
	v, err = bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, 2.5, v)
	assert.Equal(t, 2.5, p.Float(key))
}

func TestBindPreferenceInt(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-int"

	p.SetInt(key, 4)
	bind := BindPreferenceInt(key, p)
	v, err := bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, 4, v)

	err = bind.Set(7)
	assert.Nil(t, err)
	v, err = bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, 7, v)
	assert.Equal(t, 7, p.Int(key))
}

func TestBindPreferenceString(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind := BindPreferenceString(key, p)
	v, err := bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, "aString", v)

	err = bind.Set("overwritten")
	assert.Nil(t, err)
	v, err = bind.Get()
	assert.Nil(t, err)
	assert.Equal(t, "overwritten", v)
	assert.Equal(t, "overwritten", p.String(key))
}

func TestPreferenceBindingCopies(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind1 := BindPreferenceString(key, p)
	bind2 := BindPreferenceString(key, p)
	v1, err := bind1.Get()
	assert.Nil(t, err)
	v2, err := bind2.Get()
	assert.Nil(t, err)
	assert.Equal(t, v2, v1)

	err = bind1.Set("overwritten")
	assert.Nil(t, err)
	v1, err = bind1.Get()
	assert.Nil(t, err)
	v2, err = bind2.Get()
	assert.Nil(t, err)
	assert.Equal(t, v2, v1)
}

func TestPreferenceBindingTriggers(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind1 := BindPreferenceString(key, p)
	bind2 := BindPreferenceString(key, p)

	ch := make(chan interface{}, 2)
	bind1.AddListener(NewDataListener(func() {
		ch <- struct{}{}
	}))

	select {
	case <-ch: // bind1 gets initial value
	case <-time.After(time.Millisecond * 100):
		t.Errorf("Timed out waiting for data binding to send initial value")
	}

	err := bind2.Set("overwritten") // write on a different listener, preferences should trigger all
	assert.Nil(t, err)
	select {
	case <-ch: // bind1 triggered by bind2 changing the same key
	case <-time.After(time.Millisecond * 100):
		t.Errorf("Timed out waiting for data binding change to trigger")
	}

	p.SetString(key, "overwritten2") // changing preference should trigger as well
	select {
	case <-ch: // bind1 triggered by preferences changing the same key directly
	case <-time.After(time.Millisecond * 300):
		t.Errorf("Timed out waiting for data binding change to trigger")
	}
}
