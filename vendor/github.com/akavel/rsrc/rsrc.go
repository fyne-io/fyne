package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/akavel/rsrc/coff"
	"github.com/akavel/rsrc/binutil"
	"github.com/akavel/rsrc/internal"
	"github.com/akavel/rsrc/rsrc"
)

var usage = `USAGE:

%s [-manifest FILE.exe.manifest] [-ico FILE.ico[,FILE2.ico...]] -o FILE.syso
  Generates a .syso file with specified resources embedded in .rsrc section,
  aimed for consumption by Go linker when building Win32 excecutables.

The generated *.syso files should get automatically recognized by 'go build'
command and linked into an executable/library, as long as there are any *.go
files in the same directory.

OPTIONS:
`

func main() {
	//TODO: allow in options advanced specification of multiple resources, as a tree (json?)
	//FIXME: verify that data file size doesn't exceed uint32 max value
	var fnamein, fnameico, fnamedata, fnameout, arch string
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.StringVar(&fnamein, "manifest", "", "path to a Windows manifest file to embed")
	flags.StringVar(&fnameico, "ico", "", "comma-separated list of paths to .ico files to embed")
	flags.StringVar(&fnamedata, "data", "", "path to raw data file to embed [WARNING: useless for Go 1.4+]")
	flags.StringVar(&fnameout, "o", "rsrc.syso", "name of output COFF (.res or .syso) file")
	flags.StringVar(&arch, "arch", "386", "architecture of output file - one of: 386, [EXPERIMENTAL: amd64]")
	_ = flags.Parse(os.Args[1:])
	if fnameout == "" || (fnamein == "" && fnamedata == "" && fnameico == "") {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flags.PrintDefaults()
		os.Exit(1)
	}

	var err error
	switch {
	case fnamein != "" || fnameico != "":
		err = rsrc.Embed(fnameout, arch, fnamein, fnameico)
	case fnamedata != "":
		err = embedData(fnameout, arch, fnamedata)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func embedData(fnameout, arch, fnamedata string) error {
	if !strings.HasSuffix(fnameout, ".syso") {
		return fmt.Errorf("Output file name '%s' must end with '.syso'", fnameout)
	}
	symname := strings.TrimSuffix(fnameout, ".syso")
	ok, err := regexp.MatchString(`^[a-z0-9_]+$`, symname)
	if err != nil {
		return fmt.Errorf("Internal error: %s", err)
	}
	if !ok {
		return fmt.Errorf("Output file name '%s' must be composed of only lowercase letters (a-z), digits (0-9) and underscore (_)", fnameout)
	}

	dat, err := binutil.SizedOpen(fnamedata)
	if err != nil {
		return fmt.Errorf("Error opening data file '%s': %s", fnamedata, err)
	}
	defer dat.Close()

	coff := coff.NewRDATA()
	err = coff.Arch(arch)
	if err != nil {
		return err
	}
	coff.AddData("_brsrc_"+symname, dat)
	coff.AddData("_ersrc_"+symname, io.NewSectionReader(strings.NewReader("\000\000"), 0, 2)) // TODO: why? copied from as-generated
	coff.Freeze()
	err = internal.Write(coff, fnameout)
	if err != nil {
		return err
	}

	//FIXME: output a .c file
	fmt.Println(strings.Replace(`#include "runtime.h"
extern byte _brsrc_NAME[], _ersrc_NAME;

/* func get_NAME() []byte */
void ·get_NAME(Slice a) {
  a.array = _brsrc_NAME;
  a.len = a.cap = &_ersrc_NAME - _brsrc_NAME;
  FLUSH(&a);
}`, "NAME", symname, -1))

	return nil
}