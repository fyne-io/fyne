// Run a command line helper for various Fyne tools.
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"fyne.io/fyne/v2/cmd/fyne/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "fyne",
		Usage:       "A command line helper for various Fyne tools.",
		Description: "The fyne command provides tooling for fyne applications and to assist in their development.",
		Commands: []*cli.Command{
			commands.Bundle(),
			commands.Env(),
			commands.Get(),
			commands.Install(),
			commands.Package(),
			commands.Release(),
			commands.Version(),
			commands.Serve(),
			commands.Build(),

			// Deprecated: Use "go mod vendor" instead.
			commands.Vendor(), //lint:ignore SA1019 This whole tool is deprecated.
		},
	}

	if info, ok := debug.ReadBuildInfo(); !ok {
		app.Version = "could not retrieve version information (ensure module support is activated and build again)"
	} else {
		app.Version = info.Main.Version
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
