// !build ignore

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"fyne.io/fyne"
)

func bundleFile(name string, filepath string, f *os.File) {
	res, err := fyne.LoadResourceFromPath(filepath)
	if err != nil {
		fyne.LogError("Unable to load file "+filepath, err)
		return
	}

	_, err = f.WriteString(fmt.Sprintf("var %s = %#v\n", name, res))
	if err != nil {
		fyne.LogError("Unable to write to bundled file", err)
	}
}

func openFile(filename string) *os.File {
	err := os.Remove(filename)
	if err != nil {
		fyne.LogError("Unable to remove file "+filename, err)
		return nil
	}

	_, dirname, _, _ := runtime.Caller(0)
	f, err := os.Create(path.Join(path.Dir(dirname), filename))
	if err != nil {
		fyne.LogError("Unable to open file "+filename, err)
		return nil
	}

	_, err = f.WriteString("// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //\n\npackage settings\n\nimport \"fyne.io/fyne\"\n\n")
	if err != nil {
		fyne.LogError("Unable to write file "+filename, err)
		return nil
	}

	return f
}

func main() {
	f := openFile(path.Join("..", "settings", "bundled.go"))
	if f == nil {
		return
	}

	bundleFile("appearanceIcon", "appearance.svg", f)

	// TODO: Add in language support and icon
	//bundleFile("languageIcon", "language.svg", f)
}
