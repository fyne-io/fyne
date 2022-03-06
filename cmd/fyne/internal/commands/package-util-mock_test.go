package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var utilExistsMock func(path string) bool
var utilCopyFileMock func(source string, target string) error
var utilCopyExeFileMock func(src, tgt string) error
var utilWriteFileMock func(target string, data []byte) error
var utilEnsureSubDirMock func(parent, name string) string

var utilRequireAndroidSDKMock func() error
var utilAndroidBuildToolsPathMock func() string

var utilIsAndroidMock func(os string) bool
var utilIsIOSMock func(os string) bool
var utilIsMobileMock func(os string) bool

type mockUtil struct{}

func expectedTotalCount(t *testing.T, totalExpected int, totalProcessed int) {
	assert.Equal(t, totalExpected, totalProcessed)
}

type mockExist struct {
	path string
	ret  bool
}

type mockExistRuns struct {
	expected []mockExist

	current int
}

func (m *mockExistRuns) verifyExpectation(t *testing.T, path string) bool {
	defer func() {
		m.current++
	}()

	require.Less(t, m.current, len(m.expected))
	assert.Equal(t, m.expected[m.current].path, path)

	return m.expected[m.current].ret
}

func (m mockUtil) Exists(path string) bool {
	return utilExistsMock(path)
}

type mockCopyFile struct {
	source, target string
	executable     bool
	ret            error
}

type mockCopyFileRuns struct {
	expected []mockCopyFile

	current int
}

func (m *mockCopyFileRuns) verifyExpectation(t *testing.T, executable bool, source string, target string) error {
	defer func() {
		m.current++
	}()

	require.Less(t, m.current, len(m.expected))
	assert.Equal(t, m.expected[m.current].source, source)
	assert.Equal(t, m.expected[m.current].target, target)
	assert.Equal(t, m.expected[m.current].executable, executable)

	return m.expected[m.current].ret

}

func (m mockUtil) CopyFile(source string, target string) error {
	return utilCopyFileMock(source, target)
}

func (m mockUtil) CopyExeFile(src, tgt string) error {
	return utilCopyExeFileMock(src, tgt)
}

type mockWriteFile struct {
	target string
	ret    error
}

type mockWriteFileRuns struct {
	expected []mockWriteFile

	current int
}

func (m *mockWriteFileRuns) verifyExpectation(t *testing.T, target string) error {
	defer func() {
		m.current++
	}()

	require.Less(t, m.current, len(m.expected))
	assert.Equal(t, m.expected[m.current].target, target)

	return m.expected[m.current].ret
}

func (m mockUtil) WriteFile(target string, data []byte) error {
	return utilWriteFileMock(target, data)
}

type mockEnsureSubDir struct {
	parent, name string
	ret          string
}

type mockEnsureSubDirRuns struct {
	expected []mockEnsureSubDir

	current int
}

func (m *mockEnsureSubDirRuns) verifyExpectation(t *testing.T, parent, name string) string {
	defer func() {
		m.current++
	}()

	require.Less(t, m.current, len(m.expected))
	assert.Equal(t, m.expected[m.current].parent, parent)
	assert.Equal(t, m.expected[m.current].name, name)

	return m.expected[m.current].ret
}

func (m mockUtil) EnsureSubDir(parent, name string) string {
	return utilEnsureSubDirMock(parent, name)
}

func (m mockUtil) RequireAndroidSDK() error {
	return utilRequireAndroidSDKMock()
}

func (m mockUtil) AndroidBuildToolsPath() string {
	return utilAndroidBuildToolsPathMock()
}

func (m mockUtil) IsAndroid(os string) bool {
	return utilIsAndroidMock(os)
}

func (m mockUtil) IsIOS(os string) bool {
	return utilIsIOSMock(os)
}

func (m mockUtil) IsMobile(os string) bool {
	return utilIsMobileMock(os)
}
