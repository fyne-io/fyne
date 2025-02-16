package metadata

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAppMetadata(t *testing.T) {
	r, err := os.Open("./testdata/FyneApp.toml")
	require.NoError(t, err)
	data, err := Load(r)
	_ = r.Close()
	require.NoError(t, err)
	assert.Equal(t, 1, data.Details.Build)

	data.Details.Build++

	versionPath := "./testdata/Version.toml"
	w, err := os.Create(versionPath)
	require.NoError(t, err)
	err = Save(data, w)
	require.NoError(t, err)
	defer func() {
		os.Remove(versionPath)
	}()
	_ = w.Close()

	r, err = os.Open(versionPath)
	require.NoError(t, err)
	defer r.Close()

	data2, err := Load(r)
	require.NoError(t, err)
	assert.Equal(t, 2, data2.Details.Build)
}
