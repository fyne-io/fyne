package repository

import (
	"io/ioutil"
	"testing"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepositoryRegistration(t *testing.T) {
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)

	// this should never fail, and we assume it doesn't in other tests here
	// for brevity
	foo, err := storage.ParseURI("mem://foo")
	assert.Nil(t, err)

	// make sure we get the same repo back
	repo, err := repository.ForURI(foo)
	assert.Nil(t, err)
	assert.Equal(t, m, repo)

	// test that re-registration also works
	m2 := NewInMemoryRepository("mem")
	repository.Register("mem", m2)
	assert.False(t, m == m2) // this is explicitly intended to be pointer comparison
	repo, err = repository.ForURI(foo)
	assert.Nil(t, err)
	assert.Equal(t, m2, repo)
}

func TestInMemoryRepositoryParsing(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)

	// since we assume in some other places that these can be parsed
	// without error, lets also explicitly test to make sure

	foo, err := storage.ParseURI("mem:///foo")
	assert.Nil(t, err)
	assert.NotNil(t, foo)

	bar, err := storage.ParseURI("mem:///bar")
	assert.Nil(t, err)
	assert.NotNil(t, bar)

	baz, _ := storage.ParseURI("mem:///baz")
	assert.Nil(t, err)
	assert.NotNil(t, baz)
}

func TestInMemoryRepositoryExists(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo"] = []byte{}
	m.Data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

	fooExists, err := storage.Exists(foo)
	assert.True(t, fooExists)
	assert.Nil(t, err)

	barExists, err := storage.Exists(bar)
	assert.True(t, barExists)
	assert.Nil(t, err)

	bazExists, err := storage.Exists(baz)
	assert.False(t, bazExists)
	assert.Nil(t, err)
}

func TestInMemoryRepositoryReader(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo"] = []byte{}
	m.Data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

	fooReader, err := storage.Reader(foo)
	assert.Nil(t, err)
	fooData, err := ioutil.ReadAll(fooReader)
	assert.Equal(t, []byte{}, fooData)
	assert.Nil(t, err)

	barReader, err := storage.Reader(bar)
	assert.Nil(t, err)
	barData, err := ioutil.ReadAll(barReader)
	assert.Equal(t, []byte{1, 2, 3}, barData)
	assert.Nil(t, err)

	bazReader, err := storage.Reader(baz)
	assert.Nil(t, bazReader)
	assert.NotNil(t, err)
}

func TestInMemoryRepositoryCanRead(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo"] = []byte{}
	m.Data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

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

func TestInMemoryRepositoryWriter(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo"] = []byte{}
	m.Data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

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

func TestInMemoryRepositoryCanWrite(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo"] = []byte{}
	m.Data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

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

func TestInMemoryRepositoryParent(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo/bar/baz"] = []byte{}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo/bar/baz")
	fooExpectedParent, _ := storage.ParseURI("mem:///foo/bar")
	fooExists, err := storage.Exists(foo)
	assert.True(t, fooExists)
	assert.Nil(t, err)

	fooParent, err := storage.Parent(foo)
	assert.Nil(t, err)
	assert.Equal(t, fooExpectedParent.String(), fooParent.String())
}

func TestInMemoryRepositoryChild(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo/bar/baz")
	fooExpectedChild, _ := storage.ParseURI("mem:///foo/bar/baz/quux")

	fooChild, err := storage.Child(foo, "quux")
	assert.Nil(t, err)
	assert.Equal(t, fooExpectedChild.String(), fooChild.String())
}

func TestInMemoryRepositoryCopy(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo"] = []byte{1, 2, 3}

	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")

	err := storage.Copy(foo, bar)
	assert.Nil(t, err)

	assert.Equal(t, m.Data["/foo"], m.Data["/bar"])
}

func TestInMemoryRepositoryMove(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo"] = []byte{1, 2, 3}

	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")

	err := storage.Move(foo, bar)
	assert.Nil(t, err)

	assert.Equal(t, []byte{1, 2, 3}, m.Data["/bar"])

	exists, err := m.Exists(foo)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestInMemoryRepositoryListing(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)
	m.Data["/foo"] = []byte{1, 2, 3}
	m.Data["/foo/bar"] = []byte{1, 2, 3}
	m.Data["/foo/baz/"] = []byte{1, 2, 3}
	m.Data["/foo/baz/quux"] = []byte{1, 2, 3}

	foo, _ := storage.ParseURI("mem:///foo")

	canList, err := storage.CanList(foo)
	assert.Nil(t, err)
	assert.True(t, canList)

	listing, err := storage.List(foo)
	assert.Nil(t, err)
	stringListing := []string{}
	for _, u := range listing {
		stringListing = append(stringListing, u.String())
	}
	assert.ElementsMatch(t, []string{"mem:///foo/bar", "mem:///foo/baz/"}, stringListing)
}

func TestInMemoryRepositoryCreateListable(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewInMemoryRepository("mem")
	repository.Register("mem", m)

	foo, _ := storage.ParseURI("mem:///foo")

	err := storage.CreateListable(foo)
	assert.Nil(t, err)

	assert.Equal(t, m.Data["/foo"], []byte{})

	// trying to create something we already created should fail
	err = storage.CreateListable(foo)
	assert.NotNil(t, err)

	// NOTE: creating an InMemoryRepository path with a non-extant parent
	// is specifically not an error, so that case is not tested.
}
