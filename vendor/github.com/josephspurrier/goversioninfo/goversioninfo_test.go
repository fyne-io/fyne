package goversioninfo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/akavel/rsrc/coff"
)

// *****************************************************************************
// Logic Testing
// *****************************************************************************

func TestFile1(t *testing.T) {
	testFile(t, "cmd")
	testFile(t, "explorer")
	testFile(t, "control")
	testFile(t, "simple")
}

func testFile(t *testing.T, filename string) {
	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}
	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	path2, _ := filepath.Abs("./tests/" + filename + ".hex")

	// This is for easily exporting results when the algorithm improves
	/*path3, _ := filepath.Abs("./tests/" + filename + ".out")
	ioutil.WriteFile(path3, vi.Buffer.Bytes(), 0655)*/

	expected, err := ioutil.ReadFile(path2)
	if err != nil {
		t.Error("Could not load "+filename+".hex", err)
	}

	if !bytes.Equal(vi.Buffer.Bytes(), expected) {
		t.Error("Data does not match " + filename + ".hex")
	}
}

func TestWrite32(t *testing.T) {
	doTestWrite(t, "386")
}

func TestWrite64(t *testing.T) {
	doTestWrite(t, "amd64")
}

func doTestWrite(t *testing.T, arch string) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}
	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteSyso(file, arch)

	_, err = ioutil.ReadFile(file)
	if err != nil {
		t.Error("Could not load "+file, err)
	} else {
		os.Remove(file)
	}
}

func TestMalformedJSON(t *testing.T) {
	filename := "bad"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config and return false
	if err := vi.ParseJSON(jsonBytes); err == nil {
		t.Error("Application was supposed to return error, got nil")
	}
}

func TestIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.IconPath = "icon.ico"

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteSyso(file, "386")

	_, err = ioutil.ReadFile(file)
	if err != nil {
		t.Error("Could not load "+file, err)
	} else {
		os.Remove(file)
	}
}

func TestBadIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.IconPath = "icon2.ico"

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteSyso(file, "386")

	_, err = ioutil.ReadFile(file)
	if err != nil {
		os.Remove(file)
	} else {
		t.Error("File should not exist "+file, err)
	}
}

func TestTimestamp(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.Timestamp = true

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteSyso(file, "386")

	_, err = ioutil.ReadFile(file)
	if err != nil {
		t.Error("Could not load "+file, err)
	} else {
		os.Remove(file)
	}
}

func TestVersionString(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}
	if vi.FixedFileInfo.GetVersionString() != "6.3.9600.16384" {
		t.Errorf("Version String does not match: %v", vi.FixedFileInfo.GetVersionString())
	}
}

func TestWriteHex(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}
	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteHex(file)

	_, err = ioutil.ReadFile(file)
	if err != nil {
		t.Error("Could not load "+file, err)
	} else {
		os.Remove(file)
	}
}

// *****************************************************************************
// Examples
// *****************************************************************************

func ExampleUseIcon() {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		fmt.Println("Could not parse "+filename+".json", err)
	}

	vi.IconPath = "icon.ico"

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteSyso(file, "386")

	_, err = ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Could not load "+file, err)
	}
}

func ExampleUseTimestamp() {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		fmt.Println("Could not parse "+filename+".json", err)
	}

	// Write a timestamp even though it is against the spec
	vi.Timestamp = true

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteSyso(file, "386")

	_, err = ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Could not load "+file, err)
	}
}

func TestStr2Uint32(t *testing.T) {
	for _, tt := range []struct {
		in  string
		out uint32
	}{{"0", 0}, {"", 0}, {"FFEF", 65519}, {"\x00\x00", 0}} {
		log.SetOutput(ioutil.Discard)
		got := str2Uint32(tt.in)
		if got != tt.out {
			t.Errorf("%q: awaited %d, got %d.", tt.in, tt.out, got)
		}
	}
}

var unmarshals = []struct {
	in      string
	needErr bool
}{
	{"", false}, {"A", true}, {"1", false}, {`"FfeF"`, false},
	{`"FfeF`, true}, {`"FXXX"`, true},
}

func TestLangID(t *testing.T) {
	var lng LangID
	for _, tt := range unmarshals {
		if err := lng.UnmarshalJSON([]byte(tt.in)); tt.needErr && err == nil {
			t.Errorf("%q: needed error, got nil.", tt.in)
		} else if !tt.needErr && err != nil {
			t.Errorf("%q: got error: %v", tt.in, err)
		}
	}
}

func TestCharsetID(t *testing.T) {
	var cs CharsetID
	for _, tt := range unmarshals {
		if err := cs.UnmarshalJSON([]byte(tt.in)); tt.needErr && err == nil {
			t.Errorf("%q: needed error, got nil.", tt.in)
		} else if !tt.needErr && err != nil {
			t.Errorf("%q: got error: %v", tt.in, err)
		}
	}
}

func TestWriteCoff(t *testing.T) {
	tempFh, err := ioutil.TempFile("", "goversioninfo-test-")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	tempfn := tempFh.Name()
	tempFh.Close()
	defer os.Remove(tempfn)

	if err := writeCoff(nil, ""); err == nil {
		t.Errorf("needed error, got nil")
	}
	if err := writeCoff(nil, tempfn); err != nil {
		t.Errorf("got %v", err)
	}

	if err := writeCoffTo(badWriter{writeErr: io.EOF}, coff.NewRSRC()); err == nil {
		t.Errorf("needed write error, got nil")
	}
	if err := writeCoffTo(badWriter{closeErr: io.EOF}, nil); err == nil {
		t.Errorf("needed close error, got nil")
	}
}

type badWriter struct {
	writeErr, closeErr error
}

func (w badWriter) Write(p []byte) (int, error) {
	return len(p), w.writeErr
}
func (w badWriter) Close() error {
	return w.closeErr
}
