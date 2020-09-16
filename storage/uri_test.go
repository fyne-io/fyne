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

func TestURI_Name(t *testing.T) {
	assert.Equal(t, "text.txt", storage.NewURI("file:///text.txt").Name())
	assert.Equal(t, "library.a", storage.NewURI("file:///library.a").Name())
	assert.Equal(t, "image.JPEG", storage.NewURI("file://C:/image.JPEG").Name())
	assert.Equal(t, "directory", storage.NewURI("file://C:/directory/").Name())
}

func TestURI_Parent(t *testing.T) {
	// note the trailing slashes are significant, as they tend to belie a
	// directory

	parent, err := storage.Parent(storage.NewURI("file:///foo/bar/baz"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewURI("file:///foo/bar/baz/"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewURI("file://C:/foo/bar/baz/"))
	assert.Nil(t, err)
	assert.Equal(t, "file://C:/foo/bar/", parent.String())

	_, err = storage.Parent(storage.NewURI("file:///"))
	assert.Equal(t, storage.URIRootError, err)

	// This should cause an error, since this is a Windows-style path and
	// thus we can't get the parent of a drive letter.
	_, err = storage.Parent(storage.NewURI("file://C:/"))
	assert.Equal(t, storage.URIRootError, err)

	// Windows supports UNIX-style paths. /C:/ is also a valid path.
	parent, err = storage.Parent(storage.NewURI("file:///C:/"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///", parent.String())
}

func TestURI_Extension(t *testing.T) {
	assert.Equal(t, ".txt", storage.NewURI("file:///text.txt").Extension())
	assert.Equal(t, ".a", storage.NewURI("file:///library.a").Extension())
	assert.Equal(t, ".JPEG", storage.NewURI("file://C:/image.JPEG").Extension())
	assert.Equal(t, "", storage.NewURI("file://C:/directory/").Extension())
}
