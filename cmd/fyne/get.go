package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type getter struct {
	pkg string
}

// Declare conformity to command interface
var _ command = (*getter)(nil)

func goPath() string {
	cmd := exec.Command("go", "env", "GOPATH")
	out, err := cmd.CombinedOutput()

	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "go")
	}

	return strings.TrimSpace(string(out[0 : len(out)-1]))
}

func (g *getter) get() error {
	path := filepath.Join(goPath(), "src", g.pkg)

	cmd := exec.Command("go", "get", "-u", g.pkg)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(out))
		return err
	}

	if !exists(path) {
		return errors.New("Package download not found in expected location")
	}

	install := &installer{srcDir: path}
	err = install.validate()
	if err != nil {
		return errors.Wrap(err, "Failed to set up installer")
	}
	return install.install()
}

func (g *getter) addFlags() {
}

func (g *getter) printHelp(indent string) {
	fmt.Println(indent, "The get command downloads and installs a Fyne application.")
	fmt.Println(indent, "A single parameter is required to specify the Go package, as with \"go get\"")
}

func (g *getter) run(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Missing \"package\" argument to define the application to get")
		os.Exit(1)
	}
	g.pkg = args[0]

	err := g.get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
