package main

import (
	"fmt"
	"os"

	"github.com/lucor/goinfo"
	"github.com/lucor/goinfo/format"
	"github.com/lucor/goinfo/report"
	"github.com/spf13/cobra"
)

const fyneModule = "fyne.io/fyne"

// fyneReport defines a custom report for fyne
type fyneReport struct {
	*report.GoMod
}

// Summary returns the summary for the fyneReport
func (r *fyneReport) Summary() string {
	return "Fyne info"
}

// Info returns the collected info
func (r *fyneReport) Info() (goinfo.Info, error) {
	info, err := r.GoMod.Info()
	if err != nil {
		return info, err
	}
	// remove info for the report
	delete(info, "module")

	info["cli_version"] = version
	return info, nil
}

func env(ccmd *cobra.Command, args []string) error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get the path for the current working dir: %v", err)
	}

	reporters := []goinfo.Reporter{
		&fyneReport{GoMod: &report.GoMod{WorkDir: workDir, Module: fyneModule}},
		&report.GoVersion{},
		&report.OS{},
		&report.GoEnv{Filter: []string{"GOOS", "GOARCH", "CGO_ENABLED", "GO111MODULE"}},
	}

	err = goinfo.Write(os.Stdout, reporters, &format.Text{})
	if err != nil {
		return err
	}

	return nil
}

func envCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Shows Fyne module and environment information",
		Long:  "The env command shows the Fyne module and environment",
		RunE:  env,
	}
}
