package app

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingsLoad(t *testing.T) {
	testPath := filepath.Join(os.TempDir(), "fyne-settings-test.json")

	file, _ := os.Create(testPath)
	defer os.Remove(testPath)
	buf := bufio.NewWriter(file)
	buf.WriteString("{\"theme\":\"light\"}")
	buf.Flush()
	file.Close()

	settings := &settings{}
	err := settings.loadFromFile(testPath)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "light", settings.schema.ThemeName)

	os.Remove(testPath)
	file, _ = os.Create(testPath)
	buf = bufio.NewWriter(file)
	buf.WriteString("{\"theme\":\"dark\"}")
	buf.Flush()
	file.Close()

	err = settings.loadFromFile(testPath)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "dark", settings.schema.ThemeName)
}
