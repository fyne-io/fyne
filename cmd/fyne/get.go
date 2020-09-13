package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func goPath() string {
	cmd := exec.Command("go", "env", "GOPATH")
	out, err := cmd.CombinedOutput()
	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "go")
	}

	return strings.TrimSpace(string(out[:len(out)-1]))
}

func get(ccmd *cobra.Command, args []string) error {
	path := filepath.Join(goPath(), "src", args[0])

	cmd := exec.Command("go", "get", "-u", args[0])
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if !exists(path) {
		return errors.New("package download not found in expected location")
	}

	install := &installer{srcDir: path}
	err = install.validate()
	if err != nil {
		return fmt.Errorf("failed to set up installer: %v", err)
	}
	return install.install(nil, nil)

	return nil
}

func getCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Downloads and installs a Fyne application",
		Long:  `The get command downloads and installs a Fyne application. A single parameter is required to specify the Go package, as with "go get"`,
		RunE:  get,
	}
}
