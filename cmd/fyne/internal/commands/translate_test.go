package commands

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"testing"
)

const exampleSource = `
package main

func main() {
	lang.X("example", "Example")
}
`

func createTestTranslateFiles(t *testing.T, file string) string {
	dir := t.TempDir()
	src := path.Join(dir, file)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(src)
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte(exampleSource))
	f.Close()

	return dir
}

func TestUpdateTranslationsFile(t *testing.T) {
	src := "foo.go"
	createTestTranslateFiles(t, src)

	opts := translateOpts{}
	dst := "en.json"

	if err := updateTranslationsFile(&opts, dst, []string{src}); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(dst)
	if err != nil {
		t.Fatal(err)
	}

	translations := make(map[string]interface{})
	dec := json.NewDecoder(f)
	if err := dec.Decode(&translations); err != nil {
		t.Fatal(err)
	}

	key := "example"
	val, found := translations[key]
	if !found {
		t.Errorf("failed to find key: %v", key)
	}
	if val != "Example" {
		t.Errorf("invalid value for key: %v: %v", key, val)
	}
}

func TestTranslateFindFilesExt(t *testing.T) {
	src := "foo.go"
	createTestTranslateFiles(t, src)
	files, err := findFilesExt(".", ".go")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("invalid number of files: %v", len(files))
	}
}

func TestWriteTranslationsFile(t *testing.T) {
	dir := t.TempDir()
	dst := path.Join(dir, "foo.json")
	perm := fs.FileMode(0644)

	if err := writeTranslationsFile([]byte(`{"a":1}`), dst, nil); err != nil {
		t.Fatalf("failed to write translations file: %v", err)
	}

	f, err := os.Open(dst)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Mode().Perm() != perm {
		t.Errorf("invalid file permissions: %v", perm)
	}

	if err := writeTranslationsFile([]byte(`{"a":2}`), dst, f); err != nil {
		t.Fatalf("failed to write translations file: %v", err)
	}

	if fi.Mode().Perm() != perm {
		t.Errorf("invalid file permissions: %v", perm)
	}
}

func TestUpdateTranslationsHash(t *testing.T) {
	src := "foo.go"
	createTestTranslateFiles(t, src)

	opts := translateOpts{}
	translations := make(map[string]interface{})
	if err := updateTranslationsHash(&opts, translations, []string{src}); err != nil {
		t.Fatal(err)
	}

	key := "example"
	val, found := translations[key]
	if !found {
		t.Errorf("failed to find key: %v", key)
	}
	if val != "Example" {
		t.Errorf("invalid value for key: %v: %v", key, val)
	}
}

func TestTranslationsVisitor(t *testing.T) {
	src := "foo.go"
	dir := createTestTranslateFiles(t, src)

	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, path.Join(dir, src), nil, parser.AllErrors)
	if err != nil {
		t.Fatalf("failed to parse source: %v", err)
	}

	opts := translateOpts{}
	translations := make(map[string]interface{})
	ast.Walk(&visitor{opts: &opts, m: translations}, af)

	key := "example"
	val, found := translations[key]
	if !found {
		t.Errorf("failed to find key: %v", key)
	}
	if val != "Example" {
		t.Errorf("invalid value for key: %v: %v", key, val)
	}
}
