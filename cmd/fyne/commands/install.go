package commands

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/mobile"
	"github.com/urfave/cli/v2"
)

// Install returns the cli command for installing fyne applications
func Install() *cli.Command {
	i := &installer{}

	return &cli.Command{
		Name:  "install",
		Usage: "Packages an application and installs an application.",
		Description: `The install command packages an application for the current platform and copies it
		into the system location for applications. This can be overridden with installDir`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"os"},
				Usage:       "The mobile platform to target (android, android/arm, android/arm64, android/amd64, android/386, ios).",
				Destination: &i.os,
			},
			&cli.StringFlag{
				Name:        "installDir",
				Aliases:     []string{"o"},
				Usage:       "A specific location to install to, rather than the OS default.",
				Destination: &i.installDir,
			},
			&cli.StringFlag{
				Name:        "icon",
				Usage:       "The name of the application icon file.",
				Value:       "Icon.png",
				Destination: &i.icon,
			},
			&cli.StringFlag{
				Name:        "appiD",
				Aliases:     []string{"id"},
				Usage:       "For ios or darwin targets an appID in the form of a reversed domain name is required, for ios this must match a valid provisioning profile",
				Destination: &i.appID,
			},
			&cli.BoolFlag{
				Name:        "release",
				Usage:       "Enable installation in release mode (disable debug, etc).",
				Destination: &i.release,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 0 {
				return errors.New("Unexpected parameter after flags")
			}

			err := i.validate()
			if err != nil {
				return err
			}

			err = i.install()
			if err != nil {
				return err
			}

			return nil
		},
	}
}

type installer struct {
	installDir, srcDir, icon, os, appID string
	packager                            *packager
	release                             bool
}

func (i *installer) install() error {
	p := i.packager

	if i.os != "" {
		if i.os == "ios" {
			return i.installIOS()
		} else if strings.Index(i.os, "android") == 0 {
			return i.installAndroid()
		}

		return errors.New("Unsupported target operating system \"" + i.os + "\"")
	}

	if i.installDir == "" {
		switch p.os {
		case "darwin":
			i.installDir = "/Applications"
		case "linux", "openbsd", "freebsd", "netbsd":
			i.installDir = "/" // the tarball contains the structure starting at usr/local
		case "windows":
			i.installDir = filepath.Join(os.Getenv("ProgramFiles"), p.name)
			err := runAsAdminWindows("mkdir", "\"\""+filepath.Join(os.Getenv("ProgramFiles"), p.name)+"\"\"")
			if err != nil {
				fyne.LogError("Failed to run as windows administrator", err)
				return err
			}
		default:
			return errors.New("Unsupported target operating system \"" + p.os + "\"")
		}
	}

	p.dir = i.installDir
	return p.doPackage()
}

func (i *installer) installAndroid() error {
	target := mobile.AppOutputName(i.os, i.packager.name)

	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		err := i.packager.doPackage()
		if err != nil {
			return nil
		}
	}

	return i.runMobileInstall("adb", target, "install")
}

func (i *installer) installIOS() error {
	target := mobile.AppOutputName(i.os, i.packager.name)
	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		err := i.packager.doPackage()
		if err != nil {
			return nil
		}
	}

	return i.runMobileInstall("ios-deploy", target, "--bundle")
}

func (i *installer) runMobileInstall(tool, target string, args ...string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		return err
	}

	cmd := exec.Command(tool, append(args, target)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (i *installer) validate() error {
	os := i.os
	if os == "" {
		os = targetOS()
	}
	i.packager = &packager{appID: i.appID, os: os, install: true, srcDir: i.srcDir}
	i.packager.icon = i.icon
	i.packager.release = i.release
	return i.packager.validate()
}
