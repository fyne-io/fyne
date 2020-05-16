package layout_test

import (
	"testing"

	internalLayout "fyne.io/fyne/internal/layout"
	"fyne.io/fyne/layout"

	"github.com/stretchr/testify/assert"
)

func TestNewHBoxLayout(t *testing.T) {
	l := layout.NewHBoxLayout()
	if assert.IsType(t, &internalLayout.Box{}, l) {
		b := l.(*internalLayout.Box)
		assert.True(t, b.Horizontal)
	}
}

func TestNewVBoxLayout(t *testing.T) {
	l := layout.NewVBoxLayout()
	if assert.IsType(t, &internalLayout.Box{}, l) {
		b := l.(*internalLayout.Box)
		assert.False(t, b.Horizontal)
	}
}
