package fyne

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func bufferLog(t *testing.T, reason string, err error) []string {
	buf := bytes.NewBuffer([]byte{})
	write := bufio.NewWriter(buf)

	log.SetOutput(write)
	LogError(reason, err)
	log.SetOutput(os.Stdout)

	err = write.Flush()
	if err != nil {
		t.Error(err)
	}
	output := strings.TrimSpace(buf.String())
	return strings.Split(output, "\n")
}

func TestLogError(t *testing.T) {
	err := errors.New("dummy error")
	output := bufferLog(t, "Testing errors", err)

	assert.Equal(t, 3, len(output))
	assert.True(t, strings.Contains(output[0], "Testing errors"))
	assert.True(t, strings.Contains(output[1], "Cause"))
	assert.True(t, strings.Contains(output[1], "dummy"))
	assert.True(t, strings.Contains(output[2], "At"))
}

func TestLogErrorNoErr(t *testing.T) {
	output := bufferLog(t, "Testing errors", nil)

	assert.Equal(t, 2, len(output))
	assert.True(t, strings.Contains(output[0], "Testing errors"))
	assert.True(t, strings.Contains(output[1], "At"))
}
