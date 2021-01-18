package dialog

import (
	"testing"

	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestListDrives(t *testing.T) {
	driveLetters := listDrives()
	assert.GreaterOrEqual(t, len(driveLetters), 1)

	if driveLetters[0] != "C:" && driveLetters[0] != "A:" {
		t.Error("Unexpected lowest drive letter, " + driveLetters[0])
	}
}

func TestFileDialog_LoadPlaces(t *testing.T) {
	f := &fileDialog{}
	driveLetters := listDrives()
	places := f.loadPlaces()

	assert.Equal(t, len(driveLetters), len(places))
	assert.Equal(t, driveLetters[0], places[0].(*widget.Button).Text)
}
