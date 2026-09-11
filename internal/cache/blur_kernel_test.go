package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlurKernelCache(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	values, ok := GetBlurKernel(1.5)
	assert.False(t, ok)
	assert.Nil(t, values)

	tm := &timeMock{}
	tm.setTime(10, 10)
	SetBlurKernel(1.5, []float32{0.2, 0.6, 0.2})
	values, ok = GetBlurKernel(1.5)
	assert.True(t, ok)
	assert.Equal(t, []float32{0.2, 0.6, 0.2}, values)

	lastClean = time.Time{}
	tm.setTime(11, 20)
	Clean(false)

	_, ok = GetBlurKernel(1.5)
	assert.False(t, ok)
}
