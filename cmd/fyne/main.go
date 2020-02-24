// Run a command line helper for various Fyne tools.
package main

import (
	"flag"
	"fmt"
	"os"
)

var commands map[string]command
var provider command

func printUsage() {
	fmt.Println("Usage: fyne [command] [parameters], where command is one of:")
	fmt.Print("  ")

	i := 0
	for id := range commands {
		fmt.Print(id)

		if i < len(commands)-1 {
			fmt.Print(", ")
		}
		i++
	}

	fmt.Println(" or help")
	fmt.Println()

	if provider != nil {
		provider.printHelp(" ")
	} else {
		for id, provider := range commands {
			fmt.Printf("  %s\n", id)
			provider.printHelp("   ")
			fmt.Printf("    For more information run \"fyne help %s\"\n", id)
			fmt.Println("")
		}
	}
	flag.PrintDefaults()
}

func help() {
	printUsage()
	os.Exit(2) // consistent with flag.Parse() with -help
}

func loadCommands() {
	commands = make(map[string]command)

	commands["bundle"] = &bundler{}
	commands["get"] = &getter{}
	commands["package"] = &packager{}
	commands["install"] = &installer{}
	commands["vendor"] = &vendor{}
}

func main() {
	loadCommands()
	flag.Usage = printUsage
	// first let's extract the first "command"
	args := os.Args[1:]
	if len(args) < 1 {
		help()
	}

	command := args[0]
	if command[0] == '-' { // there was a parameter instead of a command
		help()
	}

	if command == "help" {
		if len(args) >= 2 {
			provider = commands[args[1]]
			provider.addFlags()
		}
		help()
	} else {
		provider = commands[command]
		if provider == nil {
			fmt.Fprintln(os.Stderr, "Unsupported command", command)
			return
		}

		provider.addFlags()

		// then parse the remaining args
		flag.CommandLine.Parse(args[1:])
		provider.run(flag.Args())
	}
}
