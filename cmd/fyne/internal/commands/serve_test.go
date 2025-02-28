package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_validatePort(t *testing.T) {
	valid := &Server{port: 80}
	require.NoError(t, valid.validate())
	assert.Equal(t, 80, valid.port)

	def := &Server{}
	require.NoError(t, def.validate())
	assert.Equal(t, 8080, def.port)

	invalidLow := &Server{port: -1}
	require.Error(t, invalidLow.validate())

	invalidHigh := &Server{port: 65536}
	require.Error(t, invalidHigh.validate())
}
