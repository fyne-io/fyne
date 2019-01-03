package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

type packager struct {
	os string
}

func (p *packager) packageDarwin() {

}

func (p *packager) packageWindows() {

}

func (p *packager) addFlags() {
	flag.StringVar(&p.os, "os", runtime.GOOS, "The operating system to target (like GOOS)")
}

func (*packager) printHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for distribution.")
	fmt.Println(indent, "Command usage: fyne package [parameters]")
}

func (p *packager) run(_ []string) {
	switch p.os {
	case "linux":
		fmt.Println("No packaging required for linux")
	case "darwin":
		p.packageDarwin()
	case "windows":
		p.packageWindows()
	default:
		fmt.Fprintln(os.Stderr, "Unsupported os parameter,", p.os)
	}
}
