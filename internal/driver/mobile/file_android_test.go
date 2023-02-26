package mobile

import (
	"net/url"
	"testing"

	"fyne.io/fyne/v2/storage"
	"github.com/stretchr/testify/assert"
)

func TestNativeFileSave(t *testing.T) {
	tests := []struct {
		name        string
		fileSave    *fileSave
		expectedErr error
	}{
		{
			name: "returns resource not found error",
			fileSave: &fileSave{
				uri: &url.URL{
					String: func() string { return "not_found" },
				},
			},
			expectedErr: storage.ErrAndroidResourceNotFound,
		},
		{
			name: "returns nil error",
			fileSave: &fileSave{
				uri: &url.URL{
					String: func() string { return "found" },
				},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := nativeFileSave(test.fileSave)
			assert.Equal(t, test.expectedErr, err, "unexpected error")
		})
	}
}

func TestNativeFileOpen(t *testing.T) {
	tests := []struct {
		name        string
		uri         *url.URL
		expectedErr error
	}{
		{
			name: "returns resource not found error",
			uri: &url.URL{
				String: func() string { return "not_found" },
			},
			expectedErr: storage.ErrAndroidResourceNotFound,
		},
		{
			name: "returns nil error",
			uri: &url.URL{
				String: func() string { return "found" },
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := nativeFileOpen(test.uri)
			assert.Equal(t, test.expectedErr, err, "unexpected error")
		})
	}
}
