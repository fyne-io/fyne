package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsMobile(t *testing.T) {
	osList := []string{"darwin", "linux", "windows", "android/arm", "android/arm64", "android/amd64", "android/386", "android", "ios"}
	osIsMobile := []bool{false, false, false, true, true, true, true, true, true}

	assert.False(t, IsMobile(""))

	for i, os := range osList {
		assert.Equal(t, IsMobile(os), osIsMobile[i])
	}

	assert.False(t, IsMobile("amiga"))
}

func Test_IsAndroid(t *testing.T) {
	osList := []string{"darwin", "linux", "windows", "android/arm", "android/arm64", "android/amd64", "android/386", "android", "ios"}
	osIsMobile := []bool{false, false, false, true, true, true, true, true, false}

	assert.False(t, IsAndroid(""))

	for i, os := range osList {
		assert.Equal(t, IsAndroid(os), osIsMobile[i])
	}

	assert.False(t, IsAndroid("amiga"))
}
