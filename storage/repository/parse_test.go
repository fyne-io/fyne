package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileURI(t *testing.T) {
	assert.Equal(t, "file:///tmp/foo.txt", NewFileURI("/tmp/foo.txt").String())
	assert.Equal(t, "file://C:/tmp/foo.txt", NewFileURI("C:/tmp/foo.txt").String())
}

func TestParseURI(t *testing.T) {
	uri, err := ParseURI("file:///tmp/foo.txt")
	require.NoError(t, err)
	assert.Equal(t, "file:///tmp/foo.txt", uri.String())

	uri, err = ParseURI("file:/tmp/foo.txt")
	require.NoError(t, err)
	assert.Equal(t, "file:///tmp/foo.txt", uri.String())

	uri, err = ParseURI("file://C:/tmp/foo.txt")
	require.NoError(t, err)
	assert.Equal(t, "file://C:/tmp/foo.txt", uri.String())
}

func TestParseInvalidURI(t *testing.T) {
	uri, err := ParseURI("/tmp/foo.txt")
	require.Error(t, err)
	assert.Nil(t, uri)

	uri, err = ParseURI("file")
	require.Error(t, err)
	assert.Nil(t, uri)

	uri, err = ParseURI("file:")
	require.Error(t, err)
	assert.Nil(t, uri)

	uri, err = ParseURI("file://")
	require.Error(t, err)
	assert.Nil(t, uri)

	uri, err = ParseURI(":foo")
	require.Error(t, err)
	assert.Nil(t, uri)

	uri, err = ParseURI("scheme://0[]/invalid")
	require.Error(t, err)
	assert.Nil(t, uri)
}
