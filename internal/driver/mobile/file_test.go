package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/storage"
)

func TestMobileFilter(t *testing.T) {
	f := mobileFilter(nil)
	assert.Empty(t, f.Extensions)
	assert.Empty(t, f.MimeTypes)

	f = mobileFilter(storage.NewExtensionFileFilter([]string{".png"}))
	assert.Len(t, f.Extensions, 1)
	assert.Empty(t, f.MimeTypes)

	f = mobileFilter(storage.NewMimeTypeFileFilter([]string{"text/plain"}))
	assert.Empty(t, f.Extensions)
	assert.Len(t, f.MimeTypes, 1)
}
