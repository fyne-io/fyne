package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

var version = "master"

func usage() {
	fmt.Printf("\n")
	pflag.Usage()
	fmt.Printf(`
You can also pipe to stdin and from stdout. 
The pipes will be detected automatically, and both --input and --output will be ignored.

		cat icon.png | icnsify | cat > icon.icns

`)
}
