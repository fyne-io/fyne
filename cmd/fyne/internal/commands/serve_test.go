package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validatePort(t *testing.T) {
	valid := &Server{port: 80}
	assert.NoError(t, valid.validate())
	assert.Equal(t, 80, valid.port)

	def := &Server{}
	assert.NoError(t, def.validate())
	assert.Equal(t, 8080, def.port)

	invalidLow := &Server{port: -1}
	assert.Error(t, invalidLow.validate())

	invalidHigh := &Server{port: 65536}
	assert.Error(t, invalidHigh.validate())
}
