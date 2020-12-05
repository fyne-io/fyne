package desktop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMouseButton_Deprecation(t *testing.T) {
	assert.Equal(t, MouseButtonPrimary, LeftMouseButton)
	assert.Equal(t, MouseButtonSecondary, RightMouseButton)

	// and just check we're not accidentally adding consts
	assert.NotEqual(t, MouseButtonTertiary, MouseButtonPrimary|MouseButtonSecondary)
}
