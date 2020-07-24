package storage_test

import (
	"testing"

	"fyne.io/fyne/storage"

	_ "fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestNewURI(t *testing.T) {
	uriStr := "file:///nothere.txt"
	u := storage.NewURI(uriStr)

	assert.NotNil(t, u)
	assert.Equal(t, uriStr, u.String())
	assert.Equal(t, "file", u.Scheme())
	assert.Equal(t, ".txt", u.Extension())
	assert.Equal(t, "text/plain", u.MimeType())
}

func TestNewURI_Content(t *testing.T) {
	uriStr := "content:///otherapp/something/pic.png"
	u := storage.NewURI(uriStr)

	assert.NotNil(t, u)
	assert.Equal(t, uriStr, u.String())
	assert.Equal(t, "content", u.Scheme())
	assert.Equal(t, ".png", u.Extension())
	assert.Equal(t, "image/png", u.MimeType())
}

func TestURI_Scheme(t *testing.T) {
	assert.Equal(t, "http", storage.NewURI("http://google.com").Scheme())
	assert.Equal(t, "http", storage.NewURI("hTtP://google.com").Scheme())
	assert.Equal(t, "file", storage.NewURI("file:///home/user/file.txt").Scheme())
	assert.Equal(t, "file", storage.NewURI("FILE://C:/autoexec.bat").Scheme())
}

func TestURI_Extension(t *testing.T) {
	assert.Equal(t, ".txt", storage.NewURI("file:///text.txt").Extension())
	assert.Equal(t, ".txt", storage.NewURI("file:///text.tXt").Extension())
	assert.Equal(t, ".a", storage.NewURI("file:///library.A").Extension())
	assert.Equal(t, ".jpeg", storage.NewURI("file://C:/image.JPEG").Extension())
}
