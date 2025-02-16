package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore_RootURI(t *testing.T) {
	id := "io.fyne.test"
	a := &fyneApp{uniqueID: id}
	d := makeStoreDocs(id, &store{a: a})

	w, err := d.Create("test")
	require.NoError(t, err)
	err = w.Close()
	require.NoError(t, err)
	err = d.Remove("test")
	require.NoError(t, err)
}
