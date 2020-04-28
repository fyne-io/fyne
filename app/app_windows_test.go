// +build !ci

// +build !android,!ios

package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeNotifyString(t *testing.T) {
	assert.Equal(t, "normal", escapeNotifyString("normal"))
	assert.Equal(t, "isn't", escapeNotifyString("isn't"))
	assert.Equal(t, "`\"mine`\"", escapeNotifyString(`"mine"`))
	assert.Equal(t, "sla\\sh", escapeNotifyString(`sla\sh`))
	assert.Equal(t, "qu``ote", escapeNotifyString("qu`ote"))
	assert.Equal(t, "escaped ```\" string", escapeNotifyString("escaped `\" string"))
}
