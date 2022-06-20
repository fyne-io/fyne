package commands

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/josephspurrier/goversioninfo"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/cmd/fyne/internal/metadata"
)

func Test_calculateExeName(t *testing.T) {
	modulesApp := calculateExeName("testdata/modules_app", "windows")
	assert.Equal(t, "module.exe", modulesApp)

	modulesShortName := calculateExeName("testdata/short_module", "linux")
	assert.Equal(t, "app", modulesShortName)

	nonModulesApp := calculateExeName("testdata", "linux")
	assert.Equal(t, "testdata", nonModulesApp)
}

func Test_fixedVersionInfo(t *testing.T) {
	tests := []struct {
		ver   string
		fixed goversioninfo.FileVersion
	}{
		{"", goversioninfo.FileVersion{Major: 0, Minor: 0, Patch: 0, Build: 1}},
		{"1.1.1.1", goversioninfo.FileVersion{Major: 1, Minor: 1, Patch: 1, Build: 1}},
		{"2.2.2", goversioninfo.FileVersion{Major: 2, Minor: 2, Patch: 2, Build: 1}},
		{"3.3.3.3.3", goversioninfo.FileVersion{Major: 3, Minor: 3, Patch: 3, Build: 3}},
	}

	for _, tt := range tests {
		parsed := fixedVersionInfo(tt.ver)
		assert.Equal(t, tt.fixed, parsed)
	}
}

func Test_isValidVersion(t *testing.T) {
	assert.True(t, isValidVersion("1"))
	assert.True(t, isValidVersion("1.2"))
	assert.True(t, isValidVersion("1.2.3"))

	assert.False(t, isValidVersion("1.2.3.4"))
	assert.False(t, isValidVersion(""))
	assert.False(t, isValidVersion("1.2-alpha3"))
	assert.False(t, isValidVersion("pre1"))
	assert.False(t, isValidVersion("1..2"))
}

func Test_MergeMetata(t *testing.T) {
	p := &Packager{appData: &appData{}}
	p.AppVersion = "v0.1"
	data := &metadata.FyneApp{
		Details: metadata.AppDetails{
			Icon:    "test.png",
			Build:   3,
			Version: "v0.0.1",
		},
	}

	mergeMetadata(p.appData, data)
	assert.Equal(t, "v0.1", p.AppVersion)
	assert.Equal(t, 3, p.AppBuild)
	assert.Equal(t, "test.png", p.icon)
}

func Test_validateAppID(t *testing.T) {
	id, err := validateAppID("myApp", "windows", "myApp", false)
	assert.Nil(t, err)
	assert.Equal(t, "myApp", id)

	id, err = validateAppID("", "darwin", "myApp", true)
	assert.Nil(t, err)
	assert.Equal(t, "com.example.myApp", id) // this was in for compatibility

	id, err = validateAppID("com.myApp", "darwin", "myApp", true)
	assert.Nil(t, err)
	assert.Equal(t, "com.myApp", id)

	_, err = validateAppID("", "ios", "myApp", false)
	assert.NotNil(t, err)

	_, err = validateAppID("myApp", "ios", "myApp", false)
	assert.NotNil(t, err)

	id, err = validateAppID("com.myApp", "android", "myApp", true)
	assert.Nil(t, err)
	assert.Equal(t, "com.myApp", id)

	_, err = validateAppID("myApp", "android", "myApp", true)
	assert.NotNil(t, err)
}

func Test_buildPackageWasm(t *testing.T) {
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

	p := &Packager{
		appData: &appData{},
		os:      "wasm",
		srcDir:  "myTest",
		release: true,
	}
	wasmBuildTest := &testCommandRuns{runs: expected, t: t}
	files, err := p.buildPackage(wasmBuildTest)
	assert.Nil(t, err)
	assert.NotNil(t, files)
	assert.Equal(t, 1, len(files))
}

