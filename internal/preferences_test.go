package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefs_SetBool(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetBool("testBool", true)

	assert.Equal(t, true, p.Bool("testBool"))
}

func TestPrefs_Bool(t *testing.T) {
	p := NewInMemoryPreferences()
	p.Values["testBool"] = true

	assert.Equal(t, true, p.Bool("testBool"))
}

func TestPrefs_Bool_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, false, p.Bool("testBool"))
}

func TestPrefs_SetFloat(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetFloat("testFloat", 1.7)

	assert.Equal(t, 1.7, p.Float("testFloat"))
}

func TestPrefs_Float(t *testing.T) {
	p := NewInMemoryPreferences()
	p.Values["testFloat"] = 1.2

	assert.Equal(t, 1.2, p.Float("testFloat"))
}

func TestPrefs_Float_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 0.0, p.Float("testFloat"))
}

func TestPrefs_SetInt(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetInt("testInt", 5)

	assert.Equal(t, 5, p.Int("testInt"))
}

func TestPrefs_Int(t *testing.T) {
	p := NewInMemoryPreferences()
	p.Values["testInt"] = 5

	assert.Equal(t, 5, p.Int("testInt"))
}

func TestPrefs_Int_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 0, p.Int("testInt"))
}

func TestPrefs_SetString(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetString("test", "value")

	assert.Equal(t, "value", p.String("test"))
}

func TestPrefs_String(t *testing.T) {
	p := NewInMemoryPreferences()
	p.Values["test"] = "value"

	assert.Equal(t, "value", p.String("test"))
}

func TestPrefs_String_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, "", p.String("test"))
}

func TestInMemoryPreferences_OnChange(t *testing.T) {
	p := NewInMemoryPreferences()
	called := false
	p.OnChange = func() {
		called = true
	}

	p.SetString("dummy", "another")

	assert.True(t, called)
}
