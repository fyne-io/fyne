// !build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"fyne.io/fyne"
)

func bundleFile(name string, filepath string, f *os.File) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fyne.LogError("Unable to load file "+filepath, err)
		return
	}
	res := fyne.NewStaticResource(path.Base(filepath), bytes)

	_, err = f.WriteString(fmt.Sprintf("var %s = %#v\n", name, res))
	if err != nil {
		fyne.LogError("Unable to write to bundled file", err)
	}
}

func openFile(filename string) *os.File {
	os.Remove(filename)
	_, dirname, _, _ := runtime.Caller(0)
	f, err := os.Create(path.Join(path.Dir(dirname), filename))
	if err != nil {
		fyne.LogError("Unable to open file "+filename, err)
		return nil
	}

	_, err = f.WriteString("// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //\n\npackage main\n\nimport \"fyne.io/fyne\"\n\n")
	if err != nil {
		fyne.LogError("Unable to write file "+filename, err)
		return nil
	}

	return f
}

func main() {
	f := openFile(path.Join("..", "bundled.go"))
	if f == nil {
		return
	}

	bundleFile("themeDarkPreview", "widgets-dark.png", f)
	bundleFile("themeLightPreview", "widgets-light.png", f)

	bundleFile("appearanceIcon", "appearance.svg", f)
	bundleFile("languageIcon", "language.svg", f)
}
