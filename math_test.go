package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	assert.Equal(t, float32(1), Min(1, 3))
	assert.Equal(t, float32(-3), Min(1, -3))
}

func TestMax(t *testing.T) {
	assert.Equal(t, float32(3), Max(1, 3))
	assert.Equal(t, float32(1), Max(1, -3))
}
