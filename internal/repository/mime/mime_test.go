package mime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitMimeType(t *testing.T) {
	main, sub := Split("text/plain")
	assert.Equal(t, "text", main)
	assert.Equal(t, "plain", sub)

	main, sub = Split("text/")
	assert.Empty(t, main)
	assert.Empty(t, sub)

	main, sub = Split("")
	assert.Empty(t, main)
	assert.Empty(t, sub)
}
