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
	"strings"

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
	os, name, dir, exe, icon string
	install                  bool
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
		"<string>com.example."+p.name+"</string>\n"+
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

	if p.install {
		runAsAdminWindows("copy", "\"\""+p.exe+"\"\"", "\"\""+filepath.Join(os.Getenv("ProgramFiles"), p.name)+"\"\"")
	}
	return nil
}

func (p *packager) addFlags() {
	flag.StringVar(&p.os, "os", runtime.GOOS, "The operating system to target (like GOOS)")
	flag.StringVar(&p.exe, "executable", "", "The path to the executable, default is the current dir main binary")
	flag.StringVar(&p.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&p.icon, "icon", "Icon.png", "The name of the application icon file")
}

func (*packager) printHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for distribution.")
	fmt.Println(indent, "Before running this command you must to build the application for release.")
	fmt.Println(indent, "Command usage: fyne package [parameters]")
}

func (p *packager) run(_ []string) {
	err := p.validate()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = p.doPackage()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	if runtime.GOOS == "windows" {
		fmt.Println("For windows releases you need to re-run go build to include the changes")
	}
}

func (p *packager) validate() error {
	baseDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Unable to get the current directory, needed to find main executable")
	}
	if p.dir == "" {
		p.dir = baseDir
	}

	if p.exe == "" {
		exeName := filepath.Base(baseDir)
		if p.os == "windows" {
			exeName = exeName + ".exe"
		}

		p.exe = filepath.Join(baseDir, exeName)
	}
	if !exists(p.exe) {
		return fmt.Errorf("Executable %s could not be found\nYou may not have run \"go build\" for the target os", p.exe)
	}

	if p.name == "" {
		name := filepath.Base(p.exe)
		p.name = strings.TrimSuffix(name, filepath.Ext(name))
	}
	if p.icon == "" {
		p.icon = filepath.Join(baseDir, "Icon.png")
	}

	return nil
}

func (p *packager) doPackage() error {
	switch p.os {
	case "darwin":
		return p.packageDarwin()
	case "linux":
		return p.packageLinux()
	case "windows":
		return p.packageWindows()
	default:
		return errors.New("Unsupported terget operating system \"" + p.os + "\"")
	}
}
