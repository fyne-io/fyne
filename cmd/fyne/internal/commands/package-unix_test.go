package commands

import (
	"bytes"
	"testing"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatDesktopFileList(t *testing.T) {
	assert.Equal(t, "", formatDesktopFileList([]string{}))
	assert.Equal(t, "One;", formatDesktopFileList([]string{"One"}))
	assert.Equal(t, "One;Two;", formatDesktopFileList([]string{"One", "Two"}))
}

func TestDesktopFileSource(t *testing.T) {
	tplData := unixData{
		Name: "Testing",
	}
	buf := &bytes.Buffer{}

	err := templates.DesktopFileUNIX.Execute(buf, tplData)
	require.NoError(t, err)
	assert.NotContains(t, buf.String(), "[X-Fyne")

	tplData.SourceRepo = "https://example.com"
	tplData.SourceDir = "cmd/name"

	err = templates.DesktopFileUNIX.Execute(buf, tplData)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "[X-Fyne")
	assert.Contains(t, buf.String(), "Repo=https://example.com")
	assert.Contains(t, buf.String(), "Dir=cmd/name")
}
