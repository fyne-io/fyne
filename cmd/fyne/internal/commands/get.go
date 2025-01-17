package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/sys/execabs"

	//lint:ignore SA1019 The recommended replacement does not solve the use-case
	"golang.org/x/tools/go/vcs"
)

// Get returns the command which downloads and installs fyne applications.
func Get() *cli.Command {
	g := &Getter{appData: &appData{}}
	return &cli.Command{
		Name:        "get",
		Usage:       "Downloads and installs a Fyne application",
		Description: "A single parameter is required to specify the Go package, as with \"go get\".",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "icon",
				Usage:       "The name of the application icon file.",
				Value:       "",
				Destination: &g.icon,
			},
			&cli.StringFlag{
				Name:        "appID",
				Aliases:     []string{"id"},
				Usage:       "For darwin and Windows targets an appID in the form of a reversed domain name is required, for ios this must match a valid provisioning profile",
				Destination: &g.AppID,
			},
			&cli.StringFlag{
				Name:        "installDir",
				Aliases:     []string{"o"},
				Usage:       "A specific location to install to, rather than the OS default.",
				Destination: &g.installDir,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("missing \"package\" argument to define the application to get")
			}

			pkg := ctx.Args().Slice()[0]
			return g.Get(pkg)
		},
	}
}

// Getter is the command that can handle downloading and installing Fyne apps to the current platform.
type Getter struct {
	*appData
	installDir string
}

// NewGetter returns a command that can handle the download and install of GUI apps built using Fyne.
// It depends on a Go and C compiler installed at this stage.
func NewGetter() *Getter {
	return &Getter{appData: &appData{}}
}

// Get automates the download and install of a named GUI app package.
func (g *Getter) Get(pkg string) error {
	wd, _ := os.Getwd()
	defer func() {
		if wd != "" {
			os.Chdir(wd)
		}
	}()

	name := filepath.Base(pkg)
	path, err := os.MkdirTemp("", fmt.Sprintf("fyne-get-%s-*", name))
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)

	repo, err := vcs.RepoRootForImportPath(pkg, false)
	if err != nil {
		return fmt.Errorf("failed to look up source control for package: %w", err)
	}
	if repo.VCS.Name != "Git" {
		return errors.New("unsupported VCS: " + repo.VCS.Name)
	}
	cmd := execabs.Command("git", "clone", repo.Repo, "--depth=1", path)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	if !util.Exists(path) { // the error above may be ignorable, unless the path was not found
		return err
	}

	if repo.Root != pkg {
		dir := strings.Replace(pkg, repo.Root, "", 1)
		path = filepath.Join(path, dir)
	}

	install := &Installer{appData: g.appData, installDir: g.installDir, srcDir: path, release: true}
	if err := install.validate(); err != nil {
		return fmt.Errorf("failed to set up installer: %w", err)
	}

	return install.install()
}

// SetAppID Allows the Get operation to specify an appID for the application that is being downloaded.
//
// Since: 2.1
func (g *Getter) SetAppID(id string) {
	g.AppID = id
}

// SetIcon allows you to set the app icon path that will be used for the next Get operation.
func (g *Getter) SetIcon(path string) {
	g.icon = path
}

// AddFlags adds available flags to the current flags parser
//
// Deprecated: Get does not define any flags.
func (g *Getter) AddFlags() {
	flag.StringVar(&g.icon, "icon", "Icon.png", "The name of the application icon file")
}

// PrintHelp prints help for this command when used in a command-line context
//
// Deprecated: Use Get() to get the cli and help messages instead.
func (g *Getter) PrintHelp(indent string) {
	fmt.Println(indent, "The get command downloads and installs a Fyne application.")
	fmt.Println(indent, "A single parameter is required to specify the Go package, as with \"go get\"")
}

// Run is the command-line version of the Get(pkg) command, the first of the passed arguments will be used
// as the package to get.
//
// Deprecated: Use Get() for the urfave/cli command or Getter.Get() to download and install an application.
func (g *Getter) Run(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Missing \"package\" argument to define the application to get")
		os.Exit(1)
	}

	if err := g.Get(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
