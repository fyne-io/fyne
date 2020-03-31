package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"fyne.io/fyne"
)

// Declare conformity to command interface
var _ command = (*env)(nil)

type version struct {
}

func (v *version) addFlags() {
}

func (v *version) printHelp(indent string) {
	fmt.Println(indent, "The version command prints version information")
	fmt.Println(indent, "Command usage: fyne version")
}

func (v *version) main() {
	ver, err := v.get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("fyne cli version", ver)
}

func (v *version) run(args []string) {
	if len(args) != 0 {
		fyne.LogError("Unexpected parameter after flags", nil)
		return
	}

	v.main()
}

func (v *version) get() (string, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", fmt.Errorf("Could not retrieve version information. Please ensure module support is activated and build again")
	}
	return info.Main.Version, nil
}
