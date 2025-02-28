package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunNative(t *testing.T) {
	err := RunNative(func(i any) error {
		native, ok := i.(*UnknownContext)

		assert.True(t, ok)
		assert.NotNil(t, native)

		return nil
	})

	require.NoError(t, err)
}
