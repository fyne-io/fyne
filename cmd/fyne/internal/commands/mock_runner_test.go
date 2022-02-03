package commands

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRunOutputValue struct {
	// Expected value
	dir   interface{}
	env   []string
	osEnv bool
	args  []string
	// Returned value
	ret []byte
	err error

	// To check if all expected call did happen
	dirSet bool
	envSet bool
}

type testCommandCall struct {
	calls []mockRunOutputValue
	index int
	t     *testing.T
}

func (t *testCommandCall) RunOutput(args ...string) ([]byte, error) {
	// Check that we have no more than the expected number of call
	assert.Less(t.t, t.index, len(t.calls))
	// Check that we have the expected number of parameters for this call
	assert.Equal(t.t, len(t.calls[t.index].args), len(args))
	// Check that each argument match our expectation
	if len(t.calls[t.index].args) == len(args) {
		for index, value := range args {
			assert.Equal(t.t, t.calls[t.index].args[index], value)
		}
	} else {
		fmt.Printf("%v\n", args)
	}

	ret, err := t.calls[t.index].ret, t.calls[t.index].err
	t.index++

	return ret, err
}

func (t *testCommandCall) DirSet(dir string) {
	// Check that we have no more than the expected number of call
	assert.Less(t.t, t.index, len(t.calls))

	assert.Equal(t.t, t.calls[t.index].dir.(string), dir)
	t.calls[t.index].dirSet = true
}

func (t *testCommandCall) EnvSet(env []string) {
	// Check that we have no more than the expected number of call
	assert.Less(t.t, t.index, len(t.calls))

	// Prepare array for comparison
	expectedEnv := t.calls[t.index].env
	if t.calls[t.index].osEnv {
		expectedEnv = append(expectedEnv, os.Environ()...)
	}
	sort.Strings(expectedEnv)
	sort.Strings(env)

	// First check length of expected and passed environment are equal
	assert.Equal(t.t, len(expectedEnv), len(env))
	if len(expectedEnv) == len(env) {
		// Check each independent environement variable match our expectation
		for index, value := range env {
			assert.Equal(t.t, expectedEnv[index], value)
		}
	}
	t.calls[t.index].envSet = true
}

func (t *testCommandCall) VerifyExpectation() {
	// Expected as many call as we got
	assert.Equal(t.t, len(t.calls), t.index)
	// Check if every call really matched our expectaction
	for _, value := range t.calls {
		if value.dir != nil {
			assert.Equal(t.t, true, value.dirSet)
		}
		if len(value.env) > 0 {
			assert.Equal(t.t, true, value.envSet)
		}
	}
}
