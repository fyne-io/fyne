package cache

import (
	"image/color"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestObjectTextures(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	c := &dummyCanvas{}
	other := &dummyCanvas{}
	obj := &dummyCanvasObject{}
	otherObj := &dummyCanvasObject{}

	tex, ok := GetTexture(obj)
	assert.False(t, ok)
	assert.Equal(t, NoTexture, tex)
	assert.False(t, IsValid(tex))

	tm := &timeMock{}
	tm.setTime(10, 10)
	SetTexture(obj, testTexture(7), c)
	tex, ok = GetTexture(obj)
	assert.True(t, ok)
	assert.Equal(t, testTexture(7), tex)
	assert.True(t, IsValid(tex))

	visited := 0
	RangeTexturesFor(c, func(_ fyne.CanvasObject) { visited++ })
	assert.Equal(t, 1, visited)

	visited = 0
	RangeTexturesFor(other, func(_ fyne.CanvasObject) { visited++ })
	assert.Equal(t, 0, visited)

	DeleteTexture(obj)
	_, ok = GetTexture(obj)
	assert.False(t, ok)

	objectTextures.Store(obj, nil)
	tex, ok = GetTexture(obj)
	assert.False(t, ok)
	assert.Equal(t, NoTexture, tex)

	SetTexture(obj, testTexture(1), c)
	SetTexture(otherObj, testTexture(2), other)
	tm.now = tm.now.Add(ValidDuration + time.Second)

	expired := 0
	RangeExpiredTexturesFor(c, func(_ fyne.CanvasObject) { expired++ })
	assert.Equal(t, 1, expired)

	expired = 0
	RangeExpiredTexturesFor(other, func(_ fyne.CanvasObject) { expired++ })
	assert.Equal(t, 1, expired)
}

func TestTextTextures(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	c := &dummyCanvas{}
	other := &dummyCanvas{}
	ent := FontCacheEntry{
		fontSizeEntry: fontSizeEntry{Text: "hi", Size: 10},
		Canvas:        c,
		Color:         color.Black,
	}
	otherEnt := FontCacheEntry{
		fontSizeEntry: fontSizeEntry{Text: "yo", Size: 10},
		Canvas:        other,
	}

	tex, ok := GetTextTexture(ent)
	assert.False(t, ok)
	assert.Equal(t, NoTexture, tex)

	textTextures.Store(ent, nil)
	tex, ok = GetTextTexture(ent)
	assert.False(t, ok)
	assert.Equal(t, NoTexture, tex)

	freed := 0
	SetTextTexture(ent, testTexture(3), c, func() { freed++ })
	SetTextTexture(otherEnt, testTexture(4), other, func() { freed++ })

	tex, ok = GetTextTexture(ent)
	assert.True(t, ok)
	assert.Equal(t, testTexture(3), tex)

	DeleteTextTexturesFor(c)
	assert.Equal(t, 1, freed)
	_, ok = GetTextTexture(ent)
	assert.False(t, ok)
	_, ok = GetTextTexture(otherEnt)
	assert.True(t, ok)

	tm := &timeMock{}
	tm.setTime(10, 10)
	SetTextTexture(ent, testTexture(5), c, func() { freed++ })
	tm.now = tm.now.Add(ValidDuration + time.Second)
	RangeExpiredTexturesFor(c, func(_ fyne.CanvasObject) {})
	assert.Equal(t, 2, freed)
	_, ok = GetTextTexture(ent)
	assert.False(t, ok)
}
