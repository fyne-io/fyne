package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"fyne.io/fyne"
	"fyne.io/fyne/cmd/fyne/commands"
)

// Declare conformity to command interface
var _ commands.Command = (*version)(nil)

type version struct {
}

func (v *version) AddFlags() {
}

func (v *version) PrintHelp(indent string) {
	fmt.Println(indent, "The version command prints version information")
	fmt.Println(indent, "Command usage: fyne version")
}

func (v *version) Run(args []string) {
	if len(args) != 0 {
		fyne.LogError("Unexpected parameter after flags", nil)
		return
	}

	v.main()
}

func (v *version) get() (string, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", fmt.Errorf("could not retrieve version information (ensure module support is activated and build again)")
	}
	return info.Main.Version, nil
}

func (v *version) main() {
	ver, err := v.get()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("fyne cli version", ver)
}
