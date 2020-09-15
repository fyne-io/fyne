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

	"github.com/spf13/cobra"
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

		err = copyFile(srcFile, targetFile)
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

func vendor(ccmd *cobra.Command, args []string) error {
	wd, _ := os.Getwd()

	if _, err := os.Stat(filepath.Join(wd, goModFile)); os.IsNotExist(err) {
		return fmt.Errorf("this program must be invoked in the project root: %v", err)
	}

	fmt.Println("Add missing and remove unused modules using 'go mod tidy'")
	cmd := exec.Command("go", "mod", "tidy", "-v")
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Reset vendor content using 'go mod vendor'")
	cmd = exec.Command("go", "mod", "vendor", "-v")
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("Parsing %s to detect dependency module path for GLFW\n", goModVendorFile)
	f, err := os.Open(filepath.Join(wd, goModVendorFile)) // #nosec
	if err != nil {
		return fmt.Errorf("cannot open %s: %v", goModVendorFile, err)
	}
	defer f.Close()

	glwfModPath, glfwModDest, err := cacheModPath(f)
	if err != nil {
		return fmt.Errorf("cannot read %s: %v", goModVendorFile, err)
	}
	if glwfModPath == "" {
		return fmt.Errorf("cannot find GLFW module in %s", goModVendorFile)
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
		return fmt.Errorf("cannot copy glfw c source code: %v", err)
	}

	fmt.Println("All set. To test using vendor dependencies: go test ./... -mod=vendor -v -count=1")
	return nil
}

func vendorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "vendor",
		Short: "Packages application dependencies to vendor folder",
		Long:  "The vendor command packages an application's dependencies and all the extra files required into it's vendor folder. Your project must have a " + goModFile + " file.",
		RunE:  vendor,
	}
}
