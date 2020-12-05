package main

import (
	"fmt"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/cmd/fyne/commands"

	"github.com/lucor/goinfo"
	"github.com/lucor/goinfo/format"
	"github.com/lucor/goinfo/report"
)

const (
	fyneModule = "fyne.io/fyne"
)

// Declare conformity to Command interface
var _ commands.Command = (*env)(nil)

type env struct {
}

func (v *env) AddFlags() {
}

func (v *env) PrintHelp(indent string) {
	fmt.Println(indent, "The env command prints the Fyne module and environment information")
	fmt.Println(indent, "Command usage: fyne env")
}

func (v *env) Run(args []string) {
	if len(args) != 0 {
		fyne.LogError("Unexpected parameter after flags", nil)
		return
	}

	v.main()
}

func (v *env) main() {
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("could not get the path for the current working dir: %v", err))
		os.Exit(1)
	}

	reporters := []goinfo.Reporter{
		&fyneReport{GoMod: &report.GoMod{WorkDir: workDir, Module: fyneModule}},
		&report.GoVersion{},
		&report.GoEnv{Filter: []string{"GOOS", "GOARCH", "CGO_ENABLED", "GO111MODULE"}},
	}

	err = goinfo.Write(os.Stdout, reporters, &format.Text{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// fyneReport defines a custom report for fyne
type fyneReport struct {
	*report.GoMod
}

// Info returns the collected info
func (r *fyneReport) Info() (goinfo.Info, error) {
	info, err := r.GoMod.Info()
	if err != nil {
		return info, err
	}
	// remove info for the report
	delete(info, "module")

	versionCmd := &version{}
	cliVersion, err := versionCmd.get()
	if err != nil {
		cliVersion = err.Error()
	}
	info["cli_version"] = cliVersion
	return info, nil
}

// Summary returns the summary for the fyneReport
func (r *fyneReport) Summary() string {
	return "Fyne info"
}
