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
	"strings"

	"fyne.io/fyne"
)

const (
	// goModFile is the go module file
	goModFile = "go.mod"
	// goModVendorFile is the vendor module file
	goModVendorFile = "vendor/modules.txt"
	// glwfMod is the glfw module
	glwfMod = "github.com/go-gl/glfw"
	// glwfModSrcDir is the glfw dir containing the c source code to copy
	glwfModSrcDir = "v3.2/glfw/glfw"
)

// Declare conformity to command interface
var _ command = (*vendor)(nil)

type vendor struct {
}

func (v *vendor) addFlags() {
}

func (v *vendor) printHelp(indent string) {
	// TODO fix
	fmt.Println(indent, "The vendor command packages an application's dependencies and all the extra")
	fmt.Println(indent, "files required into it's vendor folder. Your project must have a " + goModFile + " file.")
	fmt.Println(indent, "Command usage: fyne vendor")
}

// cacheModPath returns the cache path for a module
func cacheModPath(r io.Reader, module string) (string, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		if len(s) != 3 {
			continue
		}
		if s[1] != module {
			continue
		}

		moduleVersion := module + "@" + s[2]
		return filepath.Join(build.Default.GOPATH, "pkg/mod", moduleVersion), nil
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Cannot read content: %v", err)
	}
	return "", nil
}

func recursiveCopy(src, target string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("Cannot read dir %s: %v", src, err)
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

		err = copyFile(srcFile, targetFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err2 := os.Mkdir(dir, 0755)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func (v *vendor) main() {
	wd, _ := os.Getwd()

	if _, err := os.Stat(filepath.Join(wd, goModFile)); os.IsNotExist(err) {
		fmt.Println("This program must be invoked in the project root")
		os.Exit(1)
	}

	fmt.Println("Reset vendor content using 'go mod vendor'")
	cmd := exec.Command("go", "mod", "vendor", "-v")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsing %s to detect dependency module path for %s\n", goModVendorFile, glwfMod)
	f, err := os.Open(filepath.Join(wd, goModVendorFile))
	if err != nil {
		fmt.Printf("Cannot open %s: %v\n", goModVendorFile, err)
		os.Exit(1)
	}
	defer f.Close()

	glwfModPath, err := cacheModPath(f, glwfMod)
	if err != nil {
		fmt.Printf("Cannot read %s: %v\n", goModVendorFile, err)
		os.Exit(1)
	}
	if glwfModPath == "" {
		fmt.Printf("Cannot find module %s in %s\n", glwfMod, goModVendorFile)
		os.Exit(1)
	}
	fmt.Printf("Package module path: %s\n", glwfModPath)

	// glwfModSrc is the path containing glfw c source code
	// $GOPATH/pkg/mod/github.com/go-gl/glfw@v...../v3.3/glfw/glfw
	glwfModSrc := filepath.Join(glwfModPath, glwfModSrcDir)
	// glwfModTarget is the path under the vendor folder where glfw c source code will be copied
	// vendor/github.com/go-gl/v3.3/glfw/glfw
	glwfModTarget := filepath.Join(wd, "vendor", glwfMod, glwfModSrcDir)
	fmt.Printf("Copying glwf c source code: %s -> %s\n", glwfModSrc, glwfModTarget)
	err = recursiveCopy(glwfModSrc, glwfModTarget)
	if err != nil {
		fmt.Printf("Cannot copy glfw c source code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All set. To test using vendor dependencies: go test ./... -mod=vendor -v -count=1")
}

func (v *vendor) run(args []string) {
	if len(args) != 0 {
		fyne.LogError("Unexpected parameter after flags", nil)
		return
	}

	v.main()
}