package commands

import (
	"testing"

	"github.com/mcuadros/go-version"
	"github.com/stretchr/testify/assert"
)

func Test_CheckVersionTableTests(t *testing.T) {
	tests := map[string]struct {
		goVersion  string
		constraint string
		expectNil  bool
	}{
		"Windows 1.17.6": {goVersion: "go version go1.17.6 windows/amd64", constraint: ">=1.17", expectNil: true},
		"Linux 1.17.6":   {goVersion: "go version go1.17.6 linux/amd64", constraint: ">=1.17", expectNil: true},
		"Linux 1.18.0":   {goVersion: "go version go1.18.0 linux/amd64", constraint: ">=1.17", expectNil: true},
		"Windows 1.14.0": {goVersion: "go version go1.14.0 windows/amd64", constraint: ">=1.17", expectNil: false},
		"Windows 1.15.0": {goVersion: "go version go1.15.0 windows/amd64", constraint: ">=1.17", expectNil: false},
		"Windows 1.16.0": {goVersion: "go version go1.16.0 windows/amd64", constraint: ">=1.17", expectNil: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectNil {
				assert.Nil(t, checkVersion(tc.goVersion, version.NewConstrainGroupFromString(tc.constraint)))
			} else {
				assert.NotNil(t, checkVersion(tc.goVersion, version.NewConstrainGroupFromString(tc.constraint)))
			}
		})
	}
}
