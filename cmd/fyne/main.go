package main

import (
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

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
