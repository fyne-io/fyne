// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"fyne.io/fyne"
)

const fontFace = "NotoSans"

func formatVariable(name string) string {
	str := strings.Replace(name, "-", "", -1)
	return strings.Replace(str, "_", "", -1)
}

func bundleFile(name string, filepath string, f *os.File) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fyne.LogError("Unable to load file "+filepath, err)
		return
	}
	res := fyne.NewStaticResource(path.Base(filepath), bytes)

	_, err = f.WriteString(fmt.Sprintf("var %s = %#v\n", formatVariable(name), res))
	if err != nil {
		fyne.LogError("Unable to write to bundled file", err)
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

func iconDir() string {
	_, dirname, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(dirname), "icons")
}

func bundleIcon(name, theme string, f *os.File) {
	path := path.Join(iconDir(), name+"-"+theme+".svg")

	formatted := strings.ToLower(name) + strings.Title(strings.ToLower(theme))
	bundleFile(formatted, path, f)
}

func openFile(filename string) *os.File {
	os.Remove(filename)
	_, dirname, _, _ := runtime.Caller(0)
	f, err := os.Create(path.Join(path.Dir(dirname), filename))
	if err != nil {
		fyne.LogError("Unable to open file "+filename, err)
		return nil
	}

	_, err = f.WriteString("package theme\n\nimport \"fyne.io/fyne\"\n\n")
	if err != nil {
		fyne.LogError("Unable to write file "+filename, err)
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
	bundleFont("NotoMono", "Regular", f)
	f.Close()

	f = openFile("bundled-icons.go")
	if f == nil {
		return
	}

	icon := path.Join(iconDir(), "fyne.png")
	bundleFile("fyne-logo", icon, f)

	themes := []string{"dark", "light"}
	for _, theme := range themes {
		bundleIcon("cancel", theme, f)
		bundleIcon("check", theme, f)
		bundleIcon("delete", theme, f)
		bundleIcon("search", theme, f)
		bundleIcon("search-replace", theme, f)

		bundleIcon("check-box", theme, f)
		bundleIcon("check-box-blank", theme, f)
		bundleIcon("radio-button", theme, f)
		bundleIcon("radio-button-checked", theme, f)

		bundleIcon("content-add", theme, f)
		bundleIcon("content-remove", theme, f)
		bundleIcon("content-cut", theme, f)
		bundleIcon("content-copy", theme, f)
		bundleIcon("content-paste", theme, f)
		bundleIcon("content-redo", theme, f)
		bundleIcon("content-undo", theme, f)

		bundleIcon("document-create", theme, f)
		bundleIcon("document-print", theme, f)
		bundleIcon("document-save", theme, f)

		bundleIcon("info", theme, f)
		bundleIcon("question", theme, f)
		bundleIcon("warning", theme, f)

		bundleIcon("arrow-back", theme, f)
		bundleIcon("arrow-down", theme, f)
		bundleIcon("arrow-forward", theme, f)
		bundleIcon("arrow-up", theme, f)

		bundleIcon("folder", theme, f)
		bundleIcon("folder-new", theme, f)
		bundleIcon("folder-open", theme, f)
		bundleIcon("help", theme, f)
		bundleIcon("home", theme, f)

		bundleIcon("mail-attachment", theme, f)
		bundleIcon("mail-compose", theme, f)
		bundleIcon("mail-forward", theme, f)
		bundleIcon("mail-reply", theme, f)
		bundleIcon("mail-reply_all", theme, f)
		bundleIcon("mail-send", theme, f)

		bundleIcon("view-fullscreen", theme, f)
		bundleIcon("view-refresh", theme, f)
		bundleIcon("view-zoom-fit", theme, f)
		bundleIcon("view-zoom-in", theme, f)
		bundleIcon("view-zoom-out", theme, f)
	}
	f.Close()
}
