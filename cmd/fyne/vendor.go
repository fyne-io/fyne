package main

import (
	"bufio"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/commands"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
)

const (
	// goModFile is the go module file
	goModFile = "go.mod"
	// goModVendorFile is the vendor module file
	goModVendorFile = "vendor/modules.txt"
	// glfwModNew is the glfw module regexp for 3.3+
	glfwModNew = `github\.com/go-gl/glfw/(v(\d+\.)?(\*|\d)?(\*|\d+))/glfw`
	// glfwModOld is the glfw module regexp for 3.2
	glfwModOld = "github.com/go-gl/glfw"
	// glfwModSrcDir is the glfw dir containing the c source code to copy
	glfwModSrcDir = "glfw"
)

// Declare conformity to command interface
var _ commands.Command = (*vendor)(nil)

type vendor struct {
}

func (v *vendor) AddFlags() {
}

func (v *vendor) PrintHelp(indent string) {
	fmt.Println(indent, "The vendor command packages an application's dependencies and all the extra")
	fmt.Println(indent, "files required into it's vendor folder. Your project must have a "+goModFile+" file.")
	fmt.Println(indent, "Command usage: fyne vendor")
}

func (v *vendor) Run(args []string) {
	if len(args) != 0 {
		fyne.LogError("Unexpected parameter after flags", nil)
		return
	}

	v.main()
}

func (v *vendor) main() {
	wd, _ := os.Getwd()

	if _, err := os.Stat(filepath.Join(wd, goModFile)); os.IsNotExist(err) {
		fmt.Println("This program must be invoked in the project root")
		os.Exit(1)
	}

	fmt.Println("Add missing and remove unused modules using 'go mod tidy'")
	cmd := exec.Command("go", "mod", "tidy", "-v")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Println("Reset vendor content using 'go mod vendor'")
	cmd = exec.Command("go", "mod", "vendor", "-v")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsing %s to detect dependency module path for GLFW\n", goModVendorFile)
	f, err := os.Open(filepath.Join(wd, goModVendorFile)) // #nosec
	if err != nil {
		fmt.Printf("Cannot open %s: %v\n", goModVendorFile, err)
		os.Exit(1)
	}
	defer f.Close()

	glwfModPath, glfwModDest, err := cacheModPath(f)
	if err != nil {
		fmt.Printf("Cannot read %s: %v\n", goModVendorFile, err)
		os.Exit(1)
	}
	if glwfModPath == "" {
		fmt.Printf("Cannot find GLFW module in %s\n", goModVendorFile)
		os.Exit(1)
	}
	fmt.Printf("Package module path: %s\n", glwfModPath)

	// glwfModSrc is the path containing glfw c source code
	// $GOPATH/pkg/mod/github.com/go-gl/glfw@v...../v3.3/glfw/glfw
	glwfModSrc := filepath.Join(glwfModPath, glfwModSrcDir)
	// glwfModTarget is the path under the vendor folder where glfw c source code will be copied
	// vendor/github.com/go-gl/v3.3/glfw/glfw
	glwfModTarget := filepath.Join(wd, "vendor", glfwModDest, glfwModSrcDir)
	fmt.Printf("Copying glwf c source code: %s -> %s\n", glwfModSrc, glwfModTarget)
	err = recursiveCopy(glwfModSrc, glwfModTarget)
	if err != nil {
		fmt.Printf("Cannot copy glfw c source code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All set. To test using vendor dependencies: go test ./... -mod=vendor -v -count=1")
}

// cacheModPath returns the cache path and target directory for the GLFW module
func cacheModPath(r io.Reader) (string, string, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		if len(s) != 3 {
			continue
		}

		r, _ := regexp.Compile(glfwModNew)
		if r.Match([]byte(s[1])) {
			return filepath.Join(build.Default.GOPATH, "pkg/mod", s[1]+"@"+s[2]), s[1], nil
		}

		r, _ = regexp.Compile(glfwModOld)
		if r.Match([]byte(s[1])) {
			extra := "/v3.2/glfw"
			return filepath.Join(build.Default.GOPATH, "pkg/mod", s[1]+"@"+s[2]+extra), s[1] + extra, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("cannot read content: %v", err)
	}
	return "", "", nil
}

func recursiveCopy(src, target string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("cannot read dir %s: %v", src, err)
	}
	err = ensureDir(target)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcFile := filepath.Join(src, file.Name())
		targetFile := filepath.Join(target, file.Name())
		if file.IsDir() {
			err = recursiveCopy(srcFile, targetFile)
			if err != nil {
				return err
			}
			continue
		}

		err = util.CopyFile(srcFile, targetFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err2 := os.Mkdir(dir, 0750)
		if err2 != nil {
			return err2
		}
	}
	return nil
}
