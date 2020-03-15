// Run a command line helper for various Fyne tools.
package main

import (
	"flag"
	"fmt"
	"os"
)

var commands []idCommandPair
var provider command

type idCommandPair struct {
	id       string
	provider command
}

func printUsage() {
	fmt.Println("Usage: fyne [command] [parameters], where command is one of:")
	fmt.Print("  ")

	i := 0
	for _, c := range commands {
		fmt.Print(c.id)

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
		for _, c := range commands {
			fmt.Printf("  %s\n", c.id)
			c.provider.printHelp("   ")
			fmt.Printf("    For more information run \"fyne help %s\"\n", c.id)
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
	commands = []idCommandPair{
		{"bundle", &bundler{}},
		{"get", &getter{}},
		{"package", &packager{}},
		{"install", &installer{}},
		{"vendor", &vendor{}},
	}
}

func getCommand(id string) command {
	for _, c := range commands {
		if c.id == id {
			return c.provider
		}
	}
	return nil
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
			if provider = getCommand(args[1]); provider != nil {
				provider.addFlags()
			}
		}
		help()
	} else {
		provider = getCommand(command)
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
