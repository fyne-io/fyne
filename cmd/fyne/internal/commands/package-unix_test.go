package commands

import (
	"bytes"
	"strings"
	"testing"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)
	assert.False(t, strings.Contains(buf.String(), "[X-Fyne"))

	tplData.SourceRepo = "https://example.com"
	tplData.SourceDir = "cmd/name"

	err = templates.DesktopFileUNIX.Execute(buf, tplData)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(buf.String(), "[X-Fyne"))
	assert.True(t, strings.Contains(buf.String(), "Repo=https://example.com"))
	assert.True(t, strings.Contains(buf.String(), "Dir=cmd/name"))
}
