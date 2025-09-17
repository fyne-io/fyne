package storage_test

import (
	"io"
	"os"
	"runtime"
	"testing"

	"fyne.io/fyne/v2"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	_ "fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

type otherURI struct {
	path string
	fyne.URI
}

func (u *otherURI) String() string {
	return "file://" + u.path
}

func (u *otherURI) Scheme() string {
	return "file"
}

func (u *otherURI) Authority() string {
	return ""
}

func (u *otherURI) Path() string {
	return u.path
}

func (u *otherURI) Query() string {
	return ""
}

func (u *otherURI) Fragment() string {
	return ""
}

func TestURIEqual(t *testing.T) {
	first := storage.NewFileURI("/first")
	second := storage.NewFileURI("/second")
	assert.False(t, storage.EqualURI(first, second))
	assert.True(t, storage.EqualURI(first, first))

	assert.True(t, storage.EqualURI(first, storage.NewFileURI("/first")))

	assert.True(t, storage.EqualURI(nil, nil))
	assert.False(t, storage.EqualURI(first, nil))
	assert.False(t, storage.EqualURI(nil, second))

	otherURI := &otherURI{path: "/other"}
	assert.True(t, storage.EqualURI(otherURI, otherURI))
	assert.False(t, storage.EqualURI(otherURI, first))
	assert.True(t, storage.EqualURI(otherURI, storage.NewFileURI("/other")))
}

var equalSink bool

func BenchmarkEqualURI(b *testing.B) {
	first := storage.NewFileURI("/first")
	b.Run("same path and same object", func(b *testing.B) {
		b.ReportAllocs()

		equal := false
		for i := 0; i < b.N; i++ {
			equal = storage.EqualURI(first, first)
		}
		equalSink = equal
	})
	first2 := storage.NewFileURI("/first")
	b.Run("same path but different object", func(b *testing.B) {
		b.ReportAllocs()

		equal := false
		for i := 0; i < b.N; i++ {
			equal = storage.EqualURI(first, first2)
		}
		equalSink = equal
	})
	second := storage.NewFileURI("/second")
	b.Run("different path and different object", func(b *testing.B) {
		b.ReportAllocs()

		equal := false
		for i := 0; i < b.N; i++ {
			equal = storage.EqualURI(first, second)
		}
		equalSink = equal
	})

	otherFirst := &otherURI{path: "/first"}
	b.Run("custom uri types but same object", func(b *testing.B) {
		b.ReportAllocs()

		equal := false
		for i := 0; i < b.N; i++ {
			equal = storage.EqualURI(otherFirst, otherFirst)
		}
		equalSink = equal
	})
	b.Run("different uri types", func(b *testing.B) {
		b.ReportAllocs()

		equal := false
		for i := 0; i < b.N; i++ {
			equal = storage.EqualURI(first, otherFirst)
		}
		equalSink = equal
	})
}

func TestURIAuthority(t *testing.T) {
	// from IETF RFC 3986
	s := "foo://example.com:8042/over/there?name=ferret#nose"
	u, err := storage.ParseURI(s)
	assert.NoError(t, err)
	assert.Equal(t, "example.com:8042", u.Authority())

	// from IETF RFC 3986
	s = "urn:example:animal:ferret:nose"
	u, err = storage.ParseURI(s)
	assert.NoError(t, err)
	assert.Equal(t, "", u.Authority())
}

func TestURIPath(t *testing.T) {
	// from IETF RFC 3986
	s := "foo://example.com:8042/over/there?name=ferret#nose"
	u, err := storage.ParseURI(s)
	assert.NoError(t, err)
	assert.Equal(t, "/over/there", u.Path())

	s = "foo:///over/there"
	u, err = storage.ParseURI(s)
	assert.NoError(t, err)
	assert.Equal(t, "/over/there", u.Path())

	// NOTE: if net/url supported RFC3986, this would pass
	// s = "foo://over/there"
	// u, err = storage.ParseURI(s)
	// assert.Nil(t, err)
	// assert.Equal(t, "over/there", u.Path())

	// from IETF RFC 3986
	s = "urn:example:animal:ferret:nose"
	u, err = storage.ParseURI(s)
	assert.NoError(t, err)
	assert.Equal(t, "example:animal:ferret:nose", u.Path())
}

func TestURIQuery(t *testing.T) {
	// from IETF RFC 3986
	s := "foo://example.com:8042/over/there?name=ferret#nose"
	u, err := storage.ParseURI(s)
	assert.NoError(t, err)
	assert.Equal(t, "name=ferret", u.Query())

	// from IETF RFC 3986
	s = "urn:example:animal:ferret:nose"
	u, err = storage.ParseURI(s)
	assert.NoError(t, err)
	assert.Equal(t, "", u.Query())
}

func TestURIFragment(t *testing.T) {
	// from IETF RFC 3986
	s := "foo://example.com:8042/over/there?name=ferret#nose"
	u, err := storage.ParseURI(s)
	assert.NoError(t, err)
	assert.Equal(t, "nose", u.Fragment())

	// from IETF RFC 3986
	s = "urn:example:animal:ferret:nose"
	u, err = storage.ParseURI(s)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewFileURI("/foo/bar/baz/"))
	assert.NoError(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	if runtime.GOOS == "windows" {
		parent, err = storage.Parent(storage.NewURI("file://C:/foo/bar/baz/"))
		assert.NoError(t, err)
		assert.Equal(t, "file://C:/foo/bar/", parent.String())
	}

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
		assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	fooParent, err := storage.Parent(foo)
	assert.NoError(t, err)
	assert.Equal(t, fooExpectedParent.String(), fooParent.String())
}

func TestFileAbs(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Error("Could not get working directory")
		defer os.Chdir(pwd)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Error("Could not get user home directory")
	}

	os.Chdir(home)

	abs := storage.NewFileURI(home)
	rel := storage.NewFileURI(".")

	assert.Equal(t, abs.Path(), rel.Path())
	assert.Equal(t, abs.String(), rel.String())

	assert.Equal(t, "file:///", storage.NewFileURI("/").String())
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
	assert.NoError(t, err)
	assert.NotNil(t, fooWriter)

	barWriter, err := storage.Writer(bar)
	assert.NoError(t, err)
	assert.NotNil(t, barWriter)

	bazWriter, err := storage.Writer(baz)
	assert.NoError(t, err)
	assert.NotNil(t, bazWriter)

	n, err := fooWriter.Write([]byte{1, 2, 3, 4, 5})
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	n, err = barWriter.Write([]byte{6, 7, 8, 9})
	assert.NoError(t, err)
	assert.Equal(t, 4, n)

	n, err = bazWriter.Write([]byte{5, 4, 3, 2, 1, 0})
	assert.NoError(t, err)
	assert.Equal(t, 6, n)

	fooWriter.Close()
	barWriter.Close()
	bazWriter.Close()

	bazAppender, err := storage.Appender(baz)
	assert.NoError(t, err)
	n, err = bazAppender.Write([]byte{1, 2, 3, 4, 5})
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	bazAppender.Close()

	// now make sure we can read the data back correctly
	fooReader, err := storage.Reader(foo)
	assert.NoError(t, err)
	fooData, err := io.ReadAll(fooReader)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, fooData)
	assert.NoError(t, err)

	barReader, err := storage.Reader(bar)
	assert.NoError(t, err)
	barData, err := io.ReadAll(barReader)
	assert.Equal(t, []byte{6, 7, 8, 9}, barData)
	assert.NoError(t, err)

	bazReader, err := storage.Reader(baz)
	assert.NoError(t, err)
	bazData, err := io.ReadAll(bazReader)
	assert.Equal(t, []byte{5, 4, 3, 2, 1, 0, 1, 2, 3, 4, 5}, bazData)
	assert.NoError(t, err)

	// now let's test deletion
	err = storage.Delete(foo)
	assert.NoError(t, err)

	err = storage.Delete(bar)
	assert.NoError(t, err)

	err = storage.Delete(baz)
	assert.NoError(t, err)

	fooExists, err := storage.Exists(foo)
	assert.False(t, fooExists)
	assert.NoError(t, err)

	barExists, err := storage.Exists(bar)
	assert.False(t, barExists)
	assert.NoError(t, err)

	bazExists, err := storage.Exists(baz)
	assert.False(t, bazExists)
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	barCanWrite, err := storage.CanWrite(bar)
	assert.True(t, barCanWrite)
	assert.NoError(t, err)

	bazCanWrite, err := storage.CanWrite(baz)
	assert.True(t, bazCanWrite)
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	barCanRead, err := storage.CanRead(bar)
	assert.True(t, barCanRead)
	assert.NoError(t, err)

	bazCanRead, err := storage.CanRead(baz)
	assert.False(t, bazCanRead)
	assert.NoError(t, err)
}

