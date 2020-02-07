// +build ignore

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"fyne.io/fyne"
)

const fontFace = "NotoSans"

const fileHeader = "// auto-generated\n" + // to exclude this file in goreportcard
	"// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //"

func formatVariable(name string) string {
	str := strings.Replace(name, "-", "", -1)
	return strings.Replace(str, "_", "", -1)
}

func bundleFile(name string, filepath string, f *os.File) {
	res, err := fyne.LoadResourceFromPath(filepath)
	if err != nil {
		fyne.LogError("Unable to load file "+filepath, err)
		return
	}

	_, err = f.WriteString(fmt.Sprintf("var %s = %#v\n", formatVariable(name), res))
	if err != nil {
		fyne.LogError("Unable to write to bundled file", err)
	}
}

func bundleFont(font, name string, f *os.File) {
	_, dirname, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(dirname), "font", fmt.Sprintf("%s-%s.ttf", font, name))

	if name == "Regular" && font != fontFace {
		name = "Monospace"
	}

	bundleFile(strings.ToLower(name), path, f)
}

func iconDir() string {
	_, dirname, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(dirname), "icons")
}

func bundleIcon(name string, f *os.File) {
	path := path.Join(iconDir(), fmt.Sprintf("%s.svg", name))

	formatted := fmt.Sprintf("%sIconRes", strings.ToLower(name))
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

	_, err = f.WriteString(fileHeader + "\n\npackage theme\n\nimport \"fyne.io/fyne\"\n\n")
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

	bundleIcon("cancel", f)
	bundleIcon("check", f)
	bundleIcon("delete", f)
	bundleIcon("search", f)
	bundleIcon("search-replace", f)
	bundleIcon("menu", f)
	bundleIcon("menu-expand", f)

	bundleIcon("check-box", f)
	bundleIcon("check-box-blank", f)
	bundleIcon("radio-button", f)
	bundleIcon("radio-button-checked", f)

	bundleIcon("content-add", f)
	bundleIcon("content-remove", f)
	bundleIcon("content-cut", f)
	bundleIcon("content-copy", f)
	bundleIcon("content-paste", f)
	bundleIcon("content-redo", f)
	bundleIcon("content-undo", f)

	bundleIcon("document-create", f)
	bundleIcon("document-print", f)
	bundleIcon("document-save", f)

	bundleIcon("info", f)
	bundleIcon("question", f)
	bundleIcon("warning", f)

	bundleIcon("arrow-back", f)
	bundleIcon("arrow-down", f)
	bundleIcon("arrow-forward", f)
	bundleIcon("arrow-up", f)
	bundleIcon("arrow-drop-down", f)
	bundleIcon("arrow-drop-up", f)

	bundleIcon("folder", f)
	bundleIcon("folder-new", f)
	bundleIcon("folder-open", f)
	bundleIcon("help", f)
	bundleIcon("home", f)
	bundleIcon("settings", f)

	bundleIcon("mail-attachment", f)
	bundleIcon("mail-compose", f)
	bundleIcon("mail-forward", f)
	bundleIcon("mail-reply", f)
	bundleIcon("mail-reply_all", f)
	bundleIcon("mail-send", f)

	bundleIcon("media-fast-forward", f)
	bundleIcon("media-fast-rewind", f)
	bundleIcon("media-pause", f)
	bundleIcon("media-play", f)
	bundleIcon("media-record", f)
	bundleIcon("media-replay", f)
	bundleIcon("media-skip-next", f)
	bundleIcon("media-skip-previous", f)

	bundleIcon("view-fullscreen", f)
	bundleIcon("view-refresh", f)
	bundleIcon("view-zoom-fit", f)
	bundleIcon("view-zoom-in", f)
	bundleIcon("view-zoom-out", f)

	bundleIcon("volume-down", f)
	bundleIcon("volume-mute", f)
	bundleIcon("volume-up", f)

	bundleIcon("visibility", f)
	bundleIcon("visibility-off", f)

	f.Close()
}
