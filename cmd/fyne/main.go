// Run a command line helper for various Fyne tools.
package main

import (
	"flag"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/commands"
)

var commandList []idCommandPair
var provider commands.Command

type idCommandPair struct {
	id       string
	provider commands.Command
}

func getCommand(id string) commands.Command {
	for _, c := range commandList {
		if c.id == id {
			return c.provider
		}
	}
	return nil
}

func help() {
	printUsage()
	os.Exit(2) // consistent with flag.Parse() with -help
}

func loadCommands() {
	commandList = []idCommandPair{
		{"bundle", commands.NewBundler()},
		{"get", commands.NewGetter()},
		{"env", &env{}},
		{"package", commands.NewPackager()},
		{"install", commands.NewInstaller()},
		{"release", commands.NewReleaser()},
		{"vendor", &vendor{}},
		{"version", &version{}},
	}
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
				provider.AddFlags()
			}
		}
		help()
	} else {
		provider = getCommand(command)
		if provider == nil {
			fmt.Fprintln(os.Stderr, "Unsupported command", command)
			return
		}

		provider.AddFlags()

		// then parse the remaining args
		err := flag.CommandLine.Parse(args[1:])
		if err != nil {
			fyne.LogError("Failed to parse flags", err)
			return
		}

		provider.Run(flag.Args())
	}
}

func printUsage() {
	fmt.Println("Usage: fyne [command] [parameters], where command is one of:")
	fmt.Print("  ")

	i := 0
	for _, c := range commandList {
		fmt.Print(c.id)

		if i < len(commandList)-1 {
			fmt.Print(", ")
		}
		i++
	}

	fmt.Println(" or help")
	fmt.Println()

	if provider != nil {
		provider.PrintHelp(" ")
	} else {
		for _, c := range commandList {
			fmt.Printf("  %s\n", c.id)
			c.provider.PrintHelp("   ")
			fmt.Printf("    For more information run \"fyne help %s\"\n", c.id)
			fmt.Println("")
		}
	}
	flag.PrintDefaults()
}
