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

func TestLogError(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	write := bufio.NewWriter(buf)

	log.SetOutput(write)
	err := errors.New("dummy error")
	LogError("Testing errors", err)
	log.SetOutput(os.Stdout)

	err = write.Flush()
	if err != nil {
		t.Error(err)
	}
	output := strings.TrimSpace(buf.String())

	assert.Equal(t, 3, len(strings.Split(output, "\n")))
}

func TestLogErrorNoErr(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	write := bufio.NewWriter(buf)

	log.SetOutput(write)
	LogError("Testing errors", nil)
	log.SetOutput(os.Stdout)

	err := write.Flush()
	if err != nil {
		t.Error(err)
	}
	output := strings.TrimSpace(buf.String())

	assert.Equal(t, 2, len(strings.Split(output, "\n")))
}
