package metadata

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAppMetadata(t *testing.T) {
	r, err := os.Open("./testdata/FyneApp.toml")
	require.NoError(t, err)
	defer r.Close()

	data, err := Load(r)
	require.NoError(t, err)
	assert.Equal(t, "https://apps.fyne.io", data.Website)
	assert.Equal(t, "io.fyne.fyne", data.Details.ID)
	assert.Equal(t, "1.0.0", data.Details.Version)
	assert.Equal(t, 1, data.Details.Build)
	assert.Equal(t, "Value1", data.Release["Test"])
	assert.Equal(t, "Value3", data.Release["InReleaseOnly"])
	assert.NotContains(t, data.Release, "InDevelopmentOnly")
	assert.Equal(t, "Value2", data.Development["Test"])
	assert.Equal(t, "Value4", data.Development["InDevelopmentOnly"])
	assert.NotContains(t, data.Development, "InReleaseOnly")

	assert.NotNil(t, data.Source)
	assert.Equal(t, "https://github.com/fyne-io/fyne", data.Source.Repo)
	assert.Equal(t, "internal/metadata/testdata", data.Source.Dir)
}
