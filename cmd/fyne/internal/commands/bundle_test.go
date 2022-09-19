package commands

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitiseName(t *testing.T) {
	file := "example.png"
	prefix := "file"

	assert.Equal(t, "fileExamplePng", sanitiseName(file, prefix))
}

func TestSanitiseName_Exported(t *testing.T) {
	file := "example.png"
	prefix := "Public"

	assert.Equal(t, "PublicExamplePng", sanitiseName(file, prefix))
}

func TestSanitiseName_Special(t *testing.T) {
	file := "a longer name (with-syms).png"
	prefix := "file"

	assert.Equal(t, "fileALongerNameWithSymsPng", sanitiseName(file, prefix))
}

func TestWriteResource(t *testing.T) {
	f, err := ioutil.TempFile("", "*.go")
	if err != nil {
		t.Fatal("Unable to create temp file:", err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	writeHeader("test", f)
	writeResource("testdata/bundle/content.txt", "contentTxt", f)

	const header = fileHeader + "\n\npackage test\n\nimport \"fyne.io/fyne/v2\"\n\n"
	const expected = header + "var contentTxt = &fyne.StaticResource{\n\tStaticName: \"content.txt\",\n\tStaticContent: []byte(\n\t\t\"I am bundled :)\"),\n}\n"

	// Seek file to start so we can read the written data.
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal("Unable to seek temp file:", err)
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal("Unable to read temp file:", err)
	}

	assert.Equal(t, expected, string(content))
}

func BenchmarkWriteResource(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f, _ := ioutil.TempFile("", "*.go")
		b.StartTimer()

		writeHeader("test", f)
		writeResource("testdata/bundle/content.txt", "contentTxt", f)

		b.StopTimer()
		f.Close()
		os.Remove(f.Name())
		b.StartTimer()
	}
}
