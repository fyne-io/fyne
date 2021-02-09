package commands

import (
	"os"
	"os/exec"

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
