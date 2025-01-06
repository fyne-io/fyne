package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/storage"
)

func TestMobileFilter(t *testing.T) {
	f := mobileFilter(nil)
	assert.Len(t, f.Extensions, 0)
	assert.Len(t, f.MimeTypes, 0)

	f = mobileFilter(storage.NewExtensionFileFilter([]string{".png"}))
	assert.Len(t, f.Extensions, 1)
	assert.Len(t, f.MimeTypes, 0)

	f = mobileFilter(storage.NewMimeTypeFileFilter([]string{"text/plain"}))
	assert.Len(t, f.Extensions, 0)
	assert.Len(t, f.MimeTypes, 1)
}
