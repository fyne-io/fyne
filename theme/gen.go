// +build ignore

package main

import "fmt"
import "io/ioutil"
import "os"
import "path"
import "runtime"
import "strings"

import "github.com/fyne-io/fyne"

const fontFace = "NotoSans"

func bundleFile(name string, filepath string, f *os.File) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Unable to load file " + filepath)
		return
	}
	res := fyne.NewResource(path.Base(filepath), bytes)

	_, err = f.WriteString("var " + strings.Replace(name, "-", "", -1) + " = " + res.ToGo() + "\n")
	if err != nil {
		fmt.Println("Unable to write to bundled file")
	}
}

func bundleFont(font, name string, f *os.File) {
	_, dirname, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(dirname), "font", font+"-"+name+".ttf")

	if name == "Regular" && font != fontFace {
		name = "Monospace"
	}

	bundleFile(strings.ToLower(name), path, f)
}

func bundleIcon(name, theme string, f *os.File) {
	_, dirname, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(dirname), "icons", name+"-"+theme+".svg")

	formatted := strings.ToLower(name) + strings.Title(strings.ToLower(theme))
	bundleFile(formatted, path, f)
}

func openFile(filename string) *os.File {
	os.Remove(filename)
	_, dirname, _, _ := runtime.Caller(0)
	f, err := os.Create(path.Join(path.Dir(dirname), filename))
	if err != nil {
		fmt.Println("Unable to open file " + filename)
		return nil
	}

	_, err = f.WriteString("package theme\n\nimport \"github.com/fyne-io/fyne\"\n\n")
	if err != nil {
		fmt.Println("Unable to write file " + filename)
		return nil
	}

	return f
}

func main() {
	f := openFile("bundled-fonts.go")
	if f == nil {
		return
	}
	bundleFont(fontFace, "Regular", f)
	bundleFont(fontFace, "Bold", f)
	bundleFont(fontFace, "Italic", f)
	bundleFont(fontFace, "BoldItalic", f)
	bundleFont("SourceCodePro", "Regular", f)
	f.Close()

	f = openFile("bundled-icons.go")
	if f == nil {
		return
	}

	themes := []string{"dark", "light"}
	for _, theme := range themes {
		bundleIcon("cancel", theme, f)
		bundleIcon("check", theme, f)
		bundleIcon("check-box", theme, f)
		bundleIcon("check-box-blank", theme, f)
	}
	f.Close()
}
