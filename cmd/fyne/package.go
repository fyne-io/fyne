package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/cmd/fyne/internal/mobile"
	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/jackmordaunt/icns"
	"github.com/josephspurrier/goversioninfo"
	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func ensureSubDir(parent, name string) string {
	path := filepath.Join(parent, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fyne.LogError("Failed to create dirrectory", err)
		}
	}
	return path
}

func copyFile(src, tgt string) error {
	return copyFileMode(src, tgt, 0644)
}

func copyExeFile(src, tgt string) error {
	return copyFileMode(src, tgt, 0755)
}

func copyFileMode(src, tgt string, perm os.FileMode) error {
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	source, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.OpenFile(tgt, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return err
}

func runAsAdminWindows(args ...string) error {
	cmd := "\"/c\""

	for _, arg := range args {
		cmd += ",\"" + arg + "\""
	}

	return exec.Command("powershell.exe", "Start-Process", "cmd.exe", "-Verb", "runAs", "-ArgumentList", cmd).Run()
}

// Declare conformity to command interface
var _ command = (*packager)(nil)

type packager struct {
	name, srcDir, dir, exe, icon string
	os, appID                    string
	install, release             bool
}

func (p *packager) packageLinux() error {
	var prefixDir string
	local := "local/"
	tempDir := "tmp-pkg"

	if p.install {
		tempDir = ""
	}

	if _, err := os.Stat(filepath.Join("/", "usr", "local")); os.IsNotExist(err) {
		prefixDir = ensureSubDir(ensureSubDir(p.dir, tempDir), "usr")
		local = ""
	} else {
		prefixDir = ensureSubDir(ensureSubDir(ensureSubDir(p.dir, tempDir), "usr"), "local")
	}

	shareDir := ensureSubDir(prefixDir, "share")

	binDir := ensureSubDir(prefixDir, "bin")
	binName := filepath.Join(binDir, filepath.Base(p.exe))
	err := copyExeFile(p.exe, binName)
	if err != nil {
		return errors.Wrap(err, "Failed to copy application binary file")
	}

	iconDir := ensureSubDir(shareDir, "pixmaps")
	iconPath := filepath.Join(iconDir, p.name+filepath.Ext(p.icon))
	err = copyFile(p.icon, iconPath)
	if err != nil {
		return errors.Wrap(err, "Failed to copy icon")
	}

	appsDir := ensureSubDir(shareDir, "applications")
	desktop := filepath.Join(appsDir, p.name+".desktop")
	deskFile, _ := os.Create(desktop)
	_, err = io.WriteString(deskFile, "[Desktop Entry]\n"+
		"Type=Application\n"+
		"Name="+p.name+"\n"+
		"Exec="+filepath.Base(p.exe)+"\n"+
		"Icon="+p.name+"\n")
	if err != nil {
		return errors.Wrap(err, "Failed to write desktop entry string")
	}

	if !p.install {
		defer os.RemoveAll(filepath.Join(p.dir, tempDir))

		makefile, _ := os.Create(filepath.Join(p.dir, tempDir, "Makefile"))
		_, err := io.WriteString(makefile, `ifneq ("$(wildcard /usr/local)","")`+"\n"+
			"PREFIX ?= /usr/local\n"+
			"else\n"+
			"PREFIX ?= /usr\n"+
			"endif\n"+
			"\n"+
			"default:\n"+
			`	# Run "sudo make install" to install the application.`+"\n"+
			"install:\n"+
			"	install -Dm00644 usr/"+local+"share/applications/"+p.name+`.desktop $(DESTDIR)$(PREFIX)/share/applications/`+p.name+".desktop\n"+
			"	install -Dm00755 usr/"+local+"bin/"+filepath.Base(p.exe)+` $(DESTDIR)$(PREFIX)/bin/`+filepath.Base(p.exe)+"\n"+
			"	install -Dm00644 usr/"+local+"share/pixmaps/"+p.name+filepath.Ext(p.icon)+` $(DESTDIR)$(PREFIX)/share/pixmaps/`+p.name+filepath.Ext(p.icon)+"\n")
		if err != nil {
			return errors.Wrap(err, "Failed to write Makefile string")
		}

		tarCmd := exec.Command("tar", "zcf", p.name+".tar.gz", "-C", tempDir, "usr", "Makefile")
		err = tarCmd.Run()
		if err != nil {
			return errors.Wrap(err, "Failed to create archive with tar")
		}
	}

	return nil
}

func (p *packager) packageDarwin() error {
	appDir := ensureSubDir(p.dir, p.name+".app")
	exeName := filepath.Base(p.exe)

	contentsDir := ensureSubDir(appDir, "Contents")
	info := filepath.Join(contentsDir, "Info.plist")
	infoFile, _ := os.Create(info)
	_, err := io.WriteString(infoFile, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n"+
		"<plist version=\"1.0\">\n"+
		"<dict>\n"+
		"<key>CFBundleName</key>\n"+
		"<string>"+p.name+"</string>\n"+
		"<key>CFBundleExecutable</key>\n"+
		"<string>"+exeName+"</string>\n"+
		"<key>CFBundleIdentifier</key>\n"+
		"<string>"+p.appID+"</string>\n"+
		"<key>CFBundleIconFile</key>\n"+
		"<string>icon.icns</string>\n"+
		"<key>CFBundleShortVersionString</key>\n"+
		"<string>1.0</string>\n"+
		"<key>NSHighResolutionCapable</key>\n"+
		"<true/>\n"+
		"<key>CFBundleInfoDictionaryVersion</key>\n"+
		"<string>6.0</string>\n"+
		"<key>CFBundlePackageType</key>\n"+
		"<string>APPL</string>\n"+
		"</dict>\n"+
		"</plist>\n")
	if err != nil {
		return errors.Wrap(err, "Failed to write infoFile string")
	}

	macOSDir := ensureSubDir(contentsDir, "MacOS")
	binName := filepath.Join(macOSDir, exeName)
	err = copyExeFile(p.exe, binName)
	if err != nil {
		return errors.Wrap(err, "Failed to copy exe file")
	}

	resDir := ensureSubDir(contentsDir, "Resources")
	icnsPath := filepath.Join(resDir, "icon.icns")

	img, err := os.Open(p.icon)
	if err != nil {
		return errors.Wrapf(err, "Failed to open source image \"%s\"", p.icon)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return errors.Wrapf(err, "Failed to decode source image")
	}
	dest, err := os.Create(icnsPath)
	if err != nil {
		return errors.Wrap(err, "Failed to open destination file")
	}
	defer dest.Close()
	if err := icns.Encode(dest, srcImg); err != nil {
		return errors.Wrap(err, "Failed to encode icns")
	}

	return nil
}

func (p *packager) packageWindows() error {
	exePath := filepath.Dir(p.exe)

	// convert icon
	img, err := os.Open(p.icon)
	if err != nil {
		return errors.Wrap(err, "Failed to open source image")
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return errors.Wrap(err, "Failed to decode source image")
	}

	icoPath := filepath.Join(exePath, p.name+".ico")
	file, err := os.Create(icoPath)
	if err != nil {
		return errors.Wrap(err, "Failed to open image file")
	}

	err = ico.Encode(file, srcImg)
	if err != nil {
		return errors.Wrap(err, "Failed to encode icon")
	}

	err = file.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close image file")
	}

	// write manifest
	manifest := p.exe + ".manifest"
	manifestGenerated := false
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		manifestGenerated = true
		manifestFile, _ := os.Create(manifest)
		_, err := io.WriteString(manifestFile, "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n"+
			"<assembly xmlns=\"urn:schemas-microsoft-com:asm.v1\" manifestVersion=\"1.0\" xmlns:asmv3=\"urn:schemas-microsoft-com:asm.v3\">\n"+
			"<assemblyIdentity version=\"1.0.0.0\" processorArchitecture=\"*\" name=\""+p.name+"\" type=\"win32\"/>\n"+
			"</assembly>")
		if err != nil {
			return errors.Wrap(err, "Failed to write manifest string")
		}
	}

	// launch rsrc to generate the object file
	outPath := filepath.Join(exePath, "fyne.syso")

	vi := &goversioninfo.VersionInfo{}
	vi.ProductName = p.name
	vi.IconPath = icoPath
	vi.ManifestPath = manifest

	vi.Build()
	vi.Walk()
	err = vi.WriteSyso(outPath, runtime.GOARCH)
	if err != nil {
		return errors.Wrap(err, "Failed to write .syso file")
	}

	err = os.Remove(icoPath)
	if err != nil {
		return errors.Wrap(err, "Failed to remove icon")
	} else if manifestGenerated {
		err := os.Remove(manifest)
		if err != nil {
			return errors.Wrap(err, "Failed to remove manifest")
		}
	}

	err = p.buildPackage()
	if err != nil {
		return errors.Wrap(err, "Failed to rebuild after adding metadata")
	}

	if p.install {
		err := runAsAdminWindows("copy", "\"\""+p.exe+"\"\"", "\"\""+filepath.Join(os.Getenv("ProgramFiles"), p.name)+"\"\"")
		if err != nil {
			return errors.Wrap(err, "Failed to run as administrator")
		}
	}
	return nil
}

