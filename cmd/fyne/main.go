// Run a command line helper for various Fyne tools.
//
// Deprecated: Install fyne.io/tools/cmd/fyne for latest version.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(`NOTE: This tool is deprecated and has migrated to fyne.io/tools/cmd/fyne.
The new tool can be installed by running the following command:\n\tgo install fyne.io/tools/cmd/fyne@latest`)
	os.Exit(1)
}
