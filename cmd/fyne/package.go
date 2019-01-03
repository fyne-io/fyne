package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

type packager struct {
	os, name, exe string
}

func (p *packager) packageLinux() {
	fmt.Println("No packaging required for linux")
}

func (p *packager) packageDarwin() {

}

func (p *packager) packageWindows() {

}

func (p *packager) addFlags() {
	flag.StringVar(&p.os, "os", runtime.GOOS, "The operating system to target (like GOOS)")
	flag.StringVar(&p.exe, "executable", "", "The path to the executable, default is the current dir main binary")
	flag.StringVar(&p.name, "name", "", "The name of the application, default is the executable file name")
}

func (*packager) printHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for distribution.")
	fmt.Println(indent, "Command usage: fyne package [parameters]")
}

func (p *packager) run(_ []string) {
	if p.exe == "" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to get the current directory, needed to find main executable")
		}

		exeName := filepath.Base(dir)
		if p.os == "windows" {
			exeName = exeName + ".exe"
		}

		p.exe = filepath.Join(dir, exeName)
	}

	if !exists(p.exe) {
		fmt.Fprintln(os.Stderr, "Executable", p.exe, "could not be found")
		os.Exit(1)
	}
	if p.name == "" {
		name := filepath.Base(p.exe)
		p.name = strings.TrimSuffix(name, filepath.Ext(name))
	}

	switch p.os {
	case "linux":
		p.packageLinux()
	case "darwin":
		p.packageDarwin()
	case "windows":
		p.packageWindows()
	default:
		fmt.Fprintln(os.Stderr, "Unsupported os parameter,", p.os)
	}
}
