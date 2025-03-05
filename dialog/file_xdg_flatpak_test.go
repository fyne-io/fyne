//go:build flatpak && !windows && !android && !ios && !wasm && !js

package dialog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFilterName(t *testing.T) {
	actual := formatFilterName([]string{"1", "2", "3"}, 4)
	assert.Equal(t, "1, 2, 3", actual)

	actual = formatFilterName([]string{"1", "2", "3"}, 3)
	assert.Equal(t, "1, 2, 3", actual)

	actual = formatFilterName([]string{"1", "2", "3"}, 2)
	assert.Equal(t, "1, 2…", actual)
}
