package main

import (
	"fmt"
	"go/build"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const pkg = "fyne.io/fyne"

// Declare conformity to command interface
var _ command = (*versioner)(nil)

type versioner struct {
	context build.Context
}

func (v *versioner) run(_ []string) {
	err := v.validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = v.doVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func (v *versioner) validate() error {
	return nil
}

func (v *versioner) doVersion() error {
	c := build.Default
	p, e := c.Import(pkg, ".", build.FindOnly)
	if e != nil {
		return e
	}
	ve, e := parsePkgString(p.PkgObj)
	if e != nil {
		return e
	}
	fmt.Println(ve)
	return nil
}

func parsePkgString(s string) (string, error) {
	// @ separates the package path from the version
	ps := strings.Split(s, "@")
	if len(ps) < 2 {
		return "", errors.New(fmt.Sprintf("No version information in %s; not using modules?", s))
	}
	// PkgObj contains some junk after the version -- "junk," in this case,
	// meaning some internal information I don't understand. It's not part of
	// the version, in any case.
	vs := strings.Split(ps[1], "/")
	if len(vs) < 1 {
		return "", errors.New(fmt.Sprintf("Missing version information %s; not using modules?", s))
	}
	return vs[0], nil
}

func (v *versioner) printHelp(indent string) {
	fmt.Println(indent, "Print the version of the Fyne library that would be used to build the package.")
}

func (v *versioner) addFlags() {
}