func (p *packager) packageAndroid(arch string) error {
	return mobile.RunNewBuild(arch, p.appID, p.icon, p.name, p.release)
}

func (p *packager) packageIOS() error {
	err := mobile.RunNewBuild("ios", p.appID, p.icon, p.name, p.release)
	if err != nil {
		return err
	}

	appDir := filepath.Join(p.dir, p.name+".app")
	iconPath := filepath.Join(appDir, "Icon.png")
	return copyFile(p.icon, iconPath)
}

func (p *packager) addFlags() {
	flag.StringVar(&p.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, ios, linux, windows)")
	flag.StringVar(&p.exe, "executable", "", "The path to the executable, default is the current dir main binary")
	flag.StringVar(&p.srcDir, "sourceDir", "", "The directory to package, if executable is not set")
	flag.StringVar(&p.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&p.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&p.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.BoolVar(&p.release, "release", false, "Should this package be prepared for release? (disable debug etc)")
}

func (*packager) printHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for distribution.")
	fmt.Println(indent, "You may specify the -executable to package, otherwise -sourceDir will be built.")
	fmt.Println(indent, "Command usage: fyne package [parameters]")
}

func (p *packager) buildPackage() error {
	b := &builder{
		os:     p.os,
		srcdir: p.srcDir,
	}

	return b.build()
}

