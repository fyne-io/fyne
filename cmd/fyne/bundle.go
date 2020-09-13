package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"fyne.io/fyne"
	"github.com/spf13/cobra"
)

type bundler struct {
	name, pkg string
	prefix    string
	noheader  bool
}

func writeResource(file, name string) {
	res, err := fyne.LoadResourceFromPath(file)
	if err != nil {
		fyne.LogError("Unable to load file "+file, err)
		return
	}

	fmt.Printf("var %s = %#v\n", name, res)
}

func writeHeader(pkg string) {
	fmt.Println("// auto-generated")
	fmt.Println()
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
		fyne.LogError("Error reading specified directory", err)
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

func (b *bundler) bundle(ccmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Missing required file or directory parameter after flags")
	}

	switch stat, err := os.Stat(args[0]); {
	case os.IsNotExist(err):
		fyne.LogError("Specified file could not be found", err)
	case stat.IsDir():
		dirBundle(b.pkg, b.prefix, b.noheader, args[0])
	case b.name != "":
		b.prefix = ""
		fallthrough
	default:
		doBundle(b.name, b.pkg, b.prefix, b.noheader, args[0])
	}

	return nil
}

func bundleCmd() *cobra.Command {
	b := bundler{}
	cmd := &cobra.Command{
		Use:   "bundle",
		Short: "Embeds static content into your go application",
		Long:  "The bundle command embeds static content into your go application. Each resource will have a generated filename unless specified.",
		RunE:  b.bundle,
	}

	cmd.PersistentFlags().StringVar(&b.name, "name", "", "the variable name to assign the resource (file mode only)")
	cmd.PersistentFlags().StringVar(&b.pkg, "package", "main", "the package to output in headers (if not appending)")
	cmd.PersistentFlags().StringVar(&b.prefix, "prefix", "resource", "a prefix for variables (ignored if name is set)")
	cmd.PersistentFlags().BoolVar(&b.noheader, "append", false, "Append an existing go file (don't output headers)")

	return cmd
}
