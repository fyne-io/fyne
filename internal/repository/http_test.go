package repository

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)

	httpRepo, err := repository.ForURI(httpurl)
	assert.Nil(t, err)
	assert.Equal(t, http, httpRepo)

	// test HTTPRepository registration for https scheme
	httpsurl, err := storage.ParseURI("https://foo.com/bar")
	assert.Nil(t, err)

	httpsRepo, err := repository.ForURI(httpsurl)
	assert.Nil(t, err)
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
	assert.Nil(t, err)

	exists, err := storage.Exists(resExists)
	assert.Nil(t, err)
	assert.Equal(t, true, exists)

	// test a valid host with an invalid path
	resNotExists, err := storage.ParseURI(validHost + invalidPath)
	assert.Nil(t, err)

	exists, err = storage.Exists(resNotExists)
	assert.Nil(t, err)
	assert.Equal(t, false, exists)

	// test an invalid host
	resInvalid, err := storage.ParseURI(invalidHost + invalidPath)
	assert.Nil(t, err)

	exists, err = storage.Exists(resInvalid)
	assert.NotNil(t, err)
	assert.Equal(t, false, exists)
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
	assert.Nil(t, err)

	data := make([]byte, len(testData))
	// the reader should be obtained without an error
	reader, err := storage.Reader(resExists)
	assert.Nil(t, err)

	// data read from the reader should be equal to testData
	n, err := reader.Read(data)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, len(data), n)
	assert.Equal(t, testData, data)

	// test a invalid url returning an error
	resInvalid, err := storage.ParseURI(invalidHost)
	assert.Nil(t, err)

	data = make([]byte, len(testData))
	// the reader should not have any data and an error should be obtained
	reader, err = storage.Reader(resInvalid)
	assert.NotNil(t, err)

	// no data should be read from the reader
	n, err = reader.Read(data)
	assert.Nil(t, err)
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
	assert.Nil(t, err)

	canRead, err := storage.CanRead(resExists)
	assert.Nil(t, err)
	assert.Equal(t, true, canRead)

	// test a invalid url for a valid host
	resNotExists, err := storage.ParseURI(validHost + invalidPath)
	assert.Nil(t, err)

	canRead, err = storage.CanRead(resNotExists)
	assert.NotNil(t, err)
	assert.Equal(t, false, canRead)

	// test a invalid host
	resInvalid, err := storage.ParseURI(invalidHost)
	assert.Nil(t, err)

	canRead, err = storage.CanRead(resInvalid)
	assert.NotNil(t, err)
	assert.Equal(t, false, canRead)
}
