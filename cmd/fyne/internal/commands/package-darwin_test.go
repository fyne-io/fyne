package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDarwinLangs(t *testing.T) {
	langs := darwinLangs([]string{"en-GB", "de-CH"})

	assert.Equal(t, 2, len(langs))
	assert.Equal(t, "en_GB", langs[0])
	assert.Equal(t, "de_CH", langs[1])
}
