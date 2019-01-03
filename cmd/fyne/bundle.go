// Run a command line helper for various Fyne tools.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"fyne.io/fyne"
)

type bundler struct {
	name, pkg string
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
	fmt.Println("import \"github.com/fyne-io/fyne\"")
}

func sanitiseName(file string) string {
	name := strings.Replace(file, ".", "-", -1)
	name = strings.Replace(name, " ", "", -1)

	return strings.ToLower(name)
}

// Bundle takes a file (at filepath) and serialises it into Go to be output into
// a generated bundle file. The go file will be part of the specified package
// (pkg) and the data will be assigned to variable named "name". If you are
// appending an existing resource file then pass true to noheader as the headers
// should only be output once per file.
func doBundle(name, pkg string, noheader bool, filepath string) {
	if !noheader {
		writeHeader(pkg)
	}
	fmt.Println()

	if name == "" {
		name = sanitiseName(path.Base(filepath))
	}
	writeResource(filepath, name)
}

func (b *bundler) addFlags() {
	flag.StringVar(&b.name, "name", "", "The variable name to assign the serialised resource")
	flag.StringVar(&b.pkg, "package", "main", "The package to output in headers (if not appending)")
	flag.BoolVar(&b.noheader, "append", false, "Append an existing go file (don't output headers)")
}

func (b *bundler) printHelp(indent string) {
	fmt.Println(indent, "The bundle command embeds static content into your go application.")
	fmt.Println(indent, "Command usage: fyne bundle [parameters] file")
}

func (b *bundler) run(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Missing required file parameter after flags")
		return
	}

	doBundle(b.name, b.pkg, b.noheader, args[0])
}
