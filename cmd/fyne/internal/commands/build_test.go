package commands

import (
	"fmt"
	"testing"

	"github.com/mcuadros/go-version"
	"github.com/stretchr/testify/assert"
)

func Test_CheckGoVersionNoGo(t *testing.T) {
	commandNil := &testCommandCall{calls: []mockRunOutputValue{}, t: t}
	assert.Nil(t, checkGoVersion(commandNil, nil))
	commandNil.verifyExpectation()

	commandNotNil := &testCommandCall{calls: []mockRunOutputValue{{
		args: []string{"version"},
		ret:  nil,
		err:  fmt.Errorf("file not found"),
	}}, t: t}
	assert.NotNil(t, checkGoVersion(commandNotNil, version.NewConstrainGroupFromString(">=1.17")))
	commandNotNil.verifyExpectation()
}

func Test_CheckGoVersionValidValue(t *testing.T) {
	expected := []mockRunOutputValue{
		{
			args: []string{"version"},
			ret:  []byte("go version go1.17.6 windows/amd64"),
			err:  nil,
		},
		{
			args: []string{"version"},
			ret:  []byte("go version go1.17.6 linux/amd64"),
			err:  nil,
		}}

	versionOk := &testCommandCall{calls: expected, t: t}
	assert.Nil(t, checkGoVersion(versionOk, version.NewConstrainGroupFromString(">=1.17")))
	assert.Nil(t, checkGoVersion(versionOk, version.NewConstrainGroupFromString(">=1.17")))
	versionOk.verifyExpectation()

	versionNotOk := &testCommandCall{calls: expected, t: t}
	assert.NotNil(t, checkGoVersion(versionNotOk, version.NewConstrainGroupFromString("<1.17")))
	assert.NotNil(t, checkGoVersion(versionNotOk, version.NewConstrainGroupFromString("<1.17")))
	versionNotOk.verifyExpectation()
}

func Test_CheckGoVersionInvalidValue(t *testing.T) {
	tests := map[string]struct {
		calls []mockRunOutputValue
	}{
		"Wrong go command output": {calls: []mockRunOutputValue{{
			args: []string{"version"},
			ret:  []byte("got noversion go1.17.6 windows/amd64"),
			err:  nil,
		}}},
		"Wrong go command length": {calls: []mockRunOutputValue{{
			args: []string{"version"},
			ret:  []byte("go version really"),
			err:  nil,
		}}},
		"Wrong version string": {calls: []mockRunOutputValue{{
			args: []string{"version"},
			ret:  []byte("go version 1.17.6 windows/amd64"),
			err:  nil,
		}}},
		"No version number": {calls: []mockRunOutputValue{{
			args: []string{"version"},
			ret:  []byte("go version gowrong windows/amd64"),
			err:  nil,
		}}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			versionInvalid := &testCommandCall{calls: tc.calls, t: t}
			assert.NotNil(t, checkGoVersion(versionInvalid, version.NewConstrainGroupFromString(">=1.17")))
			versionInvalid.verifyExpectation()
		})
	}
}

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

func Test_BuildWasmVersion(t *testing.T) {
	expected := []mockRunOutputValue{
		{
			args: []string{"version"},
			ret:  []byte("go version go1.17.6 windows/amd64"),
			err:  nil,
		},
		{
			args:  []string{"build"},
			ret:   []byte(""),
			err:   nil,
			env:   []string{"GOARCH=wasm", "GOOS=js"},
			osEnv: true,
			dir:   "myTest",
		},
	}

	wasmBuildTest := &testCommandCall{calls: expected, t: t}
	b := &builder{os: "wasm", srcdir: "myTest", runner: wasmBuildTest}
	assert.Nil(t, b.build())
	wasmBuildTest.verifyExpectation()
}

func Test_BuildWasmReleaseVersion(t *testing.T) {
	expected := []mockRunOutputValue{
		{
			args: []string{"version"},
			ret:  []byte("go version go1.17.6 windows/amd64"),
			err:  nil,
		},
		{
			args:  []string{"build", "-ldflags", "-s -w", "-tags", "release"},
			ret:   []byte(""),
			err:   nil,
			env:   []string{"GOARCH=wasm", "GOOS=js"},
			osEnv: true,
			dir:   "myTest",
		},
	}

	wasmBuildTest := &testCommandCall{calls: expected, t: t}
	b := &builder{os: "wasm", srcdir: "myTest", release: true, runner: wasmBuildTest}
	assert.Nil(t, b.build())
	wasmBuildTest.verifyExpectation()
}

func Test_BuildWasmOldVersion(t *testing.T) {
	expected := []mockRunOutputValue{
		{
			args: []string{"version"},
			ret:  []byte("go version go1.16.0 windows/amd64"),
			err:  nil,
		},
	}

	wasmBuildTest := &testCommandCall{calls: expected, t: t}
	b := &builder{os: "wasm", srcdir: "myTest", runner: wasmBuildTest}
	assert.NotNil(t, b.build())
	wasmBuildTest.verifyExpectation()
}
