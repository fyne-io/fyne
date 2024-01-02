//go:build !windows

package dialog

import (
	"testing"

	"fyne.io/fyne/v2/storage"
	"github.com/stretchr/testify/assert"
)

func TestIsHidden(t *testing.T) {
	assert.True(t, isHidden(storage.NewFileURI(".nothing")))
	assert.False(t, isHidden(storage.NewFileURI("file.txt")))
}

func TestFileDialog_LoadPlaces(t *testing.T) {
	f := &fileDialog{}
	places := f.getPlaces()

	assert.Equal(t, 1, len(places))
	assert.Equal(t, "Computer", places[0].locName)
}
