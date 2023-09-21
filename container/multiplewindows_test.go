package container

import (
	"testing"

	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestMultipleWindows_Add(t *testing.T) {
	m := NewMultipleWindows()
	assert.Zero(t, len(m.Windows))

	m.Add(NewInnerWindow("1", widget.NewLabel("Inside")))
	assert.Equal(t, 1, len(m.Windows))
}
