package storage_test

import (
	"io/ioutil"
	"runtime"
	"testing"

	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	_ "fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestURIAuthority(t *testing.T) {

	// from IETF RFC 3986
	s := "foo://example.com:8042/over/there?name=ferret#nose"
	u, err := storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "example.com:8042", u.Authority())

	// from IETF RFC 3986
	s = "urn:example:animal:ferret:nose"
	u, err = storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "", u.Authority())
}

func TestURIPath(t *testing.T) {
	// from IETF RFC 3986
	s := "foo://example.com:8042/over/there?name=ferret#nose"
	u, err := storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "/over/there", u.Path())

	s = "foo:///over/there"
	u, err = storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "/over/there", u.Path())

	// NOTE: if net/url supported RFC3986, this would pass
	// s = "foo://over/there"
	// u, err = storage.ParseURI(s)
	// assert.Nil(t, err)
	// assert.Equal(t, "over/there", u.Path())

	// from IETF RFC 3986
	s = "urn:example:animal:ferret:nose"
	u, err = storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "example:animal:ferret:nose", u.Path())
}

func TestURIQuery(t *testing.T) {
	// from IETF RFC 3986
	s := "foo://example.com:8042/over/there?name=ferret#nose"
	u, err := storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "name=ferret", u.Query())

	// from IETF RFC 3986
	s = "urn:example:animal:ferret:nose"
	u, err = storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "", u.Query())
}

func TestURIFragment(t *testing.T) {
	// from IETF RFC 3986
	s := "foo://example.com:8042/over/there?name=ferret#nose"
	u, err := storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "nose", u.Fragment())

	// from IETF RFC 3986
	s = "urn:example:animal:ferret:nose"
	u, err = storage.ParseURI(s)
	assert.Nil(t, err)
	assert.Equal(t, "", u.Fragment())
}

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
	assert.Equal(t, "file", storage.NewURI("FILE://C:/autoexec.bat").Scheme())
	assert.Equal(t, "file", storage.NewFileURI("/home/user/file.txt").Scheme())
}

func TestURI_Name(t *testing.T) {
	assert.Equal(t, "text.txt", storage.NewURI("file:///text.txt").Name())
	assert.Equal(t, "library.a", storage.NewURI("file:///library.a").Name())
	assert.Equal(t, "directory", storage.NewURI("file://C:/directory/").Name())
	assert.Equal(t, "image.JPEG", storage.NewFileURI("C:/image.JPEG").Name())
}

