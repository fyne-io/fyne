package report

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lucor/goinfo"
	"golang.org/x/mod/modfile"
)

const (
	gomod = "go.mod"
)

// GoMod returns the info about a Go module
// The info returned are defined into goModInfo
type GoMod struct {
	//Module is the module to look for into WorkDir and GOPATH as fallback.
	// If not specified the module present in WorkDir will be returned, if any.
	Module string
	// WorkDir is the directory where to look for the go module. If not
	// specified, default to current directory.
	WorkDir string
}

// Summary return the summary
func (i *GoMod) Summary() string {
	return "Go module info"
}

// Info returns the collected info
func (i *GoMod) Info() (goinfo.Info, error) {
	// Start to look for go module into WorkDir
	info, err := getModInfo(i.Module, i.WorkDir)
	if info == nil && i.Module != "" {
		// try to fallback with GOPATH
		path := filepath.Join(build.Default.GOPATH, "src", i.Module)
		info, err = getModInfo(i.Module, path)
	}
	return info, err
}

// getModInfo returns info about module in dir
func getModInfo(module string, searchDir string) (goinfo.Info, error) {

	info := goinfo.Info{
		// module is the module found in searchDir, if any
		"module": module,
		// path	 is the directory where the module is found
		"path": searchDir,
		// go_mod is true if a go module is found in searchDir
		"go_mod": false,
		// go_path is true if the go module is found under GOPATH
		"go_path": false,
		// imported is true if the go module is imported and used as third-party library
		"imported": false,
		// version is the version detected for Module reported as:
		// a) VCS info (i.e. v1.2.3-116-ga5fcb7eb-dirty) when the code in searchDir contains a copy of the module
		// b) defined in go.mod (i.e. v1.2.2) when the code is imported and used as third-party library
		"version": "",
	}

	// read the go.mod content, if exists
	data, err := ioutil.ReadFile(filepath.Join(searchDir, gomod))
	if err != nil {
		return nil, fmt.Errorf("go.mod not found in %s", searchDir)
	}

	// parse the go.mod content
	f, err := modfile.Parse("", data, nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse the go.mod at %s: %v", searchDir, err)
	}

	// A module was found in searchDir
	info["go_mod"] = true

	// No module was specified let's set as the one found in go.mod
	if module == "" {
		module = f.Module.Mod.Path
	}

	// when go.mod is related to module returns the commit version info
	if module == f.Module.Mod.Path {
		info["module"] = module
		commitVersion, err := vcsInfo(searchDir)
		if err != nil {
			return info, err
		}
		info["version"] = commitVersion
		return info, nil
	}

	// when the module is used in go.mod as third-party returns the version found in the go.mod, if any
	for _, v := range f.Require {
		if module == v.Mod.Path {
			info["imported"] = true
			info["version"] = v.Mod.Version
			return info, nil
		}
	}
	return nil, fmt.Errorf("could not found info for %s", module)
}

func vcsInfo(workDir string) (string, error) {
	if _, err := os.Stat(filepath.Join(workDir, ".git")); !os.IsNotExist(err) {
		info, err := gitInfo(workDir)
		if err != nil {
			return "", fmt.Errorf("could not detect commit info in %s using git describe: %w", workDir, err)
		}
		return info, nil
	}
	return "", nil
}

// gitInfo returns the "human-readable" commit name using git describe
func gitInfo(workDir string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--dirty", "--always")
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	s := strings.Trim(string(out), "\n")
	if err != nil {
		return "", fmt.Errorf("%s %s", s, err)
	}
	return s, nil
}
