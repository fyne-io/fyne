// Run a command line helper for various Fyne tools.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"fyne.io/fyne"
)

type bundler struct {
	name, pkg string
	prefix    string
	noheader  bool
}

func writeResource(file, name string) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Unable to load file " + file)
		return
	}

	res := fyne.NewStaticResource(path.Base(file), bytes)

	fmt.Printf("var %s = %#v\n", name, res)
}

func writeHeader(pkg string) {
	fmt.Println("package", pkg)
	fmt.Println()
	fmt.Println("import \"fyne.io/fyne\"")
}

func sanitiseName(file, prefix string) string {
	titled := strings.Title(file)

	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	name := reg.ReplaceAllString(titled, "")

	return prefix + name
}

// Bundle takes a file (at filepath) and serialises it into Go to be output into
// a generated bundle file. The go file will be part of the specified package
// (pkg) and the data will be assigned to variable named "name". If you are
// appending an existing resource file then pass true to noheader as the headers
// should only be output once per file.
func doBundle(name, pkg, prefix string, noheader bool, filepath string) {
	if !noheader {
		writeHeader(pkg)
	}
	fmt.Println()

	if name == "" {
		name = sanitiseName(path.Base(filepath), prefix)
	}
	writeResource(filepath, name)
}

func dirBundle(pkg, prefix string, noheader bool, dirpath string) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading specified directory", err)
		return
	}

	omitHeader := noheader
	for _, file := range files {
		filename := file.Name()
		if path.Ext(filename) == ".go" {
			continue
		}

		doBundle("", pkg, prefix, omitHeader, path.Join(dirpath, filename))
		omitHeader = true // only output header once at most
	}
}

func (b *bundler) addFlags() {
	flag.StringVar(&b.name, "name", "", "The variable name to assign the resource (file mode only)")
	flag.StringVar(&b.pkg, "package", "main", "The package to output in headers (if not appending)")
	flag.StringVar(&b.prefix, "prefix", "resource", "A prefix for variables (ignored if name is set)")
	flag.BoolVar(&b.noheader, "append", false, "Append an existing go file (don't output headers)")
}

func (b *bundler) printHelp(indent string) {
	fmt.Println(indent, "The bundle command embeds static content into your go application.")
	fmt.Println(indent, "Each resource will have a generated filename unless specified")
	fmt.Println(indent, "Command usage: fyne bundle [parameters] file|directory")
}

func (b *bundler) run(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Missing required file or directory parameter after flags")
		return
	}

	stat, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Specified file could not be found")
	} else if stat.IsDir() {
		dirBundle(b.pkg, b.prefix, b.noheader, args[0])
	} else {
		if b.name != "" {
			b.prefix = ""
		}
		doBundle(b.name, b.pkg, b.prefix, b.noheader, args[0])
	}
}
