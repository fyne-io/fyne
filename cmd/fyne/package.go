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

	"fyne.io/fyne"
	"github.com/Kodeworks/golang-image-ico"
	"github.com/jackmordaunt/icns"
	"github.com/josephspurrier/goversioninfo"
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

type packager struct {
	os, name, dir, exe, icon string
}

func (p *packager) packageLinux() {
	prefixDir := ensureSubDir(ensureSubDir(p.dir, "usr"), "local")
	shareDir := ensureSubDir(prefixDir, "share")

	appsDir := ensureSubDir(shareDir, "applications")
	desktop := filepath.Join(appsDir, p.name+".desktop")
	deskFile, _ := os.Create(desktop)
	io.WriteString(deskFile, "[Desktop Entry]\n"+
		"Type=Application\n"+
		"Name="+p.name+"\n"+
		"Exec="+filepath.Base(p.exe)+"\n"+
		"Icon="+filepath.Base(p.icon)+"\n")

	binDir := ensureSubDir(prefixDir, "bin")
	binName := filepath.Join(binDir, filepath.Base(p.exe))
	copyExeFile(p.exe, binName)

	iconThemeDir := ensureSubDir(ensureSubDir(shareDir, "icons"), "hicolor")
	iconDir := ensureSubDir(ensureSubDir(iconThemeDir, "512x512"), "apps")
	iconPath := filepath.Join(iconDir, filepath.Base(p.icon))
	copyFile(p.icon, iconPath)

	tarCmd := exec.Command("tar", "cf", p.name+".tar.gz", "usr")
	tarCmd.Run()
}

func (p *packager) packageDarwin() {
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
		fyne.LogError("Failed to open source image", err)
		panic(err)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		fyne.LogError("Faied to decode source image", err)
		panic(err)
	}
	dest, err := os.Create(icnsPath)
	if err != nil {
		fyne.LogError("Failed to open destination file", err)
		panic(err)
	}
	defer dest.Close()
	if err := icns.Encode(dest, srcImg); err != nil {
		fyne.LogError("Failed to encode icns", err)
		panic(err)
	}
}

func (p *packager) packageWindows() {
	// convert icon
	img, err := os.Open(p.icon)
	if err != nil {
		fyne.LogError("Failed to open source image", err)
		panic(err)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		fyne.LogError("Failed to decode source image", err)
		panic(err)
	}

	icoPath := filepath.Join(p.dir, p.name+".ico")
	file, err := os.Create(icoPath)
	if err != nil {
		fyne.LogError("Failed to open image file", err)
		panic(err)
	}
	ico.Encode(file, srcImg)
	file.Close()

	// write manifest
	manifest := p.exe + ".manifest"
	manifestFile, _ := os.Create(manifest)
	io.WriteString(manifestFile, "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n"+
		"<assembly xmlns=\"urn:schemas-microsoft-com:asm.v1\" manifestVersion=\"1.0\" xmlns:asmv3=\"urn:schemas-microsoft-com:asm.v3\">\n"+
		"<assemblyIdentity version=\"1.0.0.0\" processorArchitecture=\"*\" name=\""+p.name+"\" type=\"win32\"/>\n"+
		"</assembly>")

	// launch rsrc to generate the object file
	outPath := filepath.Join(p.dir, "fyne.syso")

	vi := &goversioninfo.VersionInfo{}
	vi.ProductName = p.name
	vi.IconPath = icoPath
	vi.ManifestPath = manifest

	vi.Build()
	vi.Walk()
	vi.WriteSyso(outPath, runtime.GOARCH)

	os.Remove(icoPath)
	os.Remove(manifest)
	fmt.Println("For windows releases you need to re-run go build to include the changes")
}

func (p *packager) addFlags() {
	flag.StringVar(&p.os, "os", runtime.GOOS, "The operating system to target (like GOOS)")
	flag.StringVar(&p.exe, "executable", "", "The path to the executable, default is the current dir main binary")
	flag.StringVar(&p.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&p.icon, "icon", "", "The name of the application icon file")
}

func (*packager) printHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for distribution.")
	fmt.Println(indent, "Before running this command you must to build the application for release.")
	fmt.Println(indent, "Command usage: fyne package [parameters]")
}

func (p *packager) run(_ []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get the current directory, needed to find main executable")
	}
	p.dir = dir

	if p.exe == "" {
		exeName := filepath.Base(dir)
		if p.os == "windows" {
			exeName = exeName + ".exe"
		}

		p.exe = filepath.Join(dir, exeName)
	}

	if !exists(p.exe) {
		fmt.Fprintln(os.Stderr, "Executable", p.exe, "could not be found")
		fmt.Fprintln(os.Stderr, "You may not have run \"go build\" for the target os")
		os.Exit(1)
	}
	if p.name == "" {
		name := filepath.Base(p.exe)
		p.name = strings.TrimSuffix(name, filepath.Ext(name))
	}
	if p.icon == "" {
		fmt.Fprintln(os.Stderr, "The -icon parameter is required for packaging")
		os.Exit(1)
	}

	switch p.os {
	case "linux":
		p.packageLinux()
	case "darwin":
		p.packageDarwin()
	case "windows":
		p.packageWindows()
	default:
		fmt.Fprintln(os.Stderr, "Unsupported os parameter,", p.os)
	}
}
