package fyne

import "testing"

import "github.com/stretchr/testify/assert"

// TestGetMissingCanvas checks for the nil case as no canvas was preset.
// What we want to do is check for the right canvas but we can't call the test
// package that provides this due to import loops :-/
func TestGetMissingCanvas(t *testing.T) {
	box := new(dummyObject)

	setDriver(new(dummyDriver))
	assert.Equal(t, nil, GetCanvas(box))
}
