package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Device(t *testing.T) {
	dev := &device{}

	assert.Equal(t, true, dev.IsMobile())
	assert.Equal(t, false, dev.IsBrowser())
	assert.Equal(t, false, dev.HasKeyboard())
}
