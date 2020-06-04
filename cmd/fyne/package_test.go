package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calculateExeName(t *testing.T) {
	modulesApp := calculateExeName("testdata/modules_app", "windows")
	assert.Equal(t, "module.exe", modulesApp)

	modulesShortName := calculateExeName("testdata/short_module", "linux")
	assert.Equal(t, "app", modulesShortName)

	nonModulesApp := calculateExeName("testdata", "linux")
	assert.Equal(t, "testdata", nonModulesApp)
}

func Test_isMobile(t *testing.T) {
	oslist := []string{"darwin", "linux", "windows", "android/arm", "android/386", "android", "ios"}
	osIsMobile := []bool{false, false, false, true, true, true, true}
	p := &packager{os: "", install: false, srcDir: ""}

	for i, os := range oslist {
		p.os = os
		assert.Equal(t, p.isMobile(), osIsMobile[i])
	}
}
