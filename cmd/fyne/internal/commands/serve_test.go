package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validatePort(t *testing.T) {
	valid := &Server{port: 80}
	assert.Nil(t, valid.validate())
	assert.Equal(t, 80, valid.port)

	def := &Server{}
	assert.Nil(t, def.validate())
	assert.Equal(t, 8080, def.port)

	invalidLow := &Server{port: -1}
	assert.NotNil(t, invalidLow.validate())

	invalidHigh := &Server{port: 65536}
	assert.NotNil(t, invalidHigh.validate())
}
