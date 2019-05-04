package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestNewResource(t *testing.T) {
	name := "test"
	content := []byte{1, 2, 3, 4}

	res := NewStaticResource(name, content, true)
	assert.Equal(t, name, res.Name())
	assert.Equal(t, content, res.Content())
	assert.True(t, res.AdoptIconColor())
}
