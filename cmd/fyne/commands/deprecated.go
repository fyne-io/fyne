package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

// Vendor returns the vendor cli command.
//
// Deprecated: Use "go mod vendor" instead.
func Vendor() *cli.Command {
	return &cli.Command{
		Name:  "vendor",
		Usage: "Deprecated: Use \"go mod vendor\" instead.",
		Action: func(_ *cli.Context) error {
			cmd := exec.Command("go", "mod", "vendor")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			return cmd.Run()
		},
	}
}

// Version returns the cli command for the program version.
//
// Deprecated: Use "fyne --version" or "fyne -v" instead.
func Version() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Deprecated: Use \"fyne --version\" or \"fyne -v\" instead.",
		Action: func(_ *cli.Context) error {
			if info, ok := debug.ReadBuildInfo(); ok {
				fmt.Println("fyne cli version:", info.Main.Version)
				return nil
			}

			return errors.New("could not retrieve version information (ensure module support is activated and build again)")
		},
	}
}

// Command defines the required functionality to provide a subcommand to the "fyne" tool.
//
// Deprecated: Use the functions that return urfave/cli commands instead.
type Command interface {
	AddFlags()
	PrintHelp(string)
	Run(args []string)
}
