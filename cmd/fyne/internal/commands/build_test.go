package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mcuadros/go-version"
	"github.com/stretchr/testify/assert"
)

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
			expectedValue: expectedValue{args: []string{"mod", "edit", "-json"}},
			mockReturn: mockReturn{
				ret: []byte("{ \"Module\": { \"Path\": \"fyne.io/fyne/v2\"} }"),
			},
		},
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
	b := &Builder{appData: &appData{}, os: "wasm", srcdir: "myTest", runner: wasmBuildTest}
	err := b.build()
	assert.Nil(t, err)
	wasmBuildTest.verifyExpectation()
}

func Test_BuildWasmReleaseVersion(t *testing.T) {
	expected := []mockRunner{
		{
			expectedValue: expectedValue{args: []string{"mod", "edit", "-json"}},
			mockReturn: mockReturn{
				ret: []byte("{ \"Module\": { \"Path\": \"fyne.io/fyne/v2\"} }"),
			},
		},
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
	b := &Builder{appData: &appData{}, os: "wasm", srcdir: "myTest", release: true, runner: wasmBuildTest}
	err := b.build()
	assert.Nil(t, err)
	wasmBuildTest.verifyExpectation()
}

func Test_BuildGopherJSReleaseVersion(t *testing.T) {
	expected := []mockRunner{}

	if runtime.GOOS != "windows" {
		expected = []mockRunner{
			{
				expectedValue: expectedValue{args: []string{"mod", "edit", "-json"}},
				mockReturn: mockReturn{
					ret: []byte("{ \"Module\": { \"Path\": \"fyne.io/fyne/v2\"} }"),
				},
			},
			{
				expectedValue: expectedValue{
					args:  []string{"version"},
					osEnv: true,
				},
				mockReturn: mockReturn{
					ret: []byte(""),
					err: nil,
				},
			},
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
	}

	gopherJSBuildTest := &testCommandRuns{runs: expected, t: t}
	b := &Builder{appData: &appData{}, os: "gopherjs", srcdir: "myTest", release: true, runner: gopherJSBuildTest}
	err := b.build()
	if runtime.GOOS == "windows" {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
	}
	gopherJSBuildTest.verifyExpectation()
}

func Test_BuildWasmOldVersion(t *testing.T) {
	expected := []mockRunner{
		{
			expectedValue: expectedValue{args: []string{"mod", "edit", "-json"}},
			mockReturn: mockReturn{
				ret: []byte("{ \"Module\": { \"Path\": \"fyne.io/fyne/v2\"} }"),
			},
		},
		{
			expectedValue: expectedValue{args: []string{"version"}},
			mockReturn:    mockReturn{ret: []byte("go version go1.16.0 windows/amd64")},
		},
	}

	wasmBuildTest := &testCommandRuns{runs: expected, t: t}
	b := &Builder{appData: &appData{}, os: "wasm", srcdir: "myTest", runner: wasmBuildTest}
	err := b.build()
	assert.NotNil(t, err)
	wasmBuildTest.verifyExpectation()
}

func Test_BuildLinuxReleaseVersion(t *testing.T) {
	relativePath := "." + string(os.PathSeparator) + filepath.Join("cmd", "terminal")

	expected := []mockRunner{
		{
			expectedValue: expectedValue{args: []string{"mod", "edit", "-json"}},
			mockReturn: mockReturn{
				ret: []byte("{ \"Module\": { \"Path\": \"fyne.io/fyne/v2\"} }"),
			},
		},
		{
			expectedValue: expectedValue{
				args:  []string{"build", "-trimpath", "-ldflags", "-s -w", "-tags", "release", relativePath},
				env:   []string{"CGO_ENABLED=1", "GOOS=linux"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
	}

	linuxBuildTest := &testCommandRuns{runs: expected, t: t}
	b := &Builder{appData: &appData{}, os: "linux", srcdir: "myTest", release: true, runner: linuxBuildTest, goPackage: relativePath}
	err := b.build()
	assert.Nil(t, err)
	linuxBuildTest.verifyExpectation()
}

type jsonTest struct {
	expected bool
	json     []byte
}

func Test_FyneGoMod(t *testing.T) {
	jsonTests := []jsonTest{
		{false, []byte(`{"Module": {"Path": "github.com/fyne-io/calculator"},"Go": "1.14",	"Require": [ { "Path": "fyne.io/fyne/v2","Version": "v2.1.4"} ] }`)},
		{true, []byte(`{ "Module": {"Path": "fyne.io/fyne/v2"},"Require": [{ "Path": "test","Version": "v2.1.4"} ] }`)},
		{true, []byte(`{"Module": {"Path": "github.com/fyne-io/calculator"},"Go": "1.14",	"Require": [ { "Path": "fyne.io/fyne/v2","Version": "v2.2.0"} ] }`)},
	}

	for _, j := range jsonTests {
		expected := []mockRunner{
			{
				expectedValue: expectedValue{args: []string{"mod", "edit", "-json"}},
				mockReturn:    mockReturn{ret: j.json},
			},
		}

		called := false

		fyneGoModTest := &testCommandRuns{runs: expected, t: t}
		injectMetadataIfPossible(fyneGoModTest, "myTest", &appData{},
			func(string, *appData) (func(), error) {
				called = true
				return func() {}, nil
			})

		assert.Equal(t, j.expected, called)
	}
}

func Test_AppendEnv(t *testing.T) {
	env := []string{
		"foo=bar",
		"bar=baz",
		"foo1=bar=baz",
	}

	appendEnv(&env, "foo2", "baz2")
	appendEnv(&env, "foo", "baz")
	appendEnv(&env, "foo1", "-bar")

	if assert.Len(t, env, 4) {
		assert.Equal(t, "foo=bar baz", env[0])
		assert.Equal(t, "bar=baz", env[1])
		assert.Equal(t, "foo1=bar=baz -bar", env[2])
		assert.Equal(t, "foo2=baz2", env[3])
	}
}

type extractTest struct {
	value       string
	wantLdFlags string
	wantGoFlags string
}

func Test_ExtractLdFlags(t *testing.T) {
	goFlagsTests := []extractTest{
		{"-ldflags=-w", "-w", ""},
		{"-ldflags=-s", "-s", ""},
		{"-ldflags=-w -ldflags=-s", "-w -s", ""},
		{"-mod=vendor", "", "-mod=vendor"},
	}

	for _, test := range goFlagsTests {
		ldFlags, goFlags := extractLdFlags(test.value)
		assert.Equal(t, test.wantLdFlags, ldFlags)
		assert.Equal(t, test.wantGoFlags, goFlags)
	}
}

func Test_NormaliseVersion(t *testing.T) {
	assert.Equal(t, "master", normaliseVersion("master"))
	assert.Equal(t, "2.3.0.0", normaliseVersion("v2.3"))
	assert.Equal(t, "2.4.0.0", normaliseVersion("v2.4.0"))
	assert.Equal(t, "2.3.6.0-dev", normaliseVersion("v2.3.6-0.20230711180435-d4b95e1cb1eb"))
}