func TestURI_Parent(t *testing.T) {
	// note the trailing slashes are significant, as they tend to belie a
	// directory

	parent, err := storage.Parent(storage.NewURI("file:///foo/bar/baz"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewFileURI("/foo/bar/baz/"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewURI("file://C:/foo/bar/baz/"))
	assert.Nil(t, err)
	assert.Equal(t, "file://C:/foo/bar/", parent.String())

	// TODO hook in an http/https handler
	//parent, err = storage.Parent(storage.NewURI("http://foo/bar/baz/"))
	//assert.Nil(t, err)
	//assert.Equal(t, "http://foo/bar/", parent.String())
	//
	//parent, err = storage.Parent(storage.NewURI("http:////foo/bar/baz/"))
	//assert.Nil(t, err)
	//assert.Equal(t, "http://foo/bar/", parent.String())
	//
	//_, err = storage.Parent(storage.NewURI("http://foo"))
	//assert.Equal(t, repository.ErrURIRoot, err)
	//
	//_, err = storage.Parent(storage.NewURI("http:///"))
	//assert.Equal(t, repository.ErrURIRoot, err)
	//
	//parent, err = storage.Parent(storage.NewURI("https://///foo/bar/"))
	//assert.Nil(t, err)
	//assert.Equal(t, "https:///foo/", parent.String())

	if runtime.GOOS == "windows" {
		// Only the Windows version of filepath will know how to handle
		// backslashes.
		uri := storage.NewURI("file://C:\\foo\\bar\\baz\\")
		assert.Equal(t, "file://C:/foo/bar/baz/", uri.String())
		uri = storage.NewFileURI("C:\\foo\\bar\\baz\\")
		assert.Equal(t, "file://C:/foo/bar/baz/", uri.String())

		parent, err = storage.Parent(uri)
		assert.Nil(t, err)
		assert.Equal(t, "file://C:/foo/bar/", parent.String())
	}

	_, err = storage.Parent(storage.NewURI("file:///"))
	assert.Equal(t, repository.ErrURIRoot, err)

	if runtime.GOOS == "windows" {
		// This is only an error under Windows, on *NIX this is a
		// relative path to a directory named "C:", which is completely
		// valid.

		// This should cause an error, since this is a Windows-style
		// path and thus we can't get the parent of a drive letter.
		_, err = storage.Parent(storage.NewURI("file://C:/"))
		assert.Equal(t, repository.ErrURIRoot, err)
	}

	// Windows supports UNIX-style paths. /C:/ is also a valid path.
	parent, err = storage.Parent(storage.NewURI("file:///C:/"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///", parent.String())
}

func TestURI_Extension(t *testing.T) {
	assert.Equal(t, ".txt", storage.NewURI("file:///text.txt").Extension())
	assert.Equal(t, ".JPEG", storage.NewURI("file://C:/image.JPEG").Extension())
	assert.Equal(t, "", storage.NewURI("file://C:/directory/").Extension())
	assert.Equal(t, ".a", storage.NewFileURI("/library.a").Extension())
}

func TestURI_Child(t *testing.T) {
	p, _ := storage.Child(storage.NewURI("file:///foo/bar"), "baz")
	assert.Equal(t, "file:///foo/bar/baz", p.String())

	p, _ = storage.Child(storage.NewURI("file:///foo/bar/"), "baz")
	assert.Equal(t, "file:///foo/bar/baz", p.String())

	if runtime.GOOS == "windows" {
		// Only the Windows version of filepath will know how to handle
		// backslashes.
		uri := storage.NewURI("file://C:\\foo\\bar\\")
		assert.Equal(t, "file://C:/foo/bar/", uri.String())

		p, _ = storage.Child(uri, "baz")
		assert.Equal(t, "file://C:/foo/bar/baz", p.String())
	}
}

func TestExists(t *testing.T) {
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo/bar/baz"] = []byte{}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("uritest:///foo/bar/baz")
	fooExpectedParent, _ := storage.ParseURI("uritest:///foo/bar")
	fooExists, err := storage.Exists(foo)
	assert.True(t, fooExists)
	assert.Nil(t, err)

	fooParent, err := storage.Parent(foo)
	assert.Nil(t, err)
	assert.Equal(t, fooExpectedParent.String(), fooParent.String())

}

func TestWriteAndDelete(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo"] = []byte{}
	m.Data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("uritest:///foo")
	bar, _ := storage.ParseURI("uritest:///bar")
	baz, _ := storage.ParseURI("uritest:///baz")

	// write some data and assert there are no errors
	fooWriter, err := storage.Writer(foo)
	assert.Nil(t, err)
	assert.NotNil(t, fooWriter)

	barWriter, err := storage.Writer(bar)
	assert.Nil(t, err)
	assert.NotNil(t, barWriter)

	bazWriter, err := storage.Writer(baz)
	assert.Nil(t, err)
	assert.NotNil(t, bazWriter)

	n, err := fooWriter.Write([]byte{1, 2, 3, 4, 5})
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = barWriter.Write([]byte{6, 7, 8, 9})
	assert.Nil(t, err)
	assert.Equal(t, 4, n)

	n, err = bazWriter.Write([]byte{5, 4, 3, 2, 1, 0})
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	fooWriter.Close()
	barWriter.Close()
	bazWriter.Close()

	// now make sure we can read the data back correctly
	fooReader, err := storage.Reader(foo)
	assert.Nil(t, err)
	fooData, err := ioutil.ReadAll(fooReader)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, fooData)
	assert.Nil(t, err)

	barReader, err := storage.Reader(bar)
	assert.Nil(t, err)
	barData, err := ioutil.ReadAll(barReader)
	assert.Equal(t, []byte{6, 7, 8, 9}, barData)
	assert.Nil(t, err)

	bazReader, err := storage.Reader(baz)
	assert.Nil(t, err)
	bazData, err := ioutil.ReadAll(bazReader)
	assert.Equal(t, []byte{5, 4, 3, 2, 1, 0}, bazData)
	assert.Nil(t, err)

	// now let's test deletion
	err = storage.Delete(foo)
	assert.Nil(t, err)

	err = storage.Delete(bar)
	assert.Nil(t, err)

	err = storage.Delete(baz)
	assert.Nil(t, err)

	fooExists, err := storage.Exists(foo)
	assert.False(t, fooExists)
	assert.Nil(t, err)

	barExists, err := storage.Exists(bar)
	assert.False(t, barExists)
	assert.Nil(t, err)

	bazExists, err := storage.Exists(baz)
	assert.False(t, bazExists)
	assert.Nil(t, err)

}

func TestCanWrite(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo"] = []byte{}
	m.Data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("uritest:///foo")
	bar, _ := storage.ParseURI("uritest:///bar")
	baz, _ := storage.ParseURI("uritest:///baz")

	fooCanWrite, err := storage.CanWrite(foo)
	assert.True(t, fooCanWrite)
	assert.Nil(t, err)

	barCanWrite, err := storage.CanWrite(bar)
	assert.True(t, barCanWrite)
	assert.Nil(t, err)

	bazCanWrite, err := storage.CanWrite(baz)
	assert.True(t, bazCanWrite)
	assert.Nil(t, err)
}

func TestCanRead(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo"] = []byte{}
	m.Data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("uritest:///foo")
	bar, _ := storage.ParseURI("uritest:///bar")
	baz, _ := storage.ParseURI("uritest:///baz")

	fooCanRead, err := storage.CanRead(foo)
	assert.True(t, fooCanRead)
	assert.Nil(t, err)

	barCanRead, err := storage.CanRead(bar)
	assert.True(t, barCanRead)
	assert.Nil(t, err)

	bazCanRead, err := storage.CanRead(baz)
	assert.False(t, bazCanRead)
	assert.Nil(t, err)
}

func TestCopy(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo"] = []byte{1, 2, 3}

	foo, _ := storage.ParseURI("uritest:///foo")
	bar, _ := storage.ParseURI("uritest:///bar")

	err := storage.Copy(foo, bar)
	assert.Nil(t, err)

	assert.Equal(t, m.Data["/foo"], m.Data["/bar"])
}

func TestRepositoryMove(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo"] = []byte{1, 2, 3}

	foo, _ := storage.ParseURI("uritest:///foo")
	bar, _ := storage.ParseURI("uritest:///bar")

	err := storage.Move(foo, bar)
	assert.Nil(t, err)

	assert.Equal(t, []byte{1, 2, 3}, m.Data["/bar"])

	exists, err := m.Exists(foo)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestRepositoryListing(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo"] = []byte{1, 2, 3}
	m.Data["/foo/bar"] = []byte{1, 2, 3}
	m.Data["/foo/baz/"] = []byte{1, 2, 3}
	m.Data["/foo/baz/quux"] = []byte{1, 2, 3}

	foo, _ := storage.ParseURI("uritest:///foo")

	canList, err := storage.CanList(foo)
	assert.Nil(t, err)
	assert.True(t, canList)

	listing, err := storage.List(foo)
	assert.Nil(t, err)
	stringListing := []string{}
	for _, u := range listing {
		stringListing = append(stringListing, u.String())
	}
	assert.ElementsMatch(t, []string{"uritest:///foo/bar", "uritest:///foo/baz/"}, stringListing)
}

func TestCreateListable(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)

	foo, _ := storage.ParseURI("uritest:///foo")

	err := storage.CreateListable(foo)
	assert.Nil(t, err)

	assert.Equal(t, m.Data["/foo"], []byte{})

	// trying to create something we already created should fail
	err = storage.CreateListable(foo)
	assert.NotNil(t, err)

	// NOTE: creating an InMemoryRepository path with a non-extant parent
	// is specifically not an error, so that case is not tested.
}
