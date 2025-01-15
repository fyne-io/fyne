package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_RootURI(t *testing.T) {
	id := "io.fyne.test"
	a := &fyneApp{uniqueID: id}
	d := makeStoreDocs(id, &store{a: a})

	w, err := d.Create("test")
	assert.NoError(t, err)
	err = w.Close()
	assert.NoError(t, err)
	err = d.Remove("test")
	assert.NoError(t, err)
}
