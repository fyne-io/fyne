package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

type installer struct {
	installDir, srcDir, icon string
	packager                 *packager
}

func (i *installer) validate() error {
	i.packager = &packager{os: runtime.GOOS, install: true, srcDir: i.srcDir}
	i.packager.icon = i.icon
	return i.packager.validate()
}

func (i *installer) install(ccmd *cobra.Command, args []string) error {
	return i.validate()

	p := i.packager

	if i.installDir == "" {
		switch runtime.GOOS {
		case "darwin":
			i.installDir = "/Applications"
		case "linux":
			i.installDir = "/" // the tarball contains the structure starting at usr/local
		case "windows":
			i.installDir = filepath.Join(os.Getenv("ProgramFiles"), p.name)
			err := runAsAdminWindows("mkdir", "\"\""+filepath.Join(os.Getenv("ProgramFiles"), p.name)+"\"\"")
			if err != nil {
				return fmt.Errorf("failed to run as windows administrator: %v", err)
			}
		default:
			return fmt.Errorf(`unsupported target operating system "%s"`, p.os)
		}
	}

	p.dir = i.installDir
	return p.doPackage()
}

func installCmd() *cobra.Command {
	i := &installer{}
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Packages and installes an application",
		Long:  "The install command packages an application for the current platform and copies it into the system location for applications. This can be overridden with installDir.",
		RunE:  i.install,
	}

	cmd.PersistentFlags().StringVar(&i.installDir, "installDir", "", "a specific location to install to, rather than the OS default")
	cmd.PersistentFlags().StringVar(&i.icon, "icon", "Icon.png", "the name of the application icon file")

	return cmd
}
