package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BuildWasmVersion(t *testing.T) {
	expected := []mockRunner{
		{
			expectedValue: expectedValue{args: []string{"mod", "edit", "-json"}},
			mockReturn: mockReturn{
				ret: []byte("{ \"Module\": { \"Path\": \"fyne.io/fyne/v2\"} }"),
			},
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
	assert.Equal(t, "2.4.1.0-dev", normaliseVersion("v2.4.1-rc7.0.20230711180435-d4b95e1cb1eb"))
}
