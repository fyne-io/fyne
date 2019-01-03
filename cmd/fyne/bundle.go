// Run a command line helper for various Fyne tools.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"fyne.io/fyne"
)

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
func Bundle(name, pkg string, noheader bool, filepath string) {
	if !noheader {
		writeHeader(pkg)
	}
	fmt.Println()

	if name == "" {
		name = sanitiseName(path.Base(filepath))
	}
	writeResource(filepath, name)
}

var name, pkg string
var noheader bool

func bundleFlags() {
	flag.StringVar(&name, "name", "", "The variable name to assign the serialised resource")
	flag.StringVar(&pkg, "package", "main", "The package to output in headers (if not appending)")
	flag.BoolVar(&noheader, "append", false, "Append an existing go file (don't output headers)")
}

func bundleRun(args []string) {
	if len(args) != 1 {
		fmt.Println("Missing required file parameter after flags")
		return
	}

	Bundle(name, pkg, noheader, args[0])
}
