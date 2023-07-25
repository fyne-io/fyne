package mobile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"

	"golang.org/x/mod/semver"
	"golang.org/x/sys/execabs"
)

// General mobile build environment. Initialized by envInit.
var (
	gomobilepath string // $GOPATH/pkg/gomobile

	androidEnv map[string][]string // android arch -> []string

	darwinEnv map[string][]string

	darwinArmNM string

	allArchs = map[string][]string{
		"android":      {"arm", "arm64", "386", "amd64"},
		"ios":          {"arm64"},
		"iossimulator": {"arm64", "amd64"},
	}

	bitcodeEnabled bool
)

func buildEnvInit() (cleanup func(), err error) {
	// Find gomobilepath.
	gopath := goEnv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		gomobilepath = filepath.Join(p, "pkg", "gomobile")
		if _, err := os.Stat(gomobilepath); buildN || err == nil {
			break
		}
	}

	if buildX {
		fmt.Fprintln(xout, "GOMOBILE="+gomobilepath)
	}

	// Check the toolchain is in a good state.
	// Pick a temporary directory for assembling an apk/app.
	if gomobilepath == "" {
		return nil, errors.New("toolchain not installed, run `gomobile init`")
	}

	cleanupFn := func() {
		if buildWork {
			fmt.Printf("WORK=%s\n", tmpdir)
			return
		}
		err := removeAll(tmpdir)
		if err != nil {
			fyne.LogError("Failed to remove all in cleanup function", err)
		}
	}
	if buildN {
		tmpdir = "$WORK"
		cleanupFn = func() {}
	} else {
		tmpdir, err = os.MkdirTemp("", "fyne-work-")
		if err != nil {
			return nil, err
		}
	}
	if buildX {
		fmt.Fprintln(xout, "WORK="+tmpdir)
	}

	if err := envInit(); err != nil {
		return nil, err
	}

	return cleanupFn, nil
}

var (
	before115 = false
	before116 = false
)

