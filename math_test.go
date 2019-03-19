package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	assert.Equal(t, Min(1, 3), 1)
	assert.Equal(t, Min(1, -3), -3)
}

func TestMax(t *testing.T) {
	assert.Equal(t, Max(1, 3), 3)
	assert.Equal(t, Max(1, -3), 1)
}
