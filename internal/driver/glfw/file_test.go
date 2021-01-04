// +build !mobile

package glfw

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

func TestGLDriver_FileReaderForURI(t *testing.T) {
	uri, _ := testURI("text.txt")
	read, err := NewGLDriver().FileReaderForURI(uri)
	assert.Nil(t, err)

	defer read.Close()

	assert.Equal(t, uri, read.URI())

	data, err := ioutil.ReadAll(read)
	assert.Nil(t, err)
	assert.Equal(t, "content", string(data))
}

func TestGLDriver_FileReaderForURI_Again(t *testing.T) { // a bug blanked files
	TestGLDriver_FileReaderForURI(t)
}

func TestGLDriver_FileWriterForURI(t *testing.T) {
	uri, path := testURI("text2.txt")
	defer os.Remove(path)
	write, err := NewGLDriver().FileWriterForURI(uri)
	assert.Nil(t, err)

	defer write.Close()

	assert.Equal(t, uri, write.URI())

	count, err := write.Write([]byte("another"))
	assert.Equal(t, 7, count)
	assert.Nil(t, err)
}

func testURI(file string) (fyne.URI, string) {
	path, _ := filepath.Abs(filepath.Join("testdata", file))
	uri := storage.NewURI("file://" + path)

	return uri, path
}

func Test_ListableURI(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "file_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tempdir)

	content := []byte("Test Content!")

	err = ioutil.WriteFile(filepath.Join(tempdir, "aaa"), content, 0666)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(tempdir, "bbb"), content, 0666)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(tempdir, "ccc"), content, 0666)
	if err != nil {
		t.Fatal(err)
	}

	uri := storage.NewURI("file://" + tempdir)

	luri, err := NewGLDriver().ListerForURI(uri)
	if err != nil {
		t.Fatal(err)
	}

	listing, err := luri.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(listing) != 3 {
		t.Errorf("Expected directory listing to contain exactly 3 items")
	}

	listingStrings := []string{}
	for _, v := range listing {
		t.Logf("stringify URI %v to %s", v, v.String())
		listingStrings = append(listingStrings, v.String())
	}

	expect := []string{
		"file://" + filepath.ToSlash(filepath.Join(tempdir, "aaa")),
		"file://" + filepath.ToSlash(filepath.Join(tempdir, "bbb")),
		"file://" + filepath.ToSlash(filepath.Join(tempdir, "ccc")),
	}

	assert.ElementsMatch(t, listingStrings, expect)

}
