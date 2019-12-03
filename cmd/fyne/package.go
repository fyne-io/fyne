package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/cmd/fyne/internal/mobile"
	"github.com/Kodeworks/golang-image-ico"
	"github.com/jackmordaunt/icns"
	"github.com/josephspurrier/goversioninfo"
	"github.com/pkg/errors"
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
		os.Mkdir(path, os.ModePerm)
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

	source, err := os.Open(src)
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
	install                      bool
}

func (p *packager) packageLinux() error {
	prefixDir := ensureSubDir(ensureSubDir(p.dir, "usr"), "local")
	shareDir := ensureSubDir(prefixDir, "share")

	appsDir := ensureSubDir(shareDir, "applications")
	desktop := filepath.Join(appsDir, p.name+".desktop")
	deskFile, _ := os.Create(desktop)
	io.WriteString(deskFile, "[Desktop Entry]\n"+
		"Type=Application\n"+
		"Name="+p.name+"\n"+
		"Exec="+filepath.Base(p.exe)+"\n"+
		"Icon="+p.name+"\n")

	binDir := ensureSubDir(prefixDir, "bin")
	binName := filepath.Join(binDir, filepath.Base(p.exe))
	copyExeFile(p.exe, binName)

	iconThemeDir := ensureSubDir(ensureSubDir(shareDir, "icons"), "hicolor")
	iconDir := ensureSubDir(ensureSubDir(iconThemeDir, "512x512"), "apps")
	iconPath := filepath.Join(iconDir, p.name+filepath.Ext(p.icon))
	copyFile(p.icon, iconPath)

	if !p.install {
		tarCmd := exec.Command("tar", "cf", p.name+".tar.gz", "usr")
		tarCmd.Run()
	}

	return nil
}

func (p *packager) packageDarwin() error {
	appDir := ensureSubDir(p.dir, p.name+".app")
	exeName := filepath.Base(p.exe)

	contentsDir := ensureSubDir(appDir, "Contents")
	info := filepath.Join(contentsDir, "Info.plist")
	infoFile, _ := os.Create(info)
	io.WriteString(infoFile, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
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

	macOSDir := ensureSubDir(contentsDir, "MacOS")
	binName := filepath.Join(macOSDir, exeName)
	copyExeFile(p.exe, binName)

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
	ico.Encode(file, srcImg)
	file.Close()

	// write manifest
	manifest := p.exe + ".manifest"
	manifestGenerated := false
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		manifestGenerated = true
		manifestFile, _ := os.Create(manifest)
		io.WriteString(manifestFile, "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n"+
			"<assembly xmlns=\"urn:schemas-microsoft-com:asm.v1\" manifestVersion=\"1.0\" xmlns:asmv3=\"urn:schemas-microsoft-com:asm.v3\">\n"+
			"<assemblyIdentity version=\"1.0.0.0\" processorArchitecture=\"*\" name=\""+p.name+"\" type=\"win32\"/>\n"+
			"</assembly>")
	}

	// launch rsrc to generate the object file
	outPath := filepath.Join(exePath, "fyne.syso")

	vi := &goversioninfo.VersionInfo{}
	vi.ProductName = p.name
	vi.IconPath = icoPath
	vi.ManifestPath = manifest

	vi.Build()
	vi.Walk()
	vi.WriteSyso(outPath, runtime.GOARCH)

	os.Remove(icoPath)
	if manifestGenerated {
		os.Remove(manifest)
	}

	err = p.buildPackage()
	if err != nil {
		return errors.Wrap(err, "Failed to rebuild after adding metadata")
	}

	if p.install {
		runAsAdminWindows("copy", "\"\""+p.exe+"\"\"", "\"\""+filepath.Join(os.Getenv("ProgramFiles"), p.name)+"\"\"")
	}
	return nil
}

func (p *packager) packageAndroid() error {
	return mobile.RunNewBuild("android", p.appID, p.icon)
}

func (p *packager) packageIOS() error {
	err := mobile.RunNewBuild("ios", p.appID, p.icon)
	if err != nil {
		return err
	}

	appDir := filepath.Join(p.dir, p.name+".app")
	iconPath := filepath.Join(appDir, "Icon.png")
	return copyFile(p.icon, iconPath)
}

func (p *packager) addFlags() {
	flag.StringVar(&p.os, "os", "", "The operating system to target (android, darwin, ios, linux, windows)")
	flag.StringVar(&p.exe, "executable", "", "The path to the executable, default is the current dir main binary")
	flag.StringVar(&p.srcDir, "sourceDir", "", "The directory to package, if executable is not set")
	flag.StringVar(&p.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&p.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&p.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
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
	}

	exeName := filepath.Base(p.srcDir)
	if p.os == "windows" {
		exeName = exeName + ".exe"
	}
	if p.exe == "" {
		p.exe = filepath.Join(p.srcDir, exeName)
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
		} else if p.os == "ios" {
			return errors.New("Missing appID parameter for ios package")
		}
	}

	return nil
}

func (p *packager) doPackage() error {
	if !exists(p.exe) {
		err := p.buildPackage()
		if err != nil {
			return errors.Wrap(err, "Error building application")
		}
		if !exists(p.exe) {
			return fmt.Errorf("Unable to build directory to expected executable, %s" + p.exe)
		}
	}

	switch p.os {
	case "darwin":
		return p.packageDarwin()
	case "linux":
		return p.packageLinux()
	case "windows":
		return p.packageWindows()
	case "android":
		return p.packageAndroid()
	case "ios":
		return p.packageIOS()
	default:
		return errors.New("Unsupported terget operating system \"" + p.os + "\"")
	}
}
