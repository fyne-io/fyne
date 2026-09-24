package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestTextCacheGet(t *testing.T) {
	ResetThemeCaches()
	assert.Equal(t, 0, fontSizeCache.Len())

	bound, base := GetFontMetrics("hi", 10, fyne.TextStyle{}, nil)
	assert.True(t, bound.IsZero())
	assert.Equal(t, float32(0), base)

	SetFontMetrics("hi", 10, fyne.TextStyle{}, nil, fyne.NewSize(10, 10), 8)
	assert.Equal(t, 1, fontSizeCache.Len())

	bound, base = GetFontMetrics("hi", 10, fyne.TextStyle{}, nil)
	assert.Equal(t, fyne.NewSize(10, 10), bound)
	assert.Equal(t, float32(8), base)
}

func TestFontMetricsWithSourceAndExpiry(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	res := fyne.NewStaticResource("font.ttf", []byte("x"))
	SetFontMetrics("hi", 10, fyne.TextStyle{}, res, fyne.NewSize(12, 12), 7)
	bound, base := GetFontMetrics("hi", 10, fyne.TextStyle{}, res)
	assert.Equal(t, fyne.NewSize(12, 12), bound)
	assert.Equal(t, float32(7), base)

	ClearFontMetrics()
	bound, base = GetFontMetrics("hi", 10, fyne.TextStyle{}, res)
	assert.True(t, bound.IsZero())
	assert.Equal(t, float32(0), base)

	tm := &timeMock{}
	tm.setTime(10, 10)
	SetFontMetrics("hi", 10, fyne.TextStyle{}, res, fyne.NewSize(12, 12), 7)
	lastClean = time.Time{}
	tm.setTime(11, 20)
	Clean(false)
	bound, _ = GetFontMetrics("hi", 10, fyne.TextStyle{}, res)
	assert.True(t, bound.IsZero())
}
