package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestApp_CloudProvider(t *testing.T) {
	a := NewApp()
	c := &mockCloud{}
	a.SetCloudProvider(c)

	assert.Equal(t, c, a.CloudProvider())
}
