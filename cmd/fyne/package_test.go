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

type numericVersionTest struct {
	Version string
	Min     int
	Max     int
	Expect  string
}

type numericVersionErrorTest struct {
	Version string
	Min     int
	Max     int
}

func TestNumericVersion(t *testing.T) {
	versionTests := []numericVersionTest{
		{"2", 1, 1, "2"},
		{"2", 2, 3, "2.0"},
		{"2", 3, 3, "2.0.0"},
		{"2", 4, 4, "2.0.0.0"},
		{"2", 4, 100, "2.0.0.0"},
		{"1.2.3", 1, 3, "1.2.3"},
		{"1.2", 3, 3, "1.2.0"},
		{"1.2", 0, 100, "1.2"},
		{"1.2.3", 0, 100, "1.2.3"},
	}
	for _, v := range versionTests {
		res, err := buildNumericVersion(v.Version, v.Min, v.Max)
		if err != nil {
			t.Error(err)
			continue
		}
		if res != v.Expect {
			t.Errorf("error building version: expected %q but got %q", v.Expect, res)
		}
	}
}

func TestNumericVersionErrors(t *testing.T) {
	versionTests := []numericVersionErrorTest{
		{"", 0, 0},
		{"", 1, 1},
		{"", 0, 1},
		{"2.3.4", 2, 2},
		{"2.3.0-alpha", 2, 3},
		{"beta", 1, 1},
		{"1.2.x", 1, 1},
	}
	for _, v := range versionTests {
		res, err := buildNumericVersion(v.Version, v.Min, v.Max)
		if err == nil {
			t.Errorf("expecting an error when parsing (%q, [%d, %d]), but got %q instead", v.Version, v.Min, v.Max, res)
		}
	}
}
