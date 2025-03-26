package repository

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkExistence(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func TestFileRepositoryRegistration(t *testing.T) {
	f := NewFileRepository()
	repository.Register("file", f)

	// this should never fail, and we assume it doesn't in other tests here
	// for brevity
	foo, err := storage.ParseURI("file:///foo")
	require.NoError(t, err)

	// make sure we get the same repo back
	repo, err := repository.ForURI(foo)
	require.NoError(t, err)
	assert.Equal(t, f, repo)
}

func TestFileRepositoryExists(t *testing.T) {
	dir := t.TempDir()

	existsPath := path.Join(dir, "exists")
	notExistsPath := path.Join(dir, "notExists")

	err := os.WriteFile(existsPath, []byte{1, 2, 3, 4}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	ex, err := storage.Exists(storage.NewFileURI(existsPath))
	require.NoError(t, err)
	assert.True(t, ex)

	ex, err = storage.Exists(storage.NewFileURI(notExistsPath))
	require.NoError(t, err)
	assert.False(t, ex)
}

func TestFileRepositoryReader(t *testing.T) {
	dir := t.TempDir()

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	bazPath := path.Join(dir, "baz")
	err := os.WriteFile(fooPath, []byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(barPath, []byte{1, 2, 3}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Set up our repository - it's OK if we already registered it...
	f := NewFileRepository()
	repository.Register("file", f)

	// ...and some URIs - we know that they will not fail parsing
	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)
	baz := storage.NewFileURI(bazPath)

	// Make sure we can read the empty file.
	fooReader, err := storage.Reader(foo)
	require.NoError(t, err)
	fooData, err := io.ReadAll(fooReader)
	assert.Equal(t, []byte{}, fooData)
	require.NoError(t, err)
	fooReader.Close()

	// Make sure we can read the file with data.
	barReader, err := storage.Reader(bar)
	require.NoError(t, err)
	barData, err := io.ReadAll(barReader)
	assert.Equal(t, []byte{1, 2, 3}, barData)
	require.NoError(t, err)
	barReader.Close()

	// Make sure we get an error if the file doesn't exist.
	bazReader, err := storage.Reader(baz)
	require.Error(t, err)
	bazReader.Close()

	// Also test that CanRead returns the expected results.
	fooCanRead, err := storage.CanRead(foo)
	assert.True(t, fooCanRead)
	require.NoError(t, err)

	barCanRead, err := storage.CanRead(bar)
	assert.True(t, barCanRead)
	require.NoError(t, err)

	bazCanRead, err := storage.CanRead(baz)
	assert.False(t, bazCanRead)
	require.NoError(t, err)
}

func TestFileRepositoryWriter(t *testing.T) {
	dir := t.TempDir()

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	bazPath := path.Join(dir, "baz")
	spamHamPath := path.Join(dir, "spam", "ham")
	err := os.WriteFile(fooPath, []byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(barPath, []byte{1, 2, 3}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Set up our repository - it's OK if we already registered it...
	f := NewFileRepository()
	repository.Register("file", f)

	// ...and some URIs - we know that they will not fail parsing
	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)
	baz := storage.NewFileURI(bazPath)
	spamHam := storage.NewFileURI(spamHamPath)

	// Make sure that spamHam errors, since writing to a non-existent
	// parent directory should be an error.
	spamHamWriter, err := storage.Writer(spamHam)
	require.Error(t, err)
	if err == nil {
		// Keep this from bodging up the Windows tests if this is
		// created in error, and then we try to delete it while there
		// is an open handle.
		spamHamWriter.Close()
	}

	// write some data and assert there are no errors
	fooWriter, err := storage.Writer(foo)
	require.NoError(t, err)
	assert.NotNil(t, fooWriter)

	barWriter, err := storage.Writer(bar)
	require.NoError(t, err)
	assert.NotNil(t, barWriter)

	bazWriter, err := storage.Writer(baz)
	require.NoError(t, err)
	assert.NotNil(t, bazWriter)

	n, err := fooWriter.Write([]byte{1, 2, 3, 4, 5})
	require.NoError(t, err)
	assert.Equal(t, 5, n)

	n, err = barWriter.Write([]byte{6, 7, 8, 9})
	require.NoError(t, err)
	assert.Equal(t, 4, n)

	n, err = bazWriter.Write([]byte{5, 4, 3, 2, 1, 0})
	require.NoError(t, err)
	assert.Equal(t, 6, n)

	fooWriter.Close()
	barWriter.Close()
	bazWriter.Close()

	bazAppender, err := storage.Appender(baz)
	require.NoError(t, err)
	n, err = bazAppender.Write([]byte{1, 2, 3, 4, 5})
	require.NoError(t, err)
	assert.Equal(t, 5, n)

	bazAppender.Close()

	// now make sure we can read the data back correctly
	fooReader, err := storage.Reader(foo)
	require.NoError(t, err)
	fooData, err := io.ReadAll(fooReader)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, fooData)
	require.NoError(t, err)

	barReader, err := storage.Reader(bar)
	require.NoError(t, err)
	barData, err := io.ReadAll(barReader)
	assert.Equal(t, []byte{6, 7, 8, 9}, barData)
	require.NoError(t, err)

	bazReader, err := storage.Reader(baz)
	require.NoError(t, err)
	bazData, err := io.ReadAll(bazReader)
	assert.Equal(t, []byte{5, 4, 3, 2, 1, 0, 1, 2, 3, 4, 5}, bazData)
	require.NoError(t, err)

	// close the readers, since Windows won't let us delete things with
	// open handles to them
	fooReader.Close()
	barReader.Close()
	bazReader.Close()

	// now let's test deletion
	err = storage.Delete(foo)
	require.NoError(t, err)

	err = storage.Delete(bar)
	require.NoError(t, err)

	err = storage.Delete(baz)
	require.NoError(t, err)

	fooExists, err := storage.Exists(foo)
	assert.False(t, fooExists)
	require.NoError(t, err)

	barExists, err := storage.Exists(bar)
	assert.False(t, barExists)
	require.NoError(t, err)

	bazExists, err := storage.Exists(baz)
	assert.False(t, bazExists)
	require.NoError(t, err)
}

func TestFileRepositoryCanWrite(t *testing.T) {
	dir := t.TempDir()

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	bazPath := path.Join(dir, "baz")
	err := os.WriteFile(fooPath, []byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(barPath, []byte{1, 2, 3}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Set up our repository - it's OK if we already registered it...
	f := NewFileRepository()
	repository.Register("file", f)

	// ...and some URIs - we know that they will not fail parsing
	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)
	baz := storage.NewFileURI(bazPath)

	fooCanWrite, err := storage.CanWrite(foo)
	assert.True(t, fooCanWrite)
	require.NoError(t, err)

	barCanWrite, err := storage.CanWrite(bar)
	assert.True(t, barCanWrite)
	require.NoError(t, err)

	bazCanWrite, err := storage.CanWrite(baz)
	assert.True(t, bazCanWrite)
	require.NoError(t, err)
}

func TestFileRepositoryParent(t *testing.T) {
	// Set up our repository - it's OK if we already registered it.
	f := NewFileRepository()
	repository.Register("file", f)

	// note the trailing slashes are significant, as they tend to belie a
	// directory

	parent, err := storage.Parent(storage.NewFileURI("/foo/bar/baz"))
	require.NoError(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewFileURI("/foo/bar/baz/"))
	require.NoError(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewFileURI("C:/foo/bar/baz/"))
	require.NoError(t, err)
	assert.Equal(t, "file://C:/foo/bar/", parent.String())

	if runtime.GOOS == "windows" {
		// Only the Windows version of filepath will know how to handle
		// backslashes.
		uri := storage.NewFileURI("C:\\foo\\bar\\baz\\")
		assert.Equal(t, "file://C:/foo/bar/baz/", uri.String())
		uri = storage.NewFileURI("C:\\foo\\bar\\baz\\")
		assert.Equal(t, "file://C:/foo/bar/baz/", uri.String())

		parent, err = storage.Parent(uri)
		require.NoError(t, err)
		assert.Equal(t, "file://C:/foo/bar/", parent.String())
	}

	_, err = storage.Parent(storage.NewFileURI("/"))
	assert.Equal(t, repository.ErrURIRoot, err)

	if runtime.GOOS == "windows" {
		// This is only an error under Windows, on *NIX this is a
		// relative path to a directory named "C:", which is completely
		// valid.

		// This should cause an error, since this is a Windows-style
		// path and thus we can't get the parent of a drive letter.
		_, err = storage.Parent(storage.NewFileURI("C:/"))
		assert.Equal(t, repository.ErrURIRoot, err)
	}

	// Windows supports UNIX-style paths. /C:/ is also a valid path.
	parent, err = storage.Parent(storage.NewFileURI("/C:/"))
	require.NoError(t, err)
	assert.Equal(t, "file:///", parent.String())
}

func TestFileRepositoryChild(t *testing.T) {
	// Set up our repository - it's OK if we already registered it.
	f := NewFileRepository()
	repository.Register("file", f)

	p, _ := storage.Child(storage.NewFileURI("/foo/bar"), "baz")
	assert.Equal(t, "file:///foo/bar/baz", p.String())

	p, _ = storage.Child(storage.NewFileURI("/foo/bar/"), "baz")
	assert.Equal(t, "file:///foo/bar/baz", p.String())

	if runtime.GOOS == "windows" {
		// Only the Windows version of filepath will know how to handle
		// backslashes.
		uri := storage.NewFileURI("C:\\foo\\bar\\")
		assert.Equal(t, "file://C:/foo/bar/", uri.String())

		p, _ = storage.Child(uri, "baz")
		assert.Equal(t, "file://C:/foo/bar/baz", p.String())
	}
}

func TestFileRepositoryCopy(t *testing.T) {
	dir := t.TempDir()

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	err := os.WriteFile(fooPath, []byte{1, 2, 3, 4, 5}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)

	err = storage.Copy(foo, bar)
	require.NoError(t, err)

	fooData, err := os.ReadFile(fooPath)
	require.NoError(t, err)

	barData, err := os.ReadFile(barPath)
	require.NoError(t, err)

	assert.Equal(t, fooData, barData)
}

func TestFileRepositoryMove(t *testing.T) {
	dir := t.TempDir()

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	err := os.WriteFile(fooPath, []byte{1, 2, 3, 4, 5}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)

	err = storage.Move(foo, bar)
	require.NoError(t, err)

	barData, err := os.ReadFile(barPath)
	require.NoError(t, err)

	assert.Equal(t, []byte{1, 2, 3, 4, 5}, barData)

	// Make sure that the source doesn't exist anymore.
	ex, err := storage.Exists(foo)
	require.NoError(t, err)
	assert.False(t, ex)
}

func TestFileRepositoryMoveDirectory(t *testing.T) {
	dir := t.TempDir()

	// Create a file in a dir to test with
	parentPath := path.Join(dir, "parent")
	fooPath := path.Join(parentPath, "foo")
	newParentPath := path.Join(dir, "newParent")
	newFooPath := path.Join(newParentPath, "foo")

	_ = os.Mkdir(parentPath, 0755)
	err := os.WriteFile(fooPath, []byte{1, 2, 3, 4, 5}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	parent := storage.NewFileURI(parentPath)
	foo := storage.NewFileURI(fooPath)
	newParent := storage.NewFileURI(newParentPath)

	err = storage.Move(parent, newParent)
	require.NoError(t, err)

	newData, err := os.ReadFile(newFooPath)
	require.NoError(t, err)

	assert.Equal(t, []byte{1, 2, 3, 4, 5}, newData)

	// Make sure that the source doesn't exist anymore.
	ex, err := storage.Exists(foo)
	require.NoError(t, err)
	assert.False(t, ex)
}

func TestFileRepositoryListing(t *testing.T) {
	dir := t.TempDir()

	// Create some files to tests with.
	fooPath := path.Join(dir, "foo")
	os.Mkdir(fooPath, 0755)
	os.Mkdir(path.Join(fooPath, "bar"), 0755)
	os.Mkdir(path.Join(fooPath, "baz"), 0755)
	os.Mkdir(path.Join(fooPath, "baz", "quux"), 0755)

	foo := storage.NewFileURI(fooPath)

	canList, err := storage.CanList(foo)
	require.NoError(t, err)
	assert.True(t, canList)

	// also check the empty dir
	childDir := storage.NewFileURI(path.Join(fooPath, "baz", "quux"))
	canList, err = storage.CanList(childDir)
	require.NoError(t, err)
	assert.True(t, canList)

	listing, err := storage.List(foo)
	require.NoError(t, err)
	assert.Len(t, listing, 2)
	stringListing := []string{}
	for _, u := range listing {
		stringListing = append(stringListing, u.String())
	}
	assert.ElementsMatch(t, []string{"file://" + filepath.ToSlash(path.Join(dir, "foo", "bar")), "file://" + filepath.ToSlash(path.Join(dir, "foo", "baz"))}, stringListing)
}

func TestFileRepositoryCreateListable(t *testing.T) {
	dir := t.TempDir()

	f := NewFileRepository()
	repository.Register("file", f)

	fooPath := path.Join(dir, "foo")
	fooBarPath := path.Join(dir, "foo", "bar")
	foo := storage.NewFileURI(fooPath)
	fooBar := storage.NewFileURI(fooBarPath)

	// Creating a dir with no parent should fail
	err := storage.CreateListable(fooBar)
	require.Error(t, err)

	// Creating foo should work though
	err = storage.CreateListable(foo)
	require.NoError(t, err)

	// and now we should be able to create fooBar
	err = storage.CreateListable(fooBar)
	require.NoError(t, err)

	// make sure the OS thinks these dirs really exist
	assert.True(t, checkExistence(fooPath))
	assert.True(t, checkExistence(fooBarPath))
}