func (p *packager) run(_ []string) {
	err := p.validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = p.doPackage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func (p *packager) validate() error {
	if p.os == "" {
		p.os = runtime.GOOS
	}
	baseDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Unable to get the current directory, needed to find main executable")
	}
	if p.dir == "" {
		p.dir = baseDir
	}
	if p.srcDir == "" {
		p.srcDir = baseDir
	} else if p.os == "ios" || p.os == "android" {
		return errors.New("Parameter -sourceDir is currently not supported for mobile builds.\n" +
			"Change directory to the main package and try again.")
	}

	exeName := calculateExeName(p.srcDir, p.os)

	if p.exe == "" {
		p.exe = filepath.Join(p.srcDir, exeName)
	} else if p.os == "ios" || p.os == "android" {
		_, _ = fmt.Fprint(os.Stderr, "Parameter -executable is ignored for mobile builds.\n")
	}

	if p.name == "" {
		p.name = exeName
	}
	if p.icon == "" || p.icon == "Icon.png" {
		p.icon = filepath.Join(p.srcDir, "Icon.png")
	}
	if !exists(p.icon) {
		return errors.New("Missing application icon at \"" + p.icon + "\"")
	}

	// only used for iOS and macOS
	if p.appID == "" {
		if p.os == "darwin" {
			p.appID = "com.example." + p.name
		} else if p.os == "ios" || p.os == "android" {
			return errors.New("Missing appID parameter for mobile package")
		}
	}

	return nil
}

func calculateExeName(sourceDir, os string) string {
	exeName := filepath.Base(sourceDir)
	/* #nosec */
	if data, err := ioutil.ReadFile(filepath.Join(sourceDir, "go.mod")); err == nil {
		modulePath := modfile.ModulePath(data)
		moduleName, _, ok := module.SplitPathVersion(modulePath)
		if ok {
			paths := strings.Split(moduleName, "/")
			name := paths[len(paths)-1]
			if name != "" {
				exeName = name
			}
		}
	}

	if os == "windows" {
		exeName = exeName + ".exe"
	}

	return exeName
}

func (p *packager) doPackage() error {
	if !exists(p.exe) && !p.isMobile() {
		err := p.buildPackage()
		if err != nil {
			return errors.Wrap(err, "Error building application")
		}
		if !exists(p.exe) {
			return fmt.Errorf("unable to build directory to expected executable, %s", p.exe)
		}
	}

	switch p.os {
	case "darwin":
		return p.packageDarwin()
	case "linux":
		return p.packageLinux()
	case "windows":
		return p.packageWindows()
	case "android/arm", "android/arm64", "android/amd64", "android/386", "android":
		return p.packageAndroid(p.os)
	case "ios":
		return p.packageIOS()
	default:
		return fmt.Errorf("unsupported target operating system \"%s\"", p.os)
	}
}

func (p *packager) isMobile() bool {
	return p.os == "ios" || strings.HasPrefix(p.os, "android")
}
