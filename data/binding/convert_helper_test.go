package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/storage"
)

func TestURIFromStringHelper(t *testing.T) {
	str := "file:///tmp/test.txt"
	u, err := uriFromString(str)

	assert.Nil(t, err)
	assert.Equal(t, str, u.String())
}

func TestURIToStringHelper(t *testing.T) {
	u := storage.NewFileURI("/tmp/test.txt")
	str, err := uriToString(u)

	assert.Nil(t, err)
	assert.Equal(t, u.String(), str)

	str, err = uriToString(nil)
	assert.Nil(t, err)
	assert.Equal(t, "", str)
}
