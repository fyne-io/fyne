package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunNative(t *testing.T) {
	err := RunNative(func(i interface{}) error {
		native, ok := i.(*UnknownContext)

		assert.True(t, ok)
		assert.NotNil(t, native)

		return nil
	})

	assert.Nil(t, err)
}
