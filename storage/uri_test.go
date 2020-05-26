package storage

import (
	"testing"

	_ "fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestNewURI(t *testing.T) {
	uriStr := "file:///nothere.txt"
	u := NewURI(uriStr)

	assert.NotNil(t, u)
	assert.Equal(t, uriStr, u.String())
	assert.Equal(t, "file", u.Scheme())
	assert.Equal(t, ".txt", u.Extension())
	assert.Equal(t, "text/plain", u.MimeType())
}

func TestNewURI_Content(t *testing.T) {
	uriStr := "content:///otherapp/something/pic.png"
	u := NewURI(uriStr)

	assert.NotNil(t, u)
	assert.Equal(t, uriStr, u.String())
	assert.Equal(t, "content", u.Scheme())
	assert.Equal(t, ".png", u.Extension())
	assert.Equal(t, "image/png", u.MimeType())
}
