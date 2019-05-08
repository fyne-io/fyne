package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSelect(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})

	assert.Equal(t, 2, len(combo.Options))
	assert.Equal(t, "", combo.Selected)
}

func TestSelect_SetSelected(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("2")

	assert.Equal(t, "2", combo.Selected)
}

func TestSelect_SetSelected_Invalid(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("3")

	assert.Equal(t, "", combo.Selected)
}

func TestSelect_SetSelected_InvalidReplace(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("2")
	combo.SetSelected("3")

	assert.Equal(t, "2", combo.Selected)
}

func TestSelect_SetSelected_Callback(t *testing.T) {
	selected := ""
	combo := NewSelect([]string{"1", "2"}, func(s string) {
		selected = s
	})
	combo.SetSelected("2")

	assert.Equal(t, "2", selected)
}
