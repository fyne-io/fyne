// Run a command line helper for various Fyne tools.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// first let's extract the first "command"
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("The fyne command requires a command parameter, one of \"bundle\" or \"package\"")
		return
	}

	command := args[0]
	if command == "bundle" {
		bundleFlags()

		// then parse the remaining args
		flag.CommandLine.Parse(args[1:])
		bundleRun(flag.Args())
		return
	} else if command == "package" {
		packageFlags()

		// then parse the remaining args
		flag.CommandLine.Parse(args[1:])
		packageRun(flag.Args())
		return
	}

	fmt.Println("Unsupported command", command)
	return
}