func envInit() (err error) {
	// Check the current Go version by go-list.
	// An arbitrary standard package ('runtime' here) is given to go-list.
	// This is because go-list tries to analyze the module at the current directory if no packages are given,
	// and if the module doesn't have any Go file, go-list fails. See golang/go#36668.

	ver, err := execabs.Command("go", "version").Output()
	if err == nil && string(ver) != "" {
		fields := strings.Split(string(ver), " ")
		if len(fields) >= 3 {
			goVer := strings.TrimPrefix(fields[2], "go")

			// If a go command is a development version, the version
			// information may only appears in the third elements.
			// For instance:
			// go version devel go1.18-527609d47b Wed Aug 25 17:07:58 2021 +0200 darwin/arm64
			if goVer == "devel" && len(fields) >= 4 {
				prefix := strings.Split(fields[3], "-")
				// a go command may miss version information. If that happens
				// we just use the environment version.
				if len(prefix) > 0 {
					goVer = strings.TrimPrefix(prefix[0], "go")
				} else {
					goVer = runtime.Version()
				}
			}

			before115 = semver.Compare("v"+goVer, "v1.15.0") < 0
			before116 = semver.Compare("v"+goVer, "v1.16.0") < 0
		}
	}
	if before115 {
		allArchs["ios"] = []string{"arm64", "amd64", "arm"}
	}

	// TODO re-enable once we find out what broke after September event 2020
	//cmd := exec.Command("go", "list", "-e", "-f", `{{range context.ReleaseTags}}{{if eq . "go1.14"}}{{.}}{{end}}{{end}}`, "runtime")
	//cmd.Stderr = os.Stderr
	//out, err := cmd.Output()
	//if err != nil {
	//	return err
	//}
	//if len(strings.TrimSpace(string(out))) > 0 {
	//		bitcodeEnabled = true
	//}

	// Setup the cross-compiler environments.
	if ndkRoot, err := ndkRoot(); err == nil {
		androidEnv = make(map[string][]string)
		if buildAndroidAPI < minAndroidAPI {
			return fmt.Errorf("gomobile requires Android API level >= %d", minAndroidAPI)
		}
		for arch, toolchain := range ndk {
			clang := toolchain.Path(ndkRoot, "clang")
			clangpp := toolchain.Path(ndkRoot, "clang++")
			if !buildN {
				tools := []string{clang, clangpp}
				if runtime.GOOS == "windows" {
					// Because of https://github.com/android-ndk/ndk/issues/920,
					// we require r19c, not just r19b. Fortunately, the clang++.cmd
					// script only exists in r19c.
					tools = append(tools, clangpp+".cmd")
				}
				for _, tool := range tools {
					_, err = os.Stat(tool)
					if err != nil {
						return fmt.Errorf("no compiler for %s was found in the NDK (tried %s). Make sure your NDK version is >= r19c. Use `sdkmanager --update` to update it", arch, tool)
					}
				}
			}
			androidEnv[arch] = []string{
				"GOOS=android",
				"GOARCH=" + arch,
				"CC=" + clang,
				"CXX=" + clangpp,
				"CGO_ENABLED=1",
			}
			if arch == "arm" {
				androidEnv[arch] = append(androidEnv[arch], "GOARM=7")
			}
		}
	}

	if !xcodeAvailable() || !util.IsIOS(buildTarget) {
		return nil
	}

	darwinArmNM = "nm"
	darwinEnv = make(map[string][]string)
	for _, arch := range allArchs[buildTarget] {
		var env []string
		var err error
		var clang, cflags string
		switch arch {
		case "arm":
			env = append(env, "GOARM=7")
			fallthrough
		case "arm64":
			if buildTarget == "ios" {
				clang, cflags, err = envClang("iphoneos")
				cflags += " -miphoneos-version-min=" + buildIOSVersion
			} else { // iossimulator
				clang, cflags, err = envClang("iphonesimulator")
				cflags += " -mios-simulator-version-min=" + buildIOSVersion
			}
		case "386", "amd64":
			clang, cflags, err = envClang("iphonesimulator")
			cflags += " -mios-simulator-version-min=" + buildIOSVersion
		default:
			panic(fmt.Errorf("unknown GOARCH: %q", arch))
		}
		if err != nil {
			return err
		}

		if bitcodeEnabled {
			cflags += " -fembed-bitcode"
		}

		os := "ios"
		if before116 {
			os = "darwin"
		}
		env = append(env,
			"GOOS="+os,
			"GOARCH="+arch,
			"CC="+clang,
			"CXX="+clang+"++",
			"CGO_CFLAGS="+cflags+" -arch "+archClang(arch),
			"CGO_CXXFLAGS="+cflags+" -arch "+archClang(arch),
			"CGO_LDFLAGS="+cflags+" -arch "+archClang(arch),
			"CGO_ENABLED=1",
		)
		darwinEnv[arch] = env
	}

	return nil
}

func ndkRoot() (string, error) {
	if buildN {
		return "$NDK_PATH", nil
	}

	ndkRoot := os.Getenv("ANDROID_NDK_HOME")
	if ndkRoot != "" {
		_, err := os.Stat(ndkRoot)
		if err == nil {
			return ndkRoot, nil
		}
	}

	androidHome := os.Getenv("ANDROID_HOME")
	if androidHome != "" {
		ndkRoot := filepath.Join(androidHome, "ndk-bundle")
		_, err := os.Stat(ndkRoot)
		if err == nil {
			return ndkRoot, nil
		}
	}

	return "", fmt.Errorf("no Android NDK found in $ANDROID_HOME/ndk-bundle nor in $ANDROID_NDK_HOME")
}

func envClang(sdkName string) (clang, cflags string, err error) {
	if buildN {
		return sdkName + "-clang", "-isysroot=" + sdkName, nil
	}
	cmd := execabs.Command("xcrun", "--sdk", sdkName, "--find", "clang")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("xcrun --find: %v\n%s", err, out)
	}
	clang = strings.TrimSpace(string(out))

	cmd = execabs.Command("xcrun", "--sdk", sdkName, "--show-sdk-path")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("xcrun --show-sdk-path: %v\n%s", err, out)
	}
	sdk := strings.TrimSpace(string(out))
	return clang, "-isysroot " + sdk, nil
}

func archClang(goarch string) string {
	switch goarch {
	case "arm":
		return "armv7"
	case "arm64":
		return "arm64"
	case "386":
		return "i386"
	case "amd64":
		return "x86_64"
	default:
		panic(fmt.Sprintf("unknown GOARCH: %q", goarch))
	}
}

