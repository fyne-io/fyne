package glfw

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
)

func TestAnimationCurve(t *testing.T) {
	assert.Equal(t, float32(0.25), animationCurve(0.25, fyne.AnimationLinear))
	assert.True(t, 0.25 > animationCurve(0.25, fyne.AnimationEaseInOut))
}

func TestAnimationEaseInOut(t *testing.T) {
	assert.Equal(t, float32(0.125), animationEaseInOut(0.25))
	assert.Equal(t, float32(0.5), animationEaseInOut(0.5))
	assert.Equal(t, float32(0.875), animationEaseInOut(0.75))
}
