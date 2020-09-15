package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fyne",
	Short: "A command line helper for various Fyne tools",
	Long:  "The fyne command provides tooling for fyne applications and their development",
}

var version string

// Deprecated: Replace with --version in the future
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Deprecated: Prints version information",
		RunE: func(ccmd *cobra.Command, args []string) error {
			fmt.Println("fyne cli version", version)
			fmt.Println("Deprecated: use --version or -v in the future")
			return nil
		},
	}
}

func main() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatalln("could not retrieve version information (ensure module support is activated and build again)")
	}
	version = info.Main.Version
	rootCmd.Version = version

	rootCmd.AddCommand(getCmd())
	rootCmd.AddCommand(vendorCmd())
	rootCmd.AddCommand(envCmd())
	rootCmd.AddCommand(installCmd())
	rootCmd.AddCommand(bundleCmd())
	rootCmd.AddCommand(packageCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