// environ merges os.Environ and the given "key=value" pairs.
// If a key is in both os.Environ and kv, kv takes precedence.
func environ(kv []string) []string {
	cur := os.Environ()
	new := make([]string, 0, len(cur)+len(kv))

	envs := make(map[string]string, len(cur))
	for _, ev := range cur {
		elem := strings.SplitN(ev, "=", 2)
		if len(elem) != 2 || elem[0] == "" {
			// pass the env var of unusual form untouched.
			// e.g. Windows may have env var names starting with "=".
			new = append(new, ev)
			continue
		}
		if goos == "windows" {
			elem[0] = strings.ToUpper(elem[0])
		}
		envs[elem[0]] = elem[1]
	}
	for _, ev := range kv {
		elem := strings.SplitN(ev, "=", 2)
		if len(elem) != 2 || elem[0] == "" {
			panic(fmt.Sprintf("malformed env var %q from input", ev))
		}
		if goos == "windows" {
			elem[0] = strings.ToUpper(elem[0])
		}
		envs[elem[0]] = elem[1]
	}
	for k, v := range envs {
		new = append(new, k+"="+v)
	}
	return new
}

func archNDK() string {
	if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
		return "windows"
	}

	var arch string
	switch runtime.GOARCH {
	case "386":
		arch = "x86"
	case "amd64":
		arch = "x86_64"
	case "arm64":
		// For darwin/arm64, see https://golang.org/cl/346153.
		if runtime.GOOS == "darwin" {
			arch = "x86_64"
			break
		}
		if runtime.GOOS == "android" { // termux
			return "linux-aarch64"
		}
		fallthrough
	default:
		panic("unsupported GOARCH: " + runtime.GOARCH)
	}
	return runtime.GOOS + "-" + arch
}

type ndkToolchain struct {
	arch        string
	abi         string
	minAPI      int
	toolPrefix  string
	clangPrefix string
}

func (tc *ndkToolchain) ClangPrefix(api int) string {
	if api < tc.minAPI {
		return fmt.Sprintf("%s%d", tc.clangPrefix, tc.minAPI)
	}
	return fmt.Sprintf("%s%d", tc.clangPrefix, api)
}

func (tc *ndkToolchain) Path(ndkRoot, toolName string) string {
	for api := buildAndroidAPI; api < 99; api++ {
		var pref string
		switch toolName {
		case "clang", "clang++":
			pref = tc.ClangPrefix(api)
		case "nm":
			pref = "llvm"
		default:
			pref = tc.toolPrefix
		}

		toolPath := filepath.Join(ndkRoot, "toolchains", "llvm", "prebuilt", archNDK(), "bin", pref+"-"+toolName)
		if util.Exists(toolPath) {
			return toolPath
		} else if runtime.GOOS == "windows" {
			// On windows some of the NDK executable have a .exe extension and some don't, so try both.
			toolPath += ".exe"
			if util.Exists(toolPath) {
				return toolPath
			}
		}
	}
	return ""
}

type ndkConfig map[string]ndkToolchain // map: GOOS->androidConfig.

func (nc ndkConfig) Toolchain(arch string) ndkToolchain {
	tc, ok := nc[arch]
	if !ok {
		panic(`unsupported architecture: ` + arch)
	}
	return tc
}

var ndk = ndkConfig{
	"arm": {
		arch:        "arm",
		abi:         "armeabi-v7a",
		minAPI:      16,
		toolPrefix:  "arm-linux-androideabi",
		clangPrefix: "armv7a-linux-androideabi",
	},
	"arm64": {
		arch:        "arm64",
		abi:         "arm64-v8a",
		minAPI:      21,
		toolPrefix:  "aarch64-linux-android",
		clangPrefix: "aarch64-linux-android",
	},

	"386": {
		arch:        "x86",
		abi:         "x86",
		minAPI:      16,
		toolPrefix:  "i686-linux-android",
		clangPrefix: "i686-linux-android",
	},
	"amd64": {
		arch:        "x86_64",
		abi:         "x86_64",
		minAPI:      21,
		toolPrefix:  "x86_64-linux-android",
		clangPrefix: "x86_64-linux-android",
	},
}

func xcodeAvailable() bool {
	err := execabs.Command("xcrun", "xcodebuild", "-version").Run()
	return err == nil
}
