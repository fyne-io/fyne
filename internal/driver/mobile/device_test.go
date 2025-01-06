package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Device(t *testing.T) {
	dev := &device{}

	assert.True(t, dev.IsMobile())
	assert.False(t, dev.IsBrowser())
	assert.False(t, dev.HasKeyboard())
}
