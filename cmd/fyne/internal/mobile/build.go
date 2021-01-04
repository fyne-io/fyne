// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gendex/gendex.go -o dex.go

package mobile

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

var tmpdir string

var cmdBuild = &command{
	run:   runBuild,
	Name:  "build",
	Usage: "[-os android|ios] [-o output] [-bundleid bundleID] [build flags] [package]",
	Short: "compile android APK and iOS app",
	Long: `
Build compiles and encodes the app named by the import path.

The named package must define a main function.

The -os Flag takes a target system name, either android (the
default) or ios.

For -os=android, if an AndroidManifest.xml is defined in the
package directory, it is added to the APK output. Otherwise, a default
manifest is generated. By default, this builds a fat APK for all supported
instruction sets (arm, 386, amd64, arm64). A subset of instruction sets can
be selected by specifying target type with the architecture name. E.g.
-os=android/arm,android/386.

For -os ios, gomobile must be run on an OS X machine with Xcode
installed.

If the package directory contains an assets subdirectory, its contents
are copied into the output.

Flag -iosversion sets the minimal version of the iOS SDK to compile against.
The default version is 7.0.

Flag -androidapi sets the Android API version to compile against.
The default and minimum is 15.

The -bundleid Flag is required for -os=ios and sets the bundle ID to use
with the app.

The -o Flag specifies the output file name. If not specified, the
output file name depends on the package built.

The -v Flag provides verbose output, including the list of packages built.

The build flags -a, -i, -n, -x, -gcflags, -ldflags, -tags, -trimpath, and -work are
shared with the build command. For documentation, see 'go help build'.
`,
}

const (
	minAndroidAPI = 15
)

func runBuild(cmd *command) (err error) {
	_, err = runBuildImpl(cmd)
	return
}

// AppOutputName provides the name of a build resource for a given os - "ios" or "android".
func AppOutputName(os, name string) string {
	switch os {
	case "android":
		return androidPkgName(name) + ".apk"
	case "ios":
		return rfc1034Label(name) + ".app"
	}

	return ""
}

// runBuildImpl builds a package for mobiles based on the given commands.
// runBuildImpl returns a built package information and an error if exists.
func runBuildImpl(cmd *command) (*packages.Package, error) {
	cleanup, err := buildEnvInit()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	args := cmd.Flag.Args()

	targetOS, targetArchs, err := parseBuildTarget(buildTarget)
	if err != nil {
		return nil, fmt.Errorf(`invalid -os=%q: %v`, buildTarget, err)
	}

	var buildPath string
	switch len(args) {
	case 0:
		buildPath = "."
	case 1:
		buildPath = args[0]
	default:
		cmd.usage()
		os.Exit(1)
	}
	pkgs, err := packages.Load(packagesConfig(targetOS, targetArchs[0]), buildPath)
	if err != nil {
		return nil, err
	}
	// len(pkgs) can be more than 1 e.g., when the specified path includes `...`.
	if len(pkgs) != 1 {
		cmd.usage()
		os.Exit(1)
	}

	pkg := pkgs[0]

	if pkg.Name != "main" && buildO != "" {
		return nil, fmt.Errorf("cannot set -o when building non-main package")
	}
	if buildBundleID == "" {
		return nil, fmt.Errorf("value for -appID is required for a mobile package")
	}

	var nmpkgs map[string]bool
	switch targetOS {
	case "android":
		if pkg.Name != "main" {
			for _, arch := range targetArchs {
				if err := goBuild(pkg.PkgPath, androidEnv[arch]); err != nil {
					return nil, err
				}
			}
			return pkg, nil
		}
		nmpkgs, err = goAndroidBuild(pkg, buildBundleID, targetArchs, cmd.IconPath, cmd.AppName, cmd.Version, cmd.Build)
		if err != nil {
			return nil, err
		}
	case "darwin":
		if !xcodeAvailable() {
			return nil, fmt.Errorf("-os=ios requires XCode")
		}
		if buildRelease {
			targetArchs = []string{"arm64"}
		}

		if pkg.Name != "main" {
			for _, arch := range targetArchs {
				if err := goBuild(pkg.PkgPath, darwinEnv[arch]); err != nil {
					return nil, err
				}
			}
			return pkg, nil
		}
		nmpkgs, err = goIOSBuild(pkg, buildBundleID, targetArchs, cmd.AppName, cmd.Version, cmd.Build, cmd.Cert, cmd.Profile)
		if err != nil {
			return nil, err
		}
	}

	if !nmpkgs["github.com/fyne-io/mobile/app"] {
		return nil, fmt.Errorf(`%s does not import "github.com/fyne-io/mobile/app"`, pkg.PkgPath)
	}

	return pkg, nil
}