func Test_PackageWasm(t *testing.T) {
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
				args:  []string{"build", "-o", "myTest.wasm"},
				env:   []string{"GOARCH=wasm", "GOOS=js"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
	}

	p := &Packager{
		appData: &appData{
			Name: "myTest",
			icon: "myTest.png",
		},
		os:     "wasm",
		srcDir: "myTest",
		dir:    "myTestTarget",
		exe:    "myTest.wasm",
	}
	wasmBuildTest := &testCommandRuns{runs: expected, t: t}

	util = mockUtil{}

	utilIsMobileMock = func(_ string) bool {
		return false
	}

	expectedEnsureSubDirRuns := mockEnsureSubDirRuns{
		expected: []mockEnsureSubDir{
			{"myTestTarget", "wasm", "myTestTarget/wasm"},
		},
	}
	utilEnsureSubDirMock = func(parent, name string) string {
		return expectedEnsureSubDirRuns.verifyExpectation(t, parent, name)
	}

	expectedExistRuns := mockExistRuns{
		expected: []mockExist{
			{"myTest.wasm", false},
			{"myTest.wasm", true},
		},
	}
	utilExistsMock = func(path string) bool {
		return expectedExistRuns.verifyExpectation(t, path)
	}

	expectedWriteFileRuns := mockWriteFileRuns{
		expected: []mockWriteFile{
			{filepath.Join("myTestTarget", "wasm", "index.html"), nil},
			{filepath.Join("myTestTarget", "wasm", "spinner_light.gif"), nil},
			{filepath.Join("myTestTarget", "wasm", "spinner_dark.gif"), nil},
			{filepath.Join("myTestTarget", "wasm", "light.css"), nil},
			{filepath.Join("myTestTarget", "wasm", "dark.css"), nil},
			{filepath.Join("myTestTarget", "wasm", "webgl-debug.js"), nil},
		},
	}
	utilWriteFileMock = func(target string, _ []byte) error {
		return expectedWriteFileRuns.verifyExpectation(t, target)
	}

	expectedCopyFileRuns := mockCopyFileRuns{
		expected: []mockCopyFile{
			{source: "myTest.png", target: filepath.Join("myTestTarget", "wasm", "icon.png")},
			{source: filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"), target: filepath.Join("myTestTarget", "wasm", "wasm_exec.js")},
			{source: "myTest.wasm", target: filepath.Join("myTestTarget", "wasm", "myTest.wasm")},
		},
	}
	utilCopyFileMock = func(source, target string) error {
		return expectedCopyFileRuns.verifyExpectation(t, false, source, target)
	}

	err := p.doPackage(wasmBuildTest)
	assert.Nil(t, err)
	wasmBuildTest.verifyExpectation()
	expectedTotalCount(t, len(expectedEnsureSubDirRuns.expected), expectedEnsureSubDirRuns.current)
	expectedTotalCount(t, len(expectedExistRuns.expected), expectedExistRuns.current)
	expectedTotalCount(t, len(expectedWriteFileRuns.expected), expectedWriteFileRuns.current)
	expectedTotalCount(t, len(expectedCopyFileRuns.expected), expectedCopyFileRuns.current)
}

func Test_buildPackageGopherJS(t *testing.T) {
	expected := []mockRunner{
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
				args:  []string{"build", "-o", "myTest.js", "--tags", "release"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
	}

	p := &Packager{
		appData: &appData{},
		os:      "gopherjs",
		srcDir:  "myTest",
		exe:     "myTest.js",
		release: true,
	}
	wasmBuildTest := &testCommandRuns{runs: expected, t: t}
	files, err := p.buildPackage(wasmBuildTest)
	assert.Nil(t, err)
	assert.NotNil(t, files)
	assert.Equal(t, 1, len(files))
}

func Test_PackageGopherJS(t *testing.T) {
	expected := []mockRunner{
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
				args: []string{"build",
					"-o", "myTest.js"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
	}

	p := &Packager{
		appData: &appData{
			Name: "myTest",
			icon: "myTest.png",
		},
		os:     "gopherjs",
		srcDir: "myTest",
		dir:    "myTestTarget",
		exe:    "myTest.js",
	}
	gopherjsBuildTest := &testCommandRuns{runs: expected, t: t}

	util = mockUtil{}

	utilIsMobileMock = func(_ string) bool {
		return false
	}

	expectedEnsureSubDirRuns := mockEnsureSubDirRuns{
		expected: []mockEnsureSubDir{
			{"myTestTarget", "gopherjs", "myTestTarget/gopherjs"},
		},
	}
	utilEnsureSubDirMock = func(parent, name string) string {
		return expectedEnsureSubDirRuns.verifyExpectation(t, parent, name)
	}

	expectedExistRuns := mockExistRuns{
		expected: []mockExist{
			{"myTest.js", false},
			{"myTest.js", true},
		},
	}
	utilExistsMock = func(path string) bool {
		return expectedExistRuns.verifyExpectation(t, path)
	}

	expectedWriteFileRuns := mockWriteFileRuns{
		expected: []mockWriteFile{
			{filepath.Join("myTestTarget", "gopherjs", "index.html"), nil},
			{filepath.Join("myTestTarget", "gopherjs", "spinner_light.gif"), nil},
			{filepath.Join("myTestTarget", "gopherjs", "spinner_dark.gif"), nil},
			{filepath.Join("myTestTarget", "gopherjs", "light.css"), nil},
			{filepath.Join("myTestTarget", "gopherjs", "dark.css"), nil},
			{filepath.Join("myTestTarget", "gopherjs", "webgl-debug.js"), nil},
		},
	}
	utilWriteFileMock = func(target string, _ []byte) error {
		return expectedWriteFileRuns.verifyExpectation(t, target)
	}

	expectedCopyFileRuns := mockCopyFileRuns{
		expected: []mockCopyFile{
			{source: "myTest.png", target: filepath.Join("myTestTarget", "gopherjs", "icon.png")},
			{source: "myTest.js", target: filepath.Join("myTestTarget", "gopherjs", "myTest.js")},
		},
	}
	utilCopyFileMock = func(source, target string) error {
		return expectedCopyFileRuns.verifyExpectation(t, false, source, target)
	}

	err := p.doPackage(gopherjsBuildTest)
	assert.Nil(t, err)
	gopherjsBuildTest.verifyExpectation()
	expectedTotalCount(t, len(expectedEnsureSubDirRuns.expected), expectedEnsureSubDirRuns.current)
	expectedTotalCount(t, len(expectedExistRuns.expected), expectedExistRuns.current)
	expectedTotalCount(t, len(expectedWriteFileRuns.expected), expectedWriteFileRuns.current)
	expectedTotalCount(t, len(expectedCopyFileRuns.expected), expectedCopyFileRuns.current)
}

func Test_BuildPackageWeb(t *testing.T) {
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
				args:  []string{"build", "-o", "myTest.wasm", "-tags", "release"},
				env:   []string{"GOARCH=wasm", "GOOS=js"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
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
				args:  []string{"build", "-o", "myTest.js", "--tags", "release"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
	}

	p := &Packager{
		appData: &appData{},
		os:      "web",
		srcDir:  "myTest",
		release: true,
		exe:     "myTest",
	}
	webBuildTest := &testCommandRuns{runs: expected, t: t}
	files, err := p.buildPackage(webBuildTest)
	assert.Nil(t, err)
	assert.NotNil(t, files)
	assert.Equal(t, 2, len(files))
}

func Test_PackageWeb(t *testing.T) {
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
				args:  []string{"build", "-o", "myTest.wasm"},
				env:   []string{"GOARCH=wasm", "GOOS=js"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
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
				args: []string{"build",
					"-o", "myTest.js"},
				osEnv: true,
				dir:   "myTest",
			},
			mockReturn: mockReturn{
				ret: []byte(""),
			},
		},
	}

	p := &Packager{
		appData: &appData{},
		os:      "web",
		srcDir:  "myTest",
		dir:     "myTestTarget",
		exe:     "myTest",
	}
	p.Name = "myTest"
	p.icon = "myTest.png"
	gopherjsBuildTest := &testCommandRuns{runs: expected, t: t}

	util = mockUtil{}

	utilIsMobileMock = func(_ string) bool {
		return false
	}

	expectedEnsureSubDirRuns := mockEnsureSubDirRuns{
		expected: []mockEnsureSubDir{
			{"myTestTarget", "web", "myTestTarget/web"},
		},
	}
	utilEnsureSubDirMock = func(parent, name string) string {
		return expectedEnsureSubDirRuns.verifyExpectation(t, parent, name)
	}

	expectedExistRuns := mockExistRuns{
		expected: []mockExist{
			{"myTest", false},
		},
	}
	utilExistsMock = func(path string) bool {
		return expectedExistRuns.verifyExpectation(t, path)
	}

	expectedWriteFileRuns := mockWriteFileRuns{
		expected: []mockWriteFile{
			{filepath.Join("myTestTarget", "web", "index.html"), nil},
			{filepath.Join("myTestTarget", "web", "spinner_light.gif"), nil},
			{filepath.Join("myTestTarget", "web", "spinner_dark.gif"), nil},
			{filepath.Join("myTestTarget", "web", "light.css"), nil},
			{filepath.Join("myTestTarget", "web", "dark.css"), nil},
			{filepath.Join("myTestTarget", "web", "webgl-debug.js"), nil},
		},
	}
	utilWriteFileMock = func(target string, _ []byte) error {
		return expectedWriteFileRuns.verifyExpectation(t, target)
	}

	expectedCopyFileRuns := mockCopyFileRuns{
		expected: []mockCopyFile{
			{source: "myTest.png", target: filepath.Join("myTestTarget", "web", "icon.png")},
			{source: "myTest.js", target: filepath.Join("myTestTarget", "web", "myTest.js")},
			{source: filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"), target: filepath.Join("myTestTarget", "web", "wasm_exec.js")},
			{source: "myTest.wasm", target: filepath.Join("myTestTarget", "web", "myTest.wasm")},
		},
	}
	utilCopyFileMock = func(source, target string) error {
		return expectedCopyFileRuns.verifyExpectation(t, false, source, target)
	}

	err := p.doPackage(gopherjsBuildTest)
	assert.Nil(t, err)
	gopherjsBuildTest.verifyExpectation()
	expectedTotalCount(t, len(expectedEnsureSubDirRuns.expected), expectedEnsureSubDirRuns.current)
	expectedTotalCount(t, len(expectedExistRuns.expected), expectedExistRuns.current)
	expectedTotalCount(t, len(expectedWriteFileRuns.expected), expectedWriteFileRuns.current)
	expectedTotalCount(t, len(expectedCopyFileRuns.expected), expectedCopyFileRuns.current)
}
