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
	pkg string
}

// NewGetter returns a command that can handle the download and install of GUI apps built using Fyne.
// It depends on a Go and C compiler installed at this stage and takes a single, package, parameter to identify the app.
func NewGetter() *Getter {
	return &Getter{}
}

// Get automates the download and install of a named GUI app package.
func (g *Getter) Get(pkg string) error {
	g.pkg = pkg
	return g.get()
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
	g.pkg = args[0]

	err := g.get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func (g *Getter) get() error {
	path := filepath.Join(goPath(), "src", g.pkg)

	cmd := exec.Command("go", "get", "-u", g.pkg)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(out))
		return err
	}

	if !util.Exists(path) {
		return errors.New("Package download not found in expected location")
	}

	install := &installer{srcDir: path}
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
