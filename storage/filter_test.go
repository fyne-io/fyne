package storage_test

import (
	"testing"

	"fyne.io/fyne/v2/storage"

	_ "fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestFileFilter(t *testing.T) {
	filter := storage.NewExtensionFileFilter([]string{".jpg", ".png"})

	assert.NotNil(t, filter)
	assert.Equal(t, true, filter.Matches(storage.NewURI("content:///otherapp/something/pic.JPG")))
	assert.Equal(t, true, filter.Matches(storage.NewURI("content:///otherapp/something/pic.jpg")))

	assert.Equal(t, true, filter.Matches(storage.NewURI("content:///otherapp/something/pic.PNG")))
	assert.Equal(t, true, filter.Matches(storage.NewURI("content:///otherapp/something/pic.png")))

	assert.Equal(t, false, filter.Matches(storage.NewURI("content:///otherapp/something/pic.TIFF")))
	assert.Equal(t, false, filter.Matches(storage.NewURI("content:///otherapp/something/pic.tiff")))
}
