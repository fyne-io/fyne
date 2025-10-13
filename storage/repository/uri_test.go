package repository

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var benchString string

func BenchmarkURIString(b *testing.B) {
	var str string
	input, _ := ParseURI("foo://example.com:8042/over/there?name=ferret#nose")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = input.String()
	}

	benchString = str
}

func TestURIExtension(t *testing.T) {
	uri := NewFileURI("file")
	assert.Equal(t, "", uri.Extension())

	uri = NewFileURI("../file")
	assert.Equal(t, "", uri.Extension())

	uri = NewFileURI("file.txt")
	assert.Equal(t, ".txt", uri.Extension())

	uri = NewFileURI("file.tar.gz")
	assert.Equal(t, ".gz", uri.Extension())

	uri = NewFileURI("/path/.txt")
	assert.Equal(t, ".txt", uri.Extension())
}

func TestURIName(t *testing.T) {
	uri := NewFileURI("file")
	assert.Equal(t, "file", uri.Name())

	uri = NewFileURI("file.txt")
	assert.Equal(t, "file.txt", uri.Name())

	uri = NewFileURI("/somewhere/file.txt")
	assert.Equal(t, "file.txt", uri.Name())

	uri = NewFileURI("/path/.txt")
	assert.Equal(t, ".txt", uri.Name())

	if runtime.GOOS == "windows" {
		uri = NewFileURI("C://somewhere/file.txt")
		assert.Equal(t, "file.txt", uri.Name())

		uri = NewFileURI("C://somewhere")
		assert.Equal(t, "somewhere", uri.Name())

		uri = NewFileURI("C://")
		assert.Equal(t, "C:", uri.Name())
	}
}
