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

	assert.Len(t, output, 3)
	assert.Contains(t, output[0], "Testing errors")
	assert.Contains(t, output[1], "Cause")
	assert.Contains(t, output[1], "dummy")
	assert.Contains(t, output[2], "At")
}

func TestLogErrorNoErr(t *testing.T) {
	output := bufferLog(t, "Testing errors", nil)

	assert.Len(t, output, 2)
	assert.Contains(t, output[0], "Testing errors")
	assert.Contains(t, output[1], "At")
}
