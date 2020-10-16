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
var _ Command = (*getter)(nil)

type getter struct {
	pkg string
}

// NewGetter returns a command that can handle the download and install of GUI apps built using Fyne.
// It depends on a Go and C compiler installed at this stage and takes a single, package, parameter to identify the app.
func NewGetter() Command {
	return &getter{}
}

func (g *getter) Get(pkg string) error {
	g.pkg = pkg
	return g.get()
}

func (g *getter) AddFlags() {
}

func (g *getter) PrintHelp(indent string) {
	fmt.Println(indent, "The get command downloads and installs a Fyne application.")
	fmt.Println(indent, "A single parameter is required to specify the Go package, as with \"go get\"")
}

func (g *getter) Run(args []string) {
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

func (g *getter) get() error {
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
