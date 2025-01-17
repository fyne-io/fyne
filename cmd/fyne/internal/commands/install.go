package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/mobile"

	"github.com/urfave/cli/v2"
	"golang.org/x/sys/execabs"
)

// Install returns the cli command for installing fyne applications
func Install() *cli.Command {
	i := NewInstaller()

	return &cli.Command{
		Name:  "install",
		Usage: "Packages an application and installs an application.",
		Description: `The install command packages an application for the current platform or one of the mobile targets.
It will copy the package to the system location for applications or a location specified by installDir.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"os"},
				Usage:       "Instead of the current system, target a mobile platform (android, android/arm, android/arm64, android/amd64, android/386, ios, iossimulator).",
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
				Value:       "",
				Destination: &i.icon,
			},
			&cli.BoolFlag{
				Name:        "use-raw-icon",
				Usage:       "Skip any OS-specific icon pre-processing",
				Value:       false,
				Destination: &i.rawIcon,
			},
			&cli.StringFlag{
				Name:        "appID",
				Aliases:     []string{"id"},
				Usage:       "For Android, darwin, iOS and Windows targets an appID in the form of a reversed domain name is required, for ios this must match a valid provisioning profile",
				Destination: &i.AppID,
			},
			&cli.BoolFlag{
				Name:        "release",
				Usage:       "Enable installation in release mode (disable debug, etc).",
				Destination: &i.release,
			},
		},
		Action: i.bundleAction,
	}
}

// Installer installs locally built Fyne apps.
type Installer struct {
	*appData
	installDir, srcDir, os string
	Packager               *Packager
	release                bool
}

// NewInstaller returns a command that can install a GUI apps built using Fyne from local source code.
func NewInstaller() *Installer {
	return &Installer{appData: &appData{}}
}

// AddFlags adds the flags for interacting with the Installer.
//
// Deprecated: Access to the individual cli commands are being removed.
func (i *Installer) AddFlags() {
	flag.StringVar(&i.os, "os", "", "The mobile platform to target (android, android/arm, android/arm64, android/amd64, android/386, ios)")
	flag.StringVar(&i.installDir, "installDir", "", "A specific location to install to, rather than the OS default")
	flag.StringVar(&i.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&i.AppID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.BoolVar(&i.release, "release", false, "Should this package be installed in release mode? (disable debug etc)")
}

// PrintHelp prints the help for the install command.
//
// Deprecated: Access to the individual cli commands are being removed.
func (i *Installer) PrintHelp(indent string) {
	fmt.Println(indent, "The install command packages an application for the current platform and copies it")
	fmt.Println(indent, "into the system location for applications. This can be overridden with installDir")
	fmt.Println(indent, "Command usage: fyne install [parameters]")
}

// Run runs the install command.
//
// Deprecated: A better version will be exposed in the future.
func (i *Installer) Run(args []string) {
	if len(args) != 0 {
		fyne.LogError("Unexpected parameter after flags", nil)
		return
	}

	err := i.validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	err = i.install()
	if err != nil {
		fyne.LogError("Unable to install application", err)
		os.Exit(1)
	}
}

func (i *Installer) bundleAction(ctx *cli.Context) error {
	if ctx.Args().Len() != 0 {
		return errors.New("unexpected parameter after flags")
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
}

func (i *Installer) install() error {
	p := i.Packager

	if i.os != "" {
		if util.IsIOS(i.os) {
			return i.installIOS()
		} else if strings.Index(i.os, "android") == 0 {
			return i.installAndroid()
		}

		return errors.New("Unsupported target operating system \"" + i.os + "\".\nTo install on the current platform simply omit the target/os parameter.")
	}

	if i.installDir == "" {
		switch p.os {
		case "darwin":
			i.installDir = "/Applications"
		case "linux", "openbsd", "freebsd", "netbsd":
			i.installDir = "/" // the tarball contains the structure starting at usr/local
		case "windows":
			dirName := p.Name
			if filepath.Ext(p.Name) == ".exe" {
				dirName = p.Name[:len(p.Name)-4]
			}
			i.installDir = filepath.Join(os.Getenv("ProgramFiles"), dirName)
			err := runAsAdminWindows("mkdir", i.installDir)
			if err != nil {
				fyne.LogError("Failed to run as windows administrator", err)
				return err
			}
		default:
			return errors.New("Unsupported target operating system \"" + p.os + "\"")
		}
	}

	p.dir = i.installDir
	err := p.doPackage(nil)
	if err != nil {
		return err
	}

	return postInstall(i)
}

func (i *Installer) installAndroid() error {
	target := mobile.AppOutputName(i.os, i.Packager.Name, i.release)

	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		err := i.Packager.doPackage(nil)
		if err != nil {
			return nil
		}
	}

	return i.runMobileInstall("adb", target, "install")
}

func (i *Installer) installIOS() error {
	target := mobile.AppOutputName(i.os, i.Packager.Name, i.release)

	// Always redo the package because the codesign for ios and iossimulator
	// must be different.
	if err := i.Packager.doPackage(nil); err != nil {
		return nil
	}

	switch i.os {
	case "ios":
		return i.runMobileInstall("ios-deploy", target, "--bundle")
	case "iossimulator":
		return i.installToIOSSimulator(target)
	default:
		return fmt.Errorf("unsupported install target: %s", target)
	}
}

func (i *Installer) runMobileInstall(tool, target string, args ...string) error {
	_, err := execabs.LookPath(tool)
	if err != nil {
		return err
	}

	cmd := execabs.Command(tool, append(args, target)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (i *Installer) validate() error {
	os := i.os
	if os == "" {
		os = targetOS()
	}
	i.Packager = &Packager{appData: i.appData, os: os, install: true, srcDir: i.srcDir}
	i.Packager.AppID = i.AppID
	i.Packager.icon = i.icon
	i.Packager.release = i.release
	return i.Packager.validate()
}

func (i *Installer) installToIOSSimulator(target string) error {
	cmd := execabs.Command(
		"xcrun", "simctl", "install",
		"booted", // Install to the booted simulator.
		target)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Install to a simulator error: %s%s", out, err)
	}

	i.runInIOSSimulator()
	return nil
}

func (i *Installer) runInIOSSimulator() error {
	cmd := execabs.Command("xcrun", "simctl", "launch", "booted", i.Packager.AppID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return err
	}
	return nil
}
