package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/internal"
)

func TestStore_RootURI(t *testing.T) {
	id := "io.fyne.test"
	a := &fyneApp{uniqueID: id}
	d := makeStoreDocs(id, &internal.InMemoryPreferences{}, &store{a: a})

	w, err := d.Create("test")
	assert.Nil(t, err)
	err = w.Close()
	assert.Nil(t, err)
	err = d.Remove("test")
	assert.Nil(t, err)
}
