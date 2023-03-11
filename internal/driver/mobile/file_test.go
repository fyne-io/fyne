package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/storage"
)

func TestMobileFilter(t *testing.T) {
	f := mobileFilter(nil)
	assert.Equal(t, 0, len(f.Extensions))
	assert.Equal(t, 0, len(f.MimeTypes))

	f = mobileFilter(storage.NewExtensionFileFilter([]string{".png"}))
	assert.Equal(t, 1, len(f.Extensions))
	assert.Equal(t, 0, len(f.MimeTypes))

	f = mobileFilter(storage.NewMimeTypeFileFilter([]string{"text/plain"}))
	assert.Equal(t, 0, len(f.Extensions))
	assert.Equal(t, 1, len(f.MimeTypes))
}
