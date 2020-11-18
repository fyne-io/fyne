package binding

import (
	"testing"
	"time"

	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestBindPreferenceBool(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-bool"

	p.SetBool(key, true)
	bind := BindPreferenceBool(key, p)
	assert.True(t, bind.Get())

	bind.Set(false)
	assert.False(t, bind.Get())
	assert.False(t, p.Bool(key))
}

func TestBindPreferenceFloat(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-float"

	p.SetFloat(key, 1.3)
	bind := BindPreferenceFloat(key, p)
	assert.Equal(t, 1.3, bind.Get())

	bind.Set(2.5)
	assert.Equal(t, 2.5, bind.Get())
	assert.Equal(t, 2.5, p.Float(key))
}

func TestBindPreferenceInt(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-int"

	p.SetInt(key, 4)
	bind := BindPreferenceInt(key, p)
	assert.Equal(t, 4, bind.Get())

	bind.Set(7)
	assert.Equal(t, 7, bind.Get())
	assert.Equal(t, 7, p.Int(key))
}

func TestBindPreferenceString(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind := BindPreferenceString(key, p)
	assert.Equal(t, "aString", bind.Get())

	bind.Set("overwritten")
	assert.Equal(t, "overwritten", bind.Get())
	assert.Equal(t, "overwritten", p.String(key))
}

func TestPreferenceBindingCopies(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind1 := BindPreferenceString(key, p)
	bind2 := BindPreferenceString(key, p)
	assert.Equal(t, bind2.Get(), bind1.Get())

	bind1.Set("overwritten")
	assert.Equal(t, bind2.Get(), bind1.Get())
}

func TestPreferenceBindingTriggers(t *testing.T) {
	a := test.NewApp()
	p := a.Preferences()
	key := "test-string"

	p.SetString(key, "aString")
	bind1 := BindPreferenceString(key, p)
	bind2 := BindPreferenceString(key, p)

	ch := make(chan interface{})
	bind1.AddListener(NewDataListener(func() {
		ch <- struct{}{}
	}))

	select {
	case <-ch: // bind1 gets initial value
	case <-time.After(time.Millisecond * 100):
		t.Errorf("Timed out waiting for data binding to send initial value")
	}

	bind2.Set("overwritten") // write on a different listener, preferences should trigger all
	select {
	case <-ch: // bind1 triggered by bind2 changing the same key
	case <-time.After(time.Millisecond * 100):
		t.Errorf("Timed out waiting for data binding change to trigger")
	}
}
