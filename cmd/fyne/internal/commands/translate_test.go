package commands

import (
	"encoding/json"
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

func TestTranslate(t *testing.T) {
	dir := t.TempDir()
	src := path.Join(dir, "foo.go")

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(src)
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte(exampleSource))
	f.Close()

	opts := translateOpts{}
	dst := "en.json"

	if err := updateTranslationsFile(&opts, dst, []string{src}); err != nil {
		t.Fatal(err)
	}

	f, err = os.Open(dst)
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
