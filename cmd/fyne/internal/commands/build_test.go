package commands

import (
	"fmt"
	"testing"

	"github.com/mcuadros/go-version"
	"github.com/stretchr/testify/assert"
)

func TestBuildGenerateMetaLDFlags(t *testing.T) {
	b := &builder{}
	assert.Equal(t, "", b.generateMetaLDFlags())

	b.id = "com.example"
	assert.Equal(t, "-X 'fyne.io/fyne/v2/internal/app.MetaID=com.example'", b.generateMetaLDFlags())

	b.version = "1.2.3"
	assert.Equal(t, "-X 'fyne.io/fyne/v2/internal/app.MetaID=com.example' -X 'fyne.io/fyne/v2/internal/app.MetaVersion=1.2.3'", b.generateMetaLDFlags())
}

func Test_CheckGoVersionNoGo(t *testing.T) {
	commandNil := &testCommandRuns{runs: []mockRunner{}, t: t}
	assert.Nil(t, checkGoVersion(commandNil, nil))
	commandNil.verifyExpectation()

	commandNotNil := &testCommandRuns{runs: []mockRunner{{
		expectedValue: expectedValue{args: []string{"version"}},
		mockReturn:    mockReturn{err: fmt.Errorf("file not found")},
	}}, t: t}
	assert.NotNil(t, checkGoVersion(commandNotNil, version.NewConstrainGroupFromString(">=1.17")))
	commandNotNil.verifyExpectation()
}

func Test_CheckGoVersionValidValue(t *testing.T) {
	expected := []mockRunner{
		{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("go version go1.17.6 windows/amd64")},
		},
		{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("go version go1.17.6 linux/amd64")},
		}}

	versionOk := &testCommandRuns{runs: expected, t: t}
	assert.Nil(t, checkGoVersion(versionOk, version.NewConstrainGroupFromString(">=1.17")))
	assert.Nil(t, checkGoVersion(versionOk, version.NewConstrainGroupFromString(">=1.17")))
	versionOk.verifyExpectation()

	versionNotOk := &testCommandRuns{runs: expected, t: t}
	assert.NotNil(t, checkGoVersion(versionNotOk, version.NewConstrainGroupFromString("<1.17")))
	assert.NotNil(t, checkGoVersion(versionNotOk, version.NewConstrainGroupFromString("<1.17")))
	versionNotOk.verifyExpectation()
}

func Test_CheckGoVersionInvalidValue(t *testing.T) {
	tests := map[string]struct {
		runs []mockRunner
	}{
		"Wrong go command output": {runs: []mockRunner{{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("got noversion go1.17.6 windows/amd64")},
		}}},
		"Wrong go command length": {runs: []mockRunner{{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("go version really")},
		}}},
		"Wrong version string": {runs: []mockRunner{{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("go version 1.17.6 windows/amd64")},
		}}},
		"No version number": {runs: []mockRunner{{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("go version gowrong windows/amd64")},
		}}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			versionInvalid := &testCommandRuns{runs: tc.runs, t: t}
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
	expected := []mockRunner{
		{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("go version go1.17.6 windows/amd64")},
		},
		{
			expectedValue: expectedValue{
				args:  []string{"build"},
				env:   []string{"GOARCH=wasm", "GOOS=js"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{ret: []byte("")},
		},
	}

	wasmBuildTest := &testCommandRuns{runs: expected, t: t}
	b := &builder{os: "wasm", srcdir: "myTest", runner: wasmBuildTest}
	err := b.build()
	assert.Nil(t, err)
	wasmBuildTest.verifyExpectation()
}

func Test_BuildWasmReleaseVersion(t *testing.T) {
	expected := []mockRunner{
		{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn: mockReturn{
				ret: []byte("go version go1.17.6 windows/amd64"),
			},
		},
		{
			expectedValue: expectedValue{
				args:  []string{"build", "-tags", "release"},
				env:   []string{"GOARCH=wasm", "GOOS=js"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
	}

	wasmBuildTest := &testCommandRuns{runs: expected, t: t}
	b := &builder{os: "wasm", srcdir: "myTest", release: true, runner: wasmBuildTest}
	err := b.build()
	assert.Nil(t, err)
	wasmBuildTest.verifyExpectation()
}

func Test_BuildGopherJSReleaseVersion(t *testing.T) {
	expected := []mockRunner{
		{
			expectedValue: expectedValue{
				args:  []string{"build", "--tags", "release"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
	}

	gopherJSBuildTest := &testCommandRuns{runs: expected, t: t}
	b := &builder{os: "gopherjs", srcdir: "myTest", release: true, runner: gopherJSBuildTest}
	err := b.build()
	assert.Nil(t, err)
	gopherJSBuildTest.verifyExpectation()
}

func Test_BuildWasmOldVersion(t *testing.T) {
	expected := []mockRunner{
		{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("go version go1.16.0 windows/amd64")},
		},
	}

	wasmBuildTest := &testCommandRuns{runs: expected, t: t}
	b := &builder{os: "wasm", srcdir: "myTest", runner: wasmBuildTest}
	err := b.build()
	assert.NotNil(t, err)
	wasmBuildTest.verifyExpectation()
}
