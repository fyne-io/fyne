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
	osList := []string{"darwin", "linux", "windows", "android/arm", "android/386", "android", "ios"}
	osIsMobile := []bool{false, false, false, true, true, true, true}
	p := &packager{os: "", install: false, srcDir: ""}

	assert.Equal(t, p.isMobile(), false)

	for i, os := range osList {
		p.os = os
		assert.Equal(t, p.isMobile(), osIsMobile[i])
	}

	p.os = "amiga"
	assert.Equal(t, p.isMobile(), false)
}