var nmRE = regexp.MustCompile(`[0-9a-f]{8} t _?(?:.*/vendor/)?(github.com/fyne-io.*/[^.]*)`)

func extractPkgs(nm string, path string) (map[string]bool, error) {
	if buildN {
		return map[string]bool{"github.com/fyne-io/mobile/app": true}, nil
	}
	r, w := io.Pipe()
	cmd := exec.Command(nm, path)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	nmpkgs := make(map[string]bool)
	errc := make(chan error, 1)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			if res := nmRE.FindStringSubmatch(s.Text()); res != nil {
				nmpkgs[res[1]] = true
			}
		}
		errc <- s.Err()
	}()

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s %s: %v", nm, path, err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("%s %s: %v", nm, path, err)
	}

	if err := <-errc; err != nil {
		return nil, fmt.Errorf("%s %s: %v", nm, path, err)
	}
	return nmpkgs, nil
}

var xout io.Writer = os.Stderr

func printcmd(format string, args ...interface{}) {
	cmd := fmt.Sprintf(format+"\n", args...)
	if tmpdir != "" {
		cmd = strings.Replace(cmd, tmpdir, "$WORK", -1)
	}
	if androidHome := os.Getenv("ANDROID_HOME"); androidHome != "" {
		cmd = strings.Replace(cmd, androidHome, "$ANDROID_HOME", -1)
	}
	if gomobilepath != "" {
		cmd = strings.Replace(cmd, gomobilepath, "$GOMOBILE", -1)
	}
	if gopath := goEnv("GOPATH"); gopath != "" {
		cmd = strings.Replace(cmd, gopath, "$GOPATH", -1)
	}
	if env := os.Getenv("HOMEPATH"); env != "" {
		cmd = strings.Replace(cmd, env, "$HOMEPATH", -1)
	}
	fmt.Fprint(xout, cmd)
}

// "Build flags", used by multiple commands.
var (
	buildA          bool        // -a
	buildI          bool        // -i
	buildN          bool        // -n
	buildV          bool        // -v
	buildX          bool        // -x
	buildO          string      // -o
	buildGcflags    string      // -gcflags
	buildLdflags    string      // -ldflags
	buildRelease    bool        // -release
	buildTarget     string      // -os
	buildTrimpath   bool        // -trimpath
	buildWork       bool        // -work
	buildBundleID   string      // -bundleid
	buildIOSVersion string      // -iosversion
	buildAndroidAPI int         // -androidapi
	buildTags       stringsFlag // -tags
)

// RunNewBuild executes a new mobile build for the specified configuration
func RunNewBuild(target, appID, icon, name, version string, build int, release bool, cert, profile string) error {
	buildTarget = target
	buildBundleID = appID
	buildRelease = release

	cmd := cmdBuild
	cmd.Flag = flag.FlagSet{}
	cmd.IconPath = icon
	cmd.AppName = name
	cmd.Version = version
	cmd.Build = build
	cmd.Cert = cert
	cmd.Profile = profile
	return runBuild(cmd)
}

func addBuildFlags(cmd *command) {
	cmd.Flag.StringVar(&buildO, "o", "", "")
	cmd.Flag.StringVar(&buildGcflags, "gcflags", "", "")
	cmd.Flag.StringVar(&buildLdflags, "ldflags", "", "")
	cmd.Flag.StringVar(&buildTarget, "os", "android", "")
	cmd.Flag.StringVar(&buildBundleID, "bundleid", "", "")
	cmd.Flag.StringVar(&buildIOSVersion, "iosversion", "7.0", "")
	cmd.Flag.IntVar(&buildAndroidAPI, "androidapi", minAndroidAPI, "")

	cmd.Flag.BoolVar(&buildA, "a", false, "")
	cmd.Flag.BoolVar(&buildI, "i", false, "")
	cmd.Flag.BoolVar(&buildTrimpath, "trimpath", false, "")
	cmd.Flag.Var(&buildTags, "tags", "")
}

