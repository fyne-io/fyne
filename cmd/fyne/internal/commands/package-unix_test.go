package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDesktopFileList(t *testing.T) {
	assert.Equal(t, "", formatDesktopFileList([]string{}))
	assert.Equal(t, "One;", formatDesktopFileList([]string{"One"}))
	assert.Equal(t, "One;Two;", formatDesktopFileList([]string{"One", "Two"}))
}
