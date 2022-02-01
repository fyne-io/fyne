package commands

import (
	"fmt"
	"testing"

	"github.com/mcuadros/go-version"
	"github.com/stretchr/testify/assert"
)

type runOutputValue struct {
	// Expected value
	args []string
	// Returned value
	ret []byte
	err error
}

type testCommandCall struct {
	calls []runOutputValue
	index int
	t     *testing.T
}

func (t *testCommandCall) RunOutput(args ...string) ([]byte, error) {
	// Check that we have no more than the expected number of call
	assert.Less(t.t, t.index, len(t.calls))
	// Check that we have the expected number of parameters for this call
	assert.Equal(t.t, len(t.calls[t.index].args), len(args))
	// Check that each argument match our expectation
	for index, value := range args {
		assert.Equal(t.t, t.calls[t.index].args[index], value)
	}

	ret, err := t.calls[t.index].ret, t.calls[t.index].err
	t.index++

	return ret, err
}

func Test_CheckVersionNoGo(t *testing.T) {
	expected := []runOutputValue{{
		args: []string{"version"},
		ret:  nil,
		err:  fmt.Errorf("file not found"),
	}}

	commandNil := &testCommandCall{calls: expected, t: t}
	assert.Nil(t, checkVersion(commandNil, nil))

	commandNotNil := &testCommandCall{calls: expected, t: t}
	assert.NotNil(t, checkVersion(commandNotNil, version.NewConstrainGroupFromString(">=1.17")))
}

func Test_CheckVersionValidValue(t *testing.T) {
	expected := []runOutputValue{
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
	assert.Nil(t, checkVersion(versionOk, version.NewConstrainGroupFromString(">=1.17")))
	assert.Nil(t, checkVersion(versionOk, version.NewConstrainGroupFromString(">=1.17")))

	versionNotOk := &testCommandCall{calls: expected, t: t}
	assert.NotNil(t, checkVersion(versionNotOk, version.NewConstrainGroupFromString("<1.17")))
	assert.NotNil(t, checkVersion(versionNotOk, version.NewConstrainGroupFromString("<1.17")))
}

func Test_CheckVersionInvalidValue(t *testing.T) {
	expected := []runOutputValue{
		{
			args: []string{"version"},
			ret:  []byte("got noversion go1.17.6 windows/amd64"),
			err:  nil,
		},
		{
			args: []string{"version"},
			ret:  []byte("go version really"),
			err:  nil,
		},
		{
			args: []string{"version"},
			ret:  []byte("go version 1.17.6 windows/amd64"),
			err:  nil,
		},
		{
			args: []string{"version"},
			ret:  []byte("go version gowrong windows/amd64"),
			err:  nil,
		},
	}

	versionInvalid := &testCommandCall{calls: expected, t: t}
	for range expected {
		assert.NotNil(t, checkVersion(versionInvalid, version.NewConstrainGroupFromString(">=1.17")))
	}
}