func addBuildFlagsNVXWork(cmd *command) {
	cmd.Flag.BoolVar(&buildN, "n", false, "")
	cmd.Flag.BoolVar(&buildV, "v", false, "")
	cmd.Flag.BoolVar(&buildX, "x", false, "")
	cmd.Flag.BoolVar(&buildWork, "work", false, "")
}

func init() {
	addBuildFlags(cmdBuild)
	addBuildFlagsNVXWork(cmdBuild)

	addBuildFlagsNVXWork(cmdClean)
}

func goBuild(src string, env []string, args ...string) error {
	return goCmd("build", []string{src}, env, args...)
}

func goCmd(subcmd string, srcs []string, env []string, args ...string) error {
	return goCmdAt("", subcmd, srcs, env, args...)
}

func goCmdAt(at string, subcmd string, srcs []string, env []string, args ...string) error {
	cmd := exec.Command("go", subcmd)
	tags := buildTags
	targetOS, _, err := parseBuildTarget(buildTarget)
	if err != nil {
		return err
	}
	if targetOS == "darwin" {
		tags = append(tags, "ios")
	}
	if buildRelease {
		tags = append(tags, "release")
	}
	if len(tags) > 0 {
		cmd.Args = append(cmd.Args, "-tags", strings.Join(tags, " "))
	}
	if buildV {
		cmd.Args = append(cmd.Args, "-v")
	}
	if subcmd != "install" && buildI {
		cmd.Args = append(cmd.Args, "-i")
	}
	if buildX {
		cmd.Args = append(cmd.Args, "-x")
	}
	if buildGcflags != "" {
		cmd.Args = append(cmd.Args, "-gcflags", buildGcflags)
	}
	if buildLdflags != "" {
		cmd.Args = append(cmd.Args, "-ldflags", buildLdflags)
	}
	if buildTrimpath {
		cmd.Args = append(cmd.Args, "-trimpath")
	}
	if buildWork {
		cmd.Args = append(cmd.Args, "-work")
	}
	cmd.Args = append(cmd.Args, args...)
	cmd.Args = append(cmd.Args, srcs...)
	cmd.Env = append([]string{}, env...)
	cmd.Dir = at
	return runCmd(cmd)
}

func parseBuildTarget(buildTarget string) (os string, archs []string, _ error) {
	if buildTarget == "" {
		return "", nil, fmt.Errorf(`invalid target ""`)
	}

	all := false
	archNames := []string{}
	for i, p := range strings.Split(buildTarget, ",") {
		osarch := strings.SplitN(p, "/", 2) // len(osarch) > 0
		if osarch[0] != "android" && osarch[0] != "ios" {
			return "", nil, fmt.Errorf(`unsupported os`)
		}

		if i == 0 {
			os = osarch[0]
		}

		if os != osarch[0] {
			return "", nil, fmt.Errorf(`cannot target different OSes`)
		}

		if len(osarch) == 1 {
			all = true
		} else {
			archNames = append(archNames, osarch[1])
		}
	}

	// verify all archs are supported one while deduping.
	isSupported := func(os, arch string) bool {
		for _, a := range allArchs[os] {
			if a == arch {
				return true
			}
		}
		return false
	}

	seen := map[string]bool{}
	for _, arch := range archNames {
		if _, ok := seen[arch]; ok {
			continue
		}
		if !isSupported(os, arch) {
			return "", nil, fmt.Errorf(`unsupported arch: %q`, arch)
		}

		seen[arch] = true
		archs = append(archs, arch)
	}

	targetOS := os
	if os == "ios" {
		targetOS = "darwin"
	}
	if all {
		return targetOS, allArchs[os], nil
	}
	return targetOS, archs, nil
}
