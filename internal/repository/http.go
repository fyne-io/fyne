package repository

import (
	"errors"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage/repository"
)

// declare conformance with repository types
var _ repository.Repository = (*HTTPRepository)(nil)

type remoteFile struct {
	*http.Response
	uri fyne.URI
}

func (f *remoteFile) Close() error {
	if f.Response == nil {
		return nil
	}
	return f.Response.Body.Close()
}

func (f *remoteFile) Read(p []byte) (int, error) {
	if f.Response == nil {
		return 0, nil
	}
	return f.Response.Body.Read(p)
}

func (f *remoteFile) URI() fyne.URI {
	return f.uri
}

// HTTPRepository implements a proxy for interacting with remote resources
// using golang's net/http library.
//
// This repository is suitable to handle the http:// and https:// scheme.
//
// Since: 2.1
type HTTPRepository struct{}

// NewHTTPRepository creates a new HTTPRepository instance.
// The caller needs to call repository.Register() with the result of this function.
//
// Since: 2.1
func NewHTTPRepository() *HTTPRepository {
	return &HTTPRepository{}
}

func constructURI(u fyne.URI) string {
	uri := ""
	uri += u.Scheme() + "://"
	uri += u.Authority()
	if u.Path() != "" {
		uri += u.Path()
	}
	if u.Query() != "" {
		uri += "?" + u.Query()
	}
	if u.Fragment() != "" {
		uri += "#" + u.Fragment()
	}

	return uri
}

// Exists checks whether the the resource at u returns a
// non "404 NOT FOUND" response header.
//
// The method is a part of the implementation for repository.Repository.Exists
//
// Since: 2.1
func (r *HTTPRepository) Exists(u fyne.URI) (bool, error) {
	uri := constructURI(u)
	resp, err := http.Head(uri)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

// Reader provides a interface for reading the body of the response received
// from the request to u.
//
// The method is a part of the implementation for repository.Repository.Exists
//
// Since: 2.1
func (r *HTTPRepository) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	uri := constructURI(u)
	resp, err := http.Get(uri)

	return &remoteFile{Response: resp, uri: u}, err
}

// CanRead makes a HEAD HTTP request to analyse the headers received
// from the remote server.
// Any response status code apart from 2xx is considered to be invalid.
//
// CanRead implements repository.Repository.CanRead
//
// Since: 2.1
func (r *HTTPRepository) CanRead(u fyne.URI) (bool, error) {
	uri := constructURI(u)
	resp, err := http.Head(uri)
	if err != nil {
		return false, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return false, errors.New("remote server did not return a successful response")
	}

	return true, nil
}

// Destroy implements repository.Repository.Destroy
//
// Sine: 2.1
func (r *HTTPRepository) Destroy(string) {
	// do nothing
}
