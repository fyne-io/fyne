package cache

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestTextCacheGet(t *testing.T) {
	ResetThemeCaches()
	assert.Equal(t, 0, len(fontSizeCache))

	bound, base := GetFontMetrics("hi", 10, fyne.TextStyle{}, nil)
	assert.True(t, bound.IsZero())
	assert.Equal(t, float32(0), base)

	SetFontMetrics("hi", 10, fyne.TextStyle{}, nil, fyne.NewSize(10, 10), 8)
	assert.Equal(t, 1, len(fontSizeCache))

	bound, base = GetFontMetrics("hi", 10, fyne.TextStyle{}, nil)
	assert.Equal(t, fyne.NewSize(10, 10), bound)
	assert.Equal(t, float32(8), base)
}
