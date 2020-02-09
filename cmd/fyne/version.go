package main

import (
	"fmt"
	"go/build"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const PKG = "fyne.io/fyne"

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
	p, e := c.Import(PKG, ".", build.FindOnly)
	if e != nil {
		return e
	}
	// p.PkgObj =~ /home/ser/go/pkg/mod/fyne.io/fyne@v1.2.2/pkg/li
	ps := strings.Split(p.PkgObj, "@")
	if len(ps) < 2 {
		return errors.New(fmt.Sprintf("Malformed package info %s; not using modules?", p.PkgObj))
	}
	vs := strings.Split(ps[1], "/")
	if len(vs) < 2 {
		return errors.New(fmt.Sprintf("Missing version information %s; not using modules?", p.PkgObj))
	}
	fmt.Println(vs[0])
	return nil
}

func (v *versioner) printHelp(indent string) {
	fmt.Println(indent, "Print the version of the Fyne library that would be used to build the package.")
}

func (v *versioner) addFlags() {
}
