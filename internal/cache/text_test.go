package cache

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestTextCacheGet(t *testing.T) {
	ResetThemeCaches()
	assert.Equal(t, 0, len(sizeCache))

	bound, base := GetFontSize("hi", 10, fyne.TextStyle{})
	assert.True(t, bound.IsZero())
	assert.Equal(t, float32(0), base)

	SetFontSize("hi", 10, fyne.TextStyle{}, fyne.NewSize(10, 10), 8)
	assert.Equal(t, 1, len(sizeCache))

	bound, base = GetFontSize("hi", 10, fyne.TextStyle{})
	assert.Equal(t, fyne.NewSize(10, 10), bound)
	assert.Equal(t, float32(8), base)
}
