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
		Run: func(ccmd *cobra.Command, args []string) {
			fmt.Println("fyne cli version", version)
			fmt.Println("Deprecated: use --version or -v in the future")
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

	rootCmd.AddCommand(bundleCmd())
	rootCmd.AddCommand(envCmd())
	rootCmd.AddCommand(getCmd())
	rootCmd.AddCommand(installCmd())
	rootCmd.AddCommand(packageCmd())
	rootCmd.AddCommand(vendorCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
