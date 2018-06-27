// +build ignore

package main

import "fmt"
import "io/ioutil"
import "os"
import "path"
import "runtime"
import "strings"

import "github.com/fyne-io/fyne"

const filename = "bundled.go"

func bundleFile(name string, filepath string, f *os.File) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Unable to load file " + filepath)
		return
	}
	res := fyne.NewResource(path.Base(filepath), bytes)

	_, err = f.WriteString("var " + strings.ToLower(name) + " = " + res.ToGo() + "\n")
	if err != nil {
		fmt.Println("Unable to write to file " + filename)
	}
}

func bundleFont(name string, f *os.File) {
	_, dirname, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(dirname), "font", "NotoSans-"+name+".ttf")

	bundleFile(name, path, f)
}

func main() {
	os.Remove(filename)
	_, dirname, _, _ := runtime.Caller(0)
	f, err := os.Create(path.Join(path.Dir(dirname), filename))
	if err != nil {
		fmt.Println("Unable to open file " + filename)
		return
	}

	_, err = f.WriteString("package theme\n\nimport \"github.com/fyne-io/fyne\"\n\n")
	if err != nil {
		fmt.Println("Unable to write file " + filename)
		return
	}

	bundleFont("Regular", f)
	bundleFont("Bold", f)
	bundleFont("Italic", f)
	bundleFont("BoldItalic", f)

	f.Close()
}
