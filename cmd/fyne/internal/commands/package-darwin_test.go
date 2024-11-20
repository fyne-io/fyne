package commands

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindTranslationLanguages(t *testing.T) {
	dir := t.TempDir()

	f, err := os.Create(filepath.Join(dir, "en.json"))
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte(""))
	f.Close()

	langs, err := findTranslationLanguages(dir)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(langs))
	assert.Equal(t, "en", langs[0])
}
