package repository

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper function to register the HTTP and HTTPS repositories
func registerRepositories() (*HTTPRepository, *HTTPRepository) {
	http := NewHTTPRepository()
	repository.Register("http", http)
	https := NewHTTPRepository()
	repository.Register("https", https)

	return http, https
}

func TestHTTPRepositoryRegistration(t *testing.T) {
	http, https := registerRepositories()

	// Test HTTPRespository registration for http scheme
	httpurl, err := storage.ParseURI("http://foo.com/bar")
	require.NoError(t, err)

	httpRepo, err := repository.ForURI(httpurl)
	require.NoError(t, err)
	assert.Equal(t, http, httpRepo)

	// test HTTPRepository registration for https scheme
	httpsurl, err := storage.ParseURI("https://foo.com/bar")
	require.NoError(t, err)

	httpsRepo, err := repository.ForURI(httpsurl)
	require.NoError(t, err)
	assert.Equal(t, https, httpsRepo)
}

func TestHTTPRepositoryExists(t *testing.T) {
	registerRepositories()

	var validHost string
	invalidPath := "/path-does-not-exist"
	invalidHost := "http://invalid.host.for.test"

	// start a test server to test http calls
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == invalidPath {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}))
	defer ts.Close()

	validHost = ts.URL

	// test a valid url
	resExists, err := storage.ParseURI(validHost)
	require.NoError(t, err)

	exists, err := storage.Exists(resExists)
	require.NoError(t, err)
	assert.True(t, exists)

	// test a valid host with an invalid path
	resNotExists, err := storage.ParseURI(validHost + invalidPath)
	require.NoError(t, err)

	exists, err = storage.Exists(resNotExists)
	require.NoError(t, err)
	assert.False(t, exists)

	// test an invalid host
	resInvalid, err := storage.ParseURI(invalidHost + invalidPath)
	require.NoError(t, err)

	exists, err = storage.Exists(resInvalid)
	require.Error(t, err)
	assert.False(t, exists)
}

func TestHTTPRepositoryReader(t *testing.T) {
	registerRepositories()

	var validHost string
	invalidPath := "/path-does-not-exist"
	invalidHost := "http://invalid.host.for.test"

	testData := []byte{1, 2, 3, 4, 5}

	// start a test server to test http calls
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == invalidPath {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Write(testData)
		}
	}))
	defer ts.Close()

	validHost = ts.URL

	// test a valid url returning a valid response body
	resExists, err := storage.ParseURI(validHost)
	require.NoError(t, err)

	data := make([]byte, len(testData))
	// the reader should be obtained without an error
	reader, err := storage.Reader(resExists)
	require.NoError(t, err)

	// data read from the reader should be equal to testData
	n, err := reader.Read(data)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, len(data), n)
	assert.Equal(t, testData, data)

	// test a invalid url returning an error
	resInvalid, err := storage.ParseURI(invalidHost)
	require.NoError(t, err)

	data = make([]byte, len(testData))
	// the reader should not have any data and an error should be obtained
	reader, err = storage.Reader(resInvalid)
	require.Error(t, err)

	// no data should be read from the reader
	n, err = reader.Read(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestHTTPRepositoryCanRead(t *testing.T) {
	registerRepositories()

	var validHost string
	invalidPath := "/path-does-not-exist"
	invalidHost := "http://invalid.host.for.test/"

	// start a test server to test http calls
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == invalidPath {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}))
	defer ts.Close()

	validHost = ts.URL

	// test a valid url returning a valid response body
	resExists, err := storage.ParseURI(validHost)
	require.NoError(t, err)

	canRead, err := storage.CanRead(resExists)
	require.NoError(t, err)
	assert.True(t, canRead)

	// test a invalid url for a valid host
	resNotExists, err := storage.ParseURI(validHost + invalidPath)
	require.NoError(t, err)

	canRead, err = storage.CanRead(resNotExists)
	require.Error(t, err)
	assert.False(t, canRead)

	// test a invalid host
	resInvalid, err := storage.ParseURI(invalidHost)
	require.NoError(t, err)

	canRead, err = storage.CanRead(resInvalid)
	require.Error(t, err)
	assert.False(t, canRead)
}
