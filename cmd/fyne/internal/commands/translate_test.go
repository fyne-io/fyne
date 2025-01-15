package commands

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

const exampleSource = `
package main

func main() {
	lang.X("example", "Example")
}
`

func createTestTranslateFiles(t *testing.T, dir, file string) {
	src := filepath.Join(dir, file)

	f, err := os.Create(src)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte(exampleSource)); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTranslateCommand(t *testing.T) {
	cmd := Translate()
	if cmd.Name != "translate" {
		t.Fatalf("invalid command name: %v", cmd.Name)
	}
}

func TestFindSourcesDefault(t *testing.T) {
	dir := t.TempDir()
	createTestTranslateFiles(t, dir, "foo.go")

	sources, err := findSources(nil, ".go", dir)
	if err != nil {
		t.Fatalf("failed to find sources: %v", err)
	}
	if len(sources) != 1 {
		t.Errorf("wrong number of sources: %v", len(sources))
	}
}

func TestFindSourcesNone(t *testing.T) {
	dir := t.TempDir()
	createTestTranslateFiles(t, dir, "foo.go")

	sources, err := findSources(nil, ".go", filepath.Join(dir, "nonexistent"))
	if err != nil {
		t.Fatalf("failed to find sources: %v", err)
	}
	if len(sources) != 0 {
		t.Errorf("wrong number of sources: %v", len(sources))
	}
}

func TestUpdateTranslationsFile(t *testing.T) {
	src := "foo.go"
	dir := t.TempDir()
	createTestTranslateFiles(t, dir, src)
	srcpath := filepath.Join(dir, src)

	opts := translateOpts{}
	dst := filepath.Join(dir, "en.json")

	if err := updateTranslationsFile(dst, []string{srcpath}, &opts); err != nil {
		t.Fatal(err)
	}

	if err := updateTranslationsFile(dst, []string{srcpath}, &opts); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(dst)
	if err != nil {
		t.Fatal(err)
	}

	translations := make(map[string]any)
	dec := json.NewDecoder(f)
	if err := dec.Decode(&translations); err != nil {
		t.Fatal(err)
	}
	f.Close()

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
	dir := t.TempDir()
	createTestTranslateFiles(t, dir, src)
	files, err := findFilesExt(dir, ".go")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("invalid number of files: %v", len(files))
	}
	if len(files) == 1 && filepath.Base(files[0]) != src {
		t.Errorf("invalid file: %v", filepath.Base(files[0]))
	}
}

func TestWriteTranslationsFile(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "foo.json")

	if err := writeTranslationsFile([]byte(`{"a":1}`), dst); err != nil {
		t.Fatalf("failed to write translations file: %v", err)
	}

	if err := writeTranslationsFile([]byte(`{"a":2}`), dst); err != nil {
		t.Fatalf("failed to write translations file: %v", err)
	}
}

func TestUpdateTranslationsHash(t *testing.T) {
	src := "foo.go"
	dir := t.TempDir()
	createTestTranslateFiles(t, dir, src)
	srcpath := filepath.Join(dir, src)

	opts := translateOpts{}
	translations := make(map[string]any)
	if err := updateTranslationsHash(translations, []string{srcpath}, &opts); err != nil {
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
	dir := t.TempDir()
	createTestTranslateFiles(t, dir, src)

	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, filepath.Join(dir, src), nil, parser.AllErrors)
	if err != nil {
		t.Fatalf("failed to parse source: %v", err)
	}

	opts := translateOpts{}
	translations := make(map[string]any)
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

func TestTranslateVisitorNew(t *testing.T) {
	if translateNew(&visitor{}, ast.NewIdent("whatever")) != nil {
		t.Fatalf("failed to get the correct state: nil")
	}

	if translateNew(&visitor{}, ast.NewIdent("lang")) == nil {
		t.Fatalf("failed to get the correct state")
	}
}

func TestTranslateVisitorCall(t *testing.T) {
	if translateCall(&visitor{}, &ast.IfStmt{}) != nil {
		t.Fatalf("failed to get the correct state: nil")
	}

	if translateCall(&visitor{}, ast.NewIdent("X")) == nil {
		t.Fatalf("failed to get the correct state")
	}
}

func TestTranslateVisitorLocalize(t *testing.T) {
	translations := make(map[string]any)
	v := &visitor{
		opts: &translateOpts{},
		m:    translations,
	}
	if translateLocalize(v, &ast.IfStmt{}) != nil {
		t.Fatalf("failed to get the correct state: nil")
	}

	key := "yay"
	if translateLocalize(v, &ast.BasicLit{Value: fmt.Sprintf("%q", key)}) != nil {
		t.Fatalf("failed to get the correct state")
	}

	if v.key != key {
		t.Errorf("invalid key: %q", v.key)
	}
}

func TestTranslateVisitorKey(t *testing.T) {
	translations := make(map[string]any)
	v := &visitor{
		opts: &translateOpts{},
		m:    translations,
	}
	if translateLocalize(v, &ast.IfStmt{}) != nil {
		t.Fatalf("failed to get the correct state: nil")
	}

	key := "whee"
	if translateKey(v, &ast.BasicLit{Value: fmt.Sprintf("%q", key)}) == nil {
		t.Fatalf("failed to get the correct state")
	}

	if v.key != key {
		t.Errorf("invalid key: %q", v.key)
	}
}

func TestTranslateVisitorFallback(t *testing.T) {
	translations := make(map[string]any)
	v := &visitor{
		opts: &translateOpts{},
		m:    translations,
	}
	if translateLocalize(v, &ast.IfStmt{}) != nil {
		t.Fatalf("failed to get the correct state: nil")
	}

	fallback := "WHEE"
	if translateLocalize(v, &ast.BasicLit{Value: fmt.Sprintf("%q", fallback)}) != nil {
		t.Fatalf("failed to get the correct state")
	}

	if v.fallback != fallback {
		t.Errorf("invalid fallback: %q", v.fallback)
	}
}