func TestCopy(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo"] = []byte{1, 2, 3}

	foo, _ := storage.ParseURI("uritest:///foo")
	bar, _ := storage.ParseURI("uritest:///bar")

	err := storage.Copy(foo, bar)
	assert.NoError(t, err)

	assert.Equal(t, m.Data["/foo"], m.Data["/bar"])
}

func TestRepositoryCopyListable(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/parent1"] = []byte{}
	m.Data["/parent1/child"] = []byte("content")

	parent, _ := storage.ParseURI("uritest:///parent1")
	newParent, _ := storage.ParseURI("uritest:///parent2")

	err := storage.Copy(parent, newParent)
	assert.NoError(t, err)
	exists, err := m.Exists(parent)
	assert.NoError(t, err)
	assert.True(t, exists)
	exists, err = m.Exists(newParent)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, []byte("content"), m.Data["/parent1/child"])
	assert.Equal(t, []byte("content"), m.Data["/parent2/child"])
}

func TestRepositoryMove(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/foo"] = []byte{1, 2, 3}

	foo, _ := storage.ParseURI("uritest:///foo")
	bar, _ := storage.ParseURI("uritest:///bar")

	err := storage.Move(foo, bar)
	assert.NoError(t, err)

	assert.Equal(t, []byte{1, 2, 3}, m.Data["/bar"])

	exists, err := m.Exists(foo)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestRepositoryMoveListable(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := intRepo.NewInMemoryRepository("uritest")
	repository.Register("uritest", m)
	m.Data["/parent1"] = []byte{}
	m.Data["/parent1/child"] = []byte("content")

	parent, _ := storage.ParseURI("uritest:///parent1")
	newParent, _ := storage.ParseURI("uritest:///parent2")

	err := storage.Move(parent, newParent)
	assert.NoError(t, err)
	exists, err := m.Exists(parent)
	assert.NoError(t, err)
	assert.False(t, exists)
	exists, err = m.Exists(newParent)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, []byte("content"), m.Data["/parent2/child"])
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
	assert.NoError(t, err)
	assert.True(t, canList)

	listing, err := storage.List(foo)
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	assert.Equal(t, []byte{}, m.Data["/foo"])

	// trying to create something we already created should fail
	err = storage.CreateListable(foo)
	assert.Error(t, err)

	// NOTE: creating an InMemoryRepository path with a non-extant parent
	// is specifically not an error, so that case is not tested.
}
