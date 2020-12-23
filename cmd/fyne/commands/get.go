package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"fyne.io/fyne/cmd/fyne/internal/util"
)

// Declare conformity to Command interface
var _ Command = (*Getter)(nil)

// Getter is the command that can handle downloading and installing Fyne apps to the current platform.
type Getter struct {
	icon string
}

// NewGetter returns a command that can handle the download and install of GUI apps built using Fyne.
// It depends on a Go and C compiler installed at this stage and takes a single, package, parameter to identify the app.
func NewGetter() *Getter {
	return &Getter{}
}

// Get automates the download and install of a named GUI app package.
func (g *Getter) Get(pkg string) error {
	return get(pkg, g.icon)
}

// AddFlags adds available flags to the current flags parser
func (g *Getter) AddFlags() {
}

// PrintHelp prints help for this command when used in a command-line context
func (g *Getter) PrintHelp(indent string) {
	fmt.Println(indent, "The get command downloads and installs a Fyne application.")
	fmt.Println(indent, "A single parameter is required to specify the Go package, as with \"go get\"")
}

// Run is the command-line version of the Get(pkg) command, the first of the passed arguments will be used
// as the package to get.
func (g *Getter) Run(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Missing \"package\" argument to define the application to get")
		os.Exit(1)
	}

	err := get(args[0], g.icon)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

// SetIcon allows you to set the app icon path that will be used for the next Get operation.
func (g *Getter) SetIcon(path string) {
	g.icon = path
}

func get(pkg, icon string) error {
	path := filepath.Join(goPath(), "src", pkg)

	cmd := exec.Command("go", "get", "-u", "-d", pkg)
	cmd.Env = append(os.Environ(), "GO111MODULE=off")

	_, err := cmd.CombinedOutput()

	if !util.Exists(path) { // the error above may be ignorable, unless the path was not found
		return err
	}

	install := &installer{srcDir: path, icon: icon, release: true}
	err = install.validate()
	if err != nil {
		return errors.Wrap(err, "Failed to set up installer")
	}
	return install.install()
}

func goPath() string {
	cmd := exec.Command("go", "env", "GOPATH")
	out, err := cmd.CombinedOutput()

	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "go")
	}

	return strings.TrimSpace(string(out[0 : len(out)-1]))
}
