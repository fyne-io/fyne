package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	assert.Equal(t, float32(1), min(float32(1), float32(3)))
	assert.Equal(t, float32(-3), min(float32(1), float32(-3)))
}

func TestMax(t *testing.T) {
	assert.Equal(t, float32(3), max(float32(1), float32(3)))
	assert.Equal(t, float32(1), max(float32(1), float32(-3)))
}
