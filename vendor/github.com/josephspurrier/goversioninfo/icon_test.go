package goversioninfo

import (
	"log"
	"os"
	"runtime"
	"testing"
)

func TestTempIcon(t *testing.T) {
	icoPath := "tmp.ico"
	outPath := "resource.syso"
	vi := &VersionInfo{}
	vi.IconPath = icoPath

	f, _ := os.Create(icoPath)
	err := f.Close()
	if err != nil {
		log.Println("Unexpected error closing new file")
	}

	vi.Build()
	vi.Walk()
	err = vi.WriteSyso(outPath, runtime.GOARCH)
	if err != nil {
		log.Println("Unexpected error writing resource")
	}

	err = os.Remove(icoPath)
	if err != nil {
		t.Fatal("Error deleting temporary icon", err)
	}
}